package feidee

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

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

		if amount >= 0 {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
		} else {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
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

	utcTimezone := time.UTC
	utcTimezoneOffsetMinutes := utils.GetTimezoneOffsetMinutes(utcTimezone)

	if utils.IsValidLongDateTimeWithoutSecondFormat(str) {
		dateTime, err := utils.ParseFromLongDateTimeWithoutSecond(str, utcTimezoneOffsetMinutes)

		if err == nil {
			return utils.FormatUnixTimeToLongDateTime(dateTime.Unix(), utcTimezone)
		}
	}

	if utils.IsValidLongDateFormat(str) {
		dateTime, err := utils.ParseFromLongDateTimeWithoutSecond(str+" 00:00", utcTimezoneOffsetMinutes)

		if err == nil {
			return utils.FormatUnixTimeToLongDateTime(dateTime.Unix(), utcTimezone)
		}
	}

	return str
}

// createFeideeMymoneyTransactionDataRowParser returns feidee mymoney transaction data row parser
func createFeideeMymoneyTransactionDataRowParser() datatable.TransactionDataRowParser {
	return &feideeMymoneyTransactionDataRowParser{}
}
