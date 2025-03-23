package dsv

import (
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// customPlainTextDataTable defines the structure of custom plain text transaction data table
type customPlainTextDataTable struct {
	innerDataTable             datatable.ImportedDataTable
	columnIndexMapping         map[datatable.TransactionDataTableColumn]int
	transactionTypeNameMapping map[string]models.TransactionType
	timeFormat                 string
	timezoneFormat             string
	timeFormatIncludeTimezone  bool
	amountDecimalSeparator     string
	amountDigitGroupingSymbol  string
}

// customPlainTextDataRow defines the structure of custom plain text transaction data row
type customPlainTextDataRow struct {
	transactionDataTable *customPlainTextDataTable
	rowData              map[datatable.TransactionDataTableColumn]string
	isValid              bool
}

// customPlainTextDataRowIterator defines the structure of custom plain text transaction data row iterator
type customPlainTextDataRowIterator struct {
	transactionDataTable *customPlainTextDataTable
	innerIterator        datatable.ImportedDataRowIterator
}

// HasColumn returns whether the data table has specified column
func (t *customPlainTextDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	// custom dsv file allows no sub category, account name and related account name column mapping, but data table converter needs these columns
	if column == datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY ||
		column == datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME ||
		column == datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME {
		return true
	}

	// timezone column will be added when original time format contains timezone
	if t.timeFormatIncludeTimezone && column == datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE {
		return true
	}

	_, exists := t.columnIndexMapping[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *customPlainTextDataTable) TransactionRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *customPlainTextDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &customPlainTextDataRowIterator{
		transactionDataTable: t,
		innerIterator:        t.innerDataTable.DataRowIterator(),
	}
}

// IsValid returns whether this row is valid data for importing
func (r *customPlainTextDataRow) IsValid() bool {
	return r.isValid
}

// GetData returns the data in the specified column type
func (r *customPlainTextDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	return r.rowData[column]
}

// HasNext returns whether the iterator does not reach the end
func (t *customPlainTextDataRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// Next returns the next transaction data row
func (t *customPlainTextDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	importedRow := t.innerIterator.Next()

	if importedRow == nil {
		return nil, nil
	}

	rowData, isValid, err := t.parseTransaction(ctx, user, importedRow)

	if err != nil {
		log.Errorf(ctx, "[custom_transaction_plain_text_data_table.Next] cannot parsing transaction in row \"%s\", because %s", t.innerIterator.CurrentRowId(), err.Error())
		return nil, err
	}

	return &customPlainTextDataRow{
		transactionDataTable: t.transactionDataTable,
		rowData:              rowData,
		isValid:              isValid,
	}, nil
}

func (t *customPlainTextDataRowIterator) parseTransaction(ctx core.Context, user *models.User, row datatable.ImportedDataRow) (map[datatable.TransactionDataTableColumn]string, bool, error) {
	rowData := make(map[datatable.TransactionDataTableColumn]string, len(t.transactionDataTable.columnIndexMapping))

	for column, columnIndex := range t.transactionDataTable.columnIndexMapping {
		if columnIndex < 0 || columnIndex >= row.ColumnCount() {
			continue
		}

		value := row.GetData(columnIndex)
		rowData[column] = value
	}

	// parse transaction type
	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] != "" {
		transactionType, exists := t.transactionDataTable.transactionTypeNameMapping[rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE]]

		if !exists {
			log.Warnf(ctx, "[custom_transaction_plain_text_data_table.parseTransaction] skip parsing this transaction, because transaction type \"%s\" mapping not defined", rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE])
			return nil, false, nil
		}

		mappedTransactionType, exists := customTransactionTypeNameMapping[transactionType]

		if !exists {
			log.Errorf(ctx, "[custom_transaction_plain_text_data_table.parseTransaction] cannot parsing transaction type \"%s\", because type \"%d\" is invalid", rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE], transactionType)
			return nil, false, errs.ErrTransactionTypeInvalid
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = mappedTransactionType
	}

	// parse date time
	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] != "" {
		dateTime, err := time.Parse(t.transactionDataTable.timeFormat, rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME])

		if err != nil {
			return nil, false, errs.ErrTransactionTimeInvalid
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(dateTime.Unix(), dateTime.Location())

		if t.transactionDataTable.timeFormatIncludeTimezone {
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(dateTime.Location())
		}
	}

	// parse timezone
	if rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] != "" {
		if t.transactionDataTable.timezoneFormat == "Z" || t.transactionDataTable.timezoneFormat == "" { // -HH:mm
			// Do Nothing
		} else if t.transactionDataTable.timezoneFormat == "ZZ" { // -HHmm
			timezone := rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE]

			if len(timezone) != 5 {
				return nil, false, errs.ErrTransactionTimeZoneInvalid
			}

			timezone = timezone[:3] + ":" + timezone[3:]
			rowData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = timezone
		} else {
			return nil, false, errs.ErrImportFileTransactionTimezoneFormatInvalid
		}
	}

	// use primary category if sub category is empty
	if rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] == "" && rowData[datatable.TRANSACTION_DATA_TABLE_CATEGORY] != "" {
		rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = rowData[datatable.TRANSACTION_DATA_TABLE_CATEGORY]
	}

	// trim trailing zero in decimal
	if rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] != "" {
		amount, err := t.parseAmount(ctx, rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT])

		if err != nil {
			log.Errorf(ctx, "[custom_transaction_plain_text_data_table.parseTransaction] cannot parsing transaction amount \"%s\", because %s", rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT], err.Error())
			return nil, false, err
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = amount
	}

	if rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] != "" {
		amount, err := t.parseAmount(ctx, rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT])

		if err != nil {
			log.Errorf(ctx, "[custom_transaction_plain_text_data_table.parseTransaction] cannot parsing transaction related amount \"%s\", because %s", rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT], err.Error())
			return nil, false, err
		}

		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = amount
	}

	if _, exists := rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY]; !exists {
		rowData[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = ""
	}

	if _, exists := rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]; !exists {
		rowData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = ""
	}

	if _, exists := rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME]; !exists {
		rowData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
	}

	return rowData, true, nil
}

func (t *customPlainTextDataRowIterator) parseAmount(ctx core.Context, amountValue string) (string, error) {
	if t.transactionDataTable.amountDigitGroupingSymbol != "" {
		amountValue = strings.ReplaceAll(amountValue, t.transactionDataTable.amountDigitGroupingSymbol, "")
	}

	if t.transactionDataTable.amountDecimalSeparator != "" && t.transactionDataTable.amountDecimalSeparator != "." {
		if strings.Contains(amountValue, ".") {
			return "", errs.ErrAmountInvalid
		}

		amountValue = strings.ReplaceAll(amountValue, t.transactionDataTable.amountDecimalSeparator, ".")
	}

	amountValue = utils.TrimTrailingZerosInDecimal(amountValue)
	amount, err := utils.ParseAmount(amountValue)

	if err != nil {
		return "", errs.ErrAmountInvalid
	}

	return utils.FormatAmount(amount), nil
}

// CreateNewCustomPlainTextDataTable returns transaction data table from imported data table
func CreateNewCustomPlainTextDataTable(dataTable datatable.ImportedDataTable, columnIndexMapping map[datatable.TransactionDataTableColumn]int, transactionTypeNameMapping map[string]models.TransactionType, timeFormat string, timezoneFormat string, amountDecimalSeparator string, amountDigitGroupingSymbol string) *customPlainTextDataTable {
	timeFormatIncludeTimezone := strings.Contains(timeFormat, "z") || strings.Contains(timeFormat, "Z")

	return &customPlainTextDataTable{
		innerDataTable:             dataTable,
		columnIndexMapping:         columnIndexMapping,
		transactionTypeNameMapping: transactionTypeNameMapping,
		timeFormat:                 getDateTimeFormat(timeFormat),
		timezoneFormat:             timezoneFormat,
		timeFormatIncludeTimezone:  timeFormatIncludeTimezone,
		amountDecimalSeparator:     amountDecimalSeparator,
		amountDigitGroupingSymbol:  amountDigitGroupingSymbol,
	}
}

func getDateTimeFormat(format string) string {
	// convert moment.js format to Go format

	format = strings.ReplaceAll(format, "YYYY", "2006")
	format = strings.ReplaceAll(format, "YY", "06")

	format = strings.ReplaceAll(format, "MMMM", "January")
	format = strings.ReplaceAll(format, "MMM", "Jan")
	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "M", "1")

	format = strings.ReplaceAll(format, "DD", "02")
	format = strings.ReplaceAll(format, "D", "2")

	format = strings.ReplaceAll(format, "dddd", "Monday")
	format = strings.ReplaceAll(format, "ddd", "Mon")

	format = strings.ReplaceAll(format, "HH", "15")
	format = strings.ReplaceAll(format, "H", "15")

	format = strings.ReplaceAll(format, "hh", "03")
	format = strings.ReplaceAll(format, "h", "3")

	format = strings.ReplaceAll(format, "mm", "04")
	format = strings.ReplaceAll(format, "m", "4")

	format = strings.ReplaceAll(format, "ss", "05")
	format = strings.ReplaceAll(format, "s", "5")

	for i := 9; i >= 1; i-- {
		format = strings.ReplaceAll(format, "."+strings.Repeat("S", i), "."+strings.Repeat("9", i))
	}

	format = strings.ReplaceAll(format, "A", "PM")
	format = strings.ReplaceAll(format, "a", "pm")

	format = strings.ReplaceAll(format, "zz", "MST")
	format = strings.ReplaceAll(format, "z", "MST")

	if strings.Contains(format, "ZZ") {
		format = strings.ReplaceAll(format, "ZZ", "Z0700")
	} else if strings.Contains(format, "Z") {
		format = strings.ReplaceAll(format, "Z", "Z07:00")
	}

	return format
}
