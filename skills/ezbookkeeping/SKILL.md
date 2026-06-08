---
name: ezbookkeeping
description: Use ezBookkeeping API Tools to query and manage a self-hosted ezBookkeeping instance, including accounts, transaction categories, tag groups, tags, transactions, exchange rates, sessions, and server version. Use when an agent needs to list, inspect, create, modify, move, hide, delete, or batch-update personal finance records through ezBookkeeping API commands.
---

# ezBookkeeping API Tools

Use the bundled scripts in `scripts/` to call ezBookkeeping `/api/v1` endpoints. Prefer the script over hand-written HTTP calls because it handles authentication, timezone headers, parameter typing, formatted output, and dry-run previews.

## Setup

Set credentials in environment variables, or in a `.env` file in the current directory, parent directory, or home directory:

```text
EBKTOOL_SERVER_BASEURL=http://localhost:8080
EBKTOOL_TOKEN=YOUR_TOKEN
```

The scripts load command definitions from `scripts/api-configs.json`. To discover current coverage and exact parameters, run:

```powershell
scripts\ebktools.ps1 list
scripts\ebktools.ps1 help transactions-add
```

```bash
sh scripts/ebktools.sh list
sh scripts/ebktools.sh help transactions-add
```

## Command Use

Use PowerShell on Windows:

```powershell
scripts\ebktools.ps1 [global-options] <command> [command-options]
```

Use POSIX shell on Linux/macOS:

```bash
sh scripts/ebktools.sh [global-options] <command> [command-options]
```

Global options:

| PowerShell | Shell | Purpose |
| --- | --- | --- |
| `-tzName <name>` | `--tz-name <name>` | Send an IANA timezone name, such as `Asia/Shanghai`. |
| `-tzOffset <minutes>` | `--tz-offset <minutes>` | Send timezone offset minutes, such as `480`. |
| `-rawResponse` | `--raw-response` | Print raw JSON instead of tables. |
| `-dryRun` | `--dry-run` | Print method, URL, headers, and body without sending the request. |

Use dry-run before destructive or broad write operations, then ask the user to confirm before executing the real command. This especially applies to `tokens-revoke`, `accounts-delete`, `accounts-sub-account-delete`, category/tag/tag-group deletes, `transactions-delete`, `transactions-batch-delete`, `transactions-move-all`, transaction batch updates, and `exchangerates-custom-delete`.

## Parameter Rules

Pass amounts as integer minor units: `1234` means `12.34`; expenses or liabilities may require negative values, such as `-1234`.

Pass comma-separated IDs for `string_array` parameters such as `tagIds`, `pictureIds`, `transactionIds`, and `ids`:

```powershell
scripts\ebktools.ps1 -dryRun transactions-batch-add-tags -transactionIds 1001,1002 -tagIds 8,9
```

Pass `geo_location` as `longitude,latitude`:

```powershell
scripts\ebktools.ps1 -dryRun transactions-add -type 3 -categoryId 12 -time 1710000000 -utcOffset 480 -sourceAccountId 1 -sourceAmount -1234 -geoLocation 116.33,39.93
```

Pass object or array parameters as JSON strings for commands such as `accounts-move`, `transaction-categories-add-batch`, `transaction-tags-add-batch`, and tag/category/account move commands:

```powershell
scripts\ebktools.ps1 -dryRun accounts-move -newDisplayOrders '[{"id":"1","displayOrder":1}]'
```

```bash
sh scripts/ebktools.sh --dry-run accounts-move --newDisplayOrders '[{"id":"1","displayOrder":1}]'
```

## Common Workflows

Inspect current bookkeeping data before writing:

```powershell
scripts\ebktools.ps1 accounts-list
scripts\ebktools.ps1 transaction-categories-list
scripts\ebktools.ps1 transaction-tags-list
scripts\ebktools.ps1 transactions-list -count 20
```

Create or modify a transaction:

```powershell
scripts\ebktools.ps1 -dryRun transactions-add -type 3 -categoryId 12 -time 1710000000 -utcOffset 480 -sourceAccountId 1 -sourceAmount -1234 -tagIds 8,9 -comment "Lunch"
scripts\ebktools.ps1 help transactions-modify
```

Update supporting metadata:

```powershell
scripts\ebktools.ps1 -dryRun transaction-tags-add -name "Business" -groupId 0
scripts\ebktools.ps1 -dryRun transaction-categories-hide -id 12 -hidden true
```

Manage exchange rates:

```powershell
scripts\ebktools.ps1 exchangerates-latest
scripts\ebktools.ps1 -dryRun exchangerates-custom-update -currency USD -rate 7.2
```
