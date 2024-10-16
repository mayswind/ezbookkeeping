package wechat

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

const wechatPayTransactionDataCsvFileHeader = "微信支付账单明细"
const wechatPayTransactionDataCsvFileHeaderWithUtf8Bom = "\xEF\xBB\xBF" + wechatPayTransactionDataCsvFileHeader
const wechatPayTransactionDataHeaderStartContentBeginning = "----------------------微信支付账单明细列表--------------------"

const wechatPayTransactionTimeColumnName = "交易时间"
const wechatPayTransactionCategoryColumnName = "交易类型"
const wechatPayTransactionProductNameColumnName = "商品"
const wechatPayTransactionTypeColumnName = "收/支"
const wechatPayTransactionAmountColumnName = "金额(元)"
const wechatPayTransactionRelatedAccountColumnName = "支付方式"
const wechatPayTransactionStatusColumnName = "当前状态"
const wechatPayTransactionDescriptionColumnName = "备注"

const wechatPayTransactionDataCategoryTransferToWeChatWallet = "零钱充值"
const wechatPayTransactionDataCategoryTransferFromWeChatWallet = "零钱提现"

const wechatPayTransactionDataStatusRefundName = "退款"

var wechatPayTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

// wechatPayTransactionDataTable defines the structure of wechat pay transaction plain text data table
type wechatPayTransactionDataTable struct {
	innerDataTable datatable.CommonDataTable
}

// wechatPayTransactionDataRow defines the structure of wechat pay transaction plain text data row
type wechatPayTransactionDataRow struct {
	isValid    bool
	finalItems map[datatable.TransactionDataTableColumn]string
}

// wechatPayTransactionDataRowIterator defines the structure of wechat pay transaction plain text data row iterator
type wechatPayTransactionDataRowIterator struct {
	dataTable     *wechatPayTransactionDataTable
	innerIterator datatable.CommonDataRowIterator
}

// HasColumn returns whether the transaction data table has specified column
func (t *wechatPayTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := wechatPayTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *wechatPayTransactionDataTable) TransactionRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *wechatPayTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &wechatPayTransactionDataRowIterator{
		dataTable:     t,
		innerIterator: t.innerDataTable.DataRowIterator(),
	}
}

// IsValid returns whether this row is valid data for importing
func (r *wechatPayTransactionDataRow) IsValid() bool {
	return r.isValid
}

// GetData returns the data in the specified column type
func (r *wechatPayTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := wechatPayTransactionSupportedColumns[column]

	if !exists {
		return ""
	}

	return r.finalItems[column]
}

// HasNext returns whether the iterator does not reach the end
func (t *wechatPayTransactionDataRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// Next returns the next imported data row
func (t *wechatPayTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	importedRow := t.innerIterator.Next()

	if importedRow == nil {
		return nil, nil
	}

	finalItems, isValid, err := t.dataTable.parseTransactionData(ctx, user, importedRow, t.innerIterator.CurrentRowId())

	if err != nil {
		return nil, err
	}

	return &wechatPayTransactionDataRow{
		isValid:    isValid,
		finalItems: finalItems,
	}, nil
}

func (t *wechatPayTransactionDataTable) hasOriginalColumn(columnName string) bool {
	return columnName != "" && t.innerDataTable.HasColumn(columnName)
}

func (t *wechatPayTransactionDataTable) parseTransactionData(ctx core.Context, user *models.User, dataRow datatable.CommonDataRow, rowId string) (map[datatable.TransactionDataTableColumn]string, bool, error) {
	if dataRow.GetData(wechatPayTransactionTypeColumnName) != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
		dataRow.GetData(wechatPayTransactionTypeColumnName) != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
		dataRow.GetData(wechatPayTransactionTypeColumnName) != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		log.Warnf(ctx, "[wechat_pay_transaction_csv_data_table.parseTransactionData] skip parsing transaction in row \"%s\", because type is \"%s\"", rowId, dataRow.GetData(wechatPayTransactionTypeColumnName))
		return nil, false, nil
	}

	data := make(map[datatable.TransactionDataTableColumn]string, len(wechatPayTransactionSupportedColumns))

	if t.hasOriginalColumn(wechatPayTransactionTimeColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = dataRow.GetData(wechatPayTransactionTimeColumnName)
	}

	if t.hasOriginalColumn(wechatPayTransactionCategoryColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(wechatPayTransactionCategoryColumnName)
	}

	if t.hasOriginalColumn(wechatPayTransactionAmountColumnName) {
		amount, success := utils.ParseFirstConsecutiveNumber(dataRow.GetData(wechatPayTransactionAmountColumnName))

		if !success {
			log.Errorf(ctx, "[wechat_pay_transaction_csv_data_table.parseTransactionData] cannot parse amount \"%s\" of transaction in row \"%s\"", dataRow.GetData(wechatPayTransactionAmountColumnName), rowId)
			return nil, false, errs.ErrAmountInvalid
		}

		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = amount
	}

	if t.hasOriginalColumn(wechatPayTransactionDescriptionColumnName) && dataRow.GetData(wechatPayTransactionDescriptionColumnName) != "" && dataRow.GetData(wechatPayTransactionDescriptionColumnName) != "/" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(wechatPayTransactionDescriptionColumnName)
	} else if t.hasOriginalColumn(wechatPayTransactionProductNameColumnName) && dataRow.GetData(wechatPayTransactionProductNameColumnName) != "" && dataRow.GetData(wechatPayTransactionProductNameColumnName) != "/" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(wechatPayTransactionProductNameColumnName)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	relatedAccountName := ""

	if t.hasOriginalColumn(wechatPayTransactionRelatedAccountColumnName) {
		relatedAccountName = dataRow.GetData(wechatPayTransactionRelatedAccountColumnName)
	}

	statusName := ""

	if t.hasOriginalColumn(wechatPayTransactionStatusColumnName) {
		statusName = dataRow.GetData(wechatPayTransactionStatusColumnName)
	}

	locale := user.Language

	if locale == "" {
		locale = ctx.GetClientLocale()
	}

	localeTextItems := locales.GetLocaleTextItems(locale)

	if t.hasOriginalColumn(wechatPayTransactionTypeColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = dataRow.GetData(wechatPayTransactionTypeColumnName)

		if dataRow.GetData(wechatPayTransactionTypeColumnName) == wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
			if relatedAccountName == "" || relatedAccountName == "/" {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.WeChatWallet
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
			}
		} else if dataRow.GetData(wechatPayTransactionTypeColumnName) == wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
			if data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == wechatPayTransactionDataCategoryTransferToWeChatWallet {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.WeChatWallet
			} else if data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == wechatPayTransactionDataCategoryTransferFromWeChatWallet {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.WeChatWallet
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = relatedAccountName
			} else {
				log.Warnf(ctx, "[wechat_pay_transaction_csv_data_table.parseTransactionData] skip parsing transaction in row \"%s\", because unkown transfer transaction category \"%s\"", rowId, data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY])
				return nil, false, nil
			}
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
		}
	}

	if data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] && statusName != "" {
		if strings.Index(statusName, wechatPayTransactionDataStatusRefundName) >= 0 {
			amount, err := utils.ParseAmount(data[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

			if err == nil {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
				data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
			}
		}
	}

	return data, true, nil
}

func createNewWeChatPayTransactionDataTable(ctx core.Context, reader io.Reader) (*wechatPayTransactionDataTable, error) {
	dataTable, err := createNewWeChatPayImportedDataTable(ctx, reader)

	if err != nil {
		return nil, err
	}

	commonDataTable := datatable.CreateNewImportedCommonDataTable(dataTable)

	if !commonDataTable.HasColumn(wechatPayTransactionTimeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionCategoryColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionTypeColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionAmountColumnName) ||
		!commonDataTable.HasColumn(wechatPayTransactionStatusColumnName) {
		log.Errorf(ctx, "[wechat_pay_transaction_csv_data_table.createNewWeChatPayTransactionDataTable] cannot parse wechat pay csv data, because missing essential columns in header row")
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	return &wechatPayTransactionDataTable{
		innerDataTable: commonDataTable,
	}, nil
}

func createNewWeChatPayImportedDataTable(ctx core.Context, reader io.Reader) (datatable.ImportedDataTable, error) {
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
			log.Errorf(ctx, "[wechat_pay_transaction_csv_data_table.createNewWeChatPayImportedDataTable] cannot parse wechat pay csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], wechatPayTransactionDataCsvFileHeader) == 0 || strings.Index(items[0], wechatPayTransactionDataCsvFileHeaderWithUtf8Bom) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[wechat_pay_transaction_csv_data_table.createNewWeChatPayImportedDataTable] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
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
				log.Errorf(ctx, "[wechat_pay_transaction_csv_data_table.createNewWeChatPayImportedDataTable] cannot parse row \"index:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", len(allOriginalLines), len(items), len(allOriginalLines[0]))
				return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
			}

			allOriginalLines = append(allOriginalLines, items)
		}
	}

	if !hasFileHeader || !foundContentBeforeDataHeaderLine {
		return nil, errs.ErrInvalidFileHeader
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[wechat_pay_transaction_csv_data_table.createNewWeChatPayImportedDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	dataTable := csvdatatable.CreateNewCustomCsvImportedDataTable(allOriginalLines)

	return dataTable, nil
}
