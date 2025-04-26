package models

import "encoding/json"

// LevelOneAccountParentId represents the parent id of level-one account
const LevelOneAccountParentId = 0

// AccountCategory represents account category
type AccountCategory byte

// Account categories
const (
	ACCOUNT_CATEGORY_CASH                   AccountCategory = 1
	ACCOUNT_CATEGORY_CHECKING_ACCOUNT       AccountCategory = 2
	ACCOUNT_CATEGORY_CREDIT_CARD            AccountCategory = 3
	ACCOUNT_CATEGORY_VIRTUAL                AccountCategory = 4
	ACCOUNT_CATEGORY_DEBT                   AccountCategory = 5
	ACCOUNT_CATEGORY_RECEIVABLES            AccountCategory = 6
	ACCOUNT_CATEGORY_INVESTMENT             AccountCategory = 7
	ACCOUNT_CATEGORY_SAVINGS_ACCOUNT        AccountCategory = 8
	ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT AccountCategory = 9
)

var assetAccountCategory = map[AccountCategory]bool{
	ACCOUNT_CATEGORY_CASH:                   true,
	ACCOUNT_CATEGORY_CHECKING_ACCOUNT:       true,
	ACCOUNT_CATEGORY_CREDIT_CARD:            false,
	ACCOUNT_CATEGORY_VIRTUAL:                true,
	ACCOUNT_CATEGORY_DEBT:                   false,
	ACCOUNT_CATEGORY_RECEIVABLES:            true,
	ACCOUNT_CATEGORY_INVESTMENT:             true,
	ACCOUNT_CATEGORY_SAVINGS_ACCOUNT:        true,
	ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT: true,
}

var liabilityAccountCategory = map[AccountCategory]bool{
	ACCOUNT_CATEGORY_CASH:                   false,
	ACCOUNT_CATEGORY_CHECKING_ACCOUNT:       false,
	ACCOUNT_CATEGORY_CREDIT_CARD:            true,
	ACCOUNT_CATEGORY_VIRTUAL:                false,
	ACCOUNT_CATEGORY_DEBT:                   true,
	ACCOUNT_CATEGORY_RECEIVABLES:            false,
	ACCOUNT_CATEGORY_INVESTMENT:             false,
	ACCOUNT_CATEGORY_SAVINGS_ACCOUNT:        false,
	ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT: false,
}

// AccountType represents account type
type AccountType byte

// Account types
const (
	ACCOUNT_TYPE_SINGLE_ACCOUNT     AccountType = 1
	ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS AccountType = 2
)

var defaultCreditCardAccountStatementDate = 0

// Account represents account data stored in database
type Account struct {
	AccountId       int64           `xorm:"PK"`
	Uid             int64           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Deleted         bool            `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Category        AccountCategory `xorm:"NOT NULL"`
	Type            AccountType     `xorm:"NOT NULL"`
	ParentAccountId int64           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Name            string          `xorm:"VARCHAR(64) NOT NULL"`
	DisplayOrder    int32           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Icon            int64           `xorm:"NOT NULL"`
	Color           string          `xorm:"VARCHAR(6) NOT NULL"`
	Currency        string          `xorm:"VARCHAR(3) NOT NULL"`
	Balance         int64           `xorm:"NOT NULL"`
	Comment         string          `xorm:"VARCHAR(255) NOT NULL"`
	Extend          *AccountExtend  `xorm:"BLOB"`
	Hidden          bool            `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// AccountExtend represents account extend data stored in database
type AccountExtend struct {
	CreditCardStatementDate *int `json:"creditCardStatementDate"`
}

// AccountCreateRequest represents all parameters of account creation request
type AccountCreateRequest struct {
	Name                    string                  `json:"name" binding:"required,notBlank,max=64"`
	Category                AccountCategory         `json:"category" binding:"required"`
	Type                    AccountType             `json:"type" binding:"required"`
	Icon                    int64                   `json:"icon,string" binding:"required,min=1"`
	Color                   string                  `json:"color" binding:"required,len=6,validHexRGBColor"`
	Currency                string                  `json:"currency" binding:"required,len=3,validCurrency"`
	Balance                 int64                   `json:"balance"`
	BalanceTime             int64                   `json:"balanceTime"`
	Comment                 string                  `json:"comment" binding:"max=255"`
	CreditCardStatementDate int                     `json:"creditCardStatementDate" binding:"min=0,max=28"`
	SubAccounts             []*AccountCreateRequest `json:"subAccounts" binding:"omitempty"`
	ClientSessionId         string                  `json:"clientSessionId"`
}

// AccountModifyRequest represents all parameters of account modification request
type AccountModifyRequest struct {
	Id                      int64                   `json:"id,string" binding:"required,min=0"`
	Name                    string                  `json:"name" binding:"required,notBlank,max=64"`
	Category                AccountCategory         `json:"category" binding:"required"`
	Icon                    int64                   `json:"icon,string" binding:"min=1"`
	Color                   string                  `json:"color" binding:"required,len=6,validHexRGBColor"`
	Currency                *string                 `json:"currency" binding:"omitempty,len=3,validCurrency"`
	Balance                 *int64                  `json:"balance" binding:"omitempty"`
	BalanceTime             *int64                  `json:"balanceTime" binding:"omitempty"`
	Comment                 string                  `json:"comment" binding:"max=255"`
	CreditCardStatementDate int                     `json:"creditCardStatementDate" binding:"min=0,max=28"`
	Hidden                  bool                    `json:"hidden"`
	SubAccounts             []*AccountModifyRequest `json:"subAccounts" binding:"omitempty"`
	ClientSessionId         string                  `json:"clientSessionId"`
}

// AccountListRequest represents all parameters of account listing request
type AccountListRequest struct {
	VisibleOnly bool `form:"visible_only"`
}

// AccountGetRequest represents all parameters of account getting request
type AccountGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// AccountHideRequest represents all parameters of account hiding request
type AccountHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

// AccountMoveRequest represents all parameters of account moving request
type AccountMoveRequest struct {
	NewDisplayOrders []*AccountNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// AccountNewDisplayOrderRequest represents a data pair of id and display order
type AccountNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// AccountDeleteRequest represents all parameters of account deleting request
type AccountDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// AccountInfoResponse represents a view-object of account
type AccountInfoResponse struct {
	Id                      int64                    `json:"id,string"`
	Name                    string                   `json:"name"`
	ParentId                int64                    `json:"parentId,string"`
	Category                AccountCategory          `json:"category"`
	Type                    AccountType              `json:"type"`
	Icon                    int64                    `json:"icon,string"`
	Color                   string                   `json:"color"`
	Currency                string                   `json:"currency"`
	Balance                 int64                    `json:"balance"`
	Comment                 string                   `json:"comment"`
	CreditCardStatementDate *int                     `json:"creditCardStatementDate,omitempty"`
	DisplayOrder            int32                    `json:"displayOrder"`
	IsAsset                 bool                     `json:"isAsset,omitempty"`
	IsLiability             bool                     `json:"isLiability,omitempty"`
	Hidden                  bool                     `json:"hidden"`
	SubAccounts             AccountInfoResponseSlice `json:"subAccounts,omitempty"`
}

// ToAccountInfoResponse returns a view-object according to database model
func (a *Account) ToAccountInfoResponse() *AccountInfoResponse {
	var creditCardStatementDate *int

	if a.ParentAccountId == LevelOneAccountParentId && a.Category == ACCOUNT_CATEGORY_CREDIT_CARD {
		if a.Extend != nil {
			creditCardStatementDate = a.Extend.CreditCardStatementDate
		} else {
			creditCardStatementDate = &defaultCreditCardAccountStatementDate
		}
	}

	return &AccountInfoResponse{
		Id:                      a.AccountId,
		Name:                    a.Name,
		ParentId:                a.ParentAccountId,
		Category:                a.Category,
		Type:                    a.Type,
		Icon:                    a.Icon,
		Color:                   a.Color,
		Currency:                a.Currency,
		Balance:                 a.Balance,
		Comment:                 a.Comment,
		CreditCardStatementDate: creditCardStatementDate,
		DisplayOrder:            a.DisplayOrder,
		IsAsset:                 assetAccountCategory[a.Category],
		IsLiability:             liabilityAccountCategory[a.Category],
		Hidden:                  a.Hidden,
	}
}

// FromDB fills the fields from the data stored in database
func (a *AccountExtend) FromDB(data []byte) error {
	return json.Unmarshal(data, a)
}

// ToDB returns the actual stored data in database
func (a *AccountExtend) ToDB() ([]byte, error) {
	return json.Marshal(a)
}

// AccountInfoResponseSlice represents the slice data structure of AccountInfoResponse
type AccountInfoResponseSlice []*AccountInfoResponse

// Len returns the count of items
func (a AccountInfoResponseSlice) Len() int {
	return len(a)
}

// Swap swaps two items
func (a AccountInfoResponseSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less reports whether the first item is less than the second one
func (a AccountInfoResponseSlice) Less(i, j int) bool {
	if a[i].Category != a[j].Category {
		return a[i].Category < a[j].Category
	}

	return a[i].DisplayOrder < a[j].DisplayOrder
}
