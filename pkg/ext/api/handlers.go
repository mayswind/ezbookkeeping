package extapi

import (
	"github.com/gin-gonic/gin/binding"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// Handlers groups all ext HTTP handlers
type Handlers struct{}

// API is the handler singleton
var API = &Handlers{}

func bindBody(c *core.WebContext, req any) *errs.Error {
	if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
		return errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	return nil
}

func bindQuery(c *core.WebContext, req any) *errs.Error {
	if err := c.ShouldBindQuery(req); err != nil {
		return errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	return nil
}

func fail(c *core.WebContext, where string, err error) *errs.Error {
	final := errs.Or(err, errs.ErrOperationFailed)

	if !errs.IsCustomError(err) {
		log.Errorf(c, "[ext.%s] failed, because %s", where, err.Error())
	}

	return final
}

// ---- business membership (identity routes: always act as the logged-in person)

// MyBusinessesHandler lists the businesses the caller can work in, starting with their own
func (h *Handlers) MyBusinessesHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetActualUid()
	result := []*BusinessView{}

	if me, err := services.Users.GetUserById(c, uid); err == nil {
		result = append(result, &BusinessView{OwnerUid: uid, Name: me.Nickname, Role: extmodels.RoleOwner.String(), Status: "owner"})
	}

	memberships, err := extservices.Memberships.ListByStaff(c, uid)

	if err != nil {
		return nil, fail(c, "me.businesses", err)
	}

	for _, m := range memberships {
		name := ""

		if owner, err := services.Users.GetUserById(c, m.OwnerUid); err == nil {
			name = owner.Nickname
		}

		status := "active"

		if m.Status == extmodels.MembershipStatusPending {
			status = "pending"
		}

		result = append(result, &BusinessView{OwnerUid: m.OwnerUid, Name: name, Role: m.Role.String(), Status: status})
	}

	return result, nil
}

// StaffInviteHandler invites an existing user to work in the caller's business
func (h *Handlers) StaffInviteHandler(c *core.WebContext) (any, *errs.Error) {
	var req InviteStaffRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	membership, err := extservices.Memberships.Invite(c, c.GetActualUid(), req.Email, roleFromString(req.Role))

	if err != nil {
		return nil, fail(c, "staff.invite", err)
	}

	return h.staffView(c, membership), nil
}

// StaffListHandler lists the caller's staff
func (h *Handlers) StaffListHandler(c *core.WebContext) (any, *errs.Error) {
	memberships, err := extservices.Memberships.ListByOwner(c, c.GetActualUid())

	if err != nil {
		return nil, fail(c, "staff.list", err)
	}

	result := make([]*StaffView, 0, len(memberships))

	for _, m := range memberships {
		result = append(result, h.staffView(c, m))
	}

	return result, nil
}

func (h *Handlers) staffView(c *core.WebContext, m *extmodels.Membership) *StaffView {
	view := &StaffView{StaffUid: m.StaffUid, Role: m.Role.String(), Status: "active"}

	if m.Status == extmodels.MembershipStatusPending {
		view.Status = "pending"
	}

	if user, err := services.Users.GetUserById(c, m.StaffUid); err == nil {
		view.Username, view.Nickname, view.Email = user.Username, user.Nickname, user.Email
	}

	return view
}

// StaffSetRoleHandler changes the role of a member
func (h *Handlers) StaffSetRoleHandler(c *core.WebContext) (any, *errs.Error) {
	var req StaffRoleRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Memberships.SetRole(c, c.GetActualUid(), req.StaffUid, roleFromString(req.Role)); err != nil {
		return nil, fail(c, "staff.set_role", err)
	}

	return true, nil
}

// StaffRemoveHandler removes a member from the caller's business
func (h *Handlers) StaffRemoveHandler(c *core.WebContext) (any, *errs.Error) {
	var req StaffRemoveRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Memberships.Remove(c, c.GetActualUid(), req.StaffUid); err != nil {
		return nil, fail(c, "staff.remove", err)
	}

	return true, nil
}

// StaffRespondHandler accepts or declines an invitation addressed to the caller
func (h *Handlers) StaffRespondHandler(c *core.WebContext) (any, *errs.Error) {
	var req RespondInviteRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Memberships.Respond(c, c.GetActualUid(), req.OwnerUid, req.Accept); err != nil {
		return nil, fail(c, "staff.respond", err)
	}

	return true, nil
}

// StaffLeaveHandler lets a member leave a business
func (h *Handlers) StaffLeaveHandler(c *core.WebContext) (any, *errs.Error) {
	var req LeaveBusinessRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Memberships.Remove(c, req.OwnerUid, c.GetActualUid()); err != nil {
		return nil, fail(c, "staff.leave", err)
	}

	return true, nil
}

// AuditListHandler lists what managers and staff did in the caller's business
func (h *Handlers) AuditListHandler(c *core.WebContext) (any, *errs.Error) {
	var req PageRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	entries, err := extservices.Audit.List(c, c.GetActualUid(), req.Limit, req.BeforeId)

	if err != nil {
		return nil, fail(c, "audit.list", err)
	}

	result := make([]*AuditView, 0, len(entries))

	for _, e := range entries {
		result = append(result, &AuditView{Id: e.AuditId, Action: e.Action, EntityType: e.EntityType, EntityId: e.EntityId, ActorUid: e.ActorUid, Role: e.Role.String(), Method: e.Method, Path: e.Path, Status: e.Status, Time: e.CreatedUnix})
	}

	return result, nil
}

var _ = exterrs.SubcategoryExt

// MySettingsHandler returns the caller's ext preferences
func (h *Handlers) MySettingsHandler(c *core.WebContext) (any, *errs.Error) {
	setting, configured, err := extservices.UserSettings.Get(c, c.GetActualUid())

	if err != nil {
		return nil, fail(c, "me.settings", err)
	}

	return &UserSettingsView{BusinessFeatures: setting.BusinessFeatures, Configured: configured}, nil
}

// MySettingsUpdateHandler changes the caller's ext preferences
func (h *Handlers) MySettingsUpdateHandler(c *core.WebContext) (any, *errs.Error) {
	var req UserSettingsRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	setting, err := extservices.UserSettings.SetBusinessFeatures(c, c.GetActualUid(), req.BusinessFeatures)

	if err != nil {
		return nil, fail(c, "me.settings.update", err)
	}

	return &UserSettingsView{BusinessFeatures: setting.BusinessFeatures, Configured: true}, nil
}
