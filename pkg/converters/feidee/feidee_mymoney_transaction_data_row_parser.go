package feidee

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var feideeMymoneyTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: "余额变更",
	models.TRANSACTION_TYPE_INCOME:         "收入",
	models.TRANSACTION_TYPE_EXPENSE:        "支出",
	models.TRANSACTION_TYPE_TRANSFER:       "转账",
}

var feideeMymoneyTransactionTypeModifyOutstandingBalanceName = "负债变更"

// feideeMymoneyTransactionDataRowParser defines the structure of feidee mymoney transaction data row parser
type feideeMymoneyTransactionDataRowParser struct {
}

// GetAddedColumns returns the added columns after converting the data row
func (p *feideeMymoneyTransactionDataRowParser) GetAddedColumns() []datatable.TransactionDataTableColumn {
	return nil
}

// Parse returns the converted transaction data row
func (p *feideeMymoneyTransactionDataRowParser) Parse(data map[datatable.TransactionDataTableColumn]string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	rowData = make(map[datatable.TransactionDataTableColumn]string, len(data))

	for column, value := range data {
		rowData[column] = value
	}

	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] != "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = p.getLongDateTime(rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME])
	}

	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE] {
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		// balance modification transaction in feidee mymoney app is not the opening balance transaction, it can be added many times
		if amount >= 0 {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
		} else {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
			rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
		}
	} else if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == feideeMymoneyTransactionTypeModifyOutstandingBalanceName {
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		// outstanding balance modification transaction in feidee mymoney app is not the opening balance transaction, it can be added many times
		if amount >= 0 {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
		} else {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
			rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
		}
	}

	return rowData, true, nil
}

// Parse returns the converted transaction data row
func (p *feideeMymoneyTransactionDataRowParser) getLongDateTime(str string) string {
	if utils.IsValidLongDateTimeFormat(str) {
		return str
	}

	if utils.IsValidLongDateTimeWithoutSecondFormat(str) {
		return str + ":00"
	}

	if utils.IsValidLongDateFormat(str) {
		return str + " 00:00:00"
	}

	return str
}

// createFeideeMymoneyTransactionDataRowParser returns feidee mymoney transaction data row parser
func createFeideeMymoneyTransactionDataRowParser() datatable.TransactionDataRowParser {
	return &feideeMymoneyTransactionDataRowParser{}
}
