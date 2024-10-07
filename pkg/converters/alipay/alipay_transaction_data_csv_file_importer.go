package alipay

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

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

const alipayTransactionDataProductNameRechargePrefix = "充值-"
const alipayTransactionDataProductNameCashWithdrawalPrefix = "提现-"
const alipayTransactionDataProductNameTransferInText = "转入"
const alipayTransactionDataProductNameTransferOutText = "转出"
const alipayTransactionDataProductNameRepaymentText = "还款"

var alipayTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   "收入",
	models.TRANSACTION_TYPE_EXPENSE:  "支出",
	models.TRANSACTION_TYPE_TRANSFER: "不计收支",
}

// alipayTransactionDataCsvImporter defines the structure of alipay csv importer for transaction data
type alipayTransactionDataCsvImporter struct {
	fileHeaderLine           string
	dataHeaderStartContent   string
	dataBottomEndLineRune    rune
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

// ParseImportedData returns the imported data by parsing the alipay transaction csv data
func (c *alipayTransactionDataCsvImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]*models.TransactionCategory, incomeCategoryMap map[string]*models.TransactionCategory, transferCategoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	enc := simplifiedchinese.GB18030
	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())
	allLines, err := c.parseAllLinesFromCsvData(ctx, reader)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if len(allLines) <= 1 {
		log.Errorf(ctx, "[alipayTransactionDataCsvImporter.ParseImportedData] cannot parse import data for user \"uid:%d\", because data table row count is less 1", user.Uid)
		return nil, nil, nil, nil, nil, nil, errs.ErrNotFoundTransactionDataInFile
	}

	headerLineItems := allLines[0]
	headerItemMap := make(map[string]int)

	for i := 0; i < len(headerLineItems); i++ {
		headerItemMap[headerLineItems[i]] = i
	}

	timeColumnIdx, timeColumnExists := headerItemMap[c.timeColumnName]
	categoryColumnIdx, categoryColumnExists := headerItemMap[c.categoryColumnName]
	targetNameColumnIdx, targetNameColumnExists := headerItemMap[c.targetNameColumnName]
	productNameColumnIdx, productNameColumnExists := headerItemMap[c.productNameColumnName]
	amountColumnIdx, amountColumnExists := headerItemMap[c.amountColumnName]
	typeColumnIdx, typeColumnExists := headerItemMap[c.typeColumnName]
	relatedAccountColumnIdx, relatedAccountColumnExists := headerItemMap[c.relatedAccountColumnName]
	statusColumnIdx, statusColumnExists := headerItemMap[c.statusColumnName]
	descriptionColumnIdx, descriptionColumnExists := headerItemMap[c.descriptionColumnName]

	if !timeColumnExists || !amountColumnExists || !typeColumnExists || !statusColumnExists {
		log.Errorf(ctx, "[alipayTransactionDataCsvImporter.ParseImportedData] cannot parse import data for user \"uid:%d\", because missing essential columns in header row", user.Uid)
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	newColumns := make([]datatable.DataTableColumn, 0, 7)
	newColumns = append(newColumns, datatable.DATA_TABLE_TRANSACTION_TYPE)
	newColumns = append(newColumns, datatable.DATA_TABLE_TRANSACTION_TIME)
	newColumns = append(newColumns, datatable.DATA_TABLE_SUB_CATEGORY)
	newColumns = append(newColumns, datatable.DATA_TABLE_ACCOUNT_NAME)
	newColumns = append(newColumns, datatable.DATA_TABLE_AMOUNT)
	newColumns = append(newColumns, datatable.DATA_TABLE_RELATED_ACCOUNT_NAME)
	newColumns = append(newColumns, datatable.DATA_TABLE_DESCRIPTION)

	dataTable := datatable.CreateNewWritableDataTable(newColumns)

	for i := 1; i < len(allLines); i++ {
		items := allLines[i]

		if len(items) < len(headerLineItems) {
			log.Errorf(ctx, "[alipayTransactionDataCsvImporter.ParseImportedData] cannot parse row \"index:%d\" for user \"uid:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", i, user.Uid, len(items), len(headerLineItems))
			return nil, nil, nil, nil, nil, nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
		}

		if items[typeColumnIdx] != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
			items[typeColumnIdx] != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
			items[typeColumnIdx] != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
			log.Warnf(ctx, "[alipayTransactionDataCsvImporter.ParseImportedData] skip parsing transaction in row \"index:%d\" for user \"uid:%d\", because type is \"%s\"", i, user.Uid, items[typeColumnIdx])
			continue
		}

		if items[statusColumnIdx] != alipayTransactionDataStatusSuccessName &&
			items[statusColumnIdx] != alipayTransactionDataStatusPaymentSuccessName &&
			items[statusColumnIdx] != alipayTransactionDataStatusRepaymentSuccessName &&
			items[statusColumnIdx] != alipayTransactionDataStatusClosedName &&
			items[statusColumnIdx] != alipayTransactionDataStatusRefundSuccessName &&
			items[statusColumnIdx] != alipayTransactionDataStatusTaxRefundSuccessName {
			log.Warnf(ctx, "[alipayTransactionDataCsvImporter.ParseImportedData] skip parsing transaction in row \"index:%d\" for user \"uid:%d\", because status is \"%s\"", i, user.Uid, items[statusColumnIdx])
			continue
		}

		data, errMsg := c.parseTransactionData(ctx,
			user,
			items,
			timeColumnIdx,
			timeColumnExists,
			categoryColumnIdx,
			categoryColumnExists,
			targetNameColumnIdx,
			targetNameColumnExists,
			productNameColumnIdx,
			productNameColumnExists,
			amountColumnIdx,
			amountColumnExists,
			typeColumnIdx,
			typeColumnExists,
			relatedAccountColumnIdx,
			relatedAccountColumnExists,
			statusColumnIdx,
			statusColumnExists,
			descriptionColumnIdx,
			descriptionColumnExists,
		)

		if data == nil {
			log.Warnf(ctx, "[alipayTransactionDataCsvImporter.ParseImportedData] skip parsing transaction in row \"index:%d\" for user \"uid:%d\", because %s", i, user.Uid, errMsg)
			continue
		}

		dataTable.Add(data)
	}

	dataTableImporter := datatable.CreateNewSimpleImporterFromWritableDataTable(
		dataTable,
		alipayTransactionTypeNameMapping,
	)

	return dataTableImporter.ParseImportedData(ctx, user, dataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (c *alipayTransactionDataCsvImporter) parseAllLinesFromCsvData(ctx core.Context, reader io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	allLines := make([][]string, 0)
	hasFileHeader := false
	foundContentBeforeDataHeaderLine := false

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[alipayTransactionDataCsvImporter.parseAllLinesFromCsvData] cannot parse alipay csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], c.fileHeaderLine) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[alipayTransactionDataCsvImporter.parseAllLinesFromCsvData] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
			}
		}

		if !foundContentBeforeDataHeaderLine {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], c.dataHeaderStartContent) >= 0 {
				foundContentBeforeDataHeaderLine = true
				continue
			} else {
				continue
			}
		}

		if foundContentBeforeDataHeaderLine {
			if len(items) <= 0 {
				continue
			} else if len(items) == 1 && c.dataBottomEndLineRune > 0 && utils.ContainsOnlyOneRune(items[0], c.dataBottomEndLineRune) {
				break
			}

			for i := 0; i < len(items); i++ {
				items[i] = strings.Trim(items[i], " ")
			}

			allLines = append(allLines, items)
		}
	}

	if !hasFileHeader || !foundContentBeforeDataHeaderLine {
		return nil, errs.ErrInvalidFileHeader
	}

	return allLines, nil
}

func (c *alipayTransactionDataCsvImporter) parseTransactionData(
	ctx core.Context,
	user *models.User,
	items []string,
	timeColumnIdx int,
	timeColumnExists bool,
	categoryColumnIdx int,
	categoryColumnExists bool,
	targetNameColumnIdx int,
	targetNameColumnExists bool,
	productNameColumnIdx int,
	productNameColumnExists bool,
	amountColumnIdx int,
	amountColumnExists bool,
	typeColumnIdx int,
	typeColumnExists bool,
	relatedAccountColumnIdx int,
	relatedAccountColumnExists bool,
	statusColumnIdx int,
	statusColumnExists bool,
	descriptionColumnIdx int,
	descriptionColumnExists bool,
) (map[datatable.DataTableColumn]string, string) {
	data := make(map[datatable.DataTableColumn]string, 11)

	if timeColumnExists && timeColumnIdx < len(items) {
		data[datatable.DATA_TABLE_TRANSACTION_TIME] = items[timeColumnIdx]
	}

	if categoryColumnExists && categoryColumnIdx < len(items) {
		data[datatable.DATA_TABLE_SUB_CATEGORY] = items[categoryColumnIdx]
	} else {
		data[datatable.DATA_TABLE_SUB_CATEGORY] = ""
	}

	if amountColumnExists && amountColumnIdx < len(items) {
		data[datatable.DATA_TABLE_AMOUNT] = items[amountColumnIdx]
	}

	if descriptionColumnExists && descriptionColumnIdx < len(items) && items[descriptionColumnIdx] != "" {
		data[datatable.DATA_TABLE_DESCRIPTION] = items[descriptionColumnIdx]
	} else if productNameColumnExists && productNameColumnIdx < len(items) && items[productNameColumnIdx] != "" {
		data[datatable.DATA_TABLE_DESCRIPTION] = items[productNameColumnIdx]
	} else {
		data[datatable.DATA_TABLE_DESCRIPTION] = ""
	}

	relatedAccountName := ""

	if relatedAccountColumnExists && relatedAccountColumnIdx < len(items) {
		relatedAccountName = items[relatedAccountColumnIdx]
	}

	statusName := ""

	if statusColumnExists && statusColumnIdx < len(items) {
		statusName = items[statusColumnIdx]
	}

	locale := user.Language

	if locale == "" {
		locale = ctx.GetClientLocale()
	}

	localeTextItems := locales.GetLocaleTextItems(locale)

	if typeColumnExists && typeColumnIdx < len(items) {
		data[datatable.DATA_TABLE_TRANSACTION_TYPE] = items[typeColumnIdx]

		if items[typeColumnIdx] == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
			if statusName == alipayTransactionDataStatusClosedName {
				return nil, fmt.Sprintf("income transaction is closed")
			}

			if statusName == alipayTransactionDataStatusSuccessName {
				data[datatable.DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
				data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				data[datatable.DATA_TABLE_ACCOUNT_NAME] = ""
				data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			}
		} else if items[typeColumnIdx] == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
			if statusName == alipayTransactionDataStatusClosedName {
				return nil, fmt.Sprintf("non-income/expense transaction is closed")
			}

			targetName := ""
			productName := ""

			if targetNameColumnExists && targetNameColumnIdx < len(items) {
				targetName = items[targetNameColumnIdx]
			}

			if productNameColumnExists && productNameColumnIdx < len(items) {
				productName = items[productNameColumnIdx]
			}

			if statusName == alipayTransactionDataStatusRefundSuccessName {
				data[datatable.DATA_TABLE_TRANSACTION_TYPE] = alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
				data[datatable.DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
				data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				if strings.Index(productName, alipayTransactionDataProductNameRechargePrefix) == 0 { // transfer to alipay wallet
					data[datatable.DATA_TABLE_ACCOUNT_NAME] = ""
					data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
				} else if strings.Index(productName, alipayTransactionDataProductNameCashWithdrawalPrefix) == 0 { // transfer from alipay wallet
					data[datatable.DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
					data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else if strings.Index(productName, alipayTransactionDataProductNameTransferInText) >= 0 { // transfer in
					data[datatable.DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
					data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else if strings.Index(productName, alipayTransactionDataProductNameTransferOutText) >= 0 { // transfer out
					data[datatable.DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
					data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else if strings.Index(productName, alipayTransactionDataProductNameRepaymentText) >= 0 { // repayment
					data[datatable.DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
					data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = targetName
				} else {
					return nil, fmt.Sprintf("product name (\"%s\") is unknown", productName)
				}
			}
		} else {
			data[datatable.DATA_TABLE_ACCOUNT_NAME] = relatedAccountName
			data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
		}
	}

	if data[datatable.DATA_TABLE_TRANSACTION_TYPE] == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] && statusName != "" {
		if statusName == alipayTransactionDataStatusRefundSuccessName || statusName == alipayTransactionDataStatusTaxRefundSuccessName {
			amount, err := utils.ParseAmount(data[datatable.DATA_TABLE_AMOUNT])

			if err == nil {
				data[datatable.DATA_TABLE_TRANSACTION_TYPE] = alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
				data[datatable.DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
			}
		}
	}

	return data, ""
}
