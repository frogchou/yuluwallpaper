$logDir = [Environment]::GetFolderPath("ApplicationData")
$logPath = Join-Path $logDir "yuluwallpaper\yuluwallpaper.log"

if (Test-Path $logPath) {
    Get-Content -Path $logPath -Tail 200
} else {
    Write-Host "Log not found: $logPath"
}