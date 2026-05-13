import base64
import json
import time
from datetime import datetime, timedelta, timezone

import pytest
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import padding, rsa

from lykn.exceptions import (
    FeatureNotLicensedError,
    HardwareMismatchError,
    LicenseExpiredError,
    LicenseNotYetValidError,
)
from lykn.schemas import Hardware
from lykn.validator import LicenseValidator


def build_signed_license(*, private_key, payload_overrides: dict | None = None) -> str:
    now = datetime.now(timezone.utc)
    payload = {
        "id": "lic-001",
        "version": 1,
        "subject": {"name": "Demo", "email": "demo@example.com", "organization": "Acme"},
        "plan": "pro",
        "plan_name": "Pro Plan",
        "issued_at": now.isoformat(),
        "not_before": (now - timedelta(minutes=5)).isoformat(),
        "not_after": (now + timedelta(days=30)).isoformat(),
        "features": ["reports", "exports"],
        "limits": {"max_users": 20, "max_devices": 2},
        "metadata": {"env": "test"},
    }
    if payload_overrides:
        payload.update(payload_overrides)

    payload_bytes = json.dumps(payload).encode()
    signature = private_key.sign(payload_bytes, padding.PKCS1v15(), hashes.SHA256())
    return json.dumps(
        {
            "payload": base64.b64encode(payload_bytes).decode(),
            "signature": base64.b64encode(signature).decode(),
        }
    )


@pytest.fixture
def rsa_materials():
    private_key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    public_pem = private_key.public_key().public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return private_key, public_pem.decode()


def test_verify_returns_license_data(rsa_materials):
    private_key, public_pem = rsa_materials
    lic_content = build_signed_license(private_key=private_key)

    validator = LicenseValidator(public_key=public_pem, license_content=lic_content)
    result = validator.verify(required_features=["reports"])

    assert result.id == "lic-001"
    assert result.plan == "pro"
    assert result.plan_name == "Pro Plan"
    assert result.limits.max_users == 20
    assert result.limits.max_devices == 2
    assert result.metadata == {"env": "test"}
    assert validator.has_feature("exports") is True


def test_verify_accepts_backend_style_hardware_payload(rsa_materials, monkeypatch):
    private_key, public_pem = rsa_materials
    lic_content = build_signed_license(
        private_key=private_key,
        payload_overrides={
            "hardware": {
                "hostname": "prod-host",
                "cpu_id": "CPU-1",
                "disk_serial": "DISK-1",
                "mac_addresses": ["AA:BB:CC:DD:EE:FF"],
                "ip_addresses": ["10.0.0.5"],
            },
        },
    )

    monkeypatch.setattr(
        "lykn.validator.collect_hardware",
        lambda: Hardware(
            hostname="prod-host",
            cpu_id="CPU-1",
            disk_serial="DISK-1",
            mac_addresses=["aa-bb-cc-dd-ee-ff"],
            ip_addresses=["10.0.0.5"],
        ),
    )

    validator = LicenseValidator(public_key=public_pem, license_content=lic_content)
    result = validator.verify(required_features=["reports", "exports"])

    assert result.hardware is not None
    assert result.hardware.hostname == "prod-host"
    assert result.hardware.cpu_id == "CPU-1"
    assert result.hardware.disk_serial == "DISK-1"
    assert result.hardware.mac_addresses == ["AA:BB:CC:DD:EE:FF"]
    assert result.hardware.ip_addresses == ["10.0.0.5"]


def test_verify_raises_when_license_expired(rsa_materials):
    private_key, public_pem = rsa_materials
    now = datetime.now(timezone.utc)
    lic_content = build_signed_license(
        private_key=private_key,
        payload_overrides={
            "not_before": (now - timedelta(days=10)).isoformat(),
            "not_after": (now - timedelta(days=1)).isoformat(),
        },
    )

    validator = LicenseValidator(public_key=public_pem, license_content=lic_content)

    with pytest.raises(LicenseExpiredError):
        validator.verify()


def test_verify_raises_when_license_not_yet_valid(rsa_materials):
    private_key, public_pem = rsa_materials
    now = datetime.now(timezone.utc)
    lic_content = build_signed_license(
        private_key=private_key,
        payload_overrides={
            "not_before": (now + timedelta(days=1)).isoformat(),
            "not_after": (now + timedelta(days=30)).isoformat(),
        },
    )

    validator = LicenseValidator(public_key=public_pem, license_content=lic_content)

    with pytest.raises(LicenseNotYetValidError):
        validator.verify()


def test_verify_raises_when_hardware_mismatched(rsa_materials, monkeypatch):
    private_key, public_pem = rsa_materials
    lic_content = build_signed_license(
        private_key=private_key,
        payload_overrides={
            "hardware": {"hostname": "prod-host", "mac_addresses": ["AA:BB:CC:DD:EE:FF"]},
        },
    )

    monkeypatch.setattr(
        "lykn.validator.collect_hardware",
        lambda: Hardware(hostname="dev-host", mac_addresses=["11:22:33:44:55:66"]),
    )

    validator = LicenseValidator(public_key=public_pem, license_content=lic_content)

    with pytest.raises(HardwareMismatchError):
        validator.verify()


def test_verify_raises_when_feature_missing(rsa_materials):
    private_key, public_pem = rsa_materials
    lic_content = build_signed_license(private_key=private_key)

    validator = LicenseValidator(public_key=public_pem, license_content=lic_content)

    with pytest.raises(FeatureNotLicensedError):
        validator.verify(required_features=["billing"])


def test_start_without_background_thread_returns_license(rsa_materials):
    private_key, public_pem = rsa_materials
    lic_content = build_signed_license(private_key=private_key)

    validator = LicenseValidator(public_key=public_pem, license_content=lic_content, check_interval=0)
    result = validator.start()

    assert result.id == "lic-001"
    assert validator._thread is None


def test_background_failure_marks_invalid_and_calls_callback(rsa_materials, monkeypatch):
    private_key, public_pem = rsa_materials
    lic_content = build_signed_license(private_key=private_key)
    validator = LicenseValidator(public_key=public_pem, license_content=lic_content, check_interval=1)

    calls: list[str] = []
    original = validator._verify_once
    state = {"count": 0}

    def flaky_verify(required_features=None):
        state["count"] += 1
        if state["count"] == 1:
            return original(required_features)
        raise LicenseExpiredError("license expired in background")

    monkeypatch.setattr(validator, "_verify_once", flaky_verify)

    @validator.on_invalid
    def handle_invalid(reason: str):
        calls.append(reason)

    validator.start()
    time.sleep(1.2)
    validator.stop()

    assert validator._invalid is True
    assert calls and "license expired in background" in calls[0]

    with pytest.raises(LicenseExpiredError):
        validator.verify()
