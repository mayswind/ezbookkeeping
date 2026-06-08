#!/usr/bin/env pwsh

# ezBookkeeping API Tools
# A command-line tool for calling ezBookkeeping APIs

param(
    [Parameter(Position=0)]
    [string]$Command = "",

    [Parameter(Mandatory=$false)]
    [string]$tzName = "",

    [Parameter(Mandatory=$false)]
    [string]$tzOffset = "",

    [Parameter(Mandatory=$false)]
    [switch]$rawResponse = $false,

    [Parameter(Mandatory=$false)]
    [switch]$dryRun = $false,

    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$CommandArgs
)

$script:EBKTOOL_SERVER_BASEURL = $env:EBKTOOL_SERVER_BASEURL
$script:EBKTOOL_TOKEN = $env:EBKTOOL_TOKEN

$script:API_CONFIGS = @()

function Load-ApiConfigs {
    $configPath = Join-Path -Path $PSScriptRoot -ChildPath 'api-configs.json'

    if (-not (Test-Path -LiteralPath $configPath -PathType Leaf)) {
        Write-Red "Error: API configuration file not found: $configPath"
        exit 1
    }

    try {
        $script:API_CONFIGS = Get-Content -Raw -LiteralPath $configPath -ErrorAction Stop | ConvertFrom-Json -AsHashtable -ErrorAction Stop
    } catch {
        Write-Red "Error: Failed to load API configuration file: $($_.Exception.Message)"
        exit 1
    }
}

function Write-Red($msg) {
    Write-Host $msg -ForegroundColor Red
}

function Write-Yellow($msg) {
    Write-Host $msg -ForegroundColor Yellow
}

function Import-DotEnvFile {
    param(
        [string]$Path
    )

    if (-not (Test-Path -Path $Path -PathType Leaf)) {
        return $false
    }

    try {
        Get-Content -Path $Path -ErrorAction Stop | ForEach-Object {
            $line = $_.Trim()

            if ([string]::IsNullOrWhiteSpace($line) -or $line.StartsWith('#')) {
                return
            }

            if ($line -match '^([^=]+)=(.*)$') {
                $key = $matches[1].Trim()
                $value = $matches[2].Trim()

                if ($value -match '^["''](.*)["'']$') {
                    $value = $matches[1]
                }

                if ($key -eq 'EBKTOOL_SERVER_BASEURL' -or $key -eq 'EBKTOOL_TOKEN') {
                    Set-Variable -Name $key -Value $value -Scope Script -Force
                }
            }
        }
        return $true
    } catch {
        return $false
    }
}

function Initialize-EnvironmentVariables {
    if ($script:EBKTOOL_SERVER_BASEURL -and $script:EBKTOOL_TOKEN) {
        return
    }

    $currentDir = Get-Location
    $parentDir = Split-Path -Path $currentDir -Parent
    $homeDir = if ($IsWindows -or $env:OS -match 'Windows') {
        $env:USERPROFILE
    } else {
        $env:HOME
    }

    if (-not $script:EBKTOOL_SERVER_BASEURL -or -not $script:EBKTOOL_TOKEN) {
        $envPath = Join-Path -Path $currentDir -ChildPath '.env'
        if (Import-DotEnvFile -Path $envPath) {
            if ($script:EBKTOOL_SERVER_BASEURL -and $script:EBKTOOL_TOKEN) {
                return
            }
        }
    }

    if (-not $script:EBKTOOL_SERVER_BASEURL -or -not $script:EBKTOOL_TOKEN) {
        if ($parentDir) {
            $envPath = Join-Path -Path $parentDir -ChildPath '.env'
            if (Import-DotEnvFile -Path $envPath) {
                if ($script:EBKTOOL_SERVER_BASEURL -and $script:EBKTOOL_TOKEN) {
                    return
                }
            }
        }
    }

    if (-not $script:EBKTOOL_SERVER_BASEURL -or -not $script:EBKTOOL_TOKEN) {
        if ($homeDir) {
            $envPath = Join-Path -Path $homeDir -ChildPath '.env'
            if (Import-DotEnvFile -Path $envPath) {
                if ($script:EBKTOOL_SERVER_BASEURL -and $script:EBKTOOL_TOKEN) {
                    return
                }
            }
        }
    }
}

function Url-Encode {
    param([string]$text)
    return [System.Uri]::EscapeDataString($text)
}

function Format-ParamValue {
    param($Value)

    if ($Value -is [bool]) {
        return $Value.ToString().ToLower()
    }

    if ($Value -is [Array]) {
        return ($Value -join ",")
    }

    if ($Value -is [Hashtable] -or $Value -is [PSCustomObject]) {
        return ($Value | ConvertTo-Json -Depth 10 -Compress)
    }

    return [string]$Value
}

function Format-Json {
    # Reference: https://jonathancrozier.com/blog/formatting-json-with-proper-indentation-using-powershell
    param(
        [Parameter(Mandatory = $true, ValueFromPipeline = $true)]
        [String]$Json,

        [ValidateRange(1, 1024)]
        [Int]$Indentation = 2
    )

    $indentationLevel = 0
    $insideString = $false
    $previousCharacterWasEscape = $false
    $stringBuilder = New-Object System.Text.StringBuilder
    $characters = $Json.ToCharArray()

    for ($i = 0; $i -lt $characters.Length; $i++) {
        $character = $characters[$i]

        if ($insideString) {
            [void]$stringBuilder.Append($character)

            if ($previousCharacterWasEscape) {
                $previousCharacterWasEscape = $false
            } elseif ($character -eq '\') {
                if ($i + 1 -lt $characters.Length) {
                    $nextCharacter = $characters[$i + 1]

                    if ($nextCharacter -in @('"', '\', '/', 'b', 'f', 'n', 'r', 't', 'u')) {
                        $previousCharacterWasEscape = $true
                    }
                }
            } elseif ($character -eq '"') {
                $insideString = $false
            }
        } else {
            switch ($character) {
                '"'{
                    $insideString = $true
                    [void]$stringBuilder.Append($character)
                }
                '{'{
                    [void]$stringBuilder.Append($character)

                    if ($i + 1 -lt $characters.Length -and $characters[$i + 1] -eq '}') {
                        [void]$stringBuilder.Append('}')
                        $i++
                        continue
                    }

                    $indentationLevel++
                    [void]$stringBuilder.Append("`n" + (' ' * ($indentationLevel * $Indentation)))
                }
                '['{
                    [void]$stringBuilder.Append($character)

                    if ($i + 1 -lt $characters.Length -and $characters[$i + 1] -eq ']') {
                        [void]$stringBuilder.Append(']')
                        $i++
                        continue
                    }

                    $indentationLevel++
                    [void]$stringBuilder.Append("`n" + (' ' * ($indentationLevel * $Indentation)))
                }
                '}'{
                    $indentationLevel--
                    [void]$stringBuilder.Append("`n" + (' ' * ($indentationLevel * $Indentation)) + $character)
                }
                ']'{
                    $indentationLevel--
                    [void]$stringBuilder.Append("`n" + (' ' * ($indentationLevel * $Indentation)) + $character)
                }
                ','{
                    [void]$stringBuilder.Append($character)
                    [void]$stringBuilder.Append("`n" + (' ' * ($indentationLevel * $Indentation)))
                }
                ':'{
                    [void]$stringBuilder.Append(": ")
                }
                default{
                    if (-not [char]::IsWhiteSpace($character)) {
                        [void]$stringBuilder.Append($character)
                    }
                }
            }
        }
    }

    return $stringBuilder.ToString()
}

function Get-SystemTimezoneName {
    try {
        $tz = [System.TimeZoneInfo]::Local
        $windowsId = $tz.Id

        if ($windowsId) {
            try {
                $ianaId = $null
                $result = [System.TimeZoneInfo]::TryConvertWindowsIdToIanaId($windowsId, [ref]$ianaId)
                if ($result -and $ianaId) {
                    return $ianaId
                }
            } catch {
                # Do Nothing
            }

            # Fallback mapping for common Windows timezone IDs
            if ($TIMEZONE_IANA_NAMES.ContainsKey($windowsId)) {
                return $TIMEZONE_IANA_NAMES[$windowsId]
            }
        }
    } catch {
        # Do Nothing
    }
    return $null
}

function Get-ExampleTimezoneName {
    $name = Get-SystemTimezoneName

    if ($null -ne $name) {
        return "$name"
    } else {
        return "Asia/Shanghai"
    }
}

function Get-SystemTimezoneOffset {
    try {
        $offset = [System.TimeZoneInfo]::Local.BaseUtcOffset
        return [string]$offset.TotalMinutes.ToString()
    } catch {
        # Do Nothing
    }
    return $null
}

function Get-ExampleTimezoneOffset {
    $offset = Get-SystemTimezoneOffset

    if ($null -ne $offset) {
        return "$offset"
    } else {
        return "480"
    }
}

function Get-ApiConfig {
    param([string]$apiName)

    foreach ($config in $API_CONFIGS) {
        if ($config.Name -eq $apiName) {
            return $config
        }
    }

    return $null
}

function Get-PrettyResponseConfig {
    param([string]$commandName)

    foreach ($config in $API_CONFIGS) {
        if ($config.Name -eq $commandName) {
            return $config.PrettyResponse
        }
    }

    return $null
}

function Flatten-HierarchicalData {
    param(
        [Parameter(Mandatory=$true)]
        $Data,
        [string]$ChildKey
    )

    $result = @()
    $items = @()

    if ($Data -is [Array]) {
        $items = $Data
    } elseif ($Data -is [PSCustomObject] -or $Data -is [Hashtable]) {
        foreach ($prop in $Data.PSObject.Properties) {
            if ($prop.Value -is [Array]) {
                $items += $prop.Value
            }
        }
    }

    foreach ($item in $items) {
        $parent = @{}
        foreach ($prop in $item.PSObject.Properties) {
            if ($prop.Name -ne $ChildKey) {
                $parent[$prop.Name] = $prop.Value
            }
        }
        $result += [PSCustomObject]$parent

        if ($item.PSObject.Properties[$ChildKey] -and $item.$ChildKey) {
            foreach ($child in $item.$ChildKey) {
                $result += $child
            }
        }
    }

    return $result
}

function Write-Markdown-Table {
    param(
        [Parameter(Mandatory=$true)]
        $Data,
        [string[]]$Columns
    )

    if (-not $Data -or ($Data -is [Array] -and $Data.Count -eq 0)) {
        Write-Host "No data to display"
        return
    }

    if (-not $Columns -or $Columns.Count -eq 0) {
        $Data | ConvertTo-Json -Depth 10 -Compress | Format-Json
        return
    }

    $tableData = @()

    if ($Data -is [Array]) {
        foreach ($item in $Data) {
            $row = [ordered]@{}
            foreach ($col in $Columns) {
                if ($item.PSObject.Properties[$col]) {
                    $value = $item.$col
                    if ($value -is [bool]) {
                        $row[$col] = $value.ToString().ToLower()
                    } elseif ($value -is [string] -and $value -eq "") {
                        $row[$col] = ""
                    } elseif ($value -is [string]) {
                        $row[$col] = $value -replace "`r", "\n" -replace "`n", "\n"
                    } elseif ($value -is [Array] -and $value.Count -eq 0) {
                        $row[$col] = "[]"
                    } elseif ($value -is [Array] -or $value -is [PSCustomObject] -or $value -is [Hashtable]) {
                        $row[$col] = ($value | ConvertTo-Json -Depth 10 -Compress)
                    } elseif ($null -eq $value) {
                        $row[$col] = "-"
                    } else {
                        $row[$col] = $value
                    }
                } else {
                    $row[$col] = "-"
                }
            }
            $tableData += [PSCustomObject]$row
        }
    } else {
        $row = [ordered]@{}
        foreach ($col in $Columns) {
            if ($Data.PSObject.Properties[$col]) {
                $value = $Data.$col
                if ($value -is [bool]) {
                    $row[$col] = $value.ToString().ToLower()
                } elseif ($value -is [string] -and $value -eq "") {
                    $row[$col] = ""
                } elseif ($value -is [string]) {
                    $row[$col] = $value -replace "`r", "\n" -replace "`n", "\n"
                } elseif ($value -is [Array] -and $value.Count -eq 0) {
                    $row[$col] = "[]"
                } elseif ($value -is [Array] -or $value -is [PSCustomObject] -or $value -is [Hashtable]) {
                    $row[$col] = ($value | ConvertTo-Json -Depth 10 -Compress)
                } elseif ($null -eq $value) {
                    $row[$col] = "-"
                } else {
                    $row[$col] = $value
                }
            } else {
                $row[$col] = "-"
            }
        }
        $tableData += [PSCustomObject]$row
    }

    if ($tableData.Count -gt 0) {
        $header = "| " + (($Columns -join " | ")) + " |"
        Write-Host $header

        $separator = "| " + ((1..$Columns.Count | ForEach-Object { "---" }) -join " | ") + " |"
        Write-Host $separator

        foreach ($item in $tableData) {
            $values = @()
            foreach ($col in $Columns) {
                $values += $item.$col
            }
            $row = "| " + (($values -join " | ")) + " |"
            Write-Host $row
        }
    }
}

function Write-Result {
    param(
        [string]$CommandName,
        $ResultData,
        [bool]$RawResponse = $false
    )

    if ($RawResponse) {
        $ResultData | ConvertTo-Json -Depth 10 -Compress | Format-Json
        return
    }

    $prettyConfig = Get-PrettyResponseConfig -commandName $CommandName

    if (-not $prettyConfig) {
        $ResultData | ConvertTo-Json -Depth 10 -Compress | Format-Json
        return
    }

    $displayType = $prettyConfig.Type
    $columns = $prettyConfig.Columns

    switch ($displayType) {
        "simple_array_to_markdown_table" {
            Write-Markdown-Table -Data $ResultData -Columns $columns
        }
        "hierarchical_array_to_markdown_table" {
            $childKey = $prettyConfig.ChildKey
            $flattened = Flatten-HierarchicalData -Data $ResultData -ChildKey $childKey
            Write-Markdown-Table -Data $flattened -Columns $columns
        }
        "hierarchical_object_to_markdown_table" {
            $childKey = $prettyConfig.ChildKey
            $flattened = Flatten-HierarchicalData -Data $ResultData -ChildKey $childKey
            Write-Markdown-Table -Data $flattened -Columns $columns
        }
        "nested_array_to_markdown_table" {
            $dataPath = $prettyConfig.DataPath
            if ($dataPath) {
                if ($dataPath.StartsWith(".")) {
                    $dataPath = $dataPath.Substring(1)
                }
                $nestedData = $ResultData.$dataPath
            } else {
                $nestedData = $ResultData
            }

            if ($prettyConfig.Metadata) {
                foreach ($meta in $prettyConfig.Metadata) {
                    $value = $ResultData.($meta.Field)
                    if ($null -ne $value) {
                        Write-Host "$($meta.Label): $value"
                    }
                }

                Write-Host ""
            }

            Write-Markdown-Table -Data $nestedData -Columns $columns
        }
        default {
            $ResultData | ConvertTo-Json -Depth 10 -Compress | Format-Json
        }
    }
}

function Show-Help {
    $exampleTimezoneName = Get-ExampleTimezoneName
    $exampleTimezoneOffset = Get-ExampleTimezoneOffset

    Write-Host "ezBookkeeping API Tools"
    Write-Host ""
    Write-Host "A command-line tool for calling ezBookkeeping APIs"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "    ebktools.ps1 [-tzName <name>] [-tzOffset <offset>] [-rawResponse] [-dryRun] <command> [command-options]"
    Write-Host ""
    Write-Host "Environment Variables (Required):"
    Write-Host "    EBKTOOL_SERVER_BASEURL      ezBookkeeping server base URL (e.g., http://localhost:8080)"
    Write-Host "    EBKTOOL_TOKEN               ezBookkeeping API token"
    Write-Host ""
    Write-Host "    You can also set the above environment variables in a '.env' file located in the current directory, parent directory or home directory."
    Write-Host ""
    Write-Host "Global Options:"
    Write-Host "    -tzName <name>              The IANA timezone name of current timezone. For example, for Beijing Time it is 'Asia/Shanghai'."
    Write-Host "    -tzOffset <offset>          The offset in minutes of the current timezone from UTC. For example, for Beijing Time which is UTC+8, the value is '480'. If both '-tzName' and '-tzOffset' are set, '-tzOffset' takes priority. If neither is set, the current system time zone is used by default."
    Write-Host "    -rawResponse                Display the response in raw JSON format instead of formatted table."
    Write-Host "    -dryRun                     Print the request method, URL, headers, and JSON body without sending it."
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "    list                        List all available API commands"
    Write-Host "    help <api-command>          Show help for a specific API command"
    Write-Host "    <api-command>               Execute an API command"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "    # Set environment variables"
    Write-Host "    `$env:EBKTOOL_SERVER_BASEURL = 'http://localhost:8080'"
    Write-Host "    `$env:EBKTOOL_TOKEN = 'YOUR_TOKEN'"
    Write-Host ""
    Write-Host "    # List all available commands"
    Write-Host "    ebktools.ps1 list"
    Write-Host ""
    Write-Host "    # Show help for a specific command"
    Write-Host "    ebktools.ps1 help server-version"
    Write-Host ""
    Write-Host "    # Call server-version API"
    Write-Host "    ebktools.ps1 server-version"
    Write-Host ""
    Write-Host "    # Call API with timezone name"
    Write-Host "    ebktools.ps1 -tzName $exampleTimezoneName transactions-list -count 10"
    Write-Host ""
    Write-Host "    # Call API with timezone offset"
    Write-Host "    ebktools.ps1 -tzOffset $exampleTimezoneOffset transactions-list -count 10"
    Write-Host ""
    Write-Host "    # Preview a request without sending it"
    Write-Host "    ebktools.ps1 -dryRun transactions-add -type 3 -categoryId 0 -time 1710000000 -utcOffset 480 -sourceAccountId 1 -sourceAmount -1234"
}

function Show-CommandList {
    Write-Host "Available API Commands:"
    Write-Host ""

    $nameWidth = (($API_CONFIGS | ForEach-Object { $_.Name.Length } | Measure-Object -Maximum).Maximum + 2)

    foreach ($config in $API_CONFIGS) {
        $name = $config.Name.PadRight($nameWidth)
        Write-Host "  $name$($config.Description)"
    }

    Write-Host ""
    Write-Host "Use 'ebktools.ps1 help <api-command>' to see detailed information about an API command."
}

function Show-CommandHelp {
    param([string]$commandName)

    $config = Get-ApiConfig $commandName

    if (-not $config) {
        Write-Red "Error: Unknown command '$commandName'"
        Write-Host ""
        Write-Host "Use 'ebktools.ps1 list' to see all available commands."
        exit 1
    }

    Write-Host "Command: $($config.Name)"
    Write-Host "Description: $($config.Description)"
    Write-Host "Method: $($config.Method)"
    Write-Host "Path: $($config.Path)"
    Write-Host ("Require current time zone: " + ($(if ($config.RequiresTimezone) { 'Yes' } else { 'No' })))
    Write-Host ""

    if ($config.RequiredParams.Count -gt 0) {
        Write-Host "Required Parameters:"
        foreach ($param in $config.RequiredParams) {
            $desc = $config.ParamDescriptions[$param]
            Write-Host "  -$($param.PadRight(26)) $desc"
        }
        Write-Host ""
    }

    if ($config.OptionalParams.Count -gt 0) {
        Write-Host "Optional Parameters:"
        foreach ($param in $config.OptionalParams) {
            $desc = $config.ParamDescriptions[$param]
            Write-Host "  -$($param.PadRight(26)) $desc"
        }
        Write-Host ""
    }

    if ($config.ResponseStructure) {
        Write-Host "Response Structure:"
        foreach ($line in $config.ResponseStructure) {
            Write-Host "  $line"
        }
        Write-Host ""
    }

    Write-Host "Example:"
    $currentTzName = Get-SystemTimezoneName
    $currentTzOffset = Get-SystemTimezoneOffset
    if ($config.RequiresTimezone -and $null -ne $currentTzName) {
        Write-Host "  ebktools.ps1 -tzName $currentTzName $($config.Name)"
    } elseif ($config.RequiresTimezone -and $null -ne $currentTzOffset) {
        Write-Host "  ebktools.ps1 -tzOffset $currentTzOffset $($config.Name)"
    } elseif ($config.RequiresTimezone) {
        Write-Host "  ebktools.ps1 -tzName <name> $($config.Name)"
    } else {
        Write-Host "  ebktools.ps1 $($config.Name)"
    }
}

function Parse-CommandArgs {
    param(
        [string[]]$commandArgs,
        [hashtable]$paramTypes
    )

    $params = @{}
    $i = 0

    while ($i -lt $commandArgs.Count) {
        $arg = $commandArgs[$i]

        if ($arg.StartsWith("-")) {
            $paramName = $arg.Substring(1)

            if ($i + 1 -lt $commandArgs.Count) {
                $paramType = "string"
                $paramValue = $commandArgs[$i + 1]

                if ($paramTypes -and $paramTypes.ContainsKey($paramName)) {
                    $paramType = $paramTypes[$paramName]
                } elseif ($paramTypes) {
                    Write-Red "Error: Unknown parameter '-$paramName'"
                    exit 1
                }

                if ($paramValue.StartsWith("-")) {
                    $possibleParamName = $paramValue.Substring(1)
                    if ($paramTypes -and $paramTypes.ContainsKey($possibleParamName)) {
                        Write-Red "Error: Parameter '-$paramName' requires a value"
                        exit 1
                    }
                }

                try {
                    switch ($paramType) {
                        "integer" {
                            $numValue = [long]::Parse($paramValue)
                            $params[$paramName] = $numValue
                        }
                        "boolean" {
                            if ($paramValue -match "^(true|false|1|0)$") {
                                $params[$paramName] = ($paramValue -eq "true" -or $paramValue -eq "1")
                            } else {
                                Write-Red "Error: Parameter '-$paramName' must be a boolean value (true/false or 1/0)"
                                exit 1
                            }
                        }
                        "string_array" {
                            $arrayValues = $paramValue.Split(",")
                            $params[$paramName] = $arrayValues
                        }
                        "geo_location" {
                            $coords = $paramValue.Split(",")
                            if ($coords.Count -ne 2) {
                                Write-Red "Error: Parameter '-$paramName' must be in format 'longitude,latitude'"
                                exit 1
                            }

                            $longitude = [double]::Parse($coords[0])
                            $latitude = [double]::Parse($coords[1])

                            $geoLocation = @{
                                latitude = $latitude
                                longitude = $longitude
                            }

                            $params[$paramName] = $geoLocation
                        }
                        "json" {
                            $params[$paramName] = $paramValue | ConvertFrom-Json -AsHashtable -NoEnumerate
                        }
                        default {
                            $params[$paramName] = $paramValue
                        }
                    }
                } catch {
                    Write-Red "Error: Parameter '-$paramName' has invalid $paramType value: '$paramValue'"
                    exit 1
                }

                $i += 2
            } else {
                Write-Red "Error: Parameter '-$paramName' requires a value"
                exit 1
            }
        } else {
            Write-Red "Error: Invalid parameter format: '$arg'"
            exit 1
        }
    }

    return $params
}

function Invoke-Api {
    param(
        [string]$commandName,
        [string[]]$commandArgs
    )

    $config = Get-ApiConfig $commandName

    if (-not $config) {
        Write-Red "Error: Unknown command '$commandName'"
        Write-Host ""
        Write-Host "Use 'ebktools.ps1 list' to see all available commands."
        exit 1
    }

    $serverBaseUrl = $script:EBKTOOL_SERVER_BASEURL
    $authToken = $script:EBKTOOL_TOKEN

    if (-not $serverBaseUrl) {
        if ($script:dryRun) {
            $serverBaseUrl = "http://example.local"
        } else {
            Write-Red "Error: Environment variable 'EBKTOOL_SERVER_BASEURL' is not set."
            Write-Host "Please set it to your ezBookkeeping server base URL (e.g., http://localhost:8080)"
            exit 1
        }
    }

    if (-not $authToken) {
        if ($script:dryRun) {
            $authToken = "DRY_RUN_TOKEN"
        } else {
            Write-Red "Error: Environment variable 'EBKTOOL_TOKEN' is not set."
            Write-Host "Please set it to your API token."
            exit 1
        }
    }

    $currentTimezoneName = Get-SystemTimezoneName
    $currentTimezoneOffset = Get-SystemTimezoneOffset

    if ($script:tzName -ne $null -and $script:tzName -ne "") {
        $currentTimezoneName = $script:tzName
    }

    if ($script:tzOffset -ne $null -and $script:tzOffset -ne "") {
        $currentTimezoneName = ""
        $currentTimezoneOffset = $script:tzOffset
    }

    if ($config.RequiresTimezone -and -not $currentTimezoneName -and -not $currentTimezoneOffset) {
        Write-Red "Error: Command '$commandName' requires timezone information."
        Write-Host "Please provide either '-tzName' or '-tzOffset' parameter."
        Write-Host ""
        Write-Host "Examples:"
        Write-Host "  ebktools.ps1 -tzName <name> $commandName ..."
        Write-Host "  ebktools.ps1 -tzOffset <offset> $commandName ..."
        exit 1
    }

    $paramTypes = @{}
    if ($config.ParamTypes) {
        $paramTypes = $config.ParamTypes
    }

    $params = Parse-CommandArgs -commandArgs $commandArgs -paramTypes $paramTypes

    foreach ($requiredParam in $config.RequiredParams) {
        if (-not $params.ContainsKey($requiredParam)) {
            Write-Red "Error: Required parameter '-$requiredParam' is missing"
            exit 1
        }
    }

    if ($serverBaseUrl.EndsWith("/")) {
        $serverBaseUrl = $serverBaseUrl.Substring(0, $serverBaseUrl.Length - 1)
    }

    $url = "$serverBaseUrl/api/v1/$($config.Path)"

    try {
        $headers = @{
            "Authorization" = "Bearer $authToken"
        }

        if ($currentTimezoneName) {
            $headers["X-Timezone-Name"] = $currentTimezoneName
        } elseif ($currentTimezoneOffset) {
            $headers["X-Timezone-Offset"] = $currentTimezoneOffset
        }

        if ($config.Method -eq "POST") {
            $headers["Content-Type"] = "application/json"

            Write-Yellow "$(if ($script:dryRun) { 'Dry run' } else { 'Calling API' }): $($config.Method) $url"
            Write-Host ""

            if ($params.Count -gt 0) {
                $body = ConvertTo-Json -Depth 10 $params
                if ($script:dryRun) {
                    Write-Host "Headers:"
                    Write-Host "  Authorization: Bearer ***"
                    Write-Host "  Content-Type: application/json"
                    if ($headers.ContainsKey("X-Timezone-Name")) {
                        Write-Host "  X-Timezone-Name: $($headers["X-Timezone-Name"])"
                    } elseif ($headers.ContainsKey("X-Timezone-Offset")) {
                        Write-Host "  X-Timezone-Offset: $($headers["X-Timezone-Offset"])"
                    }
                    Write-Host ""
                    Write-Host "Body:"
                    Write-Host ($body | Format-Json)
                    return
                }
                $response = Invoke-WebRequest -Uri $url -Method POST -Headers $headers -Body $body -ErrorAction Stop -UseBasicParsing
            } else {
                if ($script:dryRun) {
                    Write-Host "Headers:"
                    Write-Host "  Authorization: Bearer ***"
                    Write-Host "  Content-Type: application/json"
                    if ($headers.ContainsKey("X-Timezone-Name")) {
                        Write-Host "  X-Timezone-Name: $($headers["X-Timezone-Name"])"
                    } elseif ($headers.ContainsKey("X-Timezone-Offset")) {
                        Write-Host "  X-Timezone-Offset: $($headers["X-Timezone-Offset"])"
                    }
                    Write-Host ""
                    Write-Host "Body:"
                    Write-Host "{}"
                    return
                }
                $response = Invoke-WebRequest -Uri $url -Method POST -Headers $headers -ErrorAction Stop -UseBasicParsing
            }

            $response = ConvertFrom-Json $response.Content
        } else {
            if ($params.Count -gt 0) {
                $queryString = ($params.GetEnumerator() | ForEach-Object { "$($_.Key)=$(Url-Encode (Format-ParamValue $_.Value))" }) -join "&"
                $url = "${url}?$queryString"
            }

            Write-Yellow "$(if ($script:dryRun) { 'Dry run' } else { 'Calling API' }): $($config.Method) $url"
            Write-Host ""

            if ($script:dryRun) {
                Write-Host "Headers:"
                Write-Host "  Authorization: Bearer ***"
                if ($headers.ContainsKey("X-Timezone-Name")) {
                    Write-Host "  X-Timezone-Name: $($headers["X-Timezone-Name"])"
                } elseif ($headers.ContainsKey("X-Timezone-Offset")) {
                    Write-Host "  X-Timezone-Offset: $($headers["X-Timezone-Offset"])"
                }
                return
            }

            $response = Invoke-WebRequest -Uri $url -Method $config.Method -Headers $headers -ErrorAction Stop -UseBasicParsing
            $response = ConvertFrom-Json $response.Content
        }

        if ($response.PSObject.Properties.Name -contains "success") {
            if ($response.success -eq $true) {
                Write-Host "Response Result:"
                if ($response.PSObject.Properties.Name -contains "result") {
                    Write-Result -CommandName $commandName -ResultData $response.result -RawResponse $script:rawResponse
                } else {
                    Write-Host "Success: true (No result data)"
                }
            } else {
                Write-Host "Raw Response:"
                $jsonOutput = ConvertTo-Json -Depth 10 -Compress $response | Format-Json
                Write-Host $jsonOutput
            }
        } else {
            Write-Host "Raw Response:"
            $jsonOutput = ConvertTo-Json -Depth 10 -Compress $response | Format-Json
            Write-Host $jsonOutput
        }
    } catch {
        if ($_.ErrorDetails.Message) {
            Write-Host "Raw Response:"
            try {
                $errorJson = ConvertFrom-Json $_.ErrorDetails.Message
                $formattedJson = ConvertTo-Json -Depth 10 -Compress $errorJson | Format-Json
                Write-Host $formattedJson
            } catch {
                Write-Host $_.ErrorDetails.Message
            }
        } else {
            $exceptionMessage = $_.Exception.Message
            Write-Red "Error: API call failed ($exceptionMessage)"
            exit 1
        }
    }
}

function Main {
    Load-ApiConfigs
    Initialize-EnvironmentVariables

    if ($Command -eq "list") {
        Show-CommandList
        exit 0
    }

    if ($Command -eq "help") {
        if ($CommandArgs.Count -eq 0) {
            Show-Help
        } else {
            Show-CommandHelp $CommandArgs[0]
        }
        exit 0
    }

    if (-not $Command) {
        Show-Help
        exit 0
    }

    Invoke-Api -commandName $Command -commandArgs $CommandArgs
}

Main
