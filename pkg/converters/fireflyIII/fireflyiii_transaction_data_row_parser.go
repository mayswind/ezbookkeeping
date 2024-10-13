package fireflyIII

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// fireflyIIITransactionDataRowParser defines the structure of firefly III transaction data row parser
type fireflyIIITransactionDataRowParser struct {
}

// GetAddedColumns returns the added columns after converting the data row
func (p *fireflyIIITransactionDataRowParser) GetAddedColumns() []datatable.TransactionDataTableColumn {
	return []datatable.TransactionDataTableColumn{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE,
	}
}

// Parse returns the converted transaction data row
func (p *fireflyIIITransactionDataRowParser) Parse(data map[datatable.TransactionDataTableColumn]string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	rowData = make(map[datatable.TransactionDataTableColumn]string, len(data))

	rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""

	for column, value := range data {
		rowData[column] = value
	}

	// parse long date time and timezone
	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] != "" {
		dateTime, err := utils.ParseFromLongDateTimeWithTimezone(rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME])

		if err != nil {
			return nil, false, errs.ErrTransactionTimeInvalid
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(dateTime.Unix(), dateTime.Location())
		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(dateTime.Location())
	}

	// trim trailing zero in decimal
	if rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] != "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.TrimTrailingZerosInDecimal(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
	}

	if rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] != "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.TrimTrailingZerosInDecimal(rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT])
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(-amount)
	} else {
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT]
	}

	// the related account currency field is foreign currency in firefly III actually
	if rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] == "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY]
	}

	// the destination account of modify balance transaction in firefly III is the asset account
	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE] {
		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME]
	}

	// the destination account of income transaction in firefly III is the asset account
	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME]
	}

	return rowData, true, nil
}

// createFireflyIIITransactionDataRowParser returns firefly III transaction data row parser
func createFireflyIIITransactionDataRowParser() datatable.TransactionDataRowParser {
	return &fireflyIIITransactionDataRowParser{}
}
