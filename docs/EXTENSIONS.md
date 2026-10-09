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
| Frontend: sales screen (cart, cash / part / credit, sales history, void) | Done: type-check, lint, build, unit tests and an API smoke test pass; **not yet clicked through in a browser** |
| Frontend: customers page (balances, history, edit, **repayments**) | Done: type-check, lint, build, unit tests and an API smoke test pass; **not yet clicked through in a browser** |
| Frontend: opt-in Business Features setting, i18n of all labels | Done |
| Reports: stock value, low stock, who owes what (with CSV download) | Done: backend tests pass on SQLite, API smoke test passes, UI **not yet clicked through in a browser** |
| Terms of Service and Privacy Policy pages, recorded acceptance, notices on login and signup | Done, see `docs/LEGAL.md`; **the wording is generic and `legal/details.json` must be filled in before deploying** |
| Export of business data (ZIP of CSV files) | Done, see below |
| Receipts: printable sale and repayment receipts (80 mm or A4), Save as PDF, copy as text, business details on them | Done: unit tests, API smoke test; **printing not yet tried on a real printer or in a browser** |
| Who did what: "by <name>" mark in transaction comments, "recorded by" columns, stock history, activity log with the thing changed | Done, same caveat |
| Paywall, signup codes, billing | Not started (separate track) |
| Reports (stock valuation, profit, ageing) | Not started |

Tests: `go test ./pkg/ext/...` (133 tests, all passing on SQLite). The service tests take several seconds each because every test boots a fresh database and syncs all tables, so the package needs several minutes; run a single test with `-run`.

**Other databases:** the same tests run unchanged on PostgreSQL or MySQL by setting `EXT_TEST_DB` (`postgres` or `mysql`), `EXT_TEST_DB_HOST`, `EXT_TEST_DB_USER` and `EXT_TEST_DB_PASSWD` (see `pkg/ext/testdb`); each test gets its own throw-away database. Results so far: PostgreSQL 17 passed all 54 service tests and the HTTP suites, and MySQL 8.4 passed the HTTP suites and 45 of 54 service tests with no failures (that run was stopped early). The tests added since (reports, attribution, people) have not yet been run on PostgreSQL or MySQL. The whole backend suite passes except
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
| `src/components/desktop/MainPageLayout.vue` | imports and tags for `<ext-top-nav />` (toolbar), `<ext-profile-menu-items />` (avatar menu) and `<ext-business-banner />` (top of page), tagged `[ext]` |
| `src/views/desktop/settings/SettingsPageLayout.vue` | one import and one `<ext-settings-nav-item />` line (the "Business Features" settings entry), tagged `[ext]` |

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
| GET `/ext/me/settings.json`, POST `/ext/me/settings/update.json` `{businessFeatures}` | any | My ext preferences (business features on or off); always about me, never about a business |
| POST `/ext/staff/invite.json` `{email, role}` | owner | Invite an existing user as `manager` or `staff` |
| GET `/ext/staff/list.json` | owner | My team |
| POST `/ext/staff/set_role.json` `{staffUid, role}` | owner | Change a role |
| POST `/ext/staff/remove.json` `{staffUid}` | owner | Remove a member |
| POST `/ext/staff/respond.json` `{ownerUid, accept}` | any | Accept or decline an invitation |
| POST `/ext/staff/leave.json` `{ownerUid}` | any | Leave a business |
| GET `/ext/audit/list.json?limit&beforeId` | owner | What managers and staff changed |
| GET `/ext/export/business.zip` | owner | ZIP of CSV files with the business records (always the caller's own business) |
| POST `/ext/me/terms/accept.json` `{version}` | any | Record that I accepted a version of the Terms and Privacy Policy (settings also return the latest) |
| GET `/ext/business/profile.json` | staff | What the business prints on receipts (name falls back to the owner's name) |
| GET `/ext/me/business_profile.json`, POST `/ext/me/business_profile/update.json` | any | My own business's receipt details |
| GET `/ext/repayments/get.json?id` | staff | One repayment with the sales it paid off |
| GET `/ext/people/list.json` | staff | The owner and everyone who has worked in the business (name, role, still active) |
| GET `/ext/staff/people.json` | any | The same for my own business, whichever business I am working in |
| GET `/ext/reports/stock_value.json`, `low_stock.json` (`?locationId`), `receivables.json` | manager | Reports |
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

### Who did what

Every record made through the ext module remembers its person: sales, repayments and stock movements have an `actor`, and the audit log
says what a manager or staff member did and to which sale, transaction, item and so on (`action`, `entityType`, `entityId`, found in the
response for things created and in the request for things changed).

For the books themselves, a transaction recorded by somebody other than the owner gets **"by <name>" at the end of its comment**:
- sales and repayments recorded by staff or managers (both the paid and the on-credit transaction): `Sale #12: weekend · by Sam`,
- a transaction added by a manager through the app's own "add transaction" call (the delegation middleware edits the `comment` of the request
  before it reaches the handler, using `StampText`, which shortens the original text and never the mark, up to the 255 character limit).

So the attribution shows everywhere the comment does, including the app's own transaction list and exports, with no change to upstream screens.
**No mark means the owner did it**: every non-owner write goes through delegation, which always marks. Limits: it is plain text, so
a manager can edit the comment later and remove it (the audit log still has the entry); only *adding* a transaction is marked, not editing one;
and transactions made before this feature are unmarked.

Screens: the sales history has a **Recorded by** column, the customer dialog **Received by** on repayments, each stock-tracked item a **History**
button (managers) listing movements with who and why, and the Team page's **Recent activity** reads like "Sam (Staff) - Recorded a sale #45".
Names come from `GET /ext/people/list.json` (the owner, members and removed members, so old records keep a name) and, for the Team page,
`GET /ext/staff/people.json`, which is always about the caller's own business.

### Receipts

A receipt opens automatically when a sale or a repayment is recorded, and from a **Receipt** button on every sale in the history and every
repayment in a customer's dialog (so any receipt can be reprinted).
- **Content** (`src/ext/receipt.ts`, plain data first and so unit tested): the business name, address and phone, receipt number (`#` and the sale
  or repayment id), date, customer, location (only when there are several), who served, where the money went, the lines (`2.5 kg x 1,500.00`),
  subtotal and discount when there is one, total, paid, **balance owed** when part is on credit, a note for credit or cancelled sales, a diagonal
  VOID mark for voided sales, and the footer message. A repayment receipt lists the sales it paid off, the amount received and what is still owed.
- **Print:** the dialog offers a small receipt for 80 mm thermal printers or a full A4 page. *Print* copies the receipt into its own element and hides
  everything else while the browser's print dialog is open, with a matching `@page` size, so only the receipt comes out; *Save as PDF* is in the same dialog.
  The paper choice is remembered in the browser.
- **Copy as text** puts a 32 (small) or 42 (A4) column plain-text version on the clipboard for WhatsApp, SMS or email, amounts lined up on the right.
- **Business details** (name on receipts, address, phone, footer message) are edited by the owner in *Settings > Business Features > Receipt details* and
  stored in `ext_business_profile` (one row per business). Staff and managers print with the owner's details (`GET /ext/business/profile.json`);
  `GET/POST /ext/me/business_profile...` always act on the caller's own business. Without a name the owner's name is used.
- A repayment receipt can be fetched again with `GET /ext/repayments/get.json?id` (it includes the sales it paid off and the accounts).
- Limits: the amount of "still owed" on a repayment receipt is the balance *at the time it is printed*, so a reprint after later sales or
  repayments shows the newer balance; the sale receipt shows what that sale still owes now. Receipts are not stored as documents; they are
  rebuilt from the sale each time, so a reprint always reflects the sale as it is now (including a later void).

### Export of business data

*Settings > Business Features > Export your business data* downloads a ZIP of CSV files for the owner's own business (`GET /ext/export/business.zip`):
items, locations, stock on hand, every stock movement, customers (with what they owe), sales and their lines, repayments and what each paid off, the team,
the activity log and the receipt details, plus a `README.txt` explaining the columns. It covers what the app's own export (accounts, categories, tags,
transactions) leaves out.
- Amounts are whole numbers in minor currency units (columns end in `_minor_units`), quantities are decimals, times are UTC, and names are written out next
  to ids (including for deleted items and people who left) so the files read on their own. Files start with a byte order mark so Excel shows accents correctly.
- Free text that starts with `=`, `+`, `-` or `@` gets a quote in front, so a customer name typed as a formula cannot run in a spreadsheet.
- It is **owner-only and always about the caller's own business**: the route is on the identity list, so a manager's business header cannot redirect it to the
  owner's records (tested at HTTP level and in the live smoke test).
- It streams in pages of 1000 rows (`exportChunkSize`), so a large business does not have to fit in memory; the response is `no-store`.

### Reports

All read-only and computed from the stock ledger and the sales, so they cannot disagree with them (managers and owners; staff get a 403).
- **Stock value** (`GET /ext/reports/stock_value.json?locationId`): per stock-tracked item, quantity x current cost price and x selling price, totals,
  and the stock per location. Items without stock are left out. Valuation uses the item's *current* prices, not what was paid for each purchase.
- **Low stock** (`GET /ext/reports/low_stock.json?locationId`): items at or below their reorder level (items without a level are never listed), most
  urgent first, with how much they are short by.
- **Who owes what** (`GET /ext/reports/receivables.json`): every customer with unpaid sales, largest debt first, aged from the sale date into up to 30,
  31 to 60, 61 to 90 and over 90 days, with the oldest unpaid sale. Voided sales and repaid amounts are excluded; the total always equals the sum of customer balances.

The Reports page (`/ext/reports`, chart icon in the toolbar for managers and owners) has a tab for each, a location filter for the stock reports, and **Download CSV**.
CSV cells that start with `=`, `+`, `-` or `@` get a quote in front so a spreadsheet never runs them as formulas.

## 8. Database

Tables (all `ext_` prefixed, created by `SyncTables()` during `ezbookkeeping database update` or on start when `auto_update_database` is true; additive only):

- User database (looked up across owners): `ext_terms_acceptance` (one row per acceptance), `ext_membership`, `ext_audit_log` (with `action`, `entity_type`, `entity_id`, added later as extra columns), `ext_user_setting` (one row per person: whether the business features are on).
- Business data (next to the owner's data): `ext_business_profile` (receipt details), `ext_location`, `ext_item`, `ext_stock_movement`, `ext_customer`, `ext_sale`, `ext_sale_line`, `ext_repayment`, `ext_repayment_allocation`.

Works on SQLite, MySQL and PostgreSQL through the same ORM as upstream. Tests run on SQLite only.
Deployment: nothing to change; the next deploy creates the tables. Back up the database first, as always.

## 9. What is left (next steps, in order)

1. **Frontend, remaining.** Built so far in `src/ext/` (no upstream file holds ext UI logic):
   - `api.ts` typed client over the app's shared axios, `business.ts` (the `X-Business-Id` header via its own axios interceptor, per-user remembered choice, switching reloads the page because every store holds the previous business's data), `qty.ts` (exact fixed-point quantities, unit tested).
   - `ExtTopNav.vue` (business menu in the toolbar, with an invitations dot), `TeamPage.vue` (invitations, businesses I work in, my team with roles, recent activity), `InventoryPage.vue` with item, stock (receive / adjust / transfer) and location dialogs. Managers and owners see the buttons, staff see a read-only view.
   - Routes `/ext/team` and `/ext/inventory`.

   `SalesPage.vue` (route `/ext/sales`, toolbar cart button): pick items, edit quantities (exact fixed-point, checked against stock at the chosen location), choose paid in full / part payment / on credit, customer (quick add), payment account, receivables account (only accounts of type Receivables in the payment account's currency) and income category, then complete the sale; recent sales with void for managers. Choices are remembered per business. Managers and owners can override line prices and give discounts; staff cannot (the server enforces it too). Totals use `src/ext/money.ts`, which mirrors the server's rounding and is unit tested.

   `CustomersPage.vue` (route `/ext/customers`): customers with what each owes (largest debt first), search, an "only people who owe me" filter, add / edit / delete (delete is refused while a customer owes money), a detail dialog with unpaid sales and repayment history, and a **Record repayment** dialog: amount (defaults to everything owed, never more), apply to the oldest sales or one chosen sale, the account the money arrived in, the owed-money (Receivables) account, and a transfer category (repayments are recorded as transfers, so the books need one). Choices are remembered per business.

   **Opt-in:** Sales, Customers, Inventory and Team are hidden until the person switches on *Settings > Business Features* (`BusinessSettingsPage.vue`, stored **on the server** in the new `ext_user_setting` table so it follows the person to every device; `features.ts` keeps a browser copy only to avoid a flicker while the page loads). Being invited to work in somebody else's business turns the features on automatically while the person is a member. Team also shows when there is an invitation to answer. The routes are guarded the same way (`router.ts`) and send people to the settings page.  A choice made before the setting moved to the server (kept in the browser) is carried over once.

   **Translations:** all labels are registered in `src/ext/locales/en.json`, with the English sentence as the key like the app's own `en.json`. Other languages: add `src/ext/locales/<code>.json` (underscore in the file name for a hyphenated code, e.g. `zh_Hans.json`) with the same keys; missing keys fall back to English. After adding a label run `python3 scripts/ext-extract-i18n.py`; a unit test (`src/ext/__tests__/locales.test.ts`) fails if a label or placeholder is not registered.

   Still to build: a sale detail view (a printable receipt exists), sending receipts by email. Possible next steps for reports: sales by day / by person, profit using cost of goods sold, charts.
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

API smoke test against a running local server: `python3 scripts/ext-smoke.py` (66 checks: stock, credit sale, repayment, void rules, staff role limits, audit log). It creates throw-away users, so use a development database.

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
