import base64
import json
from pathlib import Path

import click


@click.command("inspect")
@click.option("--license", "license_path", required=True, type=click.Path(exists=True, dir_okay=False))
def inspect_command(license_path: str):
    content = json.loads(Path(license_path).read_text())
    payload = json.loads(base64.b64decode(content["payload"]))
    signature = content.get("signature", "")
    data = {
        "payload": payload,
        "signature_prefix": signature[:16],
        "signature_length": len(signature),
    }
    click.echo(json.dumps(data, indent=2, ensure_ascii=False))
