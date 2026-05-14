# lykn Go SDK

Go license validation SDK for applications protected by Lykn.

## What this SDK provides

- Runtime license verification with `Validator`.
- RSA signature, validity period, hardware binding, and feature checks.
- CLI tools for verification, license inspection, and hardware collection.
- Gin middleware for route-level license and feature enforcement.

## Installation

Install the SDK:

```bash
go get github.com/wu-clan/lykn/sdk/go
```

Install the CLI:

```bash
go install github.com/wu-clan/lykn/sdk/go/cmd/lykn@latest
```

For local development from this repository:

```bash
cd sdk/go
go test ./...
```

## Verify a license

Use `Validator` with a project public key and a signed license file:

```go
package main

import (
    "log"
    "time"

    lykn "github.com/wu-clan/lykn/sdk/go"
)

func main() {
    validator, err := lykn.NewValidator(lykn.ValidatorOptions{
        PublicKeyPath: "./public.pem",
        LicensePath: "./license.lic",
        CheckInterval: 5 * time.Minute,
    })
    if err != nil {
        log.Fatal(err)
    }

    license, err := validator.Verify("reports")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("license", license.ID)
}
```

`Verify` returns `*lykn.LicenseData`, including:

- `ID`: license ID.
- `Subject`: licensed user or organization information.
- `Plan`: license plan code.
- `PlanName`: license plan display name.
- `IssuedAt`, `NotBefore`, `NotAfter`: license time window.
- `Hardware`: optional hardware binding data.
- `Features`: enabled feature list.
- `Limits`: structured license limits such as `MaxUsers` and `MaxDevices`.
- `Metadata`: legacy custom metadata, kept for older licenses.

Validation can return errors that support `errors.Is`:

- `ErrLicenseFile`
- `ErrLicenseSignature`
- `ErrLicenseExpired`
- `ErrLicenseNotYetValid`
- `ErrHardwareMismatch`
- `ErrFeatureNotLicensed`

## CLI

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

## Gin

Use `RequireLicense` to require any valid license, or `RequireFeatures` to require specific licensed features:

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    lykn "github.com/wu-clan/lykn/sdk/go"
    lykngin "github.com/wu-clan/lykn/sdk/go/middleware/gin"
)

func main() {
    validator, err := lykn.NewValidator(lykn.ValidatorOptions{
        PublicKeyPath: "./public.pem",
        LicensePath: "./license.lic",
    })
    if err != nil {
        panic(err)
    }

    r := gin.Default()
    r.GET("/health/license", lykngin.RequireLicense(validator), func(c *gin.Context) {
        license := lykngin.MustLicense(c)
        c.JSON(http.StatusOK, gin.H{"license": license.ID})
    })
    r.GET("/reports", lykngin.RequireFeatures(validator, "reports"), func(c *gin.Context) {
        license := lykngin.MustLicense(c)
        c.JSON(http.StatusOK, gin.H{"plan": license.Plan})
    })
    _ = r.Run()
}
```

Invalid licenses or missing features return HTTP 403 from the middleware.

## Compatibility

- Go 1.22 or newer.
- Optional Gin integration uses `github.com/gin-gonic/gin`.
- License files must use Lykn's existing `payload` and `signature` JSON format.

## Development

```bash
cd sdk/go
go test ./...
```
