package wechat

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

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

// wechatPayTransactionDataTable defines the structure of wechatPay transaction plain text data table
type wechatPayTransactionDataTable struct {
	allOriginalLines                  [][]string
	originalHeaderLineColumnNames     []string
	originalTimeColumnIndex           int
	originalCategoryColumnIndex       int
	originalTargetNameColumnIndex     int
	originalProductNameColumnIndex    int
	originalTypeColumnIndex           int
	originalAmountColumnIndex         int
	originalRelatedAccountColumnIndex int
	originalStatusColumnIndex         int
	originalDescriptionColumnIndex    int
}

// wechatPayTransactionDataRow defines the structure of wechatPay transaction plain text data row
type wechatPayTransactionDataRow struct {
	dataTable     *wechatPayTransactionDataTable
	isValid       bool
	originalItems []string
	finalItems    map[datatable.TransactionDataTableColumn]string
}

// wechatPayTransactionDataRowIterator defines the structure of wechatPay transaction plain text data row iterator
type wechatPayTransactionDataRowIterator struct {
	dataTable    *wechatPayTransactionDataTable
	currentIndex int
}

// HasColumn returns whether the transaction data table has specified column
func (t *wechatPayTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := wechatPayTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *wechatPayTransactionDataTable) TransactionRowCount() int {
	if len(t.allOriginalLines) < 1 {
		return 0
	}

	return len(t.allOriginalLines) - 1
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *wechatPayTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &wechatPayTransactionDataRowIterator{
		dataTable:    t,
		currentIndex: 0,
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
	return t.currentIndex+1 < len(t.dataTable.allOriginalLines)
}

// Next returns the next imported data row
func (t *wechatPayTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	if t.currentIndex+1 >= len(t.dataTable.allOriginalLines) {
		return nil, nil
	}

	t.currentIndex++

	rowItems := t.dataTable.allOriginalLines[t.currentIndex]
	isValid := true

	if t.dataTable.originalTypeColumnIndex >= 0 &&
		rowItems[t.dataTable.originalTypeColumnIndex] != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
		rowItems[t.dataTable.originalTypeColumnIndex] != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
		rowItems[t.dataTable.originalTypeColumnIndex] != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		log.Warnf(ctx, "[wechat_pay_transaction_data_plain_text_data_table.Next] skip parsing transaction in row \"index:%d\", because type is \"%s\"", t.currentIndex, rowItems[t.dataTable.originalTypeColumnIndex])
		isValid = false
	}

	var finalItems map[datatable.TransactionDataTableColumn]string
	var errMsg string

	if isValid {
		finalItems, errMsg = t.dataTable.parseTransactionData(ctx, user, rowItems)

		if finalItems == nil {
			log.Warnf(ctx, "[wechat_pay_transaction_data_plain_text_data_table.Next] skip parsing transaction in row \"index:%d\", because %s", t.currentIndex, errMsg)
			isValid = false
		}
	}

	return &wechatPayTransactionDataRow{
		dataTable:     t.dataTable,
		isValid:       isValid,
		originalItems: rowItems,
		finalItems:    finalItems,
	}, nil
}

func (t *wechatPayTransactionDataTable) parseTransactionData(ctx core.Context, user *models.User, items []string) (map[datatable.TransactionDataTableColumn]string, string) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(wechatPayTransactionSupportedColumns))

	if t.originalTimeColumnIndex >= 0 && t.originalTimeColumnIndex < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = items[t.originalTimeColumnIndex]
	}

	if t.originalCategoryColumnIndex >= 0 && t.originalCategoryColumnIndex < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = items[t.originalCategoryColumnIndex]
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
	}

	if t.originalAmountColumnIndex >= 0 && t.originalAmountColumnIndex < len(items) {
		amount, success := utils.ParseFirstConsecutiveNumber(items[t.originalAmountColumnIndex])

		if success {
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = amount
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = items[t.originalAmountColumnIndex]
		}
	}

	if t.originalDescriptionColumnIndex >= 0 && t.originalDescriptionColumnIndex < len(items) && items[t.originalDescriptionColumnIndex] != "" && items[t.originalDescriptionColumnIndex] != "/" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = items[t.originalDescriptionColumnIndex]
	} else if t.originalProductNameColumnIndex >= 0 && t.originalProductNameColumnIndex < len(items) && items[t.originalProductNameColumnIndex] != "" && items[t.originalProductNameColumnIndex] != "/" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = items[t.originalProductNameColumnIndex]
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	relatedAccountName := ""

	if t.originalRelatedAccountColumnIndex >= 0 && t.originalRelatedAccountColumnIndex < len(items) {
		relatedAccountName = items[t.originalRelatedAccountColumnIndex]
	}

	statusName := ""

	if t.originalStatusColumnIndex >= 0 && t.originalStatusColumnIndex < len(items) {
		statusName = items[t.originalStatusColumnIndex]
	}

	locale := user.Language

	if locale == "" {
		locale = ctx.GetClientLocale()
	}

	localeTextItems := locales.GetLocaleTextItems(locale)

	if t.originalTypeColumnIndex >= 0 && t.originalTypeColumnIndex < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = items[t.originalTypeColumnIndex]

		if items[t.originalTypeColumnIndex] == wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
			if relatedAccountName == "" || relatedAccountName == "/" {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.WeChatWallet
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
			}
		} else if items[t.originalTypeColumnIndex] == wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
			if data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == wechatPayTransactionDataCategoryTransferToWeChatWallet {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.WeChatWallet
			} else if data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == wechatPayTransactionDataCategoryTransferFromWeChatWallet {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.WeChatWallet
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = relatedAccountName
			} else {
				return nil, fmt.Sprintf("unkown transfer transaction category")
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

	return data, ""
}

func createNewWeChatPayTransactionDataTable(ctx core.Context, reader io.Reader) (*wechatPayTransactionDataTable, error) {
	allOriginalLines, err := parseAllLinesFromWechatPayTransactionPlainText(ctx, reader)

	if err != nil {
		return nil, err
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[wechat_pay_transaction_data_plain_text_data_table.createNewwechatPayTransactionPlainTextDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	originalHeaderItems := allOriginalLines[0]
	originalHeaderItemMap := make(map[string]int)

	for i := 0; i < len(originalHeaderItems); i++ {
		originalHeaderItemMap[originalHeaderItems[i]] = i
	}

	timeColumnIdx, timeColumnExists := originalHeaderItemMap["交易时间"]
	categoryColumnIdx, categoryColumnExists := originalHeaderItemMap["交易类型"]
	targetNameColumnIdx, targetNameColumnExists := originalHeaderItemMap["交易对方"]
	productNameColumnIdx, productNameColumnExists := originalHeaderItemMap["商品"]
	typeColumnIdx, typeColumnExists := originalHeaderItemMap["收/支"]
	amountColumnIdx, amountColumnExists := originalHeaderItemMap["金额(元)"]
	relatedAccountColumnIdx, relatedAccountColumnExists := originalHeaderItemMap["支付方式"]
	statusColumnIdx, statusColumnExists := originalHeaderItemMap["当前状态"]
	descriptionColumnIdx, descriptionColumnExists := originalHeaderItemMap["备注"]

	if !timeColumnExists || !categoryColumnExists || !typeColumnExists || !amountColumnExists || !statusColumnExists {
		log.Errorf(ctx, "[wechat_pay_transaction_data_plain_text_data_table.createNewwechatPayTransactionPlainTextDataTable] cannot parse wechat pay csv data, because missing essential columns in header row")
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if !targetNameColumnExists {
		targetNameColumnIdx = -1
	}

	if !productNameColumnExists {
		productNameColumnIdx = -1
	}

	if !relatedAccountColumnExists {
		relatedAccountColumnIdx = -1
	}

	if !descriptionColumnExists {
		descriptionColumnIdx = -1
	}

	return &wechatPayTransactionDataTable{
		allOriginalLines:                  allOriginalLines,
		originalHeaderLineColumnNames:     originalHeaderItems,
		originalTimeColumnIndex:           timeColumnIdx,
		originalCategoryColumnIndex:       categoryColumnIdx,
		originalTargetNameColumnIndex:     targetNameColumnIdx,
		originalProductNameColumnIndex:    productNameColumnIdx,
		originalAmountColumnIndex:         amountColumnIdx,
		originalTypeColumnIndex:           typeColumnIdx,
		originalRelatedAccountColumnIndex: relatedAccountColumnIdx,
		originalStatusColumnIndex:         statusColumnIdx,
		originalDescriptionColumnIndex:    descriptionColumnIdx,
	}, nil
}

func parseAllLinesFromWechatPayTransactionPlainText(ctx core.Context, reader io.Reader) ([][]string, error) {
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
			log.Errorf(ctx, "[wechat_pay_transaction_data_plain_text_data_table.parseAllLinesFromWechatPayTransactionPlainText] cannot parse wechat pay csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], wechatPayTransactionDataCsvFileHeader) == 0 || strings.Index(items[0], wechatPayTransactionDataCsvFileHeaderWithUtf8Bom) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[wechat_pay_transaction_data_plain_text_data_table.parseAllLinesFromWechatPayTransactionPlainText] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
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
				log.Errorf(ctx, "[wechat_pay_transaction_data_plain_text_data_table.parseAllLinesFromWechatPayTransactionPlainText] cannot parse row \"index:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", len(allOriginalLines), len(items), len(allOriginalLines[0]))
				return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
			}

			allOriginalLines = append(allOriginalLines, items)
		}
	}

	if !hasFileHeader || !foundContentBeforeDataHeaderLine {
		return nil, errs.ErrInvalidFileHeader
	}

	return allOriginalLines, nil
}
