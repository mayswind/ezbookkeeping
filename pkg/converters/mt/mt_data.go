package mt

type mtCreditDebitMark string

const (
	MT_MARK_CREDIT          mtCreditDebitMark = "C"
	MT_MARK_DEBIT           mtCreditDebitMark = "D"
	MT_MARK_REVERSAL_CREDIT mtCreditDebitMark = "RC"
	MT_MARK_REVERSAL_DEBIT  mtCreditDebitMark = "RD"
)

// mt940Data defines the structure of mt940 data
type mt940Data struct {
	StatementReferenceNumber string
	RelatedReference         string
	AccountId                string
	SequentialNumber         string
	OpeningBalance           *mtBalance
	ClosingBalance           *mtBalance
	ClosingAvailableBalance  *mtBalance
	Statements               []*mtStatement
}

// mtStatement defines the structure of mt940 statement
type mtStatement struct {
	ValueDate                              string
	EntryDate                              string
	CreditDebitMark                        mtCreditDebitMark
	FundsCode                              string
	Amount                                 string
	TransactionTypeIdentificationCode      string
	ReferenceForAccountOwner               string
	ReferenceOfAccountServicingInstitution string
	AdditionalInformation                  []string
}

// mtBalance defines the structure of mt940 balance
type mtBalance struct {
	DebitCreditMark mtCreditDebitMark
	Date            string
	Currency        string
	Amount          string
}
