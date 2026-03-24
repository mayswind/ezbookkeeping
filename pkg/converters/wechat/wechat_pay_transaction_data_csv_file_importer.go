package wechat

import (
	"bytes"
	"time"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/Paxtiny/oscar/pkg/converters/converter"
	"github.com/Paxtiny/oscar/pkg/converters/csv"
	"github.com/Paxtiny/oscar/pkg/converters/datatable"
	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/errs"
	"github.com/Paxtiny/oscar/pkg/log"
	"github.com/Paxtiny/oscar/pkg/models"
)

// wechatPayTransactionDataCsvFileImporter defines the structure of wechatPay csv importer for transaction data
type wechatPayTransactionDataCsvFileImporter struct {
	fileHeaderLineBeginning         string
	dataHeaderStartContentBeginning string
}

// Initialize a webchat pay transaction data csv file importer singleton instance
var (
	WeChatPayTransactionDataCsvFileImporter = &wechatPayTransactionDataCsvFileImporter{}
)

// ParseImportedData returns the imported data by parsing the wechat pay transaction csv data
func (c *wechatPayTransactionDataCsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	fallback := unicode.UTF8.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(fallback))

	csvDataTable, err := csv.CreateNewCsvBasicDataTable(ctx, reader, false)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTable, err := createNewWeChatPayTransactionBasicDataTable(ctx, csvDataTable)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewCommonDataTableFromBasicDataTable(dataTable)

	if !commonDataTable.HasColumn(wechatPayTransactionTimeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionCategoryColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionTypeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionAmountColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionStatusColumnName) {
		log.Errorf(ctx, "[wechat_pay_transaction_data_csv_file_importer.ParseImportedData] cannot parse wechat pay csv data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionRowParser := createWeChatPayTransactionDataRowParser(dataTable.HeaderColumnNames())
	transactionDataTable := datatable.CreateNewTransactionDataTableFromCommonDataTable(commonDataTable, wechatPayTransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(wechatPayTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
