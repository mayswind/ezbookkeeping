package qif

import (
	"bufio"
	"bytes"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

const qifBankTransactionHeader = "!Type:Bank"
const qifCashTransactionHeader = "!Type:Cash"
const qifCreditCardTransactionHeader = "!Type:CCard"
const qifAssetAccountTransactionHeader = "!Type:Oth A"
const qifLiabilityAccountTransactionHeader = "!Type:Oth L"
const qifMemorizedTransactionHeader = "!Type:Memorized"
const qifMemorisedTransactionHeader = "!Type:Memorised"
const qifInvestmentTransactionHeader = "!Type:Invst"
const qifAccountHeader = "!Account"
const qifCategoryHeader = "!Type:Cat"
const qifClassHeader = "!Type:Class"
const qifTypeHeaderPrefix = "!Type:"

const qifEntryStartRune = '!'
const qifEntryEnd = '^'

// qifDataReader defines the structure of quicken interchange format (qif) data reader
type qifDataReader struct {
	allLines []string
}

// read returns the imported qif data
// Reference: https://www.w3.org/2000/10/swap/pim/qif-doc/QIF-doc.htm
func (r *qifDataReader) read(ctx core.Context) (*qifData, error) {
	if len(r.allLines) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	data := &qifData{}
	var currentEntryHeader string
	var currentEntryData []string
	var currentAccount *qifAccountData

	for i := 0; i < len(r.allLines); i++ {
		line := r.allLines[i]

		if len(line) < 1 {
			continue
		}

		if line[0] == qifEntryStartRune {
			if len(currentEntryData) > 0 {
				log.Errorf(ctx, "[qif_data_reader.read] read new entry header \"%s\" after unclosed entry", line)
				return nil, errs.ErrInvalidQIFFile
			}

			line = strings.TrimRight(line, " ")

			if line == qifBankTransactionHeader ||
				line == qifCashTransactionHeader ||
				line == qifCreditCardTransactionHeader ||
				line == qifAssetAccountTransactionHeader ||
				line == qifLiabilityAccountTransactionHeader ||
				line == qifMemorizedTransactionHeader ||
				line == qifMemorisedTransactionHeader ||
				line == qifInvestmentTransactionHeader ||
				line == qifAccountHeader ||
				line == qifCategoryHeader ||
				line == qifClassHeader {
				currentEntryHeader = line
			} else if strings.Index(line, qifTypeHeaderPrefix) == 0 {
				currentEntryHeader = line
				log.Warnf(ctx, "[qif_data_reader.read] read unsupported entry header line \"%s\" and skip the following entries", line)
			} else {
				log.Warnf(ctx, "[qif_data_reader.read] read unsupported entry header line \"%s\" and skip this line", line)
			}
		} else if line[0] == qifEntryEnd {
			entryData := currentEntryData
			currentEntryData = nil

			if currentEntryHeader == qifBankTransactionHeader ||
				currentEntryHeader == qifCashTransactionHeader ||
				currentEntryHeader == qifCreditCardTransactionHeader ||
				currentEntryHeader == qifAssetAccountTransactionHeader ||
				currentEntryHeader == qifLiabilityAccountTransactionHeader {
				transactionData, err := r.parseTransaction(ctx, entryData, false)

				if err != nil {
					return nil, err
				}

				if transactionData == nil {
					continue
				}

				transactionData.Account = currentAccount

				if currentEntryHeader == qifBankTransactionHeader {
					data.BankAccountTransactions = append(data.BankAccountTransactions, transactionData)
				} else if currentEntryHeader == qifCashTransactionHeader {
					data.CashAccountTransactions = append(data.CashAccountTransactions, transactionData)
				} else if currentEntryHeader == qifCreditCardTransactionHeader {
					data.CreditCardAccountTransactions = append(data.CreditCardAccountTransactions, transactionData)
				} else if currentEntryHeader == qifAssetAccountTransactionHeader {
					data.AssetAccountTransactions = append(data.AssetAccountTransactions, transactionData)
				} else if currentEntryHeader == qifLiabilityAccountTransactionHeader {
					data.LiabilityAccountTransactions = append(data.LiabilityAccountTransactions, transactionData)
				}
			} else if currentEntryHeader == qifMemorizedTransactionHeader || currentEntryHeader == qifMemorisedTransactionHeader {
				transactionData, err := r.parseMemorizedTransaction(ctx, entryData)

				if err != nil {
					return nil, err
				}

				if transactionData == nil {
					continue
				}

				transactionData.Account = currentAccount
				data.MemorizedTransactions = append(data.MemorizedTransactions, transactionData)
			} else if currentEntryHeader == qifInvestmentTransactionHeader {
				transactionData, err := r.parseInvestmentTransaction(ctx, entryData)

				if err != nil {
					return nil, err
				}

				if transactionData == nil {
					continue
				}

				transactionData.Account = currentAccount
				data.InvestmentAccountTransactions = append(data.InvestmentAccountTransactions, transactionData)
			} else if currentEntryHeader == qifAccountHeader {
				accountData, err := r.parseAccount(ctx, entryData)

				if err != nil {
					return nil, err
				}

				if accountData == nil {
					continue
				}

				currentAccount = accountData
				data.Accounts = append(data.Accounts, accountData)
			} else if currentEntryHeader == qifCategoryHeader {
				categoryData, err := r.parseCategory(ctx, entryData)

				if err != nil {
					return nil, err
				}

				if categoryData == nil {
					continue
				}

				data.Categories = append(data.Categories, categoryData)
			} else if currentEntryHeader == qifClassHeader {
				classData, err := r.parseClass(ctx, entryData)

				if err != nil {
					return nil, err
				}

				if classData == nil {
					continue
				}

				data.Classes = append(data.Classes, classData)
			} else {
				log.Warnf(ctx, "[qif_data_reader.read] read unsupported entry header \"%s\" and skip this entry", currentEntryHeader)
			}
		} else if currentEntryHeader != "" {
			currentEntryData = append(currentEntryData, line)
		} else {
			log.Warnf(ctx, "[qif_data_reader.read] read unsupported line \"%s\" and skip this line", line)
		}
	}

	return data, nil
}

func (r *qifDataReader) parseTransaction(ctx core.Context, data []string, ignoreUnknown bool) (*qifTransactionData, error) {
	if len(data) < 1 {
		return nil, nil
	}

	transactionData := &qifTransactionData{}

	for i := 0; i < len(data); i++ {
		line := data[i]

		if len(line) < 1 {
			continue
		}

		if line[0] == 'D' {
			transactionData.Date = line[1:]
		} else if line[0] == 'T' {
			transactionData.Amount = line[1:]
		} else if line[0] == 'C' {
			transactionData.ClearedStatus = r.parseClearedStatus(ctx, line[1:])
		} else if line[0] == 'N' {
			transactionData.Num = line[1:]
		} else if line[0] == 'P' {
			transactionData.Payee = line[1:]
		} else if line[0] == 'M' {
			transactionData.Memo = line[1:]
		} else if line[0] == 'A' {
			transactionData.Addresses = append(transactionData.Addresses, line[1:])
		} else if line[0] == 'L' {
			transactionData.Category = line[1:]
		} else if line[0] == 'S' {
			transactionData.SubTransactionCategory = append(transactionData.SubTransactionCategory, line[1:])
		} else if line[0] == 'E' {
			transactionData.SubTransactionMemo = append(transactionData.SubTransactionMemo, line[1:])
		} else if line[0] == '$' {
			transactionData.SubTransactionAmount = append(transactionData.SubTransactionAmount, line[1:])
		} else {
			if !ignoreUnknown {
				log.Warnf(ctx, "[qif_data_reader.parseTransaction] read unsupported line \"%s\" and skip this line", line)
				continue
			}
		}
	}

	return transactionData, nil
}

func (r *qifDataReader) parseMemorizedTransaction(ctx core.Context, data []string) (*qifMemorizedTransactionData, error) {
	if len(data) < 1 {
		return nil, nil
	}

	baseTransactionData, err := r.parseTransaction(ctx, data, true)

	if err != nil {
		return nil, err
	}

	transactionData := &qifMemorizedTransactionData{
		qifTransactionData: *baseTransactionData,
		Amortization:       qifMemorizedTransactionAmortizationData{},
	}

	for i := 0; i < len(data); i++ {
		line := data[i]

		if len(line) < 1 {
			continue
		}

		// these lines has been already processed in parseTransaction
		if line[0] == 'D' || line[0] == 'T' || line[0] == 'C' || line[0] == 'N' ||
			line[0] == 'P' || line[0] == 'M' || line[0] == 'A' || line[0] == 'L' ||
			line[0] == 'S' || line[0] == 'E' || line[0] == '$' {
			continue
		}

		if line[0] == 'K' {
			if line == string(qifCheckTransactionType) {
				transactionData.TransactionType = qifCheckTransactionType
			} else if line == string(qifDepositTransactionType) {
				transactionData.TransactionType = qifDepositTransactionType
			} else if line == string(qifPaymentTransactionType) {
				transactionData.TransactionType = qifPaymentTransactionType
			} else if line == string(qifInvestmentTransactionType) {
				transactionData.TransactionType = qifInvestmentTransactionType
			} else if line == string(qifElectronicPayeeTransactionType) {
				transactionData.TransactionType = qifElectronicPayeeTransactionType
			} else {
				log.Warnf(ctx, "[qif_data_reader.parseMemorizedTransaction] read unsupported transaction type \"%s\" and skip this line", line)
				continue
			}
		} else if line[0] == '1' {
			transactionData.Amortization.FirstPaymentDate = line[1:]
		} else if line[0] == '2' {
			transactionData.Amortization.TotalYearsForLoan = line[1:]
		} else if line[0] == '3' {
			transactionData.Amortization.NumberOfPayments = line[1:]
		} else if line[0] == '4' {
			transactionData.Amortization.NumberOfPeriodsPerYear = line[1:]
		} else if line[0] == '5' {
			transactionData.Amortization.InterestRate = line[1:]
		} else if line[0] == '6' {
			transactionData.Amortization.CurrentLoanBalance = line[1:]
		} else if line[0] == '7' {
			transactionData.Amortization.OriginalLoanAmount = line[1:]
		} else {
			log.Warnf(ctx, "[qif_data_reader.parseMemorizedTransaction] read unsupported line \"%s\" and skip this line", line)
			continue
		}
	}

	return transactionData, nil
}

func (r *qifDataReader) parseInvestmentTransaction(ctx core.Context, data []string) (*qifInvestmentTransactionData, error) {
	if len(data) < 1 {
		return nil, nil
	}

	transactionData := &qifInvestmentTransactionData{}

	for i := 0; i < len(data); i++ {
		line := data[i]

		if len(line) < 1 {
			continue
		}

		if line[0] == 'D' {
			transactionData.Date = line[1:]
		} else if line[0] == 'N' {
			transactionData.Action = line[1:]
		} else if line[0] == 'Y' {
			transactionData.Security = line[1:]
		} else if line[0] == 'I' {
			transactionData.Price = line[1:]
		} else if line[0] == 'Q' {
			transactionData.Quantity = line[1:]
		} else if line[0] == 'T' {
			transactionData.Amount = line[1:]
		} else if line[0] == 'C' {
			transactionData.ClearedStatus = r.parseClearedStatus(ctx, line[1:])
		} else if line[0] == 'P' {
			transactionData.Text = line[1:]
		} else if line[0] == 'M' {
			transactionData.Memo = line[1:]
		} else if line[0] == 'O' {
			transactionData.Commission = line[1:]
		} else if line[0] == 'L' {
			transactionData.AccountForTransfer = line[1:]
		} else if line[0] == '$' {
			transactionData.AmountTransferred = line[1:]
		} else {
			log.Warnf(ctx, "[qif_data_reader.parseInvestmentTransaction] read unsupported line \"%s\" and skip this line", line)
			continue
		}
	}

	return transactionData, nil
}

func (r *qifDataReader) parseAccount(ctx core.Context, data []string) (*qifAccountData, error) {
	if len(data) < 1 {
		return nil, nil
	}

	accountData := &qifAccountData{}

	for i := 0; i < len(data); i++ {
		line := data[i]

		if len(line) < 1 {
			continue
		}

		if line[0] == 'N' {
			accountData.Name = line[1:]
		} else if line[0] == 'T' {
			accountData.AccountType = line[1:]
		} else if line[0] == 'D' {
			accountData.Description = line[1:]
		} else if line[0] == 'L' {
			accountData.CreditLimit = line[1:]
		} else if line[0] == '/' {
			accountData.StatementBalanceDate = line[1:]
		} else if line[0] == '$' {
			accountData.StatementBalanceAmount = line[1:]
		} else {
			log.Warnf(ctx, "[qif_data_reader.parseAccount] read unsupported line \"%s\" and skip this line", line)
			continue
		}
	}

	return accountData, nil
}

func (r *qifDataReader) parseCategory(ctx core.Context, data []string) (*qifCategoryData, error) {
	if len(data) < 1 {
		return nil, nil
	}

	categoryData := &qifCategoryData{}

	for i := 0; i < len(data); i++ {
		line := data[i]

		if len(line) < 1 {
			continue
		}

		if line[0] == 'N' {
			categoryData.Name = line[1:]
		} else if line[0] == 'D' {
			categoryData.Description = line[1:]
		} else if line[0] == 'T' {
			categoryData.TaxRelated = true
		} else if line[0] == 'I' {
			categoryData.CategoryType = qifIncomeTransaction
		} else if line[0] == 'E' {
			categoryData.CategoryType = qifExpenseTransaction
		} else if line[0] == 'B' {
			categoryData.BudgetAmount = line[1:]
		} else if line[0] == 'R' {
			categoryData.TaxScheduleInformation = line[1:]
		} else {
			log.Warnf(ctx, "[qif_data_reader.parseCategory] read unsupported line \"%s\" and skip this line", line)
			continue
		}
	}

	if categoryData.CategoryType == "" {
		categoryData.CategoryType = qifExpenseTransaction
	}

	return categoryData, nil
}

func (r *qifDataReader) parseClass(ctx core.Context, data []string) (*qifClassData, error) {
	if len(data) < 1 {
		return nil, nil
	}

	classData := &qifClassData{}

	for i := 0; i < len(data); i++ {
		line := data[i]

		if len(line) < 1 {
			continue
		}

		if line[0] == 'N' {
			classData.Name = line[1:]
		} else if line[0] == 'D' {
			classData.Description = line[1:]
		} else {
			log.Warnf(ctx, "[qif_data_reader.parseClass] read unsupported line \"%s\" and skip this line", line)
			continue
		}
	}

	return classData, nil
}

func (r *qifDataReader) parseClearedStatus(ctx core.Context, value string) qifTransactionClearedStatus {
	if value == "" {
		return qifClearedStatusUnreconciled
	} else if value == "*" || strings.ToUpper(value) == "C" {
		return qifClearedStatusCleared
	} else if strings.ToUpper(value) == "R" || strings.ToUpper(value) == "X" {
		return qifClearedStatusReconciled
	} else {
		log.Warnf(ctx, "[qif_data_reader.parseClearedStatus] read unsupported transaction cleared status \"%s\" and skip this value", value)
		return qifClearedStatusUnreconciled
	}
}

func createNewQifDataReader(data []byte) *qifDataReader {
	fallback := unicode.UTF8.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(fallback))
	scanner := bufio.NewScanner(reader)
	allLines := make([]string, 0)

	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	return &qifDataReader{
		allLines: allLines,
	}
}
