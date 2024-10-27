package wechat

import (
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

// weChatPayTransactionDataRowParser defines the structure of wechat pay transaction data row parser
type weChatPayTransactionDataRowParser struct {
}

// Parse returns the converted transaction data row
func (t *weChatPayTransactionDataRowParser) Parse(ctx core.Context, user *models.User, dataTable *datatable.CommonTransactionDataTable, dataRow datatable.CommonDataRow, rowId string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	if dataRow.GetData(wechatPayTransactionTypeColumnName) != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
		dataRow.GetData(wechatPayTransactionTypeColumnName) != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
		dataRow.GetData(wechatPayTransactionTypeColumnName) != wechatPayTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		log.Warnf(ctx, "[wechat_pay_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because type is \"%s\"", rowId, dataRow.GetData(wechatPayTransactionTypeColumnName))
		return nil, false, nil
	}

	data := make(map[datatable.TransactionDataTableColumn]string, len(wechatPayTransactionSupportedColumns))

	if dataTable.HasOriginalColumn(wechatPayTransactionTimeColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = dataRow.GetData(wechatPayTransactionTimeColumnName)
	}

	if dataTable.HasOriginalColumn(wechatPayTransactionCategoryColumnName) {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(wechatPayTransactionCategoryColumnName)
	}

	if dataTable.HasOriginalColumn(wechatPayTransactionAmountColumnName) {
		amount, success := utils.ParseFirstConsecutiveNumber(dataRow.GetData(wechatPayTransactionAmountColumnName))

		if !success {
			log.Errorf(ctx, "[wechat_pay_transaction_data_row_parser.Parse] cannot parse amount \"%s\" of transaction in row \"%s\"", dataRow.GetData(wechatPayTransactionAmountColumnName), rowId)
			return nil, false, errs.ErrAmountInvalid
		}

		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = amount
	}

	if dataTable.HasOriginalColumn(wechatPayTransactionDescriptionColumnName) && dataRow.GetData(wechatPayTransactionDescriptionColumnName) != "" && dataRow.GetData(wechatPayTransactionDescriptionColumnName) != "/" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(wechatPayTransactionDescriptionColumnName)
	} else if dataTable.HasOriginalColumn(wechatPayTransactionProductNameColumnName) && dataRow.GetData(wechatPayTransactionProductNameColumnName) != "" && dataRow.GetData(wechatPayTransactionProductNameColumnName) != "/" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(wechatPayTransactionProductNameColumnName)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	relatedAccountName := ""

	if dataTable.HasOriginalColumn(wechatPayTransactionRelatedAccountColumnName) {
		relatedAccountName = dataRow.GetData(wechatPayTransactionRelatedAccountColumnName)
	}

	statusName := ""

	if dataTable.HasOriginalColumn(wechatPayTransactionStatusColumnName) {
		statusName = dataRow.GetData(wechatPayTransactionStatusColumnName)
	}

	locale := user.Language

	if locale == "" {
		locale = ctx.GetClientLocale()
	}

	localeTextItems := locales.GetLocaleTextItems(locale)

	if dataTable.HasOriginalColumn(wechatPayTransactionTypeColumnName) {
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
				log.Warnf(ctx, "[wechat_pay_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because unkown transfer transaction category \"%s\"", rowId, data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY])
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

// createWeChatPayTransactionDataRowParser returns wechat pay transaction data row parser
func createWeChatPayTransactionDataRowParser() datatable.CommonTransactionDataRowParser {
	return &weChatPayTransactionDataRowParser{}
}