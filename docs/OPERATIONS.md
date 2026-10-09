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
