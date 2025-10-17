param(
    [string]$Type,
    [switch]$NoLint,
    [switch]$NoTest,
    [string]$Output,
    [switch]$Release,
    [switch]$Help
)

$script:SkipTests = $env:SKIP_TESTS
$script:ReleaseType = "unknown"
$script:Version = ""
$script:CommitHash = ""
$script:BuildUnixTime = $env:BUILD_UNIXTIME
$script:BuildDate = $env:BUILD_DATE

function Write-Red($msg) {
    Write-Host $msg -ForegroundColor Red
}

function Check-Dependency {
    param([string[]]$commands)
    foreach ($cmd in $commands) {
        if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
            Write-Red "Error: `"$cmd`" is required."
            exit 127
        }
    }
}

function Show-Help {
    Write-Host "ezBookkeeping build script for Windows PowerShell"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "    build.ps1 type [options]"
    Write-Host ""
    Write-Host "Types:"
    Write-Host "    backend            Build backend binary file"
    Write-Host "    frontend           Build frontend files"
    Write-Host "    package            Build package archive"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "    -Release           Build release (The script will use environment variable `"RELEASE_BUILD`" to detect whether this is release building by default)"
    Write-Host "    -Output <filename> Package file name (for `"package`" type only)"
    Write-Host "    -NoLint            Do not execute lint check before building"
    Write-Host "    -NoTest            Do not execute unit testing before building (You can use environment variable `"SKIP_TESTS`" to skip specified tests)"
    Write-Host "    -Help              Show help"
}

function Parse-Args {
    if (-not $Type) {
        Show-Help
        exit 0
    }

    if ($Release -or $env:RELEASE_BUILD) {
        $script:ReleaseType = "release"
    } else {
        $script:ReleaseType = "snapshot"
    }
}

function Check-Type-Dependencies {
    Check-Dependency "git"

    switch ($Type.ToLower()) {
        "backend"  {
            Check-Dependency "go","gcc"
        }
        "frontend" {
            Check-Dependency "node","npm"
        }
        "package"  {
            Check-Dependency "go","gcc","node","npm","7z"
        }
    }
}

function Set-Build-Parameters {
    $script:Version = (Get-Content package.json | ConvertFrom-Json).version
    $script:CommitHash = git rev-parse --short=7 HEAD

    if (-not $BuildUnixTime) {
        $script:BuildUnixTime = [int][double]::Parse((Get-Date -UFormat %s))
    }

    if (-not $BuildDate) {
        $script:BuildDate = Get-Date -Format "yyyyMMdd"
    }
}

function Build-Backend {
    Write-Host "Pulling backend dependencies..."
    go get .

    if (-not $NoLint) {
        Write-Host "Executing backend lint checking..."
        go vet -v .\...

        if ($LASTEXITCODE -ne 0) {
            Write-Red "Error: Failed to pass lint checking"
            exit 1
        }
    }

    if (-not $NoTest) {
        Write-Host "Executing backend unit testing..."
        go clean -cache

        if (-not $SkipTests) {
            go test .\... -v
        } else {
            Write-Host "(Skip unit test `"$SkipTests`")"
            go test .\... -v -skip "$SkipTests"
        }

        if ($LASTEXITCODE -ne 0) {
            Write-Red "Error: Failed to pass unit testing"
            exit 1
        }
    }

    $backend_build_extra_arguments = "-X main.Version=$Version "
    $backend_build_extra_arguments = "$backend_build_extra_arguments -X main.CommitHash=$CommitHash"

    if (-not $Release) {
        $backend_build_extra_arguments += " -X main.BuildUnixTime=$BuildUnixTime"
    }

    Write-Host "Building backend binary file ($ReleaseType)..."

    $env:CGO_ENABLED = 1
    go build -a -v -trimpath -tags timetzdata -ldflags "-w -s -linkmode external -extldflags '-static' $backend_build_extra_arguments" -o ezbookkeeping.exe ezbookkeeping.go

    Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue
}

function Build-Frontend {
    Write-Host "Pulling frontend dependencies..."
    npm install

    if (-not $NoLint) {
        Write-Host "Executing frontend lint checking..."
        npm run lint

        if ($LASTEXITCODE -ne 0) {
            Write-Red "Error: Failed to pass lint checking"
            exit 1
        }
    }

    if (-not $NoTest) {
        Write-Host "Executing frontend unit testing..."

        npm run test

        if ($LASTEXITCODE -ne 0) {
            Write-Red "Error: Failed to pass unit testing"
            exit 1
        }
    }

    Write-Host "Building frontend files ($ReleaseType)..."

    if (-not $Release) {
        $env:buildUnixTime = $BuildUnixTime
        npm run build
        Remove-Item Env:\buildUnixTime -ErrorAction SilentlyContinue
    } else {
        npm run build
    }
}

function Build-Package {
    $packageFileName = "ezbookkeeping-$Version"

    if (-not $Release) {
        $packageFileName = "$packageFileName-$BuildDate"
    }

    $packageFileName = "$packageFileName-windows.zip"

    if ($Output) {
        $packageFileName = $Output
    }

    Write-Host "Building package archive '$packageFileName' ($ReleaseType)..."

    Build-Backend
    Build-Frontend

    Remove-Item package -Recurse -Force -ErrorAction SilentlyContinue
    New-Item -ItemType Directory -Path "package"
    New-Item -ItemType Directory -Path "package\data"
    New-Item -ItemType Directory -Path "package\storage"
    New-Item -ItemType Directory -Path "package\log"

    Copy-Item ezbookkeeping.exe package\
    Copy-Item dist package\public -Recurse
    Copy-Item conf package\conf -Recurse
    Copy-Item templates package\templates -Recurse
    Copy-Item LICENSE package\

    Push-Location package
    7z a -r -tzip -mx9 "..\$packageFileName" *
    Pop-Location
}

function Main {
    if ($Help) {
        Show-Help
        exit 0
    }

    Parse-Args
    Check-Type-Dependencies
    Set-Build-Parameters

    switch ($Type) {
        "backend"  {
            Build-Backend
        }
        "frontend" {
            Build-Frontend
        }
        "package"  {
            Build-Package
        }
        default    {
            Write-Red "Invalid type: $Type"
            Show-Help
            exit 2
        }
    }
}

Main
