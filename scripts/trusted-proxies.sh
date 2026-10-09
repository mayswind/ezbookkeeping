#!/usr/bin/env sh
# Prints the value for EBK_SECURITY_TRUSTED_PROXY_IPS on Render: the app's default private networks, the machine's own IPv6
# loopback, and Cloudflare's published ranges. Cloudflare sits in front of Render, and without its ranges the app takes
# Cloudflare's server for the visitor, so everybody behind it shares one login-failure counter. See docs/OPERATIONS.md.
#
#   scripts/trusted-proxies.sh        then paste the line into the Render dashboard (and into render.yaml)
# Refresh it now and then: Cloudflare changes its list rarely, but it does change.

set -eu
DEFAULT="10.0.0.0/8,169.254.0.0/16,127.0.0.0/8,172.16.0.0/12,192.168.0.0/16,::1/128"
RANGES="$(curl -fsS https://www.cloudflare.com/ips-v4; echo; curl -fsS https://www.cloudflare.com/ips-v6)"
printf '%s' "$DEFAULT"
printf '%s\n' "$RANGES" | grep -E '^[0-9a-fA-F:.]+/[0-9]+$' | while read -r range; do printf ',%s' "$range"; done
echo
