package qif

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var qifTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: utils.IntToString(int(models.TRANSACTION_TYPE_MODIFY_BALANCE)),
	models.TRANSACTION_TYPE_INCOME:         utils.IntToString(int(models.TRANSACTION_TYPE_INCOME)),
	models.TRANSACTION_TYPE_EXPENSE:        utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE)),
	models.TRANSACTION_TYPE_TRANSFER:       utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER)),
}

// qifTransactionDataImporter defines the structure of quicken interchange format (qif) importer for transaction data
type qifTransactionDataImporter struct {
	dateFormatType qifDateFormatType
}

// Initialize a quicken interchange format (qif) transaction data importer singleton instance
var (
	QifYearMonthDayTransactionDataImporter = &qifTransactionDataImporter{
		dateFormatType: qifYearMonthDayDateFormat,
	}

	QifMonthDayYearTransactionDataImporter = &qifTransactionDataImporter{
		dateFormatType: qifMonthDayYearDateFormat,
	}

	QifDayMonthYearTransactionDataImporter = &qifTransactionDataImporter{
		dateFormatType: qifDayMonthYearDateFormat,
	}
)

// ParseImportedData returns the imported data by parsing the quicken interchange format (qif) transaction data
func (c *qifTransactionDataImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	qifDataReader := createNewQifDataReader(data)
	qifData, err := qifDataReader.read(ctx)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionDataTable, err := createNewQifTransactionDataTable(c.dateFormatType, qifData)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(qifTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
