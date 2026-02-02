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

    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$CommandArgs
)

# API Configuration Structure
$API_CONFIGS = @(
    @{
        Name = "tokens-list"
        Description = "Get available sessions information"
        Method = "GET"
        Path = "tokens/list.json"
        RequiresTimezone = $false
        RequiredParams = @()
        OptionalParams = @()
        ParamTypes = @{}
        ParamDescriptions = @{}
        ResponseStructure = @(
            "["
            "  {"
            "    `"tokenId`": `"string (Token ID)`","
            "    `"tokenType`": `"integer (Token type, 1: Normal Token, 5: MCP Token, 8: API Token)`","
            "    `"userAgent`": `"string (The User Agent when the session created)`","
            "    `"lastSeen`": `"integer (Last refresh unix time of the session)`","
            "    `"isCurrent`": `"boolean (Whether the session is current)`""
            "  }"
            "]"
        )
    }
    @{
        Name = "tokens-revoke"
        Description = "Revoke token"
        Method = "POST"
        Path = "tokens/revoke.json"
        RequiresTimezone = $false
        RequiredParams = @("tokenId")
        OptionalParams = @()
        ParamTypes = @{
            "tokenId" = "string"
        }
        ParamDescriptions = @{
            "tokenId" = "string (Token ID)"
        }
        ResponseStructure = @(
            "boolean (Whether the token is revoked successfully)"
        )
    }
    @{
        Name = "accounts-list"
        Description = "Get all accounts list"
        Method = "GET"
        Path = "accounts/list.json"
        RequiresTimezone = $false
        RequiredParams = @()
        OptionalParams = @()
        ParamTypes = @{}
        ParamDescriptions = @{}
        ResponseStructure = @(
            "["
            "  {"
            "    `"id`": `"string (Account ID)`","
            "    `"name`": `"string (Account name)`","
            "    `"parentId`": `"string (Parent account ID)`","
            "    `"category`": `"integer (Account category, 1: Cash, 2: Checking Account, 3: Credit Card, 4: Virtual Account, 5: Debt Account, 6: Receivables, 7: Investment Account, 8: Savings Account, 9: Certificate of Deposit)`","
            "    `"type`": `"integer (Account type, 1: Single Account, 2: Multiple Sub-accounts)`","
            "    `"icon`": `"string (Account icon ID)`","
            "    `"color`": `"string (Account icon color, hex color code RRGGBB)`","
            "    `"currency`": `"string (Account currency code)`","
            "    `"balance`": `"integer (Account balance, supports up to two decimals. For example, a value of '1234' represents an amount of '12.34')`","
            "    `"comment`": `"string (Account description)`","
            "    `"creditCardStatementDate`": `"integer (The statement date of the credit card account)`","
            "    `"displayOrder`": `"integer (The display order of the account)`","
            "    `"isAsset`": `"boolean (Whether the account is an asset account)`","
            "    `"isLiability`": `"boolean (Whether the account is a liability account)`","
            "    `"hidden`": `"boolean (Whether the account is hidden)`","
            "    `"subAccounts`": [`"each sub-account object like an account object`"]"
            "  }"
            "]"
        )
    }
    @{
        Name = "accounts-add"
        Description = "Add account"
        Method = "POST"
        Path = "accounts/add.json"
        RequiresTimezone = $true
        RequiredParams = @("name", "category", "type", "icon", "color", "currency")
        OptionalParams = @("balance", "balanceTime", "comment", "creditCardStatementDate")
        ParamTypes = @{
            "name" = "string"
            "category" = "integer"
            "type" = "integer"
            "icon" = "string"
            "color" = "string"
            "currency" = "string"
            "balance" = "integer"
            "balanceTime" = "integer"
            "comment" = "string"
            "creditCardStatementDate" = "integer"
        }
        ParamDescriptions = @{
            "name" = "string (Account name)"
            "category" = "integer (Account category, 1: Cash, 2: Checking Account, 3: Credit Card, 4: Virtual Account, 5: Debt Account, 6: Receivables, 7: Investment Account, 8: Savings Account, 9: Certificate of Deposit)"
            "type" = "integer (Account type, 1: Single Account, 2: Multiple Sub-accounts)"
            "icon" = "string (Account icon ID)"
            "color" = "string (Account icon color, hex color code RRGGBB)"
            "currency" = "string (Account currency code, ISO 4217 code, `"---`" for the parent account)"
            "balance" = "integer (Account balance, supports up to two decimals. For example, a value of `"1234`" represents an amount of `"12.34`". Liability account should set to negative amount)"
            "balanceTime" = "integer (The unix time when the account balance is the set value. This field is required when balance is set)"
            "comment" = "string (Account description)"
            "creditCardStatementDate" = "integer (The statement date of the credit card account)"
        }
        ResponseStructure = @(
            "{"
            "  `"id`": `"string (Account ID)`","
            "  `"name`": `"string (Account name)`","
            "  `"parentId`": `"string (Parent account ID)`","
            "  `"category`": `"integer (Account category)`","
            "  `"type`": `"integer (Account type)`","
            "  `"icon`": `"string (Account icon ID)`","
            "  `"color`": `"string (Account icon color)`","
            "  `"currency`": `"string (Account currency code)`","
            "  `"balance`": `"integer (Account balance)`","
            "  `"comment`": `"string (Account description)`","
            "  `"creditCardStatementDate`": `"integer (The statement date of the credit card account)`","
            "  `"displayOrder`": `"integer (The display order of the account)`","
            "  `"isAsset`": `"boolean (Whether the account is an asset account)`","
            "  `"isLiability`": `"boolean (Whether the account is a liability account)`","
            "  `"hidden`": `"boolean (Whether the account is hidden)`","
            "  `"subAccounts`": [`"every sub-account object like account object`"]"
            "}"
        )
    }
    @{
        Name = "transaction-categories-list"
        Description = "Get all transaction categories"
        Method = "GET"
        Path = "transaction/categories/list.json"
        RequiresTimezone = $false
        RequiredParams = @()
        OptionalParams = @()
        ParamTypes = @{}
        ParamDescriptions = @{}
        ResponseStructure = @(
            "{"
            "  `"transaction category type (1: Income, 2: Expense, 3:Transfer)`": ["
            "    {"
            "      `"id`": `"string (Transaction category ID)`","
            "      `"name`": `"string (Transaction category name)`","
            "      `"parentId`": `"string (Parent transaction category ID)`","
            "      `"type`": `"integer (Transaction category type)`","
            "      `"icon`": `"string (Transaction category icon ID)`","
            "      `"color`": `"string (Transaction category icon color, hex color code RRGGBB)`","
            "      `"comment`": `"string (Transaction category description)`","
            "      `"displayOrder`": `"integer (The display order of the transaction category)`","
            "      `"hidden`": `"boolean (Whether the transaction category is hidden)`","
            "      `"subCategories`": [`"each sub-category object like a transaction category object`"]"
            "    }"
            "  ]"
            "}"
        )
    }
    @{
        Name = "transaction-categories-add"
        Description = "Add transaction category"
        Method = "POST"
        Path = "transaction/categories/add.json"
        RequiresTimezone = $false
        RequiredParams = @("name", "type", "icon", "color")
        OptionalParams = @("parentId", "comment")
        ParamTypes = @{
            "name" = "string"
            "type" = "integer"
            "parentId" = "string"
            "icon" = "string"
            "color" = "string"
            "comment" = "string"
        }
        ParamDescriptions = @{
            "name" = "string (Transaction category name)"
            "type" = "integer (Transaction category type, 1: Income, 2: Expense, 3: Transfer)"
            "parentId" = "string (Parent transaction category ID, 0 for primary category)"
            "icon" = "string (Transaction category icon ID)"
            "color" = "string (Transaction category icon color, hex color code RRGGBB)"
            "comment" = "string (Transaction category description)"
        }
        ResponseStructure = @(
            "{"
            "  `"id`": `"string (Transaction category ID)`","
            "  `"name`": `"string (Transaction category name)`","
            "  `"parentId`": `"string (Parent transaction category ID)`","
            "  `"type`": `"integer (Transaction category type)`","
            "  `"icon`": `"string (Transaction category icon ID)`","
            "  `"color`": `"string (Transaction category icon color)`","
            "  `"comment`": `"string (Transaction category description)`","
            "  `"displayOrder`": `"integer (The display order of the transaction category)`","
            "  `"hidden`": `"boolean (Whether the transaction category is hidden)`","
            "  `"subCategories`": [`"each sub-category object like a transaction category object`"]"
            "}"
        )
    }
    @{
        Name = "transaction-tags-list"
        Description = "Get all transaction tags list"
        Method = "GET"
        Path = "transaction/tags/list.json"
        RequiresTimezone = $false
        RequiredParams = @()
        OptionalParams = @()
        ParamTypes = @{}
        ParamDescriptions = @{}
        ResponseStructure = @(
            "["
            "  {"
            "    `"id`": `"string (Transaction tag ID)`","
            "    `"name`": `"string (Transaction tag name)`","
            "    `"groupId`": `"string (Transaction tag group ID)`","
            "    `"displayOrder`": `"integer (The display order of the transaction tag)`","
            "    `"hidden`": `"boolean (Whether the transaction tag is hidden)`""
            "  }"
            "]"
        )
    }
    @{
        Name = "transaction-tags-add"
        Description = "Add transaction tag"
        Method = "POST"
        Path = "transaction/tags/add.json"
        RequiresTimezone = $false
        RequiredParams = @("name")
        OptionalParams = @("groupId")
        ParamTypes = @{
            "name" = "string"
            "groupId" = "string"
        }
        ParamDescriptions = @{
            "name" = "string (Transaction tag name)"
            "groupId" = "string (Transaction tag group ID, 0 means default group)"
        }
        ResponseStructure = @(
            "{"
            "  `"id`": `"string (Transaction tag ID)`","
            "  `"name`": `"string (Transaction tag name)`","
            "  `"groupId`": `"string (Transaction tag group ID)`","
            "  `"displayOrder`": `"integer (The display order of the transaction tag)`","
            "  `"hidden`": `"boolean (Whether the transaction tag is hidden)`""
            "}"
        )
    }
    @{
        Name = "transactions-list"
        Description = "Get transactions list"
        Method = "GET"
        Path = "transactions/list.json"
        RequiresTimezone = $true
        RequiredParams = @("count")
        OptionalParams = @("type", "category_ids", "account_ids", "tag_filter", "amount_filter", "keyword", "max_time", "min_time", "page", "with_count", "with_pictures", "trim_account", "trim_category", "trim_tag")
        ParamTypes = @{
            "count" = "integer"
            "type" = "integer"
            "category_ids" = "string"
            "account_ids" = "string"
            "tag_filter" = "string"
            "amount_filter" = "string"
            "keyword" = "string"
            "max_time" = "integer"
            "min_time" = "integer"
            "page" = "integer"
            "with_count" = "boolean"
            "with_pictures" = "boolean"
            "trim_account" = "boolean"
            "trim_category" = "boolean"
            "trim_tag" = "boolean"
        }
        ParamDescriptions = @{
            "count" = "integer (The count of transactions per page, maximum is 50)"
            "type" = "integer (Filter transaction by type, 1: Balance modification, 2: Income, 3: Expense, 4: Transfer)"
            "category_ids" = "string (Filter by category IDs, separated by comma)"
            "account_ids" = "string (Filter by account IDs, separated by comma)"
            "tag_filter" = "string (Filter by tags)"
            "amount_filter" = "string (Filter by amount)"
            "keyword" = "string (Filter by keyword)"
            "max_time" = "integer (The maximum time sequence ID, Set to 0 for latest)"
            "min_time" = "integer (The minimum time sequence ID)"
            "page" = "integer (Specified page integer)"
            "with_count" = "boolean (Whether to get total count)"
            "with_pictures" = "boolean (Whether to get picture IDs)"
            "trim_account" = "boolean (Whether to get account ID only)"
            "trim_category" = "boolean (Whether to get category ID only)"
            "trim_tag" = "boolean (Whether to get tag IDs only)"
        }
        ResponseStructure = @(
            "{"
            "  `"items`": ["
            "    {"
            "      `"id`": `"string (Transaction ID)`","
            "      `"timeSequenceId`": `"string (Transaction time sequence ID)`","
            "      `"type`": `"integer (Transaction type)`","
            "      `"categoryId`": `"string (Transaction category ID)`","
            "      `"category`": `"object (Transaction category object)`","
            "      `"time`": `"integer (Transaction unix time)`","
            "      `"utcOffset`": `"integer (Transaction time zone offset minutes)`","
            "      `"sourceAccountId`": `"string (Source account ID)`","
            "      `"sourceAccount`": `"object (Source account object)`","
            "      `"destinationAccountId`": `"string (Destination account ID)`","
            "      `"destinationAccount`": `"object (Destination account object)`","
            "      `"sourceAmount`": `"integer (Source amount, supports up to two decimals. For example, a value of '1234' represents an amount of '12.34')`","
            "      `"destinationAmount`": `"integer (Destination amount, supports up to two decimals. For example, a value of '1234' represents an amount of '12.34')`","
            "      `"hideAmount`": `"boolean (Whether to hide the amount)`","
            "      `"tagIds`": [`"each string representing a transaction tag ID`"],"
            "      `"tags`": [`"each object representing a transaction tag object`"],"
            "      `"pictures`": [`"each object representing a transaction picture object`"],"
            "      `"comment`": `"string (Transaction description)`","
            "      `"geoLocation`": `"object (Transaction geographic location)`","
            "      `"editable`": `"boolean (Whether the transaction is editable)`""
            "    }"
            "  ],"
            "  `"nextTimeSequenceId`": `"integer (The next cursor 'max_time' parameter when requesting older data)`","
            "  `"totalCount`": `"integer (The total count of transactions)`""
            "}"
        )
    }
    @{
        Name = "transactions-list-all"
        Description = "Get all transactions list"
        Method = "GET"
        Path = "transactions/list/all.json"
        RequiresTimezone = $true
        RequiredParams = @()
        OptionalParams = @("type", "category_ids", "account_ids", "tag_filter", "amount_filter", "keyword", "start_time", "end_time", "with_pictures", "trim_account", "trim_category", "trim_tag")
        ParamTypes = @{
            "type" = "integer"
            "category_ids" = "string"
            "account_ids" = "string"
            "tag_filter" = "string"
            "amount_filter" = "string"
            "keyword" = "string"
            "start_time" = "integer"
            "end_time" = "integer"
            "with_pictures" = "boolean"
            "trim_account" = "boolean"
            "trim_category" = "boolean"
            "trim_tag" = "boolean"
        }
        ParamDescriptions = @{
            "type" = "integer (Filter transaction by type, 1: Balance modification, 2: Income, 3: Expense, 4: Transfer)"
            "category_ids" = "string (Filter by category IDs, separated by comma)"
            "account_ids" = "string (Filter by account IDs, separated by comma)"
            "tag_filter" = "string (Filter by tags)"
            "amount_filter" = "string (Filter by amount)"
            "keyword" = "string (Filter by keyword)"
            "start_time" = "integer (Transaction list start unix time)"
            "end_time" = "integer (Transaction list end unix time)"
            "with_pictures" = "boolean (Whether to get picture IDs)"
            "trim_account" = "boolean (Whether to get account ID only)"
            "trim_category" = "boolean (Whether to get category ID only)"
            "trim_tag" = "boolean (Whether to get tag IDs only)"
        }
        ResponseStructure = @(
            "["
            "  {"
            "    `"id`": `"string (Transaction ID)`","
            "    `"timeSequenceId`": `"string (Transaction time sequence ID)`","
            "    `"type`": `"integer (Transaction type)`","
            "    `"categoryId`": `"string (Transaction category ID)`","
            "    `"category`": `"object (Transaction category object)`","
            "    `"time`": `"integer (Transaction unix time)`","
            "    `"utcOffset`": `"integer (Transaction time zone offset minutes)`","
            "    `"sourceAccountId`": `"string (Source account ID)`","
            "    `"sourceAccount`": `"object (Source account object)`","
            "    `"destinationAccountId`": `"string (Destination account ID)`","
            "    `"destinationAccount`": `"object (Destination account object)`","
            "    `"sourceAmount`": `"integer (Source amount, supports up to two decimals. For example, a value of '1234' represents an amount of '12.34')`","
            "    `"destinationAmount`": `"integer (Destination amount, supports up to two decimals. For example, a value of '1234' represents an amount of '12.34')`","
            "    `"hideAmount`": `"boolean (Whether to hide the amount)`","
            "    `"tagIds`": [`"each string representing a transaction tag ID`"],"
            "    `"tags`": [`"each object representing a transaction tag object`"],"
            "    `"pictures`": [`"each object representing a transaction picture object`"],"
            "    `"comment`": `"string (Transaction description)`","
            "    `"geoLocation`": `"object (Transaction geographic location)`","
            "    `"editable`": `"boolean (Whether the transaction is editable)`""
            "  }"
            "]"
        )
    }
    @{
        Name = "transactions-add"
        Description = "Add transaction"
        Method = "POST"
        Path = "transactions/add.json"
        RequiresTimezone = $true
        RequiredParams = @("type", "categoryId", "time", "utcOffset", "sourceAccountId", "sourceAmount")
        OptionalParams = @("destinationAccountId", "destinationAmount", "hideAmount", "tagIds", "pictureIds", "comment", "geoLocation")
        ParamTypes = @{
            "type" = "integer"
            "categoryId" = "string"
            "time" = "integer"
            "utcOffset" = "integer"
            "sourceAccountId" = "string"
            "sourceAmount" = "integer"
            "destinationAccountId" = "string"
            "destinationAmount" = "integer"
            "hideAmount" = "boolean"
            "tagIds" = "string_array"
            "pictureIds" = "string_array"
            "comment" = "string"
            "geoLocation" = "geo_location"
        }
        ParamDescriptions = @{
            "type" = "integer (Transaction type, 1: Balance Modification, 2: Income, 3: Expense, 4: Transfer)"
            "categoryId" = "string (Transaction category ID)"
            "time" = "integer (Transaction unix time)"
            "utcOffset" = "integer (Transaction time zone offset minutes)"
            "sourceAccountId" = "string (Source account ID)"
            "sourceAmount" = "integer (Source amount, supports up to two decimals. For example, a value of `"1234`" represents an amount of `"12.34`")"
            "destinationAccountId" = "string (Destination account ID)"
            "destinationAmount" = "integer (Destination amount, supports up to two decimals. For example, a value of `"1234`" represents an amount of `"12.34`")"
            "hideAmount" = "boolean (Whether to hide amount)"
            "tagIds" = "string (Transaction tag IDs, separated by comma, e.g. `"tagid1,tagid2`")"
            "pictureIds" = "string (Transaction picture IDs, separated by comma, e.g. `"picid1,picid2`")"
            "comment" = "string (Transaction description)"
            "geoLocation" = "string (Transaction geographic location, format: longitude,latitude, e.g. `"116.33,39.93`")"
        }
        ResponseStructure = @(
            "{"
            "  `"id`": `"string (Transaction ID)`","
            "  `"timeSequenceId`": `"string (Transaction time sequence ID)`","
            "  `"type`": `"integer (Transaction type)`","
            "  `"categoryId`": `"string (Transaction category ID)`","
            "  `"category`": `"object (Transaction category object)`","
            "  `"time`": `"integer (Transaction unix time)`","
            "  `"utcOffset`": `"integer (Transaction time zone offset minutes)`","
            "  `"sourceAccountId`": `"string (Source account ID)`","
            "  `"sourceAccount`": `"object (Source account object)`","
            "  `"destinationAccountId`": `"string (Destination account ID)`","
            "  `"destinationAccount`": `"object (Destination account object)`","
            "  `"sourceAmount`": `"integer (Source amount)`","
            "  `"destinationAmount`": `"integer (Destination amount)`","
            "  `"hideAmount`": `"boolean (Whether to hide the amount)`","
            "  `"tagIds`": [`"each string representing a transaction tag ID`"],"
            "  `"tags`": [`"each object representing a transaction tag object`"],"
            "  `"pictures`": [`"each object representing a transaction picture object`"],"
            "  `"comment`": `"string (Transaction description)`","
            "  `"geoLocation`": `"object (Transaction geographic location)`","
            "  `"editable`": `"boolean (Whether the transaction is editable)`""
            "}"
        )
    }
)

# Reference: https://github.com/unicode-org/cldr/blob/main/common/supplemental/windowsZones.xml
$TIMEZONE_IANA_NAMES = @{
    "Dateline Standard Time" = "Etc/GMT+12"
    "UTC-11" = "Etc/GMT+11"
    "Aleutian Standard Time" = "America/Adak"
    "Hawaiian Standard Time" = "Pacific/Honolulu"
    "Marquesas Standard Time" = "Pacific/Marquesas"
    "Alaskan Standard Time" = "America/Anchorage"
    "UTC-09" = "Etc/GMT+9"
    "Pacific Standard Time (Mexico)" = "America/Tijuana"
    "UTC-08" = "Etc/GMT+8"
    "Pacific Standard Time" = "America/Los_Angeles"
    "US Mountain Standard Time" = "America/Phoenix"
    "Mountain Standard Time (Mexico)" = "America/Mazatlan"
    "Mountain Standard Time" = "America/Denver"
    "Yukon Standard Time" = "America/Whitehorse"
    "Central America Standard Time" = "America/Guatemala"
    "Central Standard Time" = "America/Chicago"
    "Easter Island Standard Time" = "Pacific/Easter"
    "Central Standard Time (Mexico)" = "America/Mexico_City"
    "Canada Central Standard Time" = "America/Regina"
    "SA Pacific Standard Time" = "America/Bogota"
    "Eastern Standard Time (Mexico)" = "America/Cancun"
    "Eastern Standard Time" = "America/New_York"
    "Haiti Standard Time" = "America/Port-au-Prince"
    "Cuba Standard Time" = "America/Havana"
    "US Eastern Standard Time" = "America/Indianapolis"
    "Turks And Caicos Standard Time" = "America/Grand_Turk"
    "Paraguay Standard Time" = "America/Asuncion"
    "Atlantic Standard Time" = "America/Halifax"
    "Venezuela Standard Time" = "America/Caracas"
    "Central Brazilian Standard Time" = "America/Cuiaba"
    "SA Western Standard Time" = "America/La_Paz"
    "Pacific SA Standard Time" = "America/Santiago"
    "Newfoundland Standard Time" = "America/St_Johns"
    "Tocantins Standard Time" = "America/Araguaina"
    "E. South America Standard Time" = "America/Sao_Paulo"
    "SA Eastern Standard Time" = "America/Cayenne"
    "Argentina Standard Time" = "America/Buenos_Aires"
    "Greenland Standard Time" = "America/Godthab"
    "Montevideo Standard Time" = "America/Montevideo"
    "Magallanes Standard Time" = "America/Punta_Arenas"
    "Saint Pierre Standard Time" = "America/Miquelon"
    "Bahia Standard Time" = "America/Bahia"
    "UTC-02" = "Etc/GMT+2"
    "Azores Standard Time" = "Atlantic/Azores"
    "Cape Verde Standard Time" = "Atlantic/Cape_Verde"
    "UTC" = "Etc/UTC"
    "GMT Standard Time" = "Europe/London"
    "Greenwich Standard Time" = "Atlantic/Reykjavik"
    "Sao Tome Standard Time" = "Africa/Sao_Tome"
    "Morocco Standard Time" = "Africa/Casablanca"
    "W. Europe Standard Time" = "Europe/Berlin"
    "Central Europe Standard Time" = "Europe/Budapest"
    "Romance Standard Time" = "Europe/Paris"
    "Central European Standard Time" = "Europe/Warsaw"
    "W. Central Africa Standard Time" = "Africa/Lagos"
    "Jordan Standard Time" = "Asia/Amman"
    "GTB Standard Time" = "Europe/Bucharest"
    "Middle East Standard Time" = "Asia/Beirut"
    "Egypt Standard Time" = "Africa/Cairo"
    "E. Europe Standard Time" = "Europe/Chisinau"
    "Syria Standard Time" = "Asia/Damascus"
    "West Bank Standard Time" = "Asia/Hebron"
    "South Africa Standard Time" = "Africa/Johannesburg"
    "FLE Standard Time" = "Europe/Kiev"
    "Israel Standard Time" = "Asia/Jerusalem"
    "South Sudan Standard Time" = "Africa/Juba"
    "Kaliningrad Standard Time" = "Europe/Kaliningrad"
    "Sudan Standard Time" = "Africa/Khartoum"
    "Libya Standard Time" = "Africa/Tripoli"
    "Namibia Standard Time" = "Africa/Windhoek"
    "Arabic Standard Time" = "Asia/Baghdad"
    "Turkey Standard Time" = "Europe/Istanbul"
    "Arab Standard Time" = "Asia/Riyadh"
    "Belarus Standard Time" = "Europe/Minsk"
    "Russian Standard Time" = "Europe/Moscow"
    "E. Africa Standard Time" = "Africa/Nairobi"
    "Iran Standard Time" = "Asia/Tehran"
    "Arabian Standard Time" = "Asia/Dubai"
    "Astrakhan Standard Time" = "Europe/Astrakhan"
    "Azerbaijan Standard Time" = "Asia/Baku"
    "Russia Time Zone 3" = "Europe/Samara"
    "Mauritius Standard Time" = "Indian/Mauritius"
    "Saratov Standard Time" = "Europe/Saratov"
    "Georgian Standard Time" = "Asia/Tbilisi"
    "Volgograd Standard Time" = "Europe/Volgograd"
    "Caucasus Standard Time" = "Asia/Yerevan"
    "Afghanistan Standard Time" = "Asia/Kabul"
    "West Asia Standard Time" = "Asia/Tashkent"
    "Ekaterinburg Standard Time" = "Asia/Yekaterinburg"
    "Pakistan Standard Time" = "Asia/Karachi"
    "Qyzylorda Standard Time" = "Asia/Qyzylorda"
    "India Standard Time" = "Asia/Calcutta"
    "Sri Lanka Standard Time" = "Asia/Colombo"
    "Nepal Standard Time" = "Asia/Katmandu"
    "Central Asia Standard Time" = "Asia/Bishkek"
    "Bangladesh Standard Time" = "Asia/Dhaka"
    "Omsk Standard Time" = "Asia/Omsk"
    "Myanmar Standard Time" = "Asia/Rangoon"
    "SE Asia Standard Time" = "Asia/Bangkok"
    "Altai Standard Time" = "Asia/Barnaul"
    "W. Mongolia Standard Time" = "Asia/Hovd"
    "North Asia Standard Time" = "Asia/Krasnoyarsk"
    "N. Central Asia Standard Time" = "Asia/Novosibirsk"
    "Tomsk Standard Time" = "Asia/Tomsk"
    "China Standard Time" = "Asia/Shanghai"
    "North Asia East Standard Time" = "Asia/Irkutsk"
    "Singapore Standard Time" = "Asia/Singapore"
    "W. Australia Standard Time" = "Australia/Perth"
    "Taipei Standard Time" = "Asia/Taipei"
    "Ulaanbaatar Standard Time" = "Asia/Ulaanbaatar"
    "Aus Central W. Standard Time" = "Australia/Eucla"
    "Transbaikal Standard Time" = "Asia/Chita"
    "Tokyo Standard Time" = "Asia/Tokyo"
    "North Korea Standard Time" = "Asia/Pyongyang"
    "Korea Standard Time" = "Asia/Seoul"
    "Yakutsk Standard Time" = "Asia/Yakutsk"
    "Cen. Australia Standard Time" = "Australia/Adelaide"
    "AUS Central Standard Time" = "Australia/Darwin"
    "E. Australia Standard Time" = "Australia/Brisbane"
    "AUS Eastern Standard Time" = "Australia/Sydney"
    "West Pacific Standard Time" = "Pacific/Port_Moresby"
    "Tasmania Standard Time" = "Australia/Hobart"
    "Vladivostok Standard Time" = "Asia/Vladivostok"
    "Lord Howe Standard Time" = "Australia/Lord_Howe"
    "Bougainville Standard Time" = "Pacific/Bougainville"
    "Russia Time Zone 10" = "Asia/Srednekolymsk"
    "Magadan Standard Time" = "Asia/Magadan"
    "Norfolk Standard Time" = "Pacific/Norfolk"
    "Sakhalin Standard Time" = "Asia/Sakhalin"
    "Central Pacific Standard Time" = "Pacific/Guadalcanal"
    "Russia Time Zone 11" = "Asia/Kamchatka"
    "New Zealand Standard Time" = "Pacific/Auckland"
    "UTC+12" = "Etc/GMT-12"
    "Fiji Standard Time" = "Pacific/Fiji"
    "Chatham Islands Standard Time" = "Pacific/Chatham"
    "UTC+13" = "Etc/GMT-13"
    "Tonga Standard Time" = "Pacific/Tongatapu"
    "Samoa Standard Time" = "Pacific/Apia"
    "Line Islands Standard Time" = "Pacific/Kiritimati"
}

function Write-Red($msg) {
    Write-Host $msg -ForegroundColor Red
}

function Write-Yellow($msg) {
    Write-Host $msg -ForegroundColor Yellow
}

function Url-Encode {
    param([string]$text)
    return [System.Uri]::EscapeDataString($text)
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

function Show-Help {
    $exampleTimezoneName = Get-ExampleTimezoneName
    $exampleTimezoneOffset = Get-ExampleTimezoneOffset

    Write-Host "ezBookkeeping API Tools"
    Write-Host ""
    Write-Host "A command-line tool for calling ezBookkeeping APIs"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "    ebktools.ps1 [-tzName <name>] [-tzOffset <offset>] <command> [command-options]"
    Write-Host ""
    Write-Host "Environment Variables (Required):"
    Write-Host "    EBKTOOL_SERVER_BASEURL      ezBookkeeping server base URL (e.g., http://localhost:8080)"
    Write-Host "    EBKTOOL_TOKEN               ezBookkeeping API token"
    Write-Host ""
    Write-Host "Global Options:"
    Write-Host "    -tzName <name>              The IANA timezone name of current timezone. For example, for Beijing Time it is 'Asia/Shanghai'."
    Write-Host "    -tzOffset <offset>          The offset in minutes of the current timezone from UTC. For example, for Beijing Time which is UTC+8, the value is '480'. If both '-tzName' and '-tzOffset' are set, '-tzName' takes priority. If neither is set, the current system time zone is used by default."
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
    Write-Host "    ebktools.ps1 help transactions-add"
    Write-Host ""
    Write-Host "    # Call accounts-list API"
    Write-Host "    ebktools.ps1 accounts-list"
    Write-Host ""
    Write-Host "    # Call API with timezone name"
    Write-Host "    ebktools.ps1 -tzName $exampleTimezoneName transactions-list -count 10"
    Write-Host ""
    Write-Host "    # Call API with timezone offset"
    Write-Host "    ebktools.ps1 -tzOffset $exampleTimezoneOffset transactions-list -count 10"
}

function Show-CommandList {
    Write-Host "Available API Commands:"
    Write-Host ""

    foreach ($config in $API_CONFIGS) {
        $name = $config.Name.PadRight(31)
        Write-Host "  $name$($config.Description)"
    }

    Write-Host ""
    Write-Host "Use 'ebktools.ps1 help <api-command>' to see detailed information about a API command."
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

            if ($i + 1 -lt $commandArgs.Count -and -not $commandArgs[$i + 1].StartsWith("-")) {
                $paramType = "string"
                $paramValue = $commandArgs[$i + 1]

                if ($paramTypes -and $paramTypes.ContainsKey($paramName)) {
                    $paramType = $paramTypes[$paramName]
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

    $serverBaseUrl = $env:EBKTOOL_SERVER_BASEURL
    $authToken = $env:EBKTOOL_TOKEN

    if (-not $serverBaseUrl) {
        Write-Red "Error: Environment variable 'EBKTOOL_SERVER_BASEURL' is not set."
        Write-Host "Please set it to your ezBookkeeping server base URL (e.g., http://localhost:8080)"
        exit 1
    }

    if (-not $authToken) {
        Write-Red "Error: Environment variable 'EBKTOOL_TOKEN' is not set."
        Write-Host "Please set it to your API token."
        exit 1
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

            Write-Yellow "Calling API: $($config.Method) $url"
            Write-Host ""

            if ($params.Count -gt 0) {
                $body = ConvertTo-Json -Depth 10 $params
                $response = Invoke-WebRequest -Uri $url -Method POST -Headers $headers -Body $body -ErrorAction Stop -UseBasicParsing
            } else {
                $response = Invoke-WebRequest -Uri $url -Method POST -Headers $headers -ErrorAction Stop -UseBasicParsing
            }

            $response = ConvertFrom-Json $response.Content
        } else {
            if ($params.Count -gt 0) {
                $queryString = ($params.GetEnumerator() | ForEach-Object { "$($_.Key)=$(Url-Encode $_.Value)" }) -join "&"
                $url = "${url}?$queryString"
            }

            Write-Yellow "Calling API: $($config.Method) $url"
            Write-Host ""

            $response = Invoke-WebRequest -Uri $url -Method $config.Method -Headers $headers -ErrorAction Stop -UseBasicParsing
            $response = ConvertFrom-Json $response.Content
        }

        if ($response.PSObject.Properties.Name -contains "success") {
            if ($response.success -eq $true) {
                Write-Host "Response Result:"
                if ($response.PSObject.Properties.Name -contains "result") {
                    $jsonOutput = ConvertTo-Json -Depth 10 -Compress $response.result | Format-Json
                    Write-Host $jsonOutput
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
