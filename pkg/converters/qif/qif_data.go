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
	BankAccountTransactions       []*qifTransactionData
	CashAccountTransactions       []*qifTransactionData
	CreditCardAccountTransactions []*qifTransactionData
	AssetAccountTransactions      []*qifTransactionData
	LiabilityAccountTransactions  []*qifTransactionData
	MemorizedTransactions         []*qifMemorizedTransactionData
	InvestmentAccountTransactions []*qifInvestmentTransactionData
	Accounts                      []*qifAccountData
	Categories                    []*qifCategoryData
	Classes                       []*qifClassData
}

// qifTransactionData defines the structure of quicken interchange format (qif) transaction data
type qifTransactionData struct {
	Date                   string
	Amount                 string
	ClearedStatus          qifTransactionClearedStatus
	Num                    string
	Payee                  string
	Memo                   string
	Addresses              []string
	Category               string
	SubTransactionCategory []string
	SubTransactionMemo     []string
	SubTransactionAmount   []string
	Account                *qifAccountData
}

// qifInvestmentTransactionData defines the structure of quicken interchange format (qif) investment transaction data
type qifInvestmentTransactionData struct {
	Date               string
	Action             string
	Security           string
	Price              string
	Quantity           string
	Amount             string
	ClearedStatus      qifTransactionClearedStatus
	Text               string
	Memo               string
	Commission         string
	AccountForTransfer string
	AmountTransferred  string
	Account            *qifAccountData
}

// qifMemorizedTransactionData defines the structure of quicken interchange format (qif) memorized transaction data
type qifMemorizedTransactionData struct {
	qifTransactionData
	TransactionType qifTransactionType
	Amortization    qifMemorizedTransactionAmortizationData
}

// qifMemorizedTransactionAmortizationData defines the structure of quicken interchange format (qif) memorized transaction amortization data
type qifMemorizedTransactionAmortizationData struct {
	FirstPaymentDate       string
	TotalYearsForLoan      string
	NumberOfPayments       string
	NumberOfPeriodsPerYear string
	InterestRate           string
	CurrentLoanBalance     string
	OriginalLoanAmount     string
}

// qifAccountData defines the structure of quicken interchange format (qif) account data
type qifAccountData struct {
	Name                   string
	AccountType            string
	Description            string
	CreditLimit            string
	StatementBalanceDate   string
	StatementBalanceAmount string
}

// qifCategoryData defines the structure of quicken interchange format (qif) category data
type qifCategoryData struct {
	Name                   string
	Description            string
	TaxRelated             bool
	CategoryType           qifCategoryType
	BudgetAmount           string
	TaxScheduleInformation string
}

// qifClassData defines the structure of quicken interchange format (qif) class data
type qifClassData struct {
	Name        string
	Description string
}
