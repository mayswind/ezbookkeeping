package iif

import (
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
	currentSplitDataIndex int
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
		datasets := t.transactionDatasets[i]

		for j := 0; j < len(datasets.transactions); j++ {
			transaction := datasets.transactions[j]

			if transaction.splitData != nil {
				totalDataRowCount += len(transaction.splitData)
			}
		}
	}

	return totalDataRowCount
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *iifTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &iifTransactionDataRowIterator{
		dataTable:             t,
		currentDatasetIndex:   0,
		currentIndexInDataset: 0,
		currentSplitDataIndex: -1,
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
	} else if t.currentIndexInDataset < len(currentDataset.transactions) &&
		t.currentSplitDataIndex+1 < len(currentDataset.transactions[t.currentIndexInDataset].splitData) {
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

// Next returns the next transaction data row
func (t *iifTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	allDatasets := t.dataTable.transactionDatasets

	for i := t.currentDatasetIndex; i < len(allDatasets); i++ {
		foundNextRow := false
		dataset := allDatasets[i]

		for j := t.currentIndexInDataset; j < len(dataset.transactions); j++ {
			if t.currentSplitDataIndex+1 < len(dataset.transactions[j].splitData) {
				t.currentSplitDataIndex++
				foundNextRow = true
				break
			}

			t.currentIndexInDataset++
			t.currentSplitDataIndex = -1
		}

		if foundNextRow {
			break
		}

		t.currentDatasetIndex++
		t.currentIndexInDataset = 0
		t.currentSplitDataIndex = -1
	}

	if t.currentDatasetIndex >= len(allDatasets) {
		return nil, nil
	}

	currentDataset := allDatasets[t.currentDatasetIndex]

	if t.currentIndexInDataset >= len(currentDataset.transactions) {
		return nil, nil
	}

	data := currentDataset.transactions[t.currentIndexInDataset]

	if len(data.splitData) < 1 {
		log.Errorf(ctx, "[iif_transaction_data_table.Next] cannot parsing transaction in row#%d (dataset#%d), because split data is empty", t.currentIndexInDataset, t.currentDatasetIndex)
		return nil, errs.ErrInvalidIIFFile
	}

	if t.currentSplitDataIndex >= len(data.splitData) {
		return nil, nil
	}

	if len(data.splitData) > 1 {
		_, err := t.isSplitTransactionSupported(ctx, currentDataset, data)

		if err != nil {
			return nil, err
		}
	}

	rowItems, err := t.parseTransaction(ctx, user, currentDataset, data, t.currentSplitDataIndex)

	if err != nil {
		log.Errorf(ctx, "[iif_transaction_data_table.Next] cannot parsing transaction in row#%d-split#%d (dataset#%d), because %s", t.currentIndexInDataset, t.currentSplitDataIndex, t.currentDatasetIndex, err.Error())
		return nil, err
	}

	return &iifTransactionDataRow{
		dataTable:  t.dataTable,
		finalItems: rowItems,
	}, nil
}

func (t *iifTransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, dataset *iifTransactionDataset, transactionData *iifTransactionData, splitDataIndex int) (map[datatable.TransactionDataTableColumn]string, error) {
	var err error

	data := make(map[datatable.TransactionDataTableColumn]string, len(iifTransactionSupportedColumns))
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME], err = t.parseTransactionTime(dataset, transactionData)

	if err != nil {
		return nil, err
	}

	transactionType, _ := dataset.getSplitDataItemValue(transactionData.splitData[splitDataIndex], iifTransactionTypeColumnName)
	mainAccountName, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionAccountNameColumnName)
	splitAccountName, _ := dataset.getSplitDataItemValue(transactionData.splitData[splitDataIndex], iifTransactionAccountNameColumnName)
	mainAmount, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionAmountColumnName)
	splitAmount, _ := dataset.getSplitDataItemValue(transactionData.splitData[splitDataIndex], iifTransactionAmountColumnName)
	mainAmountNum, err := parseAmount(mainAmount)

	if err != nil {
		return nil, errs.ErrAmountInvalid
	}

	splitAmountNum, err := parseAmount(splitAmount)

	if err != nil {
		return nil, errs.ErrAmountInvalid
	}

	if transactionType == iifTransactionTypeBeginningBalance { // balance modification
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = mainAccountName
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(mainAmountNum)
	} else if (t.dataTable.incomeAccountNames[mainAccountName] && !t.dataTable.incomeAccountNames[splitAccountName] && !t.dataTable.expenseAccountNames[splitAccountName]) ||
		(t.dataTable.incomeAccountNames[splitAccountName] && !t.dataTable.incomeAccountNames[mainAccountName] && !t.dataTable.expenseAccountNames[mainAccountName]) { // income
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
		categoryName := ""
		accountName := ""
		amountNum := int64(0)

		if t.dataTable.incomeAccountNames[mainAccountName] && !t.dataTable.incomeAccountNames[splitAccountName] {
			categoryName = mainAccountName
			accountName = splitAccountName

			if len(transactionData.splitData) > 1 {
				amountNum = splitAmountNum
			} else {
				amountNum = -mainAmountNum
			}
		} else if t.dataTable.incomeAccountNames[splitAccountName] && !t.dataTable.incomeAccountNames[mainAccountName] {
			categoryName = splitAccountName
			accountName = mainAccountName

			if len(transactionData.splitData) > 1 {
				amountNum = -splitAmountNum
			} else {
				amountNum = mainAmountNum
			}
		} else {
			log.Errorf(ctx, "[iif_transaction_data_table.parseTransaction] cannot parse transaction, because main account \"%s\" and split account \"%s\" are all income account", mainAccountName, splitAccountName)
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
	} else if (t.dataTable.expenseAccountNames[mainAccountName] && !t.dataTable.expenseAccountNames[splitAccountName] && !t.dataTable.incomeAccountNames[splitAccountName]) ||
		(t.dataTable.expenseAccountNames[splitAccountName] && !t.dataTable.expenseAccountNames[mainAccountName] && !t.dataTable.incomeAccountNames[mainAccountName]) { // expense
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
		categoryName := ""
		accountName := ""
		amountNum := int64(0)

		if t.dataTable.expenseAccountNames[mainAccountName] && !t.dataTable.expenseAccountNames[splitAccountName] {
			categoryName = mainAccountName
			accountName = splitAccountName

			if len(transactionData.splitData) > 1 {
				amountNum = -splitAmountNum
			} else {
				amountNum = mainAmountNum
			}
		} else if t.dataTable.expenseAccountNames[splitAccountName] && !t.dataTable.expenseAccountNames[mainAccountName] {
			categoryName = splitAccountName
			accountName = mainAccountName

			if len(transactionData.splitData) > 1 {
				amountNum = splitAmountNum
			} else {
				amountNum = -mainAmountNum
			}
		} else {
			log.Errorf(ctx, "[iif_transaction_data_table.parseTransaction] cannot parse transaction, because main account \"%s\" and split account \"%s\" are all expense account", mainAccountName, splitAccountName)
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
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = iifTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
		amountNum := int64(0)
		relatedAmountNum := int64(0)
		mainAccountTransferToSplitAccount := false

		if len(transactionData.splitData) > 1 {
			amountNum = splitAmountNum
			relatedAmountNum = splitAmountNum
			mainAccountTransferToSplitAccount = amountNum >= 0
		} else {
			if mainAmountNum >= 0 {
				amountNum = splitAmountNum
				relatedAmountNum = mainAmountNum
				mainAccountTransferToSplitAccount = false
			} else if splitAmountNum >= 0 {
				amountNum = mainAmountNum
				relatedAmountNum = splitAmountNum
				mainAccountTransferToSplitAccount = true
			}
		}

		if mainAccountTransferToSplitAccount {
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = mainAccountName
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = splitAccountName
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = splitAccountName
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = mainAccountName
		}

		if amountNum >= 0 {
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amountNum)
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amountNum)
		}

		if relatedAmountNum >= 0 {
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(relatedAmountNum)
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(-relatedAmountNum)
		}
	}

	if splitMemo, _ := dataset.getSplitDataItemValue(transactionData.splitData[splitDataIndex], iifTransactionMemoColumnName); splitMemo != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = splitMemo
	} else if memo, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionMemoColumnName); memo != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = memo
	} else if splitName, _ := dataset.getSplitDataItemValue(transactionData.splitData[splitDataIndex], iifTransactionNameColumnName); splitName != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = splitName
	} else if name, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionNameColumnName); name != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = name
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	return data, nil
}

func (t *iifTransactionDataRowIterator) isSplitTransactionSupported(ctx core.Context, dataset *iifTransactionDataset, transactionData *iifTransactionData) (bool, error) {
	supportSplitTransactions := true
	transactionType, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionTypeColumnName)

	if transactionType == iifTransactionTypeBeginningBalance { // balance modification
		supportSplitTransactions = false
		log.Errorf(ctx, "[iif_transaction_data_table.isSplitTransactionSupported] cannot parse split balance modification transaction#%d (dataset#%d)", t.currentIndexInDataset, t.currentDatasetIndex)
	} else {
		transactionAmountStr, _ := dataset.getTransactionDataItemValue(transactionData, iifTransactionAmountColumnName)
		transactionAmount, err := parseAmount(transactionAmountStr)

		if err != nil {
			log.Errorf(ctx, "[iif_transaction_data_table.isSplitTransactionSupported] cannot parsing transaction in row#%d (dataset#%d), because transaction amount \"%s\" is invalid", t.currentIndexInDataset, t.currentDatasetIndex, transactionAmountStr)
			return false, errs.ErrAmountInvalid
		}

		splitTotalAmount := int64(0)

		for i := 0; i < len(transactionData.splitData); i++ {
			splitAmountStr, _ := dataset.getSplitDataItemValue(transactionData.splitData[i], iifTransactionAmountColumnName)
			splitAmount, err := parseAmount(splitAmountStr)

			if err != nil {
				log.Errorf(ctx, "[iif_transaction_data_table.isSplitTransactionSupported] cannot parsing transaction in row#%d-split#%d (dataset#%d), because split amount \"%s\" is invalid", t.currentIndexInDataset, i, t.currentDatasetIndex, splitAmountStr)
				return false, errs.ErrAmountInvalid
			}

			splitTotalAmount += splitAmount
		}

		if splitTotalAmount != -transactionAmount {
			supportSplitTransactions = false
			log.Errorf(ctx, "[iif_transaction_data_table.isSplitTransactionSupported] cannot parse split transaction#%d (dataset#%d), because the sum amount of each split data \"%d\" not equal to the transaction amount \"%d\"", t.currentIndexInDataset, t.currentDatasetIndex, splitTotalAmount, -transactionAmount)
		}
	}

	if len(transactionData.splitData) > 1 && !supportSplitTransactions {
		return false, errs.ErrNotSupportedSplitTransactions
	}

	return true, nil
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

	if utils.IsValidYearMonthDayLongOrShortDateFormat(date) && !utils.IsValidMonthDayYearLongOrShortDateFormat(date) {
		year = dateParts[0]
		month = dateParts[1]
		day = dateParts[2]
	}

	return utils.FormatYearMonthDayToLongDateTime(year, month, day)
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

func parseAmount(amount string) (int64, error) {
	return utils.ParseAmount(strings.ReplaceAll(amount, ",", ""))
}
