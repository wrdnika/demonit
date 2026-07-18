# Demonit device agent (Windows PowerShell)
# Prefer scripts/device_agent.py for cross-platform OSS usage.
#
# Usage:
#   .\device-agent.ps1 -DeviceId "YOUR-UUID"
#   .\device-agent.ps1 -DeviceId "..." -ApiBase "http://192.168.1.10:8080" -IntervalSec 10
# From Git Bash:
#   powershell.exe -ExecutionPolicy Bypass -File ./device-agent.ps1 -DeviceId "YOUR-UUID"

param(
  [Parameter(Mandatory = $true)]
  [string]$DeviceId,

  [string]$ApiBase = "http://localhost:8080",
  [string]$DeviceApiKey = "dev-device-key-change-me",
  [int]$IntervalSec = 10
)

$ErrorActionPreference = "Stop"
Write-Host "Demonit agent started" -ForegroundColor Green
Write-Host "  device_id : $DeviceId"
Write-Host "  api       : $ApiBase"
Write-Host "  interval  : ${IntervalSec}s"
Write-Host "Ctrl+C to stop.`n"

function Get-CpuUsage {
  try {
    $p = Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average
    return [math]::Round([double]$p.Average, 1)
  }
  catch {
    return Get-Random -Minimum 5 -Maximum 55
  }
}

function Get-RamUsage {
  try {
    $os = Get-CimInstance Win32_OperatingSystem
    $total = [double]$os.TotalVisibleMemorySize
    $free = [double]$os.FreePhysicalMemory
    if ($total -le 0) { return 0 }
    return [math]::Round((($total - $free) / $total) * 100, 1)
  }
  catch {
    return Get-Random -Minimum 20 -Maximum 80
  }
}

while ($true) {
  $cpu = Get-CpuUsage
  $ram = Get-RamUsage
  $payload = @{
    device_id      = $DeviceId
    cpu_usage      = $cpu
    ram_usage      = $ram
    status_payload = @{
      hostname    = $env:COMPUTERNAME
      os          = "windows"
      agent       = "powershell"
      temperature = [math]::Round((Get-Random -Minimum 350 -Maximum 520) / 10, 1)
    }
  } | ConvertTo-Json -Depth 5

  try {
    $res = Invoke-RestMethod `
      -Uri "$ApiBase/api/v1/heartbeat" `
      -Method POST `
      -Headers @{
        "Content-Type"    = "application/json"
        "X-Device-API-Key" = $DeviceApiKey
      } `
      -Body $payload

    $ts = Get-Date -Format "HH:mm:ss"
    Write-Host "[$ts] OK  cpu=$cpu% ram=$ram%  status=$($res.data.status)"
  }
  catch {
    $ts = Get-Date -Format "HH:mm:ss"
    Write-Host "[$ts] FAIL $($_.Exception.Message)" -ForegroundColor Red
  }

  Start-Sleep -Seconds $IntervalSec
}
