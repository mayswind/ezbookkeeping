package datatable

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

type testCommonDataTable struct {
	headerColumns []string
	dataRows      []*testCommonDataTableRow
}

type testCommonDataTableRow struct {
	rowId   string
	rowData map[string]string
}

type testCommonDataTableRowIterator struct {
	dataTable    *testCommonDataTable
	currentIndex int
}

func (t *testCommonDataTable) DataRowCount() int {
	return len(t.dataRows)
}

func (t *testCommonDataTable) HeaderColumnCount() int {
	return len(t.headerColumns)
}

func (t *testCommonDataTable) HasColumn(columnName string) bool {
	for _, header := range t.headerColumns {
		if header == columnName {
			return true
		}
	}
	return false
}

func (t *testCommonDataTable) DataRowIterator() CommonDataTableRowIterator {
	return &testCommonDataTableRowIterator{
		dataTable:    t,
		currentIndex: -1,
	}
}

func (t *testCommonDataTableRow) GetData(dataKey string) string {
	return t.rowData[dataKey]
}

func (t *testCommonDataTableRow) HasData(dataKey string) bool {
	_, exists := t.rowData[dataKey]
	return exists
}

func (t *testCommonDataTableRow) ColumnCount() int {
	return len(t.rowData)
}

func (t *testCommonDataTableRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.dataRows)
}

func (t *testCommonDataTableRowIterator) Next() CommonDataTableRow {
	if !t.HasNext() {
		return nil
	}

	t.currentIndex++
	return t.dataTable.dataRows[t.currentIndex]
}

func (t *testCommonDataTableRowIterator) CurrentRowId() string {
	if t.currentIndex < 0 || t.currentIndex >= len(t.dataTable.dataRows) {
		return ""
	}

	return t.dataTable.dataRows[t.currentIndex].rowId
}

type testCommonTransactionDataRowParser struct {
	returnError bool
}

func (p *testCommonTransactionDataRowParser) Parse(ctx core.Context, user *models.User, dataRow CommonDataTableRow, rowId string) (map[TransactionDataTableColumn]string, bool, error) {
	if p.returnError {
		return nil, false, errs.ErrOperationFailed
	}

	rowData := make(map[TransactionDataTableColumn]string)
	rowData[TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = dataRow.GetData("TransactionTime")
	rowData[TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = dataRow.GetData("TransactionType")
	rowData[TRANSACTION_DATA_TABLE_AMOUNT] = dataRow.GetData("Amount")
	rowData[TRANSACTION_DATA_TABLE_DESCRIPTION] = "Test Description"
	return rowData, true, nil
}

func TestCommonDataTableToTransactionDataTableWrapper_HasColumn(t *testing.T) {
	basicDataTable := &testCommonDataTable{
		headerColumns: []string{"TransactionTime", "TransactionType", "Amount"},
		dataRows:      []*testCommonDataTableRow{},
	}

	supportedColumns := map[TransactionDataTableColumn]bool{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: true,
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: true,
		TRANSACTION_DATA_TABLE_AMOUNT:           true,
	}

	transactionDataTable := CreateNewTransactionDataTableFromCommonDataTable(basicDataTable, supportedColumns, &testCommonTransactionDataRowParser{})

	assert.True(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.True(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
	assert.True(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.False(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_CATEGORY))
	assert.False(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_DESCRIPTION))
}

func TestCommonDataTableToTransactionDataTableWrapper_TransactionRowCount(t *testing.T) {
	rows := []*testCommonDataTableRow{
		{
			rowId: "1",
			rowData: map[string]string{
				"TransactionTime": "2024-01-01",
				"TransactionType": "1",
				"Amount":          "100",
			},
		},
		{
			rowId: "2",
			rowData: map[string]string{
				"TransactionTime": "2024-01-02",
				"TransactionType": "2",
				"Amount":          "200",
			},
		},
		{
			rowId: "3",
			rowData: map[string]string{
				"TransactionTime": "2024-01-03",
				"TransactionType": "1",
				"Amount":          "300",
			},
		},
	}

	basicDataTable := &testCommonDataTable{
		headerColumns: []string{"TransactionTime", "TransactionType", "Amount"},
		dataRows:      rows,
	}

	supportedColumns := map[TransactionDataTableColumn]bool{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: true,
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: true,
		TRANSACTION_DATA_TABLE_AMOUNT:           true,
	}

	transactionDataTable := CreateNewTransactionDataTableFromCommonDataTable(basicDataTable, supportedColumns, &testCommonTransactionDataRowParser{})
	assert.Equal(t, len(rows), transactionDataTable.TransactionRowCount())
}

func TestCommonDataTableToTransactionDataTableWrapper_TransactionRowIterator(t *testing.T) {
	rows := []*testCommonDataTableRow{
		{
			rowId: "1",
			rowData: map[string]string{
				"TransactionTime": "2024-01-01",
				"TransactionType": "1",
				"Amount":          "100",
			},
		},
		{
			rowId: "2",
			rowData: map[string]string{
				"TransactionTime": "2024-01-02",
				"TransactionType": "2",
				"Amount":          "200",
			},
		},
	}

	basicDataTable := &testCommonDataTable{
		headerColumns: []string{"TransactionTime", "TransactionType", "Amount"},
		dataRows:      rows,
	}

	supportedColumns := map[TransactionDataTableColumn]bool{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: true,
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: true,
		TRANSACTION_DATA_TABLE_AMOUNT:           true,
	}

	transactionDataTable := CreateNewTransactionDataTableFromCommonDataTable(basicDataTable, supportedColumns, &testCommonTransactionDataRowParser{})
	iterator := transactionDataTable.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	firstRow, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, firstRow)
	assert.True(t, firstRow.IsValid())
	assert.Equal(t, "2024-01-01", firstRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "1", firstRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
	assert.Equal(t, "100", firstRow.GetData(TRANSACTION_DATA_TABLE_AMOUNT))
	assert.Equal(t, "", firstRow.GetData(TRANSACTION_DATA_TABLE_DESCRIPTION))

	assert.True(t, iterator.HasNext())
	secondRow, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, secondRow)
	assert.True(t, secondRow.IsValid())
	assert.Equal(t, "2024-01-02", secondRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "2", secondRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
	assert.Equal(t, "200", secondRow.GetData(TRANSACTION_DATA_TABLE_AMOUNT))
	assert.Equal(t, "", secondRow.GetData(TRANSACTION_DATA_TABLE_DESCRIPTION))

	assert.False(t, iterator.HasNext())
	emptyRow, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, emptyRow)
}

func TestCommonDataTableToTransactionDataTableWrapper_TransactionRowIterator_EOF(t *testing.T) {
	rows := []*testCommonDataTableRow{
		{
			rowId: "1",
			rowData: map[string]string{
				"TransactionTime": "2024-01-01",
				"TransactionType": "1",
				"Amount":          "100",
			},
		},
	}

	basicDataTable := &testCommonDataTable{
		headerColumns: []string{"TransactionTime", "TransactionType", "Amount"},
		dataRows:      rows,
	}

	supportedColumns := map[TransactionDataTableColumn]bool{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: true,
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: true,
		TRANSACTION_DATA_TABLE_AMOUNT:           true,
	}

	transactionDataTable := CreateNewTransactionDataTableFromCommonDataTable(basicDataTable, supportedColumns, &testCommonTransactionDataRowParser{returnError: true})
	iterator := transactionDataTable.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.EqualError(t, err, errs.ErrOperationFailed.Message)
	assert.Nil(t, row)
}
