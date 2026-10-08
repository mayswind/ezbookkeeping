// Package extmodels contains the database models of the ext module.
// All tables are prefixed with "ext_" so they can never collide with upstream tables.
package extmodels

// QtyScale is the fixed-point scale of quantities: 1 unit is stored as 1000 (3 decimal places)
const QtyScale int64 = 1000

// Role represents the role of a user inside a business. The business owner is the user who owns the data.
type Role byte

// Roles, a higher value includes the permissions of all lower values
const (
	RoleStaff   Role = 1
	RoleManager Role = 2
	RoleOwner   Role = 3
)

// String returns the textual role name
func (r Role) String() string {
	switch r {
	case RoleStaff:
		return "staff"
	case RoleManager:
		return "manager"
	case RoleOwner:
		return "owner"
	default:
		return "unknown"
	}
}

// MembershipStatus represents the status of a membership
type MembershipStatus byte

// Membership statuses
const (
	MembershipStatusPending MembershipStatus = 1
	MembershipStatusActive  MembershipStatus = 2
	MembershipStatusRevoked MembershipStatus = 3
)

// StockReason represents why a stock movement happened
type StockReason byte

// Stock movement reasons
const (
	StockReasonOpening     StockReason = 1
	StockReasonPurchase    StockReason = 2
	StockReasonSale        StockReason = 3
	StockReasonAdjustment  StockReason = 4
	StockReasonTransferOut StockReason = 5
	StockReasonTransferIn  StockReason = 6
	StockReasonSaleVoid    StockReason = 7
)

// Membership links a staff user to a business owner
type Membership struct {
	MembershipId    int64            `xorm:"PK AUTOINCR"`
	OwnerUid        int64            `xorm:"UNIQUE(UQE_ext_membership_owner_staff) INDEX(IDX_ext_membership_owner) NOT NULL"`
	StaffUid        int64            `xorm:"UNIQUE(UQE_ext_membership_owner_staff) INDEX(IDX_ext_membership_staff) NOT NULL"`
	Role            Role             `xorm:"NOT NULL"`
	Status          MembershipStatus `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}

// TableName returns the table name
func (Membership) TableName() string { return "ext_membership" }

// AuditLog records every write request made by a manager or staff member on behalf of an owner
type AuditLog struct {
	AuditId     int64  `xorm:"PK AUTOINCR"`
	OwnerUid    int64  `xorm:"INDEX(IDX_ext_audit_log_owner_time) NOT NULL"`
	ActorUid    int64  `xorm:"NOT NULL"`
	Role        Role   `xorm:"NOT NULL"`
	Method      string `xorm:"VARCHAR(8) NOT NULL"`
	Path        string `xorm:"VARCHAR(255) NOT NULL"`
	Status      int    `xorm:"NOT NULL"`
	ClientIp    string `xorm:"VARCHAR(45)"`
	CreatedUnix int64  `xorm:"INDEX(IDX_ext_audit_log_owner_time) NOT NULL"`
}

// TableName returns the table name
func (AuditLog) TableName() string { return "ext_audit_log" }

// Location is a store, warehouse or any place where stock is kept
type Location struct {
	LocationId      int64  `xorm:"PK AUTOINCR"`
	OwnerUid        int64  `xorm:"UNIQUE(UQE_ext_location_owner_name) INDEX(IDX_ext_location_owner_deleted) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_ext_location_owner_deleted) NOT NULL"`
	Name            string `xorm:"VARCHAR(64) UNIQUE(UQE_ext_location_owner_name) NOT NULL"`
	IsDefault       bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64 `xorm:"UNIQUE(UQE_ext_location_owner_name) NOT NULL DEFAULT 0"` // 0 while active, so two active locations never share a name
}

// TableName returns the table name
func (Location) TableName() string { return "ext_location" }

// Item is a product (or service) that can be sold
type Item struct {
	ItemId          int64  `xorm:"PK AUTOINCR"`
	OwnerUid        int64  `xorm:"UNIQUE(UQE_ext_item_owner_sku) INDEX(IDX_ext_item_owner_deleted) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_ext_item_owner_deleted) NOT NULL"`
	Sku             string `xorm:"VARCHAR(64) UNIQUE(UQE_ext_item_owner_sku) NOT NULL"`
	Name            string `xorm:"VARCHAR(128) NOT NULL"`
	Unit            string `xorm:"VARCHAR(16) NOT NULL"`
	CostPrice       int64  `xorm:"NOT NULL"` // minor currency units per whole unit
	SalePrice       int64  `xorm:"NOT NULL"` // minor currency units per whole unit
	ReorderLevel    int64  `xorm:"NOT NULL"` // quantity, scaled by QtyScale
	TrackStock      bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// TableName returns the table name
func (Item) TableName() string { return "ext_item" }

// StockMovement is an immutable entry of the stock ledger. Stock on hand is the sum of movements.
type StockMovement struct {
	MovementId   int64       `xorm:"PK AUTOINCR"`
	OwnerUid     int64       `xorm:"INDEX(IDX_ext_stock_movement_owner_item_location) INDEX(IDX_ext_stock_movement_owner_time) NOT NULL"`
	ItemId       int64       `xorm:"INDEX(IDX_ext_stock_movement_owner_item_location) NOT NULL"`
	LocationId   int64       `xorm:"INDEX(IDX_ext_stock_movement_owner_item_location) NOT NULL"`
	QtyChange    int64       `xorm:"NOT NULL"` // scaled by QtyScale, negative when stock leaves
	Reason       StockReason `xorm:"NOT NULL"`
	RefType      string      `xorm:"VARCHAR(16) NOT NULL"` // "sale", "transfer" or ""
	RefId        int64       `xorm:"NOT NULL"`
	UnitCost     int64       `xorm:"NOT NULL"`
	Note         string      `xorm:"VARCHAR(255) NOT NULL"`
	ActorUid     int64       `xorm:"NOT NULL"`
	MovementTime int64       `xorm:"INDEX(IDX_ext_stock_movement_owner_time) NOT NULL"`
	CreatedUnix  int64
}

// TableName returns the table name
func (StockMovement) TableName() string { return "ext_stock_movement" }

// Customer is somebody who buys, possibly on credit
type Customer struct {
	CustomerId      int64  `xorm:"PK AUTOINCR"`
	OwnerUid        int64  `xorm:"INDEX(IDX_ext_customer_owner_deleted) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_ext_customer_owner_deleted) NOT NULL"`
	Name            string `xorm:"VARCHAR(128) NOT NULL"`
	Phone           string `xorm:"VARCHAR(32) NOT NULL"`
	Email           string `xorm:"VARCHAR(128) NOT NULL"`
	Note            string `xorm:"VARCHAR(255) NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// TableName returns the table name
func (Customer) TableName() string { return "ext_customer" }

// Sale is a sale of one or more items, fully or partly paid at the time of sale
type Sale struct {
	SaleId              int64  `xorm:"PK AUTOINCR"`
	OwnerUid            int64  `xorm:"INDEX(IDX_ext_sale_owner_time) INDEX(IDX_ext_sale_owner_customer) NOT NULL"`
	LocationId          int64  `xorm:"NOT NULL"`
	CustomerId          int64  `xorm:"INDEX(IDX_ext_sale_owner_customer) NOT NULL"` // 0 for walk-in customers
	SaleTime            int64  `xorm:"INDEX(IDX_ext_sale_owner_time) NOT NULL"`
	Subtotal            int64  `xorm:"NOT NULL"`
	Discount            int64  `xorm:"NOT NULL"`
	Total               int64  `xorm:"NOT NULL"`
	Paid                int64  `xorm:"NOT NULL"` // paid at sale time plus later repayments
	Voided              bool   `xorm:"NOT NULL"`
	PaymentAccountId    int64  `xorm:"NOT NULL"`
	ReceivableAccountId int64  `xorm:"NOT NULL"`
	CategoryId          int64  `xorm:"NOT NULL"`
	PaidTransactionId   int64  `xorm:"NOT NULL"` // transaction for the amount paid at sale time, 0 if none
	CreditTransactionId int64  `xorm:"NOT NULL"` // transaction for the amount put on credit, 0 if none
	Note                string `xorm:"VARCHAR(255) NOT NULL"`
	ActorUid            int64  `xorm:"NOT NULL"`
	CreatedUnixTime     int64
	VoidedUnixTime      int64
}

// TableName returns the table name
func (Sale) TableName() string { return "ext_sale" }

// Outstanding returns the amount still owed on the sale
func (s *Sale) Outstanding() int64 {
	if s.Voided {
		return 0
	}

	return s.Total - s.Paid
}

// SaleLine is one line of a sale
type SaleLine struct {
	LineId    int64 `xorm:"PK AUTOINCR"`
	OwnerUid  int64 `xorm:"NOT NULL"`
	SaleId    int64 `xorm:"INDEX(IDX_ext_sale_line_sale) NOT NULL"`
	ItemId    int64 `xorm:"NOT NULL"`
	Qty       int64 `xorm:"NOT NULL"` // scaled by QtyScale
	UnitPrice int64 `xorm:"NOT NULL"` // minor currency units per whole unit
	LineTotal int64 `xorm:"NOT NULL"`
}

// TableName returns the table name
func (SaleLine) TableName() string { return "ext_sale_line" }

// Repayment is money received from a customer against what they owe
type Repayment struct {
	RepaymentId         int64  `xorm:"PK AUTOINCR"`
	OwnerUid            int64  `xorm:"INDEX(IDX_ext_repayment_owner_customer) NOT NULL"`
	CustomerId          int64  `xorm:"INDEX(IDX_ext_repayment_owner_customer) NOT NULL"`
	Amount              int64  `xorm:"NOT NULL"`
	PaymentAccountId    int64  `xorm:"NOT NULL"`
	ReceivableAccountId int64  `xorm:"NOT NULL"`
	TransactionId       int64  `xorm:"NOT NULL"`
	RepaymentTime       int64  `xorm:"NOT NULL"`
	Note                string `xorm:"VARCHAR(255) NOT NULL"`
	ActorUid            int64  `xorm:"NOT NULL"`
	CreatedUnixTime     int64
}

// TableName returns the table name
func (Repayment) TableName() string { return "ext_repayment" }

// RepaymentAllocation records which sale a repayment (or part of it) paid off
type RepaymentAllocation struct {
	AllocationId int64 `xorm:"PK AUTOINCR"`
	OwnerUid     int64 `xorm:"NOT NULL"`
	RepaymentId  int64 `xorm:"INDEX(IDX_ext_repayment_allocation_repayment) NOT NULL"`
	SaleId       int64 `xorm:"INDEX(IDX_ext_repayment_allocation_sale) NOT NULL"`
	Amount       int64 `xorm:"NOT NULL"`
}

// TableName returns the table name
func (RepaymentAllocation) TableName() string { return "ext_repayment_allocation" }

// GlobalTables lists the tables that live in the user database because they are looked up across owners
func GlobalTables() []any {
	return []any{
		new(Membership),
		new(AuditLog),
	}
}

// OwnerTables lists the tables that live next to the owner's business data
func OwnerTables() []any {
	return []any{
		new(Location),
		new(Item),
		new(StockMovement),
		new(Customer),
		new(Sale),
		new(SaleLine),
		new(Repayment),
		new(RepaymentAllocation),
	}
}

// StockLevel is the stock on hand of one item at one location, derived from the stock ledger
type StockLevel struct {
	ItemId     int64 `json:"-"`
	LocationId int64 `json:"-"`
	Qty        int64 `json:"-"`
}

// SaleTotals is the sum of sales of one customer, used to compute what the customer owes
type SaleTotals struct {
	CustomerId int64 `json:"-"`
	Total      int64 `json:"-"`
	Paid       int64 `json:"-"`
}
