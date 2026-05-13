from click.testing import CliRunner

from lykn.cli.main import cli


class FakePackResult:
    engine = "pyinstaller"
    command = ["pyinstaller", "--noconfirm", "main.py"]
    output_path = "dist/demo-app"
    resources = []


def test_pack_command_uses_entry_script(monkeypatch, tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    (project_dir / "main.py").write_text("print('hi')")

    captured = {}

    def fake_load_pack_config(project_dir, overrides):
        captured["project_dir"] = project_dir
        captured["overrides"] = overrides
        return "CONFIG"

    monkeypatch.setattr("lykn.cli.commands.pack.load_pack_config", fake_load_pack_config)
    monkeypatch.setattr("lykn.cli.commands.pack.build_staging", lambda config: "STAGING")
    monkeypatch.setattr("lykn.cli.commands.pack.run_pack", lambda config, staging: FakePackResult())

    result = CliRunner().invoke(
        cli,
        [
            "pack",
            "--project",
            str(project_dir),
            "--engine",
            "pyinstaller",
            "--entry",
            "main.py",
            "--out",
            "dist",
        ],
    )

    assert result.exit_code == 0
    assert "dist/demo-app" in result.output
    assert captured["project_dir"] == project_dir
    assert captured["overrides"].entry == "main.py"


def test_pack_command_uses_module_entry(monkeypatch, tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()

    monkeypatch.setattr("lykn.cli.commands.pack.load_pack_config", lambda project_dir, overrides: "CONFIG")
    monkeypatch.setattr("lykn.cli.commands.pack.build_staging", lambda config: "STAGING")
    monkeypatch.setattr("lykn.cli.commands.pack.run_pack", lambda config, staging: FakePackResult())

    result = CliRunner().invoke(
        cli,
        [
            "pack",
            "--project",
            str(project_dir),
            "--engine",
            "nuitka",
            "--module",
            "app.main",
        ],
    )

    assert result.exit_code == 0
    assert "pyinstaller" in result.output or "nuitka" in result.output


def test_pack_command_reports_failure(monkeypatch, tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()

    monkeypatch.setattr("lykn.cli.commands.pack.load_pack_config", lambda project_dir, overrides: "CONFIG")
    monkeypatch.setattr("lykn.cli.commands.pack.build_staging", lambda config: "STAGING")

    def fail(config, staging):
        raise RuntimeError("pack failed")

    monkeypatch.setattr("lykn.cli.commands.pack.run_pack", fail)

    result = CliRunner().invoke(
        cli,
        [
            "pack",
            "--project",
            str(project_dir),
            "--engine",
            "pyinstaller",
            "--entry",
            "main.py",
        ],
    )

    assert result.exit_code != 0
    assert "pack failed" in result.output


def test_pack_command_rejects_entry_and_module_together(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()

    result = CliRunner().invoke(
        cli,
        [
            "pack",
            "--project",
            str(project_dir),
            "--entry",
            "main.py",
            "--module",
            "app.main",
        ],
    )

    assert result.exit_code != 0
    assert "exactly one of entry_script or entry_module" in result.output
