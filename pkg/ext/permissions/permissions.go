// Package extperm decides what a business role may do. Paths are relative to "/api/v1".
// The matrix is deny-by-default: a route that is not listed is owner-only, so new upstream
// endpoints stay closed to managers and staff until someone explicitly opens them here.
package extperm

import (
	"strings"

	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// APIPrefix is the prefix of all authenticated API routes
const APIPrefix = "/api/v1"

// selfOnlyPrefixes are routes about the logged-in person (profile, security, tokens, business membership).
// A business context header is ignored for these so they can never operate on the owner's identity.
var selfOnlyPrefixes = []string{
	"/users/",
	"/tokens/",
	"/systems/",
	"/ext/me/",
	"/ext/staff/",
	"/ext/audit/",
}

type rule struct {
	method string
	path   string
	prefix bool // path is a prefix, otherwise an exact match
	min    extmodels.Role
}

func exact(method string, min extmodels.Role, paths ...string) []rule {
	rules := make([]rule, 0, len(paths))

	for _, p := range paths {
		rules = append(rules, rule{method: method, path: p, min: min})
	}

	return rules
}

func prefixed(method string, min extmodels.Role, paths ...string) []rule {
	rules := make([]rule, 0, len(paths))

	for _, p := range paths {
		rules = append(rules, rule{method: method, path: p, prefix: true, min: min})
	}

	return rules
}

var rules = buildRules()

func buildRules() []rule {
	var r []rule
	add := func(more []rule) { r = append(r, more...) }

	// ---- staff: take sales and look things up
	add(exact("GET", extmodels.RoleStaff,
		"/accounts/list.json", "/accounts/get.json",
		"/transaction/categories/list.json", "/transaction/categories/get.json",
		"/transaction/tags/list.json", "/transaction/tags/get.json",
		"/transaction/tags/groups/list.json", "/transaction/tags/groups/get.json",
		"/transaction/templates/list.json", "/transaction/templates/get.json",
		"/transactions/list.json", "/transactions/list/by_month.json", "/transactions/get.json", "/transactions/count.json",
		"/exchange_rates/latest.json",
		"/ext/locations/list.json", "/ext/people/list.json",
		"/ext/items/list.json", "/ext/items/stock.json",
		"/ext/customers/list.json", "/ext/customers/balances.json",
		"/ext/sales/list.json", "/ext/sales/get.json",
		"/ext/repayments/list.json", "/ext/repayments/get.json", "/ext/business/profile.json",
	))
	// Staff record money through sales and repayments only. A raw "add transaction" could be an expense or a transfer
	// between the owner's accounts, so it needs a manager.
	add(exact("POST", extmodels.RoleStaff,
		"/ext/customers/add.json",
		"/ext/sales/add.json",
		"/ext/repayments/add.json",
	))

	// ---- manager: run the business day to day
	add(prefixed("GET", extmodels.RoleManager,
		"/accounts/", "/transaction/", "/transactions/", "/insights/explorers/get", "/insights/explorers/list",
		"/ext/",
	))
	add(exact("GET", extmodels.RoleManager, "/data/statistics.json"))
	add(exact("POST", extmodels.RoleManager,
		"/transactions/add.json", "/transactions/modify.json", "/transactions/delete.json", "/transactions/batch_delete.json",
		"/accounts/add.json", "/accounts/modify.json", "/accounts/hide.json", "/accounts/move.json",
		"/transaction/categories/add.json", "/transaction/categories/add_batch.json", "/transaction/categories/modify.json",
		"/transaction/categories/hide.json", "/transaction/categories/move.json",
		"/transaction/tags/add.json", "/transaction/tags/add_batch.json", "/transaction/tags/modify.json",
		"/transaction/tags/hide.json", "/transaction/tags/move.json",
		"/transaction/templates/add.json", "/transaction/templates/modify.json",
		"/transaction/templates/hide.json", "/transaction/templates/move.json",
		"/ext/locations/add.json", "/ext/locations/modify.json", "/ext/locations/delete.json",
		"/ext/items/add.json", "/ext/items/modify.json", "/ext/items/delete.json",
		"/ext/stock/receive.json", "/ext/stock/adjust.json", "/ext/stock/transfer.json",
		"/ext/customers/modify.json", "/ext/customers/delete.json",
		"/ext/sales/void.json",
	))

	// Everything else (data export / import / clear, deleting accounts, categories and tags,
	// custom icons, exchange rate overrides, AI features, ...) is owner-only: no rule, so denied.
	return r
}

// IsSelfOnlyPath returns whether the path is about the logged-in person and must never be delegated
func IsSelfOnlyPath(path string) bool {
	for _, p := range selfOnlyPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

// Allowed returns whether the role may call the method and path. The owner may do anything.
func Allowed(role extmodels.Role, method string, path string) bool {
	if role >= extmodels.RoleOwner {
		return true
	}

	for _, r := range rules {
		if r.method != method {
			continue
		}

		if r.prefix && !strings.HasPrefix(path, r.path) {
			continue
		}

		if !r.prefix && path != r.path {
			continue
		}

		if role >= r.min {
			return true
		}
	}

	return false
}

// RelativePath strips the API prefix from a full request path
func RelativePath(fullPath string) string {
	return strings.TrimPrefix(fullPath, APIPrefix)
}
