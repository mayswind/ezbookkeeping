package excel

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestExcelFileImportedDataTableDataRowCount(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	assert.Nil(t, err)
	assert.Equal(t, 2, datatable.DataRowCount())
}

func TestExcelFileImportedDataTableDataRowCount_MultipleSheets(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	assert.Nil(t, err)
	assert.Equal(t, 5, datatable.DataRowCount())
}

func TestExcelFileImportedDataTableDataRowCount_OnlyHeaderLine(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/only_one_row_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	assert.Nil(t, err)
	assert.Equal(t, 0, datatable.DataRowCount())
}

func TestExcelFileImportedDataTableDataRowCount_EmptyContent(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/empty_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	assert.Nil(t, err)
	assert.Equal(t, 0, datatable.DataRowCount())
}

func TestExcelFileImportedDataTableHeaderColumnNames(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	assert.EqualValues(t, []string{"A1", "B1", "C1"}, datatable.HeaderColumnNames())
}

func TestExcelFileImportedDataTableHeaderColumnNames_EmptyContent(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/empty_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	assert.Nil(t, datatable.HeaderColumnNames())
}

func TestExcelFileDataRowIterator(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
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

func TestExcelFileDataRowIterator_MultipleSheets(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	iterator := datatable.DataRowIterator()
	assert.True(t, iterator.HasNext())

	// sheet 1 data row 1
	assert.NotNil(t, iterator.Next())
	assert.True(t, iterator.HasNext())

	// sheet 1 data row 2
	assert.NotNil(t, iterator.Next())
	assert.True(t, iterator.HasNext())

	// sheet 3 data row 1
	assert.NotNil(t, iterator.Next())
	assert.True(t, iterator.HasNext())

	// sheet 5 data row 1
	assert.NotNil(t, iterator.Next())
	assert.True(t, iterator.HasNext())

	// sheet 5 data row 2
	assert.NotNil(t, iterator.Next())
	assert.False(t, iterator.HasNext())

	// not existed data row
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())

	// not existed data row
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())
}

func TestExcelFileDataRowIterator_OnlyHeaderLine(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/only_one_row_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	iterator := datatable.DataRowIterator()
	assert.False(t, iterator.HasNext())

	// not existed data row 1
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())

	// not existed data row 2
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())
}

func TestExcelFileDataRowIterator_EmptyContent(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/empty_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	iterator := datatable.DataRowIterator()
	assert.False(t, iterator.HasNext())

	// not existed data row 1
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())

	// not existed data row 2
	assert.Nil(t, iterator.Next())
	assert.False(t, iterator.HasNext())
}

func TestExcelFileDataRowColumnCount(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.EqualValues(t, 4, row1.ColumnCount())

	row2 := iterator.Next()
	assert.EqualValues(t, 4, row2.ColumnCount())
}

func TestExcelFileDataRowGetData(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
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

func TestExcelFileDataRowGetData_GetNotExistedColumnData(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	iterator := datatable.DataRowIterator()

	row1 := iterator.Next()
	assert.Equal(t, "", row1.GetData(3))
}

func TestExcelFileDataRowGetData_MultipleSheets(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_excel_file.xls")
	assert.Nil(t, err)

	datatable, err := CreateNewExcelFileImportedDataTable(testdata)
	iterator := datatable.DataRowIterator()

	sheet1Row1 := iterator.Next()
	assert.Equal(t, "1-A2", sheet1Row1.GetData(0))
	assert.Equal(t, "1-B2", sheet1Row1.GetData(1))
	assert.Equal(t, "1-C2", sheet1Row1.GetData(2))

	sheet1Row2 := iterator.Next()
	assert.Equal(t, "1-A3", sheet1Row2.GetData(0))
	assert.Equal(t, "1-B3", sheet1Row2.GetData(1))
	assert.Equal(t, "1-C3", sheet1Row2.GetData(2))

	// skip empty sheet2

	sheet3Row1 := iterator.Next()
	assert.Equal(t, "3-A2", sheet3Row1.GetData(0))
	assert.Equal(t, "3-B2", sheet3Row1.GetData(1))
	assert.Equal(t, "", sheet3Row1.GetData(2))

	// skip no data row sheet4

	sheet5Row1 := iterator.Next()
	assert.Equal(t, "5-A2", sheet5Row1.GetData(0))
	assert.Equal(t, "5-B2", sheet5Row1.GetData(1))
	assert.Equal(t, "5-C2", sheet5Row1.GetData(2))

	sheet5Row2 := iterator.Next()
	assert.Equal(t, "5-A3", sheet5Row2.GetData(0))
	assert.Equal(t, "5-B3", sheet5Row2.GetData(1))
	assert.Equal(t, "5-C3", sheet5Row2.GetData(2))
}

func TestCreateNewExcelFileImportedDataTable_MultipleSheetsWithDifferentHeaders(t *testing.T) {
	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_with_different_header_row_excel_file.xls")
	assert.Nil(t, err)

	_, err = CreateNewExcelFileImportedDataTable(testdata)
	assert.EqualError(t, err, errs.ErrFieldsInMultiTableAreDifferent.Message)
}
