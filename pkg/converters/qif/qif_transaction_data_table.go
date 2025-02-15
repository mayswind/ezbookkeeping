package qif

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const qifOpeningBalancePayeeText = "Opening Balance"

var qifTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:             true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

// qifDateFormatType represents the quicken interchange format (qif) date format type
type qifDateFormatType byte

const (
	qifYearMonthDayDateFormat qifDateFormatType = 0
	qifMonthDayYearDateFormat qifDateFormatType = 1
	qifDayMonthYearDateFormat qifDateFormatType = 2
)

// qifTransactionDataTable defines the structure of quicken interchange format (qif) transaction data table
type qifTransactionDataTable struct {
	dateFormatType qifDateFormatType
	allData        []*qifTransactionData
}

// qifTransactionDataRow defines the structure of quicken interchange format (qif) transaction data row
type qifTransactionDataRow struct {
	dataTable  *qifTransactionDataTable
	data       *qifTransactionData
	finalItems map[datatable.TransactionDataTableColumn]string
}

// qifTransactionDataRowIterator defines the structure of quicken interchange format (qif) transaction data row iterator
type qifTransactionDataRowIterator struct {
	dataTable    *qifTransactionDataTable
	currentIndex int
}

// HasColumn returns whether the transaction data table has specified column
func (t *qifTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := qifTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *qifTransactionDataTable) TransactionRowCount() int {
	return len(t.allData)
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *qifTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &qifTransactionDataRowIterator{
		dataTable:    t,
		currentIndex: -1,
	}
}

// IsValid returns whether this row is valid data for importing
func (r *qifTransactionDataRow) IsValid() bool {
	return true
}

// GetData returns the data in the specified column type
func (r *qifTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := qifTransactionSupportedColumns[column]

	if exists {
		return r.finalItems[column]
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *qifTransactionDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allData)
}

// Next returns the next imported data row
func (t *qifTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	if t.currentIndex+1 >= len(t.dataTable.allData) {
		return nil, nil
	}

	t.currentIndex++

	data := t.dataTable.allData[t.currentIndex]
	rowItems, err := t.parseTransaction(ctx, user, data)

	if err != nil {
		log.Errorf(ctx, "[qif_transaction_data_table.Next] cannot parsing transaction in row#%d, because %s", t.currentIndex, err.Error())
		return nil, err
	}

	return &qifTransactionDataRow{
		dataTable:  t.dataTable,
		data:       data,
		finalItems: rowItems,
	}, nil
}

func (t *qifTransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, qifTransaction *qifTransactionData) (map[datatable.TransactionDataTableColumn]string, error) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(qifTransactionSupportedColumns))

	if qifTransaction.date == "" {
		return nil, errs.ErrMissingTransactionTime
	}

	transactionTime, err := t.parseTransactionTime(ctx, qifTransaction.date)

	if err != nil {
		return nil, err
	}

	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = transactionTime

	if qifTransaction.amount == "" {
		return nil, errs.ErrAmountInvalid
	}

	amount, err := utils.ParseAmount(strings.ReplaceAll(qifTransaction.amount, ",", "")) // trim thousands separator

	if err != nil {
		return nil, errs.ErrAmountInvalid
	}

	if qifTransaction.account != nil {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = qifTransaction.account.name
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = ""
	}

	if len(qifTransaction.category) > 0 && qifTransaction.category[0] == '[' && qifTransaction.category[len(qifTransaction.category)-1] == ']' {
		if qifTransaction.payee == qifOpeningBalancePayeeText { // balance modification
			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = qifTransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
			data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = qifTransaction.category[1 : len(qifTransaction.category)-1]
		} else { // transfer
			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = qifTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]

			if amount >= 0 { // transfer from [account name]
				data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]
				data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = qifTransaction.category[1 : len(qifTransaction.category)-1]
			} else { // transfer to [account name]
				data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = qifTransaction.category[1 : len(qifTransaction.category)-1]
			}
		}
	} else { // income/expense
		if amount >= 0 {
			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = qifTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = qifTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
			data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
		}

		if strings.Index(qifTransaction.category, ":") > 0 { // category:subcategory
			categories := strings.Split(qifTransaction.category, ":")
			data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = categories[0]
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = categories[len(categories)-1]
		} else {
			data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = ""
			data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = qifTransaction.category
		}
	}

	if qifTransaction.memo != "" {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = qifTransaction.memo
	} else if qifTransaction.payee != "" && qifTransaction.payee != qifOpeningBalancePayeeText {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = qifTransaction.payee
	}

	return data, nil
}

func (t *qifTransactionDataRowIterator) parseTransactionTime(ctx core.Context, date string) (string, error) {
	var year, month, day string

	if (t.dataTable.dateFormatType == qifYearMonthDayDateFormat && utils.IsValidYearMonthDayLongOrShortDateFormat(date)) ||
		(t.dataTable.dateFormatType == qifMonthDayYearDateFormat && utils.IsValidMonthDayYearLongOrShortDateFormat(date)) ||
		(t.dataTable.dateFormatType == qifDayMonthYearDateFormat && utils.IsValidDayMonthYearLongOrShortDateFormat(date)) {
		date = strings.ReplaceAll(date, ".", "-")
		date = strings.ReplaceAll(date, "/", "-")
		date = strings.ReplaceAll(date, "'", "-")
		items := strings.Split(date, "-")

		if t.dataTable.dateFormatType == qifYearMonthDayDateFormat {
			year = items[0]
			month = items[1]
			day = items[2]
		} else if t.dataTable.dateFormatType == qifMonthDayYearDateFormat {
			month = items[0]
			day = items[1]
			year = items[2]
		} else if t.dataTable.dateFormatType == qifDayMonthYearDateFormat {
			day = items[0]
			month = items[1]
			year = items[2]
		}
	}

	if year == "" || month == "" || day == "" {
		log.Errorf(ctx, "[qif_transaction_data_table.parseTransactionTime] cannot parse date \"%s\"", date)
		return "", errs.ErrTransactionTimeInvalid
	}

	if len(month) < 2 {
		month = "0" + month
	}

	if len(day) < 2 {
		day = "0" + day
	}

	return fmt.Sprintf("%s-%s-%s 00:00:00", year, month, day), nil
}

func createNewQifTransactionDataTable(dateFormatType qifDateFormatType, qifData *qifData) (*qifTransactionDataTable, error) {
	if qifData == nil {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	allData := make([]*qifTransactionData, 0)
	allData = append(allData, qifData.bankAccountTransactions...)
	allData = append(allData, qifData.cashAccountTransactions...)
	allData = append(allData, qifData.creditCardAccountTransactions...)
	allData = append(allData, qifData.assetAccountTransactions...)
	allData = append(allData, qifData.liabilityAccountTransactions...)

	return &qifTransactionDataTable{
		dateFormatType: dateFormatType,
		allData:        allData,
	}, nil
}
