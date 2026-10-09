package ext

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmw "github.com/mayswind/ezbookkeeping/pkg/ext/middleware"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	extperm "github.com/mayswind/ezbookkeeping/pkg/ext/permissions"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
	"github.com/mayswind/ezbookkeeping/pkg/ext/testdb"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

type world struct {
	t       *testing.T
	router  *gin.Engine
	owner   int64
	manager int64
	staff   int64
	other   int64
}

// newWorld wires the real delegation middleware and ext routes behind a stand-in for JWT authorization:
// the "X-Test-Uid" header plays the role of the logged-in user.
func newWorld(t *testing.T) *world {
	t.Helper()
	gin.SetMode(gin.TestMode)

	config := testdb.Config(t)
	require.NoError(t, datastore.InitializeDataStore(config))
	require.NoError(t, uuid.InitializeUuidGenerator(config))
	require.NoError(t, datastore.Container.UserStore.SyncStructs(new(models.User)))
	require.NoError(t, SyncTables())

	w := &world{t: t}
	c := core.NewNullContext()
	newUser := func(name string) int64 {
		user := &models.User{Username: name, Email: name + "@example.com", Nickname: name, Language: "en", DefaultCurrency: "NGN"}
		require.NoError(t, services.Users.CreateUser(c, user, true))

		return user.Uid
	}
	w.owner, w.manager, w.staff, w.other = newUser("owner"), newUser("manager"), newUser("staff"), newUser("other")

	for uid, role := range map[int64]extmodels.Role{w.manager: extmodels.RoleManager, w.staff: extmodels.RoleStaff} {
		_, err := extservices.Memberships.Invite(c, w.owner, map[int64]string{w.manager: "manager@example.com", w.staff: "staff@example.com"}[uid], role)
		require.NoError(t, err)
		require.NoError(t, extservices.Memberships.Respond(c, uid, w.owner, true))
	}

	bindApi := func(fn core.ApiHandlerFunc) gin.HandlerFunc {
		return func(ginCtx *gin.Context) {
			wc := core.WrapWebContext(ginCtx, nil)
			result, err := fn(wc)

			if err != nil {
				utils.PrintJsonErrorResult(wc, err)
			} else {
				utils.PrintJsonSuccessResult(wc, result)
			}
		}
	}
	bindMiddleware := func(fn core.MiddlewareHandlerFunc) gin.HandlerFunc {
		return func(ginCtx *gin.Context) { fn(core.WrapWebContext(ginCtx, nil)) }
	}

	w.router = gin.New()
	v1 := w.router.Group("/api/v1")
	v1.Use(func(ginCtx *gin.Context) { // stand-in for JWT authorization
		uid, _ := strconv.ParseInt(ginCtx.GetHeader("X-Test-Uid"), 10, 64)
		wc := core.WrapWebContext(ginCtx, nil)
		wc.SetTokenClaims(&core.UserTokenClaims{Uid: uid})
	})
	v1.Use(DelegationMiddleware(bindMiddleware))
	RegisterRoutes(v1, bindApi)

	// stand-ins for upstream routes: one data route and one identity route, both reporting whose uid they see
	whoami := func(wc *core.WebContext) (any, *errs.Error) {
		return map[string]int64{"current": wc.GetCurrentUid(), "actual": wc.GetActualUid()}, nil
	}
	v1.GET("/accounts/list.json", bindApi(whoami))
	v1.POST("/transactions/delete.json", bindApi(whoami))
	v1.POST("/transactions/add.json", bindApi(func(wc *core.WebContext) (any, *errs.Error) {
		var req struct {
			Comment string `json:"comment"`
		}

		if err := wc.ShouldBindJSON(&req); err != nil {
			return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}

		return map[string]any{"id": "77", "comment": req.Comment, "uid": wc.GetCurrentUid()}, nil
	}))
	v1.POST("/data/clear/all.json", bindApi(whoami))
	v1.POST("/users/profile/update.json", bindApi(whoami))

	return w
}

func (w *world) call(uid int64, businessId int64, method string, path string, body any) (int, map[string]any) {
	w.t.Helper()

	var reader *bytes.Reader

	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(w.t, err)
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, "/api/v1"+path, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-Uid", strconv.FormatInt(uid, 10))

	if businessId != 0 {
		req.Header.Set(extmw.BusinessHeaderName, strconv.FormatInt(businessId, 10))
	}

	rec := httptest.NewRecorder()
	w.router.ServeHTTP(rec, req)

	var parsed map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &parsed)

	return rec.Code, parsed
}

func result(t *testing.T, parsed map[string]any) map[string]any {
	t.Helper()

	data, ok := parsed["result"].(map[string]any)
	require.True(t, ok, "unexpected response: %v", parsed)

	return data
}

func TestHTTP_StaffWorkOnTheOwnersData(t *testing.T) {
	w := newWorld(t)

	code, body := w.call(w.staff, w.owner, "GET", "/accounts/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, float64(w.owner), result(t, body)["current"], "data routes see the owner's uid")
	assert.Equal(t, float64(w.staff), result(t, body)["actual"], "the real caller is still known")
}

func TestHTTP_WithoutHeaderEveryoneStaysOnTheirOwnData(t *testing.T) {
	w := newWorld(t)

	code, body := w.call(w.staff, 0, "GET", "/accounts/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, float64(w.staff), result(t, body)["current"])
}

func TestHTTP_IdentityRoutesIgnoreTheBusinessHeader(t *testing.T) {
	w := newWorld(t)

	code, body := w.call(w.manager, w.owner, "POST", "/users/profile/update.json", map[string]string{})
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, float64(w.manager), result(t, body)["current"], "a manager must never edit the owner's profile")
}

func TestHTTP_RolesAreEnforcedOnEveryRoute(t *testing.T) {
	w := newWorld(t)

	cases := []struct {
		name    string
		uid     int64
		method  string
		path    string
		allowed bool
	}{
		{"staff reads accounts", w.staff, "GET", "/accounts/list.json", true},
		{"staff cannot delete transactions", w.staff, "POST", "/transactions/delete.json", false},
		{"staff cannot post raw transactions", w.staff, "POST", "/transactions/add.json", false},
		{"manager can post raw transactions", w.manager, "POST", "/transactions/add.json", true},
		{"manager can delete transactions", w.manager, "POST", "/transactions/delete.json", true},
		{"manager cannot clear all data", w.manager, "POST", "/data/clear/all.json", false},
		{"staff cannot clear all data", w.staff, "POST", "/data/clear/all.json", false},
		{"staff cannot add items", w.staff, "POST", "/ext/items/add.json", false},
		{"staff can list items", w.staff, "GET", "/ext/items/list.json", true},
		{"stranger cannot read accounts", w.other, "GET", "/accounts/list.json", false},
		{"stranger cannot list items", w.other, "GET", "/ext/items/list.json", false},
	}

	for _, tc := range cases {
		code, _ := w.call(tc.uid, w.owner, tc.method, tc.path, map[string]any{})

		if tc.allowed {
			assert.Equal(t, http.StatusOK, code, tc.name)
		} else {
			assert.Equal(t, http.StatusForbidden, code, tc.name)
		}
	}
}

func TestHTTP_InventoryFlowThroughRoles(t *testing.T) {
	w := newWorld(t)

	// the manager sets up the catalog and stock in the owner's business
	code, body := w.call(w.manager, w.owner, "POST", "/ext/items/add.json", map[string]any{"sku": "A1", "name": "Rice", "salePrice": 1500, "trackStock": true})
	require.Equal(t, http.StatusOK, code, body)
	itemId := result(t, body)["id"].(string)

	code, body = w.call(w.manager, w.owner, "POST", "/ext/stock/receive.json", map[string]any{"itemId": itemId, "qty": 10000, "unitCost": 700})
	require.Equal(t, http.StatusOK, code, body)
	assert.Equal(t, strconv.FormatInt(w.manager, 10), result(t, body)["actorUid"], "the manager, not the owner, is recorded as the actor")

	// staff can see stock, and the owner sees the same data without any header
	code, body = w.call(w.staff, w.owner, "GET", "/ext/items/stock.json", nil)
	require.Equal(t, http.StatusOK, code)
	levels := body["result"].([]any)
	require.Len(t, levels, 1)
	assert.Equal(t, float64(10000), levels[0].(map[string]any)["qty"])

	code, body = w.call(w.owner, 0, "GET", "/ext/items/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Len(t, body["result"].([]any), 1)

	// another user's business is empty and separate
	code, body = w.call(w.other, 0, "GET", "/ext/items/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Empty(t, body["result"].([]any))
}

func TestHTTP_WritesInABusinessAreAudited(t *testing.T) {
	w := newWorld(t)

	code, _ := w.call(w.manager, w.owner, "POST", "/ext/locations/add.json", map[string]any{"name": "Shop 2"})
	require.Equal(t, http.StatusOK, code)
	code, _ = w.call(w.manager, w.owner, "GET", "/ext/locations/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	code, _ = w.call(w.staff, w.owner, "POST", "/ext/items/add.json", map[string]any{"sku": "X", "name": "X"}) // forbidden
	require.Equal(t, http.StatusForbidden, code)

	code, body := w.call(w.owner, 0, "GET", "/ext/audit/list.json", nil)
	require.Equal(t, http.StatusOK, code)

	entries := body["result"].([]any)
	require.Len(t, entries, 1, "reads and rejected requests are not part of the audit log of writes")
	entry := entries[0].(map[string]any)
	assert.Equal(t, strconv.FormatInt(w.manager, 10), entry["actorUid"])
	assert.Equal(t, "manager", entry["role"])
	assert.Equal(t, "POST", entry["method"])
	assert.Equal(t, float64(200), entry["status"])
}

func TestHTTP_RemovedStaffLoseAccessImmediately(t *testing.T) {
	w := newWorld(t)

	code, _ := w.call(w.staff, w.owner, "GET", "/ext/items/list.json", nil)
	require.Equal(t, http.StatusOK, code)

	code, _ = w.call(w.owner, 0, "POST", "/ext/staff/remove.json", map[string]any{"staffUid": strconv.FormatInt(w.staff, 10)})
	require.Equal(t, http.StatusOK, code)

	code, _ = w.call(w.staff, w.owner, "GET", "/ext/items/list.json", nil)
	assert.Equal(t, http.StatusForbidden, code)
}

func TestHTTP_StaffCannotManageStaffEvenWithTheHeader(t *testing.T) {
	w := newWorld(t)

	// staff management is an identity route: the header is ignored, so the manager only ever manages their own (empty) team
	code, body := w.call(w.manager, w.owner, "GET", "/ext/staff/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Empty(t, body["result"].([]any), "a manager must not see or edit the owner's team")

	code, _ = w.call(w.manager, w.owner, "POST", "/ext/staff/invite.json", map[string]any{"email": "other@example.com", "role": "manager"})
	require.Equal(t, http.StatusOK, code)

	code, body = w.call(w.owner, 0, "GET", "/ext/staff/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Len(t, body["result"].([]any), 2, "the owner's team is unchanged by what the manager did")
}

func TestHTTP_InvalidBusinessHeader(t *testing.T) {
	w := newWorld(t)

	req := httptest.NewRequest("GET", "/api/v1/accounts/list.json", nil)
	req.Header.Set("X-Test-Uid", strconv.FormatInt(w.staff, 10))
	req.Header.Set(extmw.BusinessHeaderName, "not-a-number")
	rec := httptest.NewRecorder()
	w.router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// Every ext route must either be an identity route or be reachable by at least a manager, so a new route
// cannot silently end up owner-only because somebody forgot the permission matrix.
func TestEveryExtRouteHasAPermissionDecision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("/api/v1")
	RegisterRoutes(group, func(core.ApiHandlerFunc) gin.HandlerFunc { return func(*gin.Context) {} })

	count := 0

	for _, route := range router.Routes() {
		path := extperm.RelativePath(route.Path)
		count++

		if extperm.IsSelfOnlyPath(path) {
			continue
		}

		assert.True(t, extperm.Allowed(extmodels.RoleManager, route.Method, path), "%s %s has no permission rule", route.Method, path)
	}

	assert.Greater(t, count, 30)
}

func TestHTTP_StaffCannotOverridePricesOrGiveDiscounts(t *testing.T) {
	w := newWorld(t)
	price := 1
	line := map[string]any{"itemId": "1", "qty": 1000}
	priced := map[string]any{"itemId": "1", "qty": 1000, "unitPrice": price}

	for name, body := range map[string]map[string]any{
		"discount":       {"categoryId": "1", "lines": []any{line}, "discount": 100, "amountPaid": 1},
		"price override": {"categoryId": "1", "lines": []any{priced}, "amountPaid": 1},
	} {
		code, resp := w.call(w.staff, w.owner, "POST", "/ext/sales/add.json", body)
		assert.Equal(t, http.StatusForbidden, code, name)
		assert.Equal(t, float64(exterrs.ErrNotPermittedForRole.Code()), resp["errorCode"], name)
	}

	// a manager passes the role check and then fails on the unknown item, proving the check is by role
	code, _ := w.call(w.manager, w.owner, "POST", "/ext/sales/add.json", map[string]any{"categoryId": "1", "lines": []any{priced}, "amountPaid": 1})
	assert.NotEqual(t, http.StatusForbidden, code)
}

func TestHTTP_BusinessFeaturesSettingIsPerPersonAndIgnoresTheBusinessHeader(t *testing.T) {
	w := newWorld(t)

	code, body := w.call(w.manager, 0, "GET", "/ext/me/settings.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, false, result(t, body)["businessFeatures"], "off until the person chooses")
	assert.Equal(t, false, result(t, body)["configured"])

	// a manager changing their own setting while a business header is present must not touch the owner's
	code, body = w.call(w.manager, w.owner, "POST", "/ext/me/settings/update.json", map[string]any{"businessFeatures": true})
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, true, result(t, body)["businessFeatures"])

	code, body = w.call(w.manager, 0, "GET", "/ext/me/settings.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, true, result(t, body)["businessFeatures"])
	assert.Equal(t, true, result(t, body)["configured"])

	code, body = w.call(w.owner, 0, "GET", "/ext/me/settings.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, false, result(t, body)["businessFeatures"], "the owner's setting is untouched")
	assert.Equal(t, false, result(t, body)["configured"])

	// switching off again is an explicit, remembered choice
	code, _ = w.call(w.manager, 0, "POST", "/ext/me/settings/update.json", map[string]any{"businessFeatures": false})
	require.Equal(t, http.StatusOK, code)
	code, body = w.call(w.manager, 0, "GET", "/ext/me/settings.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Equal(t, false, result(t, body)["businessFeatures"])
	assert.Equal(t, true, result(t, body)["configured"])
}

func TestHTTP_TransactionsAddedByAManagerAreMarkedWithTheirName(t *testing.T) {
	w := newWorld(t)

	code, body := w.call(w.manager, w.owner, "POST", "/transactions/add.json", map[string]any{"comment": "lunch", "sourceAmount": 100})
	require.Equal(t, http.StatusOK, code, body)
	assert.Equal(t, "lunch · by manager", result(t, body)["comment"], "the upstream handler sees the marked comment")
	assert.Equal(t, float64(w.owner), result(t, body)["uid"], "and it is booked to the owner")

	code, body = w.call(w.manager, w.owner, "POST", "/transactions/add.json", map[string]any{"sourceAmount": 100})
	require.Equal(t, http.StatusOK, code, body)
	assert.Equal(t, "by manager", result(t, body)["comment"], "a transaction without a comment still gets the mark")

	// the owner's own entries are left exactly as typed
	code, body = w.call(w.owner, 0, "POST", "/transactions/add.json", map[string]any{"comment": "lunch"})
	require.Equal(t, http.StatusOK, code, body)
	assert.Equal(t, "lunch", result(t, body)["comment"])

	// and so is a manager's entry in their own books
	code, body = w.call(w.manager, 0, "POST", "/transactions/add.json", map[string]any{"comment": "mine"})
	require.Equal(t, http.StatusOK, code, body)
	assert.Equal(t, "mine", result(t, body)["comment"])
}

func TestHTTP_AuditLogSaysWhatWasDoneAndToWhich(t *testing.T) {
	w := newWorld(t)

	code, body := w.call(w.manager, w.owner, "POST", "/transactions/add.json", map[string]any{"comment": "x"})
	require.Equal(t, http.StatusOK, code, body)

	code, body = w.call(w.manager, w.owner, "POST", "/transactions/delete.json", map[string]any{"id": "123"})
	require.Equal(t, http.StatusOK, code, body)

	code, body = w.call(w.manager, w.owner, "POST", "/ext/locations/add.json", map[string]any{"name": "Shop 2"})
	require.Equal(t, http.StatusOK, code, body)
	locationId := result(t, body)["id"].(string)

	code, body = w.call(w.owner, 0, "GET", "/ext/audit/list.json", nil)
	require.Equal(t, http.StatusOK, code)

	entries := body["result"].([]any) // newest first
	require.Len(t, entries, 3)

	location := entries[0].(map[string]any)
	assert.Equal(t, "locations.add", location["action"])
	assert.Equal(t, "locations", location["entityType"])
	assert.Equal(t, locationId, location["entityId"], "created things are found in the response")

	deleted := entries[1].(map[string]any)
	assert.Equal(t, "transactions.delete", deleted["action"])
	assert.Equal(t, "123", deleted["entityId"], "changed things are found in the request")

	added := entries[2].(map[string]any)
	assert.Equal(t, "transactions.add", added["action"])
	assert.Equal(t, "77", added["entityId"])
	assert.Equal(t, strconv.FormatInt(w.manager, 10), added["actorUid"])
}

func TestHTTP_PeopleAreListedForEveryRole(t *testing.T) {
	w := newWorld(t)

	for _, uid := range []int64{w.owner, w.manager, w.staff} {
		business := int64(0)
		if uid != w.owner {
			business = w.owner
		}

		code, body := w.call(uid, business, "GET", "/ext/people/list.json", nil)
		require.Equal(t, http.StatusOK, code)

		people := body["result"].([]any)
		require.Len(t, people, 3)
		assert.Equal(t, "owner", people[0].(map[string]any)["role"])
		assert.Equal(t, strconv.FormatInt(w.owner, 10), people[0].(map[string]any)["uid"])
	}

	code, _ := w.call(w.other, w.owner, "GET", "/ext/people/list.json", nil)
	assert.Equal(t, http.StatusForbidden, code, "strangers cannot list a business's people")
}

func TestHTTP_ReportsAreForManagersAndOwners(t *testing.T) {
	w := newWorld(t)

	for _, path := range []string{"/ext/reports/stock_value.json", "/ext/reports/low_stock.json", "/ext/reports/receivables.json"} {
		code, _ := w.call(w.owner, 0, "GET", path, nil)
		assert.Equal(t, http.StatusOK, code, "owner "+path)

		code, _ = w.call(w.manager, w.owner, "GET", path, nil)
		assert.Equal(t, http.StatusOK, code, "manager "+path)

		code, _ = w.call(w.staff, w.owner, "GET", path, nil)
		assert.Equal(t, http.StatusForbidden, code, "staff "+path)
	}
}

func TestHTTP_MyPeopleIsAlwaysAboutTheCallersOwnBusiness(t *testing.T) {
	w := newWorld(t)

	// a manager working in the owner's business still gets their own business (just themselves)
	code, body := w.call(w.manager, w.owner, "GET", "/ext/staff/people.json", nil)
	require.Equal(t, http.StatusOK, code)
	people := body["result"].([]any)
	require.Len(t, people, 1)
	assert.Equal(t, strconv.FormatInt(w.manager, 10), people[0].(map[string]any)["uid"])

	// while the business-scoped list is the owner's team
	code, body = w.call(w.manager, w.owner, "GET", "/ext/people/list.json", nil)
	require.Equal(t, http.StatusOK, code)
	assert.Len(t, body["result"].([]any), 3)
}
