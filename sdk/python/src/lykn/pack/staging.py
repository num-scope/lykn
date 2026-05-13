import fnmatch
import shutil
import tempfile
from pathlib import Path

from lykn.pack.models import PackConfig, ResourceSpec, StagingResult

DEFAULT_EXCLUDES = [
    ".env",
    ".env.*",
    ".git",
    ".git/*",
    "__pycache__",
    "*.pyc",
    ".pytest_cache",
]

DEFAULT_RESOURCE_CANDIDATES = [
    Path("license.lic"),
    Path("public.pem"),
]


def _matches_any(path: str, patterns: list[str]) -> bool:
    return any(fnmatch.fnmatch(path, pattern) for pattern in patterns)


def _copy_resource(project_dir: Path, staging_dir: Path, resource: ResourceSpec) -> None:
    source = project_dir / resource.source
    if not source.exists():
        raise FileNotFoundError(str(resource.source))

    target = staging_dir / (Path(resource.target) if resource.target else resource.source)
    target.parent.mkdir(parents=True, exist_ok=True)
    if source.is_dir():
        shutil.copytree(source, target, dirs_exist_ok=True)
    else:
        shutil.copy2(source, target)


def _default_resources(project_dir: Path) -> list[ResourceSpec]:
    items: list[ResourceSpec] = []
    for candidate in DEFAULT_RESOURCE_CANDIDATES:
        if (project_dir / candidate).exists():
            items.append(ResourceSpec(source=candidate))
    return items


def build_staging(config: PackConfig) -> StagingResult:
    staging_dir = Path(tempfile.mkdtemp(prefix="lykn-pack-"))
    excludes = list(dict.fromkeys(DEFAULT_EXCLUDES + config.exclude))
    resources = _default_resources(config.project_dir)

    existing = {item.source for item in resources}
    for item in config.resources:
        if item.source not in existing:
            resources.append(item)
            existing.add(item.source)

    for child in config.project_dir.iterdir():
        rel = child.relative_to(config.project_dir).as_posix()
        if _matches_any(rel, excludes):
            continue
        if child.name in {"license.lic", "public.pem"}:
            continue
        target = staging_dir / child.name
        if child.is_dir():
            shutil.copytree(child, target, ignore=shutil.ignore_patterns(*excludes), dirs_exist_ok=True)
        else:
            shutil.copy2(child, target)

    for resource in resources:
        _copy_resource(config.project_dir, staging_dir, resource)

    if config.entry_module:
        (staging_dir / "__main__.py").write_text(
            "import runpy\n"
            f"runpy.run_module({config.entry_module!r}, run_name='__main__')\n"
        )

    return StagingResult(staging_dir=staging_dir, resources=resources)
