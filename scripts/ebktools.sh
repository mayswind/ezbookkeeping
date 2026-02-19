#!/usr/bin/env sh

# ezBookkeeping API Tools
# A command-line tool for calling ezBookkeeping APIs

# API Configuration Structure
API_CONFIGS='[
    {
    "Name": "tokens-list",
    "Description": "Get available sessions information",
    "Method": "GET",
    "Path": "tokens/list.json",
    "RequiresTimezone": false,
    "RequiredParams": [],
    "OptionalParams": [],
    "ParamTypes": {},
    "ParamDescriptions": {},
    "ResponseStructure": [
      "[",
      "  {",
      "    \"tokenId\": \"string (Token ID)\",",
      "    \"tokenType\": \"integer (Token type, 1: Normal Token, 5: MCP Token, 8: API Token)\",",
      "    \"userAgent\": \"string (The User Agent when the session created)\",",
      "    \"lastSeen\": \"integer (Last refresh unix time of the session)\",",
      "    \"isCurrent\": \"boolean (Whether the session is current)\"",
      "  }",
      "]"
    ],
    "PrettyResponse": {
      "Type": "simple_array_to_markdown_table",
      "Columns": ["tokenId", "tokenType", "userAgent", "lastSeen", "isCurrent"]
    }
  },
  {
    "Name": "tokens-revoke",
    "Description": "Revoke token",
    "Method": "POST",
    "Path": "tokens/revoke.json",
    "RequiresTimezone": false,
    "RequiredParams": ["tokenId"],
    "OptionalParams": [],
    "ParamTypes": {
      "tokenId": "string"
    },
    "ParamDescriptions": {
      "tokenId": "string (Token ID)"
    },
    "ResponseStructure": [
      "boolean (Whether the token is revoked successfully)"
    ]
  },
  {
    "Name": "accounts-list",
    "Description": "Get all accounts list",
    "Method": "GET",
    "Path": "accounts/list.json",
    "RequiresTimezone": false,
    "RequiredParams": [],
    "OptionalParams": [],
    "ParamTypes": {},
    "ParamDescriptions": {},
    "ResponseStructure": [
      "[",
      "  {",
      "    \"id\": \"string (Account ID)\",",
      "    \"name\": \"string (Account name)\",",
      "    \"parentId\": \"string (Parent account ID, 0 for primary account)\",",
      "    \"category\": \"integer (Account category, 1: Cash, 2: Checking Account, 3: Credit Card, 4: Virtual Account, 5: Debt Account, 6: Receivables, 7: Investment Account, 8: Savings Account, 9: Certificate of Deposit)\",",
      "    \"type\": \"integer (Account type, 1: Single Account, 2: Multiple Sub-accounts)\",",
      "    \"icon\": \"string (Account icon ID)\",",
      "    \"color\": \"string (Account icon color, hex color code RRGGBB)\",",
      "    \"currency\": \"string (Account currency code)\",",
      "    \"balance\": \"integer (Account balance, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"')\",",
      "    \"comment\": \"string (Account description)\",",
      "    \"creditCardStatementDate\": \"integer (The statement date of the credit card account)\",",
      "    \"displayOrder\": \"integer (The display order of the account)\",",
      "    \"isAsset\": \"boolean (Whether the account is an asset account)\",",
      "    \"isLiability\": \"boolean (Whether the account is a liability account)\",",
      "    \"hidden\": \"boolean (Whether the account is hidden)\",",
      "    \"subAccounts\": [\"each sub-account object like an account object\"]",
      "  }",
      "]"
    ],
    "PrettyResponse": {
      "Type": "hierarchical_array_to_markdown_table",
      "Columns": ["category", "type", "parentId", "id", "name", "currency", "balance", "hidden", "comment"],
      "ChildKey": "subAccounts"
    }
  },
  {
    "Name": "accounts-add",
    "Description": "Add account",
    "Method": "POST",
    "Path": "accounts/add.json",
    "RequiresTimezone": true,
    "RequiredParams": ["name", "category", "type", "icon", "color", "currency"],
    "OptionalParams": ["balance", "balanceTime", "comment", "creditCardStatementDate"],
    "ParamTypes": {
      "name": "string",
      "category": "integer",
      "type": "integer",
      "icon": "string",
      "color": "string",
      "currency": "string",
      "balance": "integer",
      "balanceTime": "integer",
      "comment": "string",
      "creditCardStatementDate": "integer"
    },
    "ParamDescriptions": {
      "name": "string (Account name)",
      "category": "integer (Account category, 1: Cash, 2: Checking Account, 3: Credit Card, 4: Virtual Account, 5: Debt Account, 6: Receivables, 7: Investment Account, 8: Savings Account, 9: Certificate of Deposit)",
      "type": "integer (Account type, 1: Single Account, 2: Multiple Sub-accounts)",
      "icon": "string (Account icon ID)",
      "color": "string (Account icon color, hex color code RRGGBB)",
      "currency": "string (Account currency code, ISO 4217 code, '"'"'---'"'"' for the parent account)",
      "balance": "integer (Account balance, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"'. Liability account should set to negative amount)",
      "balanceTime": "integer (The unix time when the account balance is the set value. This field is required when balance is set)",
      "comment": "string (Account description)",
      "creditCardStatementDate": "integer (The statement date of the credit card account)"
    },
    "ResponseStructure": [
      "{",
      "  \"id\": \"string (Account ID)\",",
      "  \"name\": \"string (Account name)\",",
      "  \"parentId\": \"string (Parent account ID)\",",
      "  \"category\": \"integer (Account category)\",",
      "  \"type\": \"integer (Account type)\",",
      "  \"icon\": \"string (Account icon ID)\",",
      "  \"color\": \"string (Account icon color)\",",
      "  \"currency\": \"string (Account currency code)\",",
      "  \"balance\": \"integer (Account balance)\",",
      "  \"comment\": \"string (Account description)\",",
      "  \"creditCardStatementDate\": \"integer (The statement date of the credit card account)\",",
      "  \"displayOrder\": \"integer (The display order of the account)\",",
      "  \"isAsset\": \"boolean (Whether the account is an asset account)\",",
      "  \"isLiability\": \"boolean (Whether the account is a liability account)\",",
      "  \"hidden\": \"boolean (Whether the account is hidden)\",",
      "  \"subAccounts\": [\"every sub-account object like account object\"]",
      "}"
    ]
  },
  {
    "Name": "transaction-categories-list",
    "Description": "Get all transaction categories",
    "Method": "GET",
    "Path": "transaction/categories/list.json",
    "RequiresTimezone": false,
    "RequiredParams": [],
    "OptionalParams": [],
    "ParamTypes": {},
    "ParamDescriptions": {},
    "ResponseStructure": [
      "{",
      "  \"transaction category type (1: Income, 2: Expense, 3:Transfer)\": [",
      "    {",
      "      \"id\": \"string (Transaction category ID)\",",
      "      \"name\": \"string (Transaction category name)\",",
      "      \"parentId\": \"string (Parent transaction category ID, 0 for primary category)\",",
      "      \"type\": \"integer (Transaction category type, 1: Income, 2: Expense, 3: Transfer)\",",
      "      \"icon\": \"string (Transaction category icon ID)\",",
      "      \"color\": \"string (Transaction category icon color, hex color code RRGGBB)\",",
      "      \"comment\": \"string (Transaction category description)\",",
      "      \"displayOrder\": \"integer (The display order of the transaction category)\",",
      "      \"hidden\": \"boolean (Whether the transaction category is hidden)\",",
      "      \"subCategories\": [\"each sub-category object like a transaction category object\"]",
      "    }",
      "  ]",
      "}"
    ],
    "PrettyResponse": {
      "Type": "hierarchical_object_to_markdown_table",
      "Columns": ["type", "parentId", "id", "name", "hidden", "comment"],
      "ChildKey": "subCategories"
    }
  },
  {
    "Name": "transaction-categories-add",
    "Description": "Add transaction category",
    "Method": "POST",
    "Path": "transaction/categories/add.json",
    "RequiresTimezone": false,
    "RequiredParams": ["name", "type", "icon", "color"],
    "OptionalParams": ["parentId", "comment"],
    "ParamTypes": {
      "name": "string",
      "type": "integer",
      "parentId": "string",
      "icon": "string",
      "color": "string",
      "comment": "string"
    },
    "ParamDescriptions": {
      "name": "string (Transaction category name)",
      "type": "integer (Transaction category type, 1: Income, 2: Expense, 3: Transfer)",
      "parentId": "string (Parent transaction category ID, 0 for primary category)",
      "icon": "string (Transaction category icon ID)",
      "color": "string (Transaction category icon color, hex color code RRGGBB)",
      "comment": "string (Transaction category description)"
    },
    "ResponseStructure": [
      "{",
      "  \"id\": \"string (Transaction category ID)\",",
      "  \"name\": \"string (Transaction category name)\",",
      "  \"parentId\": \"string (Parent transaction category ID)\",",
      "  \"type\": \"integer (Transaction category type)\",",
      "  \"icon\": \"string (Transaction category icon ID)\",",
      "  \"color\": \"string (Transaction category icon color)\",",
      "  \"comment\": \"string (Transaction category description)\",",
      "  \"displayOrder\": \"integer (The display order of the transaction category)\",",
      "  \"hidden\": \"boolean (Whether the transaction category is hidden)\",",
      "  \"subCategories\": [\"each sub-category object like a transaction category object\"]",
      "}"
    ]
  },
  {
    "Name": "transaction-tags-list",
    "Description": "Get all transaction tags list",
    "Method": "GET",
    "Path": "transaction/tags/list.json",
    "RequiresTimezone": false,
    "RequiredParams": [],
    "OptionalParams": [],
    "ParamTypes": {},
    "ParamDescriptions": {},
    "ResponseStructure": [
      "[",
      "  {",
      "    \"id\": \"string (Transaction tag ID)\",",
      "    \"name\": \"string (Transaction tag name)\",",
      "    \"groupId\": \"string (Transaction tag group ID)\",",
      "    \"displayOrder\": \"integer (The display order of the transaction tag)\",",
      "    \"hidden\": \"boolean (Whether the transaction tag is hidden)\"",
      "  }",
      "]"
    ],
    "PrettyResponse": {
      "Type": "simple_array_to_markdown_table",
      "Columns": ["groupId", "id", "name", "hidden"]
    }
  },
  {
    "Name": "transaction-tags-add",
    "Description": "Add transaction tag",
    "Method": "POST",
    "Path": "transaction/tags/add.json",
    "RequiresTimezone": false,
    "RequiredParams": ["name"],
    "OptionalParams": ["groupId"],
    "ParamTypes": {
      "name": "string",
      "groupId": "string"
    },
    "ParamDescriptions": {
      "name": "string (Transaction tag name)",
      "groupId": "string (Transaction tag group ID, 0 means default group)"
    },
    "ResponseStructure": [
      "{",
      "  \"id\": \"string (Transaction tag ID)\",",
      "  \"name\": \"string (Transaction tag name)\",",
      "  \"groupId\": \"string (Transaction tag group ID)\",",
      "  \"displayOrder\": \"integer (The display order of the transaction tag)\",",
      "  \"hidden\": \"boolean (Whether the transaction tag is hidden)\"",
      "}"
    ]
  },
  {
    "Name": "transactions-list",
    "Description": "Get transactions list with pagination",
    "Method": "GET",
    "Path": "transactions/list.json",
    "RequiresTimezone": true,
    "RequiredParams": ["count"],
    "OptionalParams": ["type", "category_ids", "account_ids", "tag_filter", "amount_filter", "keyword", "max_time", "min_time", "page", "with_count", "with_pictures", "trim_account", "trim_category", "trim_tag"],
    "ParamTypes": {
      "count": "integer",
      "type": "integer",
      "category_ids": "string",
      "account_ids": "string",
      "tag_filter": "string",
      "amount_filter": "string",
      "keyword": "string",
      "max_time": "integer",
      "min_time": "integer",
      "page": "integer",
      "with_count": "boolean",
      "with_pictures": "boolean",
      "trim_account": "boolean",
      "trim_category": "boolean",
      "trim_tag": "boolean"
    },
    "ParamDescriptions": {
      "count": "integer (The count of transactions per page, maximum is 50)",
      "type": "integer (Filter transaction by type, 1: Balance modification, 2: Income, 3: Expense, 4: Transfer)",
      "category_ids": "string (Filter by category IDs, separated by comma)",
      "account_ids": "string (Filter by account IDs, separated by comma)",
      "tag_filter": "string (Filter by tags)",
      "amount_filter": "string (Filter by amount)",
      "keyword": "string (Filter by keyword)",
      "max_time": "integer (The maximum time sequence ID, Set to 0 for latest)",
      "min_time": "integer (The minimum time sequence ID)",
      "page": "integer (Specified page integer)",
      "with_count": "boolean (Whether to get total count)",
      "with_pictures": "boolean (Whether to get picture IDs)",
      "trim_account": "boolean (Whether to get account ID only)",
      "trim_category": "boolean (Whether to get category ID only)",
      "trim_tag": "boolean (Whether to get tag IDs only)"
    },
    "ResponseStructure": [
      "{",
      "  \"items\": [",
      "    {",
      "      \"id\": \"string (Transaction ID)\",",
      "      \"timeSequenceId\": \"string (Transaction time sequence ID)\",",
      "      \"type\": \"integer (Transaction type)\",",
      "      \"categoryId\": \"string (Transaction category ID)\",",
      "      \"category\": \"object (Transaction category object)\",",
      "      \"time\": \"integer (Transaction unix time)\",",
      "      \"utcOffset\": \"integer (Transaction time zone offset minutes)\",",
      "      \"sourceAccountId\": \"string (Source account ID)\",",
      "      \"sourceAccount\": \"object (Source account object)\",",
      "      \"destinationAccountId\": \"string (Destination account ID)\",",
      "      \"destinationAccount\": \"object (Destination account object)\",",
      "      \"sourceAmount\": \"integer (Source amount, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"')\",",
      "      \"destinationAmount\": \"integer (Destination amount, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"')\",",
      "      \"hideAmount\": \"boolean (Whether to hide the amount)\",",
      "      \"tagIds\": [\"each string representing a transaction tag ID\"],",
      "      \"tags\": [\"each object representing a transaction tag object\"],",
      "      \"pictures\": [\"each object representing a transaction picture object\"],",
      "      \"comment\": \"string (Transaction description)\",",
      "      \"geoLocation\": \"object (Transaction geographic location)\",",
      "      \"editable\": \"boolean (Whether the transaction is editable)\"",
      "    }",
      "  ],",
      "  \"nextTimeSequenceId\": \"integer (The next cursor '"'"'max_time'"'"' parameter when requesting older data)\",",
      "  \"totalCount\": \"integer (The total count of transactions)\"",
      "}"
    ],
    "PrettyResponse": {
      "Type": "nested_array_to_markdown_table",
      "Columns": ["id", "type", "time", "utcOffset", "categoryId", "sourceAccountId", "sourceAmount", "destinationAccountId", "destinationAmount", "tagIds", "geoLocation", "comment"],
      "DataPath": ".items",
      "Metadata": [
        {"Field": "totalCount", "Label": "Total Count"},
        {"Field": "nextTimeSequenceId", "Label": "Next Time Sequence ID"}
      ]
    }
  },
  {
    "Name": "transactions-list-all",
    "Description": "Get all transactions list",
    "Method": "GET",
    "Path": "transactions/list/all.json",
    "RequiresTimezone": true,
    "RequiredParams": [],
    "OptionalParams": ["type", "category_ids", "account_ids", "tag_filter", "amount_filter", "keyword", "start_time", "end_time", "with_pictures", "trim_account", "trim_category", "trim_tag"],
    "ParamTypes": {
      "type": "integer",
      "category_ids": "string",
      "account_ids": "string",
      "tag_filter": "string",
      "amount_filter": "string",
      "keyword": "string",
      "start_time": "integer",
      "end_time": "integer",
      "with_pictures": "boolean",
      "trim_account": "boolean",
      "trim_category": "boolean",
      "trim_tag": "boolean"
    },
    "ParamDescriptions": {
      "type": "integer (Filter transaction by type, 1: Balance modification, 2: Income, 3: Expense, 4: Transfer)",
      "category_ids": "string (Filter by category IDs, separated by comma)",
      "account_ids": "string (Filter by account IDs, separated by comma)",
      "tag_filter": "string (Filter by tags)",
      "amount_filter": "string (Filter by amount)",
      "keyword": "string (Filter by keyword)",
      "start_time": "integer (Transaction list start unix time)",
      "end_time": "integer (Transaction list end unix time)",
      "with_pictures": "boolean (Whether to get picture IDs)",
      "trim_account": "boolean (Whether to get account ID only)",
      "trim_category": "boolean (Whether to get category ID only)",
      "trim_tag": "boolean (Whether to get tag IDs only)"
    },
    "ResponseStructure": [
      "[",
      "  {",
      "    \"id\": \"string (Transaction ID)\",",
      "    \"timeSequenceId\": \"string (Transaction time sequence ID)\",",
      "    \"type\": \"integer (Transaction type)\",",
      "    \"categoryId\": \"string (Transaction category ID)\",",
      "    \"category\": \"object (Transaction category object)\",",
      "    \"time\": \"integer (Transaction unix time)\",",
      "    \"utcOffset\": \"integer (Transaction time zone offset minutes)\",",
      "    \"sourceAccountId\": \"string (Source account ID)\",",
      "    \"sourceAccount\": \"object (Source account object)\",",
      "    \"destinationAccountId\": \"string (Destination account ID)\",",
      "    \"destinationAccount\": \"object (Destination account object)\",",
      "    \"sourceAmount\": \"integer (Source amount, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"')\",",
      "    \"destinationAmount\": \"integer (Destination amount, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"')\",",
      "    \"hideAmount\": \"boolean (Whether to hide the amount)\",",
      "    \"tagIds\": [\"each string representing a transaction tag ID\"],",
      "    \"tags\": [\"each object representing a transaction tag object\"],",
      "    \"pictures\": [\"each object representing a transaction picture object\"],",
      "    \"comment\": \"string (Transaction description)\",",
      "    \"geoLocation\": \"object (Transaction geographic location)\",",
      "    \"editable\": \"boolean (Whether the transaction is editable)\"",
      "  }",
      "]"
    ],
    "PrettyResponse": {
      "Type": "simple_array_to_markdown_table",
      "Columns": ["id", "type", "time", "utcOffset", "categoryId", "sourceAccountId", "sourceAmount", "destinationAccountId", "destinationAmount", "tagIds", "geoLocation", "comment"]
    }
  },
  {
    "Name": "transactions-add",
    "Description": "Add transaction",
    "Method": "POST",
    "Path": "transactions/add.json",
    "RequiresTimezone": true,
    "RequiredParams": ["type", "categoryId", "time", "utcOffset", "sourceAccountId", "sourceAmount"],
    "OptionalParams": ["destinationAccountId", "destinationAmount", "hideAmount", "tagIds", "pictureIds", "comment", "geoLocation"],
    "ParamTypes": {
      "type": "integer",
      "categoryId": "string",
      "time": "integer",
      "utcOffset": "integer",
      "sourceAccountId": "string",
      "sourceAmount": "integer",
      "destinationAccountId": "string",
      "destinationAmount": "integer",
      "hideAmount": "boolean",
      "tagIds": "string_array",
      "pictureIds": "string_array",
      "comment": "string",
      "geoLocation": "geo_location"
    },
    "ParamDescriptions": {
      "type": "integer (Transaction type, 1: Balance Modification, 2: Income, 3: Expense, 4: Transfer)",
      "categoryId": "string (Transaction category ID, supports secondary category)",
      "time": "integer (Transaction unix time)",
      "utcOffset": "integer (Transaction time zone offset minutes)",
      "sourceAccountId": "string (Source account ID, supports account without sub-accounts or sub-account)",
      "sourceAmount": "integer (Source amount, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"')",
      "destinationAccountId": "string (Destination account ID, supports account without sub-accounts or sub-account)",
      "destinationAmount": "integer (Destination amount, supports up to two decimals. For example, a value of '"'"'1234'"'"' represents an amount of '"'"'12.34'"'"')",
      "hideAmount": "boolean (Whether to hide amount)",
      "tagIds": "string (Transaction tag IDs, separated by comma, e.g. '"'"'tagid1,tagid2'"'"')",
      "pictureIds": "string (Transaction picture IDs, separated by comma, e.g. '"'"'picid1,picid2'"'"')",
      "comment": "string (Transaction description)",
      "geoLocation": "string (Transaction geographic location, format: longitude,latitude, e.g. '"'"'116.33,39.93'"'"')"
    },
    "ResponseStructure": [
      "{",
      "  \"id\": \"string (Transaction ID)\",",
      "  \"timeSequenceId\": \"string (Transaction time sequence ID)\",",
      "  \"type\": \"integer (Transaction type)\",",
      "  \"categoryId\": \"string (Transaction category ID)\",",
      "  \"category\": \"object (Transaction category object)\",",
      "  \"time\": \"integer (Transaction unix time)\",",
      "  \"utcOffset\": \"integer (Transaction time zone offset minutes)\",",
      "  \"sourceAccountId\": \"string (Source account ID)\",",
      "  \"sourceAccount\": \"object (Source account object)\",",
      "  \"destinationAccountId\": \"string (Destination account ID)\",",
      "  \"destinationAccount\": \"object (Destination account object)\",",
      "  \"sourceAmount\": \"integer (Source amount)\",",
      "  \"destinationAmount\": \"integer (Destination amount)\",",
      "  \"hideAmount\": \"boolean (Whether to hide the amount)\",",
      "  \"tagIds\": [\"each string representing a transaction tag ID\"],",
      "  \"tags\": [\"each object representing a transaction tag object\"],",
      "  \"pictures\": [\"each object representing a transaction picture object\"],",
      "  \"comment\": \"string (Transaction description)\",",
      "  \"geoLocation\": \"object (Transaction geographic location)\",",
      "  \"editable\": \"boolean (Whether the transaction is editable)\"",
      "}"
    ]
  },
  {
    "Name": "exchangerates-latest",
    "Description": "Get latest exchange rates",
    "Method": "GET",
    "Path": "exchange_rates/latest.json",
    "RequiresTimezone": false,
    "RequiredParams": [],
    "OptionalParams": [],
    "ParamTypes": {},
    "ParamDescriptions": {},
    "ResponseStructure": [
      "{",
      "  \"dataSource\": \"string (Exchange rate data source name)\",",
      "  \"referenceUrl\": \"string (Exchange rate data reference URL)\",",
      "  \"updateTime\": \"integer (Exchange rate data update unix time)\",",
      "  \"baseCurrency\": \"string (Base currency code)\",",
      "  \"exchangeRates\": [",
      "    {",
      "      \"currency\": \"string (Currency code)\",",
      "      \"rate\": \"string (Exchange rate, 1 unit of base currency equals to how many units of this currency)\"",
      "    }",
      "  ]",
      "}"
    ],
    "PrettyResponse": {
      "Type": "nested_array_to_markdown_table",
      "Columns": ["currency", "rate"],
      "DataPath": ".exchangeRates",
      "Metadata": [
        {"Field": "dataSource", "Label": "Data Source"},
        {"Field": "baseCurrency", "Label": "Base Currency"},
        {"Field": "updateTime", "Label": "Update Time"}
      ]
    }
  },
  {
    "Name": "server-version",
    "Description": "Get ezBookkeeping server version information",
    "Method": "GET",
    "Path": "systems/version.json",
    "RequiresTimezone": false,
    "RequiredParams": [],
    "OptionalParams": [],
    "ParamTypes": {},
    "ParamDescriptions": {},
    "ResponseStructure": [
      "{",
      "  \"version\": \"string (Server version)\",",
      "  \"commitHash\": \"string (Git commit hash)\"",
      "}"
    ]
  }
]'

EBKTOOL_SERVER_BASEURL="${EBKTOOL_SERVER_BASEURL}"
EBKTOOL_TOKEN="${EBKTOOL_TOKEN}"
TIMEZONE_NAME=""
TIMEZONE_OFFSET=""
RAW_RESPONSE="false"

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
            EBKTOOL_SERVER_BASEURL|EBKTOOL_TOKEN)
                eval "$key=\"\$value\""
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
    ebktools.sh [--tz-name <name>] [--tz-offset <offset>] [--raw-response] <command> [command-options]

Environment Variables (Required):
    EBKTOOL_SERVER_BASEURL      ezBookkeeping server base URL (e.g., http://localhost:8080)
    EBKTOOL_TOKEN               ezBookkeeping API token

    You can also set the above environment variables in a .env file located in the current directory, parent directory or home directory.

Global Options:
    --tz-name <name>            The IANA timezone name of current timezone. For example, for Beijing Time it is 'Asia/Shanghai'.
    --tz-offset <offset>        The offset in minutes of the current timezone from UTC. For example, for Beijing Time which is UTC+8, the value is '480'. If both '--tz-name' and '--tz-offset' are set, '--tz-name' takes priority. If neither is set, the current system time zone is used by default.
    --raw-response              Display the response in raw JSON format instead of formatted table.

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
EOF
}

list_commands() {
    echo "Available API Commands:"
    echo ""

    echo "$API_CONFIGS" | jq -r '.[] | "\(.Name)|\(.Description)"' | while IFS='|' read -r name desc; do
        printf "  %-30s %s\n" "$name" "$desc"
    done

    echo ""
    echo "Use 'ebktools.sh help <api-command>' to see detailed information about a API command."
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
        echo_red "Error: Environment variable 'EBKTOOL_SERVER_BASEURL' is not set."
        echo "Please set it to your ezBookkeeping server base URL (e.g., http://localhost:8080)"
        exit 1
    fi

    if [ -z "$authToken" ]; then
        echo_red "Error: Environment variable 'EBKTOOL_TOKEN' is not set."
        echo "Please set it to your API token."
        exit 1
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

    validate_param() {
        param_name="$1"
        param_value="$2"
        param_type="$3"

        case "$param_type" in
            integer)
                if ! echo "$param_value" | grep -Eq '^-?[0-9]+$'; then
                    echo_red "Error: Parameter '--${param_name}' must be a integer value"
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
        esac
    }

    params=""
    json_params="{}"
    while [ $# -gt 0 ]; do
        case "${1}" in
            --*)
                param_name="$(echo "${1}" | sed 's/^--//')"

                if [ $# -lt 2 ]; then
                    echo_red "Error: Parameter '--${param_name}' requires a value"
                    exit 1
                fi

                case "$2" in
                    --*)
                        echo_red "Error: Parameter '--${param_name}' requires a value"
                        exit 1
                        ;;
                esac

                param_value="$2"
                param_type="$(get_param_type "$param_name")"

                validate_param "$param_name" "$param_value" "$param_type"

                encoded_value="$(url_encode "$param_value")"
                if [ -z "$params" ]; then
                    params="${param_name}=${encoded_value}"
                else
                    params="${params}&${param_name}=${encoded_value}"
                fi

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
            if ! echo "$params" | grep -q "${param}="; then
                echo_red "Error: Required parameter '--${param}' is missing"
                exit 1
            fi
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
        echo_yellow "Calling API: $method $url"
        echo ""

        if [ "$json_params" != "{}" ]; then
            if [ -n "$timezone_headers" ]; then
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $EBKTOOL_TOKEN" \
                    -H "Content-Type: application/json" \
                    -H "$timezone_headers" \
                    -d "$json_params" \
                    "$url")"
                curl_exit_code=$?
            else
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $EBKTOOL_TOKEN" \
                    -H "Content-Type: application/json" \
                    -d "$json_params" \
                    "$url")"
                curl_exit_code=$?
            fi
        else
            if [ -n "$timezone_headers" ]; then
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $EBKTOOL_TOKEN" \
                    -H "$timezone_headers" \
                    "$url")"
                curl_exit_code=$?
            else
                response="$(curl -s -X "POST" \
                    -H "Authorization: Bearer $EBKTOOL_TOKEN" \
                    "$url")"
                curl_exit_code=$?
            fi
        fi
    else
        if [ -n "$params" ]; then
            url="${url}?${params}"
        fi

        echo_yellow "Calling API: $method $url"
        echo ""

        if [ -n "$timezone_headers" ]; then
            response="$(curl -s -X "$method" \
                -H "Authorization: Bearer $EBKTOOL_TOKEN" \
                -H "$timezone_headers" \
                "$url")"
        else
            response="$(curl -s -X "$method" \
                -H "Authorization: Bearer $EBKTOOL_TOKEN" \
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
