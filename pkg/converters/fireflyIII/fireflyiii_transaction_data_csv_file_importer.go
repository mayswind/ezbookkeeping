package fireflyIII

import (
	"bytes"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var fireflyIIITransactionDataColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         "date",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         "type",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:             "category",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             "source_name",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         "currency_code",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   "amount",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     "destination_name",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: "foreign_currency_code",
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           "foreign_amount",
	datatable.TRANSACTION_DATA_TABLE_TAGS:                     "tags",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:              "description",
}

var fireflyIIITransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: "Opening balance",
	models.TRANSACTION_TYPE_INCOME:         "Deposit",
	models.TRANSACTION_TYPE_EXPENSE:        "Withdrawal",
	models.TRANSACTION_TYPE_TRANSFER:       "Transfer",
}

// fireflyIIITransactionDataCsvFileImporter defines the structure of firefly III csv importer for transaction data
type fireflyIIITransactionDataCsvFileImporter struct{}

// Initialize a firefly III transaction data csv file importer singleton instance
var (
	FireflyIIITransactionDataCsvFileImporter = &fireflyIIITransactionDataCsvFileImporter{}
)

// ParseImportedData returns the imported data by parsing the firefly III transaction csv data
func (c *fireflyIIITransactionDataCsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	reader := bytes.NewReader(data)
	dataTable, err := csv.CreateNewCsvImportedDataTable(ctx, reader)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionRowParser := createFireflyIIITransactionDataRowParser()
	transactionDataTable := datatable.CreateNewImportedTransactionDataTableWithRowParser(dataTable, fireflyIIITransactionDataColumnNameMapping, transactionRowParser)
	dataTableImporter := converter.CreateNewImporterWithTypeNameMapping(fireflyIIITransactionTypeNameMapping, "", ",")

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
