import base64
import json
from pathlib import Path

from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import padding

from lykn.exceptions import LicenseFileError, LicenseSignatureError


def load_public_key(source: str | Path):
    """Load an RSA public key from a file path or PEM string."""
    if isinstance(source, Path):
        pem_data = source.read_bytes()
    elif isinstance(source, str) and not source.strip().startswith("-----"):
        pem_data = Path(source).read_bytes()
    else:
        pem_data = source.encode() if isinstance(source, str) else source

    return serialization.load_pem_public_key(pem_data)


def verify_license_file(lic_content: str | bytes, public_key) -> dict:
    """Verify a .lic file's signature and return the decoded payload."""
    if isinstance(lic_content, bytes):
        lic_content = lic_content.decode()

    try:
        lic = json.loads(lic_content)
    except json.JSONDecodeError as exc:
        raise LicenseFileError(f"Invalid .lic JSON format: {exc}") from exc

    if "payload" not in lic or "signature" not in lic:
        raise LicenseFileError("Missing 'payload' or 'signature' field")

    try:
        payload_bytes = base64.b64decode(lic["payload"])
        signature_bytes = base64.b64decode(lic["signature"])
    except Exception as exc:
        raise LicenseFileError(f"Base64 decode error: {exc}") from exc

    try:
        public_key.verify(
            signature_bytes,
            payload_bytes,
            padding.PKCS1v15(),
            hashes.SHA256(),
        )
    except Exception as exc:
        raise LicenseSignatureError("License signature verification failed") from exc

    try:
        return json.loads(payload_bytes)
    except json.JSONDecodeError as exc:
        raise LicenseFileError(f"Invalid license payload JSON: {exc}") from exc
