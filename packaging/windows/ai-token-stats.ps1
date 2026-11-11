# Start at logon, from an elevated PowerShell:
#   powershell -ExecutionPolicy Bypass -File packaging/windows/ai-token-stats.ps1 -Binary C:\Tools\ai-token-stats.exe

param(
    [string]$Binary = "$PSScriptRoot\ai-token-stats.exe",
    [int]$Port = 7679
)

$action = New-ScheduledTaskAction -Execute $Binary -Argument "-port $Port"
$trigger = New-ScheduledTaskTrigger -AtLogOn
$settings = New-ScheduledTaskSettingsSet -RestartCount 3 -RestartInterval (New-TimeSpan -Minutes 1) -AllowStartIfOnBatteries
Register-ScheduledTask -TaskName "ai-token-stats" -Action $action -Trigger $trigger -Settings $settings -Force | Out-Null
Write-Output "ai-token-stats will listen on 127.0.0.1:$Port at next logon"
