# End-to-end check of the ext API against a running local server (default http://localhost:8080).
# It registers two throw-away users, so only run it against a development database.
# Usage: python3 scripts/ext-smoke.py
import json, sys, time, urllib.request, urllib.error
BASE = "http://localhost:8080/api/"
sfx = str(int(time.time()))[-6:]
fails = []

def call(method, path, body=None, token=None, business=None):
    req = urllib.request.Request(BASE + path, method=method, data=json.dumps(body).encode() if body is not None else None)
    req.add_header("Content-Type", "application/json")
    req.add_header("X-Timezone-Offset", "60"); req.add_header("X-Timezone-Name", "Africa/Lagos")
    if token: req.add_header("Authorization", "Bearer " + token)
    if business: req.add_header("X-Business-Id", business)
    try:
        with urllib.request.urlopen(req) as r: return r.status, json.loads(r.read())
    except urllib.error.HTTPError as e:
        try: return e.code, json.loads(e.read())
        except Exception: return e.code, {}

def check(name, cond, extra=""):
    print(("PASS " if cond else "FAIL ") + name + (" " + str(extra) if (extra and not cond) else ""))
    if not cond: fails.append(name)

def register(name):
    s, b = call("POST", "register.json", {"username": name, "email": name + "@example.com", "nickname": name.upper(), "password": "secret123", "language": "en", "defaultCurrency": "NGN", "firstDayOfWeek": 0, "categories": []})
    assert s == 200, (s, b)
    return b["result"]["token"]

owner_name, staff_name = "smkown" + sfx, "smkstf" + sfx
owner = register(owner_name); staff = register(staff_name); mgr_name = "smkmgr" + sfx; mgr = register(mgr_name)
o = lambda m, p, b=None: call(m, "v1/" + p, b, owner)

def account(name, category):
    s, b = o("POST", "accounts/add.json", {"name": name, "category": category, "type": 1, "icon": "1", "iconType": 0, "color": "000000", "currency": "NGN", "balance": "0", "comment": ""})
    assert s == 200, (s, b); return b["result"]["id"]
def category(name, ctype, parent="0"):
    s, b = o("POST", "transaction/categories/add.json", {"name": name, "type": ctype, "parentId": parent, "icon": "1", "iconType": 0, "color": "000000", "comment": ""})
    assert s == 200, (s, b); return b["result"]["id"]
def balance(acc):
    s, b = o("GET", "accounts/list.json?visible_only=false")
    return next(int(a["balance"]) for a in b["result"] if a["id"] == acc)

cash, recv = account("Cash", 1), account("Customers owe", 6)
inc = category("Sales", 1, category("Sales group", 1)); xfer = category("Repayments", 3, category("Repay group", 3))
s, b = o("POST", "ext/items/add.json", {"sku": "RICE", "name": "Rice 1kg", "unit": "kg", "costPrice": 70000, "salePrice": 150000, "reorderLevel": 0, "trackStock": True})
check("add item", s == 200, b); item = b["result"]["id"]
s, b = o("POST", "ext/stock/receive.json", {"itemId": item, "locationId": "0", "qty": 10000, "unitCost": 70000})
check("receive stock", s == 200, b)
s, b = o("POST", "ext/customers/add.json", {"name": "Ada Obi"}); check("add customer", s == 200, b); cust = b["result"]["id"]

sale_req = lambda **kw: {"locationId": "0", "customerId": "0", "time": int(time.time()), "utcOffset": 60, "lines": [{"itemId": item, "qty": 2500}], "discount": 0,
                          "amountPaid": 0, "paymentAccountId": "0", "receivableAccountId": "0", "categoryId": inc, "note": "smoke", **kw}
# 2.5 kg x 1500.00 = 3750.00 -> 375000 minor units
s, b = o("POST", "ext/sales/add.json", sale_req(customerId=cust, amountPaid=100000, paymentAccountId=cash, receivableAccountId=recv))
check("credit sale accepted", s == 200, b); sale = b["result"]
check("sale total = 375000", sale["total"] == 375000, sale); check("sale owes 275000", sale["outstanding"] == 275000, sale)
check("cash balance 100000", balance(cash) == 100000, balance(cash)); check("receivables balance 275000", balance(recv) == 275000, balance(recv))
s, b = o("GET", "ext/items/stock.json"); check("stock now 7.5", b["result"][0]["qty"] == 7500, b)
s, b = o("GET", "ext/customers/balances.json"); check("customer owes 275000", b["result"][0]["outstanding"] == 275000, b)

s, b = o("POST", "ext/sales/add.json", sale_req(amountPaid=0)); check("walk-in on credit refused", s == 400, (s, b))
s, b = o("POST", "ext/sales/add.json", sale_req(lines=[{"itemId": item, "qty": 99000}], amountPaid=1, paymentAccountId=cash)); check("oversell refused", s == 400, (s, b))

s, b = o("POST", "ext/repayments/add.json", {"customerId": cust, "amount": 75000, "saleId": "0", "time": int(time.time()), "utcOffset": 60, "paymentAccountId": cash, "receivableAccountId": recv, "categoryId": xfer, "note": ""})
check("repayment accepted", s == 200, b); rep_id = b["result"]["id"]; check("receivables now 200000", balance(recv) == 200000, balance(recv)); check("cash now 175000", balance(cash) == 175000, balance(cash))
s, b = o("POST", "ext/sales/void.json", {"id": sale["id"]}); check("void refused after repayment", s == 400, (s, b))

# what receipts rely on
s, b = o("GET", "ext/repayments/get.json?id=" + rep_id); check("repayment with its allocations", s == 200 and b["result"]["amount"] == 75000 and len(b["result"]["allocations"]) == 1 and b["result"]["paymentAccountId"] == cash, (s, b))
s, b = o("GET", "ext/business/profile.json"); check("receipt name falls back to the owner's name", s == 200 and b["result"]["name"] == owner_name.upper() and b["result"]["receiptName"] == "", (s, b))
s, b = o("POST", "ext/me/business_profile/update.json", {"receiptName": "Ada Stores", "address": "12 Market Road", "phone": "0800", "footer": "Thank you"}); check("save receipt details", s == 200 and b["result"]["name"] == "Ada Stores", (s, b))

# legal pages, the recorded acceptance and the business export
import urllib.request as _u
status = _u.urlopen("http://localhost:8080/legal/terms.html").status; page = _u.urlopen("http://localhost:8080/legal/privacy.html").read().decode()
check("legal pages are public", status == 200 and "Privacy Policy" in page, status)
s, b = o("GET", "ext/me/settings.json"); check("terms not accepted yet", s == 200 and b["result"]["acceptedTermsVersion"] == "", (s, b))
s, b = o("POST", "ext/me/terms/accept.json", {"version": "smoke-1"}); check("accept the terms", s == 200, (s, b))
s, b = o("GET", "ext/me/settings.json"); check("acceptance is recorded", s == 200 and b["result"]["acceptedTermsVersion"] == "smoke-1" and b["result"]["acceptedTermsTime"] > 0, (s, b))
import io, zipfile
def export_names(token, business=None):
    req = urllib.request.Request(BASE + "v1/ext/export/business.zip", headers={"Authorization": "Bearer " + token, **({"X-Business-Id": business} if business else {})})
    with urllib.request.urlopen(req) as r:
        data = r.read(); ctype = r.headers.get("Content-Type")
    z = zipfile.ZipFile(io.BytesIO(data)); return ctype, {n: z.read(n).decode("utf-8-sig") for n in z.namelist()}
ctype, files = export_names(owner)
check("owner's export is a zip with the business in it", ctype == "application/zip" and "RICE" in files["items.csv"] and "Ada Obi" in files["customers.csv"] and "README.txt" in files, ctype)
check("export amounts are exact", "375000" in files["sales.csv"] and "2.5" in files["sale_lines.csv"], files["sales.csv"][:200])

# what the customers page relies on
s, b = o("GET", "ext/sales/list.json?customerId=" + cust + "&onlyOpen=true"); check("open sales of a customer", s == 200 and len(b["result"]) == 1 and b["result"][0]["outstanding"] == 200000, (s, b))
s, b = o("GET", "ext/repayments/list.json?customerId=" + cust); check("repayment history of a customer", s == 200 and len(b["result"]) == 1 and b["result"][0]["amount"] == 75000, (s, b))
s, b = o("POST", "ext/customers/modify.json", {"id": cust, "name": "Ada Obi Jr", "phone": "0800", "email": "", "note": ""}); check("edit customer", s == 200 and b["result"]["name"] == "Ada Obi Jr" and b["result"]["outstanding"] == 200000, (s, b))
s, b = o("POST", "ext/customers/delete.json", {"id": cust}); check("cannot delete a customer who owes", s == 400, (s, b))

# per-person business features setting (server side)
s, b = o("GET", "ext/me/settings.json"); check("features off until chosen", s == 200 and (b["result"]["businessFeatures"], b["result"]["configured"]) == (False, False), (s, b))
s, b = o("POST", "ext/me/settings/update.json", {"businessFeatures": True}); check("turn features on", s == 200 and b["result"]["businessFeatures"] is True, (s, b))
s, b = o("GET", "ext/me/settings.json"); check("features stay on", s == 200 and (b["result"]["businessFeatures"], b["result"]["configured"]) == (True, True), (s, b))
s, b = call("GET", "v1/ext/me/settings.json", None, staff); check("another person's setting is separate", s == 200 and b["result"]["businessFeatures"] is False, (s, b))

# staff
uid_owner = call("GET", "v1/ext/me/businesses.json", None, owner)[1]["result"][0]["ownerUid"]
s, b = o("POST", "ext/staff/invite.json", {"email": staff_name + "@example.com", "role": "staff"}); check("invite staff", s == 200, b)
st = lambda m, p, b=None, biz=uid_owner: call(m, "v1/" + p, b, staff, biz)
s, b = call("GET", "v1/ext/me/businesses.json", None, staff); check("staff sees pending invitation", any(x["status"] == "pending" for x in b["result"]), b)
s, b = call("POST", "v1/ext/staff/respond.json", {"ownerUid": uid_owner, "accept": True}, staff); check("staff accepts", s == 200, b)
s, b = st("GET", "ext/items/list.json"); check("staff lists owner's items", s == 200 and len(b["result"]) == 1, (s, b))
s, b = st("POST", "ext/sales/add.json", sale_req(amountPaid=150000 * 1, lines=[{"itemId": item, "qty": 1000}], paymentAccountId=cash))
check("staff records a cash sale", s == 200, (s, b)); staff_sale = b["result"]
check("sale belongs to owner, actor is staff", staff_sale["actorUid"] != uid_owner, staff_sale)
# who did it: the books say so, and so do the records
s, b = o("GET", "transactions/get.json?id=" + staff_sale["paidTransactionId"]); check("staff sale is marked in the books", s == 200 and ("by " + staff_name.upper()) in b["result"]["comment"], (s, b))
s, b = o("GET", "transactions/get.json?id=" + sale["paidTransactionId"]); check("owner sale carries no mark", s == 200 and " by " not in b["result"]["comment"], (s, b))
s, b = o("GET", "ext/people/list.json"); check("people: owner first, then staff", s == 200 and [p["role"] for p in b["result"]] == ["owner", "staff"] and b["result"][1]["name"] == staff_name.upper(), (s, b))
s, b = st("GET", "ext/people/list.json"); check("staff can see who works there", s == 200 and len(b["result"]) == 2, (s, b))
s, b = o("GET", "ext/audit/list.json"); check("audit names the sale", s == 200 and any(e["action"] == "sales.add" and e["entityId"] == staff_sale["id"] for e in b["result"]), b)

# a manager adds a transaction through the app's own API: it is booked to the owner and marked with their name
s, b = o("POST", "ext/staff/invite.json", {"email": mgr_name + "@example.com", "role": "manager"}); check("invite manager", s == 200, b)
s, b = call("POST", "v1/ext/staff/respond.json", {"ownerUid": uid_owner, "accept": True}, mgr); check("manager accepts", s == 200, b)
s, b = call("POST", "v1/transactions/add.json", {"type": 2, "categoryId": inc, "time": int(time.time()), "utcOffset": 60, "sourceAccountId": cash, "destinationAccountId": "0",
            "sourceAmount": 5000, "destinationAmount": 0, "hideAmount": False, "tagIds": [], "pictureIds": [], "comment": "petty cash"}, mgr, uid_owner)
check("manager adds a transaction to the owner's books", s == 200, (s, b))
if s == 200:
    s, b = o("GET", "transactions/get.json?id=" + b["result"]["id"]); check("manager's transaction is marked with their name", s == 200 and b["result"]["comment"] == "petty cash · by " + mgr_name.upper(), (s, b))

# reports (6.5 kg left: 6.5 x 700.00 cost, 6.5 x 1500.00 selling price)
s, b = o("GET", "ext/reports/stock_value.json"); check("stock value report", s == 200 and b["result"]["totalCostValue"] == 455000 and b["result"]["totalRetailValue"] == 975000 and b["result"]["rows"][0]["qty"] == 6500, (s, b))
s, b = o("GET", "ext/reports/low_stock.json"); check("nothing low without a reorder level", s == 200 and b["result"] == [], (s, b))
s, b = o("GET", "ext/reports/receivables.json"); check("who owes what", s == 200 and b["result"]["totalOutstanding"] == 200000 and b["result"]["current"] == 200000 and len(b["result"]["rows"]) == 1 and b["result"]["rows"][0]["customer"]["name"] == "Ada Obi Jr", (s, b))
s, b = st("GET", "ext/reports/stock_value.json"); check("staff cannot read reports (403)", s == 403, (s, b))

s, b = st("POST", "ext/sales/add.json", sale_req(discount=100, amountPaid=1, lines=[{"itemId": item, "qty": 1000}], paymentAccountId=cash)); check("staff discount refused (403)", s == 403, (s, b))
s, b = st("POST", "ext/sales/void.json", {"id": staff_sale["id"]}); check("staff void refused (403)", s == 403, (s, b))
s, b = st("POST", "ext/items/add.json", {"sku": "X", "name": "X", "unit": "", "costPrice": 0, "salePrice": 0, "reorderLevel": 0, "trackStock": False}); check("staff add item refused (403)", s == 403, (s, b))
s, b = st("POST", "ext/customers/modify.json", {"id": cust, "name": "Hacked", "phone": "", "email": "", "note": ""}); check("staff edit customer refused (403)", s == 403, (s, b))
s, b = st("POST", "ext/customers/add.json", {"name": "Walk up"}); check("staff can add a customer", s == 200, (s, b))
s, b = st("POST", "ext/me/settings/update.json", {"businessFeatures": True}); check("staff changes only their own setting", s == 200, (s, b))
s, b = o("GET", "ext/me/settings.json"); check("owner setting unaffected by staff", s == 200 and b["result"]["businessFeatures"] is True, (s, b))
s, b = st("GET", "ext/business/profile.json"); check("staff print the owner's receipt details", s == 200 and b["result"]["name"] == "Ada Stores" and b["result"]["footer"] == "Thank you", (s, b))
s, b = st("GET", "ext/repayments/get.json?id=" + rep_id); check("staff can reprint a repayment receipt", s == 200, (s, b))
ctype, files = export_names(staff, uid_owner); check("a business header cannot fetch the owner's export", "RICE" not in files["items.csv"] and "Ada Obi" not in files["customers.csv"], files["items.csv"][:120])
s, b = st("GET", "ext/staff/list.json"); check("staff sees no team of the owner", s == 200 and b["result"] == [], (s, b))
s, b = o("GET", "ext/audit/list.json"); check("owner audit shows staff sale", s == 200 and any(e["method"] == "POST" and "sales/add" in e["path"] for e in b["result"]), b)
s, b = call("GET", "v1/ext/items/list.json", None, staff, "999999"); check("staff cannot use a business they don't belong to", s == 403, (s, b))
s, b = o("POST", "ext/sales/void.json", {"id": staff_sale["id"]}); check("owner voids the staff sale", s == 200, (s, b))
# the app's own "Clear All Data" also clears the business records, but only when it succeeds
s, b = o("POST", "data/clear/all.json", {"password": "not-the-password"}); check("clear all data refused with a wrong password", s == 400, (s, b))
s, b = o("GET", "ext/items/list.json"); check("a refused clear-all leaves the business alone", s == 200 and len(b["result"]) == 1, (s, b))
s, b = st("POST", "data/clear/all.json", {"password": "secret123"}); check("a staff member cannot clear the owner's data (403)", s == 403, (s, b))
s, b = o("POST", "data/clear/all.json", {"password": "secret123"}); check("clear all data accepted", s == 200, (s, b))
s, b = o("GET", "ext/items/list.json"); check("items are gone", s == 200 and b["result"] == [], (s, b))
s, b = o("GET", "ext/customers/list.json"); check("customers are gone", s == 200 and b["result"] == [], (s, b))
s, b = o("GET", "ext/sales/list.json"); check("sales are gone", s == 200 and b["result"] == [], (s, b))
s, b = o("GET", "ext/reports/receivables.json"); check("nobody owes anything", s == 200 and b["result"]["totalOutstanding"] == 0, (s, b))
s, b = o("GET", "ext/audit/list.json"); check("the clearing is on the activity log", s == 200 and b["result"][0]["action"] == "data.clear_all", b)
s, b = o("GET", "ext/staff/list.json"); check("the team survives", s == 200 and len(b["result"]) == 2, (s, b))
s, b = o("GET", "ext/business/profile.json"); check("receipt details survive", s == 200 and b["result"]["name"] == "Ada Stores", (s, b))

print("\n%d failure(s)" % len(fails)); sys.exit(1 if fails else 0)
