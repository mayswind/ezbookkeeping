package ofx

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// oFXDeclarationVersion represents the declaration version of open financial exchange (ofx) file
type oFXDeclarationVersion string

const (
	ofxVersion1 oFXDeclarationVersion = "100"
	ofxVersion2 oFXDeclarationVersion = "200"
)

const ofxDefaultTimezoneOffset = "+00:00"

// ofxAccountType represents account type in open financial exchange (ofx) file
type ofxAccountType string

// OFX account types
const (
	ofxCheckingAccount             ofxAccountType = "CHECKING"
	ofxSavingsAccount              ofxAccountType = "SAVINGS"
	ofxMoneyMarketAccount          ofxAccountType = "MONEYMRKT"
	ofxLineOfCreditAccount         ofxAccountType = "CREDITLINE"
	ofxCertificateOfDepositAccount ofxAccountType = "CD"
)

// ofxTransactionType represents transaction type in open financial exchange (ofx) file
type ofxTransactionType string

// OFX transaction types
const (
	ofxGenericCreditTransaction          ofxTransactionType = "CREDIT"
	ofxGenericDebitTransaction           ofxTransactionType = "DEBIT"
	ofxInterestTransaction               ofxTransactionType = "INT"
	ofxDividendTransaction               ofxTransactionType = "DIV"
	ofxFIFeeTransaction                  ofxTransactionType = "FEE"
	ofxServiceChargeTransaction          ofxTransactionType = "SRVCHG"
	ofxDepositTransaction                ofxTransactionType = "DEP"
	ofxATMTransaction                    ofxTransactionType = "ATM"
	ofxPOSTransaction                    ofxTransactionType = "POS"
	ofxTransferTransaction               ofxTransactionType = "XFER"
	ofxCheckTransaction                  ofxTransactionType = "CHECK"
	ofxElectronicPaymentTransaction      ofxTransactionType = "PAYMENT"
	ofxCashWithdrawalTransaction         ofxTransactionType = "CASH"
	ofxDirectDepositTransaction          ofxTransactionType = "DIRECTDEP"
	ofxMerchantInitiatedDebitTransaction ofxTransactionType = "DIRECTDEBIT"
	ofxRepeatingPaymentTransaction       ofxTransactionType = "REPEATPMT"
	ofxHoldTransaction                   ofxTransactionType = "HOLD"
	ofxOtherTransaction                  ofxTransactionType = "OTHER"
)

var ofxTransactionTypeMapping = map[ofxTransactionType]models.TransactionType{
	ofxGenericCreditTransaction:          models.TRANSACTION_TYPE_EXPENSE,
	ofxGenericDebitTransaction:           models.TRANSACTION_TYPE_EXPENSE,
	ofxDividendTransaction:               models.TRANSACTION_TYPE_INCOME,
	ofxFIFeeTransaction:                  models.TRANSACTION_TYPE_EXPENSE,
	ofxServiceChargeTransaction:          models.TRANSACTION_TYPE_EXPENSE,
	ofxDepositTransaction:                models.TRANSACTION_TYPE_INCOME,
	ofxTransferTransaction:               models.TRANSACTION_TYPE_TRANSFER,
	ofxCheckTransaction:                  models.TRANSACTION_TYPE_EXPENSE,
	ofxElectronicPaymentTransaction:      models.TRANSACTION_TYPE_EXPENSE,
	ofxCashWithdrawalTransaction:         models.TRANSACTION_TYPE_EXPENSE,
	ofxDirectDepositTransaction:          models.TRANSACTION_TYPE_INCOME,
	ofxMerchantInitiatedDebitTransaction: models.TRANSACTION_TYPE_EXPENSE,
	ofxRepeatingPaymentTransaction:       models.TRANSACTION_TYPE_EXPENSE,
}

// ofxFile represents the struct of open financial exchange (ofx) file
type ofxFile struct {
	FileHeader                  *ofxFileHeader
	BankMessageResponseV1       *ofxBankMessageResponseV1
	CreditCardMessageResponseV1 *ofxCreditCardMessageResponseV1
}

// ofxFileHeader represents the struct of open financial exchange (ofx) file header
type ofxFileHeader struct {
	OFXDeclarationVersion oFXDeclarationVersion
	OFXDataVersion        string
	Security              string
	OldFileUid            string
	NewFileUid            string
}

// ofxBankMessageResponseV1 represents the struct of open financial exchange (ofx) bank message response v1
type ofxBankMessageResponseV1 struct {
	StatementTransactionResponse *ofxBankStatementTransactionResponse
}

// ofxCreditCardMessageResponseV1 represents the struct of open financial exchange (ofx) credit card message response v1
type ofxCreditCardMessageResponseV1 struct {
	StatementTransactionResponse *ofxCreditCardStatementTransactionResponse
}

// ofxBankStatementTransactionResponse represents the struct of open financial exchange (ofx) bank statement transaction response
type ofxBankStatementTransactionResponse struct {
	StatementResponse *ofxBankStatementResponse
}

// ofxCreditCardStatementTransactionResponse represents the struct of open financial exchange (ofx) credit card statement transaction response
type ofxCreditCardStatementTransactionResponse struct {
	StatementResponse *ofxCreditCardStatementResponse
}

// ofxBankStatementResponse represents the struct of open financial exchange (ofx) bank statement response
type ofxBankStatementResponse struct {
	DefaultCurrency string
	AccountFrom     *ofxBankAccount
	TransactionList *ofxBankTransactionList
}

// ofxCreditCardStatementResponse represents the struct of open financial exchange (ofx) credit card statement response
type ofxCreditCardStatementResponse struct {
	DefaultCurrency string
	AccountFrom     *ofxCreditCardAccount
	TransactionList *ofxCreditCardTransactionList
}

// ofxBankAccount represents the struct of open financial exchange (ofx) bank account
type ofxBankAccount struct {
	BankId      string
	BranchId    string
	AccountId   string
	AccountType ofxAccountType
	AccountKey  string
}

// ofxCreditCardAccount represents the struct of open financial exchange (ofx) credit card account
type ofxCreditCardAccount struct {
	AccountId  string
	AccountKey string
}

// ofxBankTransactionList represents the struct of open financial exchange (ofx) bank transaction list
type ofxBankTransactionList struct {
	StartDate             string
	EndDate               string
	StatementTransactions []*ofxBankStatementTransaction
}

// ofxCreditCardTransactionList represents the struct of open financial exchange (ofx) credit card transaction list
type ofxCreditCardTransactionList struct {
	StartDate             string
	EndDate               string
	StatementTransactions []*ofxCreditCardStatementTransaction
}

// ofxBaseStatementTransaction represents the struct of open financial exchange (ofx) base statement transaction
type ofxBaseStatementTransaction struct {
	TransactionId    string
	TransactionType  ofxTransactionType
	PostedDate       string
	Amount           string
	Name             string
	Payee            *ofxPayee
	Memo             string
	Currency         string
	OriginalCurrency string
}

// ofxBankStatementTransaction represents the struct of open financial exchange (ofx) bank statement transaction
type ofxBankStatementTransaction struct {
	ofxBaseStatementTransaction
	AccountTo *ofxBankAccount
}

// ofxCreditCardStatementTransaction represents the struct of open financial exchange (ofx) credit card statement transaction
type ofxCreditCardStatementTransaction struct {
	ofxBaseStatementTransaction
	AccountTo *ofxCreditCardAccount
}

// ofxPayee represents the struct of open financial exchange (ofx) payee info
type ofxPayee struct {
	Name       string
	Address1   string
	Address2   string
	Address3   string
	City       string
	State      string
	PostalCode string
	Country    string
	Phone      string
}
