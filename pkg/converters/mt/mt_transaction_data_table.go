package mt

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var mt940TransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:     true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

// mt940TransactionDataTable represents the mt940 statement data dataTable
type mt940TransactionDataTable struct {
	data *mt940Data
}

// mt940TransactionDataRow represents a row in the mt940 statement data dataTable
type mt940TransactionDataRow struct {
	statement  *mtStatement
	finalItems map[datatable.TransactionDataTableColumn]string
}

// mt940TransactionDataRowIterator represents an iterator for mt940 statement data rows
type mt940TransactionDataRowIterator struct {
	dataTable    *mt940TransactionDataTable
	currentIndex int
}

// HasColumn implements TransactionDataTable.HasColumn
func (t *mt940TransactionDataTable) HasColumn(column datatable.TransactionDataTableColumn) bool {
	_, exists := mt940TransactionSupportedColumns[column]
	return exists
}

// TransactionRowCount implements TransactionDataTable.TransactionRowCount
func (t *mt940TransactionDataTable) TransactionRowCount() int {
	return len(t.data.Statements)
}

// TransactionRowIterator implements TransactionDataTable.TransactionRowIterator
func (t *mt940TransactionDataTable) TransactionRowIterator() datatable.TransactionDataRowIterator {
	return &mt940TransactionDataRowIterator{
		dataTable:    t,
		currentIndex: -1,
	}
}

// IsValid implements TransactionDataRow.IsValid
func (r *mt940TransactionDataRow) IsValid() bool {
	return true
}

// GetData implements TransactionDataRow.GetData
func (r *mt940TransactionDataRow) GetData(column datatable.TransactionDataTableColumn) string {
	_, exists := mt940TransactionSupportedColumns[column]

	if exists {
		return r.finalItems[column]
	}

	return ""
}

// HasNext implements TransactionDataRowIterator.HasNext
func (t *mt940TransactionDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.data.Statements)
}

// Next implements TransactionDataRowIterator.Next
func (t *mt940TransactionDataRowIterator) Next(ctx core.Context, user *models.User) (datatable.TransactionDataRow, error) {
	if t.currentIndex+1 >= len(t.dataTable.data.Statements) {
		return nil, nil
	}

	t.currentIndex++

	data := t.dataTable.data.Statements[t.currentIndex]
	rowItems, err := t.parseTransaction(ctx, user, t.dataTable.data, data)

	if err != nil {
		log.Errorf(ctx, "[mt_transaction_data_table.Next] cannot parsing transaction in row#%d, because %s", t.currentIndex, err.Error())
		return nil, err
	}

	return &mt940TransactionDataRow{
		statement:  data,
		finalItems: rowItems,
	}, nil
}

func (t *mt940TransactionDataRowIterator) parseTransaction(ctx core.Context, user *models.User, mt940Data *mt940Data, statement *mtStatement) (map[datatable.TransactionDataTableColumn]string, error) {
	data := make(map[datatable.TransactionDataTableColumn]string, len(mt940TransactionSupportedColumns))

	if statement.ValueDate == "" && len(statement.ValueDate) != 6 {
		return nil, errs.ErrTransactionTimeInvalid
	}

	transactionTime, err := utils.FormatYearMonthDayToLongDateTime(statement.ValueDate[0:2], statement.ValueDate[2:4], statement.ValueDate[4:6])

	if err != nil {
		log.Errorf(ctx, "[mt_transaction_data_table.parseTransaction] cannot format transaction time in row#%d, because %s", t.currentIndex, err.Error())
		return nil, errs.ErrTransactionTimeInvalid
	}

	data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = transactionTime
	data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = mt940Data.AccountId

	if mt940Data.OpeningBalance != nil && mt940Data.OpeningBalance.Currency != "" {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = mt940Data.OpeningBalance.Currency
	} else if mt940Data.ClosingBalance != nil && mt940Data.ClosingBalance.Currency != "" {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = mt940Data.ClosingBalance.Currency
	} else {
		return nil, errs.ErrAccountCurrencyInvalid
	}

	amountValue := strings.ReplaceAll(statement.Amount, ",", ".") // decimal separator is comma in mt data

	if len(amountValue) > 0 && amountValue[len(amountValue)-1] == '.' {
		amountValue = amountValue[:len(amountValue)-1]
	}

	amount, err := utils.ParseAmount(amountValue)

	if err != nil {
		log.Errorf(ctx, "[mt_transaction_data_table.parseTransaction] cannot parsing transaction amount \"%s\", because %s", statement.Amount, err.Error())
		return nil, errs.ErrAmountInvalid
	}

	data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(amount)

	if statement.CreditDebitMark == MT_MARK_CREDIT {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_INCOME))
	} else if statement.CreditDebitMark == MT_MARK_DEBIT {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE))
	} else if statement.CreditDebitMark == MT_MARK_REVERSAL_CREDIT {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE))
	} else if statement.CreditDebitMark == MT_MARK_REVERSAL_DEBIT {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = utils.IntToString(int(models.TRANSACTION_TYPE_INCOME))
	} else {
		return nil, errs.ErrTransactionTypeInvalid
	}

	informationToAccountOwnerMap := statement.GetInformationToAccountOwnerMap()

	if len(informationToAccountOwnerMap) > 0 {
		if value, exists := informationToAccountOwnerMap[MT_INFORMATION_TO_ACCOUNT_OWNER_TAG_REMITTANCE]; exists {
			data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = value
		}
	} else {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = strings.Join(statement.InformationToAccountOwner, "\n")
	}

	return data, nil
}

// createNewMT940TransactionDataTable creates a new mt940 statement data dataTable
func createNewMT940TransactionDataTable(data *mt940Data) (*mt940TransactionDataTable, error) {
	if data == nil || len(data.Statements) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return &mt940TransactionDataTable{
		data: data,
	}, nil
}
