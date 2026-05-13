from dataclasses import dataclass, field
from pathlib import Path

try:
    import tomllib
except ModuleNotFoundError:  # pragma: no cover
    import tomli as tomllib

from lykn.pack.models import PackConfig, ResourceSpec


@dataclass
class CLIOverrides:
    engine: str | None = None
    entry: str | None = None
    module: str | None = None
    output_dir: str | None = None
    bundle_mode: str | None = None
    resources: list[str] = field(default_factory=list)
    exclude: list[str] = field(default_factory=list)
    name: str | None = None
    clean: bool | None = None
    extra_args: list[str] = field(default_factory=list)


def _read_tool_pack(project_dir: Path) -> dict:
    pyproject = project_dir / "pyproject.toml"
    if not pyproject.exists():
        return {}
    data = tomllib.loads(pyproject.read_text())
    return data.get("tool", {}).get("lykn", {}).get("pack", {})


def _resource_specs(values: list[str]) -> list[ResourceSpec]:
    seen: set[Path] = set()
    items: list[ResourceSpec] = []
    for value in values:
        path = Path(value)
        if path in seen:
            continue
        seen.add(path)
        items.append(ResourceSpec(source=path))
    return items


def load_pack_config(project_dir: Path, overrides: CLIOverrides) -> PackConfig:
    tool_pack = _read_tool_pack(project_dir)
    if overrides.entry is not None and overrides.module is not None:
        raise ValueError("exactly one of entry_script or entry_module must be provided")
    if overrides.entry is not None:
        entry_script = overrides.entry
        entry_module = None
    elif overrides.module is not None:
        entry_script = None
        entry_module = overrides.module
    else:
        entry_script = tool_pack.get("entry")
        entry_module = tool_pack.get("module")
    resources = overrides.resources or tool_pack.get("resources", [])
    exclude = overrides.exclude or tool_pack.get("exclude", [])
    clean = overrides.clean if overrides.clean is not None else tool_pack.get("clean", True)
    extra_args = overrides.extra_args or tool_pack.get("extra_args", [])

    return PackConfig(
        project_dir=project_dir,
        engine=overrides.engine or tool_pack.get("engine", "pyinstaller"),
        entry_script=Path(entry_script) if entry_script else None,
        entry_module=entry_module,
        output_dir=Path(overrides.output_dir or tool_pack.get("out", "dist")),
        bundle_mode=overrides.bundle_mode or tool_pack.get("mode", "onedir"),
        resources=_resource_specs(resources),
        exclude=list(dict.fromkeys(exclude)),
        name=overrides.name or tool_pack.get("name"),
        clean=clean,
        extra_args=list(dict.fromkeys(extra_args)),
    )
