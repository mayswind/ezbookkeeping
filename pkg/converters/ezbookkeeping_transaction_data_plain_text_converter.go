package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ezBookKeepingTransactionDataPlainTextConverter defines the structure of plain file converter for transaction data
type ezBookKeepingTransactionDataPlainTextConverter struct {
	DataTableTransactionDataConverter
	columns []DataTableColumn
}

var ezbookkeepingDataColumnNameMapping = map[DataTableColumn]string{
	DATA_TABLE_TRANSACTION_TIME:         "Time",
	DATA_TABLE_TRANSACTION_TIMEZONE:     "Timezone",
	DATA_TABLE_TRANSACTION_TYPE:         "Type",
	DATA_TABLE_CATEGORY:                 "Category",
	DATA_TABLE_SUB_CATEGORY:             "Sub Category",
	DATA_TABLE_ACCOUNT_NAME:             "Account",
	DATA_TABLE_ACCOUNT_CURRENCY:         "Account Currency",
	DATA_TABLE_AMOUNT:                   "Amount",
	DATA_TABLE_RELATED_ACCOUNT_NAME:     "Account2",
	DATA_TABLE_RELATED_ACCOUNT_CURRENCY: "Account2 Currency",
	DATA_TABLE_RELATED_AMOUNT:           "Account2 Amount",
	DATA_TABLE_GEOGRAPHIC_LOCATION:      "Geographic Location",
	DATA_TABLE_TAGS:                     "Tags",
	DATA_TABLE_DESCRIPTION:              "Description",
}

var ezbookkeepingTransactionTypeNameMapping = map[models.TransactionDbType]string{
	models.TRANSACTION_DB_TYPE_MODIFY_BALANCE: "Balance Modification",
	models.TRANSACTION_DB_TYPE_INCOME:         "Income",
	models.TRANSACTION_DB_TYPE_EXPENSE:        "Expense",
	models.TRANSACTION_DB_TYPE_TRANSFER_OUT:   "Transfer",
	models.TRANSACTION_DB_TYPE_TRANSFER_IN:    "Transfer",
}

var ezbookkeepingNameTransactionTypeMapping = map[string]models.TransactionDbType{
	"Balance Modification": models.TRANSACTION_DB_TYPE_MODIFY_BALANCE,
	"Income":               models.TRANSACTION_DB_TYPE_INCOME,
	"Expense":              models.TRANSACTION_DB_TYPE_EXPENSE,
	"Transfer":             models.TRANSACTION_DB_TYPE_TRANSFER_OUT,
}

var ezbookkeepingDataColumns = []DataTableColumn{
	DATA_TABLE_TRANSACTION_TIME,
	DATA_TABLE_TRANSACTION_TIMEZONE,
	DATA_TABLE_TRANSACTION_TYPE,
	DATA_TABLE_CATEGORY,
	DATA_TABLE_SUB_CATEGORY,
	DATA_TABLE_ACCOUNT_NAME,
	DATA_TABLE_ACCOUNT_CURRENCY,
	DATA_TABLE_AMOUNT,
	DATA_TABLE_RELATED_ACCOUNT_NAME,
	DATA_TABLE_RELATED_ACCOUNT_CURRENCY,
	DATA_TABLE_RELATED_AMOUNT,
	DATA_TABLE_GEOGRAPHIC_LOCATION,
	DATA_TABLE_TAGS,
	DATA_TABLE_DESCRIPTION,
}

// ToExportedContent returns the exported plain text transaction data
func (c *ezBookKeepingTransactionDataPlainTextConverter) ToExportedContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error) {
	dataTableBuilder := createNewezbookkeepingTransactionPlainTextDataTableBuilder(len(transactions), c.columns, c.dataColumnMapping, c.columnSeparator, c.lineSeparator)
	err := c.buildExportedContent(dataTableBuilder, uid, transactions, accountMap, categoryMap, tagMap, allTagIndexes)

	if err != nil {
		return nil, err
	}

	return []byte(dataTableBuilder.String()), nil
}

func (c *ezBookKeepingTransactionDataPlainTextConverter) ParseImportedData(user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) ([]*models.Transaction, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTable, err := createNewezbookkeepingTransactionPlainTextDataTable(string(data), c.columnSeparator, c.lineSeparator)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	return c.parseImportedData(user, dataTable, defaultTimezoneOffset, accountMap, categoryMap, tagMap)
}
