package extapi

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
)

func (h *Handlers) businessProfileView(c *core.WebContext, profile *extmodels.BusinessProfile) *BusinessProfileView {
	name := profile.ReceiptName

	if name == "" {
		name = extservices.DisplayName(c, profile.OwnerUid)
	}

	return &BusinessProfileView{ReceiptName: profile.ReceiptName, Name: name, Address: profile.Address, Phone: profile.Phone, Footer: profile.Footer}
}

// BusinessProfileHandler returns what the business being worked in prints on its receipts
func (h *Handlers) BusinessProfileHandler(c *core.WebContext) (any, *errs.Error) {
	profile, err := extservices.BusinessProfiles.Get(c, c.GetCurrentUid())

	if err != nil {
		return nil, fail(c, "business.profile", err)
	}

	return h.businessProfileView(c, profile), nil
}

// MyBusinessProfileHandler returns the caller's own business profile, whichever business they are working in
func (h *Handlers) MyBusinessProfileHandler(c *core.WebContext) (any, *errs.Error) {
	profile, err := extservices.BusinessProfiles.Get(c, c.GetActualUid())

	if err != nil {
		return nil, fail(c, "me.business_profile", err)
	}

	return h.businessProfileView(c, profile), nil
}

// MyBusinessProfileUpdateHandler changes what the caller's own business prints on its receipts
func (h *Handlers) MyBusinessProfileUpdateHandler(c *core.WebContext) (any, *errs.Error) {
	var req BusinessProfileRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	profile, err := extservices.BusinessProfiles.Save(c, c.GetActualUid(), extservices.ProfileInput{
		ReceiptName: req.ReceiptName, Address: req.Address, Phone: req.Phone, Footer: req.Footer,
	})

	if err != nil {
		return nil, fail(c, "me.business_profile.update", err)
	}

	return h.businessProfileView(c, profile), nil
}

// RepaymentGetHandler returns one repayment with the sales it paid off
func (h *Handlers) RepaymentGetHandler(c *core.WebContext) (any, *errs.Error) {
	var req RepaymentGetRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	result, err := extservices.Repayments.Get(c, c.GetCurrentUid(), req.Id)

	if err != nil {
		return nil, fail(c, "repayments.get", err)
	}

	view := repaymentView(result.Repayment)

	for _, a := range result.Allocations {
		view.Allocations = append(view.Allocations, RepaymentAllocationView{SaleId: a.SaleId, Amount: a.Amount})
	}

	return view, nil
}
