from pathlib import Path

import pytest

from lykn.pack.engine import get_engine, run_pack
from lykn.pack.models import PackConfig, ResourceSpec, StagingResult


def make_config(engine: str, *, entry_script: str | None = "main.py", entry_module: str | None = None, mode: str = "onedir"):
    return PackConfig(
        project_dir=Path("/tmp/project"),
        engine=engine,
        entry_script=Path(entry_script) if entry_script else None,
        entry_module=entry_module,
        output_dir=Path("dist"),
        bundle_mode=mode,
        name="demo-app",
        extra_args=["--verbose"],
    )


def make_staging():
    return StagingResult(
        staging_dir=Path("/tmp/staging"),
        resources=[ResourceSpec(source=Path("license.lic"))],
    )


def test_pyinstaller_engine_builds_command_for_script_entry():
    engine = get_engine("pyinstaller")
    command = engine.build_command(make_config("pyinstaller"), make_staging())

    assert command[:2] == ["pyinstaller", "--noconfirm"]
    assert "--onedir" in command
    assert "--distpath" in command
    assert "/tmp/staging/main.py" in command


def test_nuitka_engine_builds_command_for_module_entry():
    engine = get_engine("nuitka")
    config = make_config("nuitka", entry_script=None, entry_module="app.main", mode="onefile")
    command = engine.build_command(config, make_staging())

    assert command[:2] == ["nuitka", "--assume-yes-for-downloads"]
    assert "--onefile" in command
    assert "--module" in command
    assert "app.main" in command


def test_run_pack_calls_subprocess_and_returns_result(monkeypatch):
    calls = {}

    def fake_run(command, check, cwd):
        calls["command"] = command
        calls["cwd"] = cwd

    monkeypatch.setattr("lykn.pack.engine.subprocess.run", fake_run)

    result = run_pack(make_config("pyinstaller"), make_staging())

    assert result.engine == "pyinstaller"
    assert result.command == calls["command"]
    assert result.output_path == Path("dist") / "demo-app"
    assert calls["cwd"] == Path("/tmp/project")


def test_get_engine_rejects_unknown_engine():
    with pytest.raises(ValueError, match="unsupported pack engine"):
        get_engine("cx_freeze")
