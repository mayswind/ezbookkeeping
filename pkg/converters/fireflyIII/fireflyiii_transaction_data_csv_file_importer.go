package fireflyIII

import (
	"bytes"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var fireflyIIITransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:             true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           true,
	datatable.TRANSACTION_DATA_TABLE_TAGS:                     true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:              true,
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
	dataTable, err := csv.CreateNewCsvBasicDataTable(ctx, reader)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewCommonDataTableFromBasicDataTable(dataTable)

	if !commonDataTable.HasColumn(fireflyIIITransactionTimeColumnName) ||
		!commonDataTable.HasColumn(fireflyIIITransactionTypeColumnName) ||
		!commonDataTable.HasColumn(fireflyIIITransactionSourceAccountNameColumnName) ||
		!commonDataTable.HasColumn(fireflyIIITransactionSourceAccountTypeColumnName) ||
		!commonDataTable.HasColumn(fireflyIIITransactionDestinationAccountNameColumnName) ||
		!commonDataTable.HasColumn(fireflyIIITransactionDestinationAccountTypeColumnName) ||
		!commonDataTable.HasColumn(fireflyIIITransactionAmountColumnName) {
		log.Errorf(ctx, "[fireflyiii_transaction_data_csv_file_importer.ParseImportedData] cannot parse Firefly III csv data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionRowParser := createFireflyIIITransactionDataRowParser()
	transactionDataTable := datatable.CreateNewTransactionDataTableFromCommonDataTable(commonDataTable, fireflyIIITransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewImporterWithTypeNameMapping(fireflyIIITransactionTypeNameMapping, "", "", ",")

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
