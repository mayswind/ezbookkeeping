package csv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func TestCsvFileImportedDataTableDataRowCount(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	})

	assert.Equal(t, 2, datatable.DataRowCount())
}

func TestCsvFileImportedDataTableDataRowCount_OnlyHeaderLine(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{
		{"A1", "B1", "C1"},
	})

	assert.Equal(t, 0, datatable.DataRowCount())
}

func TestCsvFileImportedDataTableDataRowCount_EmptyContent(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{})

	assert.Equal(t, 0, datatable.DataRowCount())
}

func TestCsvFileImportedDataTableHeaderColumnNames(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	})

	assert.EqualValues(t, []string{"A1", "B1", "C1"}, datatable.HeaderColumnNames())
}

func TestCsvFileImportedDataTableHeaderColumnNames_EmptyContent(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{})

	assert.Nil(t, datatable.HeaderColumnNames())
}

func TestCsvFileImportedDataRowIterator(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	})

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

func TestCsvFileImportedDataRowColumnCount(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	})

	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.EqualValues(t, 3, row1.ColumnCount())

	row2 := iterator.Next()
	assert.EqualValues(t, 3, row2.ColumnCount())
}

func TestCsvFileImportedDataRowGetData(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	})

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

func TestCsvFileImportedDataRowGetData_GetNotExistedColumnData(t *testing.T) {
	datatable := CreateNewCustomCsvImportedDataTable([][]string{
		{"A1", "B1", "C1"},
		{"A2", "B2", "C2"},
		{"A3", "B3", "C3"},
	})

	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.Equal(t, "", row1.GetData(3))
}

func TestCreateNewCsvImportedDataTable(t *testing.T) {
	context := core.NewNullContext()
	reader := bytes.NewReader([]byte("A1,B1,C1\n" +
		"A2,B2,C2\n" +
		"A3,B3,C3\n"))
	datatable, err := CreateNewCsvImportedDataTable(context, reader)
	assert.Nil(t, err)

	assert.Equal(t, 2, datatable.DataRowCount())

	iterator := datatable.DataRowIterator()
	assert.True(t, iterator.HasNext())

	row1 := iterator.Next()
	assert.EqualValues(t, 3, row1.ColumnCount())
	assert.Equal(t, "A2", row1.GetData(0))
	assert.Equal(t, "B2", row1.GetData(1))
	assert.Equal(t, "C2", row1.GetData(2))
	assert.True(t, iterator.HasNext())

	row2 := iterator.Next()
	assert.EqualValues(t, 3, row2.ColumnCount())
	assert.Equal(t, "A3", row2.GetData(0))
	assert.Equal(t, "B3", row2.GetData(1))
	assert.Equal(t, "C3", row2.GetData(2))
	assert.False(t, iterator.HasNext())
}

func TestCreateNewCsvImportedDataTable_SkipBlankLine(t *testing.T) {
	context := core.NewNullContext()
	reader := bytes.NewReader([]byte("\n" +
		"A1,B1,C1\n" +
		"A2,B2,C2\n" +
		"\n" +
		"A3,B3,C3\n"))
	datatable, err := CreateNewCsvImportedDataTable(context, reader)
	assert.Nil(t, err)

	assert.Equal(t, 2, datatable.DataRowCount())

	iterator := datatable.DataRowIterator()
	assert.True(t, iterator.HasNext())

	row1 := iterator.Next()
	assert.EqualValues(t, 3, row1.ColumnCount())
	assert.Equal(t, "A2", row1.GetData(0))
	assert.Equal(t, "B2", row1.GetData(1))
	assert.Equal(t, "C2", row1.GetData(2))
	assert.True(t, iterator.HasNext())

	row2 := iterator.Next()
	assert.EqualValues(t, 3, row2.ColumnCount())
	assert.Equal(t, "A3", row2.GetData(0))
	assert.Equal(t, "B3", row2.GetData(1))
	assert.Equal(t, "C3", row2.GetData(2))
	assert.False(t, iterator.HasNext())
}
