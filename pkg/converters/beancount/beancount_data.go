package beancount

import "strings"

const beancountEquityAccountNameOpeningBalance = "Opening-Balances"

// beancountDirective represents the Beancount directive
type beancountDirective string

// Beancount directives
const (
	beancountDirectiveOpen                  beancountDirective = "open"
	beancountDirectiveClose                 beancountDirective = "close"
	beancountDirectiveTransaction           beancountDirective = "txn"
	beancountDirectiveCompletedTransaction  beancountDirective = "*"
	beancountDirectiveInCompleteTransaction beancountDirective = "!"
	beancountDirectivePaddingTransaction    beancountDirective = "P"
	beancountDirectiveCommodity             beancountDirective = "commodity"
	beancountDirectivePrice                 beancountDirective = "price"
	beancountDirectiveNote                  beancountDirective = "note"
	beancountDirectiveDocument              beancountDirective = "document"
	beancountDirectiveEvent                 beancountDirective = "event"
	beancountDirectiveBalance               beancountDirective = "balance"
	beancountDirectivePad                   beancountDirective = "pad"
	beancountDirectiveQuery                 beancountDirective = "query"
	beancountDirectiveCustom                beancountDirective = "custom"
)

// beancountAccountType represents the Beancount account type
type beancountAccountType byte

// Beancount account types
const (
	beancountUnknownAccountType     beancountAccountType = 0
	beancountAssetsAccountType      beancountAccountType = 1
	beancountLiabilitiesAccountType beancountAccountType = 2
	beancountEquityAccountType      beancountAccountType = 3
	beancountIncomeAccountType      beancountAccountType = 4
	beancountExpensesAccountType    beancountAccountType = 5
)

// beancountData defines the structure of beancount data
type beancountData struct {
	Accounts     map[string]*beancountAccount
	Transactions []*beancountTransactionEntry
}

// beancountAccount defines the structure of beancount account
type beancountAccount struct {
	Name        string
	AccountType beancountAccountType
	OpenDate    string
	CloseDate   string
}

// beancountTransactionEntry defines the structure of beancount transaction entry
type beancountTransactionEntry struct {
	Date      string
	Directive beancountDirective
	Payee     string
	Narration string
	Postings  []*beancountPosting
	Tags      []string
	Links     []string
	Metadata  map[string]string
}

// beancountPosting defines the structure of beancount transaction posting
type beancountPosting struct {
	Account            string
	Amount             string
	OriginalAmount     string
	Commodity          string
	TotalCost          string
	TotalCostCommodity string
	Price              string
	PriceCommodity     string
	Metadata           map[string]string
}

func (a *beancountAccount) isOpeningBalanceEquityAccount() bool {
	if a.AccountType != beancountEquityAccountType {
		return false
	}

	nameItems := strings.Split(a.Name, string(beancountMetadataKeySuffix))

	if len(nameItems) != 2 {
		return false
	}

	return nameItems[1] == beancountEquityAccountNameOpeningBalance
}
