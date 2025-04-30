package models

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// TransactionType represents transaction type
type TransactionType byte

// Transaction types
const (
	TRANSACTION_TYPE_MODIFY_BALANCE TransactionType = 1
	TRANSACTION_TYPE_INCOME         TransactionType = 2
	TRANSACTION_TYPE_EXPENSE        TransactionType = 3
	TRANSACTION_TYPE_TRANSFER       TransactionType = 4
)

// TransactionDbType represents transaction type in database
type TransactionDbType byte

// Transaction db types
const (
	TRANSACTION_DB_TYPE_MODIFY_BALANCE TransactionDbType = 1
	TRANSACTION_DB_TYPE_INCOME         TransactionDbType = 2
	TRANSACTION_DB_TYPE_EXPENSE        TransactionDbType = 3
	TRANSACTION_DB_TYPE_TRANSFER_OUT   TransactionDbType = 4
	TRANSACTION_DB_TYPE_TRANSFER_IN    TransactionDbType = 5
)

// String returns a textual representation of the transaction types for db enum
func (s TransactionDbType) String() string {
	switch s {
	case TRANSACTION_DB_TYPE_MODIFY_BALANCE:
		return "Modify Balance"
	case TRANSACTION_DB_TYPE_INCOME:
		return "Income"
	case TRANSACTION_DB_TYPE_EXPENSE:
		return "Expense"
	case TRANSACTION_DB_TYPE_TRANSFER_OUT:
		return "Transfer Out"
	case TRANSACTION_DB_TYPE_TRANSFER_IN:
		return "Transfer In"
	default:
		return fmt.Sprintf("Invalid(%d)", int(s))
	}
}

// ToTransactionType returns the transaction type for this db enum
func (s TransactionDbType) ToTransactionType() (TransactionType, error) {
	if s == TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		return TRANSACTION_TYPE_MODIFY_BALANCE, nil
	} else if s == TRANSACTION_DB_TYPE_EXPENSE {
		return TRANSACTION_TYPE_EXPENSE, nil
	} else if s == TRANSACTION_DB_TYPE_INCOME {
		return TRANSACTION_TYPE_INCOME, nil
	} else if s == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		return TRANSACTION_TYPE_TRANSFER, nil
	} else if s == TRANSACTION_DB_TYPE_TRANSFER_IN {
		return TRANSACTION_TYPE_TRANSFER, nil
	} else {
		return 0, errs.ErrTransactionTypeInvalid
	}
}

// TransactionTagFilterType represents transaction tag filter type
type TransactionTagFilterType byte

// Transaction tag filter types
const (
	TRANSACTION_TAG_FILTER_HAS_ANY     TransactionTagFilterType = 0
	TRANSACTION_TAG_FILTER_HAS_ALL     TransactionTagFilterType = 1
	TRANSACTION_TAG_FILTER_NOT_HAS_ANY TransactionTagFilterType = 2
	TRANSACTION_TAG_FILTER_NOT_HAS_ALL TransactionTagFilterType = 3
)

// Transaction represents transaction data stored in database
type Transaction struct {
	TransactionId        int64             `xorm:"PK"`
	Uid                  int64             `xorm:"UNIQUE(UQE_transaction_uid_time) INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) INDEX(IDX_transaction_uid_deleted_time_longitude_latitude) NOT NULL"`
	Deleted              bool              `xorm:"INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) INDEX(IDX_transaction_uid_deleted_time_longitude_latitude) NOT NULL"`
	Type                 TransactionDbType `xorm:"INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) NOT NULL"`
	CategoryId           int64             `xorm:"INDEX(IDX_transaction_uid_deleted_category_id_time) NOT NULL"`
	AccountId            int64             `xorm:"INDEX(IDX_transaction_uid_deleted_account_id_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) NOT NULL"`
	TransactionTime      int64             `xorm:"UNIQUE(UQE_transaction_uid_time) INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) NOT NULL"`
	TimezoneUtcOffset    int16             `xorm:"NOT NULL"`
	Amount               int64             `xorm:"NOT NULL"`
	RelatedId            int64             `xorm:"NOT NULL"`
	RelatedAccountId     int64             `xorm:"NOT NULL"`
	RelatedAccountAmount int64             `xorm:"NOT NULL"`
	HideAmount           bool              `xorm:"NOT NULL"`
	Comment              string            `xorm:"VARCHAR(255) NOT NULL"`
	GeoLongitude         float64           `xorm:"INDEX(IDX_transaction_uid_deleted_time_longitude_latitude)"`
	GeoLatitude          float64           `xorm:"INDEX(IDX_transaction_uid_deleted_time_longitude_latitude)"`
	CreatedIp            string            `xorm:"VARCHAR(39)"`
	ScheduledCreated     bool
	CreatedUnixTime      int64
	UpdatedUnixTime      int64
	DeletedUnixTime      int64
}

// TransactionGeoLocationRequest represents all parameters of transaction geographic location info update request
type TransactionGeoLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// TransactionCreateRequest represents all parameters of transaction creation request
type TransactionCreateRequest struct {
	Type                 TransactionType                `json:"type" binding:"required"`
	CategoryId           int64                          `json:"categoryId,string"`
	Time                 int64                          `json:"time" binding:"required,min=1"`
	UtcOffset            int16                          `json:"utcOffset" binding:"min=-720,max=840"`
	SourceAccountId      int64                          `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId int64                          `json:"destinationAccountId,string" binding:"min=0"`
	SourceAmount         int64                          `json:"sourceAmount" binding:"min=-99999999999,max=99999999999"`
	DestinationAmount    int64                          `json:"destinationAmount" binding:"min=-99999999999,max=99999999999"`
	HideAmount           bool                           `json:"hideAmount"`
	TagIds               []string                       `json:"tagIds"`
	PictureIds           []string                       `json:"pictureIds"`
	Comment              string                         `json:"comment" binding:"max=255"`
	GeoLocation          *TransactionGeoLocationRequest `json:"geoLocation" binding:"omitempty"`
	ClientSessionId      string                         `json:"clientSessionId"`
}

// TransactionModifyRequest represents all parameters of transaction modification request
type TransactionModifyRequest struct {
	Id                   int64                          `json:"id,string" binding:"required,min=1"`
	CategoryId           int64                          `json:"categoryId,string"`
	Time                 int64                          `json:"time" binding:"required,min=1"`
	UtcOffset            int16                          `json:"utcOffset" binding:"min=-720,max=840"`
	SourceAccountId      int64                          `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId int64                          `json:"destinationAccountId,string" binding:"min=0"`
	SourceAmount         int64                          `json:"sourceAmount" binding:"min=-99999999999,max=99999999999"`
	DestinationAmount    int64                          `json:"destinationAmount" binding:"min=-99999999999,max=99999999999"`
	HideAmount           bool                           `json:"hideAmount"`
	TagIds               []string                       `json:"tagIds"`
	PictureIds           []string                       `json:"pictureIds"`
	Comment              string                         `json:"comment" binding:"max=255"`
	GeoLocation          *TransactionGeoLocationRequest `json:"geoLocation" binding:"omitempty"`
}

// TransactionImportRequest represents all parameters of transaction import request
type TransactionImportRequest struct {
	Transactions    []*TransactionCreateRequest `json:"transactions"`
	ClientSessionId string                      `json:"clientSessionId"`
}

// TransactionImportProcessRequest represents all parameters of transaction import process request
type TransactionImportProcessRequest struct {
	ClientSessionId string `form:"client_session_id"`
}

// TransactionCountRequest represents transaction count request
type TransactionCountRequest struct {
	Type          TransactionDbType        `form:"type" binding:"min=0,max=4"`
	CategoryIds   string                   `form:"category_ids"`
	AccountIds    string                   `form:"account_ids"`
	TagIds        string                   `form:"tag_ids"`
	TagFilterType TransactionTagFilterType `form:"tag_filter_type" binding:"min=0,max=3"`
	AmountFilter  string                   `form:"amount_filter" binding:"validAmountFilter"`
	Keyword       string                   `form:"keyword"`
	MaxTime       int64                    `form:"max_time" binding:"min=0"`
	MinTime       int64                    `form:"min_time" binding:"min=0"`
}

// TransactionListByMaxTimeRequest represents all parameters of transaction listing by max time request
type TransactionListByMaxTimeRequest struct {
	Type          TransactionDbType        `form:"type" binding:"min=0,max=4"`
	CategoryIds   string                   `form:"category_ids"`
	AccountIds    string                   `form:"account_ids"`
	TagIds        string                   `form:"tag_ids"`
	TagFilterType TransactionTagFilterType `form:"tag_filter_type" binding:"min=0,max=3"`
	AmountFilter  string                   `form:"amount_filter" binding:"validAmountFilter"`
	Keyword       string                   `form:"keyword"`
	MaxTime       int64                    `form:"max_time" binding:"min=0"`
	MinTime       int64                    `form:"min_time" binding:"min=0"`
	Page          int32                    `form:"page" binding:"min=0"`
	Count         int32                    `form:"count" binding:"required,min=1,max=50"`
	WithCount     bool                     `form:"with_count"`
	WithPictures  bool                     `form:"with_pictures"`
	TrimAccount   bool                     `form:"trim_account"`
	TrimCategory  bool                     `form:"trim_category"`
	TrimTag       bool                     `form:"trim_tag"`
}

// TransactionListInMonthByPageRequest represents all parameters of transaction listing by month request
type TransactionListInMonthByPageRequest struct {
	Year          int32                    `form:"year" binding:"required,min=1"`
	Month         int32                    `form:"month" binding:"required,min=1"`
	Type          TransactionDbType        `form:"type" binding:"min=0,max=4"`
	CategoryIds   string                   `form:"category_ids"`
	AccountIds    string                   `form:"account_ids"`
	TagIds        string                   `form:"tag_ids"`
	TagFilterType TransactionTagFilterType `form:"tag_filter_type" binding:"min=0,max=3"`
	AmountFilter  string                   `form:"amount_filter" binding:"validAmountFilter"`
	Keyword       string                   `form:"keyword"`
	WithPictures  bool                     `form:"with_pictures"`
	TrimAccount   bool                     `form:"trim_account"`
	TrimCategory  bool                     `form:"trim_category"`
	TrimTag       bool                     `form:"trim_tag"`
}

// TransactionStatisticRequest represents all parameters of transaction statistic request
type TransactionStatisticRequest struct {
	StartTime              int64                    `form:"start_time" binding:"min=0"`
	EndTime                int64                    `form:"end_time" binding:"min=0"`
	TagIds                 string                   `form:"tag_ids"`
	TagFilterType          TransactionTagFilterType `form:"tag_filter_type" binding:"min=0,max=3"`
	UseTransactionTimezone bool                     `form:"use_transaction_timezone"`
}

// TransactionStatisticTrendsRequest represents all parameters of transaction statistic trends request
type TransactionStatisticTrendsRequest struct {
	YearMonthRangeRequest
	TagIds                 string                   `form:"tag_ids"`
	TagFilterType          TransactionTagFilterType `form:"tag_filter_type" binding:"min=0,max=3"`
	UseTransactionTimezone bool                     `form:"use_transaction_timezone"`
}

// TransactionAmountsRequest represents all parameters of transaction amounts request
type TransactionAmountsRequest struct {
	Query                  string `form:"query"`
	UseTransactionTimezone bool   `form:"use_transaction_timezone"`
}

// TransactionAmountsRequestItem represents an item of transaction amounts request
type TransactionAmountsRequestItem struct {
	Name      string
	StartTime int64
	EndTime   int64
}

// TransactionGetRequest represents all parameters of transaction getting request
type TransactionGetRequest struct {
	Id           int64 `form:"id,string" binding:"required,min=1"`
	WithPictures bool  `form:"with_pictures"`
	TrimAccount  bool  `form:"trim_account"`
	TrimCategory bool  `form:"trim_category"`
	TrimTag      bool  `form:"trim_tag"`
}

// TransactionDeleteRequest represents all parameters of transaction deleting request
type TransactionDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// YearMonthRangeRequest represents all parameters of a request with year and month range
type YearMonthRangeRequest struct {
	StartYearMonth string `form:"start_year_month"`
	EndYearMonth   string `form:"end_year_month"`
}

// TransactionGeoLocationResponse represents a view-object of transaction geographic location info
type TransactionGeoLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TransactionInfoResponse represents a view-object of transaction
type TransactionInfoResponse struct {
	Id                   int64                                    `json:"id,string"`
	TimeSequenceId       int64                                    `json:"timeSequenceId,string"`
	Type                 TransactionType                          `json:"type"`
	CategoryId           int64                                    `json:"categoryId,string"`
	Category             *TransactionCategoryInfoResponse         `json:"category,omitempty"`
	Time                 int64                                    `json:"time"`
	UtcOffset            int16                                    `json:"utcOffset"`
	SourceAccountId      int64                                    `json:"sourceAccountId,string"`
	SourceAccount        *AccountInfoResponse                     `json:"sourceAccount,omitempty"`
	DestinationAccountId int64                                    `json:"destinationAccountId,string,omitempty"`
	DestinationAccount   *AccountInfoResponse                     `json:"destinationAccount,omitempty"`
	SourceAmount         int64                                    `json:"sourceAmount"`
	DestinationAmount    int64                                    `json:"destinationAmount,omitempty"`
	HideAmount           bool                                     `json:"hideAmount"`
	TagIds               []string                                 `json:"tagIds"`
	Tags                 []*TransactionTagInfoResponse            `json:"tags,omitempty"`
	Pictures             TransactionPictureInfoBasicResponseSlice `json:"pictures,omitempty"`
	Comment              string                                   `json:"comment"`
	GeoLocation          *TransactionGeoLocationResponse          `json:"geoLocation,omitempty"`
	Editable             bool                                     `json:"editable"`
}

// TransactionCountResponse represents transaction count response
type TransactionCountResponse struct {
	TotalCount int64 `json:"totalCount"`
}

// TransactionInfoPageWrapperResponse represents a response of transaction which contains items and next id
type TransactionInfoPageWrapperResponse struct {
	Items              TransactionInfoResponseSlice `json:"items"`
	NextTimeSequenceId *int64                       `json:"nextTimeSequenceId,string"`
	TotalCount         *int64                       `json:"totalCount,omitempty"`
}

// TransactionInfoPageWrapperResponse2 represents a response of transaction which contains items and count
type TransactionInfoPageWrapperResponse2 struct {
	Items      TransactionInfoResponseSlice `json:"items"`
	TotalCount int64                        `json:"totalCount"`
}

// TransactionStatisticResponse represents transaction statistic response
type TransactionStatisticResponse struct {
	StartTime int64                               `json:"startTime"`
	EndTime   int64                               `json:"endTime"`
	Items     []*TransactionStatisticResponseItem `json:"items"`
}

// TransactionStatisticResponseItem represents total amount item for a response
type TransactionStatisticResponseItem struct {
	CategoryId  int64 `json:"categoryId,string"`
	AccountId   int64 `json:"accountId,string"`
	TotalAmount int64 `json:"amount"`
}

// TransactionStatisticTrendsResponseItem represents the data within each statistic interval
type TransactionStatisticTrendsResponseItem struct {
	Year  int32                               `json:"year"`
	Month int32                               `json:"month"`
	Items []*TransactionStatisticResponseItem `json:"items"`
}

// TransactionAmountsResponseItem represents an item of transaction amounts
type TransactionAmountsResponseItem struct {
	StartTime int64                                       `json:"startTime"`
	EndTime   int64                                       `json:"endTime"`
	Amounts   []*TransactionAmountsResponseItemAmountInfo `json:"amounts"`
}

// TransactionMonthAmountsResponseItem represents an item of transaction month amounts
type TransactionMonthAmountsResponseItem struct {
	Year    int32                                       `json:"year"`
	Month   int32                                       `json:"month"`
	Amounts []*TransactionAmountsResponseItemAmountInfo `json:"amounts"`
}

// TransactionAmountsResponseItemAmountInfo represents amount info for a response item
type TransactionAmountsResponseItemAmountInfo struct {
	Currency      string `json:"currency"`
	IncomeAmount  int64  `json:"incomeAmount"`
	ExpenseAmount int64  `json:"expenseAmount"`
}

// IsEditable returns whether this transaction can be edited
func (t *Transaction) IsEditable(currentUser *User, utcOffset int16, account *Account, relatedAccount *Account) bool {
	if currentUser == nil || !currentUser.CanEditTransactionByTransactionTime(t.TransactionTime, utcOffset) {
		return false
	}

	if account == nil || account.Hidden {
		return false
	}

	if t.Type == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if relatedAccount == nil || relatedAccount.Hidden {
			return false
		}
	}

	return true
}

// ToTransactionInfoResponse returns a view-object according to database model
func (t *Transaction) ToTransactionInfoResponse(tagIds []int64, editable bool) *TransactionInfoResponse {
	transactionType, err := t.Type.ToTransactionType()

	if err != nil {
		return nil
	}

	sourceAccountId := t.AccountId
	sourceAmount := t.Amount

	destinationAccountId := int64(0)
	destinationAmount := int64(0)

	if t.Type == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		destinationAccountId = t.RelatedAccountId
		destinationAmount = t.RelatedAccountAmount
	} else if t.Type == TRANSACTION_DB_TYPE_TRANSFER_IN {
		sourceAccountId = t.RelatedAccountId
		sourceAmount = t.RelatedAccountAmount

		destinationAccountId = t.AccountId
		destinationAmount = t.Amount
	}

	geoLocation := &TransactionGeoLocationResponse{}

	if t.GeoLongitude != 0 || t.GeoLatitude != 0 {
		geoLocation.Longitude = t.GeoLongitude
		geoLocation.Latitude = t.GeoLatitude
	} else {
		geoLocation = nil
	}

	return &TransactionInfoResponse{
		Id:                   t.TransactionId,
		TimeSequenceId:       t.TransactionTime,
		Type:                 transactionType,
		CategoryId:           t.CategoryId,
		Time:                 utils.GetUnixTimeFromTransactionTime(t.TransactionTime),
		UtcOffset:            t.TimezoneUtcOffset,
		SourceAccountId:      sourceAccountId,
		DestinationAccountId: destinationAccountId,
		SourceAmount:         sourceAmount,
		DestinationAmount:    destinationAmount,
		HideAmount:           t.HideAmount,
		TagIds:               utils.Int64ArrayToStringArray(tagIds),
		Comment:              t.Comment,
		GeoLocation:          geoLocation,
		Editable:             editable,
	}
}

// GetTransactionAmountsRequestItems returns request items by query parameters
func (t *TransactionAmountsRequest) GetTransactionAmountsRequestItems() ([]*TransactionAmountsRequestItem, error) {
	items := strings.Split(t.Query, "|")
	requestItems := make([]*TransactionAmountsRequestItem, 0, len(items))

	for i := 0; i < len(items); i++ {
		itemValues := strings.Split(items[i], "_")

		if len(itemValues) != 3 {
			return nil, errs.ErrQueryItemsInvalid
		}

		startTime, err := utils.StringToInt64(itemValues[1])

		if err != nil {
			return nil, err
		}

		endTime, err := utils.StringToInt64(itemValues[2])

		if err != nil {
			return nil, err
		}

		requestItem := &TransactionAmountsRequestItem{
			Name:      itemValues[0],
			StartTime: startTime,
			EndTime:   endTime,
		}

		requestItems = append(requestItems, requestItem)
	}

	return requestItems, nil
}

// GetNumericYearMonthRange returns numeric start year, start month, end year and end month
func (t *YearMonthRangeRequest) GetNumericYearMonthRange() (int32, int32, int32, int32, error) {
	var startYear, startMonth, endYear, endMonth int32
	var err error

	if t.StartYearMonth != "" {
		startYear, startMonth, err = utils.ParseNumericYearMonth(t.StartYearMonth)

		if err != nil {
			return 0, 0, 0, 0, err
		}
	}

	if t.EndYearMonth != "" {
		endYear, endMonth, err = utils.ParseNumericYearMonth(t.EndYearMonth)

		if err != nil {
			return 0, 0, 0, 0, err
		}
	}

	return startYear, startMonth, endYear, endMonth, nil
}

// TransactionInfoResponseSlice represents the slice data structure of TransactionInfoResponse
type TransactionInfoResponseSlice []*TransactionInfoResponse

// Len returns the count of items
func (s TransactionInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionInfoResponseSlice) Less(i, j int) bool {
	if s[i].Time != s[j].Time {
		return s[i].Time > s[j].Time
	}

	return s[i].Id > s[j].Id
}

// TransactionStatisticTrendsResponseItemSlice represents the slice data structure of TransactionStatisticTrendsResponseItem
type TransactionStatisticTrendsResponseItemSlice []*TransactionStatisticTrendsResponseItem

// Len returns the count of items
func (s TransactionStatisticTrendsResponseItemSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionStatisticTrendsResponseItemSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionStatisticTrendsResponseItemSlice) Less(i, j int) bool {
	if s[i].Year != s[j].Year {
		return s[i].Year < s[j].Year
	}

	return s[i].Month < s[j].Month
}

// TransactionAmountsResponseItemAmountInfoSlice represents the slice data structure of TransactionAmountsResponseItemAmountInfo
type TransactionAmountsResponseItemAmountInfoSlice []*TransactionAmountsResponseItemAmountInfo

// Len returns the count of items
func (s TransactionAmountsResponseItemAmountInfoSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionAmountsResponseItemAmountInfoSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionAmountsResponseItemAmountInfoSlice) Less(i, j int) bool {
	return strings.Compare(s[i].Currency, s[j].Currency) < 0
}
