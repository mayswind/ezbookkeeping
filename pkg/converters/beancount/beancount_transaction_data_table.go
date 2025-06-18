package beancount

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var beancountTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:             true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:              true,
}

var BEANCOUNT_TRANSACTION_TAG_SEPARATOR = "#"

// beancountTransactionDataTable defines the structure of Beancount transaction data table
type beancountTransactionDataTable struct {
	allData    []*beancountTransactionEntry
	accountMap map[string]*beancountAccount
}

// beancountTransactionDataRow defines the structure of Beancount transaction data row
type beancountTransactionDataRow struct {
	dataTable  *beancountTransactionDataTable
	data       *beancountTransactionEntry
	finalItems map[datatable.TransactionDataTableColumn]string
}

// beancountTransactionDataRowIterator defines the structure of Beancount transaction data row iterator
type beancountTransactionDataRowIterator struct {
	dataTable    *beancountTransactionDataTable
	currentIndex int
}

// HasColumn returns whether the transaction data table has specified column
func (t *beancountTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := beancountTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *beancountTransactionDataTable) TransactionRowCount() int {
	return len(t.allData)
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *beancountTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &beancountTransactionDataRowIterator{
		dataTable:    t,
		currentIndex: -1,
	}
}

// IsValid returns whether this row is valid data for importing
func (r *beancountTransactionDataRow) IsValid() bool {
	return true
}

// GetData returns the data in the specified column type
func (r *beancountTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := beancountTransactionSupportedColumns[column]

	if exists {
		return r.finalItems[column]
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *beancountTransactionDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allData)
}

// Next returns the next transaction data row
func (t *beancountTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	if t.currentIndex+1 >= len(t.dataTable.allData) {
		return nil, nil
	}

	t.currentIndex++

	data := t.dataTable.allData[t.currentIndex]
	rowItems, err := t.parseTransaction(ctx, user, data)

	if err != nil {
		return nil, err
	}

	return &beancountTransactionDataRow{
		dataTable:  t.dataTable,
		data:       data,
		finalItems: rowItems,
	}, nil
}

func (t *beancountTransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, beancountEntry *beancountTransactionEntry) (map[datatable.TransactionDataTableColumn]string, error) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(beancountTransactionSupportedColumns))

	if beancountEntry.date == "" {
		return nil, errs.ErrMissingTransactionTime
	}

	// Beancount supports the international ISO 8601 standard format for dates, with dashes or the same ordering with slashes
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = strings.ReplaceAll(beancountEntry.date, "/", "-") + " 00:00:00"

	if len(beancountEntry.postings) == 2 {
		splitData1 := beancountEntry.postings[0]
		splitData2 := beancountEntry.postings[1]

		account1 := t.dataTable.accountMap[splitData1.account]
		account2 := t.dataTable.accountMap[splitData2.account]

		if account1 == nil || account2 == nil {
			return nil, errs.ErrMissingAccountData
		}

		amount1, err := utils.ParseAmount(splitData1.amount)

		if err != nil {
			log.Errorf(ctx, "[beancount_transaction_data_table.parseTransaction] cannot parse amount \"%s\", because %s", splitData1.amount, err.Error())
			return nil, errs.ErrAmountInvalid
		}

		amount2, err := utils.ParseAmount(splitData2.amount)

		if err != nil {
			log.Errorf(ctx, "[beancount_transaction_data_table.parseTransaction] cannot parse amount \"%s\", because %s", splitData2.amount, err.Error())
			return nil, errs.ErrAmountInvalid
		}

		if ((account1.accountType == beancountEquityAccountType || account1.accountType == beancountIncomeAccountType) && (account2.accountType == beancountAssetsAccountType || account2.accountType == beancountLiabilitiesAccountType)) ||
			((account2.accountType == beancountEquityAccountType || account2.accountType == beancountIncomeAccountType) && (account1.accountType == beancountAssetsAccountType || account1.accountType == beancountLiabilitiesAccountType)) { // income
			fromAccount := account1
			toAccount := account2
			toCurrency := splitData2.commodity
			toAmount := amount2

			if (account2.accountType == beancountEquityAccountType || account2.accountType == beancountIncomeAccountType) && (account1.accountType == beancountAssetsAccountType || account1.accountType == beancountLiabilitiesAccountType) {
				fromAccount = account2
				toAccount = account1
				toCurrency = splitData1.commodity
				toAmount = amount1
			}

			if fromAccount.isOpeningBalanceEquityAccount() {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_MODIFY_BALANCE))
			} else {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_INCOME))
			}

			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = fromAccount.name
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = toAccount.name
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = toCurrency
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(toAmount)
		} else if account1.accountType == beancountExpensesAccountType && (account2.accountType == beancountAssetsAccountType || account2.accountType == beancountLiabilitiesAccountType) ||
			(account2.accountType == beancountExpensesAccountType && (account1.accountType == beancountAssetsAccountType || account1.accountType == beancountLiabilitiesAccountType)) { // expense
			fromAccount := account1
			fromCurrency := splitData1.commodity
			fromAmount := amount1
			toAccount := account2

			if account1.accountType == beancountExpensesAccountType && (account2.accountType == beancountAssetsAccountType || account2.accountType == beancountLiabilitiesAccountType) {
				fromAccount = account2
				fromCurrency = splitData2.commodity
				fromAmount = amount2
				toAccount = account1
			}

			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE))
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = toAccount.name
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = fromAccount.name
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = fromCurrency
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-fromAmount)
		} else if (account1.accountType == beancountAssetsAccountType || account1.accountType == beancountLiabilitiesAccountType) &&
			(account2.accountType == beancountAssetsAccountType || account2.accountType == beancountLiabilitiesAccountType) {
			var fromAccount, toAccount *beancountAccount
			var fromAmount, toAmount int64
			var fromCurrency, toCurrency string

			if amount1 < 0 {
				fromAccount = account1
				fromCurrency = splitData1.commodity
				fromAmount = -amount1
				toAccount = account2
				toCurrency = splitData2.commodity
				toAmount = amount2
			} else if amount2 < 0 {
				fromAccount = account2
				fromCurrency = splitData2.commodity
				fromAmount = -amount2
				toAccount = account1
				toCurrency = splitData1.commodity
				toAmount = amount1
			} else {
				log.Errorf(ctx, "[beancount_transaction_data_table.parseTransaction] cannot parse transfer transaction, because unexcepted account amounts \"%d\" and \"%d\"", amount1, amount2)
				return nil, errs.ErrInvalidBeancountFile
			}

			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER))
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = fromAccount.name
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = fromCurrency
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(fromAmount)
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = toAccount.name
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = toCurrency
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(toAmount)
		} else {
			log.Errorf(ctx, "[beancount_transaction_data_table.parseTransaction] cannot parse transaction, because unexcepted account types \"%d\" and \"%d\"", account1.accountType, account2.accountType)
			return nil, errs.ErrThereAreNotSupportedTransactionType
		}
	} else if len(beancountEntry.postings) <= 1 {
		log.Errorf(ctx, "[beancount_transaction_data_table.parseTransaction] cannot parse transaction, because postings count is %d", len(beancountEntry.postings))
		return nil, errs.ErrInvalidBeancountFile
	} else {
		log.Errorf(ctx, "[beancount_transaction_data_table.parseTransaction] cannot parse split transaction, because postings count is %d", len(beancountEntry.postings))
		return nil, errs.ErrNotSupportedSplitTransactions
	}

	data[datatable.TRANSACTION_DATA_TABLE_TAGS] = strings.Join(beancountEntry.tags, BEANCOUNT_TRANSACTION_TAG_SEPARATOR)
	data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = beancountEntry.narration

	return data, nil
}

func createNewBeancountTransactionDataTable(beancountData *beancountData) (*beancountTransactionDataTable, error) {
	if beancountData == nil {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return &beancountTransactionDataTable{
		allData:    beancountData.transactions,
		accountMap: beancountData.accounts,
	}, nil
}
