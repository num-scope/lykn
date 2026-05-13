import json
from pathlib import Path

from lykn.crypto import load_public_key, verify_license_file
from lykn.schemas import LicenseData

FIXTURES_DIR = Path(__file__).resolve().parent.parent.parent.parent / "tests" / "fixtures"


def test_verify_go_signed_license():
    """Verify that a license signed by the Go crypto module can be validated by the Python SDK."""
    public_pem = (FIXTURES_DIR / "public.pem").read_text()
    lic_content = (FIXTURES_DIR / "license.lic").read_text()
    expected = json.loads((FIXTURES_DIR / "license.json").read_text())

    key = load_public_key(public_pem)
    result = verify_license_file(lic_content, key)

    assert result["id"] == expected["id"]
    assert result["version"] == expected["version"]
    assert result["subject"]["name"] == expected["subject"]["name"]
    assert result["subject"]["email"] == expected["subject"]["email"]
    assert result["plan"] == expected["plan"]
    assert result["plan_name"] == expected["plan_name"]
    assert result["features"] == expected["features"]
    assert result["hardware"] == expected["hardware"]
    assert result["limits"] == expected["limits"]
    assert result["metadata"] == expected["metadata"]
    assert result["not_before"] == expected["not_before"]
    assert result["not_after"] == expected["not_after"]

    data = LicenseData.model_validate(result)
    assert data.id == expected["id"]
    assert data.subject.name == expected["subject"]["name"]
    assert data.plan == expected["plan"]
    assert data.plan_name == expected["plan_name"]
    assert data.features == expected["features"]
    assert data.hardware is not None
    assert data.hardware.hostname == expected["hardware"]["hostname"]
    assert data.hardware.cpu_id == expected["hardware"]["cpu_id"]
    assert data.hardware.disk_serial == expected["hardware"]["disk_serial"]
    assert data.hardware.mac_addresses == expected["hardware"]["mac_addresses"]
    assert data.limits.max_users == expected["limits"]["max_users"]
    assert data.limits.max_devices == expected["limits"]["max_devices"]
    assert data.metadata == expected["metadata"]
