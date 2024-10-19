package qif

// qifTransactionClearedStatus represents the quicken interchange format (qif) transaction cleared status
type qifTransactionClearedStatus string

// Quicken interchange format transaction types
const (
	qifClearedStatusUnreconciled qifTransactionClearedStatus = ""
	qifClearedStatusCleared      qifTransactionClearedStatus = "C"
	qifClearedStatusReconciled   qifTransactionClearedStatus = "R"
)

// qifTransactionType represents the quicken interchange format (qif) transaction type
type qifTransactionType string

// Quicken interchange format transaction types
const (
	qifInvalidTransactionType         qifTransactionType = ""
	qifCheckTransactionType           qifTransactionType = "KC"
	qifDepositTransactionType         qifTransactionType = "KD"
	qifPaymentTransactionType         qifTransactionType = "KP"
	qifInvestmentTransactionType      qifTransactionType = "KI"
	qifElectronicPayeeTransactionType qifTransactionType = "KE"
)

// qifCategoryType represents the quicken interchange format (qif) category type
type qifCategoryType string

// Quicken interchange format category types
const (
	qifIncomeTransaction  qifCategoryType = "I"
	qifExpenseTransaction qifCategoryType = "E"
)

// qifData defines the structure of quicken interchange format (qif) data
type qifData struct {
	bankAccountTransactions       []*qifTransactionData
	cashAccountTransactions       []*qifTransactionData
	creditCardAccountTransactions []*qifTransactionData
	assetAccountTransactions      []*qifTransactionData
	liabilityAccountTransactions  []*qifTransactionData
	memorizedTransactions         []*qifMemorizedTransactionData
	investmentAccountTransactions []*qifInvestmentTransactionData
	accounts                      []*qifAccountData
	categories                    []*qifCategoryData
	classes                       []*qifClassData
}

// qifTransactionData defines the structure of quicken interchange format (qif) transaction data
type qifTransactionData struct {
	date                   string
	amount                 string
	clearedStatus          qifTransactionClearedStatus
	num                    string
	payee                  string
	memo                   string
	addresses              []string
	category               string
	subTransactionCategory []string
	subTransactionMemo     []string
	subTransactionAmount   []string
	account                *qifAccountData
}

// qifInvestmentTransactionData defines the structure of quicken interchange format (qif) investment transaction data
type qifInvestmentTransactionData struct {
	date               string
	action             string
	security           string
	price              string
	quantity           string
	amount             string
	clearedStatus      qifTransactionClearedStatus
	text               string
	memo               string
	commission         string
	accountForTransfer string
	amountTransferred  string
	account            *qifAccountData
}

// qifMemorizedTransactionData defines the structure of quicken interchange format (qif) memorized transaction data
type qifMemorizedTransactionData struct {
	qifTransactionData
	transactionType qifTransactionType
	amortization    qifMemorizedTransactionAmortizationData
}

// qifMemorizedTransactionAmortizationData defines the structure of quicken interchange format (qif) memorized transaction amortization data
type qifMemorizedTransactionAmortizationData struct {
	firstPaymentDate       string
	totalYearsForLoan      string
	numberOfPayments       string
	numberOfPeriodsPerYear string
	interestRate           string
	currentLoanBalance     string
	originalLoanAmount     string
}

// qifAccountData defines the structure of quicken interchange format (qif) account data
type qifAccountData struct {
	name                   string
	accountType            string
	description            string
	creditLimit            string
	statementBalanceDate   string
	statementBalanceAmount string
}

// qifCategoryData defines the structure of quicken interchange format (qif) category data
type qifCategoryData struct {
	name                   string
	description            string
	taxRelated             bool
	categoryType           qifCategoryType
	budgetAmount           string
	taxScheduleInformation string
}

// qifClassData defines the structure of quicken interchange format (qif) class data
type qifClassData struct {
	name        string
	description string
}
