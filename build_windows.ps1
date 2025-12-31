$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$icon = Join-Path $root "favicon.ico"
$pngIcon = Join-Path $root "favicon.png"
$trayIcon = Join-Path $root "cmd/yuluwallpaper/favicon.png"
$syso = Join-Path $root "cmd/yuluwallpaper/icon_windows.syso"

if (-not (Test-Path $icon)) {
    throw "Icon not found: $icon"
}
if (-not (Test-Path $pngIcon)) {
    throw "Icon not found: $pngIcon"
}

Copy-Item -Force $pngIcon $trayIcon

if (-not (Get-Command rsrc -ErrorAction SilentlyContinue)) {
    Write-Host "rsrc not found, installing..."
    go install github.com/akavel/rsrc@latest
}

& rsrc -ico $icon -o $syso
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

go build -ldflags "-H=windowsgui" -o yuluwallpaper.exe ./cmd/yuluwallpaper
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}
Write-Host "Built yuluwallpaper.exe"
