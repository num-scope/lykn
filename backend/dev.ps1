# Lykn backend helper for Windows PowerShell
# Usage:
#   .\dev.ps1 help
#   .\dev.ps1 run
#   .\dev.ps1 test
#   .\dev.ps1 build
#   .\dev.ps1 clean
#   .\dev.ps1 fmt
#   .\dev.ps1 demo
#   .\dev.ps1 docker-build

[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [ValidateSet('help', 'run', 'test', 'build', 'clean', 'fmt', 'demo', 'docker-build')]
    [string]$Command = 'help',

    [string]$Config = 'config/config.yaml',
    [string]$Binary = 'server.exe',
    [string]$DockerImage = 'lykn-server:dev',
    [string]$DemoOut = '..\tests\fixtures',
    [string]$Go = 'go'
)

$ErrorActionPreference = 'Stop'
Set-Location -Path $PSScriptRoot

function Show-Help {
    @"
Lykn server PowerShell targets:
  .\dev.ps1 run           Run the HTTP server with LYKN_CONFIG=$Config
  .\dev.ps1 test          Run all Go tests
  .\dev.ps1 build         Build server binary to $Binary
  .\dev.ps1 clean         Remove build output
  .\dev.ps1 fmt           Run gofmt on all tracked Go packages
  .\dev.ps1 demo          Generate complete demo fixtures in $DemoOut
  .\dev.ps1 docker-build  Build backend Docker image $DockerImage

Optional parameters:
  -Config config/config.yaml
  -Binary server.exe
  -DemoOut ..\tests\fixtures
"@ | Write-Host
}

switch ($Command) {
    'help' {
        Show-Help
    }
    'run' {
        $env:LYKN_CONFIG = $Config
        & $Go run ./cmd/server
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    }
    'test' {
        & $Go test ./...
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    }
    'build' {
        & $Go build -o $Binary ./cmd/server
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
        Write-Host "Built $Binary"
    }
    'clean' {
        @($Binary, 'server', 'server.exe') | ForEach-Object {
            if (Test-Path $_) {
                Remove-Item -Force $_
                Write-Host "Removed $_"
            }
        }
        if (Test-Path 'bin') {
            Remove-Item -Recurse -Force 'bin'
            Write-Host 'Removed bin/'
        }
    }
    'fmt' {
        & $Go fmt ./...
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    }
    'demo' {
        & $Go run ./cmd/demo $DemoOut
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    }
    'docker-build' {
        docker build -t $DockerImage .
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    }
}
