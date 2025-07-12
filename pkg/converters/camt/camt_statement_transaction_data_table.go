package camt

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var camtTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE: true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:     true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

// camtStatementTransactionDataTable defines the structure of camt statement transaction data table
type camtStatementTransactionDataTable struct {
	allStatements []*camtStatement
}

// camtStatementTransactionDataRow defines the structure of camt statement transaction data row
type camtStatementTransactionDataRow struct {
	dataTable          *camtStatementTransactionDataTable
	account            *camtAccount
	entry              *camtEntry
	transactionDetails *camtTransactionDetails
	finalItems         map[datatable.TransactionDataTableColumn]string
}

// camtStatementTransactionDataRowIterator defines the structure of camt statement transaction data row iterator
type camtStatementTransactionDataRowIterator struct {
	dataTable                      *camtStatementTransactionDataTable
	currentStatementIndex          int
	currentEntryIndex              int
	currentTransactionDetailsIndex int
}

// HasColumn returns whether the transaction data table has specified column
func (t *camtStatementTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := camtTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *camtStatementTransactionDataTable) TransactionRowCount() int {
	totalDataRowCount := 0

	for i := 0; i < len(t.allStatements); i++ {
		statement := t.allStatements[i]

		for j := 0; j < len(statement.Entries); j++ {
			entry := statement.Entries[j]

			if entry.EntryDetails != nil {
				totalDataRowCount += len(entry.EntryDetails.TransactionDetails)
			} else {
				totalDataRowCount++
			}
		}
	}

	return totalDataRowCount
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *camtStatementTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &camtStatementTransactionDataRowIterator{
		dataTable:                      t,
		currentStatementIndex:          0,
		currentEntryIndex:              0,
		currentTransactionDetailsIndex: -1,
	}
}

// IsValid returns whether this row is valid data for importing
func (r *camtStatementTransactionDataRow) IsValid() bool {
	return true
}

// GetData returns the data in the specified column type
func (r *camtStatementTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := camtTransactionSupportedColumns[column]

	if exists {
		return r.finalItems[column]
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *camtStatementTransactionDataRowIterator) HasNext() bool {
	allStatements := t.dataTable.allStatements

	if t.currentStatementIndex >= len(allStatements) {
		return false
	}

	currentStatement := allStatements[t.currentStatementIndex]

	if t.currentEntryIndex+1 < len(currentStatement.Entries) {
		return true
	} else if t.currentEntryIndex < len(currentStatement.Entries) {
		currencyEntry := currentStatement.Entries[t.currentEntryIndex]

		if currencyEntry.EntryDetails != nil {
			if t.currentTransactionDetailsIndex+1 < len(currencyEntry.EntryDetails.TransactionDetails) {
				return true
			}
		} else {
			if t.currentTransactionDetailsIndex < 0 {
				return true
			}
		}
	}

	for i := t.currentStatementIndex + 1; i < len(allStatements); i++ {
		statement := allStatements[i]

		if len(statement.Entries) < 1 {
			continue
		}

		return true
	}

	return false
}

// Next returns the next transaction data row
func (t *camtStatementTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	allStatements := t.dataTable.allStatements

	for i := t.currentStatementIndex; i < len(allStatements); i++ {
		foundNextRow := false
		statement := allStatements[i]

		for j := t.currentEntryIndex; j < len(statement.Entries); j++ {
			if statement.Entries[j].EntryDetails != nil {
				if t.currentTransactionDetailsIndex+1 < len(statement.Entries[j].EntryDetails.TransactionDetails) {
					t.currentTransactionDetailsIndex++
					foundNextRow = true
					break
				}
			} else {
				if t.currentTransactionDetailsIndex < 0 {
					t.currentTransactionDetailsIndex++
					foundNextRow = true
					break
				}
			}

			t.currentEntryIndex++
			t.currentTransactionDetailsIndex = -1
		}

		if foundNextRow {
			break
		}

		t.currentStatementIndex++
		t.currentEntryIndex = 0
		t.currentTransactionDetailsIndex = -1
	}

	if t.currentStatementIndex >= len(allStatements) {
		return nil, nil
	}

	currentStatement := allStatements[t.currentStatementIndex]

	if t.currentEntryIndex >= len(currentStatement.Entries) {
		return nil, nil
	}

	account := currentStatement.Account
	entry := currentStatement.Entries[t.currentEntryIndex]
	var transactionDetails *camtTransactionDetails

	if entry.EntryDetails != nil {
		if t.currentTransactionDetailsIndex >= len(entry.EntryDetails.TransactionDetails) {
			return nil, nil
		} else {
			transactionDetails = entry.EntryDetails.TransactionDetails[t.currentTransactionDetailsIndex]
		}
	} else {
		if t.currentTransactionDetailsIndex >= 1 {
			return nil, nil
		}
	}

	rowItems, err := t.parseTransaction(ctx, user, account, entry, transactionDetails)

	if err != nil {
		log.Errorf(ctx, "[camt_statement_transaction_data_table.Next] cannot parsing transaction in entry#%d-transaction_detail#%d (statement#%d), because %s", t.currentEntryIndex, t.currentTransactionDetailsIndex, t.currentStatementIndex, err.Error())
		return nil, err
	}

	return &camtStatementTransactionDataRow{
		dataTable:          t.dataTable,
		account:            account,
		entry:              entry,
		transactionDetails: transactionDetails,
		finalItems:         rowItems,
	}, nil
}

func (t *camtStatementTransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, account *camtAccount, entry *camtEntry, transactionDetails *camtTransactionDetails) (map[datatable.TransactionDataTableColumn]string, error) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(camtTransactionSupportedColumns))

	if account == nil {
		return nil, errs.ErrMissingAccountData
	}

	if entry.BookingDate != nil && entry.BookingDate.DateTime != "" {
		dateTime, err := utils.ParseFromLongDateTimeWithTimezoneRFC3339Format(entry.BookingDate.DateTime)

		if err != nil {
			return nil, errs.ErrTransactionTimeInvalid
		}

		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(dateTime.Unix(), dateTime.Location())
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(dateTime.Location())
	} else if entry.BookingDate != nil && entry.BookingDate.Date != "" {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = fmt.Sprintf("%s 00:00:00", entry.BookingDate.Date)
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = datatable.TRANSACTION_DATA_TABLE_TIMEZONE_NOT_AVAILABLE
	} else {
		return nil, errs.ErrMissingTransactionTime
	}

	if account.IBAN != "" {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = account.IBAN
	} else if account.OtherIdentification != "" {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = account.OtherIdentification
	}

	if transactionDetails != nil && transactionDetails.AmountDetails != nil && transactionDetails.AmountDetails.TransactionAmount != nil && transactionDetails.AmountDetails.TransactionAmount.Currency != "" {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = transactionDetails.AmountDetails.TransactionAmount.Currency
	} else if entry.Amount != nil && entry.Amount.Currency != "" {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = entry.Amount.Currency
	} else if account.Currency != "" {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = account.Currency
	} else {
		return nil, errs.ErrAccountCurrencyInvalid
	}

	amountValue := ""

	if entry.EntryDetails != nil && len(entry.EntryDetails.TransactionDetails) > 1 && transactionDetails != nil { // when there are multiple transaction details in one entry, only use the amount in the transaction details
		if transactionDetails.AmountDetails != nil && transactionDetails.AmountDetails.InstructedAmount != nil && transactionDetails.AmountDetails.InstructedAmount.Value != "" {
			amountValue = transactionDetails.AmountDetails.InstructedAmount.Value
		} else if transactionDetails.AmountDetails != nil && transactionDetails.AmountDetails.TransactionAmount != nil && transactionDetails.AmountDetails.TransactionAmount.Value != "" {
			amountValue = transactionDetails.AmountDetails.TransactionAmount.Value
		} else {
			return nil, errs.ErrAmountInvalid
		}
	} else if entry.Amount != nil && entry.Amount.Value != "" {
		amountValue = entry.Amount.Value
	}

	if amountValue == "" {
		return nil, errs.ErrAmountInvalid
	}

	amount, err := utils.ParseAmount(amountValue)

	if err != nil {
		log.Errorf(ctx, "[camt_statement_transaction_data_table.parseTransaction] cannot parsing transaction amount \"%s\", because %s", amountValue, err.Error())
		return nil, errs.ErrAmountInvalid
	}

	data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)

	if entry.CreditDebitIndicator == CAMT_INDICATOR_CREDIT {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_INCOME))
	} else if entry.CreditDebitIndicator == CAMT_INDICATOR_DEBIT {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE))
	} else {
		return nil, errs.ErrTransactionTypeInvalid
	}

	if transactionDetails != nil && transactionDetails.AdditionalTransactionInformation != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = transactionDetails.AdditionalTransactionInformation
	} else if transactionDetails != nil && transactionDetails.RemittanceInformation != nil && len(transactionDetails.RemittanceInformation.Unstructured) > 0 {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = strings.Join(transactionDetails.RemittanceInformation.Unstructured, "\n")
	} else if entry.AdditionalEntryInformation != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = entry.AdditionalEntryInformation
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	return data, nil
}

func createNewCamtStatementTransactionDataTable(file *camt053File) (*camtStatementTransactionDataTable, error) {
	if file == nil || file.BankToCustomerStatement == nil || len(file.BankToCustomerStatement.Statements) == 0 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return &camtStatementTransactionDataTable{
		allStatements: file.BankToCustomerStatement.Statements,
	}, nil
}
