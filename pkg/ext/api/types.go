// Package extapi contains the HTTP handlers of the ext module.
// Ids are serialized as strings and amounts as integers in minor currency units, like upstream.
// Quantities are integers scaled by 1000 (2.5 units is 2500).
package extapi

import (
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func roleFromString(role string) extmodels.Role {
	switch role {
	case "manager":
		return extmodels.RoleManager
	case "staff":
		return extmodels.RoleStaff
	case "owner":
		return extmodels.RoleOwner
	default:
		return 0
	}
}

// ---- requests

type IdRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

type InviteStaffRequest struct {
	Email string `json:"email" binding:"required,max=128"`
	Role  string `json:"role" binding:"required"`
}

type StaffRoleRequest struct {
	StaffUid int64  `json:"staffUid,string" binding:"required,min=1"`
	Role     string `json:"role" binding:"required"`
}

type StaffRemoveRequest struct {
	StaffUid int64 `json:"staffUid,string" binding:"required,min=1"`
}

type RespondInviteRequest struct {
	OwnerUid int64 `json:"ownerUid,string" binding:"required,min=1"`
	Accept   bool  `json:"accept"`
}

type UserSettingsRequest struct {
	BusinessFeatures bool `json:"businessFeatures"`
}

type BusinessProfileRequest struct {
	ReceiptName string `json:"receiptName" binding:"max=128"`
	Address     string `json:"address" binding:"max=255"`
	Phone       string `json:"phone" binding:"max=32"`
	Footer      string `json:"footer" binding:"max=255"`
}

type RepaymentGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

type LeaveBusinessRequest struct {
	OwnerUid int64 `json:"ownerUid,string" binding:"required,min=1"`
}

type PageRequest struct {
	BeforeId int64 `form:"beforeId,string"`
	Limit    int   `form:"limit"`
}

type LocationRequest struct {
	Id   int64  `json:"id,string"`
	Name string `json:"name" binding:"required,max=64"`
}

type ItemRequest struct {
	Id           int64  `json:"id,string"`
	Sku          string `json:"sku" binding:"required,max=64"`
	Name         string `json:"name" binding:"required,max=128"`
	Unit         string `json:"unit" binding:"max=16"`
	CostPrice    int64  `json:"costPrice" binding:"min=0"`
	SalePrice    int64  `json:"salePrice" binding:"min=0"`
	ReorderLevel int64  `json:"reorderLevel" binding:"min=0"`
	TrackStock   bool   `json:"trackStock"`
}

type StockLevelRequest struct {
	ItemId     int64 `form:"itemId,string"`
	LocationId int64 `form:"locationId,string"`
}

type StockChangeRequest struct {
	ItemId     int64  `json:"itemId,string" binding:"required,min=1"`
	LocationId int64  `json:"locationId,string"`
	Qty        int64  `json:"qty" binding:"required"`
	UnitCost   int64  `json:"unitCost" binding:"min=0"`
	Note       string `json:"note" binding:"max=255"`
	Time       int64  `json:"time"`
	Opening    bool   `json:"opening"`
}

type StockTransferRequest struct {
	ItemId         int64  `json:"itemId,string" binding:"required,min=1"`
	FromLocationId int64  `json:"fromLocationId,string"`
	ToLocationId   int64  `json:"toLocationId,string" binding:"required,min=1"`
	Qty            int64  `json:"qty" binding:"required,min=1"`
	Note           string `json:"note" binding:"max=255"`
}

type MovementListRequest struct {
	PageRequest
	ItemId     int64 `form:"itemId,string"`
	LocationId int64 `form:"locationId,string"`
}

type CustomerRequest struct {
	Id    int64  `json:"id,string"`
	Name  string `json:"name" binding:"required,max=128"`
	Phone string `json:"phone" binding:"max=32"`
	Email string `json:"email" binding:"max=128"`
	Note  string `json:"note" binding:"max=255"`
}

type SaleLineRequest struct {
	ItemId    int64  `json:"itemId,string" binding:"required,min=1"`
	Qty       int64  `json:"qty" binding:"required,min=1"`
	UnitPrice *int64 `json:"unitPrice" binding:"omitempty,min=0"`
}

type SaleCreateRequest struct {
	LocationId          int64             `json:"locationId,string"`
	CustomerId          int64             `json:"customerId,string"`
	Time                int64             `json:"time"`
	UtcOffset           int16             `json:"utcOffset" binding:"min=-720,max=840"`
	Lines               []SaleLineRequest `json:"lines" binding:"required,min=1,max=200,dive"`
	Discount            int64             `json:"discount" binding:"min=0"`
	AmountPaid          int64             `json:"amountPaid" binding:"min=0"`
	PaymentAccountId    int64             `json:"paymentAccountId,string"`
	ReceivableAccountId int64             `json:"receivableAccountId,string"`
	CategoryId          int64             `json:"categoryId,string" binding:"required,min=1"`
	Note                string            `json:"note" binding:"max=255"`
}

type SaleListRequest struct {
	PageRequest
	CustomerId int64 `form:"customerId,string"`
	OnlyOpen   bool  `form:"onlyOpen"`
}

type SaleGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

type RepaymentCreateRequest struct {
	CustomerId          int64  `json:"customerId,string" binding:"required,min=1"`
	Amount              int64  `json:"amount" binding:"required,min=1"`
	SaleId              int64  `json:"saleId,string"`
	Time                int64  `json:"time"`
	UtcOffset           int16  `json:"utcOffset" binding:"min=-720,max=840"`
	PaymentAccountId    int64  `json:"paymentAccountId,string" binding:"required,min=1"`
	ReceivableAccountId int64  `json:"receivableAccountId,string" binding:"required,min=1"`
	CategoryId          int64  `json:"categoryId,string" binding:"required,min=1"`
	Note                string `json:"note" binding:"max=255"`
}

type RepaymentListRequest struct {
	PageRequest
	CustomerId int64 `form:"customerId,string"`
}

// ---- responses

type UserSettingsView struct {
	BusinessFeatures bool `json:"businessFeatures"`
	Configured       bool `json:"configured"` // false until the person chose, so the app can carry over an older local choice
}

type BusinessView struct {
	OwnerUid int64  `json:"ownerUid,string"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Status   string `json:"status"` // "owner", "active" or "pending"
}

type StaffView struct {
	StaffUid int64  `json:"staffUid,string"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type AuditView struct {
	Id         int64  `json:"id,string"`
	Action     string `json:"action"`
	EntityType string `json:"entityType"`
	EntityId   int64  `json:"entityId,string"`
	ActorUid   int64  `json:"actorUid,string"`
	Role       string `json:"role"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	Time       int64  `json:"time"`
}

type LocationView struct {
	Id        int64  `json:"id,string"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
}

type ItemView struct {
	Id           int64  `json:"id,string"`
	Sku          string `json:"sku"`
	Name         string `json:"name"`
	Unit         string `json:"unit"`
	CostPrice    int64  `json:"costPrice"`
	SalePrice    int64  `json:"salePrice"`
	ReorderLevel int64  `json:"reorderLevel"`
	TrackStock   bool   `json:"trackStock"`
}

type StockLevelView struct {
	ItemId     int64 `json:"itemId,string"`
	LocationId int64 `json:"locationId,string"`
	Qty        int64 `json:"qty"`
}

type MovementView struct {
	Id         int64  `json:"id,string"`
	ItemId     int64  `json:"itemId,string"`
	LocationId int64  `json:"locationId,string"`
	QtyChange  int64  `json:"qtyChange"`
	Reason     int    `json:"reason"`
	RefType    string `json:"refType"`
	RefId      int64  `json:"refId,string"`
	UnitCost   int64  `json:"unitCost"`
	Note       string `json:"note"`
	ActorUid   int64  `json:"actorUid,string"`
	Time       int64  `json:"time"`
}

type CustomerView struct {
	Id          int64  `json:"id,string"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Note        string `json:"note"`
	Outstanding int64  `json:"outstanding"`
}

type SaleLineView struct {
	ItemId    int64 `json:"itemId,string"`
	Qty       int64 `json:"qty"`
	UnitPrice int64 `json:"unitPrice"`
	LineTotal int64 `json:"lineTotal"`
}

type SaleView struct {
	Id                  int64          `json:"id,string"`
	LocationId          int64          `json:"locationId,string"`
	CustomerId          int64          `json:"customerId,string"`
	Time                int64          `json:"time"`
	Subtotal            int64          `json:"subtotal"`
	Discount            int64          `json:"discount"`
	Total               int64          `json:"total"`
	Paid                int64          `json:"paid"`
	Outstanding         int64          `json:"outstanding"`
	Voided              bool           `json:"voided"`
	PaymentAccountId    int64          `json:"paymentAccountId,string"`
	ReceivableAccountId int64          `json:"receivableAccountId,string"`
	CategoryId          int64          `json:"categoryId,string"`
	PaidTransactionId   int64          `json:"paidTransactionId,string"`
	CreditTransactionId int64          `json:"creditTransactionId,string"`
	Note                string         `json:"note"`
	ActorUid            int64          `json:"actorUid,string"`
	Lines               []SaleLineView `json:"lines,omitempty"`
}

type RepaymentAllocationView struct {
	SaleId int64 `json:"saleId,string"`
	Amount int64 `json:"amount"`
}

type RepaymentView struct {
	Id                  int64                     `json:"id,string"`
	CustomerId          int64                     `json:"customerId,string"`
	Amount              int64                     `json:"amount"`
	TransactionId       int64                     `json:"transactionId,string"`
	PaymentAccountId    int64                     `json:"paymentAccountId,string"`
	ReceivableAccountId int64                     `json:"receivableAccountId,string"`
	Time                int64                     `json:"time"`
	Note                string                    `json:"note"`
	ActorUid            int64                     `json:"actorUid,string"`
	Allocations         []RepaymentAllocationView `json:"allocations,omitempty"`
}

type ReportLocationRequest struct {
	LocationId int64 `form:"locationId,string"`
}

type PersonView struct {
	Uid    int64  `json:"uid,string"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"` // false for people who were removed or left
}

type LocationQtyView struct {
	LocationId int64 `json:"locationId,string"`
	Qty        int64 `json:"qty"`
}

type StockValueRowView struct {
	Item        *ItemView         `json:"item"`
	Qty         int64             `json:"qty"`
	CostValue   int64             `json:"costValue"`
	RetailValue int64             `json:"retailValue"`
	Locations   []LocationQtyView `json:"locations,omitempty"`
}

type StockValueReportView struct {
	Rows             []StockValueRowView `json:"rows"`
	TotalCostValue   int64               `json:"totalCostValue"`
	TotalRetailValue int64               `json:"totalRetailValue"`
}

type LowStockRowView struct {
	Item      *ItemView `json:"item"`
	Qty       int64     `json:"qty"`
	Shortfall int64     `json:"shortfall"`
}

type ReceivableRowView struct {
	Customer       *CustomerView `json:"customer"`
	Outstanding    int64         `json:"outstanding"`
	Current        int64         `json:"current"`
	Days31To60     int64         `json:"days31To60"`
	Days61To90     int64         `json:"days61To90"`
	Over90         int64         `json:"over90"`
	OldestSaleTime int64         `json:"oldestSaleTime"`
	OpenSales      int           `json:"openSales"`
}

type ReceivablesReportView struct {
	Rows             []ReceivableRowView `json:"rows"`
	TotalOutstanding int64               `json:"totalOutstanding"`
	Current          int64               `json:"current"`
	Days31To60       int64               `json:"days31To60"`
	Days61To90       int64               `json:"days61To90"`
	Over90           int64               `json:"over90"`
}

// BusinessProfileView is what a business prints on its receipts. Name is the receipt name, or the owner's name when none is set.
type BusinessProfileView struct {
	ReceiptName string `json:"receiptName"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Footer      string `json:"footer"`
}
