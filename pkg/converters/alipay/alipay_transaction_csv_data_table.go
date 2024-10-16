package alipay

import (
	"encoding/csv"
	"io"
	"strings"

	csvdatatable "github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/locales"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const alipayTransactionDataStatusSuccessName = "交易成功"
const alipayTransactionDataStatusPaymentSuccessName = "支付成功"
const alipayTransactionDataStatusRepaymentSuccessName = "还款成功"
const alipayTransactionDataStatusClosedName = "交易关闭"
const alipayTransactionDataStatusRefundSuccessName = "退款成功"
const alipayTransactionDataStatusTaxRefundSuccessName = "退税成功"

const alipayTransactionDataProductNameTransferToAlipayPrefix = "充值-"
const alipayTransactionDataProductNameTransferFromAlipayPrefix = "提现-"
const alipayTransactionDataProductNameTransferInText = "转入"
const alipayTransactionDataProductNameTransferOutText = "转出"
const alipayTransactionDataProductNameRepaymentText = "还款"

var alipayTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
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

// alipayTransactionDataTable defines the structure of alipay transaction plain text data table
type alipayTransactionDataTable struct {
	innerDataTable datatable.CommonDataTable
	columns        alipayTransactionColumnNames
}

// alipayTransactionDataRow defines the structure of alipay transaction plain text data row
type alipayTransactionDataRow struct {
	isValid    bool
	finalItems map[datatable.TransactionDataTableColumn]string
}

// alipayTransactionDataRowIterator defines the structure of alipay transaction plain text data row iterator
type alipayTransactionDataRowIterator struct {
	dataTable     *alipayTransactionDataTable
	innerIterator datatable.CommonDataRowIterator
}

// HasColumn returns whether the transaction data table has specified column
func (t *alipayTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := alipayTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *alipayTransactionDataTable) TransactionRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *alipayTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &alipayTransactionDataRowIterator{
		dataTable:     t,
		innerIterator: t.innerDataTable.DataRowIterator(),
	}
}

// IsValid returns whether this row is valid data for importing
func (r *alipayTransactionDataRow) IsValid() bool {
	return r.isValid
}

// GetData returns the data in the specified column type
func (r *alipayTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := alipayTransactionSupportedColumns[column]

	if !exists {
		return ""
	}

	return r.finalItems[column]
}

// HasNext returns whether the iterator does not reach the end
func (t *alipayTransactionDataRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// Next returns the next imported data row
func (t *alipayTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	importedRow := t.innerIterator.Next()

	if importedRow == nil {
		return nil, nil
	}

	finalItems, isValid, err := t.dataTable.parseTransactionData(ctx, user, importedRow, t.innerIterator.CurrentRowId())

	if err != nil {
		return nil, err
	}

	return &alipayTransactionDataRow{
		isValid:    isValid,
		finalItems: finalItems,
	}, nil
}

func (t *alipayTransactionDataTable) hasOriginalColumn(columnName string) bool {
	return columnName != "" && t.innerDataTable.HasColumn(columnName)
}

func (t *alipayTransactionDataTable) parseTransactionData(ctx core.Context, user *models.User, dataRow datatable.CommonDataRow, rowId string) (map[datatable.TransactionDataTableColumn]string, bool, error) {
	if dataRow.GetData(t.columns.typeColumnName) != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
		dataRow.GetData(t.columns.typeColumnName) != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
		dataRow.GetData(t.columns.typeColumnName) != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		log.Warnf(ctx, "[alipay_transaction_csv_data_table.parseTransactionData] skip parsing transaction in row \"%s\", because type is \"%s\"", rowId, dataRow.GetData(t.columns.typeColumnName))
		return nil, false, nil
	}

	if dataRow.GetData(t.columns.statusColumnName) != alipayTransactionDataStatusSuccessName &&
		dataRow.GetData(t.columns.statusColumnName) != alipayTransactionDataStatusPaymentSuccessName &&
		dataRow.GetData(t.columns.statusColumnName) != alipayTransactionDataStatusRepaymentSuccessName &&
		dataRow.GetData(t.columns.statusColumnName) != alipayTransactionDataStatusClosedName &&
		dataRow.GetData(t.columns.statusColumnName) != alipayTransactionDataStatusRefundSuccessName &&
		dataRow.GetData(t.columns.statusColumnName) != alipayTransactionDataStatusTaxRefundSuccessName {
		log.Warnf(ctx, "[alipay_transaction_csv_data_table.parseTransactionData] skip parsing transaction in row \"%s\", because status is \"%s\"", rowId, dataRow.GetData(t.columns.statusColumnName))
		return nil, false, nil
	}

	data := make(map[datatable.TransactionDataTableColumn]string, len(alipayTransactionSupportedColumns))

	if t.hasOriginalColumn(t.columns.timeColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = dataRow.GetData(t.columns.timeColumnName)
	}

	if t.hasOriginalColumn(t.columns.categoryColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(t.columns.categoryColumnName)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
	}

	if t.hasOriginalColumn(t.columns.amountColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = dataRow.GetData(t.columns.amountColumnName)
	}

	if t.hasOriginalColumn(t.columns.descriptionColumnName) && dataRow.GetData(t.columns.descriptionColumnName) != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(t.columns.descriptionColumnName)
	} else if t.hasOriginalColumn(t.columns.productNameColumnName) && dataRow.GetData(t.columns.productNameColumnName) != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(t.columns.productNameColumnName)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	relatedAccountName := ""

	if t.hasOriginalColumn(t.columns.relatedAccountColumnName) {
		relatedAccountName = dataRow.GetData(t.columns.relatedAccountColumnName)
	}

	statusName := ""

	if t.hasOriginalColumn(t.columns.statusColumnName) {
		statusName = dataRow.GetData(t.columns.statusColumnName)
	}

	locale := user.Language

	if locale == "" {
		locale = ctx.GetClientLocale()
	}

	localeTextItems := locales.GetLocaleTextItems(locale)

	if t.hasOriginalColumn(t.columns.typeColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = dataRow.GetData(t.columns.typeColumnName)

		if dataRow.GetData(t.columns.typeColumnName) == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
			if statusName == alipayTransactionDataStatusClosedName {
				log.Warnf(ctx, "[wechat_pay_transaction_csv_data_table.parseTransactionData] skip parsing transaction in row \"%s\", because income transaction is closed", rowId)
				return nil, false, nil
			}

			if statusName == alipayTransactionDataStatusSuccessName {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = ""
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			}
		} else if dataRow.GetData(t.columns.typeColumnName) == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
			if statusName == alipayTransactionDataStatusClosedName {
				log.Warnf(ctx, "[wechat_pay_transaction_csv_data_table.parseTransactionData] skip parsing transaction in row \"%s\", because non-income/expense transaction is closed", rowId)
				return nil, false, nil
			}

			targetName := ""
			productName := ""

			if t.hasOriginalColumn(t.columns.targetNameColumnName) {
				targetName = dataRow.GetData(t.columns.targetNameColumnName)
			}

			if t.hasOriginalColumn(t.columns.productNameColumnName) {
				productName = dataRow.GetData(t.columns.productNameColumnName)
			}

			if statusName == alipayTransactionDataStatusRefundSuccessName {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				if strings.Index(productName, alipayTransactionDataProductNameTransferToAlipayPrefix) == 0 { // transfer to alipay wallet
					data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = ""
					data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
				} else if strings.Index(productName, alipayTransactionDataProductNameTransferFromAlipayPrefix) == 0 { // transfer from alipay wallet
					data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
					data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else if strings.Index(productName, alipayTransactionDataProductNameTransferInText) >= 0 { // transfer in
					data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
					data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else if strings.Index(productName, alipayTransactionDataProductNameTransferOutText) >= 0 { // transfer out
					data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
					data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else if strings.Index(productName, alipayTransactionDataProductNameRepaymentText) >= 0 { // repayment
					data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
					data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else {
					log.Warnf(ctx, "[wechat_pay_transaction_csv_data_table.parseTransactionData] skip parsing transaction in row \"%s\", because product name (\"%s\") is unknown", rowId, productName)
					return nil, false, nil
				}
			}
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
		}
	}

	if data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] && statusName != "" {
		if statusName == alipayTransactionDataStatusRefundSuccessName || statusName == alipayTransactionDataStatusTaxRefundSuccessName {
			amount, err := utils.ParseAmount(data[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

			if err == nil {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
				data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
			}
		}
	}

	return data, true, nil
}

func createNewAlipayTransactionDataTable(ctx core.Context, reader io.Reader, fileHeaderLine string, dataHeaderStartContent string, dataBottomEndLineRune rune, originalColumnNames alipayTransactionColumnNames) (*alipayTransactionDataTable, error) {
	dataTable, err := createNewAlipayImportedDataTable(ctx, reader, fileHeaderLine, dataHeaderStartContent, dataBottomEndLineRune)

	if err != nil {
		return nil, err
	}

	commonDataTable := datatable.CreateNewImportedCommonDataTable(dataTable)

	if !commonDataTable.HasColumn(originalColumnNames.timeColumnName) ||
		!commonDataTable.HasColumn(originalColumnNames.amountColumnName) ||
		!commonDataTable.HasColumn(originalColumnNames.typeColumnName) ||
		!commonDataTable.HasColumn(originalColumnNames.statusColumnName) {
		log.Errorf(ctx, "[alipay_transaction_csv_data_table.createNewAlipayTransactionDataTable] cannot parse alipay csv data, because missing essential columns in header row")
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	return &alipayTransactionDataTable{
		innerDataTable: commonDataTable,
		columns:        originalColumnNames,
	}, nil
}

func createNewAlipayImportedDataTable(ctx core.Context, reader io.Reader, fileHeaderLine string, dataHeaderStartContent string, dataBottomEndLineRune rune) (datatable.ImportedDataTable, error) {
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
