package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// feideeMymoneyTransactionDataXlsImporter defines the structure of feidee mymoney xls importer for transaction data
type feideeMymoneyTransactionDataXlsImporter struct {
	DataTableTransactionDataImporter
}

// Initialize a feidee mymoney transaction data xls file importer singleton instance
var (
	FeideeMymoneyTransactionDataXlsImporter = &feideeMymoneyTransactionDataXlsImporter{
		DataTableTransactionDataImporter{
			dataColumnMapping:      feideeMymoneyDataColumnNameMapping,
			transactionTypeMapping: feideeMymoneyTransactionTypeNameMapping,
			postProcessFunc:        feideeMymoneyTransactionDataImporterPostProcess,
		},
	}
)

// ParseImportedData returns the imported data by parsing the feidee mymoney transaction xls data
func (c *feideeMymoneyTransactionDataXlsImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTable, err := createNewFeideeMymoneyTransactionExcelFileDataTable(data)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	return c.parseImportedData(ctx, user, dataTable, defaultTimezoneOffset, accountMap, categoryMap, tagMap)
}