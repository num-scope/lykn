from pathlib import Path

import click

from lykn.pack.config import CLIOverrides, load_pack_config
from lykn.pack.engine import run_pack
from lykn.pack.staging import build_staging


@click.command("pack")
@click.option("--project", "project_dir", required=True, type=click.Path(exists=True, file_okay=False, path_type=Path))
@click.option("--engine", default=None, type=click.Choice(["pyinstaller", "nuitka"]))
@click.option("--entry", default=None)
@click.option("--module", default=None)
@click.option("--out", "output_dir", default=None)
@click.option("--mode", "bundle_mode", type=click.Choice(["onedir", "onefile"]), default=None)
@click.option("--resource", "resources", multiple=True)
@click.option("--exclude", "exclude", multiple=True)
@click.option("--name", default=None)
@click.option("--clean/--no-clean", default=None)
@click.option("--extra-arg", "extra_args", multiple=True)
def pack_command(
    project_dir: Path,
    engine: str | None,
    entry: str | None,
    module: str | None,
    output_dir: str | None,
    bundle_mode: str | None,
    resources: tuple[str, ...],
    exclude: tuple[str, ...],
    name: str | None,
    clean: bool | None,
    extra_args: tuple[str, ...],
):
    overrides = CLIOverrides(
        engine=engine,
        entry=entry,
        module=module,
        output_dir=output_dir,
        bundle_mode=bundle_mode,
        resources=list(resources),
        exclude=list(exclude),
        name=name,
        clean=clean,
        extra_args=list(extra_args),
    )

    try:
        config = load_pack_config(project_dir=project_dir, overrides=overrides)
        staging = build_staging(config)
        result = run_pack(config, staging)
    except Exception as exc:
        raise click.ClickException(str(exc)) from exc

    click.echo(f"Pack engine: {result.engine}")
    click.echo(f"Output: {result.output_path}")
    click.echo(f"Command: {' '.join(result.command)}")
