import base64
import json

from click.testing import CliRunner

from lykn.cli.main import cli


class FakeLicense:
    id = "lic-001"
    plan = "pro"
    features = ["reports", "exports"]


class FakeValidator:
    def __init__(self, public_key, license_path=None, license_content=None, check_interval=0):
        self.public_key = public_key
        self.license_path = license_path
        self.license_content = license_content
        self.check_interval = check_interval

    def verify(self, required_features=None):
        if required_features == ["billing"]:
            raise RuntimeError("Missing licensed features: billing")
        return FakeLicense()


def test_verify_command_success(monkeypatch, tmp_path):
    monkeypatch.setattr("lykn.cli.commands.verify.LicenseValidator", FakeValidator)
    key_file = tmp_path / "public.pem"
    lic_file = tmp_path / "license.lic"
    key_file.write_text("public-key")
    lic_file.write_text("license-content")

    result = CliRunner().invoke(
        cli,
        ["verify", "--public-key", str(key_file), "--license", str(lic_file), "--feature", "reports"],
    )

    assert result.exit_code == 0
    assert "License valid" in result.output
    assert "lic-001" in result.output


def test_verify_command_failure(monkeypatch, tmp_path):
    monkeypatch.setattr("lykn.cli.commands.verify.LicenseValidator", FakeValidator)
    key_file = tmp_path / "public.pem"
    lic_file = tmp_path / "license.lic"
    key_file.write_text("public-key")
    lic_file.write_text("license-content")

    result = CliRunner().invoke(
        cli,
        ["verify", "--public-key", str(key_file), "--license", str(lic_file), "--feature", "billing"],
    )

    assert result.exit_code != 0
    assert "Missing licensed features: billing" in result.output


def test_inspect_command_shows_payload(tmp_path):
    payload = {"id": "lic-001", "plan": "pro"}
    lic_file = tmp_path / "license.lic"
    lic_file.write_text(
        json.dumps(
            {
                "payload": base64.b64encode(json.dumps(payload).encode()).decode(),
                "signature": "YWJjMTIz",
            }
        )
    )

    result = CliRunner().invoke(cli, ["inspect", "--license", str(lic_file)])

    assert result.exit_code == 0
    assert '"id": "lic-001"' in result.output
    assert '"signature_prefix": "YWJjMTIz"' in result.output


def test_hardware_info_outputs_json(monkeypatch):
    monkeypatch.setattr(
        "lykn.cli.commands.hardware.collect_hardware",
        lambda: {"hostname": "demo-host", "mac_addresses": ["AA:BB:CC:DD:EE:FF"]},
    )

    result = CliRunner().invoke(cli, ["hardware-info", "--format", "json"])

    assert result.exit_code == 0
    assert '"hostname": "demo-host"' in result.output


def test_hardware_info_outputs_table(monkeypatch):
    monkeypatch.setattr(
        "lykn.cli.commands.hardware.collect_hardware",
        lambda: {"hostname": "demo-host", "mac_addresses": ["AA:BB:CC:DD:EE:FF"]},
    )

    result = CliRunner().invoke(cli, ["hardware-info", "--format", "table"])

    assert result.exit_code == 0
    assert "hostname" in result.output
    assert "demo-host" in result.output
