# Terms of Service and Privacy Policy

The app has a Terms of Service page and a Privacy Policy page, a recorded acceptance by every user, and a notice on the login and signup screens.
**The wording is a generic starting point, not legal advice. Have a lawyer read it for your country before customers rely on it.**

## What exists
- `legal/terms.template.html`, `legal/privacy.template.html`: the wording, with `{{placeholders}}`.
- `legal/details.json`: your company name, address, contact email, country, governing law, hosting region and a few policy choices. **This is the only
  file you normally edit.** Values that still start with `REPLACE` show highlighted in yellow on the pages.
- `scripts/build-legal.py`: builds `public/legal/terms.html`, `public/legal/privacy.html` and `src/ext/legalVersion.ts` from the two. The built files are
  committed. `--check --strict` fails while any `REPLACE` is left or the built files are out of date; the deploy workflow runs it, so the pages cannot go
  live unfinished.
- The pages are public at `/legal/terms.html` and `/legal/privacy.html` (also linked from the login and signup screens and from the profile menu).
- Acceptance: after logging in, a person who has not accepted the current version sees a box they must accept (or they are logged out). Each acceptance is
  stored (`ext_terms_acceptance`: who, which version, when, from which address) and never overwritten, so there is a history.

## Filling it in
1. Edit `legal/details.json`: replace every `REPLACE...`.
2. `python3 scripts/build-legal.py` then open `public/legal/terms.html` and `privacy.html` in a browser and read both.
3. `python3 scripts/build-legal.py --check --strict` must print "legal pages are up to date and complete".
4. Commit `legal/details.json`, `public/legal/*` and `src/ext/legalVersion.ts` together.

## Facts the Privacy Policy states that you must keep true
| Statement | Where it comes from | Check |
|---|---|---|
| Backups and uploaded files are in Stockholm (eu-north-1) | the S3 bucket `ezbookkeeping-sqlite-bkup` | still the case if you move the bucket |
| Up to `backup_days` (7) of backup history | `retention: 168h` in `deploy/render/litestream.yml` | deployed, and S3 versioning lifecycle also expires old versions |
| Deleted data is gone within `deleted_data_retention_days` (30) | backups expire after 7 days; S3 versions after the lifecycle rule | set the bucket rule to expire non-current versions after 30 days |
| Hosting region | the Render service | write the region shown in the Render dashboard |
| Server log retention | Render's log retention | write what Render keeps |
| No marketing email, no analytics or advertising trackers | the app sends only password reset mail; no tracker is in the pages | revisit if either changes |
| Email is sent through Amazon SES | the SMTP settings | update if you change provider |

## Changing the wording later
Edit the templates or `details.json`, **raise `version`** (for example `2027-02-01`) and rebuild. Everybody is asked to accept again the next time they open
the app. A change that is only a typo fix can keep the same version.

## Paid plans
`paid_plans` is `false` while the service is free; the Terms then say charges will be announced `notice_days` ahead. When you start charging set it to
`true` and fill in `refund_policy`, raise the version, rebuild.

## Not covered yet
- **Deleting an account on request.** The Privacy Policy promises deletion within 30 days. The app has no self-service "delete my account" and the
  command-line tool is not available on Render's free plan, so today a deletion has to be done by a developer with access to the database. Plan a way
  to do it before customers ask.
- **"Clear all data" does not clear business records.** The app's own clear-data button removes accounts and transactions but not items, stock, sales or
  customers. A full erasure must cover the `ext_` tables too.
- The pages are English only.
