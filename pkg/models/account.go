package models

// Level-One Account
const ACCOUNT_PARENT_ID_LEVEL_ONE = 0

type AccountCategory byte

const (
	ACCOUNT_CATEGORY_CASH        AccountCategory = 1
	ACCOUNT_CATEGORY_DEBIT_CARD  AccountCategory = 2
	ACCOUNT_CATEGORY_CREDIT_CARD AccountCategory = 3
	ACCOUNT_CATEGORY_VIRTUAL     AccountCategory = 4
	ACCOUNT_CATEGORY_DEBT        AccountCategory = 5
	ACCOUNT_CATEGORY_RECEIVABLES AccountCategory = 6
	ACCOUNT_CATEGORY_INVESTMENT  AccountCategory = 7
)

type AccountType byte

const (
	ACCOUNT_TYPE_SINGLE_ACCOUNT     AccountType = 1
	ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS AccountType = 2
)

type Account struct {
	AccountId       int64           `xorm:"PK"`
	Uid             int64           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Deleted         bool            `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Category        AccountCategory `xorm:"NOT NULL"`
	Type            AccountType     `xorm:"NOT NULL"`
	ParentAccountId int64           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Name            string          `xorm:"VARCHAR(32) NOT NULL"`
	DisplayOrder    int             `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
	Icon            int64           `xorm:"NOT NULL"`
	Currency        string          `xorm:"VARCHAR(3) NOT NULL"`
	Balance         int64           `xorm:"NOT NULL"`
	Comment         string          `xorm:"VARCHAR(255) NOT NULL"`
	Hidden          bool            `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

type AccountCreateRequest struct {
	Name        string                  `json:"name" binding:"required,notBlank,max=32"`
	Category    AccountCategory         `json:"category" binding:"required"`
	Type        AccountType             `json:"type" binding:"required"`
	Icon        int64                   `json:"icon,string" binding:"required,min=1"`
	Currency    string                  `json:"currency" binding:"required,len=3,validCurrency"`
	Comment     string                  `json:"comment" binding:"max=255"`
	SubAccounts []*AccountCreateRequest `json:"subAccounts" binding:"omitempty"`
}

type AccountMoveRequest struct {
	NewDisplayOrders []*AccountNewDisplayOrderRequest `json:"newDisplayOrders"`
}

type AccountNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int   `json:"displayOrder"`
}

type AccountDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

type AccountInfoResponse struct {
	Id           int64                    `json:"id,string"`
	Name         string                   `json:"name"`
	ParentId     int64                    `json:"parentId,string"`
	Category     AccountCategory          `json:"category"`
	Type         AccountType              `json:"type"`
	Icon         int64                    `json:"icon,string"`
	Currency     string                   `json:"currency"`
	Balance      int64                    `json:"balance"`
	Comment      string                   `json:"comment"`
	DisplayOrder int                      `json:"displayOrder"`
	Hidden       bool                     `json:"hidden"`
	SubAccounts  AccountInfoResponseSlice `json:"subAccounts,omitempty"`
}

func (a *Account) ToAccountInfoResponse() *AccountInfoResponse {
	return &AccountInfoResponse{
		Id:           a.AccountId,
		Name:         a.Name,
		ParentId:     a.ParentAccountId,
		Category:     a.Category,
		Type:         a.Type,
		Icon:         a.Icon,
		Currency:     a.Currency,
		Balance:      a.Balance,
		Comment:      a.Comment,
		DisplayOrder: a.DisplayOrder,
		Hidden:       a.Hidden,
	}
}

type AccountInfoResponseSlice []*AccountInfoResponse

func (a AccountInfoResponseSlice) Len() int {
	return len(a)
}

func (a AccountInfoResponseSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a AccountInfoResponseSlice) Less(i, j int) bool {
	return a[i].DisplayOrder < a[j].DisplayOrder
}
