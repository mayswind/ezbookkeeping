---
name: ezbookkeeping
description: Use ezBookkeeping API Tools script to record new transactions, query transactions, retrieve account information, retrieve categories, retrieve tags, and retrieve exchange rate data in the self hosted personal finance application ezBookkeeping.
---

# ezBookkeeping API Tools

## Usage

### List all supported commands

Linux / macOS

```bash
sh scripts/ebktools.sh list
```

Windows

```powershell
scripts\ebktools.ps1 list
```

### Show help for a specific command

Linux / macOS

```bash
sh scripts/ebktools.sh help <command>
```

Windows

```powershell
scripts\ebktools.ps1 help <command>
```

### Call API

Linux / macOS

```bash
sh scripts/ebktools.sh [global-options] <command> [command-options]
```

Windows

```powershell
scripts\ebktools.ps1 [global-options] <command> [command-options]
```

## Troubleshooting

If the script reports that the environment variable `EBKTOOL_SERVER_BASEURL` or `EBKTOOL_TOKEN` is not set, user can define them as system environment variables, or create a `.env` file in the user home directory that contains these two variables and place it there.

The meanings of these environment variables are as follows:

| Variable | Required | Description |
| --- | --- | --- |
| `EBKTOOL_SERVER_BASEURL` | Required | ezBookkeeping server base URL (e.g., `http://localhost:8080`) |
| `EBKTOOL_TOKEN` | Required | ezBookkeeping API token |

## Reference

ezBookkeeping: [https://ezbookkeeping.mayswind.net](https://ezbookkeeping.mayswind.net)