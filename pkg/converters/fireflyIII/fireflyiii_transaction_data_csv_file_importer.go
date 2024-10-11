package fireflyIII

import (
	"bytes"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// fireflyiiiTransactionDataCsvImporter defines the structure of firefly III csv importer for transaction data
type fireflyIIITransactionDataCsvImporter struct{}

// Initialize a firefly III transaction data csv file importer singleton instance
var (
	FireflyIIITransactionDataCsvImporter = &fireflyIIITransactionDataCsvImporter{}
)

// ParseImportedData returns the imported data by parsing the firefly iii transaction csv data
func (c *fireflyIIITransactionDataCsvImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]*models.TransactionCategory, incomeCategoryMap map[string]*models.TransactionCategory, transferCategoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	reader := bytes.NewReader(data)

	dataTable, err := createNewFireflyIIITransactionPlainTextDataTable(
		ctx,
		reader,
	)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := datatable.CreateNewImporter(
		dataTable.GetDataColumnMapping(),
		fireflyIIITransactionTypeNameMapping,
		"",
		",",
	)

	return dataTableImporter.ParseImportedData(ctx, user, dataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
