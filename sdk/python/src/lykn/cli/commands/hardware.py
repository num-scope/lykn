import json

import click
from rich.console import Console
from rich.table import Table

from lykn.hardware import collect_hardware


@click.command("hardware-info")
@click.option("--format", "output_format", type=click.Choice(["json", "table"]), default="json")
def hardware_info_command(output_format: str):
    hardware = collect_hardware()
    if hasattr(hardware, "model_dump"):
        payload = hardware.model_dump()
    else:
        payload = hardware

    if output_format == "json":
        click.echo(json.dumps(payload, indent=2, ensure_ascii=False))
        return

    table = Table(title="Hardware Info")
    table.add_column("field")
    table.add_column("value")
    for key, value in payload.items():
        text = json.dumps(value, ensure_ascii=False) if isinstance(value, list) else str(value)
        table.add_row(key, text)
    Console().print(table)
