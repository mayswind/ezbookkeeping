# Backups and the restore drill

Production runs on Render's free plan with SQLite. The disk is wiped on every restart, so the database lives in S3:
Litestream (`deploy/render/`) restores the latest copy when the service starts and copies every change to the bucket
(`ezbookkeeping-sqlite-bkup`, path `ezbookkeeping/`, region `eu-north-1`) while it runs.
Uploaded files (receipts, pictures) are separate: they go to the same bucket under `uploads/` and are not covered by Litestream's history.

## The drill

`scripts/restore-drill.sh` proves the backup can actually be restored. Run it after any change to the backup setup and about once a month.
It is **read-only**: it downloads the backup into a temporary folder, restores it, checks it, boots the app on a copy, then deletes everything
(the temporary folder holds the S3 credentials and a copy of the customers' data). It never writes to the bucket or touches the live service.

```sh
scripts/restore-drill.sh            # uses LITESTREAM_ACCESS_KEY_ID / _SECRET_ACCESS_KEY from .env
DRILL_AS_OF=2026-10-08T07:20:00Z scripts/restore-drill.sh   # restore the database as it was at that time (UTC)
```

It checks: a backup generation exists; the restore finishes; SQLite's integrity check passes; the restored database holds users;
the app starts on the restored data (running its own upgrade step) and serves pages; the data is unchanged by that step.

### Results of the first drill (9 Oct 2026)
- Restore from S3: 13 to 36 seconds for a database of about 220 KB, replaying 22 copied segments.
- Integrity check: ok. Contents: 1 user, 1 account, 3 transactions, 23 categories.
- The current code started on the restored (older) production data and added the 12 new `ext_` tables by itself, with all data intact.
- Point-in-time restore worked: restoring to a moment before two of the three transactions existed showed exactly one.

## What the drill does not cover
- Whether replication is *currently* keeping up. The last copy in the bucket was 08 Oct 22:03 UTC. An idle database writes nothing, so that is normal
  when nobody used the app overnight, but it is also what a stopped replication looks like. To check: make any change in the live app (add a
  transaction), wait about ten seconds, and run the drill again: the "last change copied to S3" time should jump to now. Do this after every deploy.
- Restoring after Render loses the whole service: the drill used the same S3 keys the service uses, so it proves those keys can read the backup,
  not that a fresh deploy would find them. Keep a copy of the keys and of `deploy/render/litestream.yml` somewhere outside Render.
- Large databases: restore time grows with the size and with the number of segments since the last snapshot.

## Risks and settings worth changing
- **History is short.** Litestream's default keeps about 24 hours of history, so a mistake noticed after two days cannot be rolled back.
  Adding `retention: 168h` (7 days) and `snapshot-interval: 24h` under the replica in `deploy/render/litestream.yml` costs almost nothing.
- **The service holds keys that can delete its own backups** (Litestream needs delete rights to tidy old data). Turn on S3 **versioning** for the bucket
  (and optionally a lifecycle rule that keeps old versions for 30 days) so a bug or a leaked key cannot erase history.
- **A missing or unreadable backup must stop the start**: `entrypoint.sh` runs with `set -e`, so a credentials error fails the deploy (good). But if the
  bucket path were changed by mistake the service would start with an empty database and begin a new history. Never change `path:` in `litestream.yml` casually.
- The free plan sleeps without traffic; a restart restores from S3 each time, so every restart is a small restore drill, and anything written in the last
  second or so before a crash may be lost.
