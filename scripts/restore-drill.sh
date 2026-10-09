#!/usr/bin/env bash
# Backup restore drill for the Render + Litestream + S3 setup. READ-ONLY: it downloads the backup into a temporary
# folder, checks it, boots the app on a copy of it, and deletes everything afterwards. It never writes to the bucket
# and never touches the live service. Run it after changes to the backup setup and every month or so.
#
#   scripts/restore-drill.sh [env-file]          (default: .env, which needs LITESTREAM_ACCESS_KEY_ID and
#                                                 LITESTREAM_SECRET_ACCESS_KEY; bucket and region default to the project's)
# Needs: docker, sqlite3, python3, curl, and for the boot test a built ./ezbookkeeping and ./dist.
# Optional: DRILL_PORT (default 18080), DRILL_AS_OF (a UTC time such as 2026-10-08T07:20:00Z for a point-in-time restore).

set -u
ENV_FILE="${1:-.env}"
PORT="${DRILL_PORT:-18080}"
IMAGE="litestream/litestream:0.3.13"   # keep in step with deploy/render/Dockerfile
fails=0
pass() { echo "PASS  $1"; }
fail() { echo "FAIL  $1"; fails=$((fails + 1)); }

WORK="$(mktemp -d)"
APP_PID=""
cleanup() {
    [ -n "$APP_PID" ] && kill "$APP_PID" 2>/dev/null
    docker run --rm --entrypoint sh -v "$WORK:/w" "$IMAGE" -c 'chown -R '"$(id -u):$(id -g)"' /w' >/dev/null 2>&1
    rm -rf "$WORK"   # holds the credentials and a copy of the customers' data
}
trap cleanup EXIT

for tool in docker sqlite3 python3 curl; do
    command -v "$tool" >/dev/null || { echo "missing tool: $tool"; exit 2; }
done

# copy only the two credentials, literally (passwords may contain characters a shell would change), never printing them
python3 - "$ENV_FILE" "$WORK/ls.env" <<'PY' || { echo "cannot read credentials from $ENV_FILE"; exit 2; }
import os, sys
wanted = {"LITESTREAM_ACCESS_KEY_ID": None, "LITESTREAM_SECRET_ACCESS_KEY": None,
          "LITESTREAM_BUCKET": "ezbookkeeping-sqlite-bkup", "LITESTREAM_REGION": "eu-north-1"}
for line in open(sys.argv[1]):
    key, _, value = line.rstrip("\n").partition("=")
    if key in wanted:
        wanted[key] = value
missing = [k for k in ("LITESTREAM_ACCESS_KEY_ID", "LITESTREAM_SECRET_ACCESS_KEY") if not wanted[k]]
if missing:
    sys.exit("missing " + ", ".join(missing))
with open(sys.argv[2], "w") as out:
    os.fchmod(out.fileno(), 0o600)
    out.write("".join(f"{k}={v}\n" for k, v in wanted.items() if v))
PY

REPO="$(cd "$(dirname "$0")/.." && pwd)"
cp "$REPO/deploy/render/litestream.yml" "$WORK/litestream.yml"
touch "$WORK/placeholder.db"; mkdir "$WORK/out"
# litestream <command> [options...] <database path>   (options must come before the path)
litestream() { local command="$1"; shift
    docker run --rm --env-file "$WORK/ls.env" -v "$WORK/litestream.yml:/etc/litestream.yml:ro" \
        -v "$WORK/placeholder.db:/ezbookkeeping/data/ezbookkeeping.db:ro" -v "$WORK/out:/out" "$IMAGE" \
        "$command" -config /etc/litestream.yml "$@"; }
DB=/ezbookkeeping/data/ezbookkeeping.db

echo "== what the bucket holds"
GEN="$(litestream generations $DB 2>&1)"; echo "$GEN"
echo "$GEN" | grep -q "^s3" && pass "a backup generation exists" || { fail "no backup generation found (is replication running?)"; exit 1; }
echo "snapshots:"; litestream snapshots $DB 2>&1 | sed 's/^/  /'
LAST="$(echo "$GEN" | awk '/^s3/ {print $NF}')"
echo "last change copied to S3: $LAST   (now: $(date -u +%Y-%m-%dT%H:%M:%SZ))"
echo "  (an idle database writes nothing, so an old time is only a problem if people have been using the app since)"

echo "== restore"
AS_OF=(); [ -n "${DRILL_AS_OF:-}" ] && AS_OF=(-timestamp "$DRILL_AS_OF") && echo "point in time: $DRILL_AS_OF"
START=$(date +%s)
if litestream restore -if-db-not-exists -if-replica-exists "${AS_OF[@]}" -o /out/restored.db $DB >"$WORK/restore.log" 2>&1; then
    pass "restore finished in $(( $(date +%s) - START )) seconds"
else
    tail -3 "$WORK/restore.log"; fail "restore failed"; exit 1
fi
docker run --rm --entrypoint sh -v "$WORK/out:/out" "$IMAGE" -c 'chown -R '"$(id -u):$(id -g)"' /out' >/dev/null 2>&1
rm -f "$WORK"/out/restored.db.tmp-*
R="$WORK/out/restored.db"

echo "== the restored database"
[ "$(sqlite3 "$R" 'PRAGMA integrity_check;')" = "ok" ] && pass "integrity check ok" || fail "integrity check failed"
for table in user account '"transaction"'; do
    printf "      %-14s %s\n" "${table//\"/}" "$(sqlite3 "$R" "select count(*) from $table;" 2>&1)"
done
USERS="$(sqlite3 "$R" 'select count(*) from user;' 2>/dev/null)"
[ "${USERS:-0}" -ge 1 ] 2>/dev/null && pass "it contains users" || fail "no users in the restored database"
echo "      newest transaction: $(sqlite3 "$R" 'select coalesce(datetime(max(created_unix_time),"unixepoch")," none") from "transaction";') UTC"

echo "== boot the app on a copy"
if [ -x "$REPO/ezbookkeeping" ] && [ -d "$REPO/dist" ]; then
    mkdir -p "$WORK/app/storage"; cp "$R" "$WORK/app/app.db"
    # exec makes the background process the app itself, so cleanup stops exactly it
    ( cd "$REPO" && exec env -i PATH="$PATH" HOME="$HOME" EBK_DATABASE_TYPE=sqlite3 EBK_DATABASE_DB_PATH="$WORK/app/app.db" \
        EBK_SERVER_HTTP_PORT="$PORT" EBK_SERVER_DOMAIN=localhost EBK_SERVER_STATIC_ROOT_PATH=dist \
        EBK_STORAGE_LOCAL_FILESYSTEM_PATH="$WORK/app/storage" EBK_LOG_MODE=console EBK_MAIL_ENABLE_SMTP=false \
        EBK_SECURITY_SECRET_KEY=drill-only-not-a-real-secret EBK_DATABASE_AUTO_UPDATE_DATABASE=true \
        ./ezbookkeeping server run >"$WORK/app/boot.log" 2>&1 ) &
    APP_PID=$!
    for _ in $(seq 1 40); do curl -s -o /dev/null "localhost:$PORT/" && break; sleep 1; done
    [ "$(curl -s -o /dev/null -w '%{http_code}' "localhost:$PORT/")" = "200" ] && pass "the app starts on the restored data and serves pages" || { tail -5 "$WORK/app/boot.log"; fail "the app did not start on the restored data"; }
    [ "$(sqlite3 "$WORK/app/app.db" 'select count(*) from user;')" = "$USERS" ] && pass "its data is unchanged after the app's own upgrade step" || fail "user count changed on startup"
else
    echo "skipped (build ./ezbookkeeping and ./dist first)"
fi

echo
[ "$fails" -eq 0 ] && echo "DRILL PASSED" || echo "DRILL FAILED ($fails problem(s))"
exit "$fails"
