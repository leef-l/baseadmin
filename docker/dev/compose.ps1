param(
    [switch]$China,
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$ComposeArgs
)

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent (Split-Path -Parent $scriptDir)
$envSource = Join-Path $scriptDir ".env"
$envTarget = Join-Path $projectRoot "admin-go\.env"
$composeFile = if ($China) {
    Join-Path $scriptDir "docker-compose.cn.yml"
} else {
    Join-Path $scriptDir "docker-compose.yml"
}

if (-not (Test-Path $envSource)) {
    throw "Missing env file: $envSource"
}

Copy-Item $envSource $envTarget -Force
Write-Host "[INFO] Synced $envSource -> $envTarget" -ForegroundColor Green

if (-not $ComposeArgs -or $ComposeArgs.Count -eq 0) {
    $ComposeArgs = @("up", "-d", "--build")
}

& docker compose --env-file $envSource -f $composeFile @ComposeArgs
exit $LASTEXITCODE
