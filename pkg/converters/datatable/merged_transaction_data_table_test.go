package datatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergedTransactionDataTable_HasColumn(t *testing.T) {
	columns1 := []string{"TransactionTime", "TransactionType", "Amount"}
	basicDataTable1 := &testBasicDataTable{
		headerColumns: columns1,
		rows:          []*testBasicDataTableRow{},
	}

	columns2 := []string{"TransactionTime", "TransactionType", "Category"}
	basicDataTable2 := &testBasicDataTable{
		headerColumns: columns2,
		rows:          []*testBasicDataTableRow{},
	}

	columnMapping1 := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	columnMapping2 := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_CATEGORY:         "Category",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping1)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping2)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2}, "sheet")

	assert.True(t, merged.HasColumn(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.True(t, merged.HasColumn(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
	assert.True(t, merged.HasColumn(TRANSACTION_DATA_TABLE_AMOUNT))
	assert.True(t, merged.HasColumn(TRANSACTION_DATA_TABLE_CATEGORY))
	assert.False(t, merged.HasColumn(TRANSACTION_DATA_TABLE_DESCRIPTION))
}

func TestMergedTransactionDataTable_HasColumn_Empty(t *testing.T) {
	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{}, "sheet")

	assert.False(t, merged.HasColumn(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.False(t, merged.HasColumn(TRANSACTION_DATA_TABLE_AMOUNT))
}

func TestMergedTransactionDataTable_TransactionRowCount(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows1 := []*testBasicDataTableRow{
		{rowId: "1", rowColumns: []string{"2026-07-01", "1", "100"}},
		{rowId: "2", rowColumns: []string{"2026-07-02", "2", "200"}},
	}
	rows2 := []*testBasicDataTableRow{
		{rowId: "3", rowColumns: []string{"2026-07-03", "1", "300"}},
	}
	rows3 := []*testBasicDataTableRow{}

	basicDataTable1 := &testBasicDataTable{headerColumns: columns, rows: rows1}
	basicDataTable2 := &testBasicDataTable{headerColumns: columns, rows: rows2}
	basicDataTable3 := &testBasicDataTable{headerColumns: columns, rows: rows3}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping)
	table3 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable3, columnMapping)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2, table3}, "sheet")

	assert.Equal(t, 3, merged.TransactionRowCount())
}

func TestMergedTransactionDataTable_TransactionRowCount_Empty(t *testing.T) {
	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{}, "sheet")
	assert.Equal(t, 0, merged.TransactionRowCount())
}

func TestMergedTransactionDataTable_TransactionRowIterator(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows1 := []*testBasicDataTableRow{
		{rowId: "1", rowColumns: []string{"2026-07-01", "1", "100"}},
		{rowId: "2", rowColumns: []string{"2026-07-02", "2", "200"}},
	}
	rows2 := []*testBasicDataTableRow{
		{rowId: "3", rowColumns: []string{"2026-07-03", "1", "300"}},
		{rowId: "4", rowColumns: []string{"2026-07-04", "2", "400"}},
	}

	basicDataTable1 := &testBasicDataTable{headerColumns: columns, rows: rows1}
	basicDataTable2 := &testBasicDataTable{headerColumns: columns, rows: rows2}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2}, "sheet")
	iterator := merged.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.True(t, row.IsValid())
	assert.Equal(t, "2026-07-01", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "100", row.GetData(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.True(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.True(t, row.IsValid())
	assert.Equal(t, "2026-07-02", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "200", row.GetData(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.True(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.True(t, row.IsValid())
	assert.Equal(t, "2026-07-03", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "300", row.GetData(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.True(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.True(t, row.IsValid())
	assert.Equal(t, "2026-07-04", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "400", row.GetData(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.False(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, row)

	assert.False(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, row)
}

func TestMergedTransactionDataTable_TransactionRowIterator_WithEmptyTable(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows1 := []*testBasicDataTableRow{
		{rowId: "1", rowColumns: []string{"2026-07-01", "1", "100"}},
	}
	rows2 := []*testBasicDataTableRow{}
	rows3 := []*testBasicDataTableRow{
		{rowId: "2", rowColumns: []string{"2026-07-02", "2", "200"}},
	}

	basicDataTable1 := &testBasicDataTable{headerColumns: columns, rows: rows1}
	basicDataTable2 := &testBasicDataTable{headerColumns: columns, rows: rows2}
	basicDataTable3 := &testBasicDataTable{headerColumns: columns, rows: rows3}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping)
	table3 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable3, columnMapping)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2, table3}, "sheet")
	iterator := merged.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.Equal(t, "2026-07-01", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))

	assert.True(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.Equal(t, "2026-07-02", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))

	assert.False(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, row)
}

func TestMergedTransactionDataTable_TransactionRowIterator_Empty(t *testing.T) {
	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{}, "sheet")
	iterator := merged.TransactionRowIterator()

	assert.False(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, row)
}

func TestMergedTransactionDataTable_TransactionRowIterator_AllEmptyTables(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	basicDataTable1 := &testBasicDataTable{headerColumns: columns, rows: []*testBasicDataTableRow{}}
	basicDataTable2 := &testBasicDataTable{headerColumns: columns, rows: []*testBasicDataTableRow{}}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2}, "sheet")
	iterator := merged.TransactionRowIterator()

	assert.False(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, row)
}

func TestMergedTransactionDataTable_TransactionRowIterator_DifferentColumns(t *testing.T) {
	columns1 := []string{"TransactionTime", "TransactionType", "Amount"}
	columns2 := []string{"TransactionTime", "TransactionType", "Category"}
	rows1 := []*testBasicDataTableRow{
		{rowId: "1", rowColumns: []string{"2026-07-01", "1", "100"}},
	}
	rows2 := []*testBasicDataTableRow{
		{rowId: "2", rowColumns: []string{"2026-07-02", "2", "Food"}},
	}

	basicDataTable1 := &testBasicDataTable{headerColumns: columns1, rows: rows1}
	basicDataTable2 := &testBasicDataTable{headerColumns: columns2, rows: rows2}

	columnMapping1 := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}
	columnMapping2 := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_CATEGORY:         "Category",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping1)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping2)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2}, "sheet")
	iterator := merged.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.Equal(t, "2026-07-01", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "100", row.GetData(TRANSACTION_DATA_TABLE_AMOUNT))
	assert.Equal(t, "", row.GetData(TRANSACTION_DATA_TABLE_CATEGORY))

	assert.True(t, iterator.HasNext())
	row, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.Equal(t, "2026-07-02", row.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "", row.GetData(TRANSACTION_DATA_TABLE_AMOUNT))
	assert.Equal(t, "Food", row.GetData(TRANSACTION_DATA_TABLE_CATEGORY))

	assert.False(t, iterator.HasNext())
}

func TestMergedTransactionDataTable_TransactionRowIterator_CurrentRowId(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows1 := []*testBasicDataTableRow{
		{rowId: "1", rowColumns: []string{"2026-07-01", "1", "100"}},
		{rowId: "2", rowColumns: []string{"2026-07-02", "2", "200"}},
	}
	rows2 := []*testBasicDataTableRow{
		{rowId: "3", rowColumns: []string{"2026-07-03", "1", "300"}},
	}

	basicDataTable1 := &testBasicDataTable{headerColumns: columns, rows: rows1}
	basicDataTable2 := &testBasicDataTable{headerColumns: columns, rows: rows2}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2}, "file")
	iterator := merged.TransactionRowIterator()

	_, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "file:1-1", iterator.CurrentRowId())

	_, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "file:1-2", iterator.CurrentRowId())

	_, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "file:2-3", iterator.CurrentRowId())
}

func TestMergedTransactionDataTable_TransactionRowIterator_CurrentRowId_WithEmptyPrefix(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows := []*testBasicDataTableRow{
		{rowId: "5", rowColumns: []string{"2026-07-01", "1", "100"}},
	}

	basicDataTable := &testBasicDataTable{headerColumns: columns, rows: rows}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	table := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable, columnMapping)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table}, "")
	iterator := merged.TransactionRowIterator()

	_, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "sheet:1-5", iterator.CurrentRowId())
}

func TestMergedTransactionDataTable_TransactionRowIterator_CurrentRowId_WithEmptyTable(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows1 := []*testBasicDataTableRow{
		{rowId: "1", rowColumns: []string{"2026-07-01", "1", "100"}},
	}
	rows2 := []*testBasicDataTableRow{}
	rows3 := []*testBasicDataTableRow{
		{rowId: "2", rowColumns: []string{"2026-07-02", "2", "200"}},
	}

	basicDataTable1 := &testBasicDataTable{headerColumns: columns, rows: rows1}
	basicDataTable2 := &testBasicDataTable{headerColumns: columns, rows: rows2}
	basicDataTable3 := &testBasicDataTable{headerColumns: columns, rows: rows3}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	table1 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable1, columnMapping)
	table2 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable2, columnMapping)
	table3 := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable3, columnMapping)

	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{table1, table2, table3}, "sheet")
	iterator := merged.TransactionRowIterator()

	_, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "sheet:1-1", iterator.CurrentRowId())

	_, err = iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "sheet:3-2", iterator.CurrentRowId())
}

func TestMergedTransactionDataTable_TransactionRowIterator_CurrentRowId_Empty(t *testing.T) {
	merged := CreateNewMergedTransactionDataTable([]TransactionDataTable{}, "sheet")
	iterator := merged.TransactionRowIterator()

	assert.Equal(t, "", iterator.CurrentRowId())
}
