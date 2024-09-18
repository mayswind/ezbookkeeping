package datatable

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestWritableDataTableAdd(t *testing.T) {
	columns := make([]DataTableColumn, 5)
	columns[0] = DATA_TABLE_TRANSACTION_TIME
	columns[1] = DATA_TABLE_TRANSACTION_TYPE
	columns[2] = DATA_TABLE_SUB_CATEGORY
	columns[3] = DATA_TABLE_ACCOUNT_NAME
	columns[4] = DATA_TABLE_AMOUNT

	writableDataTable := CreateNewWritableDataTable(columns)

	assert.Equal(t, 0, writableDataTable.DataRowCount())

	expectedTransactionUnixTime := time.Now().Unix()
	expectedTextualTransactionTime := utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTime, time.Local)
	expectedTransactionType := "Expense"
	expectedSubCategory := "Test Category"
	expectedAccountName := "Test Account"
	expectedAmount := "123.45"

	writableDataTable.Add(map[DataTableColumn]string{
		DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTime,
		DATA_TABLE_TRANSACTION_TYPE: expectedTransactionType,
		DATA_TABLE_SUB_CATEGORY:     expectedSubCategory,
		DATA_TABLE_ACCOUNT_NAME:     expectedAccountName,
		DATA_TABLE_AMOUNT:           expectedAmount,
	})
	assert.Equal(t, 1, writableDataTable.DataRowCount())

	dataRow := writableDataTable.Get(0)

	actualTransactionTime, err := dataRow.GetTime(0, utils.GetTimezoneOffsetMinutes(time.Local))
	assert.Nil(t, err)

	actualTransactionUnixTime := actualTransactionTime.Unix()
	assert.Equal(t, expectedTransactionUnixTime, actualTransactionUnixTime)

	actualTextualTransactionTime := dataRow.GetData(0)
	assert.Equal(t, expectedTextualTransactionTime, actualTextualTransactionTime)

	actualTransactionType := dataRow.GetData(1)
	assert.Equal(t, expectedTransactionType, actualTransactionType)

	actualSubCategory := dataRow.GetData(2)
	assert.Equal(t, expectedSubCategory, actualSubCategory)

	actualAccountName := dataRow.GetData(3)
	assert.Equal(t, expectedAccountName, actualAccountName)

	actualAmount := dataRow.GetData(4)
	assert.Equal(t, expectedAmount, actualAmount)
}

func TestWritableDataTableAdd_NotExistsColumn(t *testing.T) {
	columns := make([]DataTableColumn, 1)
	columns[0] = DATA_TABLE_TRANSACTION_TIME

	writableDataTable := CreateNewWritableDataTable(columns)

	expectedTransactionUnixTime := time.Now().Unix()
	expectedTextualTransactionTime := utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTime, time.Local)
	expectedTransactionType := "Expense"

	writableDataTable.Add(map[DataTableColumn]string{
		DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTime,
		DATA_TABLE_TRANSACTION_TYPE: expectedTransactionType,
	})
	assert.Equal(t, 1, writableDataTable.DataRowCount())

	dataRow := writableDataTable.Get(0)
	assert.Equal(t, 1, dataRow.ColumnCount())
}

func TestWritableDataTableGet_NotExistsRow(t *testing.T) {
	columns := make([]DataTableColumn, 1)
	columns[0] = DATA_TABLE_TRANSACTION_TIME

	writableDataTable := CreateNewWritableDataTable(columns)
	assert.Equal(t, 0, writableDataTable.DataRowCount())

	dataRow := writableDataTable.Get(0)
	assert.Nil(t, dataRow)
}

func TestWritableDataRowGetData_NotExistsColumn(t *testing.T) {
	columns := make([]DataTableColumn, 1)
	columns[0] = DATA_TABLE_TRANSACTION_TIME

	writableDataTable := CreateNewWritableDataTable(columns)

	expectedTransactionUnixTime := time.Now().Unix()
	expectedTextualTransactionTime := utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTime, time.Local)

	writableDataTable.Add(map[DataTableColumn]string{
		DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTime,
	})
	assert.Equal(t, 1, writableDataTable.DataRowCount())

	dataRow := writableDataTable.Get(0)
	assert.Equal(t, 1, dataRow.ColumnCount())
	assert.Equal(t, "", dataRow.GetData(1))
}

func TestWritableDataTableDataRowIterator(t *testing.T) {
	columns := make([]DataTableColumn, 5)
	columns[0] = DATA_TABLE_TRANSACTION_TIME
	columns[1] = DATA_TABLE_TRANSACTION_TYPE
	columns[2] = DATA_TABLE_SUB_CATEGORY
	columns[3] = DATA_TABLE_ACCOUNT_NAME
	columns[4] = DATA_TABLE_AMOUNT

	writableDataTable := CreateNewWritableDataTable(columns)
	assert.Equal(t, 0, writableDataTable.DataRowCount())

	expectedTransactionUnixTimes := make([]int64, 3)
	expectedTextualTransactionTimes := make([]string, 3)
	expectedTransactionTypes := make([]string, 3)
	expectedSubCategories := make([]string, 3)
	expectedAccountNames := make([]string, 3)
	expectedAmounts := make([]string, 3)

	expectedTransactionUnixTimes[0] = time.Now().Add(-5 * time.Hour).Unix()
	expectedTextualTransactionTimes[0] = utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTimes[0], time.Local)
	expectedTransactionTypes[0] = "Balance Modification"
	expectedSubCategories[0] = ""
	expectedAccountNames[0] = "Test Account"
	expectedAmounts[0] = "123.45"
	writableDataTable.Add(map[DataTableColumn]string{
		DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTimes[0],
		DATA_TABLE_TRANSACTION_TYPE: expectedTransactionTypes[0],
		DATA_TABLE_SUB_CATEGORY:     expectedSubCategories[0],
		DATA_TABLE_ACCOUNT_NAME:     expectedAccountNames[0],
		DATA_TABLE_AMOUNT:           expectedAmounts[0],
	})

	expectedTransactionUnixTimes[1] = time.Now().Add(-45 * time.Minute).Unix()
	expectedTextualTransactionTimes[1] = utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTimes[1], time.Local)
	expectedTransactionTypes[1] = "Expense"
	expectedSubCategories[1] = "Test Category2"
	expectedAccountNames[1] = "Test Account"
	expectedAmounts[1] = "-23.4"
	writableDataTable.Add(map[DataTableColumn]string{
		DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTimes[1],
		DATA_TABLE_TRANSACTION_TYPE: expectedTransactionTypes[1],
		DATA_TABLE_SUB_CATEGORY:     expectedSubCategories[1],
		DATA_TABLE_ACCOUNT_NAME:     expectedAccountNames[1],
		DATA_TABLE_AMOUNT:           expectedAmounts[1],
	})

	expectedTransactionUnixTimes[2] = time.Now().Unix()
	expectedTextualTransactionTimes[2] = utils.FormatUnixTimeToLongDateTime(expectedTransactionUnixTimes[2], time.Local)
	expectedTransactionTypes[2] = "Income"
	expectedSubCategories[2] = "Test Category3"
	expectedAccountNames[2] = "Test Account2"
	expectedAmounts[2] = "123"
	writableDataTable.Add(map[DataTableColumn]string{
		DATA_TABLE_TRANSACTION_TIME: expectedTextualTransactionTimes[2],
		DATA_TABLE_TRANSACTION_TYPE: expectedTransactionTypes[2],
		DATA_TABLE_SUB_CATEGORY:     expectedSubCategories[2],
		DATA_TABLE_ACCOUNT_NAME:     expectedAccountNames[2],
		DATA_TABLE_AMOUNT:           expectedAmounts[2],
	})
	assert.Equal(t, 3, writableDataTable.DataRowCount())

	index := 0
	iterator := writableDataTable.DataRowIterator()

	for iterator.HasNext() {
		dataRow := iterator.Next()

		actualTransactionTime, err := dataRow.GetTime(0, utils.GetTimezoneOffsetMinutes(time.Local))
		assert.Nil(t, err)

		actualTransactionUnixTime := actualTransactionTime.Unix()
		assert.Equal(t, expectedTransactionUnixTimes[index], actualTransactionUnixTime)

		actualTextualTransactionTime := dataRow.GetData(0)
		assert.Equal(t, expectedTextualTransactionTimes[index], actualTextualTransactionTime)

		actualTransactionType := dataRow.GetData(1)
		assert.Equal(t, expectedTransactionTypes[index], actualTransactionType)

		actualSubCategory := dataRow.GetData(2)
		assert.Equal(t, expectedSubCategories[index], actualSubCategory)

		actualAccountName := dataRow.GetData(3)
		assert.Equal(t, expectedAccountNames[index], actualAccountName)

		actualAmount := dataRow.GetData(4)
		assert.Equal(t, expectedAmounts[index], actualAmount)

		index++
	}

	assert.Equal(t, 3, index)
}
