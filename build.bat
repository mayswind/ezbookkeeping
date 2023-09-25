@echo off

set "TYPE="
set "NO_LINT=0"
set "NO_TEST=0"
set "RELEASE=%RELEASE_BUILD%"
set "RELEASE_TYPE=unknown"
set "VERSION="
set "COMMIT_HASH="
set "BUILD_UNIXTIME="
set "BUILD_DATE="
set "PACKAGE_FILENAME="
for /f %%a in ('"prompt $E$S & echo on & for %%b in (1) do rem"') do set "ESC=%%a"

if "%~1"=="" call :show_help & goto :end
goto :pre_parse_args

:echo_red
    echo %ESC%[91m%~1%ESC%[0m
    goto :eof

:set_unixtime
    setlocal enableextensions
    for /f %%x in ('wmic path win32_utctime get /format:list ^| findstr "="') do set %%x
    set /a z=(14-100%Month%%%100)/12, y=10000%Year%%%10000-z
    set /a ut=y*365+y/4-y/100+y/400+(153*(100%Month%%%100+12*z-3)+2)/5+Day-719469
    set /a ut=ut*86400+100%Hour%%%100*3600+100%Minute%%%100*60+100%Second%%%100
    endlocal & set "%1=%ut%" & goto :eof

:set_date
    setlocal enableextensions
    for /f %%x in ('wmic path win32_localtime get /format:list ^| findstr "="') do set %%x
    if %Month% lss 10 set "Month=0%Month%"
    if %Day% lss 10 set "Day=0%Day%"
    endlocal & set "%1=%Year%%Month%%Day%" & goto :eof

:check_dependency
    if "%~1"=="" goto :eof
    where /q %~1 || call :echo_red "Error: "%~1" is required." && goto :end

    shift
    goto :check_dependency

:show_help
    echo ezBookkeeping build script for Windows
    echo.
    echo Usage:
    echo     build.cmd type [options]
    echo.
    echo Types:
    echo     backend                 Build backend binary file
    echo     frontend                Build frontend files
    echo     package                 Build package archive
    echo.
    echo Options:
    echo     /r, --release           Build release (The script will use environment variable "RELEASE_BUILD" to detect whether this is release building by default)
    echo     /o, --output ^<filename^> Package file name (For "package" type only)
    echo     --no-lint               Do not execute lint check before building
    echo     --no-test               Do not execute unit testing before building
    echo     /h, --help              Show help
    goto :eof

:pre_parse_args
    if "%~1"=="" goto :post_parse_args

    if /i "%~1"=="backend"   set "TYPE=%~1" & shift
    if /i "%~1"=="frontend"  set "TYPE=%~1" & shift
    if /i "%~1"=="package"   set "TYPE=%~1" & shift

:parse_args
    if "%~1"=="" goto :post_parse_args

    if /i "%~1"=="/r"        set "RELEASE=1" & shift & goto :parse_args
    if /i "%~1"=="-r"        set "RELEASE=1" & shift & goto :parse_args
    if /i "%~1"=="--release" set "RELEASE=1" & shift & goto :parse_args

    if /i "%~1"=="/o"        set "PACKAGE_FILENAME=%~2" & shift & shift & goto :parse_args
    if /i "%~1"=="-o"        set "PACKAGE_FILENAME=%~2" & shift & shift & goto :parse_args
    if /i "%~1"=="--output"  set "PACKAGE_FILENAME=%~2" & shift & shift & goto :parse_args

    if /i "%~1"=="--no-lint" set "NO_LINT=1" & shift & goto :parse_args
    if /i "%~1"=="--no-test" set "NO_TEST=1" & shift & goto :parse_args

    if /i "%~1"=="/h"        call :show_help & goto :end
    if /i "%~1"=="-h"        call :show_help & goto :end
    if /i "%~1"=="--help"    call :show_help & goto :end

    call :echo_red "Invalid argument: %~1" & call :show_help & goto :end

:post_parse_args
    if "%RELEASE%"=="" set "RELEASE=0"

    if "%RELEASE%"=="0" (
        set "RELEASE_TYPE=snapshot"
    ) else (
        set "RELEASE_TYPE=release"
    )

:check_type_dependencies
    if not defined TYPE call :echo_red "Error: No specified type" & call :show_help & goto :end

    call :check_dependency "git"
    if "%TYPE%"=="backend"  call :check_dependency "go" "gcc"
    if "%TYPE%"=="frontend" call :check_dependency "node" "npm"
    if "%TYPE%"=="package"  call :check_dependency "go" "gcc" "node" "npm" "7z"

    if not "%errorlevel%"=="0" goto :end

:set_build_parameters
    for /f "tokens=2 delims=:" %%x in ('findstr "\"version\": \"*\"," package.json') do set "VERSION=%%x"
    set VERSION=%VERSION: =%
    set VERSION=%VERSION:,=%
    set VERSION=%VERSION:"=%
    for /f %%x in ('git rev-parse --short HEAD') do set "COMMIT_HASH=%%x"
    call :set_unixtime BUILD_UNIXTIME
    call :set_date BUILD_DATE

:main
    if "%TYPE%"=="backend"  call :build_backend & goto :end
    if "%TYPE%"=="frontend" call :build_frontend & goto :end
    if "%TYPE%"=="package"  call :build_package & goto :end
    goto :end

:build_backend
    setlocal enabledelayedexpansion
    echo Pulling backend dependencies...
    call go get .

    if "%NO_LINT%"=="0" (
        echo Executing backend lint checking...
        call go vet -v .\...

        if !errorlevel! neq 0 (
            call :echo_red "Error: Failed to pass lint checking"
            goto :end
        )
    )

    if "%NO_TEST%"=="0" (
        echo Executing backend unit testing...
        call go clean -cache
        call go test .\... -v

        if !errorlevel! neq 0 (
            call :echo_red "Error: Failed to pass unit testing"
            goto :end
        )
    )

    endlocal

    set "CGO_ENABLED=1"

    setlocal
    set "backend_build_extra_arguments=-X main.Version=%VERSION%"
    set "backend_build_extra_arguments=%backend_build_extra_arguments% -X main.CommitHash=%COMMIT_HASH%"

    if "%RELEASE%"=="0" (
        set "backend_build_extra_arguments=%backend_build_extra_arguments% -X main.BuildUnixTime=%BUILD_UNIXTIME%"
    )

    echo Building backend binary file (%RELEASE_TYPE%)...

    call go build -a -v -trimpath -tags timetzdata -ldflags "-w -s -linkmode external -extldflags '-static' %backend_build_extra_arguments%" -o ezbookkeeping.exe ezbookkeeping.go
    endlocal

    set "CGO_ENABLED="

    goto :eof

:build_frontend
    setlocal enabledelayedexpansion
    echo Pulling frontend dependencies...
    call npm install

    if "%NO_LINT%"=="0" (
        echo Executing frontend lint checking...

        call npm run lint

        if !errorlevel! neq 0 (
            call :echo_red "Error: Failed to pass lint checking"
            goto :end
        )
    )

    endlocal

    echo Building frontend files(%RELEASE_TYPE%)...

    if "%RELEASE%"=="0" (
        set "buildUnixTime=%BUILD_UNIXTIME%"
        call npm run build
        set "buildUnixTime="
    ) else (
        call npm run build
    )

    goto :eof

:build_package
    setlocal enabledelayedexpansion
    set "package_file_name=%VERSION%"

    if "%RELEASE%"=="0" (
        set "build_date="
        set "package_file_name=%package_file_name%-%build_date%"
    )

    set "package_file_name=ezbookkeeping-%package_file_name%-windows.zip"

    if defined PACKAGE_FILENAME set set "package_file_name=%PACKAGE_FILENAME%"

    echo Building package archive "%package_file_name%" (%RELEASE_TYPE%)...

    call :build_backend
    call :build_frontend

    rmdir package /s /q
    mkdir package
    mkdir package\data
    mkdir package\log
    xcopy ezbookkeeping.exe package\
    xcopy dist package\public /e /i
    xcopy conf package\conf /e /i
    xcopy templates package\templates /e /i
    xcopy LICENSE package\

    cd package

    if !errorlevel! neq 0 (
        call :echo_red "Error: Build Failed"
        goto :end
    )

    call 7z a -r -tzip -mx9 ..\%package_file_name% package *

    cd ..
    endlocal

    goto :eof

:end
    set "TYPE="
    set "NO_LINT="
    set "NO_TEST="
    set "RELEASE="
    set "RELEASE_TYPE="
    set "VERSION="
    set "COMMIT_HASH="
    set "BUILD_UNIXTIME="
    set "BUILD_DATE="
    set "PACKAGE_FILENAME="
    exit /B
