#!/usr/bin/env sh
# Fails when a branch changes upstream files other than the known "seams".
# Everything the ext module adds lives in new files, so upstream merges stay cheap.
#
# Usage: scripts/check-seams.sh [base-ref]     (default base: main)
# Run it before opening a pull request and in CI.

BASE="${1:-main}"

# Upstream files the ext module is allowed to touch. Keep each change tiny and tagged with "[ext]".
SEAMS="
cmd/webserver.go
cmd/database.go
pkg/core/context_web.go
src/router/desktop.ts
src/components/desktop/MainPageLayout.vue
src/views/desktop/settings/SettingsPageLayout.vue
src/views/desktop/LoginPage.vue
src/views/desktop/SignupPage.vue
src/views/mobile/LoginPage.vue
src/views/mobile/SignupPage.vue
"

# Paths that are ours, not upstream's
OWN_PREFIXES="
pkg/ext/
legal/
public/legal/
docs/
scripts/
deploy/
src/ext/
.github/workflows/deploy.yml
.github/workflows/ext-ci.yml
render.yaml
FORK.md
"

if ! git rev-parse --verify --quiet "$BASE" >/dev/null; then
    echo "Base ref \"$BASE\" does not exist" >&2
    exit 2
fi

status=0

for file in $(git diff --name-only "$BASE"...HEAD); do
    allowed=0

    for seam in $SEAMS; do
        [ "$file" = "$seam" ] && allowed=1
    done

    for prefix in $OWN_PREFIXES; do
        case "$file" in
            "$prefix"*) allowed=1 ;;
        esac
    done

    if [ "$allowed" -eq 0 ]; then
        echo "NOT ALLOWED: $file is an upstream file and is not a registered seam"
        status=1
    fi
done

# A seam change must stay small
for seam in $SEAMS; do
    changed=$(git diff --numstat "$BASE"...HEAD -- "$seam" | awk '{print $1 + $2}')

    if [ -n "$changed" ] && [ "$changed" -gt 40 ]; then
        echo "TOO BIG: $seam changed $changed lines; seams must stay tiny"
        status=1
    fi
done

if [ "$status" -eq 0 ]; then
    echo "Seam check passed"
fi

exit $status
