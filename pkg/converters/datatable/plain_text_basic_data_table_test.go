package datatable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlainTextBasicDataTableDataRowCount(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, false)

	assert.Equal(t, 3, datatable.DataRowCount())
}

func TestPlainTextBasicDataTableDataRowCount_HasTitleLine(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, true)

	assert.Equal(t, 2, datatable.DataRowCount())
}

func TestPlainTextBasicDataTableDataRowCount_OnlyHeaderLine(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
	}, true)

	assert.Equal(t, 0, datatable.DataRowCount())
}

func TestPlainTextBasicDataTableDataRowCount_EmptyContent(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{}, false)
	assert.Equal(t, 0, datatable.DataRowCount())

	datatable = CreateNewPlainTextDataTableWithoutLineNumbers([][]string{}, true)
	assert.Equal(t, 0, datatable.DataRowCount())
}

func TestPlainTextBasicDataTableHeaderColumnNames(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, true)

	assert.EqualValues(t, []string{"A1", "B1", "C1"}, datatable.HeaderColumnNames())
}

func TestPlainTextBasicDataTableHeaderColumnNames_EmptyContent(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{}, false)
	assert.Nil(t, datatable.HeaderColumnNames())

	datatable = CreateNewPlainTextDataTableWithoutLineNumbers([][]string{}, true)
	assert.Nil(t, datatable.HeaderColumnNames())
}

func TestPlainTextBasicDataTableRowIterator(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, false)

	iterator := datatable.DataRowIterator()
	assert.True(t, iterator.HasNext())

	// data row 1
	assert.NotNil(t, iterator.Next())
	assert.True(t, iterator.HasNext())

	// data row 2
	assert.NotNil(t, iterator.Next())
	assert.True(t, iterator.HasNext())

	// data row 3
	assert.NotNil(t, iterator.Next())
	assert.False(t, iterator.HasNext())

	// not existed data row 4
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())
}

func TestPlainTextBasicDataTableRowIterator_HasTitleLine(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, true)

	iterator := datatable.DataRowIterator()
	assert.True(t, iterator.HasNext())

	// data row 1
	assert.NotNil(t, iterator.Next())
	assert.True(t, iterator.HasNext())

	// data row 2
	assert.NotNil(t, iterator.Next())
	assert.False(t, iterator.HasNext())

	// not existed data row 3
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())

	// not existed data row 4
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())
}

func TestPlainTextBasicDataTableRowIteratorCurrentRowId_WithLineNumbers(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithLineNumbers([]*PlainTextDataTableLineData{
		{LineNumber: 1, Items: []string{"A1", "B1"}},
		{LineNumber: 3, Items: []string{"A2", "B2"}},
		{LineNumber: 6, Items: []string{"A3", "B3"}},
	}, true)

	iterator := datatable.DataRowIterator()

	iterator.Next()
	assert.Equal(t, "line:3", iterator.CurrentRowId())

	iterator.Next()
	assert.Equal(t, "line:6", iterator.CurrentRowId())
}

func TestPlainTextBasicDataTableRowIteratorCurrentRowId_WithoutLineNumbers(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, false)

	iterator := datatable.DataRowIterator()
	assert.Equal(t, "", iterator.CurrentRowId())

	iterator.Next()
	assert.Equal(t, "record:1", iterator.CurrentRowId())

	iterator.Next()
	assert.Equal(t, "record:2", iterator.CurrentRowId())

	iterator.Next()
	assert.Equal(t, "record:3", iterator.CurrentRowId())

	iterator.Next()
	assert.Equal(t, "record:3", iterator.CurrentRowId())
}

func TestPlainTextBasicDataTableRowIteratorCurrentRowId_WithoutLineNumbers_HasTitleLine(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, true)

	iterator := datatable.DataRowIterator()

	iterator.Next()
	assert.Equal(t, "record:1", iterator.CurrentRowId())

	iterator.Next()
	assert.Equal(t, "record:2", iterator.CurrentRowId())
}

func TestPlainTextBasicDataTableRowIteratorCurrentRowId_EmptyContent(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{}, false)
	iterator := datatable.DataRowIterator()
	assert.Equal(t, "", iterator.CurrentRowId())

	datatable = CreateNewPlainTextDataTableWithoutLineNumbers([][]string{}, true)
	iterator = datatable.DataRowIterator()
	assert.Equal(t, "", iterator.CurrentRowId())
}

func TestPlainTextBasicDataTableRowColumnCount(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, true)

	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.EqualValues(t, 3, row1.ColumnCount())

	row2 := iterator.Next()
	assert.EqualValues(t, 3, row2.ColumnCount())
}

func TestPlainTextBasicDataTableRowGetData(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, false)

	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.Equal(t, "A1", row1.GetData(0))
	assert.Equal(t, "B1", row1.GetData(1))
	assert.Equal(t, "C1", row1.GetData(2))

	row2 := iterator.Next()
	assert.Equal(t, "A2", row2.GetData(0))
	assert.Equal(t, "B2", row2.GetData(1))
	assert.Equal(t, "C2", row2.GetData(2))

	row3 := iterator.Next()
	assert.Equal(t, "A3", row3.GetData(0))
	assert.Equal(t, "B3", row3.GetData(1))
	assert.Equal(t, "C3", row3.GetData(2))
}

func TestPlainTextBasicDataTableRowGetData_HasTitleLine(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, true)

	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.Equal(t, "A2", row1.GetData(0))
	assert.Equal(t, "B2", row1.GetData(1))
	assert.Equal(t, "C2", row1.GetData(2))

	row2 := iterator.Next()
	assert.Equal(t, "A3", row2.GetData(0))
	assert.Equal(t, "B3", row2.GetData(1))
	assert.Equal(t, "C3", row2.GetData(2))
}

func TestPlainTextBasicDataTableRowGetData_GetNotExistedColumnData(t *testing.T) {
	datatable := CreateNewPlainTextDataTableWithoutLineNumbers([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	}, true)

	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.Equal(t, "", row1.GetData(3))
}
