# Registers stats as a scheduled task that starts at logon and restarts if it exits.
# Run from an elevated PowerShell:
#   powershell -ExecutionPolicy Bypass -File packaging/windows/stats.ps1 -Binary C:\Tools\stats.exe

param(
    [string]$Binary = "$PSScriptRoot\stats.exe",
    [int]$Port = 7679
)

$action = New-ScheduledTaskAction -Execute $Binary -Argument "-port $Port"
$trigger = New-ScheduledTaskTrigger -AtLogOn
$settings = New-ScheduledTaskSettingsSet -RestartCount 3 -RestartInterval (New-TimeSpan -Minutes 1) -AllowStartIfOnBatteries
Register-ScheduledTask -TaskName "stats" -Action $action -Trigger $trigger -Settings $settings -Force | Out-Null
Write-Output "stats will listen on 127.0.0.1:$Port at next logon"