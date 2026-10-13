package datatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSubBasicTable_WithValidInput(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"A1", "B1", "C1"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"A2", "B2", "C2"},
		},
		{
			rowId:      "3",
			rowColumns: []string{"A3", "B3", "C3"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	subTable := CreateSubBasicTable(basicDataTable, 0, 2, true)
	assert.Equal(t, 1, subTable.DataRowCount())
	assert.Equal(t, []string{"A1", "B1", "C1"}, subTable.HeaderColumnNames())

	subTable = CreateSubBasicTable(basicDataTable, 1, 2, true)
	assert.Equal(t, 0, subTable.DataRowCount())
	assert.Equal(t, []string{"A2", "B2", "C2"}, subTable.HeaderColumnNames())

	subTable = CreateSubBasicTable(basicDataTable, 1, 2, false)
	assert.Equal(t, 1, subTable.DataRowCount())
	assert.Equal(t, []string(nil), subTable.HeaderColumnNames())
}

func TestCreateSubBasicTable_WithInvalidInput(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"A1", "B1", "C1"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"A2", "B2", "C2"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	subTable := CreateSubBasicTable(basicDataTable, -1, 2, false)
	assert.Equal(t, 0, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTable(basicDataTable, 5, 2, false)
	assert.Equal(t, 2, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTable(basicDataTable, 0, 5, false)
	assert.Equal(t, 0, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTable(basicDataTable, 2, 1, false)
	assert.Equal(t, 2, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)
}

func TestCreateSubBasicTableWithParser_WithValidInput(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"A1", "B1", "C1"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"A2", "B2", "C2"},
		},
		{
			rowId:      "3",
			rowColumns: []string{"A3", "B3", "C3"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	parseFunc := func(data string) string {
		return "parsed_" + data
	}

	subTable := CreateSubBasicTableWithParser(basicDataTable, 0, 3, true, parseFunc)
	assert.Equal(t, 2, subTable.DataRowCount())
	assert.Equal(t, []string{"A1", "B1", "C1"}, subTable.HeaderColumnNames())

	subTable = CreateSubBasicTableWithParser(basicDataTable, 1, 3, false, parseFunc)
	assert.Equal(t, 2, subTable.DataRowCount())
	assert.Equal(t, []string(nil), subTable.HeaderColumnNames())
}

func TestCreateSubBasicTableWithParser_WithInvalidInput(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"A1", "B1", "C1"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"A2", "B2", "C2"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	parseFunc := func(data string) string {
		return data
	}

	subTable := CreateSubBasicTableWithParser(basicDataTable, -1, 2, false, parseFunc)
	assert.Equal(t, 0, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTableWithParser(basicDataTable, 5, 2, false, parseFunc)
	assert.Equal(t, 2, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTableWithParser(basicDataTable, 0, 5, false, parseFunc)
	assert.Equal(t, 0, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTableWithParser(basicDataTable, 2, 1, false, parseFunc)
	assert.Equal(t, 2, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)
}

func TestSubBasicDataTable_DataRowIterator(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"A1", "B1", "C1"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"A2", "B2", "C2"},
		},
		{
			rowId:      "3",
			rowColumns: []string{"A3", "B3", "C3"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	subTable := CreateSubBasicTable(basicDataTable, 1, 3, false)
	iterator := subTable.DataRowIterator()

	assert.True(t, iterator.HasNext())
	firstRow := iterator.Next()
	assert.NotNil(t, firstRow)
	assert.Equal(t, "2", iterator.CurrentRowId())
	assert.Equal(t, "A2", firstRow.GetData(0))
	assert.Equal(t, "B2", firstRow.GetData(1))
	assert.Equal(t, "C2", firstRow.GetData(2))

	assert.True(t, iterator.HasNext())
	secondRow := iterator.Next()
	assert.NotNil(t, secondRow)
	assert.Equal(t, "3", iterator.CurrentRowId())
	assert.Equal(t, "A3", secondRow.GetData(0))
	assert.Equal(t, "B3", secondRow.GetData(1))
	assert.Equal(t, "C3", secondRow.GetData(2))

	assert.False(t, iterator.HasNext())
	assert.Nil(t, iterator.Next())
}

func TestSubBasicDataTable_EmptyDataRange(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"A1", "B1", "C1"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"A2", "B2", "C2"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	subTable := CreateSubBasicTable(basicDataTable, 1, 1, false)
	assert.Equal(t, 0, subTable.DataRowCount())

	iterator := subTable.DataRowIterator()
	assert.False(t, iterator.HasNext())
	assert.Nil(t, iterator.Next())
}

func TestSubBasicDataTable_DataRowIteratorWithParser(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"A1", "B1", "C1"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"A2", "B2", "C2"},
		},
		{
			rowId:      "3",
			rowColumns: []string{"A3", "B3", "C3"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	subTable := CreateSubBasicTableWithParser(basicDataTable, 1, 3, false, func(data string) string {
		return "parsed_" + data
	})
	iterator := subTable.DataRowIterator()

	assert.True(t, iterator.HasNext())
	firstRow := iterator.Next()
	assert.NotNil(t, firstRow)
	assert.Equal(t, "2", iterator.CurrentRowId())
	assert.Equal(t, 3, firstRow.ColumnCount())
	assert.Equal(t, "parsed_A2", firstRow.GetData(0))
	assert.Equal(t, "parsed_B2", firstRow.GetData(1))
	assert.Equal(t, "parsed_C2", firstRow.GetData(2))

	assert.True(t, iterator.HasNext())
	secondRow := iterator.Next()
	assert.NotNil(t, secondRow)
	assert.Equal(t, "3", iterator.CurrentRowId())
	assert.Equal(t, "parsed_A3", secondRow.GetData(0))
	assert.Equal(t, "parsed_B3", secondRow.GetData(1))
	assert.Equal(t, "parsed_C3", secondRow.GetData(2))

	assert.False(t, iterator.HasNext())
	assert.Nil(t, iterator.Next())
}

func TestSubBasicDataTable_DataRowIteratorWithNilParser(t *testing.T) {
	basicDataTable := &testBasicDataTable{
		headerColumns: []string{"Col1", "Col2"},
		rows: []*testBasicDataTableRow{
			{
				rowId:      "1",
				rowColumns: []string{"A1", "B1"},
			},
		},
	}

	subTable := CreateSubBasicTableWithParser(basicDataTable, 0, 1, false, nil)
	iterator := subTable.DataRowIterator()

	assert.True(t, iterator.HasNext())
	row := iterator.Next()
	assert.NotNil(t, row)
	assert.Equal(t, "A1", row.GetData(0))
	assert.Equal(t, "B1", row.GetData(1))
}

