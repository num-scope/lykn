import base64
import json

import pytest
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import padding, rsa

from lykn.crypto import load_public_key, verify_license_file
from lykn.exceptions import LicenseFileError, LicenseSignatureError


@pytest.fixture
def rsa_key_pair():
    """Generate an RSA key pair for testing."""
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
    )
    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption(),
    )
    public_pem = private_key.public_key().public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return private_key, private_pem, public_pem


@pytest.fixture
def signed_license(rsa_key_pair):
    """Create a signed .lic content using Python (same protocol as Go)."""
    private_key, _, public_pem = rsa_key_pair

    license_data = {
        "id": "test-uuid",
        "version": 1,
        "subject": {"name": "Test", "email": "test@example.com", "organization": ""},
        "plan": "enterprise",
        "issued_at": "2026-01-01T00:00:00Z",
        "not_before": "2026-01-01T00:00:00Z",
        "not_after": "2027-01-01T00:00:00Z",
        "features": ["feature_a"],
        "metadata": {},
    }

    payload = json.dumps(license_data).encode()
    signature = private_key.sign(
        payload,
        padding.PKCS1v15(),
        hashes.SHA256(),
    )

    lic_content = json.dumps(
        {
            "payload": base64.b64encode(payload).decode(),
            "signature": base64.b64encode(signature).decode(),
        }
    )

    return lic_content, public_pem, license_data


def test_load_public_key_from_pem_string(rsa_key_pair):
    _, _, public_pem = rsa_key_pair
    key = load_public_key(public_pem.decode())
    assert key is not None


def test_load_public_key_from_file(rsa_key_pair, tmp_path):
    _, _, public_pem = rsa_key_pair
    key_file = tmp_path / "public.pem"
    key_file.write_bytes(public_pem)
    key = load_public_key(str(key_file))
    assert key is not None


def test_load_public_key_from_path_object(rsa_key_pair, tmp_path):
    _, _, public_pem = rsa_key_pair
    key_file = tmp_path / "public.pem"
    key_file.write_bytes(public_pem)
    key = load_public_key(key_file)
    assert key is not None


def test_verify_license_file_valid(signed_license):
    lic_content, public_pem, expected_data = signed_license
    key = load_public_key(public_pem.decode())
    result = verify_license_file(lic_content, key)
    assert result["id"] == expected_data["id"]
    assert result["plan"] == expected_data["plan"]
    assert result["features"] == expected_data["features"]


def test_verify_license_file_bytes_input(signed_license):
    lic_content, public_pem, expected_data = signed_license
    key = load_public_key(public_pem.decode())
    result = verify_license_file(lic_content.encode(), key)
    assert result["id"] == expected_data["id"]


def test_verify_license_file_invalid_signature(rsa_key_pair):
    _, _, public_pem = rsa_key_pair
    invalid_lic = json.dumps(
        {
            "payload": base64.b64encode(b'{"id":"test"}').decode(),
            "signature": base64.b64encode(b"invalid").decode(),
        }
    )
    key = load_public_key(public_pem.decode())
    with pytest.raises(LicenseSignatureError):
        verify_license_file(invalid_lic, key)


def test_verify_license_file_invalid_json():
    with pytest.raises(LicenseFileError):
        verify_license_file("not json", None)


def test_verify_license_file_missing_fields():
    with pytest.raises(LicenseFileError):
        verify_license_file('{"payload":"abc"}', None)
