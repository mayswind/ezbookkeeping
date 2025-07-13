package fireflyIII

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const fireflyIIITransactionTimeColumnName = "date"
const fireflyIIITransactionTypeColumnName = "type"
const fireflyIIITransactionCategoryColumnName = "category"
const fireflyIIITransactionSourceAccountNameColumnName = "source_name"
const fireflyIIITransactionSourceAccountTypeColumnName = "source_type"
const fireflyIIITransactionCurrencyCodeColumnName = "currency_code"
const fireflyIIITransactionAmountColumnName = "amount"
const fireflyIIITransactionDestinationAccountNameColumnName = "destination_name"
const fireflyIIITransactionDestinationAccountTypeColumnName = "destination_type"
const fireflyIIITransactionForeignCurrencyCodeColumnName = "foreign_currency_code"
const fireflyIIITransactionForeignAmountColumnName = "foreign_amount"
const fireflyIIITransactionTagsColumnName = "tags"
const fireflyIIITransactionDescriptionColumnName = "description"

const fireflyIIIAssetAccountName = "Asset account"
const fireflyIIIExpenseAccountName = "Expense account"
const fireflyIIIRevenueAccountName = "Revenue account"
const fireflyIIIDebtAccountName = "Debt"

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

	// parse transaction type, transaction category and amount
	transactionType := dataRow.GetData(fireflyIIITransactionTypeColumnName)
	sourceAccountType := dataRow.GetData(fireflyIIITransactionSourceAccountTypeColumnName)
	destinationAccountType := dataRow.GetData(fireflyIIITransactionDestinationAccountTypeColumnName)

	amount, err := utils.ParseAmount(utils.TrimTrailingZerosInDecimal(dataRow.GetData(fireflyIIITransactionAmountColumnName)))

	if err != nil {
		return nil, false, errs.ErrAmountInvalid
	}

	foreignAmount := amount

	if dataRow.HasData(fireflyIIITransactionForeignAmountColumnName) && dataRow.GetData(fireflyIIITransactionForeignAmountColumnName) != "" {
		foreignAmount, err = utils.ParseAmount(utils.TrimTrailingZerosInDecimal(dataRow.GetData(fireflyIIITransactionForeignAmountColumnName)))

		if err != nil {
			return nil, false, errs.ErrAmountInvalid
		}
	}

	rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(fireflyIIITransactionCategoryColumnName)

	if sourceAccountType == fireflyIIIRevenueAccountName && (destinationAccountType == fireflyIIIAssetAccountName || destinationAccountType == fireflyIIIDebtAccountName) { // income
		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]

		// if the category is empty, use the source account (revenue account) name as the category name
		if rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == "" {
			rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(fireflyIIITransactionSourceAccountNameColumnName)
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionDestinationAccountNameColumnName)
		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
	} else if (sourceAccountType == fireflyIIIAssetAccountName || sourceAccountType == fireflyIIIDebtAccountName) && destinationAccountType == fireflyIIIExpenseAccountName { // expense
		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]

		// if the category is empty, use the destination account (expense account) name as the category name
		if rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == "" {
			rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = dataRow.GetData(fireflyIIITransactionDestinationAccountNameColumnName)
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionSourceAccountNameColumnName)
		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
	} else if transactionType == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE] { // opening balance
		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]
		rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionDestinationAccountNameColumnName)
		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
	} else if (sourceAccountType == fireflyIIIAssetAccountName || sourceAccountType == fireflyIIIDebtAccountName) && (destinationAccountType == fireflyIIIAssetAccountName || destinationAccountType == fireflyIIIDebtAccountName) { // transfer
		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]
		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionSourceAccountNameColumnName)
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = dataRow.GetData(fireflyIIITransactionDestinationAccountNameColumnName)

		if transactionType == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE] {
			rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
			rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(-foreignAmount)
		} else {
			rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
			rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(foreignAmount)
		}
	} else {
		log.Errorf(ctx, "[fireflyiii_transaction_data_row_parser.Parse] cannot detect transaction type, source account type is \"%s\", destination account type is \"%s\", Firefly III transaction type is \"%s\"", sourceAccountType, destinationAccountType, transactionType)
		return nil, false, errs.ErrTransactionTypeInvalid
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
