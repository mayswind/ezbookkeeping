package wechat

import (
	"bytes"
	"encoding/csv"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	csvdatatable "github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var wechatPayTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

var wechatPayTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   "收入",
	models.TRANSACTION_TYPE_EXPENSE:  "支出",
	models.TRANSACTION_TYPE_TRANSFER: "/",
}

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
func (c *wechatPayTransactionDataCsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	fallback := unicode.UTF8.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(fallback))

	dataTable, err := c.createNewWeChatPayImportedDataTable(ctx, reader)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewImportedCommonDataTable(dataTable)

	if !commonDataTable.HasColumn(wechatPayTransactionTimeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionCategoryColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionTypeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionAmountColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionStatusColumnName) {
		log.Errorf(ctx, "[wechat_pay_transaction_data_csv_file_importer.ParseImportedData] cannot parse wechat pay csv data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionRowParser := createWeChatPayTransactionDataRowParser()
	transactionDataTable := datatable.CreateNewCommonTransactionDataTable(commonDataTable, wechatPayTransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(wechatPayTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (c *wechatPayTransactionDataCsvFileImporter) createNewWeChatPayImportedDataTable(ctx core.Context, reader io.Reader) (datatable.ImportedDataTable, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	allOriginalLines := make([][]string, 0)
	hasFileHeader := false
	foundContentBeforeDataHeaderLine := false

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[wechat_pay_transaction_data_csv_file_importer.createNewWeChatPayImportedDataTable] cannot parse wechat pay csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], wechatPayTransactionDataCsvFileHeader) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[wechat_pay_transaction_data_csv_file_importer.createNewWeChatPayImportedDataTable] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
				continue
			}
		}

		if !foundContentBeforeDataHeaderLine {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], wechatPayTransactionDataHeaderStartContentBeginning) == 0 {
				foundContentBeforeDataHeaderLine = true
				continue
			} else {
				continue
			}
		}

		if foundContentBeforeDataHeaderLine {
			if len(items) <= 0 {
				continue
			}

			for i := 0; i < len(items); i++ {
				items[i] = strings.Trim(items[i], " ")
			}

			if len(allOriginalLines) > 0 && len(items) < len(allOriginalLines[0]) {
				log.Errorf(ctx, "[wechat_pay_transaction_data_csv_file_importer.createNewWeChatPayImportedDataTable] cannot parse row \"index:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", len(allOriginalLines), len(items), len(allOriginalLines[0]))
				return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
			}

			allOriginalLines = append(allOriginalLines, items)
		}
	}

	if !hasFileHeader || !foundContentBeforeDataHeaderLine {
		return nil, errs.ErrInvalidFileHeader
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[wechat_pay_transaction_data_csv_file_importer.createNewWeChatPayImportedDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	dataTable := csvdatatable.CreateNewCustomCsvImportedDataTable(allOriginalLines)

	return dataTable, nil
}
