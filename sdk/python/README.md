# lykn Python SDK

Python license validation SDK for applications protected by Lykn.

## What this SDK provides

- Runtime license verification with `LicenseValidator`.
- RSA signature, validity period, hardware binding, and feature checks.
- CLI tools for verification, license inspection, hardware collection, and packaging.
- FastAPI dependencies for route-level license and feature enforcement.
- Packaging helpers for bundling protected Python applications with PyInstaller or Nuitka.

## Installation

Install the base SDK:

```bash
pip install lykn
```

Install optional integrations when needed:

```bash
pip install lykn[fastapi]
pip install lykn[pack]
```

For local development from this repository:

```bash
cd sdk/python
pip install -e ".[dev]"
```

## Verify a license

Use `LicenseValidator` with a project public key and a signed license file:

```python
from lykn import LicenseValidator

validator = LicenseValidator(
    public_key="./public.pem",
    license_path="./license.lic",
    check_interval=300,
)

license_data = validator.start()
assert validator.has_feature("reports") is True
validator.stop()
```

`LicenseValidator` accepts either `license_path` or `license_content`.
Set `check_interval` to a positive number of seconds to re-check the license in
the background after `start()` is called.

Common methods:

- `verify(required_features=None)`: validate once and return `LicenseData`.
- `start()`: validate once and start background re-checks when enabled.
- `stop()`: stop the background re-check loop.
- `has_feature(feature)`: check whether the current license includes a feature.
- `on_invalid(callback)`: register a callback for background validation failure.

## License data

Successful validation returns a `LicenseData` object with these fields:

- `id`: license ID.
- `subject`: licensed user or organization information.
- `plan`: license plan code.
- `plan_name`: license plan display name.
- `issued_at`, `not_before`, `not_after`: license time window.
- `hardware`: optional hardware binding data.
- `features`: enabled feature list.
- `limits`: structured license limits such as `max_users` and `max_devices`.
- `metadata`: legacy custom metadata, kept for older licenses.

Validation can raise these SDK errors:

- `LicenseFileError`
- `LicenseSignatureError`
- `LicenseExpiredError`
- `LicenseNotYetValidError`
- `HardwareMismatchError`
- `FeatureNotLicensedError`

## CLI

The `lykn` command is installed with the SDK.

Verify a license and require one or more features:

```bash
lykn verify --public-key ./public.pem --license ./license.lic --feature reports
```

Inspect a signed license payload without verifying it:

```bash
lykn inspect --license ./license.lic
```

Collect local hardware data for hardware-bound licenses:

```bash
lykn hardware-info --format json
lykn hardware-info --format table
```

Bundle a protected Python application:

```bash
lykn pack --project ./my-app --engine pyinstaller --entry main.py --out ./dist
lykn pack --project ./my-app --engine nuitka --module my_app --mode onefile
```

## FastAPI

Install the FastAPI extra first:

```bash
pip install lykn[fastapi]
```

Use `RequireLicense` to require any valid license, or `RequireFeatures` to
require specific licensed features:

```python
from fastapi import Depends, FastAPI
from lykn import LicenseValidator
from lykn.contrib.fastapi import RequireFeatures, RequireLicense

validator = LicenseValidator(public_key="./public.pem", license_path="./license.lic")
app = FastAPI()

@app.get("/health/license")
def license_health(license_data=Depends(RequireLicense(validator))):
    return {"license": license_data.id}

@app.get("/reports")
def reports(license_data=Depends(RequireFeatures(validator, "reports"))):
    return {"plan": license_data.plan}
```

Invalid licenses or missing features return HTTP 403 from the dependency.

## Packaging protected apps

Install the packaging extra first:

```bash
pip install lykn[pack]
```

`lykn pack` stages the project, includes common Lykn resources such as
`license.lic` and `public.pem` when present, and delegates the build to
PyInstaller or Nuitka.

Use CLI options for one-off builds:

```bash
lykn pack \
  --project ./my-app \
  --engine pyinstaller \
  --entry main.py \
  --mode onedir \
  --resource license.lic \
  --resource public.pem \
  --exclude .env \
  --name my-app
```

Or store defaults in `pyproject.toml`:

```toml
[tool.lykn.pack]
engine = "pyinstaller"
entry = "main.py"
out = "dist"
mode = "onedir"
resources = ["license.lic", "public.pem"]
exclude = [".env", ".env.*", ".git", "__pycache__"]
name = "my-app"
clean = true
extra_args = ["--noconfirm"]
```

CLI options override values from `pyproject.toml`.

## Compatibility

- Python 3.10 or newer.
- Optional FastAPI integration requires `fastapi>=0.100`.
- Optional packaging support requires `pyinstaller>=6.0` or `nuitka>=2.0`.
