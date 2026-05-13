from datetime import datetime

from lykn.schemas import Hardware, LicenseData, LicenseLimits, Subject


def test_license_data_minimal():
    data = LicenseData(
        id="test-uuid",
        subject=Subject(name="Test"),
        issued_at=datetime(2026, 1, 1),
        not_before=datetime(2026, 1, 1),
        not_after=datetime(2027, 1, 1),
    )
    assert data.id == "test-uuid"
    assert data.version == 1
    assert data.features == []
    assert data.hardware is None
    assert data.plan == ""
    assert data.plan_name == ""
    assert data.limits.max_users == 0
    assert data.limits.max_devices == 0
    assert data.metadata == {}


def test_license_data_full():
    data = LicenseData(
        id="test-uuid",
        version=1,
        subject=Subject(name="Company", email="a@b.com", organization="Org"),
        plan="enterprise",
        plan_name="Enterprise Plan",
        issued_at=datetime(2026, 1, 1),
        not_before=datetime(2026, 1, 1),
        not_after=datetime(2027, 1, 1),
        hardware=Hardware(hostname="server-01", mac_addresses=["AA:BB:CC:DD:EE:FF"]),
        features=["feature_a", "feature_b"],
        limits=LicenseLimits(max_users=20, max_devices=2),
        metadata={"key": "value"},
    )
    assert data.plan == "enterprise"
    assert data.plan_name == "Enterprise Plan"
    assert data.hardware.hostname == "server-01"
    assert len(data.hardware.mac_addresses) == 1
    assert len(data.features) == 2
    assert data.limits.max_users == 20
    assert data.limits.max_devices == 2
    assert data.metadata["key"] == "value"


def test_license_data_from_dict():
    raw = {
        "id": "test-uuid",
        "version": 1,
        "subject": {"name": "Test", "email": "", "organization": ""},
        "plan": "",
        "plan_name": "Default Plan",
        "issued_at": "2026-01-01T00:00:00",
        "not_before": "2026-01-01T00:00:00",
        "not_after": "2027-01-01T00:00:00",
        "features": [],
        "limits": {"max_users": 0, "max_devices": 1},
        "metadata": {},
    }
    data = LicenseData(**raw)
    assert data.id == "test-uuid"
    assert data.subject.name == "Test"
    assert data.plan_name == "Default Plan"
    assert data.limits.max_devices == 1


def test_hardware_defaults():
    hw = Hardware()
    assert hw.mac_addresses == []
    assert hw.ip_addresses == []
    assert hw.hostname == ""
    assert hw.cpu_id == ""
    assert hw.disk_serial == ""
