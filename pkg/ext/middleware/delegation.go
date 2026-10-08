// Package extmw contains the HTTP middleware of the ext module.
package extmw

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	extperm "github.com/mayswind/ezbookkeeping/pkg/ext/permissions"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// BusinessHeaderName is the request header with which a manager or staff member says whose business they work in.
// Without it every request operates on the caller's own data, exactly like upstream.
const BusinessHeaderName = "X-Business-Id"

const (
	contextRoleKey      = "EXT_ROLE"
	contextDelegatedKey = "EXT_DELEGATED"
)

// Decision is the outcome of resolving a request's business context
type Decision struct {
	Delegated bool
	OwnerUid  int64
	Role      extmodels.Role
}

// MembershipLookup finds the active membership of the caller in a business
type MembershipLookup func(ownerUid int64) (*extmodels.Membership, error)

// Decide resolves whose data a request operates on. It is a pure function so the security rules can be tested.
//   - no header, or the caller's own id: the caller's own data, no delegation
//   - identity routes (profile, password, tokens, 2FA, membership): header ignored, caller's own data
//   - otherwise: the caller must be an active member and the role must allow the method and path
func Decide(actualUid int64, headerValue string, method string, relativePath string, lookup MembershipLookup) (Decision, *errs.Error) {
	headerValue = strings.TrimSpace(headerValue)

	if headerValue == "" {
		return Decision{}, nil
	}

	ownerUid, err := strconv.ParseInt(headerValue, 10, 64)

	if err != nil || ownerUid <= 0 {
		return Decision{}, exterrs.ErrBusinessIdInvalid
	}

	if ownerUid == actualUid || extperm.IsSelfOnlyPath(relativePath) {
		return Decision{}, nil
	}

	membership, err := lookup(ownerUid)

	if err != nil {
		if e, ok := err.(*errs.Error); ok {
			return Decision{}, e
		}

		return Decision{}, errs.ErrOperationFailed
	}

	if !extperm.Allowed(membership.Role, method, relativePath) {
		return Decision{}, exterrs.ErrNotPermittedInBusiness
	}

	return Decision{Delegated: true, OwnerUid: ownerUid, Role: membership.Role}, nil
}

// Delegation lets managers and staff work on a business owner's data. It must run after the JWT authorization
// middleware and before the route handlers; write requests made in a business context are written to the audit log.
func Delegation() core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		actualUid := c.GetActualUid()

		decision, derr := Decide(actualUid, c.GetHeader(BusinessHeaderName), c.Request.Method, extperm.RelativePath(c.Request.URL.Path),
			func(ownerUid int64) (*extmodels.Membership, error) {
				return extservices.Memberships.GetActive(c, ownerUid, actualUid)
			})

		if derr != nil {
			log.Warnf(c, "[ext.delegation] user \"uid:%d\" denied for business \"%s\" on %s %s, because %s", actualUid, c.GetHeader(BusinessHeaderName), c.Request.Method, c.Request.URL.Path, derr.Message)
			utils.PrintJsonErrorResult(c, derr)

			return
		}

		if !decision.Delegated {
			c.Next()
			return
		}

		c.SetEffectiveUid(decision.OwnerUid)
		c.Set(contextDelegatedKey, true)
		c.Set(contextRoleKey, decision.Role)

		c.Next()

		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			return
		}

		err := extservices.Audit.Record(c, &extmodels.AuditLog{
			OwnerUid: decision.OwnerUid, ActorUid: actualUid, Role: decision.Role, Method: c.Request.Method,
			Path: c.Request.URL.Path, Status: c.Writer.Status(), ClientIp: c.ClientIP(),
		})

		if err != nil {
			log.Errorf(c, "[ext.delegation] failed to write audit log for user \"uid:%d\" in business \"uid:%d\", because %s", actualUid, decision.OwnerUid, err.Error())
		}
	}
}

// RoleOf returns the role of the caller in the current request: the delegated role, or owner when working on own data
func RoleOf(c *core.WebContext) extmodels.Role {
	if role, exists := c.Get(contextRoleKey); exists {
		if r, ok := role.(extmodels.Role); ok {
			return r
		}
	}

	return extmodels.RoleOwner
}

// IsDelegated returns whether the current request works on somebody else's business
func IsDelegated(c *core.WebContext) bool {
	delegated, exists := c.Get(contextDelegatedKey)

	return exists && delegated == true
}
