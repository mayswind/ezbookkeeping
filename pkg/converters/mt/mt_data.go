package mt

import "strings"

type mtCreditDebitMark string

const (
	MT_MARK_CREDIT          mtCreditDebitMark = "C"
	MT_MARK_DEBIT           mtCreditDebitMark = "D"
	MT_MARK_REVERSAL_CREDIT mtCreditDebitMark = "RC"
	MT_MARK_REVERSAL_DEBIT  mtCreditDebitMark = "RD"
)

const (
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_REMITTANCE string = "REMI"
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
	InformationToAccountOwner              []string
}

// mtBalance defines the structure of mt940 balance
type mtBalance struct {
	DebitCreditMark mtCreditDebitMark
	Date            string
	Currency        string
	Amount          string
}

// GetInformationToAccountOwnerMap returns a map of additional information
func (s *mtStatement) GetInformationToAccountOwnerMap() map[string]string {
	additionalInfoMap := make(map[string]string, len(s.InformationToAccountOwner))

	for _, info := range s.InformationToAccountOwner {
		items := strings.Split(info, "/")

		if len(items) < 3 {
			continue
		}

		for i := 2; i < len(items); i += 2 {
			key := strings.TrimSpace(items[i-1])
			value := strings.TrimSpace(items[i])

			if len(key) > 0 {
				additionalInfoMap[key] = value
			}
		}
	}

	return additionalInfoMap
}
