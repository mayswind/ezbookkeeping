package fireflyIII

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const fireflyIIITransactionTimeColumnName = "date"
const fireflyIIITransactionTypeColumnName = "type"
const fireflyIIITransactionCategoryColumnName = "category"
const fireflyIIITransactionSourceAccountColumnName = "source_name"
const fireflyIIITransactionCurrencyCodeColumnName = "currency_code"
const fireflyIIITransactionAmountColumnName = "amount"
const fireflyIIITransactionDestinationAccountColumnName = "destination_name"
const fireflyIIITransactionForeignCurrencyCodeColumnName = "foreign_currency_code"
const fireflyIIITransactionForeignAmountColumnName = "foreign_amount"
const fireflyIIITransactionTagsColumnName = "tags"
const fireflyIIITransactionDescriptionColumnName = "description"

// fireflyIIITransactionDataRowParser defines the structure of firefly III transaction data row parser
type fireflyIIITransactionDataRowParser struct {
}

// Parse returns the converted transaction data row
func (p *fireflyIIITransactionDataRowParser) Parse(ctx core.Context, user *models.User, dataRow datatable.CommonDataTableRow, rowId string) (rowData map[datatable.TransactionDataTableColumn]string, rowDataValid bool, err error) {
	rowData = make(map[datatable.TransactionDataTableColumn]string, len(fireflyIIITransactionSupportedColumns))

	// parse long date time and timezone
	dateTime, err := utils.ParseFromLongDateTimeWithTimezoneRFC3339Format(dataRow.GetData(fireflyIIITransactionTimeColumnName))

	if err != nil {
		return nil, false, errs.ErrTransactionTimeInvalid
	}

	rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(dateTime.Unix(), dateTime.Location())
	rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(dateTime.Location())

	// parse transaction type
	rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = dataRow.GetData(fireflyIIITransactionTypeColumnName)

	// parse transaction category
	rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(fireflyIIITransactionCategoryColumnName)

	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
		// if the category is empty, use the source account (revenue account) name as the category name
		if rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == "" {
			rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(fireflyIIITransactionSourceAccountColumnName)
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionDestinationAccountColumnName)
	} else if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] {
		// if the category is empty, use the destination account (expense account) name as the category name
		if rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == "" {
			rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(fireflyIIITransactionDestinationAccountColumnName)
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionSourceAccountColumnName)
	} else if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER] {
		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionSourceAccountColumnName)
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionDestinationAccountColumnName)
	} else if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE] {
		rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionDestinationAccountColumnName)
	} else {
		return nil, false, errs.ErrTransactionTypeInvalid
	}

	// parse amount
	rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = dataRow.GetData(fireflyIIITransactionAmountColumnName)
	rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = dataRow.GetData(fireflyIIITransactionForeignAmountColumnName)

	// trim trailing zero in decimal
	if rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] != "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.TrimTrailingZerosInDecimal(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
	}

	if rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] != "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.TrimTrailingZerosInDecimal(rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT])
		amount, err := utils.ParseAmount(rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT])

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(amount)
	} else {
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT]
	}

	// parse account currency
	rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = dataRow.GetData(fireflyIIITransactionCurrencyCodeColumnName)
	rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = dataRow.GetData(fireflyIIITransactionForeignCurrencyCodeColumnName)

	// the related account currency field is foreign currency in firefly III actually
	if rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] == "" && rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] != "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY]
	}

	// parse tags / description
	rowData[datatable.TRANSACTION_DATA_TABLE_TAGS] = dataRow.GetData(fireflyIIITransactionTagsColumnName)
	rowData[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataRow.GetData(fireflyIIITransactionDescriptionColumnName)

	return rowData, true, nil
}

// createFireflyIIITransactionDataRowParser returns firefly III transaction data row parser
func createFireflyIIITransactionDataRowParser() datatable.CommonTransactionDataRowParser {
	return &fireflyIIITransactionDataRowParser{}
}
