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
	accounts     map[string]*beancountAccount
	transactions []*beancountTransactionEntry
}

// beancountAccount defines the structure of beancount account
type beancountAccount struct {
	name        string
	accountType beancountAccountType
	openDate    string
	closeDate   string
}

// beancountTransactionEntry defines the structure of beancount transaction entry
type beancountTransactionEntry struct {
	date      string
	directive beancountDirective
	payee     string
	narration string
	postings  []*beancountPosting
	tags      []string
	links     []string
	metadata  map[string]string
}

// beancountPosting defines the structure of beancount transaction posting
type beancountPosting struct {
	account            string
	amount             string
	originalAmount     string
	commodity          string
	totalCost          string
	totalCostCommodity string
	price              string
	priceCommodity     string
	metadata           map[string]string
}

func (a *beancountAccount) isOpeningBalanceEquityAccount() bool {
	if a.accountType != beancountEquityAccountType {
		return false
	}

	nameItems := strings.Split(a.name, string(beancountMetadataKeySuffix))

	if len(nameItems) != 2 {
		return false
	}

	return nameItems[1] == beancountEquityAccountNameOpeningBalance
}
