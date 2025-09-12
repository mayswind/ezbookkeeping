package jdcom

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const jdComFinanceTransactionDataCsvFileHeader = "导出信息："

const jdComFinanceTransactionTimeColumnName = "交易时间"
const jdComFinanceTransactionMerchantNameColumnName = "商户名称"
const jdComFinanceTransactionMemoColumnName = "交易说明"
const jdComFinanceTransactionAmountColumnName = "金额"
const jdComFinanceTransactionRelatedAccountColumnName = "收/付款方式"
const jdComFinanceTransactionStatusColumnName = "交易状态"
const jdComFinanceTransactionTypeColumnName = "收/支"
const jdComFinanceTransactionCategoryColumnName = "交易分类"
const jdComFinanceTransactionDescriptionColumnName = "备注"

const jdComFinanceTransactionAmountRefundAll = "(已全额退款)"

const jdComFinanceTransactionMemoTransferToWalletPrefix = "充值"
const jdComFinanceTransactionMemoTransferFromWalletPrefix = "提现"
const jdComFinanceTransactionMemoTransferInText = "转入"
const jdComFinanceTransactionMemoTransferOutText = "转出"
const jdComFinanceTransactionMemoRepaymentText = "还款"
const jdComFinanceTransactionMemoRefundText = "退款"

const jdComFinanceTransactionDataStatusSuccessName = "交易成功"
const jdComFinanceTransactionDataStatusRefundSuccessName = "退款成功"

var jdComFinanceTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

var jdComFinanceTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   "收入",
	models.TRANSACTION_TYPE_EXPENSE:  "支出",
	models.TRANSACTION_TYPE_TRANSFER: "不计收支",
}

// jdComFinanceTransactionDataRowParser defines the structure of jd.com finance transaction data row parser
type jdComFinanceTransactionDataRowParser struct {
	existedOriginalDataColumns map[string]bool
}

// Parse returns the converted transaction data row
func (p *jdComFinanceTransactionDataRowParser) Parse(ctx core.Context, user *models.User, dataRow datatable.CommonDataTableRow, rowId string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	if dataRow.GetData(jdComFinanceTransactionTypeColumnName) != jdComFinanceTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] &&
		dataRow.GetData(jdComFinanceTransactionTypeColumnName) != jdComFinanceTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] &&
		dataRow.GetData(jdComFinanceTransactionTypeColumnName) != jdComFinanceTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		log.Warnf(ctx, "[jdcom_finance_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because type is \"%s\"", rowId, dataRow.GetData(jdComFinanceTransactionTypeColumnName))
		return nil, false, nil
	}

	statusName := dataRow.GetData(jdComFinanceTransactionStatusColumnName)

	if statusName != jdComFinanceTransactionDataStatusSuccessName &&
		statusName != jdComFinanceTransactionDataStatusRefundSuccessName {
		log.Warnf(ctx, "[jdcom_finance_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because status is \"%s\"", rowId, statusName)
		return nil, false, nil
	}

	data := make(map[datatable.TransactionDataTableColumn]string, len(jdComFinanceTransactionSupportedColumns))
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = dataRow.GetData(jdComFinanceTransactionTimeColumnName)
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = dataRow.GetData(jdComFinanceTransactionTypeColumnName)
	data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(jdComFinanceTransactionCategoryColumnName)
	data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(jdComFinanceTransactionRelatedAccountColumnName)
	data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""

	if strings.Index(dataRow.GetData(jdComFinanceTransactionAmountColumnName), "(") >= 0 {
		// If a transaction includes a refund, the original transaction amount will like "-xx.xx(已全额退款)" or "-xx.xx(已退款yy.yy)", along with another refund transaction
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = strings.Split(dataRow.GetData(jdComFinanceTransactionAmountColumnName), "(")[0]
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = dataRow.GetData(jdComFinanceTransactionAmountColumnName)
	}

	if p.hasOriginalColumn(jdComFinanceTransactionDescriptionColumnName) && dataRow.GetData(jdComFinanceTransactionDescriptionColumnName) != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(jdComFinanceTransactionDescriptionColumnName)
	} else if p.hasOriginalColumn(jdComFinanceTransactionMemoColumnName) && dataRow.GetData(jdComFinanceTransactionMemoColumnName) != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(jdComFinanceTransactionMemoColumnName)
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = ""
	}

	if data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == jdComFinanceTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		memo := dataRow.GetData(jdComFinanceTransactionMemoColumnName)

		if statusName == jdComFinanceTransactionDataStatusRefundSuccessName || strings.Index(memo, jdComFinanceTransactionMemoRefundText) >= 0 { // refund
			amount, err := utils.ParseAmount(data[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

			if err == nil {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = jdComFinanceTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
				data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
			}
		} else if strings.Index(dataRow.GetData(jdComFinanceTransactionAmountColumnName), jdComFinanceTransactionAmountRefundAll) > 0 { // expense transaction (but include a full refund)
			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = jdComFinanceTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
		} else { // transfer
			if strings.Index(memo, jdComFinanceTransactionMemoTransferToWalletPrefix) >= 0 { // transfer to jd.com finance wallet
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = dataRow.GetData(jdComFinanceTransactionMerchantNameColumnName)
			} else if strings.Index(memo, jdComFinanceTransactionMemoTransferFromWalletPrefix) >= 0 { // transfer from jd.com finance wallet
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(jdComFinanceTransactionMerchantNameColumnName)
			} else if strings.Index(memo, jdComFinanceTransactionMemoTransferInText) >= 0 { // transfer in
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = dataRow.GetData(jdComFinanceTransactionMerchantNameColumnName)
			} else if strings.Index(memo, jdComFinanceTransactionMemoTransferOutText) >= 0 { // transfer out
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(jdComFinanceTransactionMerchantNameColumnName)
			} else if strings.Index(memo, jdComFinanceTransactionMemoRepaymentText) >= 0 { // repayment
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = dataRow.GetData(jdComFinanceTransactionMerchantNameColumnName)
			} else {
				log.Warnf(ctx, "[jdcom_finance_transaction_data_row_parser.Parse] skip parsing transaction in row \"%s\", because memo (\"%s\") of this transfer transaction is unknown", rowId, memo)
				return nil, false, nil
			}
		}
	}

	return data, true, nil
}

func (p *jdComFinanceTransactionDataRowParser) hasOriginalColumn(columnName string) bool {
	_, exists := p.existedOriginalDataColumns[columnName]
	return exists
}

// createJDComFinanceTransactionDataRowParser returns jd.com finance transaction data row parser
func createJDComFinanceTransactionDataRowParser(headerColumnNames []string) datatable.CommonTransactionDataRowParser {
	existedOriginalDataColumns := make(map[string]bool, len(headerColumnNames))

	for i := 0; i < len(headerColumnNames); i++ {
		existedOriginalDataColumns[headerColumnNames[i]] = true
	}

	return &jdComFinanceTransactionDataRowParser{
		existedOriginalDataColumns: existedOriginalDataColumns,
	}
}
