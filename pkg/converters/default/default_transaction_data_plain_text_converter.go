package _default

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// defaultTransactionDataPlainTextConverter defines the structure of ezbookkeeping default plain text converter for transaction data
type defaultTransactionDataPlainTextConverter struct {
	columnSeparator string
}

const ezbookkeepingLineSeparator = "\n"
const ezbookkeepingGeoLocationSeparator = " "
const ezbookkeepingTagSeparator = ";"

var ezbookkeepingDataColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         "Time",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE:     "Timezone",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         "Type",
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:                 "Category",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:             "Sub Category",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             "Account",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         "Account Currency",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   "Amount",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     "Account2",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: "Account2 Currency",
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           "Account2 Amount",
	datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION:      "Geographic Location",
	datatable.TRANSACTION_DATA_TABLE_TAGS:                     "Tags",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:              "Description",
}

var ezbookkeepingTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: "Balance Modification",
	models.TRANSACTION_TYPE_INCOME:         "Income",
	models.TRANSACTION_TYPE_EXPENSE:        "Expense",
	models.TRANSACTION_TYPE_TRANSFER:       "Transfer",
}

var ezbookkeepingDataColumns = []datatable.TransactionDataTableColumn{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE,
	datatable.TRANSACTION_DATA_TABLE_CATEGORY,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT,
	datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION,
	datatable.TRANSACTION_DATA_TABLE_TAGS,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION,
}

// ToExportedContent returns the exported transaction plain text data
func (c *defaultTransactionDataPlainTextConverter) ToExportedContent(ctx core.Context, uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error) {
	dataTableBuilder := createNewDefaultTransactionPlainTextDataTableBuilder(
		len(transactions),
		ezbookkeepingDataColumns,
		ezbookkeepingDataColumnNameMapping,
		c.columnSeparator,
		ezbookkeepingLineSeparator,
	)

	dataTableExporter := converter.CreateNewExporter(
		ezbookkeepingTransactionTypeNameMapping,
		ezbookkeepingGeoLocationSeparator,
		ezbookkeepingTagSeparator,
	)

	err := dataTableExporter.BuildExportedContent(ctx, dataTableBuilder, uid, transactions, accountMap, categoryMap, tagMap, allTagIndexes)

	if err != nil {
		return nil, err
	}

	return []byte(dataTableBuilder.String()), nil
}

// ParseImportedData returns the imported data by parsing the transaction plain text data
func (c *defaultTransactionDataPlainTextConverter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTable, err := createNewDefaultPlainTextDataTable(
		string(data),
		c.columnSeparator,
		ezbookkeepingLineSeparator,
	)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionDataTable := datatable.CreateNewImportedTransactionDataTable(dataTable, ezbookkeepingDataColumnNameMapping)

	dataTableImporter := converter.CreateNewImporterWithTypeNameMapping(
		ezbookkeepingTransactionTypeNameMapping,
		ezbookkeepingGeoLocationSeparator,
		ezbookkeepingTagSeparator,
	)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
