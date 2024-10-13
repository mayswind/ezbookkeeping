package alipay

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

var alipayTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]any{
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
	allOriginalLines                  [][]string
	originalHeaderLineColumnNames     []string
	originalTimeColumnIndex           int
	originalCategoryColumnIndex       int
	originalTargetNameColumnIndex     int
	originalProductNameColumnIndex    int
	originalAmountColumnIndex         int
	originalTypeColumnIndex           int
	originalRelatedAccountColumnIndex int
	originalStatusColumnIndex         int
	originalDescriptionColumnIndex    int
}

// alipayTransactionDataRow defines the structure of alipay transaction plain text data row
type alipayTransactionDataRow struct {
	dataTable     *alipayTransactionDataTable
	isValid       bool
	originalItems []string
	finalItems    map[datatable.TransactionDataTableColumn]string
}

// alipayTransactionDataRowIterator defines the structure of alipay transaction plain text data row iterator
type alipayTransactionDataRowIterator struct {
	dataTable    *alipayTransactionDataTable
	currentIndex int
}

// HasColumn returns whether the transaction data table has specified column
func (t *alipayTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := alipayTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *alipayTransactionDataTable) TransactionRowCount() int {
	if len(t.allOriginalLines) < 1 {
		return 0
	}

	return len(t.allOriginalLines) - 1
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *alipayTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &alipayTransactionDataRowIterator{
		dataTable:    t,
		currentIndex: 0,
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
	return t.currentIndex+1 < len(t.dataTable.allOriginalLines)
}

// Next returns the next imported data row
func (t *alipayTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	if t.currentIndex+1 >= len(t.dataTable.allOriginalLines) {
		return nil, nil
	}

	t.currentIndex++

	rowItems := t.dataTable.allOriginalLines[t.currentIndex]
	isValid := true

	if t.dataTable.originalTypeColumnIndex >= 0 &&
		rowItems[t.dataTable.originalTypeColumnIndex] != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
		rowItems[t.dataTable.originalTypeColumnIndex] != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
		rowItems[t.dataTable.originalTypeColumnIndex] != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		log.Warnf(ctx, "[alipay_transaction_data_plain_text_data_table.Next] skip parsing transaction in row \"index:%d\", because type is \"%s\"", t.currentIndex, rowItems[t.dataTable.originalTypeColumnIndex])
		isValid = false
	}

	if t.dataTable.originalStatusColumnIndex >= 0 &&
		rowItems[t.dataTable.originalStatusColumnIndex] != alipayTransactionDataStatusSuccessName &&
		rowItems[t.dataTable.originalStatusColumnIndex] != alipayTransactionDataStatusPaymentSuccessName &&
		rowItems[t.dataTable.originalStatusColumnIndex] != alipayTransactionDataStatusRepaymentSuccessName &&
		rowItems[t.dataTable.originalStatusColumnIndex] != alipayTransactionDataStatusClosedName &&
		rowItems[t.dataTable.originalStatusColumnIndex] != alipayTransactionDataStatusRefundSuccessName &&
		rowItems[t.dataTable.originalStatusColumnIndex] != alipayTransactionDataStatusTaxRefundSuccessName {
		log.Warnf(ctx, "[alipay_transaction_data_plain_text_data_table.Next] skip parsing transaction in row \"index:%d\", because status is \"%s\"", t.currentIndex, rowItems[t.dataTable.originalStatusColumnIndex])
		isValid = false
	}

	var finalItems map[datatable.TransactionDataTableColumn]string
	var errMsg string

	if isValid {
		finalItems, errMsg = t.dataTable.parseTransactionData(ctx, user, rowItems)

		if finalItems == nil {
			log.Warnf(ctx, "[alipay_transaction_data_plain_text_data_table.Next] skip parsing transaction in row \"index:%d\", because %s", t.currentIndex, errMsg)
			isValid = false
		}
	}

	return &alipayTransactionDataRow{
		dataTable:     t.dataTable,
		isValid:       isValid,
		originalItems: rowItems,
		finalItems:    finalItems,
	}, nil
}

func (t *alipayTransactionDataTable) parseTransactionData(ctx core.Context, user *models.User, items []string) (map[datatable.TransactionDataTableColumn]string, string) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(alipayTransactionSupportedColumns))

	if t.originalTimeColumnIndex >= 0 && t.originalTimeColumnIndex < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = items[t.originalTimeColumnIndex]
	}

	if t.originalCategoryColumnIndex >= 0 && t.originalCategoryColumnIndex < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = items[t.originalCategoryColumnIndex]
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
	}

	if t.originalAmountColumnIndex >= 0 && t.originalAmountColumnIndex < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = items[t.originalAmountColumnIndex]
	}

	if t.originalDescriptionColumnIndex >= 0 && t.originalDescriptionColumnIndex < len(items) && items[t.originalDescriptionColumnIndex] != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = items[t.originalDescriptionColumnIndex]
	} else if t.originalProductNameColumnIndex >= 0 && t.originalProductNameColumnIndex < len(items) && items[t.originalProductNameColumnIndex] != "" {
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

		if items[t.originalTypeColumnIndex] == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
			if statusName == alipayTransactionDataStatusClosedName {
				return nil, fmt.Sprintf("income transaction is closed")
			}

			if statusName == alipayTransactionDataStatusSuccessName {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = ""
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			}
		} else if items[t.originalTypeColumnIndex] == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
			if statusName == alipayTransactionDataStatusClosedName {
				return nil, fmt.Sprintf("non-income/expense transaction is closed")
			}

			targetName := ""
			productName := ""

			if t.originalTargetNameColumnIndex >= 0 && t.originalTargetNameColumnIndex < len(items) {
				targetName = items[t.originalTargetNameColumnIndex]
			}

			if t.originalProductNameColumnIndex >= 0 && t.originalProductNameColumnIndex < len(items) {
				productName = items[t.originalProductNameColumnIndex]
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
					return nil, fmt.Sprintf("product name (\"%s\") is unknown", productName)
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

	return data, ""
}

func createNewAlipayTransactionDataTable(ctx core.Context, reader io.Reader, fileHeaderLine string, dataHeaderStartContent string, dataBottomEndLineRune rune, originalColumnNames alipayTransactionColumnNames) (*alipayTransactionDataTable, error) {
	allOriginalLines, err := parseAllLinesFromAlipayTransactionPlainText(ctx, reader, fileHeaderLine, dataHeaderStartContent, dataBottomEndLineRune)

	if err != nil {
		return nil, err
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[alipay_transaction_data_plain_text_data_table.createNewAlipayTransactionPlainTextDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	originalHeaderItems := allOriginalLines[0]
	originalHeaderItemMap := make(map[string]int)

	for i := 0; i < len(originalHeaderItems); i++ {
		originalHeaderItemMap[originalHeaderItems[i]] = i
	}

	timeColumnIdx, timeColumnExists := originalHeaderItemMap[originalColumnNames.timeColumnName]
	categoryColumnIdx, categoryColumnExists := originalHeaderItemMap[originalColumnNames.categoryColumnName]
	targetNameColumnIdx, targetNameColumnExists := originalHeaderItemMap[originalColumnNames.targetNameColumnName]
	productNameColumnIdx, productNameColumnExists := originalHeaderItemMap[originalColumnNames.productNameColumnName]
	amountColumnIdx, amountColumnExists := originalHeaderItemMap[originalColumnNames.amountColumnName]
	typeColumnIdx, typeColumnExists := originalHeaderItemMap[originalColumnNames.typeColumnName]
	relatedAccountColumnIdx, relatedAccountColumnExists := originalHeaderItemMap[originalColumnNames.relatedAccountColumnName]
	statusColumnIdx, statusColumnExists := originalHeaderItemMap[originalColumnNames.statusColumnName]
	descriptionColumnIdx, descriptionColumnExists := originalHeaderItemMap[originalColumnNames.descriptionColumnName]

	if !timeColumnExists || !amountColumnExists || !typeColumnExists || !statusColumnExists {
		log.Errorf(ctx, "[alipay_transaction_data_plain_text_data_table.createNewAlipayTransactionPlainTextDataTable] cannot parse alipay csv data, because missing essential columns in header row")
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if originalColumnNames.categoryColumnName == "" || !categoryColumnExists {
		categoryColumnIdx = -1
	}

	if originalColumnNames.targetNameColumnName == "" || !targetNameColumnExists {
		targetNameColumnIdx = -1
	}

	if originalColumnNames.productNameColumnName == "" || !productNameColumnExists {
		productNameColumnIdx = -1
	}

	if originalColumnNames.relatedAccountColumnName == "" || !relatedAccountColumnExists {
		relatedAccountColumnIdx = -1
	}

	if originalColumnNames.descriptionColumnName == "" || !descriptionColumnExists {
		descriptionColumnIdx = -1
	}

	return &alipayTransactionDataTable{
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

func parseAllLinesFromAlipayTransactionPlainText(ctx core.Context, reader io.Reader, fileHeaderLine string, dataHeaderStartContent string, dataBottomEndLineRune rune) ([][]string, error) {
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
			log.Errorf(ctx, "[alipay_transaction_data_plain_text_data_table.parseAllLinesFromAlipayTransactionPlainText] cannot parse alipay csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], fileHeaderLine) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[alipay_transaction_data_plain_text_data_table.parseAllLinesFromAlipayTransactionPlainText] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
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
				log.Errorf(ctx, "[alipay_transaction_data_plain_text_data_table.parseAllLinesFromAlipayTransactionPlainText] cannot parse row \"index:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", len(allOriginalLines), len(items), len(allOriginalLines[0]))
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
