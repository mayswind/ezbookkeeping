package feidee

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_MODIFY_BALANCE_NAME = "余额变更"
var FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_OUTSTANDING_MODIFY_BALANCE_NAME = "负债变更"
var FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_INCOME_NAME = "收入"
var FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_EXPENSE_NAME = "支出"

var feideeMymoneyElecloudTransactionTypeNameMapping = map[string]models.TransactionType{
	FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_MODIFY_BALANCE_NAME:             models.TRANSACTION_TYPE_MODIFY_BALANCE,
	FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_OUTSTANDING_MODIFY_BALANCE_NAME: models.TRANSACTION_TYPE_MODIFY_BALANCE,
	FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_INCOME_NAME:                     models.TRANSACTION_TYPE_INCOME,
	FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_EXPENSE_NAME:                    models.TRANSACTION_TYPE_EXPENSE,
	"转账": models.TRANSACTION_TYPE_TRANSFER,
	"借入": models.TRANSACTION_TYPE_TRANSFER,
	"借出": models.TRANSACTION_TYPE_TRANSFER,
	"收债": models.TRANSACTION_TYPE_TRANSFER,
	"还债": models.TRANSACTION_TYPE_TRANSFER,
	"代付": models.TRANSACTION_TYPE_TRANSFER,
	"报销": models.TRANSACTION_TYPE_TRANSFER,
	"退款": models.TRANSACTION_TYPE_EXPENSE,
}

// feideeMymoneyElecloudTransactionDataRowParser defines the structure of feidee mymoney (elecloud) transaction data row parser
type feideeMymoneyElecloudTransactionDataRowParser struct {
}

// GetAddedColumns returns the added columns after converting the data row
func (p *feideeMymoneyElecloudTransactionDataRowParser) GetAddedColumns() []datatable.TransactionDataTableColumn {
	return nil
}

// Parse returns the converted transaction data row
func (p *feideeMymoneyElecloudTransactionDataRowParser) Parse(data map[datatable.TransactionDataTableColumn]string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	rowData = make(map[datatable.TransactionDataTableColumn]string, len(data))

	for column, value := range data {
		rowData[column] = value
	}

	rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = strings.ReplaceAll(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT], ",", "") // remove thousand separator

	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_MODIFY_BALANCE_NAME {
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		// balance modification transaction in feidee mymoney (elecloud) is not the opening balance transaction, it can be added many times
		if amount >= 0 {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_INCOME_NAME
		} else {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_EXPENSE_NAME
			rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
		}
	} else if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_OUTSTANDING_MODIFY_BALANCE_NAME {
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		// outstanding balance modification transaction in feidee mymoney app is not the opening balance transaction, it can be added many times
		if amount >= 0 {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_EXPENSE_NAME
		} else {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = FEIDEE_MYMONEY_ELECLOUD_TRANSACTION_TYPE_INCOME_NAME
			rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
		}
	}

	return rowData, true, nil
}

// createFeideeMymoneyElecloudTransactionDataRowParser returns feidee mymoney (elecloud) transaction data row parser
func createFeideeMymoneyElecloudTransactionDataRowParser() datatable.TransactionDataRowParser {
	return &feideeMymoneyElecloudTransactionDataRowParser{}
}
