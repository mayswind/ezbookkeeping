package alipay

import (
	"bytes"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var alipayTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   "收入",
	models.TRANSACTION_TYPE_EXPENSE:  "支出",
	models.TRANSACTION_TYPE_TRANSFER: "不计收支",
}

// alipayTransactionDataCsvImporter defines the structure of alipay csv importer for transaction data
type alipayTransactionDataCsvImporter struct {
	fileHeaderLine         string
	dataHeaderStartContent string
	dataBottomEndLineRune  rune
	originalColumnNames    alipayTransactionColumnNames
}

// ParseImportedData returns the imported data by parsing the alipay transaction csv data
func (c *alipayTransactionDataCsvImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]*models.TransactionCategory, incomeCategoryMap map[string]*models.TransactionCategory, transferCategoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	enc := simplifiedchinese.GB18030
	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())

	dataTable, err := createNewAlipayTransactionPlainTextDataTable(
		ctx,
		reader,
		c.fileHeaderLine,
		c.dataHeaderStartContent,
		c.dataBottomEndLineRune,
		c.originalColumnNames,
	)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := datatable.CreateNewSimpleImporter(
		dataTable.GetDataColumnMapping(),
		alipayTransactionTypeNameMapping,
	)

	return dataTableImporter.ParseImportedData(ctx, user, dataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
