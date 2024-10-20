package gnucash

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var gnucashTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         true,
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:                 true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:             true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:              true,
}

// gnucashTransactionDataTable defines the structure of gnucash transaction data table
type gnucashTransactionDataTable struct {
	allData    []*gnucashTransactionData
	accountMap map[string]*gnucashAccountData
}

// gnucashTransactionDataRow defines the structure of gnucash transaction data row
type gnucashTransactionDataRow struct {
	dataTable  *gnucashTransactionDataTable
	data       *gnucashTransactionData
	finalItems map[datatable.TransactionDataTableColumn]string
	isValid    bool
}

// gnucashTransactionDataRowIterator defines the structure of gnucash transaction data row iterator
type gnucashTransactionDataRowIterator struct {
	dataTable    *gnucashTransactionDataTable
	currentIndex int
}

// HasColumn returns whether the transaction data table has specified column
func (t *gnucashTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := gnucashTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *gnucashTransactionDataTable) TransactionRowCount() int {
	return len(t.allData)
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *gnucashTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &gnucashTransactionDataRowIterator{
		dataTable:    t,
		currentIndex: -1,
	}
}

// IsValid returns whether this row is valid data for importing
func (r *gnucashTransactionDataRow) IsValid() bool {
	return r.isValid
}

// GetData returns the data in the specified column type
func (r *gnucashTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := gnucashTransactionSupportedColumns[column]

	if exists {
		return r.finalItems[column]
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *gnucashTransactionDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allData)
}

// Next returns the next imported data row
func (t *gnucashTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	if t.currentIndex+1 >= len(t.dataTable.allData) {
		return nil, nil
	}

	t.currentIndex++

	data := t.dataTable.allData[t.currentIndex]
	rowItems, isValid, err := t.parseTransaction(ctx, user, data)

	if err != nil {
		return nil, err
	}

	return &gnucashTransactionDataRow{
		dataTable:  t.dataTable,
		data:       data,
		finalItems: rowItems,
		isValid:    isValid,
	}, nil
}

func (t *gnucashTransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, gnucashTransaction *gnucashTransactionData) (map[datatable.TransactionDataTableColumn]string, bool, error) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(gnucashTransactionSupportedColumns))

	if gnucashTransaction.PostedDate == "" {
		return nil, false, errs.ErrMissingTransactionTime
	}

	dateTime, err := utils.ParseFromLongDateTimeWithTimezone2(gnucashTransaction.PostedDate)

	if err != nil {
		return nil, false, errs.ErrTransactionTimeInvalid
	}

	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(dateTime.Unix(), dateTime.Location())
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(dateTime.Location())

	if len(gnucashTransaction.Splits) == 2 {
		splitData1 := gnucashTransaction.Splits[0]
		splitData2 := gnucashTransaction.Splits[1]

		account1 := t.dataTable.accountMap[splitData1.Account]
		account2 := t.dataTable.accountMap[splitData2.Account]

		if account1 == nil || account2 == nil {
			return nil, false, errs.ErrMissingAccountData
		}

		amount1, err := t.parseAmount(splitData1.Quantity)

		if err != nil {
			return nil, false, err
		}

		amount2, err := t.parseAmount(splitData2.Quantity)

		if err != nil {
			return nil, false, err
		}

		if ((account1.AccountType == gnucashEquityAccountType || account1.AccountType == gnucashIncomeAccountType) && gnucashAssetOrLiabilityAccountTypes[account2.AccountType]) ||
			((account2.AccountType == gnucashEquityAccountType || account2.AccountType == gnucashIncomeAccountType) && gnucashAssetOrLiabilityAccountTypes[account1.AccountType]) { // income
			fromAccount := account1
			toAccount := account2
			toAmount := amount2

			if (account2.AccountType == gnucashEquityAccountType || account2.AccountType == gnucashIncomeAccountType) && gnucashAssetOrLiabilityAccountTypes[account1.AccountType] {
				fromAccount = account2
				toAccount = account1
				toAmount = amount1
			}

			if t.hasSpecifiedSlotKeyValue(fromAccount.Slots, gnucashSlotEquityType, gnucashSlotEquityTypeOpeningBalance) {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_MODIFY_BALANCE))
			} else {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_INCOME))
			}

			data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = t.getCategoryName(fromAccount)
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = fromAccount.Name
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = toAccount.Name

			if toAccount.Commodity != nil && toAccount.Commodity.Space == gnucashCommodityCurrencySpace {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = toAccount.Commodity.Id
			}

			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = toAmount
		} else if (account1.AccountType == gnucashExpenseAccountType && gnucashAssetOrLiabilityAccountTypes[account2.AccountType]) ||
			(account2.AccountType == gnucashExpenseAccountType && gnucashAssetOrLiabilityAccountTypes[account1.AccountType]) { // expense
			fromAccount := account1
			fromAmount := amount1
			toAccount := account2

			if account1.AccountType == gnucashExpenseAccountType && gnucashAssetOrLiabilityAccountTypes[account2.AccountType] {
				fromAccount = account2
				fromAmount = amount2
				toAccount = account1
			}

			if len(fromAmount) > 0 && fromAmount[0] == '-' {
				amount, err := utils.ParseAmount(fromAmount)

				if err != nil {
					return nil, false, errs.ErrAmountInvalid
				}

				fromAmount = utils.FormatAmount(-amount)
			}

			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE))
			data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = t.getCategoryName(toAccount)
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = toAccount.Name
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = fromAccount.Name

			if fromAccount.Commodity != nil && fromAccount.Commodity.Space == gnucashCommodityCurrencySpace {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = fromAccount.Commodity.Id
			}

			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = fromAmount
		} else if gnucashAssetOrLiabilityAccountTypes[account1.AccountType] && gnucashAssetOrLiabilityAccountTypes[account2.AccountType] {
			var fromAccount, toAccount *gnucashAccountData
			var fromAmount, toAmount string

			if len(amount1) > 0 && amount1[0] == '-' {
				fromAccount = account1
				fromAmount = amount1[1:]
				toAccount = account2
				toAmount = amount2
			} else if len(amount2) > 0 && amount2[0] == '-' {
				fromAccount = account2
				fromAmount = amount2[1:]
				toAccount = account1
				toAmount = amount1
			} else {
				log.Errorf(ctx, "[gnucash_transaction_table.parseTransaction] cannot parse transfer transaction \"id:%s\", because unexcepted account amounts \"%s\" and \"%s\"", gnucashTransaction.Id, amount1, amount2)
				return nil, false, errs.ErrInvalidGnuCashFile
			}

			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER))
			data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = ""
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = fromAccount.Name

			if fromAccount.Commodity != nil && fromAccount.Commodity.Space == gnucashCommodityCurrencySpace {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = fromAccount.Commodity.Id
			}

			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = fromAmount
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = toAccount.Name

			if toAccount.Commodity != nil && toAccount.Commodity.Space == gnucashCommodityCurrencySpace {
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = toAccount.Commodity.Id
			}

			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = toAmount
		} else {
			log.Errorf(ctx, "[gnucash_transaction_table.parseTransaction] cannot parse transaction \"id:%s\", because unexcepted account types \"%s\" and \"%s\"", gnucashTransaction.Id, account1.AccountType, account2.AccountType)
			return nil, false, errs.ErrThereAreNotSupportedTransactionType
		}
	} else if len(gnucashTransaction.Splits) == 1 {
		splitData := gnucashTransaction.Splits[0]
		account := t.dataTable.accountMap[splitData.Account]

		if account == nil {
			return nil, false, errs.ErrMissingAccountData
		}

		amount, err := t.parseAmount(splitData.Quantity)

		if err != nil {
			return nil, false, err
		}

		amountNum, err := utils.ParseAmount(amount)

		if err != nil {
			return nil, false, err
		}

		if amountNum == 0 {
			log.Warnf(ctx, "[gnucash_transaction_table.parseTransaction] skip parsing transaction \"id:%s\" with zero amount", gnucashTransaction.Id)
			return nil, false, nil
		}

		log.Errorf(ctx, "[gnucash_transaction_table.parseTransaction] cannot parse transaction \"id:%s\", because split count is %d", gnucashTransaction.Id, len(gnucashTransaction.Splits))
		return nil, false, errs.ErrThereAreNotSupportedTransactionType
	} else if len(gnucashTransaction.Splits) < 1 {
		log.Errorf(ctx, "[gnucash_transaction_table.parseTransaction] cannot parse transaction \"id:%s\", because split count is %d", gnucashTransaction.Id, len(gnucashTransaction.Splits))
		return nil, false, errs.ErrInvalidGnuCashFile
	} else {
		log.Errorf(ctx, "[gnucash_transaction_table.parseTransaction] cannot parse split transaction \"id:%s\", because split count is %d", gnucashTransaction.Id, len(gnucashTransaction.Splits))
		return nil, false, errs.ErrNotSupportedSplitTransactions
	}

	data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = gnucashTransaction.Description

	return data, true, nil
}

func (t *gnucashTransactionDataRowIterator) parseAmount(quantity string) (string, error) {
	items := strings.Split(quantity, "/")

	if len(items) != 2 {
		return "", errs.ErrAmountInvalid
	}

	value, err := utils.StringToInt64(items[0])

	if err != nil {
		return "", errs.ErrAmountInvalid
	}

	if items[1] == "100" {
		return utils.FormatAmount(value), nil
	}

	factor, err := utils.StringToInt64(items[1])

	if err != nil {
		return "", errs.ErrAmountInvalid
	}

	value = value * 100 / factor

	return utils.FormatAmount(value), nil
}

func (t *gnucashTransactionDataRowIterator) getCategoryName(accountData *gnucashAccountData) string {
	if accountData == nil || accountData.ParentId == "" {
		return ""
	}

	parentAccount := t.dataTable.accountMap[accountData.ParentId]

	if parentAccount == nil || parentAccount.AccountType == gnucashRootAccountType {
		return ""
	}

	return parentAccount.Name
}

func (t *gnucashTransactionDataRowIterator) hasSpecifiedSlotKeyValue(slots []*gnucashSlotData, key string, value string) bool {
	for i := 0; i < len(slots); i++ {
		if slots[i].Key == key && slots[i].Value == value {
			return true
		}
	}

	return false
}

func createNewGnuCashTransactionDataTable(database *gnucashDatabase) (*gnucashTransactionDataTable, error) {
	if database == nil || len(database.Books) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	allData := make([]*gnucashTransactionData, 0)
	accountMap := make(map[string]*gnucashAccountData)

	for i := 0; i < len(database.Books); i++ {
		book := database.Books[i]
		allData = append(allData, book.Transactions...)

		for j := 0; j < len(book.Accounts); j++ {
			account := book.Accounts[j]
			accountMap[account.Id] = account
		}
	}

	return &gnucashTransactionDataTable{
		allData:    allData,
		accountMap: accountMap,
	}, nil
}
