package iif

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

const iifAccountNameColumnName = "NAME"
const iifAccountTypeColumnName = "ACCNTTYPE"

const iifAccountTypeIncome = "INC"
const iifAccountTypeExpense = "EXP"

const iifTransactionTypeColumnName = "TRNSTYPE"
const iifTransactionDateColumnName = "DATE"
const iifTransactionAccountNameColumnName = "ACCNT"
const iifTransactionNameColumnName = "NAME"
const iifTransactionAmountColumnName = "AMOUNT"
const iifTransactionMemoColumnName = "MEMO"

const iifTransactionTypeBeginningBalance = "BEGINBALCHECK"

const iifTransactionCategorySeparator = ":"

var iifTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:             true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:       true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

var iifTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: utils.IntToString(int(models.TRANSACTION_TYPE_MODIFY_BALANCE)),
	models.TRANSACTION_TYPE_INCOME:         utils.IntToString(int(models.TRANSACTION_TYPE_INCOME)),
	models.TRANSACTION_TYPE_EXPENSE:        utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE)),
	models.TRANSACTION_TYPE_TRANSFER:       utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER)),
}

// iifTransactionDataTable defines the structure of intuit interchange format (iif) transaction data table
type iifTransactionDataTable struct {
	incomeAccountNames  map[string]bool
	expenseAccountNames map[string]bool
	transactionDatasets []*iifTransactionDataset
}

// iifTransactionDataRow defines the structure of intuit interchange format (iif) transaction data row
type iifTransactionDataRow struct {
	dataTable  *iifTransactionDataTable
	finalItems map[datatable.TransactionDataTableColumn]string
}

// iifTransactionDataRowIterator defines the structure of intuit interchange format (iif) transaction data row iterator
type iifTransactionDataRowIterator struct {
	dataTable             *iifTransactionDataTable
	currentDatasetIndex   int
	currentIndexInDataset int
}

// HasColumn returns whether the transaction data table has specified column
func (t *iifTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := iifTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *iifTransactionDataTable) TransactionRowCount() int {
	totalDataRowCount := 0

	for i := 0; i < len(t.transactionDatasets); i++ {
		transactions := t.transactionDatasets[i]
		totalDataRowCount += len(transactions.transactions)
	}

	return totalDataRowCount
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *iifTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &iifTransactionDataRowIterator{
		dataTable:             t,
		currentDatasetIndex:   0,
		currentIndexInDataset: -1,
	}
}

// IsValid returns whether this row is valid data for importing
func (r *iifTransactionDataRow) IsValid() bool {
	return true
}

// GetData returns the data in the specified column type
func (r *iifTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := iifTransactionSupportedColumns[column]

	if exists {
		return r.finalItems[column]
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *iifTransactionDataRowIterator) HasNext() bool {
	allDatasets := t.dataTable.transactionDatasets

	if t.currentDatasetIndex >= len(allDatasets) {
		return false
	}

	currentDataset := allDatasets[t.currentDatasetIndex]

	if t.currentIndexInDataset+1 < len(currentDataset.transactions) {
		return true
	}

	for i := t.currentDatasetIndex + 1; i < len(allDatasets); i++ {
		dataset := allDatasets[i]

		if len(dataset.transactions) < 1 {
			continue
		}

		return true
	}

	return false
}

// Next returns the next imported data row
func (t *iifTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	allDatasets := t.dataTable.transactionDatasets
	currentIndexInDataset := t.currentIndexInDataset

	for i := t.currentDatasetIndex; i < len(allDatasets); i++ {
		dataset := allDatasets[i]

		if currentIndexInDataset+1 < len(dataset.transactions) {
			t.currentIndexInDataset++
			currentIndexInDataset = t.currentIndexInDataset
			break
		}

		t.currentDatasetIndex++
		t.currentIndexInDataset = -1
		currentIndexInDataset = -1
	}

	if t.currentDatasetIndex >= len(allDatasets) {
		return nil, nil
	}

	currentDataset := allDatasets[t.currentDatasetIndex]

	if t.currentIndexInDataset >= len(currentDataset.transactions) {
		return nil, nil
	}

	data := currentDataset.transactions[t.currentIndexInDataset]
	rowItems, err := t.parseTransaction(ctx, user, currentDataset, data)

	if err != nil {
		return nil, err
	}

	return &iifTransactionDataRow{
		dataTable:  t.dataTable,
		finalItems: rowItems,
	}, nil
}

func (t *iifTransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, dataset *iifTransactionDataset, transactionData *iifTransactionData) (map[datatable.TransactionDataTableColumn]string, error) {
	if len(transactionData.splitData) < 1 {
		return nil, errs.ErrInvalidIIFFile
	} else if len(transactionData.splitData) > 1 {
		return nil, errs.ErrNotSupportedSplitTransactions
	}

	var err error

	data := make(map[datatable.TransactionDataTableColumn]string, len(iifTransactionSupportedColumns))
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME], err = t.parseTransactionTime(dataset, transactionData)

	if err != nil {
		return nil, err
	}

	transactionType, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionTypeColumnName)
	accountName1, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionAccountNameColumnName)
	accountName2, _ := dataset.getSplitDataItemValue(transactionData.splitData[0], iifTransactionAccountNameColumnName)
	amount1, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionAmountColumnName)
	amount2, _ := dataset.getSplitDataItemValue(transactionData.splitData[0], iifTransactionAmountColumnName)
	amountNum1, err := utils.ParseAmount(amount1)

	if err != nil {
		return nil, errs.ErrAmountInvalid
	}

	amountNum2, err := utils.ParseAmount(amount2)

	if err != nil {
		return nil, errs.ErrAmountInvalid
	}

	name, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionNameColumnName)

	if transactionType == iifTransactionTypeBeginningBalance { // balance modification
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = accountName1
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amountNum1)
	} else if t.dataTable.incomeAccountNames[accountName1] || t.dataTable.incomeAccountNames[accountName2] { // income
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
		categoryName := ""
		accountName := ""
		amountNum := int64(0)

		if t.dataTable.incomeAccountNames[accountName1] && !t.dataTable.incomeAccountNames[accountName2] {
			categoryName = accountName1
			accountName = accountName2
			amountNum = amountNum2
		} else if t.dataTable.incomeAccountNames[accountName2] && !t.dataTable.incomeAccountNames[accountName1] {
			categoryName = accountName2
			accountName = accountName1
			amountNum = amountNum1
		} else {
			log.Errorf(ctx, "[iif_transaction_data_table.parseTransaction] cannot parse transaction, because two accounts \"%s\" and \"%s\" are all income account", accountName1, accountName2)
			return nil, errs.ErrInvalidIIFFile
		}

		categoryNames := strings.Split(categoryName, iifTransactionCategorySeparator)

		if len(categoryNames) > 1 {
			data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = categoryNames[0]
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = categoryNames[len(categoryNames)-1]
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = categoryName
		}

		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = accountName
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amountNum)
	} else if t.dataTable.expenseAccountNames[accountName1] || t.dataTable.expenseAccountNames[accountName2] { // expense
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
		categoryName := ""
		accountName := ""
		amountNum := int64(0)

		if t.dataTable.expenseAccountNames[accountName1] && !t.dataTable.expenseAccountNames[accountName2] {
			categoryName = accountName1
			accountName = accountName2
			amountNum = amountNum2
		} else if t.dataTable.expenseAccountNames[accountName2] && !t.dataTable.expenseAccountNames[accountName1] {
			categoryName = accountName2
			accountName = accountName1
			amountNum = amountNum1
		} else {
			log.Errorf(ctx, "[iif_transaction_data_table.parseTransaction] cannot parse transaction, because two accounts \"%s\" and \"%s\" are all expense account", accountName1, accountName2)
			return nil, errs.ErrInvalidIIFFile
		}

		categoryNames := strings.Split(categoryName, iifTransactionCategorySeparator)

		if len(categoryNames) > 1 {
			data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = categoryNames[0]
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = categoryNames[len(categoryNames)-1]
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = categoryName
		}

		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = accountName
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amountNum)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = name

		if amountNum1 >= 0 {
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = accountName2
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amountNum2)
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = accountName1
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(amountNum1)
		} else if amountNum2 >= 0 {
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = accountName1
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amountNum1)
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = accountName2
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(amountNum2)
		}
	}

	memo, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionMemoColumnName)
	data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = memo

	return data, nil
}

func (t *iifTransactionDataRowIterator) parseTransactionTime(dataset *iifTransactionDataset, transactionData *iifTransactionData) (string, error) {
	date, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionDateColumnName)
	dateParts := strings.Split(date, "/")

	if len(dateParts) != 3 {
		return "", errs.ErrTransactionTimeInvalid
	}

	month := dateParts[0]
	day := dateParts[1]
	year := dateParts[2]

	if len(month) < 2 {
		month = "0" + month
	}

	if len(day) < 2 {
		day = "0" + day
	}

	return fmt.Sprintf("%s-%s-%s 00:00:00", year, month, day), nil
}

func createNewIIfTransactionDataTable(ctx core.Context, accountDatasets []*iifAccountDataset, transactionDatasets []*iifTransactionDataset) (*iifTransactionDataTable, error) {
	if len(transactionDatasets) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	incomeAccountNames, expenseAccountNames := getIncomeAndExpenseAccountNameMap(accountDatasets)

	for i := 0; i < len(transactionDatasets); i++ {
		transactionDataset := transactionDatasets[i]

		for _, requiredColumnName := range []string{
			iifTransactionDateColumnName,
			iifTransactionAccountNameColumnName,
			iifTransactionAmountColumnName,
		} {
			if _, exists := transactionDataset.transactionDataColumnIndexes[requiredColumnName]; !exists {
				return nil, errs.ErrMissingRequiredFieldInHeaderRow
			}
		}
	}

	return &iifTransactionDataTable{
		incomeAccountNames:  incomeAccountNames,
		expenseAccountNames: expenseAccountNames,
		transactionDatasets: transactionDatasets,
	}, nil
}

func getIncomeAndExpenseAccountNameMap(accountDatasets []*iifAccountDataset) (incomeAccountNames map[string]bool, expenseAccountNames map[string]bool) {
	incomeAccountNames = make(map[string]bool)
	expenseAccountNames = make(map[string]bool)

	for i := 0; i < len(accountDatasets); i++ {
		accountDataset := accountDatasets[i]
		accountNameColumnIndex, accountNameColumnExists := accountDataset.accountDataColumnIndexes[iifAccountNameColumnName]
		accountTypeColumnIndex, accountTypeColumnExists := accountDataset.accountDataColumnIndexes[iifAccountTypeColumnName]

		if !accountNameColumnExists || accountNameColumnIndex < 0 ||
			!accountTypeColumnExists || accountTypeColumnIndex < 0 {
			continue
		}

		for j := 0; j < len(accountDataset.accounts); j++ {
			items := accountDataset.accounts[j].dataItems

			if accountNameColumnIndex >= len(items) ||
				accountTypeColumnIndex >= len(items) {
				continue
			}

			accountName := items[accountNameColumnIndex]
			accountType := items[accountTypeColumnIndex]

			if accountType == iifAccountTypeIncome {
				incomeAccountNames[accountName] = true
			} else if accountType == iifAccountTypeExpense {
				expenseAccountNames[accountName] = true
			}
		}
	}

	return incomeAccountNames, expenseAccountNames
}
