// Package exterrs defines error codes of the ext (owner/staff, inventory and sales) module.
// Package name differs from the directory name on purpose to avoid clashing with pkg/errs.
package exterrs

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// SubcategoryExt is the error sub category reserved for the ext module (upstream uses 0-20 so far)
const SubcategoryExt = 100

// Error codes related to staff and business membership
var (
	ErrBusinessIdInvalid       = errs.NewNormalError(SubcategoryExt, 0, http.StatusBadRequest, "business id is invalid")
	ErrNotAMemberOfBusiness    = errs.NewNormalError(SubcategoryExt, 1, http.StatusForbidden, "you are not an active member of this business")
	ErrNotPermittedInBusiness  = errs.NewNormalError(SubcategoryExt, 2, http.StatusForbidden, "your role is not permitted to perform this action in this business")
	ErrMembershipNotFound      = errs.NewNormalError(SubcategoryExt, 3, http.StatusNotFound, "membership not found")
	ErrCannotInviteSelf        = errs.NewNormalError(SubcategoryExt, 4, http.StatusBadRequest, "cannot invite yourself")
	ErrInviteeNotFound         = errs.NewNormalError(SubcategoryExt, 5, http.StatusBadRequest, "no user is registered with this email, the invitee must register first")
	ErrRoleInvalid             = errs.NewNormalError(SubcategoryExt, 6, http.StatusBadRequest, "role is invalid, must be manager or staff")
	ErrMembershipAlreadyActive = errs.NewNormalError(SubcategoryExt, 7, http.StatusBadRequest, "this user is already a member of your business")
	ErrMembershipNotPending    = errs.NewNormalError(SubcategoryExt, 8, http.StatusBadRequest, "membership is not pending")
	ErrBusinessIsNotAccessible = errs.NewNormalError(SubcategoryExt, 9, http.StatusForbidden, "this endpoint cannot be used in a business context")
)

// Error codes related to inventory
var (
	ErrLocationNotFound         = errs.NewNormalError(SubcategoryExt, 20, http.StatusNotFound, "location not found")
	ErrLocationNameIsEmpty      = errs.NewNormalError(SubcategoryExt, 21, http.StatusBadRequest, "location name is empty")
	ErrLocationNameExists       = errs.NewNormalError(SubcategoryExt, 22, http.StatusBadRequest, "location name already exists")
	ErrCannotDeleteLastLocation = errs.NewNormalError(SubcategoryExt, 23, http.StatusBadRequest, "cannot delete the last location")
	ErrLocationHasStock         = errs.NewNormalError(SubcategoryExt, 24, http.StatusBadRequest, "location still has stock, transfer it first")
	ErrItemNotFound             = errs.NewNormalError(SubcategoryExt, 25, http.StatusNotFound, "item not found")
	ErrItemNameIsEmpty          = errs.NewNormalError(SubcategoryExt, 26, http.StatusBadRequest, "item name is empty")
	ErrItemSkuIsEmpty           = errs.NewNormalError(SubcategoryExt, 27, http.StatusBadRequest, "item sku is empty")
	ErrItemSkuExists            = errs.NewNormalError(SubcategoryExt, 28, http.StatusBadRequest, "item sku already exists")
	ErrItemPriceInvalid         = errs.NewNormalError(SubcategoryExt, 29, http.StatusBadRequest, "item price cannot be negative")
	ErrQuantityInvalid          = errs.NewNormalError(SubcategoryExt, 30, http.StatusBadRequest, "quantity is invalid")
	ErrInsufficientStock        = errs.NewNormalError(SubcategoryExt, 31, http.StatusBadRequest, "insufficient stock")
	ErrSameLocation             = errs.NewNormalError(SubcategoryExt, 32, http.StatusBadRequest, "source and destination location are the same")
	ErrItemDoesNotTrackStock    = errs.NewNormalError(SubcategoryExt, 33, http.StatusBadRequest, "item does not track stock")
	ErrStockReasonInvalid       = errs.NewNormalError(SubcategoryExt, 34, http.StatusBadRequest, "stock movement reason is invalid")
)

// Error codes related to customers, sales and repayments
var (
	ErrCustomerNotFound           = errs.NewNormalError(SubcategoryExt, 40, http.StatusNotFound, "customer not found")
	ErrCustomerNameIsEmpty        = errs.NewNormalError(SubcategoryExt, 41, http.StatusBadRequest, "customer name is empty")
	ErrCustomerHasBalance         = errs.NewNormalError(SubcategoryExt, 42, http.StatusBadRequest, "customer still owes money, cannot be deleted")
	ErrSaleNotFound               = errs.NewNormalError(SubcategoryExt, 43, http.StatusNotFound, "sale not found")
	ErrSaleHasNoLines             = errs.NewNormalError(SubcategoryExt, 44, http.StatusBadRequest, "sale has no lines")
	ErrSaleDiscountInvalid        = errs.NewNormalError(SubcategoryExt, 45, http.StatusBadRequest, "sale discount is invalid")
	ErrSaleTotalInvalid           = errs.NewNormalError(SubcategoryExt, 46, http.StatusBadRequest, "sale total must be greater than zero")
	ErrSaleAmountPaidInvalid      = errs.NewNormalError(SubcategoryExt, 47, http.StatusBadRequest, "amount paid must be between zero and the sale total")
	ErrCreditSaleRequiresCustomer = errs.NewNormalError(SubcategoryExt, 48, http.StatusBadRequest, "a sale on credit requires a customer")
	ErrAccountsRequired           = errs.NewNormalError(SubcategoryExt, 49, http.StatusBadRequest, "payment account and receivable account are required for this sale")
	ErrAccountNotUsable           = errs.NewNormalError(SubcategoryExt, 50, http.StatusBadRequest, "account does not exist, is hidden or is a parent account")
	ErrAccountCurrencyMismatch    = errs.NewNormalError(SubcategoryExt, 51, http.StatusBadRequest, "payment account and receivable account must use the same currency")
	ErrReceivableAccountInvalid   = errs.NewNormalError(SubcategoryExt, 52, http.StatusBadRequest, "receivable account must be a receivables account")
	ErrSaleAlreadyVoided          = errs.NewNormalError(SubcategoryExt, 53, http.StatusBadRequest, "sale is already voided")
	ErrSaleHasRepayments          = errs.NewNormalError(SubcategoryExt, 54, http.StatusBadRequest, "sale has repayments, they must be handled before voiding")
	ErrRepaymentAmountInvalid     = errs.NewNormalError(SubcategoryExt, 55, http.StatusBadRequest, "repayment amount must be greater than zero")
	ErrRepaymentExceedsBalance    = errs.NewNormalError(SubcategoryExt, 56, http.StatusBadRequest, "repayment amount exceeds the outstanding balance")
	ErrSaleCategoryInvalid        = errs.NewNormalError(SubcategoryExt, 57, http.StatusBadRequest, "sale category must be an income category")
	ErrConcurrentModification     = errs.NewNormalError(SubcategoryExt, 58, http.StatusConflict, "data was modified by another request, please retry")
	ErrFinancialRecordFailed      = errs.NewNormalError(SubcategoryExt, 59, http.StatusInternalServerError, "failed to record the financial transaction")
)

// Error codes added after review
var (
	ErrItemFieldTooLong    = errs.NewNormalError(SubcategoryExt, 35, http.StatusBadRequest, "item sku (max 64), name (max 128) or unit (max 16) is too long")
	ErrStockTooLarge       = errs.NewNormalError(SubcategoryExt, 36, http.StatusBadRequest, "stock quantity or unit cost is too large")
	ErrNotPermittedForRole = errs.NewNormalError(SubcategoryExt, 60, http.StatusForbidden, "your role cannot override prices or give discounts")
)

// Error codes added with receipts
var (
	ErrProfileFieldTooLong = errs.NewNormalError(SubcategoryExt, 61, http.StatusBadRequest, "business name (max 128), address (max 255), phone (max 32) or footer (max 255) is too long")
	ErrRepaymentNotFound   = errs.NewNormalError(SubcategoryExt, 62, http.StatusNotFound, "repayment not found")
)
