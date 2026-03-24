---
name: oscar
description: Use oscar API Tools script to record new transactions, query transactions, retrieve account information, retrieve categories, retrieve tags, and retrieve exchange rate data in the self hosted personal finance application oscar.
---

# oscar API Tools

## Usage

### List all supported commands

```bash
sh scripts/oscar-tools.sh list
```

### Show help for a specific command

```bash
sh scripts/oscar-tools.sh help <command>
```

### Call API

```bash
sh scripts/oscar-tools.sh [global-options] <command> [command-options]
```

## Troubleshooting

If the script reports that the environment variable `OSCAR_SERVER_BASEURL` or `OSCAR_TOKEN` is not set, user can define them as system environment variables, or create a `.env` file in the user home directory that contains these two variables and place it there.

The meanings of these environment variables are as follows:

| Variable | Required | Description |
| --- | --- | --- |
| `OSCAR_SERVER_BASEURL` | Required | oscar server base URL (e.g., `http://localhost:8080`) |
| `OSCAR_TOKEN` | Required | oscar API token |

## Reference

oscar by nicodAImus: [https://github.com/Paxtiny/oscar](https://github.com/Paxtiny/oscar)
