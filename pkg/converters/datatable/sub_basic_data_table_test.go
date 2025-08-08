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

	subTable := CreateSubBasicTable(basicDataTable, 1, 2)
	assert.Equal(t, 1, subTable.DataRowCount())
	assert.Equal(t, columns, subTable.HeaderColumnNames())
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

	subTable := CreateSubBasicTable(basicDataTable, -1, 2)
	assert.Equal(t, 0, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTable(basicDataTable, 5, 2)
	assert.Equal(t, 2, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTable(basicDataTable, 0, 5)
	assert.Equal(t, 0, subTable.fromIndex)
	assert.Equal(t, 2, subTable.toIndex)

	subTable = CreateSubBasicTable(basicDataTable, 2, 1)
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

	subTable := CreateSubBasicTable(basicDataTable, 1, 3)
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

	subTable := CreateSubBasicTable(basicDataTable, 1, 1)
	assert.Equal(t, 0, subTable.DataRowCount())

	iterator := subTable.DataRowIterator()
	assert.False(t, iterator.HasNext())
	assert.Nil(t, iterator.Next())
}
