package jdcom

import (
	"bytes"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// jdComFinanceTransactionDataCsvFileImporter defines the structure of jd.com finance csv importer for transaction data
type jdComFinanceTransactionDataCsvFileImporter struct {
	fileHeaderLineBeginning         string
	dataHeaderStartContentBeginning string
}

// Initialize a jd.com finance transaction data csv file importer singleton instance
var (
	JDComFinanceTransactionDataCsvFileImporter = &jdComFinanceTransactionDataCsvFileImporter{}
)

// ParseImportedData returns the imported data by parsing the jd.com finance transaction csv data
func (c *jdComFinanceTransactionDataCsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, []string, error) {
	fallback := unicode.UTF8.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(fallback))

	csvDataTable, err := csv.CreateNewCsvBasicDataTable(ctx, reader, false)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	dataTable, err := createNewJDComFinanceTransactionBasicDataTable(ctx, csvDataTable)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewCommonDataTableFromBasicDataTable(dataTable)

	if !commonDataTable.HasColumn(jdComFinanceTransactionTimeColumnName) ||
		!commonDataTable.HasColumn(jdComFinanceTransactionMerchantNameColumnName) ||
		!commonDataTable.HasColumn(jdComFinanceTransactionMemoColumnName) ||
		!commonDataTable.HasColumn(jdComFinanceTransactionAmountColumnName) ||
		!commonDataTable.HasColumn(jdComFinanceTransactionRelatedAccountColumnName) ||
		!commonDataTable.HasColumn(jdComFinanceTransactionStatusColumnName) ||
		!commonDataTable.HasColumn(jdComFinanceTransactionTypeColumnName) {
		log.Errorf(ctx, "[jdcom_finance_transaction_data_csv_file_importer.ParseImportedData] cannot parse jd.com finance csv data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionRowParser := createJDComFinanceTransactionDataRowParser(dataTable.HeaderColumnNames())
	transactionDataTable := datatable.CreateNewTransactionDataTableFromCommonDataTable(commonDataTable, jdComFinanceTransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(jdComFinanceTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
