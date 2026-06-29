package ai

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var aiTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:       true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

// aiRecognizedTransactionDataTable defines the structure of AI recognized transaction data table
type aiRecognizedTransactionDataTable struct {
	allData []*models.RecognizedTransactionResult
}

// aiRecognizedTransactionDataRow defines the structure of AI recognized transaction data row
type aiRecognizedTransactionDataRow struct {
	dataTable  *aiRecognizedTransactionDataTable
	data       *models.RecognizedTransactionResult
	finalItems map[datatable.TransactionDataTableColumn]string
}

// aiRecognizedTransactionDataRowIterator defines the structure of AI recognized transaction data row iterator
type aiRecognizedTransactionDataRowIterator struct {
	dataTable    *aiRecognizedTransactionDataTable
	currentIndex int
}

// HasColumn returns whether the transaction data table has specified column
func (t *aiRecognizedTransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := aiTransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *aiRecognizedTransactionDataTable) TransactionRowCount() int {
	return len(t.allData)
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *aiRecognizedTransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &aiRecognizedTransactionDataRowIterator{
		dataTable:    t,
		currentIndex: -1,
	}
}

// IsValid returns whether this row is valid data for importing
func (r *aiRecognizedTransactionDataRow) IsValid() bool {
	return true
}

// GetData returns the data in the specified column type
func (r *aiRecognizedTransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := aiTransactionSupportedColumns[column]

	if exists {
		return r.finalItems[column]
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *aiRecognizedTransactionDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allData)
}

// Next returns the next transaction data row
func (t *aiRecognizedTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow datatable.TransactionDataRow, err error) {
	if t.currentIndex+1 >= len(t.dataTable.allData) {
		return nil, nil
	}

	t.currentIndex++

	data := t.dataTable.allData[t.currentIndex]
	rowItems, err := t.parseTransaction(ctx, user, data)

	if err != nil {
		log.Errorf(ctx, "[ai_recognized_transaction_data_table.Next] cannot parsing transaction in row#%d, because %s", t.currentIndex, err.Error())
		return nil, err
	}

	return &aiRecognizedTransactionDataRow{
		dataTable:  t.dataTable,
		data:       data,
		finalItems: rowItems,
	}, nil
}

func (t *aiRecognizedTransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, result *models.RecognizedTransactionResult) (map[datatable.TransactionDataTableColumn]string, error) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(aiTransactionSupportedColumns))
	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = result.Time

	if result == nil || len(result.Type) == 0 {
		return nil, errs.ErrTransactionTypeInvalid
	}

	if result.Type == "income" {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_INCOME))
	} else if result.Type == "expense" {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE))
	} else if result.Type == "transfer" {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER))
	} else {
		return nil, errs.ErrTransactionTypeInvalid
	}

	data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = result.CategoryName
	data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = result.AccountName
	data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = result.Amount
	data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = result.DestinationAccountName
	data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = result.DestinationAmount
	data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = result.Description

	return data, nil
}

func (t *aiRecognizedTransactionDataRowIterator) getLongDateTime(dateTime string) string {
	if utils.IsValidLongDateTimeFormat(dateTime) {
		return dateTime
	}

	if utils.IsValidLongDateTimeWithoutSecondFormat(dateTime) {
		return dateTime + ":00"
	}

	if utils.IsValidLongDateFormat(dateTime) {
		return dateTime + " 00:00:00"
	}

	return dateTime
}

func createNewAIRecognizedTransactionDataTable(recognizedTransactions []*models.RecognizedTransactionResult) (*aiRecognizedTransactionDataTable, error) {
	if recognizedTransactions == nil || len(recognizedTransactions) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return &aiRecognizedTransactionDataTable{
		allData: recognizedTransactions,
	}, nil
}
