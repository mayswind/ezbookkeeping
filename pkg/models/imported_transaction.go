package models

import "github.com/mayswind/ezbookkeeping/pkg/utils"

// ImportTransaction represents the imported transaction data
type ImportTransaction struct {
	*Transaction
	TagIds                             []string
	OriginalCategoryName               string
	OriginalSourceAccountName          string
	OriginalSourceAccountCurrency      string
	OriginalDestinationAccountName     string
	OriginalDestinationAccountCurrency string
	OriginalTagNames                   []string
}

// ImportTransactionRequest represents all parameters of the imported transaction data
type ImportTransactionRequest struct {
	Transactions []*ImportTransactionRequestItem
}

// ImportTransactionRequestItem represents a single item of the imported transaction data
type ImportTransactionRequestItem struct {
	Time                   string `json:"time"`
	UtcOffset              string `json:"utcOffset"`
	Type                   string `json:"type"`
	CategoryName           string `json:"categoryName,omitempty"`
	SourceAccountName      string `json:"sourceAccountName,omitempty"`
	DestinationAccountName string `json:"destinationAccountName,omitempty"`
	SourceAmount           string `json:"sourceAmount"`
	DestinationAmount      string `json:"destinationAmount,omitempty"`
	GeoLocation            string `json:"geoLocation,omitempty"`
	TagNames               string `json:"tagNames,omitempty"`
	Comment                string `json:"comment,omitempty"`
}

// ImportTransactionResponse represents a view-object of the imported transaction data
type ImportTransactionResponse struct {
	Type                               TransactionType                 `json:"type"`
	CategoryId                         int64                           `json:"categoryId,string"`
	OriginalCategoryName               string                          `json:"originalCategoryName"`
	Time                               int64                           `json:"time"`
	UtcOffset                          int16                           `json:"utcOffset"`
	SourceAccountId                    int64                           `json:"sourceAccountId,string"`
	OriginalSourceAccountName          string                          `json:"originalSourceAccountName"`
	OriginalSourceAccountCurrency      string                          `json:"originalSourceAccountCurrency"`
	DestinationAccountId               int64                           `json:"destinationAccountId,string,omitempty"`
	OriginalDestinationAccountName     string                          `json:"originalDestinationAccountName,omitempty"`
	OriginalDestinationAccountCurrency string                          `json:"originalDestinationAccountCurrency,omitempty"`
	SourceAmount                       int64                           `json:"sourceAmount"`
	DestinationAmount                  int64                           `json:"destinationAmount,omitempty"`
	TagIds                             []string                        `json:"tagIds"`
	OriginalTagNames                   []string                        `json:"originalTagNames"`
	Comment                            string                          `json:"comment"`
	GeoLocation                        *TransactionGeoLocationResponse `json:"geoLocation,omitempty"`
}

// ImportTransactionResponsePageWrapper represents a response of imported transaction which contains items and count
type ImportTransactionResponsePageWrapper struct {
	Items      []*ImportTransactionResponse `json:"items"`
	TotalCount int64                        `json:"totalCount"`
}

// ToImportTransactionResponse returns the a view-objects according to imported transaction data
func (t ImportTransaction) ToImportTransactionResponse() *ImportTransactionResponse {
	transactionType, err := t.Type.ToTransactionType()

	if err != nil {
		return nil
	}

	geoLocation := &TransactionGeoLocationResponse{}

	if t.GeoLongitude != 0 || t.GeoLatitude != 0 {
		geoLocation.Longitude = t.GeoLongitude
		geoLocation.Latitude = t.GeoLatitude
	} else {
		geoLocation = nil
	}

	return &ImportTransactionResponse{
		Type:                               transactionType,
		CategoryId:                         t.CategoryId,
		OriginalCategoryName:               t.OriginalCategoryName,
		Time:                               utils.GetUnixTimeFromTransactionTime(t.TransactionTime),
		UtcOffset:                          t.TimezoneUtcOffset,
		SourceAccountId:                    t.AccountId,
		OriginalSourceAccountName:          t.OriginalSourceAccountName,
		OriginalSourceAccountCurrency:      t.OriginalSourceAccountCurrency,
		DestinationAccountId:               t.RelatedAccountId,
		OriginalDestinationAccountName:     t.OriginalDestinationAccountName,
		OriginalDestinationAccountCurrency: t.OriginalDestinationAccountCurrency,
		SourceAmount:                       t.Amount,
		DestinationAmount:                  t.RelatedAccountAmount,
		TagIds:                             t.TagIds,
		OriginalTagNames:                   t.OriginalTagNames,
		Comment:                            t.Comment,
		GeoLocation:                        geoLocation,
	}
}

// ImportedTransactionSlice represents the slice data structure of import transaction data
type ImportedTransactionSlice []*ImportTransaction

// Len returns the count of items
func (s ImportedTransactionSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s ImportedTransactionSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s ImportedTransactionSlice) Less(i, j int) bool {
	if s[i].Type != s[j].Type && (s[i].Type == TRANSACTION_DB_TYPE_MODIFY_BALANCE || s[j].Type == TRANSACTION_DB_TYPE_MODIFY_BALANCE) {
		if s[i].Type == TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			return true
		} else if s[j].Type == TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			return false
		}
	}

	if s[i].TransactionTime != s[j].TransactionTime {
		return s[i].TransactionTime < s[j].TransactionTime
	}

	if s[i].Type != s[j].Type {
		return s[i].Type < s[j].Type
	}

	if s[i].OriginalCategoryName != s[j].OriginalCategoryName {
		return s[i].OriginalCategoryName < s[j].OriginalCategoryName
	}

	if s[i].OriginalSourceAccountName != s[j].OriginalSourceAccountName {
		return s[i].OriginalSourceAccountName < s[j].OriginalSourceAccountName
	}

	if s[i].Amount != s[j].Amount {
		return s[i].Amount < s[j].Amount
	}

	if s[i].Comment != s[j].Comment {
		return s[i].Comment < s[j].Comment
	}

	return false
}

// ToTransactionsList returns a list of transaction models
func (s ImportedTransactionSlice) ToTransactionsList() []*Transaction {
	transactions := make([]*Transaction, s.Len())

	for i := 0; i < s.Len(); i++ {
		transactions[i] = s[i].Transaction
	}

	return transactions
}

// ToTransactionTagIdsMap returns a list of transaction tag ids
func (s ImportedTransactionSlice) ToTransactionTagIdsMap() (map[int][]int64, error) {
	transactionTagIdsMap := make(map[int][]int64, s.Len())

	for i := 0; i < s.Len(); i++ {
		tagIds, err := utils.StringArrayToInt64Array(s[i].TagIds)

		if err != nil {
			return nil, err
		}

		transactionTagIdsMap[i] = tagIds
	}

	return transactionTagIdsMap, nil
}

// ToImportTransactionResponseList returns the a list of view-objects according to imported transaction data
func (s ImportedTransactionSlice) ToImportTransactionResponseList() []*ImportTransactionResponse {
	transactionResps := make([]*ImportTransactionResponse, 0, s.Len())

	for i := 0; i < s.Len(); i++ {
		importedTransaction := s[i]
		importedTransactionResp := importedTransaction.ToImportTransactionResponse()

		if importedTransactionResp == nil {
			continue
		}

		transactionResps = append(transactionResps, importedTransactionResp)
	}

	return transactionResps
}
