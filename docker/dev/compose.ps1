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

function Test-FrontendRequested {
    param(
        [string[]]$Args
    )

    for ($i = 0; $i -lt $Args.Count; $i++) {
        $arg = $Args[$i]
        if ($arg -eq "frontend" -or $arg -eq "--profile=frontend") {
            return $true
        }
        if ($arg -eq "--profile" -and $i + 1 -lt $Args.Count -and $Args[$i + 1] -eq "frontend") {
            return $true
        }
    }

    return $false
}

function Get-AvailableMemoryMb {
    try {
        $os = Get-CimInstance Win32_OperatingSystem -ErrorAction Stop
        return [int][math]::Floor($os.FreePhysicalMemory / 1024)
    } catch {
        return $null
    }
}

function Protect-FrontendStart {
    param(
        [string[]]$Args
    )

    if (-not (Test-FrontendRequested -Args $Args)) {
        return
    }

    if ($env:ALLOW_LOW_MEMORY_FRONTEND -eq "1") {
        Write-Warning "Skipping frontend memory guard because ALLOW_LOW_MEMORY_FRONTEND=1"
        return
    }

    $minMb = 2048
    if ($env:FRONTEND_MIN_HOST_MEM_MB) {
        $parsedMinMb = 0
        if ([int]::TryParse($env:FRONTEND_MIN_HOST_MEM_MB, [ref]$parsedMinMb)) {
            $minMb = $parsedMinMb
        }
    }

    $availableMb = Get-AvailableMemoryMb
    if ($null -eq $availableMb) {
        Write-Warning "Unable to determine host available memory, continuing without frontend guard"
        return
    }

    if ($availableMb -lt $minMb) {
        throw "Refusing to start frontend: host available memory ${availableMb}MB is below FRONTEND_MIN_HOST_MEM_MB=${minMb}MB. Run backend-only compose, or set ALLOW_LOW_MEMORY_FRONTEND=1 if you accept the risk."
    }
}

if (-not (Test-Path $envSource)) {
    throw "Missing env file: $envSource"
}

Copy-Item $envSource $envTarget -Force
Write-Host "[INFO] Synced $envSource -> $envTarget" -ForegroundColor Green

if (-not $ComposeArgs -or $ComposeArgs.Count -eq 0) {
    $ComposeArgs = @("up", "-d", "--build")
}

Protect-FrontendStart -Args $ComposeArgs

& docker compose --env-file $envSource -f $composeFile @ComposeArgs
exit $LASTEXITCODE
