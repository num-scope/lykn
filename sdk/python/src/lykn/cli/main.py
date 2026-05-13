import click

from lykn.cli.commands.hardware import hardware_info_command
from lykn.cli.commands.inspect_cmd import inspect_command
from lykn.cli.commands.pack import pack_command
from lykn.cli.commands.verify import verify_command


@click.group()
def cli():
    """lykn command line interface."""


cli.add_command(verify_command)
cli.add_command(inspect_command)
cli.add_command(hardware_info_command)
cli.add_command(pack_command)


if __name__ == "__main__":
    cli()
