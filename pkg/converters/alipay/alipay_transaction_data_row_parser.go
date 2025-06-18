package alipay

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
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

// alipayTransactionDataRowParser defines the structure of alipay transaction data row parser
type alipayTransactionDataRowParser struct {
	columns                    alipayTransactionColumnNames
	existedOriginalDataColumns map[string]bool
}

// Parse returns the converted transaction data row
func (p *alipayTransactionDataRowParser) Parse(ctx core.Context, user *models.User, dataRow datatable.CommonDataTableRow, rowId string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	if dataRow.GetData(p.columns.typeColumnName) != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
		dataRow.GetData(p.columns.typeColumnName) != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
		dataRow.GetData(p.columns.typeColumnName) != alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		log.Warnf(ctx, "[alipay_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because type is \"%s\"", rowId, dataRow.GetData(p.columns.typeColumnName))
		return nil, false, nil
	}

	if dataRow.GetData(p.columns.statusColumnName) != alipayTransactionDataStatusSuccessName &&
		dataRow.GetData(p.columns.statusColumnName) != alipayTransactionDataStatusPaymentSuccessName &&
		dataRow.GetData(p.columns.statusColumnName) != alipayTransactionDataStatusRepaymentSuccessName &&
		dataRow.GetData(p.columns.statusColumnName) != alipayTransactionDataStatusClosedName &&
		dataRow.GetData(p.columns.statusColumnName) != alipayTransactionDataStatusRefundSuccessName &&
		dataRow.GetData(p.columns.statusColumnName) != alipayTransactionDataStatusTaxRefundSuccessName {
		log.Warnf(ctx, "[alipay_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because status is \"%s\"", rowId, dataRow.GetData(p.columns.statusColumnName))
		return nil, false, nil
	}

	data := make(map[datatable.TransactionDataTableColumn]string, len(alipayTransactionSupportedColumns))

	if p.hasOriginalColumn(p.columns.timeColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = dataRow.GetData(p.columns.timeColumnName)
	}

	if p.hasOriginalColumn(p.columns.categoryColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(p.columns.categoryColumnName)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
	}

	if p.hasOriginalColumn(p.columns.amountColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = dataRow.GetData(p.columns.amountColumnName)
	}

	if p.hasOriginalColumn(p.columns.descriptionColumnName) && dataRow.GetData(p.columns.descriptionColumnName) != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(p.columns.descriptionColumnName)
	} else if p.hasOriginalColumn(p.columns.productNameColumnName) && dataRow.GetData(p.columns.productNameColumnName) != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(p.columns.productNameColumnName)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	relatedAccountName := ""

	if p.hasOriginalColumn(p.columns.relatedAccountColumnName) {
		relatedAccountName = dataRow.GetData(p.columns.relatedAccountColumnName)
	}

	statusName := ""

	if p.hasOriginalColumn(p.columns.statusColumnName) {
		statusName = dataRow.GetData(p.columns.statusColumnName)
	}

	locale := user.Language

	if locale == "" {
		locale = ctx.GetClientLocale()
	}

	localeTextItems := locales.GetLocaleTextItems(locale)

	if p.hasOriginalColumn(p.columns.typeColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = dataRow.GetData(p.columns.typeColumnName)

		if dataRow.GetData(p.columns.typeColumnName) == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
			if statusName == alipayTransactionDataStatusClosedName {
				log.Warnf(ctx, "[alipay_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because income transaction is closed", rowId)
				return nil, false, nil
			}

			if statusName == alipayTransactionDataStatusSuccessName {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = localeTextItems.DataConverterTextItems.Alipay
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			} else {
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = ""
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			}
		} else if dataRow.GetData(p.columns.typeColumnName) == alipayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
			if statusName == alipayTransactionDataStatusClosedName {
				log.Warnf(ctx, "[alipay_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because non-income/expense transaction is closed", rowId)
				return nil, false, nil
			}

			targetName := ""
			productName := ""

			if p.hasOriginalColumn(p.columns.targetNameColumnName) {
				targetName = dataRow.GetData(p.columns.targetNameColumnName)
			}

			if p.hasOriginalColumn(p.columns.productNameColumnName) {
				productName = dataRow.GetData(p.columns.productNameColumnName)
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
					log.Warnf(ctx, "[alipay_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because product name (\"%s\") is unknown", rowId, productName)
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

func (p *alipayTransactionDataRowParser) hasOriginalColumn(columnName string) bool {
	_, exists := p.existedOriginalDataColumns[columnName]
	return exists
}

// createAlipayTransactionDataRowParser returns alipay transaction data row parser
func createAlipayTransactionDataRowParser(originalColumnNames alipayTransactionColumnNames, headerColumnNames []string) datatable.CommonTransactionDataRowParser {
	existedOriginalDataColumns := make(map[string]bool, len(headerColumnNames))

	for i := 0; i < len(headerColumnNames); i++ {
		existedOriginalDataColumns[headerColumnNames[i]] = true
	}

	return &alipayTransactionDataRowParser{
		columns:                    originalColumnNames,
		existedOriginalDataColumns: existedOriginalDataColumns,
	}
}
