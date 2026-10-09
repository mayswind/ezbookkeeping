# Operations checks

See also `docs/BACKUP.md` (backups and the restore drill).

## Login limits behind Render's proxy

The app slows down password guessing: after 5 failed password or token checks in a minute from one address (`max_failures_per_ip_per_minute`)
or for one account (`max_failures_per_user_per_minute`), further attempts are refused for that minute. The address comes from gin's `ClientIP`
using `trusted_proxy_ips` (`EBK_SECURITY_TRUSTED_PROXY_IPS`; default: private networks only), which works like this
(`pkg/ext/clientip_test.go` records it):

- Forwarding headers are believed only when the connection comes from a trusted proxy; otherwise they are ignored, so nobody can fake an address.
- `X-Forwarded-For` is read **from the right**, skipping trusted proxies; the first address that is not trusted is "the visitor". Extra addresses a visitor
  adds on the left are ignored.
- The risk: if the chain behind your host contains a proxy that is *not* trusted (for example a CDN's public address), that proxy's address is taken as the
  visitor. Everybody arriving through it then shares one counter, so five wrong passwords from anyone lock everyone out for a minute.

### Finding on 9 Oct 2026: every visitor was `::1`
Render's log showed every request, for every visitor, from `::1`: traffic reaches the app through something on the same machine, over IPv6 loopback.
The default trusted list only has `127.0.0.0/8` (IPv4 loopback), so `::1` was not trusted, the forwarding header was ignored, and all visitors shared one
failure counter (so five wrong passwords from anyone could block everyone's login attempts for a minute).
Fix: add `::1/128` to `EBK_SECURITY_TRUSTED_PROXY_IPS` (done in `render.yaml`; set the same value in the Render dashboard to apply it without waiting for a deploy):

    10.0.0.0/8,169.254.0.0/16,127.0.0.0/8,172.16.0.0/12,192.168.0.0/16,::1/128

Trusting loopback is safe: only a connection from the machine itself has that address, so an outside visitor cannot use it to fake an address
(`pkg/ext/clientip_test.go`). After the change, repeat the check below. If the log then shows your own address, it is fixed; if it shows a CDN's
address, add the CDN ranges (next section); if it still shows `::1`, the proxy sends no forwarding header and another approach is needed.

### Second finding, same day: Cloudflare is the next hop
After `::1` was trusted, the log showed `172.71.150.174`, which is inside Cloudflare's `172.64.0.0/13`. Cloudflare sits in front of Render, its
server is the rightmost address in `X-Forwarded-For` that is not trusted, so the app took it for the visitor, and everybody arriving through that server
still shared one counter. Fix: also trust Cloudflare's published ranges. `scripts/trusted-proxies.sh` prints the whole value (the default private
networks, `::1/128` and Cloudflare's IPv4 and IPv6 ranges, 28 entries); paste its output into `EBK_SECURITY_TRUSTED_PROXY_IPS` in the Render dashboard
(it is also in `render.yaml`). `pkg/ext/clientip_test.go` shows the real visitor is found with it and cannot be faked. Run the script again every few
months and update the setting if Cloudflare has added ranges; a visitor arriving through an unlisted Cloudflare range would again look like Cloudflare.

### How to check (about five minutes, nothing is changed)
1. Find your public address (search "what is my IP", or open https://ifconfig.me).
2. Log in to the live site, then open the service's **Logs** in Render and find the line for that login (`POST /api/authorize.json`). Request lines look like
   `[REQUEST] [id] 200 0 <user id> <ADDRESS> POST /api/authorize.json 170ms`; the address is the field before the method.
3. Read the result:
   - It equals your public address: correct. Nothing to do.
   - It is a private address (10.x, 172.16 to 172.31, 192.168.x): the real address is not being passed on; tell me and we will look at the headers.
   - It is a public address that is not yours (for example starting 104., 172.64 to 172.71, 162.158., 188.114. or 2606:4700:): that is a CDN's address and
     every visitor through it shares one counter. Fix below.
4. Repeat from another network (your phone on mobile data). The logged address should be different from the first. If it is the same, the counter is shared.
5. Optional, only when nobody else is using the site: enter a wrong password for a made-up user name six times within a minute. The sixth attempt should be
   refused ("failure count limit"). Then log in normally from the other network: it must still work. If it is refused too, the counter is shared.

### Fix if the address is a CDN's
Add the CDN's published ranges to the trusted list in Render: set `EBK_SECURITY_TRUSTED_PROXY_IPS` to the default list plus the ranges from
https://www.cloudflare.com/ips-v4 and https://www.cloudflare.com/ips-v6 (comma separated, no spaces), redeploy, and repeat the check. Refresh the list now and
then, since the CDN changes it. Only add ranges you trust: anything in the list can set the visitor's address.

## S3 versioning
Versioning on the backup bucket means deleted or overwritten files are kept as old versions. Litestream deletes old segments as part of its retention
cleanup, so with versioning on those deletions keep costing storage until a lifecycle rule expires them. Add a rule to the bucket that deletes
**non-current versions after 30 days** (and removes expired delete markers), so a leaked key or a bug can still be undone for a month without the bucket
growing forever. The restore drill does not test versioning; check the setting in the S3 console (Bucket > Properties > Bucket Versioning).
