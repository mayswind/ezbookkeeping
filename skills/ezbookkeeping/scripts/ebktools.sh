#!/usr/bin/env sh

# ezBookkeeping API Tools
# A command-line tool for calling ezBookkeeping APIs

SCRIPT_DIR=$(CDPATH= cd "$(dirname "$0")" && pwd)
API_CONFIGS=""

EBKTOOL_SERVER_BASEURL="${EBKTOOL_SERVER_BASEURL}"
EBKTOOL_TOKEN="${EBKTOOL_TOKEN}"
TIMEZONE_NAME=""
TIMEZONE_OFFSET=""
RAW_RESPONSE="false"
DRY_RUN="false"

echo_red() {
    printf '\033[31m%s\033[0m\n' "$1"
}

echo_yellow() {
    printf '\033[33m%s\033[0m\n' "$1"
}

check_dependency() {
    for cmd in $1
    do
        if ! command -v "$cmd" > /dev/null 2>&1; then
            echo_red "Error: \"$cmd\" is required."
            exit 127
        fi
    done
}

load_api_configs() {
    config_file="$SCRIPT_DIR/api-configs.json"

    if [ ! -f "$config_file" ]; then
        echo_red "Error: API configuration file not found: $config_file"
        exit 1
    fi

    if ! API_CONFIGS="$(jq -c '.' "$config_file" 2>/dev/null)"; then
        echo_red "Error: Failed to load API configuration file: $config_file"
        exit 1
    fi
}

load_env_file() {
    env_file="$1"

    if [ ! -f "$env_file" ]; then
        return 1
    fi

    while IFS= read -r line || [ -n "$line" ]; do
        case "$line" in
            ''|'#'*) continue ;;
        esac

        line="$(echo "$line" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')"

        if ! echo "$line" | grep -q '='; then
            continue
        fi

        key="$(echo "$line" | cut -d'=' -f1)"
        value="$(echo "$line" | cut -d'=' -f2-)"
        value="$(echo "$value" | sed -e 's/^["'"'"']//' -e 's/["'"'"']$//')"

        case "$key" in
            EBKTOOL_SERVER_BASEURL)
                EBKTOOL_SERVER_BASEURL="$value"
                ;;
            EBKTOOL_TOKEN)
                EBKTOOL_TOKEN="$value"
                ;;
        esac
    done < "$env_file"

    return 0
}

load_env_from_paths() {
    if [ -n "$EBKTOOL_SERVER_BASEURL" ] && [ -n "$EBKTOOL_TOKEN" ]; then
        return 0
    fi

    current_dir="$(pwd)"
    parent_dir="$(dirname "$current_dir")"
    home_dir="$HOME"

    if [ -z "$EBKTOOL_SERVER_BASEURL" ] || [ -z "$EBKTOOL_TOKEN" ]; then
        if load_env_file "$current_dir/.env"; then
            if [ -n "$EBKTOOL_SERVER_BASEURL" ] && [ -n "$EBKTOOL_TOKEN" ]; then
                return 0
            fi
        fi
    fi

    if [ -z "$EBKTOOL_SERVER_BASEURL" ] || [ -z "$EBKTOOL_TOKEN" ]; then
        if load_env_file "$parent_dir/.env"; then
            if [ -n "$EBKTOOL_SERVER_BASEURL" ] && [ -n "$EBKTOOL_TOKEN" ]; then
                return 0
            fi
        fi
    fi

    if [ -z "$EBKTOOL_SERVER_BASEURL" ] || [ -z "$EBKTOOL_TOKEN" ]; then
        if load_env_file "$home_dir/.env"; then
            if [ -n "$EBKTOOL_SERVER_BASEURL" ] && [ -n "$EBKTOOL_TOKEN" ]; then
                return 0
            fi
        fi
    fi

    return 0
}

url_encode() {
    text="$1"
    printf "%s" "$text" | jq -sRr @uri
}

get_system_timezone_name() {
    if [ -f /etc/timezone ]; then
        cat /etc/timezone
    elif [ -L /etc/localtime ]; then
        readlink /etc/localtime | sed 's|.*/zoneinfo/||'
    elif command -v timedatectl > /dev/null 2>&1; then
        timedatectl | grep "Time zone" | awk '{print $3}'
    fi
}

get_example_timezone_name() {
    tz_name="$(get_system_timezone_name)"

    if [ -n "$tz_name" ]; then
        echo "$tz_name"
    else
        echo "Asia/Shanghai"
    fi
}

get_system_timezone_offset() {
    offset_str="$(date +%z 2>/dev/null)"

    if [ -z "$offset_str" ]; then
        return
    fi

    sign="${offset_str%????}"
    hours="${offset_str#?}"
    hours="${hours%??}"
    minutes="${offset_str#???}"

    hours=$(echo "$hours" | sed 's/^0*//')
    minutes=$(echo "$minutes" | sed 's/^0*//')

    [ -z "$hours" ] && hours=0
    [ -z "$minutes" ] && minutes=0

    if [ "$sign" = "+" ]; then
        echo "$((hours * 60 + minutes))"
    else
        echo "$((-(hours * 60 + minutes)))"
    fi
}

get_example_timezone_offset() {
    tz_offset="$(get_system_timezone_offset)"

    if [ -n "$tz_offset" ]; then
        echo "$tz_offset"
    else
        echo "480"
    fi
}

parse_api_config() {
    api_name="$1"
    echo "$API_CONFIGS" | jq -r --arg name "$api_name" '.[] | select(.Name == $name)'
}

get_pretty_response_config() {
    command_name="$1"
    echo "$API_CONFIGS" | jq -r --arg name "$command_name" '.[] | select(.Name == $name) | .PrettyResponse // null'
}

flatten_hierarchical_data() {
    data="$1"
    child_key="$2"

    printf "%s\n" "$data" | jq -r --arg childKey "$child_key" '
        if type == "array" then
            [.[] | . as $parent | [$parent | del(.[$childKey])] + (.[$childKey] // [])]
        elif type == "object" then
            [.[] | .[] | . as $parent | [$parent | del(.[$childKey])] + (.[$childKey] // [])]
        else
            []
        end | flatten
    '
}

print_markdown_table() {
    data="$1"
    columns="$2"

    if [ -z "$data" ] || [ "$data" = "null" ] || [ "$data" = "[]" ]; then
        echo "No data to display"
        return
    fi

    cols_array="$columns"

    if [ -z "$cols_array" ] || [ "$cols_array" = "null" ]; then
        printf "%s\n" "$data" | jq '.'
        return
    fi

    header="$(echo "$cols_array" | jq -r 'join(" | ")')"
    separator="$(echo "$cols_array" | jq -r '[.[] | "---"] | join(" | ")')"

    rows="$(printf "%s\n" "$data" | jq -r --argjson cols "$cols_array" '
        if type == "array" then
            .[] | [$cols[] as $col | .[$col] | if . == null then "-" elif type == "string" then gsub("\r"; "\\r") | gsub("\n"; "\\n") else tostring end] | join(" | ")
        else
            [$cols[] as $col | .[$col] | if . == null then "-" elif type == "string" then gsub("\r"; "\\r") | gsub("\n"; "\\n") else tostring end] | join(" | ")
        end
    ' 2>/dev/null)"

    if [ -z "$rows" ]; then
        printf "%s\n" "$data" | jq '.'
        return
    fi

    echo "| $header |"
    echo "| $separator |"
    printf "%s\n" "$rows" | while IFS= read -r row; do
        printf "%s\n" "| $row |"
    done
}

print_result() {
    command_name="$1"
    result_data="$2"

    if [ "$RAW_RESPONSE" = "true" ]; then
        printf "%s\n" "$result_data" | jq '.'
        return
    fi

    pretty_config="$(get_pretty_response_config "$command_name")"

    if [ -z "$pretty_config" ] || [ "$pretty_config" = "null" ]; then
        printf "%s\n" "$result_data" | jq '.'
        return
    fi

    display_type="$(echo "$pretty_config" | jq -r '.Type')"
    columns="$(echo "$pretty_config" | jq -c '.Columns')"

    case "$display_type" in
        simple_array_to_markdown_table)
            print_markdown_table "$result_data" "$columns"
            ;;
        hierarchical_array_to_markdown_table)
            child_key="$(echo "$pretty_config" | jq -r '.ChildKey')"
            flattened="$(flatten_hierarchical_data "$result_data" "$child_key")"
            print_markdown_table "$flattened" "$columns"
            ;;
        hierarchical_object_to_markdown_table)
            child_key="$(echo "$pretty_config" | jq -r '.ChildKey')"
            flattened="$(flatten_hierarchical_data "$result_data" "$child_key")"
            print_markdown_table "$flattened" "$columns"
            ;;
        nested_array_to_markdown_table)
            data_path="$(echo "$pretty_config" | jq -r '.DataPath // "."')"
            case "$data_path" in
                .*) ;;
                *) data_path=".$data_path" ;;
            esac
            nested_data="$(printf "%s\n" "$result_data" | jq -r "$data_path")"

            metadata="$(echo "$pretty_config" | jq -r '.Metadata // null')"
            if [ -n "$metadata" ] && [ "$metadata" != "null" ]; then
                echo "$metadata" | jq -r --arg result "$result_data" '
                    ($result | fromjson) as $data |
                    .[] | select($data[.Field] != null) | "\(.Label): \($data[.Field])"
                '
                echo ""
            fi

            print_markdown_table "$nested_data" "$columns"
            ;;
        *)
            printf "%s\n" "$result_data" | jq '.'
            ;;
    esac
}

show_help() {
    example_timezone_name="$(get_example_timezone_name)"
    example_timezone_offset="$(get_example_timezone_offset)"

    cat <<-EOF
ezBookkeeping API Tools

A command-line tool for calling ezBookkeeping APIs

Usage:
    ebktools.sh [--tz-name <name>] [--tz-offset <offset>] [--raw-response] [--dry-run] <command> [command-options]

Environment Variables (Required):
    EBKTOOL_SERVER_BASEURL      ezBookkeeping server base URL (e.g., http://localhost:8080)
    EBKTOOL_TOKEN               ezBookkeeping API token

    You can also set the above environment variables in a '.env' file located in the current directory, parent directory or home directory.

Global Options:
    --tz-name <name>            The IANA timezone name of current timezone. For example, for Beijing Time it is 'Asia/Shanghai'.
    --tz-offset <offset>        The offset in minutes of the current timezone from UTC. For example, for Beijing Time which is UTC+8, the value is '480'. If both '--tz-name' and '--tz-offset' are set, '--tz-offset' takes priority. If neither is set, the current system time zone is used by default.
    --raw-response              Display the response in raw JSON format instead of formatted table.
    --dry-run                   Print the request method, URL, headers, and JSON body without sending it.

Commands:
    list                        List all available API commands
    help <api-command>          Show help for a specific API command
    <api-command>               Execute an API command

Examples:
    # Set environment variables
    export EBKTOOL_SERVER_BASEURL="http://localhost:8080"
    export EBKTOOL_TOKEN="YOUR_TOKEN"

    # List all available commands
    ebktools.sh list

    # Show help for a specific command
    ebktools.sh help server-version

    # Call server-version API
    ebktools.sh server-version

    # Call API with timezone name
    ebktools.sh --tz-name ${example_timezone_name} transactions-list --count 10

    # Call API with timezone offset
    ebktools.sh --tz-offset ${example_timezone_offset} transactions-list --count 10

    # Preview a request without sending it
    ebktools.sh --dry-run transactions-add --type 3 --categoryId 0 --time 1710000000 --utcOffset 480 --sourceAccountId 1 --sourceAmount -1234
EOF
}

list_commands() {
    echo "Available API Commands:"
    echo ""

    name_width="$(echo "$API_CONFIGS" | jq '[.[].Name | length] | max + 2')"

    echo "$API_CONFIGS" | jq -r '.[] | "\(.Name)|\(.Description)"' | while IFS='|' read -r name desc; do
        printf "  %-*s %s\n" "$name_width" "$name" "$desc"
    done

    echo ""
    echo "Use 'ebktools.sh help <api-command>' to see detailed information about an API command."
}

show_command_help() {
    command_name="$1"
    config="$(parse_api_config "$command_name")"

    if [ -z "$config" ]; then
        echo_red "Error: Unknown command '$command_name'"
        echo ""
        echo "Use 'ebktools.sh list' to see all available commands."
        exit 1
    fi

    name="$(echo "$config" | jq -r '.Name')"
    desc="$(echo "$config" | jq -r '.Description')"
    method="$(echo "$config" | jq -r '.Method')"
    path="$(echo "$config" | jq -r '.Path')"
    requires_timezone="$(echo "$config" | jq -r '.RequiresTimezone // false')"
    response_struct="$(echo "$config" | jq -r '.ResponseStructure')"

    echo "Command: $name"
    echo "Description: $desc"
    echo "Method: $method"
    echo "Path: $path"
    echo "Require current time zone: $([ "$requires_timezone" = "true" ] && echo Yes || echo No)"
    echo ""

    required_count="$(echo "$config" | jq '.RequiredParams | length')"
    if [ "$required_count" -gt 0 ]; then
        echo "Required Parameters:"
        echo "$config" | jq -r '.RequiredParams[]' | while read -r param; do
            param_desc="$(echo "$config" | jq -r --arg p "$param" '.ParamDescriptions[$p] // ""')"
            printf "  --%-25s %s\n" "$param" "$param_desc"
        done
        echo ""
    fi

    optional_count="$(echo "$config" | jq '.OptionalParams | length')"
    if [ "$optional_count" -gt 0 ]; then
        echo "Optional Parameters:"
        echo "$config" | jq -r '.OptionalParams[]' | while read -r param; do
            param_desc="$(echo "$config" | jq -r --arg p "$param" '.ParamDescriptions[$p] // ""')"
            printf "  --%-25s %s\n" "$param" "$param_desc"
        done
        echo ""
    fi

    response_struct="$(echo "$config" | jq -r '.ResponseStructure')"

    if [ -n "$response_struct" ] && [ "$response_struct" != "null" ]; then
        echo "Response Structure:"
        echo "$config" | jq -r '.ResponseStructure[]' | while IFS= read -r line; do
            echo "  $line"
        done
        echo ""
    fi

    echo "Example:"
    current_tz_name="$(get_system_timezone_name)"
    current_tz_offset="$(get_system_timezone_offset)"
    if [ "$requires_timezone" = "true" ] && [ -n "$current_tz_name" ]; then
        echo "  ebktools.sh --tz-name ${current_tz_name} $name"
    elif [ "$requires_timezone" = "true" ] && [ -n "$current_tz_offset" ]; then
        echo "  ebktools.sh --tz-offset ${current_tz_offset} $name"
    elif [ "$requires_timezone" = "true" ]; then
        echo "  ebktools.sh --tz-name <name> $name"
    else
        echo "  ebktools.sh $name"
    fi
}

call_api() {
    command_name="$1"
    shift

    config="$(parse_api_config "$command_name")"

    if [ -z "$config" ]; then
        echo_red "Error: Unknown command '$command_name'"
        echo ""
        echo "Use 'ebktools.sh list' to see all available commands."
        exit 1
    fi

    serverBaseUrl="$EBKTOOL_SERVER_BASEURL"
    authToken="$EBKTOOL_TOKEN"

    if [ -z "$serverBaseUrl" ]; then
        if [ "$DRY_RUN" = "true" ]; then
            serverBaseUrl="http://example.local"
        else
            echo_red "Error: Environment variable 'EBKTOOL_SERVER_BASEURL' is not set."
            echo "Please set it to your ezBookkeeping server base URL (e.g., http://localhost:8080)"
            exit 1
        fi
    fi

    if [ -z "$authToken" ]; then
        if [ "$DRY_RUN" = "true" ]; then
            authToken="DRY_RUN_TOKEN"
        else
            echo_red "Error: Environment variable 'EBKTOOL_TOKEN' is not set."
            echo "Please set it to your API token."
            exit 1
        fi
    fi

    requires_timezone="$(echo "$config" | jq -r '.RequiresTimezone // false')"
    current_tz_name="$(get_system_timezone_name)"
    current_tz_offset="$(get_system_timezone_offset)"

    if [ -n "$TIMEZONE_NAME" ]; then
        current_tz_name="$TIMEZONE_NAME"
    fi

    if [ -n "$TIMEZONE_OFFSET" ]; then
        current_tz_name=""
        current_tz_offset="$TIMEZONE_OFFSET"
    fi

    if [ "$requires_timezone" = "true" ] && [ -z "$current_tz_name" ] && [ -z "$current_tz_offset" ]; then
        echo_red "Error: Command '$command_name' requires timezone information."
        echo "Please provide either '--tz-name' or '--tz-offset' parameter."
        echo ""
        echo "Examples:"
        echo "  ebktools.sh --tz-name <name> $command_name ..."
        echo "  ebktools.sh --tz-offset <offset> $command_name ..."
        exit 1
    fi

    method="$(echo "$config" | jq -r '.Method')"
    path="$(echo "$config" | jq -r '.Path')"

    get_param_type() {
        param_name="$1"
        param_type="$(echo "$config" | jq -r --arg p "$param_name" '.ParamTypes[$p] // "string"')"
        echo "$param_type"
    }

    has_param_type() {
        param_name="$1"
        echo "$config" | jq -e --arg p "$param_name" '.ParamTypes | has($p)' >/dev/null 2>&1
    }

    validate_param() {
        param_name="$1"
        param_value="$2"
        param_type="$3"

        case "$param_type" in
            integer)
                if ! echo "$param_value" | grep -Eq '^-?[0-9]+$'; then
                    echo_red "Error: Parameter '--${param_name}' must be an integer value"
                    exit 1
                fi
                ;;
            boolean)
                if ! echo "$param_value" | grep -Eq '^(true|false|1|0)$'; then
                    echo_red "Error: Parameter '--${param_name}' must be a boolean value (true/false or 1/0)"
                    exit 1
                fi
                ;;
            geo_location)
                if ! echo "$param_value" | grep -Eq '^-?[0-9]+(\.[0-9]+)?,-?[0-9]+(\.[0-9]+)?$'; then
                    echo_red "Error: Parameter '--${param_name}' must be in format 'longitude,latitude'"
                    exit 1
                fi
                ;;
            json)
                if ! printf "%s\n" "$param_value" | jq -e '.' >/dev/null 2>&1; then
                    echo_red "Error: Parameter '--${param_name}' must be valid JSON"
                    exit 1
                fi
                ;;
        esac
    }

    params=""
    json_params="{}"
    present_params="|"
    while [ $# -gt 0 ]; do
        case "${1}" in
            --*)
                param_name="$(echo "${1}" | sed 's/^--//')"

                if [ $# -lt 2 ]; then
                    echo_red "Error: Parameter '--${param_name}' requires a value"
                    exit 1
                fi

                param_value="$2"
                if has_param_type "$param_name"; then
                    param_type="$(get_param_type "$param_name")"
                else
                    echo_red "Error: Unknown parameter '--${param_name}'"
                    exit 1
                fi

                if [ "${param_value#--}" != "$param_value" ]; then
                    echo_red "Error: Parameter '--${param_name}' requires a value"
                    exit 1
                fi

                validate_param "$param_name" "$param_value" "$param_type"

                encoded_value="$(url_encode "$param_value")"
                if [ -z "$params" ]; then
                    params="${param_name}=${encoded_value}"
                else
                    params="${params}&${param_name}=${encoded_value}"
                fi
                present_params="${present_params}${param_name}|"

                case "$param_type" in
                    integer)
                        json_params="$(echo "$json_params" | jq --arg k "$param_name" --argjson v "$param_value" '. + {($k): $v}')"
                        ;;
                    boolean)
                        if [ "$param_value" = "true" ] || [ "$param_value" = "1" ]; then
                            bool_value="true"
                        else
                            bool_value="false"
                        fi
                        json_params="$(echo "$json_params" | jq --arg k "$param_name" --argjson v "$bool_value" '. + {($k): $v}')"
                        ;;
                    string_array)
                        json_array="[]"
                        old_ifs="$IFS"
                        IFS=','
                        for val in $param_value; do
                            json_array="$(echo "$json_array" | jq --arg v "$val" '. + [$v]')"
                        done
                        IFS="$old_ifs"
                        json_params="$(echo "$json_params" | jq --arg k "$param_name" --argjson v "$json_array" '. + {($k): $v}')"
                        ;;
                    geo_location)
                        longitude="$(echo "$param_value" | cut -d',' -f1)"
                        latitude="$(echo "$param_value" | cut -d',' -f2)"
                        geo_json="$(jq -n --arg lat "$latitude" --arg lon "$longitude" '{latitude: ($lat | tonumber), longitude: ($lon | tonumber)}')"
                        json_params="$(echo "$json_params" | jq --arg k "$param_name" --argjson v "$geo_json" '. + {($k): $v}')"
                        ;;
                    json)
                        json_params="$(echo "$json_params" | jq --arg k "$param_name" --argjson v "$param_value" '. + {($k): $v}')"
                        ;;
                    *)
                        json_params="$(echo "$json_params" | jq --arg k "$param_name" --arg v "$param_value" '. + {($k): $v}')"
                        ;;
                esac
                shift 2
                ;;
            *)
                echo_red "Error: Invalid parameter: '$1'"
                exit 1
                ;;
        esac
    done

    required_count="$(echo "$config" | jq '.RequiredParams | length')"
    if [ "$required_count" -gt 0 ]; then
        i=0
        while [ "$i" -lt "$required_count" ]; do
            param="$(echo "$config" | jq -r --argjson idx "$i" '.RequiredParams[$idx]')"
            case "$present_params" in
                *"|${param}|"*) ;;
                *)
                    echo_red "Error: Required parameter '--${param}' is missing"
                    exit 1
                    ;;
            esac
            i=$((i + 1))
        done
    fi

    if [ "${serverBaseUrl%/}" != "$serverBaseUrl" ]; then
        serverBaseUrl="${serverBaseUrl%/}"
    fi

    url="${serverBaseUrl}/api/v1/${path}"
    timezone_headers=""

    if [ -n "$current_tz_name" ]; then
        timezone_headers="X-Timezone-Name: $current_tz_name"
    elif [ -n "$current_tz_offset" ]; then
        timezone_headers="X-Timezone-Offset: $current_tz_offset"
    fi

    if [ "$method" = "POST" ]; then
        if [ "$DRY_RUN" = "true" ]; then
            echo_yellow "Dry run: $method $url"
        else
            echo_yellow "Calling API: $method $url"
        fi
        echo ""

        if [ "$DRY_RUN" = "true" ]; then
            echo "Headers:"
            echo "  Authorization: Bearer ***"
            echo "  Content-Type: application/json"
            if [ -n "$timezone_headers" ]; then
                echo "  $timezone_headers"
            fi
            echo ""
            echo "Body:"
            printf "%s\n" "$json_params" | jq '.'
            return
        fi

        if [ "$json_params" != "{}" ]; then
            if [ -n "$timezone_headers" ]; then
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $authToken" \
                    -H "Content-Type: application/json" \
                    -H "$timezone_headers" \
                    -d "$json_params" \
                    "$url")"
                curl_exit_code=$?
            else
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $authToken" \
                    -H "Content-Type: application/json" \
                    -d "$json_params" \
                    "$url")"
                curl_exit_code=$?
            fi
        else
            if [ -n "$timezone_headers" ]; then
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $authToken" \
                    -H "$timezone_headers" \
                    "$url")"
                curl_exit_code=$?
            else
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $authToken" \
                    "$url")"
                curl_exit_code=$?
            fi
        fi
    else
        if [ -n "$params" ]; then
            url="${url}?${params}"
        fi

        if [ "$DRY_RUN" = "true" ]; then
            echo_yellow "Dry run: $method $url"
        else
            echo_yellow "Calling API: $method $url"
        fi
        echo ""

        if [ "$DRY_RUN" = "true" ]; then
            echo "Headers:"
            echo "  Authorization: Bearer ***"
            if [ -n "$timezone_headers" ]; then
                echo "  $timezone_headers"
            fi
            return
        fi

        if [ -n "$timezone_headers" ]; then
            response="$(curl -s -X "$method" \
                -H "Authorization: Bearer $authToken" \
                -H "$timezone_headers" \
                "$url")"
        else
            response="$(curl -s -X "$method" \
                -H "Authorization: Bearer $authToken" \
                "$url")"
        fi
        curl_exit_code=$?
    fi

    if [ "$curl_exit_code" -ne 0 ]; then
        echo_red "Error: API call failed (curl exit code: $curl_exit_code)"
        exit 1
    fi

    success=$(printf "%s\n" "$response" | jq -r '.success // "null"')

    if [ "$success" = "true" ]; then
        echo "Response Result:"
        result=$(printf "%s\n" "$response" | jq '.result // "null"')
        if [ "$result" != "null" ]; then
            print_result "$command_name" "$result"
        else
            echo "Success: true (No result data)"
        fi
    elif [ "$success" = "false" ]; then
        echo "Raw Response:"
        printf "%s\n" "$response" | jq '.'
    else
        echo "Raw Response:"
        printf "%s\n" "$response" | jq '.'
    fi
}

main() {
    check_dependency "grep sed awk date curl jq"
    load_api_configs

    load_env_from_paths

    COMMAND=""

    while [ $# -gt 0 ]; do
        case "${1}" in
            --tz-name)
                if [ $# -lt 2 ]; then
                    echo_red "Error: '--tz-name' requires a value"
                    exit 1
                fi
                TIMEZONE_NAME="$2"
                shift 2
                ;;
            --tz-offset)
                if [ $# -lt 2 ]; then
                    echo_red "Error: '--tz-offset' requires a value"
                    exit 1
                fi
                TIMEZONE_OFFSET="$2"
                shift 2
                ;;
            --raw-response)
                RAW_RESPONSE="true"
                shift
                ;;
            --dry-run)
                DRY_RUN="true"
                shift
                ;;
            --help | -h)
                show_help
                exit 0
                ;;
            list)
                list_commands
                exit 0
                ;;
            help)
                if [ -z "$2" ]; then
                    show_help
                else
                    show_command_help "$2"
                fi
                exit 0
                ;;
            --*)
                echo_red "Error: Unknown option: $1"
                show_help
                exit 1
                ;;
            *)
                COMMAND="$1"
                shift
                break
                ;;
        esac
    done

    if [ -z "$COMMAND" ]; then
        show_help
        exit 0
    fi

    call_api "$COMMAND" "$@"
}

main "$@"
