from pathlib import Path

import pytest

from lykn.pack.models import PackConfig, ResourceSpec
from lykn.pack.staging import build_staging


def make_config(project_dir: Path, resources: list[ResourceSpec] | None = None, exclude: list[str] | None = None):
    return PackConfig(
        project_dir=project_dir,
        engine="pyinstaller",
        entry_script=Path("main.py"),
        resources=resources or [],
        exclude=exclude or [],
    )


def test_build_staging_collects_default_resources(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    (project_dir / "main.py").write_text("print('hi')")
    (project_dir / "license.lic").write_text("license")
    (project_dir / "public.pem").write_text("public")

    result = build_staging(make_config(project_dir))

    assert (result.staging_dir / "license.lic").read_text() == "license"
    assert (result.staging_dir / "public.pem").read_text() == "public"
    assert [item.source for item in result.resources] == [Path("license.lic"), Path("public.pem")]


def test_build_staging_appends_explicit_resources(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    (project_dir / "main.py").write_text("print('hi')")
    (project_dir / "runtime").mkdir()
    (project_dir / "runtime" / "config.json").write_text("{}")

    config = make_config(
        project_dir,
        resources=[ResourceSpec(source=Path("runtime/config.json"))],
    )

    result = build_staging(config)

    assert (result.staging_dir / "runtime" / "config.json").read_text() == "{}"


def test_build_staging_writes_module_bootstrap_when_entry_module_used(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()

    config = PackConfig(
        project_dir=project_dir,
        engine="pyinstaller",
        entry_module="app.main",
    )

    result = build_staging(config)

    assert (result.staging_dir / "__main__.py").read_text() == (
        "import runpy\n"
        "runpy.run_module('app.main', run_name='__main__')\n"
    )


def test_build_staging_applies_default_and_custom_excludes(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    (project_dir / "main.py").write_text("print('hi')")
    (project_dir / ".env").write_text("SECRET=1")
    (project_dir / "notes.tmp").write_text("skip")

    result = build_staging(make_config(project_dir, exclude=["*.tmp"]))

    assert not (result.staging_dir / ".env").exists()
    assert not (result.staging_dir / "notes.tmp").exists()


def test_build_staging_fails_for_missing_explicit_resource(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    (project_dir / "main.py").write_text("print('hi')")

    config = make_config(
        project_dir,
        resources=[ResourceSpec(source=Path("missing/license.lic"))],
    )

    with pytest.raises(FileNotFoundError, match="missing/license.lic"):
        build_staging(config)
