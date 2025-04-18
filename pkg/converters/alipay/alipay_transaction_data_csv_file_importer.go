package alipay

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	csvdatatable "github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var alipayTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

var alipayTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   "收入",
	models.TRANSACTION_TYPE_EXPENSE:  "支出",
	models.TRANSACTION_TYPE_TRANSFER: "不计收支",
}

// alipayTransactionColumnNames defines the structure of alipay transaction plain text header names
type alipayTransactionColumnNames struct {
	timeColumnName           string
	categoryColumnName       string
	targetNameColumnName     string
	productNameColumnName    string
	amountColumnName         string
	typeColumnName           string
	relatedAccountColumnName string
	statusColumnName         string
	descriptionColumnName    string
}

// alipayTransactionDataCsvFileImporter defines the structure of alipay csv importer for transaction data
type alipayTransactionDataCsvFileImporter struct {
	fileHeaderLine         string
	dataHeaderStartContent string
	dataBottomEndLineRune  rune
	originalColumnNames    alipayTransactionColumnNames
}

// ParseImportedData returns the imported data by parsing the alipay transaction csv data
func (c *alipayTransactionDataCsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	enc := simplifiedchinese.GB18030
	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())

	dataTable, err := c.createNewAlipayImportedDataTable(ctx, reader, c.fileHeaderLine, c.dataHeaderStartContent, c.dataBottomEndLineRune)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewImportedCommonDataTable(dataTable)

	if !commonDataTable.HasColumn(c.originalColumnNames.timeColumnName) ||
		!commonDataTable.HasColumn(c.originalColumnNames.amountColumnName) ||
		!commonDataTable.HasColumn(c.originalColumnNames.typeColumnName) ||
		!commonDataTable.HasColumn(c.originalColumnNames.statusColumnName) {
		log.Errorf(ctx, "[alipay_transaction_data_csv_file_importer.ParseImportedData] cannot parse alipay csv data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionRowParser := createAlipayTransactionDataRowParser(c.originalColumnNames)
	transactionDataTable := datatable.CreateNewCommonTransactionDataTable(commonDataTable, alipayTransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(alipayTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (c *alipayTransactionDataCsvFileImporter) createNewAlipayImportedDataTable(ctx core.Context, reader io.Reader, fileHeaderLine string, dataHeaderStartContent string, dataBottomEndLineRune rune) (datatable.ImportedDataTable, error) {
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
			log.Errorf(ctx, "[alipay_transaction_csv_data_table.createNewAlipayImportedDataTable] cannot parse alipay csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], fileHeaderLine) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[alipay_transaction_csv_data_table.createNewAlipayImportedDataTable] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
				continue
			}
		}

		if !foundContentBeforeDataHeaderLine {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], dataHeaderStartContent) >= 0 {
				foundContentBeforeDataHeaderLine = true
				continue
			} else {
				continue
			}
		}

		if foundContentBeforeDataHeaderLine {
			if len(items) <= 0 {
				continue
			} else if len(items) == 1 && dataBottomEndLineRune > 0 && utils.ContainsOnlyOneRune(items[0], dataBottomEndLineRune) {
				break
			}

			for i := 0; i < len(items); i++ {
				items[i] = strings.Trim(items[i], " ")
			}

			if len(allOriginalLines) > 0 && len(items) < len(allOriginalLines[0]) {
				log.Errorf(ctx, "[alipay_transaction_csv_data_table.createNewAlipayImportedDataTable] cannot parse row \"index:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", len(allOriginalLines), len(items), len(allOriginalLines[0]))
				return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
			}

			allOriginalLines = append(allOriginalLines, items)
		}
	}

	if !hasFileHeader || !foundContentBeforeDataHeaderLine {
		return nil, errs.ErrInvalidFileHeader
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[alipay_transaction_csv_data_table.createNewAlipayImportedDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	dataTable := csvdatatable.CreateNewCustomCsvImportedDataTable(allOriginalLines)

	return dataTable, nil
}
