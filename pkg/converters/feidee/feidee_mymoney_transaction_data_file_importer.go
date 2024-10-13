package feidee

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var feideeMymoneyDataColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     "日期",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     "交易类型",
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:             "分类",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         "子分类",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         "账户1",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               "金额",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: "账户2",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          "备注",
}

var feideeMymoneyTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: "余额变更",
	models.TRANSACTION_TYPE_INCOME:         "收入",
	models.TRANSACTION_TYPE_EXPENSE:        "支出",
	models.TRANSACTION_TYPE_TRANSFER:       "转账",
}
