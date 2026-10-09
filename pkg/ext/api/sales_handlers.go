package extapi

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmw "github.com/mayswind/ezbookkeeping/pkg/ext/middleware"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
)

func customerView(cu *extmodels.Customer, outstanding int64) *CustomerView {
	return &CustomerView{Id: cu.CustomerId, Name: cu.Name, Phone: cu.Phone, Email: cu.Email, Note: cu.Note, Outstanding: outstanding}
}

func saleView(s *extmodels.Sale, lines []*extmodels.SaleLine) *SaleView {
	view := &SaleView{
		Id: s.SaleId, LocationId: s.LocationId, CustomerId: s.CustomerId, Time: s.SaleTime, Subtotal: s.Subtotal, Discount: s.Discount,
		Total: s.Total, Paid: s.Paid, Outstanding: s.Outstanding(), Voided: s.Voided, PaymentAccountId: s.PaymentAccountId,
		ReceivableAccountId: s.ReceivableAccountId, CategoryId: s.CategoryId, PaidTransactionId: s.PaidTransactionId,
		CreditTransactionId: s.CreditTransactionId, Note: s.Note, ActorUid: s.ActorUid,
	}

	for _, l := range lines {
		view.Lines = append(view.Lines, SaleLineView{ItemId: l.ItemId, Qty: l.Qty, UnitPrice: l.UnitPrice, LineTotal: l.LineTotal})
	}

	return view
}

// CustomerListHandler lists customers with what each one owes
func (h *Handlers) CustomerListHandler(c *core.WebContext) (any, *errs.Error) {
	customers, err := extservices.Customers.List(c, c.GetCurrentUid())

	if err != nil {
		return nil, fail(c, "customers.list", err)
	}

	balances, err := extservices.Customers.Balances(c, c.GetCurrentUid())

	if err != nil {
		return nil, fail(c, "customers.list", err)
	}

	result := make([]*CustomerView, 0, len(customers))

	for _, cu := range customers {
		result = append(result, customerView(cu, balances[cu.CustomerId]))
	}

	return result, nil
}

// CustomerBalancesHandler lists only the customers who owe money
func (h *Handlers) CustomerBalancesHandler(c *core.WebContext) (any, *errs.Error) {
	customers, err := extservices.Customers.List(c, c.GetCurrentUid())

	if err != nil {
		return nil, fail(c, "customers.balances", err)
	}

	balances, err := extservices.Customers.Balances(c, c.GetCurrentUid())

	if err != nil {
		return nil, fail(c, "customers.balances", err)
	}

	result := make([]*CustomerView, 0)

	for _, cu := range customers {
		if balances[cu.CustomerId] > 0 {
			result = append(result, customerView(cu, balances[cu.CustomerId]))
		}
	}

	return result, nil
}

func customerInput(req *CustomerRequest) extservices.CustomerInput {
	return extservices.CustomerInput{Name: req.Name, Phone: req.Phone, Email: req.Email, Note: req.Note}
}

// CustomerAddHandler adds a customer
func (h *Handlers) CustomerAddHandler(c *core.WebContext) (any, *errs.Error) {
	var req CustomerRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	customer, err := extservices.Customers.Add(c, c.GetCurrentUid(), customerInput(&req))

	if err != nil {
		return nil, fail(c, "customers.add", err)
	}

	return customerView(customer, 0), nil
}

// CustomerModifyHandler updates a customer
func (h *Handlers) CustomerModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req CustomerRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	customer, err := extservices.Customers.Modify(c, c.GetCurrentUid(), req.Id, customerInput(&req))

	if err != nil {
		return nil, fail(c, "customers.modify", err)
	}

	owed, err := extservices.Customers.Outstanding(c, c.GetCurrentUid(), req.Id)

	if err != nil {
		return nil, fail(c, "customers.modify", err)
	}

	return customerView(customer, owed), nil
}

// CustomerDeleteHandler deletes a customer who owes nothing
func (h *Handlers) CustomerDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req IdRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Customers.Delete(c, c.GetCurrentUid(), req.Id); err != nil {
		return nil, fail(c, "customers.delete", err)
	}

	return true, nil
}

// SaleCreateHandler records a sale, books the money and moves the stock
func (h *Handlers) SaleCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var req SaleCreateRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	// staff sell at the listed price; changing prices or giving discounts is for managers and owners
	if extmw.RoleOf(c) < extmodels.RoleManager && (req.Discount != 0 || hasPriceOverride(req.Lines)) {
		return nil, exterrs.ErrNotPermittedForRole
	}

	lines := make([]extservices.SaleLineInput, 0, len(req.Lines))

	for _, l := range req.Lines {
		lines = append(lines, extservices.SaleLineInput{ItemId: l.ItemId, Qty: l.Qty, UnitPrice: l.UnitPrice})
	}

	detail, err := extservices.Sales.Create(c, c.GetCurrentUid(), c.GetActualUid(), extservices.SaleInput{
		LocationId: req.LocationId, CustomerId: req.CustomerId, Time: req.Time, UtcOffset: req.UtcOffset, Lines: lines,
		Discount: req.Discount, AmountPaid: req.AmountPaid, PaymentAccountId: req.PaymentAccountId,
		ReceivableAccountId: req.ReceivableAccountId, CategoryId: req.CategoryId, Note: req.Note,
	})

	if err != nil {
		return nil, fail(c, "sales.add", err)
	}

	return saleView(detail.Sale, detail.Lines), nil
}

// SaleListHandler lists sales
func (h *Handlers) SaleListHandler(c *core.WebContext) (any, *errs.Error) {
	var req SaleListRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	sales, err := extservices.Sales.List(c, c.GetCurrentUid(), extservices.SaleFilter{CustomerId: req.CustomerId, OnlyOpen: req.OnlyOpen, BeforeId: req.BeforeId, Limit: req.Limit})

	if err != nil {
		return nil, fail(c, "sales.list", err)
	}

	result := make([]*SaleView, 0, len(sales))

	for _, s := range sales {
		result = append(result, saleView(s, nil))
	}

	return result, nil
}

// SaleGetHandler returns one sale with its lines
func (h *Handlers) SaleGetHandler(c *core.WebContext) (any, *errs.Error) {
	var req SaleGetRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	detail, err := extservices.Sales.Get(c, c.GetCurrentUid(), req.Id)

	if err != nil {
		return nil, fail(c, "sales.get", err)
	}

	return saleView(detail.Sale, detail.Lines), nil
}

// SaleVoidHandler cancels a sale
func (h *Handlers) SaleVoidHandler(c *core.WebContext) (any, *errs.Error) {
	var req IdRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Sales.Void(c, c.GetCurrentUid(), c.GetActualUid(), req.Id); err != nil {
		return nil, fail(c, "sales.void", err)
	}

	return true, nil
}

// RepaymentCreateHandler records money received from a customer
func (h *Handlers) RepaymentCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var req RepaymentCreateRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	result, err := extservices.Repayments.Add(c, c.GetCurrentUid(), c.GetActualUid(), extservices.RepaymentInput{
		CustomerId: req.CustomerId, Amount: req.Amount, SaleId: req.SaleId, Time: req.Time, UtcOffset: req.UtcOffset,
		PaymentAccountId: req.PaymentAccountId, ReceivableAccountId: req.ReceivableAccountId, CategoryId: req.CategoryId, Note: req.Note,
	})

	if err != nil {
		return nil, fail(c, "repayments.add", err)
	}

	view := repaymentView(result.Repayment)

	for _, a := range result.Allocations {
		view.Allocations = append(view.Allocations, RepaymentAllocationView{SaleId: a.SaleId, Amount: a.Amount})
	}

	return view, nil
}

func repaymentView(r *extmodels.Repayment) *RepaymentView {
	return &RepaymentView{Id: r.RepaymentId, CustomerId: r.CustomerId, Amount: r.Amount, TransactionId: r.TransactionId, PaymentAccountId: r.PaymentAccountId, ReceivableAccountId: r.ReceivableAccountId, Time: r.RepaymentTime, Note: r.Note, ActorUid: r.ActorUid}
}

// RepaymentListHandler lists repayments
func (h *Handlers) RepaymentListHandler(c *core.WebContext) (any, *errs.Error) {
	var req RepaymentListRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	repayments, err := extservices.Repayments.List(c, c.GetCurrentUid(), req.CustomerId, req.BeforeId, req.Limit)

	if err != nil {
		return nil, fail(c, "repayments.list", err)
	}

	result := make([]*RepaymentView, 0, len(repayments))

	for _, r := range repayments {
		result = append(result, repaymentView(r))
	}

	return result, nil
}

func hasPriceOverride(lines []SaleLineRequest) bool {
	for _, l := range lines {
		if l.UnitPrice != nil {
			return true
		}
	}

	return false
}
