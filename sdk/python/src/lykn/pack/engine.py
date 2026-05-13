import subprocess
from pathlib import Path
from typing import Protocol

from lykn.pack.models import PackConfig, PackResult, StagingResult


class PackEngine(Protocol):
    def build_command(self, config: PackConfig, staging: StagingResult) -> list[str]: ...
    def expected_output(self, config: PackConfig) -> Path: ...


class PyInstallerEngine:
    def build_command(self, config: PackConfig, staging: StagingResult) -> list[str]:
        entry = str(staging.staging_dir / config.entry_script) if config.entry_script else config.entry_module
        command = [
            "pyinstaller",
            "--noconfirm",
            "--name",
            config.name or "app",
            "--distpath",
            str(config.output_dir),
            "--workpath",
            str(config.project_dir / ".lykn-build"),
        ]
        command.append("--onefile" if config.bundle_mode == "onefile" else "--onedir")
        if config.entry_module:
            command.extend(["--hidden-import", config.entry_module, str(staging.staging_dir / "__main__.py")])
        else:
            command.append(entry)
        command.extend(config.extra_args)
        return command

    def expected_output(self, config: PackConfig) -> Path:
        return config.output_dir / (config.name or "app")


class NuitkaEngine:
    def build_command(self, config: PackConfig, staging: StagingResult) -> list[str]:
        command = [
            "nuitka",
            "--assume-yes-for-downloads",
            f"--output-dir={config.output_dir}",
        ]
        command.append("--onefile" if config.bundle_mode == "onefile" else "--standalone")
        if config.name:
            command.append(f"--output-filename={config.name}")
        if config.entry_module:
            command.extend(["--module", config.entry_module])
        else:
            command.append(str(staging.staging_dir / config.entry_script))
        command.extend(config.extra_args)
        return command

    def expected_output(self, config: PackConfig) -> Path:
        return config.output_dir / (config.name or "app")


def get_engine(name: str) -> PackEngine:
    if name == "pyinstaller":
        return PyInstallerEngine()
    if name == "nuitka":
        return NuitkaEngine()
    raise ValueError(f"unsupported pack engine: {name}")


def run_pack(config: PackConfig, staging: StagingResult) -> PackResult:
    engine = get_engine(config.engine)
    command = engine.build_command(config, staging)
    subprocess.run(command, check=True, cwd=config.project_dir)
    return PackResult(
        engine=config.engine,
        command=command,
        output_path=engine.expected_output(config),
        resources=staging.resources,
    )
