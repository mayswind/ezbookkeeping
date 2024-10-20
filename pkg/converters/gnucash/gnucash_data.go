package gnucash

import "encoding/xml"

const gnucashCommodityCurrencySpace = "CURRENCY"
const gnucashRootAccountType = "ROOT"
const gnucashEquityAccountType = "EQUITY"
const gnucashIncomeAccountType = "INCOME"
const gnucashExpenseAccountType = "EXPENSE"

const gnucashSlotEquityType = "equity-type"
const gnucashSlotEquityTypeOpeningBalance = "opening-balance"

var gnucashAssetOrLiabilityAccountTypes = map[string]bool{
	"ASSET":      true,
	"BANK":       true,
	"CASH":       true,
	"CREDIT":     true,
	"LIABILITY":  true,
	"MUTUAL":     true,
	"PAYABLE":    true,
	"RECEIVABLE": true,
	"STOCK":      true,
}

// gnucashDatabase represents the struct of gnucash database file
type gnucashDatabase struct {
	XMLName xml.Name            `xml:"gnc-v2"`
	Counts  []*gnucashCountData `xml:"count-data"`
	Books   []*gnucashBookData  `xml:"book"`
}

// gnucashCountData represents the struct of gnucash count data
type gnucashCountData struct {
	Key   string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// gnucashBookData represents the struct of gnucash book data
type gnucashBookData struct {
	Id           string                    `xml:"id"`
	Counts       []*gnucashCountData       `xml:"count-data"`
	Accounts     []*gnucashAccountData     `xml:"account"`
	Transactions []*gnucashTransactionData `xml:"transaction"`
}

// gnucashCommodityData represents the struct of gnucash commodity data
type gnucashCommodityData struct {
	Space string `xml:"space"`
	Id    string `xml:"id"`
}

// gnucashSlotData represents the struct of gnucash slot data
type gnucashSlotData struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

// gnucashAccountData represents the struct of gnucash account data
type gnucashAccountData struct {
	Name        string                `xml:"name"`
	Id          string                `xml:"id"`
	AccountType string                `xml:"type"`
	Description string                `xml:"description"`
	ParentId    string                `xml:"parent"`
	Commodity   *gnucashCommodityData `xml:"commodity"`
	Slots       []*gnucashSlotData    `xml:"slots>slot"`
}

// gnucashTransactionData represents the struct of gnucash transaction data
type gnucashTransactionData struct {
	Id          string                         `xml:"id"`
	Currency    *gnucashCommodityData          `xml:"currency"`
	PostedDate  string                         `xml:"date-posted>date"`
	EnteredDate string                         `xml:"date-entered>date"`
	Description string                         `xml:"description"`
	Splits      []*gnucashTransactionSplitData `xml:"splits>split"`
}

// gnucashTransactionSplitData represents the struct of gnucash transaction split data
type gnucashTransactionSplitData struct {
	Id              string `xml:"id"`
	ReconciledState string `xml:"reconciled-state"`
	Value           string `xml:"value"`
	Quantity        string `xml:"quantity"`
	Account         string `xml:"account"`
}
