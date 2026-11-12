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
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_RETURN_REASON          string = "RTRN"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_CLIENT_REFERENCE       string = "CREF"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_END_TO_END_REFERENCE   string = "EREF"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_PAYMENT_INFORMATION_ID string = "PREF"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_INSTRUCTION_ID         string = "IREF"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_MANDATE_REFERENCE      string = "MARF"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_CREDITOR_REFERENCE     string = "CSID"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_COUNTERPARTY_ID        string = "CNTP"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_REMITTANCE             string = "REMI"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_PURPOSE                string = "PURP"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_ULTIMATE_CREDITOR      string = "ULTC"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_ULTIMATE_DEBTOR        string = "ULTD"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_EXCHANGE_RATE          string = "EXCH"
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_CHARGES                string = "CHGS"
)

const (
	MT_REMITTANCE_CODE_UNSTRUCTURED string = "USTD"
	MT_REMITTANCE_CODE_STRUCTURED   string = "STRD"
)

// Reference: https://assets.ing.com/m/5f4411ca5643dd0d/original/MT940-MT942-Format-description-for-IBP.pdf
var mtInformationToAccountOwnerTagSubFieldCounts = map[string]int{
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_RETURN_REASON:          1, // Return reason: [4a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_CLIENT_REFERENCE:       1, // Customer reference: [35a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_END_TO_END_REFERENCE:   1, // End to end ID: [35a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_PAYMENT_INFORMATION_ID: 1, // Payment information ID: [35a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_INSTRUCTION_ID:         1, // Instruction ID: [35a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_MANDATE_REFERENCE:      1, // Mandate reference: [35a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_CREDITOR_REFERENCE:     1, // Creditor reference: [35a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_COUNTERPARTY_ID:        4, // Counterparty ID: [35a] (Account number) / [11a] (BIC) / [50a] (Name) / [35a] (City)
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_REMITTANCE:             3, // Remittance information: [255a] (Code /Issuer / Remittance info)
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_PURPOSE:                1, // Purpose code: [4a]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_ULTIMATE_CREDITOR:      2, // Ultimate creditor: [70a] (Name) / [35a] (ID)
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_ULTIMATE_DEBTOR:        2, // Ultimate debtor: [70a] (Name) / [35a] (ID)
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_EXCHANGE_RATE:          1, // Exchange rate: [12d]
	MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_CHARGES:                1, // Charges: [15d] (ISO Currency code,Charges amount)
}

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
	information := s.InformationToAccountOwner
	structured := s.isInformationToAccountOwnerStructured()

	if structured {
		information = []string{strings.Join(information, "")}
	}

	for _, info := range information {
		items := strings.Split(info, "/")

		if len(items) < 3 {
			continue
		}

		for i := 1; i+1 < len(items); {
			key := strings.TrimSpace(items[i])

			if structured && len(key) < 1 {
				i++
				continue
			}

			valueStartIndex := i + 1
			valueEndIndex := valueStartIndex + 1
			subFieldCount := mtInformationToAccountOwnerTagSubFieldCounts[key]

			if key == MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_REMITTANCE && items[valueStartIndex] != MT_REMITTANCE_CODE_UNSTRUCTURED && items[valueStartIndex] != MT_REMITTANCE_CODE_STRUCTURED {
				subFieldCount = 1
			}

			if structured && subFieldCount > 0 {
				valueEndIndex = valueStartIndex

				for valueEndIndex < len(items) {
					if valueEndIndex-valueStartIndex >= subFieldCount {
						if key == MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_COUNTERPARTY_ID {
							break
						}

						// A field terminator followed by /<code word>/ separates values.
						if len(items[valueEndIndex-1]) < 1 {
							if _, exists := mtInformationToAccountOwnerTagSubFieldCounts[items[valueEndIndex]]; exists {
								break
							}
						}
					}

					valueEndIndex++
				}
			}

			nextIndex := valueEndIndex

			if structured && subFieldCount > 0 && valueEndIndex-valueStartIndex > subFieldCount && len(items[valueEndIndex-1]) < 1 {
				valueEndIndex-- // Remove only the field terminator, preserving empty subfields.
			}

			value := strings.TrimSpace(strings.Join(items[valueStartIndex:valueEndIndex], "/"))

			if len(key) > 0 {
				additionalInfoMap[key] = value
			}

			i = nextIndex
		}
	}

	return additionalInfoMap
}

func (s *mtStatement) isInformationToAccountOwnerStructured() bool {
	information := strings.Join(s.InformationToAccountOwner, "")

	for key, subFieldCount := range mtInformationToAccountOwnerTagSubFieldCounts {
		if subFieldCount < 2 {
			continue
		}

		tag := "/" + key + "/"

		if key == MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_REMITTANCE {
			if strings.Contains(information, tag+MT_REMITTANCE_CODE_UNSTRUCTURED+"/") || strings.Contains(information, tag+MT_REMITTANCE_CODE_STRUCTURED+"/") {
				return true
			}
		} else if strings.Contains(information, tag) {
			return true
		}
	}

	return false
}
