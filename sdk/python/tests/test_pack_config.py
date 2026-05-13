from pathlib import Path

import pytest

from lykn.pack.config import CLIOverrides, load_pack_config


def write_pyproject(path: Path, body: str) -> None:
    path.write_text(body)


def test_load_pack_config_reads_tool_section(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    write_pyproject(
        project_dir / "pyproject.toml",
        """
[tool.lykn.pack]
engine = "pyinstaller"
entry = "main.py"
out = "dist"
mode = "onedir"
resources = ["license.lic", "public.pem"]
exclude = [".env", ".git"]
name = "demo-app"
clean = true
extra_args = ["--noconfirm"]
""".strip(),
    )

    config = load_pack_config(project_dir=project_dir, overrides=CLIOverrides())

    assert config.project_dir == project_dir
    assert config.engine == "pyinstaller"
    assert config.entry_script == Path("main.py")
    assert config.entry_module is None
    assert config.output_dir == Path("dist")
    assert config.bundle_mode == "onedir"
    assert [item.source for item in config.resources] == [Path("license.lic"), Path("public.pem")]
    assert config.exclude == [".env", ".git"]
    assert config.name == "demo-app"
    assert config.clean is True
    assert config.extra_args == ["--noconfirm"]


def test_load_pack_config_cli_overrides_replace_tool_values(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    write_pyproject(
        project_dir / "pyproject.toml",
        """
[tool.lykn.pack]
engine = "pyinstaller"
entry = "main.py"
out = "dist"
mode = "onedir"
""".strip(),
    )

    overrides = CLIOverrides(
        engine="nuitka",
        module="app.main",
        output_dir="build",
        bundle_mode="onefile",
        resources=["runtime/license.lic"],
        exclude=[".env.*"],
        name="cli-name",
        clean=False,
        extra_args=["--nofollow-import-to=tkinter"],
    )

    config = load_pack_config(project_dir=project_dir, overrides=overrides)

    assert config.engine == "nuitka"
    assert config.entry_script is None
    assert config.entry_module == "app.main"
    assert config.output_dir == Path("build")
    assert config.bundle_mode == "onefile"
    assert [item.source for item in config.resources] == [Path("runtime/license.lic")]
    assert config.exclude == [".env.*"]
    assert config.name == "cli-name"
    assert config.clean is False
    assert config.extra_args == ["--nofollow-import-to=tkinter"]


def test_load_pack_config_rejects_entry_and_module_together(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    write_pyproject(
        project_dir / "pyproject.toml",
        """
[tool.lykn.pack]
engine = "pyinstaller"
entry = "main.py"
module = "app.main"
""".strip(),
    )

    with pytest.raises(ValueError, match="exactly one of entry_script or entry_module"):
        load_pack_config(project_dir=project_dir, overrides=CLIOverrides())


def test_load_pack_config_rejects_cli_entry_and_module_together(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()

    with pytest.raises(ValueError, match="exactly one of entry_script or entry_module"):
        load_pack_config(
            project_dir=project_dir,
            overrides=CLIOverrides(entry="main.py", module="app.main"),
        )


def test_load_pack_config_requires_entry_when_not_provided_anywhere(tmp_path):
    project_dir = tmp_path / "demo"
    project_dir.mkdir()
    write_pyproject(
        project_dir / "pyproject.toml",
        """
[tool.lykn.pack]
engine = "pyinstaller"
""".strip(),
    )

    with pytest.raises(ValueError, match="exactly one of entry_script or entry_module"):
        load_pack_config(project_dir=project_dir, overrides=CLIOverrides())
