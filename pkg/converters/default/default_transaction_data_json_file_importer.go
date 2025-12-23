package _default

import (
	"encoding/json"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var allJsonDataSupportedColumns = []datatable.TransactionDataTableColumn{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT,
	datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION,
	datatable.TRANSACTION_DATA_TABLE_TAGS,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION,
}

// defaultTransactionDataJsonImporter defines the structure of ezbookkeeping default json importer for transaction data
type defaultTransactionDataJsonImporter struct{}

// Initialize an ezbookkeeping default transaction data json file importer singleton instance
var (
	DefaultTransactionDataJsonFileImporter = &defaultTransactionDataJsonImporter{}
)

// ParseImportedData returns the imported data by parsing the transaction json data
func (c *defaultTransactionDataJsonImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	var importRequest models.ImportTransactionRequest

	if err := json.Unmarshal(data, &importRequest); err != nil {
		return nil, nil, nil, nil, nil, nil, errs.ErrInvalidJSONFile
	}

	transactionDataTable, err := c.createNewDefaultTransactionDataTable(importRequest)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewImporterWithTypeNameMapping(
		ezbookkeepingTransactionTypeNameMapping,
		ezbookkeepingGeoLocationSeparator,
		ezbookkeepingGeoLocationOrder,
		ezbookkeepingTagSeparator,
	)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (c *defaultTransactionDataJsonImporter) createNewDefaultTransactionDataTable(importRequest models.ImportTransactionRequest) (datatable.TransactionDataTable, error) {
	transactionDataTable := datatable.CreateNewWritableTransactionDataTable(allJsonDataSupportedColumns)

	if importRequest.Transactions == nil || len(importRequest.Transactions) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	for i := 0; i < len(importRequest.Transactions); i++ {
		transaction := importRequest.Transactions[i]

		utcOffset, err := utils.StringToInt(transaction.UtcOffset)

		if err != nil {
			return nil, errs.ErrTransactionTimeZoneInvalid
		}

		timezone := time.FixedZone("Transaction Timezone", utcOffset*60)
		timezoneOffset := utils.FormatTimezoneOffset(time.Now().Unix(), timezone)

		row := make(map[datatable.TransactionDataTableColumn]string, len(allJsonDataSupportedColumns))
		row[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = transaction.Time
		row[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = timezoneOffset
		row[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = transaction.Type
		row[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = transaction.CategoryName
		row[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = transaction.SourceAccountName
		row[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = transaction.SourceAmount
		row[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = transaction.DestinationAccountName
		row[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = transaction.DestinationAmount
		row[datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION] = transaction.GeoLocation
		row[datatable.TRANSACTION_DATA_TABLE_TAGS] = transaction.TagNames
		row[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = transaction.Comment

		transactionDataTable.Add(row)
	}

	return transactionDataTable, nil
}
