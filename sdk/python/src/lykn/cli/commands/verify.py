import click

from lykn.validator import LicenseValidator


@click.command("verify")
@click.option("--public-key", required=True, type=click.Path(exists=True, dir_okay=False))
@click.option("--license", "license_path", required=True, type=click.Path(exists=True, dir_okay=False))
@click.option("--feature", "features", multiple=True)
def verify_command(public_key: str, license_path: str, features: tuple[str, ...]):
    validator = LicenseValidator(public_key=public_key, license_path=license_path)
    try:
        result = validator.verify(required_features=list(features) or None)
    except Exception as exc:
        raise click.ClickException(str(exc)) from exc

    click.echo(f"License valid: {result.id}")
    click.echo(f"Plan: {result.plan}")
    click.echo(f"Features: {', '.join(result.features)}")
