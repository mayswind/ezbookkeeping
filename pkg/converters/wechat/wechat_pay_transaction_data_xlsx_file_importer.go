package wechat

import (
	"time"

	"github.com/Paxtiny/oscar/pkg/converters/converter"
	"github.com/Paxtiny/oscar/pkg/converters/datatable"
	"github.com/Paxtiny/oscar/pkg/converters/excel"
	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/errs"
	"github.com/Paxtiny/oscar/pkg/log"
	"github.com/Paxtiny/oscar/pkg/models"
)

// wechatPayTransactionDataXlsxFileImporter defines the structure of wechatPay xlsx importer for transaction data
type wechatPayTransactionDataXlsxFileImporter struct {
	dataHeaderStartContentBeginning string
}

// Initialize a webchat pay transaction data xlsx file importer singleton instance
var (
	WeChatPayTransactionDataXlsxFileImporter = &wechatPayTransactionDataXlsxFileImporter{}
)

// ParseImportedData returns the imported data by parsing the wechat pay transaction csv data
func (c *wechatPayTransactionDataXlsxFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	xlsxDataTable, err := excel.CreateNewExcelOOXMLFileBasicDataTable(data, false)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTable, err := createNewWeChatPayTransactionBasicDataTable(ctx, xlsxDataTable)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewCommonDataTableFromBasicDataTable(dataTable)

	if !commonDataTable.HasColumn(wechatPayTransactionTimeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionCategoryColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionTypeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionAmountColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionStatusColumnName) {
		log.Errorf(ctx, "[wechat_pay_transaction_data_xlsx_file_importer.ParseImportedData] cannot parse wechat pay xlsx data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionRowParser := createWeChatPayTransactionDataRowParser(dataTable.HeaderColumnNames())
	transactionDataTable := datatable.CreateNewTransactionDataTableFromCommonDataTable(commonDataTable, wechatPayTransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(wechatPayTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
