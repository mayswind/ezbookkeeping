package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// feideeMymoneyTransactionDataXlsImporter defines the structure of feidee mymoney xls importer for transaction data
type feideeMymoneyTransactionDataXlsImporter struct {
	DataTableTransactionDataImporter
}

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

// Initialize a feidee mymoney transaction data xls file importer singleton instance
var (
	FeideeMymoneyTransactionDataXlsImporter = &feideeMymoneyTransactionDataXlsImporter{
		DataTableTransactionDataImporter{
			dataColumnMapping:      feideeMymoneyDataColumnNameMapping,
			transactionTypeMapping: feideeMymoneyTransactionTypeNameMapping,
		},
	}
)

// ParseImportedData returns the imported data by parsing the feidee mymoney transaction xls data
func (c *feideeMymoneyTransactionDataXlsImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTable, err := createNewFeideeMymoneyTransactionExcelFileDataTable(data)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	return c.parseImportedData(ctx, user, dataTable, defaultTimezoneOffset, accountMap, categoryMap, tagMap)
}
