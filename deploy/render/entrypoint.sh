#!/bin/sh
set -e

DB_PATH=/ezbookkeeping/data/ezbookkeeping.db

# On a fresh disk (every Render restart or deploy) pull the latest copy from S3.
# A first-ever start has no replica yet, so this is a no-op and the app creates a new database.
litestream restore -if-db-not-exists -if-replica-exists -config /etc/litestream.yml "$DB_PATH"

# Replicate continuously while the app runs; Litestream stops the app and flushes on shutdown.
exec litestream replicate -config /etc/litestream.yml -exec "/ezbookkeeping/ezbookkeeping server run"
