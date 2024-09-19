package _default

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ezBookKeepingTransactionDataPlainTextConverter defines the structure of ezbookkeeping plain text converter for transaction data
type ezBookKeepingTransactionDataPlainTextConverter struct {
	columnSeparator string
}

const ezbookkeepingLineSeparator = "\n"
const ezbookkeepingGeoLocationSeparator = " "
const ezbookkeepingTagSeparator = ";"

var ezbookkeepingDataColumnNameMapping = map[datatable.DataTableColumn]string{
	datatable.DATA_TABLE_TRANSACTION_TIME:         "Time",
	datatable.DATA_TABLE_TRANSACTION_TIMEZONE:     "Timezone",
	datatable.DATA_TABLE_TRANSACTION_TYPE:         "Type",
	datatable.DATA_TABLE_CATEGORY:                 "Category",
	datatable.DATA_TABLE_SUB_CATEGORY:             "Sub Category",
	datatable.DATA_TABLE_ACCOUNT_NAME:             "Account",
	datatable.DATA_TABLE_ACCOUNT_CURRENCY:         "Account Currency",
	datatable.DATA_TABLE_AMOUNT:                   "Amount",
	datatable.DATA_TABLE_RELATED_ACCOUNT_NAME:     "Account2",
	datatable.DATA_TABLE_RELATED_ACCOUNT_CURRENCY: "Account2 Currency",
	datatable.DATA_TABLE_RELATED_AMOUNT:           "Account2 Amount",
	datatable.DATA_TABLE_GEOGRAPHIC_LOCATION:      "Geographic Location",
	datatable.DATA_TABLE_TAGS:                     "Tags",
	datatable.DATA_TABLE_DESCRIPTION:              "Description",
}

var ezbookkeepingTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: "Balance Modification",
	models.TRANSACTION_TYPE_INCOME:         "Income",
	models.TRANSACTION_TYPE_EXPENSE:        "Expense",
	models.TRANSACTION_TYPE_TRANSFER:       "Transfer",
}

var ezbookkeepingDataColumns = []datatable.DataTableColumn{
	datatable.DATA_TABLE_TRANSACTION_TIME,
	datatable.DATA_TABLE_TRANSACTION_TIMEZONE,
	datatable.DATA_TABLE_TRANSACTION_TYPE,
	datatable.DATA_TABLE_CATEGORY,
	datatable.DATA_TABLE_SUB_CATEGORY,
	datatable.DATA_TABLE_ACCOUNT_NAME,
	datatable.DATA_TABLE_ACCOUNT_CURRENCY,
	datatable.DATA_TABLE_AMOUNT,
	datatable.DATA_TABLE_RELATED_ACCOUNT_NAME,
	datatable.DATA_TABLE_RELATED_ACCOUNT_CURRENCY,
	datatable.DATA_TABLE_RELATED_AMOUNT,
	datatable.DATA_TABLE_GEOGRAPHIC_LOCATION,
	datatable.DATA_TABLE_TAGS,
	datatable.DATA_TABLE_DESCRIPTION,
}

// ToExportedContent returns the exported transaction plain text data
func (c *ezBookKeepingTransactionDataPlainTextConverter) ToExportedContent(ctx core.Context, uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error) {
	dataTableBuilder := createNewezbookkeepingTransactionPlainTextDataTableBuilder(
		len(transactions),
		ezbookkeepingDataColumns,
		ezbookkeepingDataColumnNameMapping,
		c.columnSeparator,
		ezbookkeepingLineSeparator,
	)

	dataTableExporter := datatable.CreateNewExporter(
		ezbookkeepingDataColumnNameMapping,
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
func (c *ezBookKeepingTransactionDataPlainTextConverter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTable, err := createNewezbookkeepingTransactionPlainTextDataTable(
		string(data),
		c.columnSeparator,
		ezbookkeepingLineSeparator,
	)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	dataTableImporter := datatable.CreateNewImporter(
		ezbookkeepingDataColumnNameMapping,
		ezbookkeepingTransactionTypeNameMapping,
		ezbookkeepingGeoLocationSeparator,
		ezbookkeepingTagSeparator,
	)

	return dataTableImporter.ParseImportedData(ctx, user, dataTable, defaultTimezoneOffset, accountMap, categoryMap, tagMap)
}
