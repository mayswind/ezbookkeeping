# Extensions: staff roles, inventory, credit sales

Handoff document for the `ext` module. It records the plan, what is built, how it fits together, and what is left.
Branch: `feature/ext-staff-inventory`. Last updated: 2026-10-08.

## 1. Goal

Turn the upstream personal-finance app into a multi-user service for small businesses (and still serve individuals):

- **Owner / manager / staff** roles: staff log in with their own account and work on the owner's books.
- **Inventory** with a single store or several locations.
- **Sales** that automatically create transactions, including **sales on credit** and later **credit repayments**.
- All of it added so that **upstream (`mayswind/ezbookkeeping`) can still be merged** without constant conflicts.

Individuals are unaffected: with no business header every request behaves exactly like upstream.

## 2. Status

| Area | State |
|---|---|
| Backend: roles, delegation, audit log | Done, tested |
| Backend: locations, items, stock ledger | Done, tested |
| Backend: customers, sales, credit, repayments | Done, tested |
| HTTP API and routes | Done, tested end to end through the real middleware |
| Upstream seams (3 files, ~39 lines) | Done |
| Seam check script and CI workflow | Done (CI workflow untested on GitHub) |
| Frontend: foundation, business switcher, Team page, Inventory page | Done: type-check, lint, build and unit tests pass; **not yet clicked through in a browser** |
| Frontend: sales screen, customers and repayments | Not started (see section 9) |
| Paywall, signup codes, billing | Not started (separate track) |
| Reports (stock valuation, profit, ageing) | Not started |

Tests: `go test ./pkg/ext/...` (74 tests, all passing). The service tests take about 4 seconds each because every test boots a fresh SQLite database and syncs all tables, so the package needs several minutes; run a single test with `-run`. The whole backend suite passes except
`TestExchangeRatesApiLatestExchangeRateHandler_NationalBankOfUkraineDataSource`, which calls a live third-party API and fails
on `main` too (not related to this work).

## 3. Architecture

Everything new lives under `pkg/ext/` (no upstream file holds business logic). Package names differ from directory names on purpose to avoid clashing with upstream packages of the same name.

```
pkg/ext/
  ext.go            SyncTables()                 called by cmd/database.go
  routes.go         RegisterRoutes(), DelegationMiddleware()   called by cmd/webserver.go
  errors/           package exterrs   error codes, sub category 100
  models/           package extmodels database tables (all prefixed ext_)
  permissions/      package extperm   role permission matrix (deny by default)
  middleware/       package extmw     delegation + audit
  services/         package extservices  business logic
  api/              package extapi    HTTP handlers and request/response types
```

### The seams (the only upstream files changed)

| File | Change |
|---|---|
| `pkg/core/context_web.go` | `SetEffectiveUid`, `GetActualUid`, and `GetCurrentUid` returns the delegated owner's uid when set |
| `cmd/webserver.go` | one `Use(ext.DelegationMiddleware(...))` line and one `ext.RegisterRoutes(...)` call, both tagged `[ext]` |
| `cmd/database.go` | one `ext.SyncTables()` call, tagged `[ext]` |
| `src/router/desktop.ts` | one import and one `...extRoutes` line, tagged `[ext]` |
| `src/components/desktop/MainPageLayout.vue` | one import and one `<ext-top-nav />` line in the top toolbar, tagged `[ext]` |

`scripts/check-seams.sh [base]` fails if a branch changes any other upstream file or makes a seam grow past 40 lines.
Run it before merging and in CI (`.github/workflows/ext-ci.yml`).

## 4. Owner / manager / staff

- A **membership** (`ext_membership`) links an owner and a user: role `manager` or `staff`, status `pending`, `active` or `revoked`.
  The owner is simply the user who owns the data; there is no membership row for owners.
- The owner **invites an existing user by email**; the invitee must **accept** (`/ext/staff/respond.json`).
  Removal (by owner) or leaving (by member) revokes access immediately.
- A member works in a business by sending the header **`X-Business-Id: <owner uid>`**. The delegation middleware then:
  1. checks an active membership,
  2. checks the role against the permission matrix for that method and path,
  3. sets the *effective uid* to the owner, so every existing upstream handler that calls `GetCurrentUid()` operates on the owner's data,
  4. after the request, writes an **audit entry** for every non-GET request (including the response status).
- **Identity routes ignore the header** and always act as the logged-in person: `/users/*`, `/tokens/*`, `/systems/*`, `/ext/me/*`,
  `/ext/staff/*`, `/ext/audit/*`. A manager can never change the owner's password or tokens, or see or edit the owner's team.
- **Deny by default:** a route with no rule in `pkg/ext/permissions/permissions.go` is owner-only. New upstream endpoints stay closed
  to managers and staff until someone opens them there.

### Permission matrix (summary)

| Capability | Staff | Manager | Owner |
|---|---|---|---|
| Read accounts, categories, tags, templates, transaction list/get/count | yes | yes | yes |
| Record a sale at the listed price, record a repayment, add a customer | yes | yes | yes |
| Override a sale price or give a discount | no | yes | yes |
| Add a raw transaction (expense, transfer, income) | no | yes | yes |
| List items, stock, customers, sales, repayments | yes | yes | yes |
| Statistics, amounts, reports, all other reads | no | yes | yes |
| Modify / delete transactions | no | yes | yes |
| Add / modify / hide accounts, categories, tags, templates | no | yes | yes |
| Items, locations, stock receive / adjust / transfer, stock ledger | no | yes | yes |
| Modify / delete customers, void a sale | no | yes | yes |
| Delete accounts, categories, tags; data export, import, clear; tokens; custom icons; AI features | no | no | yes |
| Manage staff, read the audit log | no | no | yes (own business only) |

## 5. Inventory and stock

- **Locations**: a default location "Main" is created on first use, so single-store businesses never see the concept.
  A location can be renamed, and deleted only if empty and not the last one.
- **Items**: `sku` unique per business, `track_stock` flag (services and untracked goods never have stock), prices in minor currency units per whole unit.
  Deleting soft-deletes and renames the sku to `sku~id`, so the sku can be reused and history is preserved.
- **Stock ledger** (`ext_stock_movement`): immutable rows; **stock on hand is the sum of movements**. Reasons: opening, purchase, sale, adjustment, transfer out/in, sale void.
- **Quantities** are integers scaled by 1000 (`2500` = 2.5 units); prices are per whole unit; line total is rounded half up.
- Stock can never go negative (receive/adjust/transfer/sale are all checked inside the database transaction that writes them).

## 6. Sales, credit and repayments: how they hit the books

Accounting uses only existing upstream transaction types and the existing **Receivables** account category, so nothing in upstream tables or reports changes.

| Event | Ledger effect |
|---|---|
| Cash / paid part of a sale | **Income** transaction into the payment account, using the chosen income category |
| Credit part of a sale | **Income** transaction into a **Receivables** account (its balance is "what customers owe") |
| Credit repayment | **Transfer** from the Receivables account to the payment account (needs a transfer category) |
| Void a sale | Its transactions are deleted and the stock returns to the location |

Rules enforced: a walk-in customer (no customer) must pay in full; payment and receivable accounts must share a currency; the receivable account must be of category Receivables;
repayments cannot exceed what is owed; repayments pay the oldest open sale first unless a sale is chosen; a sale with repayments cannot be voided.
Per-customer balance = sum of `total - paid` over non-voided sales. Repayment allocations are stored (`ext_repayment_allocation`).

Sales are recorded in this order, with compensation if a step fails: (a) sale + lines + stock movements in one DB transaction
(row locks on the location and items on MySQL/PostgreSQL), (b) the finance transaction(s) through the upstream `TransactionService`,
(c) save the transaction ids on the sale. If (b) or (c) fails, the finance transactions are deleted and the sale rows removed, so the client sees a clean failure and a retry cannot double-sell.
Repayments create the transfer first and the DB rows second, and delete the transfer if the DB step fails.

**Void** is retryable: first the sale is flipped to voided and restocked in one DB transaction (which also refuses if repayments exist),
then the finance transactions are deleted one at a time, clearing each id on the sale. If a step fails, calling void again finishes the remaining work.
Stock returns to the sale's location, or to the default location if that location was deleted since.

## 7. API

All routes are under `/api/v1` and need a normal login token. Ids are strings, amounts are integers (minor units), quantities are scaled by 1000.
Send `X-Business-Id` to work in somebody else's business.

| Method and path | Min role | Purpose |
|---|---|---|
| GET `/ext/me/businesses.json` | any | My own business and those I was invited to |
| POST `/ext/staff/invite.json` `{email, role}` | owner | Invite an existing user as `manager` or `staff` |
| GET `/ext/staff/list.json` | owner | My team |
| POST `/ext/staff/set_role.json` `{staffUid, role}` | owner | Change a role |
| POST `/ext/staff/remove.json` `{staffUid}` | owner | Remove a member |
| POST `/ext/staff/respond.json` `{ownerUid, accept}` | any | Accept or decline an invitation |
| POST `/ext/staff/leave.json` `{ownerUid}` | any | Leave a business |
| GET `/ext/audit/list.json?limit&beforeId` | owner | What managers and staff changed |
| GET `/ext/locations/list.json` | staff | Locations |
| POST `/ext/locations/add.json`, `modify.json`, `delete.json` | manager | Manage locations |
| GET `/ext/items/list.json`, `/ext/items/stock.json?itemId&locationId` | staff | Items, stock on hand |
| POST `/ext/items/add.json`, `modify.json`, `delete.json` | manager | Manage items |
| POST `/ext/stock/receive.json`, `adjust.json`, `transfer.json` | manager | Change stock |
| GET `/ext/stock/movements/list.json` | manager | Stock ledger |
| GET `/ext/customers/list.json`, `balances.json` | staff | Customers and what they owe |
| POST `/ext/customers/add.json` | staff | Add a customer |
| POST `/ext/customers/modify.json`, `delete.json` | manager | Manage customers |
| POST `/ext/sales/add.json` | staff | Record a sale |
| GET `/ext/sales/list.json`, `get.json?id` | staff | Sales |
| POST `/ext/sales/void.json` `{id}` | manager | Cancel a sale |
| POST `/ext/repayments/add.json` | staff | Record a credit repayment |
| GET `/ext/repayments/list.json` | staff | Repayments |

Example sale on credit (500.00 paid now, the rest owed; amounts in minor units):

```json
POST /api/v1/ext/sales/add.json
X-Business-Id: 123456789
{
  "customerId": "42", "locationId": "0", "categoryId": "777",
  "lines": [{ "itemId": "9", "qty": 5000 }],
  "amountPaid": 50000, "paymentAccountId": "11", "receivableAccountId": "12"
}
```

Errors use upstream's format. Ext errors are in sub category `100` (`pkg/ext/errors/errors.go`).

## 8. Database

Tables (all `ext_` prefixed, created by `SyncTables()` during `ezbookkeeping database update` or on start when `auto_update_database` is true; additive only):

- User database (looked up across owners): `ext_membership`, `ext_audit_log`.
- Business data (next to the owner's data): `ext_location`, `ext_item`, `ext_stock_movement`, `ext_customer`, `ext_sale`, `ext_sale_line`, `ext_repayment`, `ext_repayment_allocation`.

Works on SQLite, MySQL and PostgreSQL through the same ORM as upstream. Tests run on SQLite only.
Deployment: nothing to change; the next deploy creates the tables. Back up the database first, as always.

## 9. What is left (next steps, in order)

1. **Frontend, remaining.** Built so far in `src/ext/` (no upstream file holds ext UI logic):
   - `api.ts` typed client over the app's shared axios, `business.ts` (the `X-Business-Id` header via its own axios interceptor, per-user remembered choice, switching reloads the page because every store holds the previous business's data), `qty.ts` (exact fixed-point quantities, unit tested).
   - `ExtTopNav.vue` (business menu in the toolbar, with an invitations dot), `TeamPage.vue` (invitations, businesses I work in, my team with roles, recent activity), `InventoryPage.vue` with item, stock (receive / adjust / transfer) and location dialogs. Managers and owners see the buttons, staff see a read-only view.
   - Routes `/ext/team` and `/ext/inventory`.

   Still to build: the **sale screen** (cart, customer, amount paid, account pickers, receipt), **customers and repayments** (balances, record a repayment), sales history with void, and pickers that default the payment, receivable and income-category accounts.
   Known gaps: prices are shown in the *user's* default currency because the API does not yet return the business currency; new strings use English text as the key and are not translated; the pages have been type-checked and built but not exercised in a browser, so expect layout fixes; the mobile app has no ext screens.
2. **Registration and invitations for new people.** Invitees must already have an account, and public registration is closed (paywall plan).
   Decide how a new staff member gets an account: invitation links that allow registration, or the owner creates the account.
3. **Billing / paywall** per owner; staff seats belong to the owner's subscription.
4. **Reports**: stock valuation, low stock (`reorder_level` exists, no endpoint yet), profit with cost of goods sold (`unit_cost` is stored on movements), receivables ageing, sales by item/day.
5. Returns / refunds, partial voids, purchase orders and suppliers (payables), sale-level price lists.
6. A "staff see only their own entries" mode (today staff can list all transactions of the business).

## 10. Known limitations and risks

Code review findings were fixed in this branch (see git history); what remains:

- **Not atomic across the books and the ext tables.** They are written in separate database transactions with compensating deletes.
  If the process dies, or a compensating delete itself fails, a sale can exist without its finance transaction (or the reverse).
  Failures are logged with the ids needed to repair them. Fixing this properly needs an exported upstream method that accepts an open session; that would be one more seam.
  A reconciliation report ("sales whose transactions are missing") would be a cheap safety net.
- **Row locking is only verified on SQLite.** On MySQL/PostgreSQL stock changes take `SELECT ... FOR UPDATE` on the location and item rows
  (always locations first, then items in ascending id order). That code path cannot run in the SQLite tests, so test concurrent sales on the real database before relying on it.
- Payment and receivable accounts must share a currency; there is one currency per sale.
- Staff cannot change prices or give discounts, and cannot post raw transactions. Both are role rules in code (`permissions.go`, `SaleCreateHandler`); make them configurable per business if owners want different policies.
- Staff can still list all of the business's transactions and sales (no "own entries only" mode yet).
- MCP routes and cron jobs are separate code paths and do not use delegation.
- Delegated requests include **API tokens**: a token belonging to a staff user can use the header too, with the same role limits.
- Prefix rules (`GET /accounts/`, `/transaction/`, `/transactions/`, `/ext/`) give managers read access to any *new* upstream GET route under those prefixes. Review the matrix after each upstream merge.
- Identity routes (`/users/`, `/tokens/`, `/systems/`, `/ext/me|staff|audit/`) deliberately ignore the business header and act as the logged-in person. A new upstream route under those prefixes that should be business-scoped would silently use the caller's own data.
- A location name can be reused after deletion, but deleting two locations with the same name within one second collides on the unique index.

## 11. Working on it

```sh
export PATH=$HOME/sdk/go/bin:$PATH GOTOOLCHAIN=local     # Go 1.27.1 (see go.mod)
go build ./...
go test ./pkg/ext/...                                    # the ext tests (SQLite, temp files)
scripts/check-seams.sh main                              # only registered seams may differ from upstream
```

### Pulling upstream

1. `git remote add upstream https://github.com/mayswind/ezbookkeeping.git` (once), `git config rerere.enabled true`.
2. Keep `main` a clean mirror of upstream where possible, merge it into this branch, never rebase it.
3. Likely conflicts are only the three seam files. Resolve by keeping upstream's code and re-adding the `[ext]` lines.
4. After merging: `go build ./... && go test ./pkg/ext/... && scripts/check-seams.sh upstream/main`, then re-read the permission matrix for new upstream routes.
5. If upstream changes `GetCurrentUid`, the authorization middleware, or the order of `apiV1Route.Use(...)`, re-check that delegation still runs after authorization and before the routes.

### Adding a feature to the module

Add models to `extmodels` (and to `GlobalTables`/`OwnerTables`), logic to `extservices`, handlers to `extapi`, routes to `routes.go`,
**and a rule to `permissions.go`** (a test fails if an ext route has no permission decision). Write a service test and, for anything security related, an HTTP test in `pkg/ext/http_test.go`.
