package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestWritableDataTableAdd(t *testing.T) {
	columns := make([]TransactionDataTableColumn, 5)
	columns[0] = TRANSACTION_DATA_TABLE_TRANSACTION_TIME
	columns[1] = TRANSACTION_DATA_TABLE_TRANSACTION_TYPE
	columns[2] = TRANSACTION_DATA_TABLE_SUB_CATEGORY
	columns[3] = TRANSACTION_DATA_TABLE_ACCOUNT_NAME
	columns[4] = TRANSACTION_DATA_TABLE_AMOUNT

	writableDataTable := CreateNewWritableTransactionDataTable(columns)

	assert.Equal(t, 0, writableDataTable.TransactionRowCount())

	expectedTransactionTime := "2024-09-01 01:23:45"
	expectedTransactionType := "Expense"
	expectedSubCategory := "Test Category"
	expectedAccountName := "Test Account"
	expectedAmount := "123.45"

	writableDataTable.Add(map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: expectedTransactionTime,
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: expectedTransactionType,
		TRANSACTION_DATA_TABLE_SUB_CATEGORY:     expectedSubCategory,
		TRANSACTION_DATA_TABLE_ACCOUNT_NAME:     expectedAccountName,
		TRANSACTION_DATA_TABLE_AMOUNT:           expectedAmount,
	})
	assert.Equal(t, 1, writableDataTable.TransactionRowCount())

	dataRow := writableDataTable.Get(0)

	actualTransactionTime := dataRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME)
	assert.Equal(t, expectedTransactionTime, actualTransactionTime)

	actualTransactionType := dataRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE)
	assert.Equal(t, expectedTransactionType, actualTransactionType)

	actualSubCategory := dataRow.GetData(TRANSACTION_DATA_TABLE_SUB_CATEGORY)
	assert.Equal(t, expectedSubCategory, actualSubCategory)

	actualAccountName := dataRow.GetData(TRANSACTION_DATA_TABLE_ACCOUNT_NAME)
	assert.Equal(t, expectedAccountName, actualAccountName)

	actualAmount := dataRow.GetData(TRANSACTION_DATA_TABLE_AMOUNT)
	assert.Equal(t, expectedAmount, actualAmount)
}

func TestWritableDataTableAdd_NotExistsColumn(t *testing.T) {
	columns := make([]TransactionDataTableColumn, 1)
	columns[0] = TRANSACTION_DATA_TABLE_TRANSACTION_TIME

	writableDataTable := CreateNewWritableTransactionDataTable(columns)

	expectedTransactionUnixTime := time.Now().Unix()
	expectedTextualTransactionTime := utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTime, time.Local)
	expectedTransactionType := "Expense"

	writableDataTable.Add(map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTime,
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: expectedTransactionType,
	})
	assert.Equal(t, 1, writableDataTable.TransactionRowCount())

	dataRow := writableDataTable.Get(0)
	assert.Equal(t, 1, dataRow.ColumnCount())
}

func TestWritableDataTableGet_NotExistsRow(t *testing.T) {
	columns := make([]TransactionDataTableColumn, 1)
	columns[0] = TRANSACTION_DATA_TABLE_TRANSACTION_TIME

	writableDataTable := CreateNewWritableTransactionDataTable(columns)
	assert.Equal(t, 0, writableDataTable.TransactionRowCount())

	dataRow := writableDataTable.Get(0)
	assert.Nil(t, dataRow)
}

func TestWritableDataRowGetData_NotExistsColumn(t *testing.T) {
	columns := make([]TransactionDataTableColumn, 1)
	columns[0] = TRANSACTION_DATA_TABLE_TRANSACTION_TIME

	writableDataTable := CreateNewWritableTransactionDataTable(columns)

	expectedTransactionUnixTime := time.Now().Unix()
	expectedTextualTransactionTime := utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTime, time.Local)

	writableDataTable.Add(map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTime,
	})
	assert.Equal(t, 1, writableDataTable.TransactionRowCount())

	dataRow := writableDataTable.Get(0)
	assert.Equal(t, 1, dataRow.ColumnCount())
	assert.Equal(t, "", dataRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
}

func TestWritableDataTableDataRowIterator(t *testing.T) {
	columns := make([]TransactionDataTableColumn, 5)
	columns[0] = TRANSACTION_DATA_TABLE_TRANSACTION_TIME
	columns[1] = TRANSACTION_DATA_TABLE_TRANSACTION_TYPE
	columns[2] = TRANSACTION_DATA_TABLE_SUB_CATEGORY
	columns[3] = TRANSACTION_DATA_TABLE_ACCOUNT_NAME
	columns[4] = TRANSACTION_DATA_TABLE_AMOUNT

	writableDataTable := CreateNewWritableTransactionDataTable(columns)
	assert.Equal(t, 0, writableDataTable.TransactionRowCount())

	expectedTransactionUnixTimes := make([]int64, 3)
	expectedTransactionTimes := make([]string, 3)
	expectedTransactionTypes := make([]string, 3)
	expectedSubCategories := make([]string, 3)
	expectedAccountNames := make([]string, 3)
	expectedAmounts := make([]string, 3)

	expectedTransactionUnixTimes[0] = time.Now().Add(-5 * time.Hour).Unix()
	expectedTransactionTimes[0] = utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTimes[0], time.Local)
	expectedTransactionTypes[0] = "Balance Modification"
	expectedSubCategories[0] = ""
	expectedAccountNames[0] = "Test Account"
	expectedAmounts[0] = "123.45"
	writableDataTable.Add(map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: expectedTransactionTimes[0],
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: expectedTransactionTypes[0],
		TRANSACTION_DATA_TABLE_SUB_CATEGORY:     expectedSubCategories[0],
		TRANSACTION_DATA_TABLE_ACCOUNT_NAME:     expectedAccountNames[0],
		TRANSACTION_DATA_TABLE_AMOUNT:           expectedAmounts[0],
	})

	expectedTransactionUnixTimes[1] = time.Now().Add(-45 * time.Minute).Unix()
	expectedTransactionTimes[1] = utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTimes[1], time.Local)
	expectedTransactionTypes[1] = "Expense"
	expectedSubCategories[1] = "Test Category2"
	expectedAccountNames[1] = "Test Account"
	expectedAmounts[1] = "-23.4"
	writableDataTable.Add(map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: expectedTransactionTimes[1],
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: expectedTransactionTypes[1],
		TRANSACTION_DATA_TABLE_SUB_CATEGORY:     expectedSubCategories[1],
		TRANSACTION_DATA_TABLE_ACCOUNT_NAME:     expectedAccountNames[1],
		TRANSACTION_DATA_TABLE_AMOUNT:           expectedAmounts[1],
	})

	expectedTransactionUnixTimes[2] = time.Now().Unix()
	expectedTransactionTimes[2] = utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTimes[2], time.Local)
	expectedTransactionTypes[2] = "Income"
	expectedSubCategories[2] = "Test Category3"
	expectedAccountNames[2] = "Test Account2"
	expectedAmounts[2] = "123"
	writableDataTable.Add(map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: expectedTransactionTimes[2],
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: expectedTransactionTypes[2],
		TRANSACTION_DATA_TABLE_SUB_CATEGORY:     expectedSubCategories[2],
		TRANSACTION_DATA_TABLE_ACCOUNT_NAME:     expectedAccountNames[2],
		TRANSACTION_DATA_TABLE_AMOUNT:           expectedAmounts[2],
	})
	assert.Equal(t, 3, writableDataTable.TransactionRowCount())

	index := 0
	iterator := writableDataTable.TransactionRowIterator()

	for iterator.HasNext() {
		dataRow, err := iterator.Next(core.NewNullContext(), &models.User{})
		assert.Nil(t, err)

		actualTransactionTime := dataRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME)
		assert.Equal(t, expectedTransactionTimes[index], actualTransactionTime)

		actualTransactionType := dataRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE)
		assert.Equal(t, expectedTransactionTypes[index], actualTransactionType)

		actualSubCategory := dataRow.GetData(TRANSACTION_DATA_TABLE_SUB_CATEGORY)
		assert.Equal(t, expectedSubCategories[index], actualSubCategory)

		actualAccountName := dataRow.GetData(TRANSACTION_DATA_TABLE_ACCOUNT_NAME)
		assert.Equal(t, expectedAccountNames[index], actualAccountName)

		actualAmount := dataRow.GetData(TRANSACTION_DATA_TABLE_AMOUNT)
		assert.Equal(t, expectedAmounts[index], actualAmount)

		index++
	}

	assert.Equal(t, 3, index)
}
