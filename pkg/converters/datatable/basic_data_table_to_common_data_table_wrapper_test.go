package datatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicDataTableToCommonDataTableWrapper_HeaderColumnCount(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          []*testBasicDataTableRow{},
	}

	commonDataTable := CreateNewCommonDataTableFromBasicDataTable(basicDataTable)
	assert.Equal(t, len(columns), commonDataTable.HeaderColumnCount())
}

func TestBasicDataTableToCommonDataTableWrapper_HasColumn(t *testing.T) {
	columns := []string{"Col1", "Col2", "Col3"}
	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          []*testBasicDataTableRow{},
	}

	commonDataTable := CreateNewCommonDataTableFromBasicDataTable(basicDataTable)

	assert.True(t, commonDataTable.HasColumn("Col1"))
	assert.True(t, commonDataTable.HasColumn("Col2"))
	assert.True(t, commonDataTable.HasColumn("Col3"))

	assert.False(t, commonDataTable.HasColumn("Col4"))
	assert.False(t, commonDataTable.HasColumn(""))
}

func TestBasicDataTableToCommonDataTableWrapper_DataRowCount(t *testing.T) {
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

	commonDataTable := CreateNewCommonDataTableFromBasicDataTable(basicDataTable)
	assert.Equal(t, len(rows), commonDataTable.DataRowCount())
}

func TestBasicDataTableToCommonDataTableWrapper_DataRowIterator(t *testing.T) {
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

	commonDataTable := CreateNewCommonDataTableFromBasicDataTable(basicDataTable)
	iterator := commonDataTable.DataRowIterator()

	assert.True(t, iterator.HasNext())
	firstRow := iterator.Next()
	assert.NotNil(t, firstRow)
	assert.Equal(t, len(columns), firstRow.ColumnCount())
	assert.True(t, firstRow.HasData("Col1"))
	assert.True(t, firstRow.HasData("Col2"))
	assert.True(t, firstRow.HasData("Col3"))
	assert.Equal(t, "A1", firstRow.GetData("Col1"))
	assert.Equal(t, "B1", firstRow.GetData("Col2"))
	assert.Equal(t, "C1", firstRow.GetData("Col3"))

	assert.True(t, iterator.HasNext())
	secondRow := iterator.Next()
	assert.NotNil(t, secondRow)
	assert.Equal(t, len(columns), secondRow.ColumnCount())
	assert.True(t, secondRow.HasData("Col1"))
	assert.True(t, secondRow.HasData("Col2"))
	assert.True(t, secondRow.HasData("Col3"))
	assert.Equal(t, "A2", secondRow.GetData("Col1"))
	assert.Equal(t, "B2", secondRow.GetData("Col2"))
	assert.Equal(t, "C2", secondRow.GetData("Col3"))

	assert.False(t, iterator.HasNext())
	assert.Nil(t, iterator.Next())
}
