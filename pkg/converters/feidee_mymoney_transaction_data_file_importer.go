package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var feideeMymoneyDataColumnNameMapping = map[DataTableColumn]string{
	DATA_TABLE_TRANSACTION_TIME:     "日期",
	DATA_TABLE_TRANSACTION_TYPE:     "交易类型",
	DATA_TABLE_CATEGORY:             "分类",
	DATA_TABLE_SUB_CATEGORY:         "子分类",
	DATA_TABLE_ACCOUNT_NAME:         "账户1",
	DATA_TABLE_AMOUNT:               "金额",
	DATA_TABLE_RELATED_ACCOUNT_NAME: "账户2",
	DATA_TABLE_DESCRIPTION:          "备注",
}

var feideeMymoneyTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: "余额变更",
	models.TRANSACTION_TYPE_INCOME:         "收入",
	models.TRANSACTION_TYPE_EXPENSE:        "支出",
	models.TRANSACTION_TYPE_TRANSFER:       "转账",
}

func feideeMymoneyTransactionDataImporterPostProcess(ctx core.Context, transaction *models.ImportTransaction) error {
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.Amount >= 0 {
			transaction.Type = models.TRANSACTION_DB_TYPE_INCOME
		} else if transaction.Amount < 0 {
			transaction.Amount = -transaction.Amount
			transaction.Type = models.TRANSACTION_DB_TYPE_EXPENSE
		}
	}

	return nil
}
