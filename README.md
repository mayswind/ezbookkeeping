# nicodAImus oscar

Privacy-first AI expense tracker. Fork of [ezBookkeeping](https://github.com/mayswind/ezbookkeeping) with client-side encryption, multi-user groups, and budgeting.

## What it does

- **Scan invoices** with AI (8 LLM providers) or on-device OCR (Tesseract)
- **Import bank data** from CSV (52 encodings), OFX, QFX, QIF, Beancount, and more
- **Track expenses** with two-level categories, tags, multi-currency, and analytics
- **Encrypt everything** client-side (AES-256-GCM) - the server never sees your data
- **Bring your own storage** - local, S3/MinIO, or WebDAV per user

## Privacy by design

oscar uses user-passphrase-derived encryption. Your passphrase never leaves your device. The server stores only encrypted blobs. Lost passphrase = lost data (by design, no recovery).

Not even the operator can read your data.

## Quick start (Docker)

```bash
docker run -d --name oscar \
  -p 8080:8080 \
  -v oscar-data:/oscar/data \
  oscar:latest
```

## Configuration

oscar supports SQLite, MySQL, and PostgreSQL. Configure in `conf/oscar.ini`.

### nicodAImus AI integration

```ini
[llm_image_recognition]
llm_provider = openai_compatible
openai_compatible_base_url = https://chat.nicodaimus.com/v1
openai_compatible_api_key = <your API key>
openai_compatible_model_id = auto
```

## Development

### Prerequisites

- Go 1.25+
- Node.js 24+
- PostgreSQL 17+ (or Docker)

### Build from source

```bash
# Backend
go build -o oscar oscar.go

# Frontend
npm install && npm run build

# Docker
docker build . -t oscar:local
```

### Run tests

```bash
go test ./... -v   # Backend
npm test           # Frontend
```

## Project structure

```
oscar.go            # Main entry point
cmd/                # CLI commands (server, database, user-data)
pkg/
  api/              # REST API endpoints
  services/         # Business logic
  models/           # Database models (xorm)
  storage/          # Object storage abstraction (local/S3/WebDAV)
  llm/              # LLM integration (receipt scanning)
  mcp/              # MCP server (7 tools)
src/                # Vue.js frontend (PWA)
conf/oscar.ini      # Configuration
```

## Upstream

oscar is a fork of [ezBookkeeping](https://github.com/mayswind/ezbookkeeping) v1.4.0 by MaysWind. We inherit its excellent transaction management, multi-currency support, import/export, analytics, i18n (19 languages), and MCP server. See the [ezBookkeeping documentation](https://ezbookkeeping.mayswind.net/) for inherited features.

## License

MIT - see [LICENSE](LICENSE)

Built by [nicodAImus](https://nicodaimus.com)
