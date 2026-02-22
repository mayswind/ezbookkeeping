package custom

import (
	"os"
	"testing"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestIsCustomExcelFileType(t *testing.T) {
	assert.True(t, IsCustomExcelFileType("custom_xlsx"))
	assert.True(t, IsCustomExcelFileType("custom_xls"))

	assert.False(t, IsCustomExcelFileType("xlsx"))
	assert.False(t, IsCustomExcelFileType("xls"))
	assert.False(t, IsCustomExcelFileType("excel"))
}

func TestCustomTransactionDataParser_ParseOOXMLExcelDataLines_EmptyData(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xlsx")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/empty_excel_file.xlsx")
	assert.Nil(t, err)

	allLines, err := importer.ParseDataLines(context, testdata)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(allLines))
}

func TestCustomTransactionDataParser_ParseOOXMLExcelDataLines_SingleSheet(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xlsx")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xlsx")
	assert.Nil(t, err)

	allLines, err := importer.ParseDataLines(context, testdata)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(allLines))

	assert.Equal(t, 3, len(allLines[0]))
	assert.Equal(t, "A1", allLines[0][0])
	assert.Equal(t, "B1", allLines[0][1])
	assert.Equal(t, "C1", allLines[0][2])

	assert.Equal(t, 3, len(allLines[1]))
	assert.Equal(t, "A2", allLines[1][0])
	assert.Equal(t, "B2", allLines[1][1])
	assert.Equal(t, "C2", allLines[1][2])

	assert.Equal(t, 3, len(allLines[2]))
	assert.Equal(t, "A3", allLines[2][0])
	assert.Equal(t, "B3", allLines[2][1])
	assert.Equal(t, "C3", allLines[2][2])
}

func TestCustomTransactionDataParser_ParseOOXMLExcelDataLines_MultipleSheet(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xlsx")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_excel_file.xlsx")
	assert.Nil(t, err)

	allLines, err := importer.ParseDataLines(context, testdata)
	assert.Nil(t, err)

	assert.Equal(t, 9, len(allLines))

	assert.Equal(t, 3, len(allLines[0]))
	assert.Equal(t, "A1", allLines[0][0])
	assert.Equal(t, "B1", allLines[0][1])
	assert.Equal(t, "C1", allLines[0][2])

	assert.Equal(t, 3, len(allLines[1]))
	assert.Equal(t, "1-A2", allLines[1][0])
	assert.Equal(t, "1-B2", allLines[1][1])
	assert.Equal(t, "1-C2", allLines[1][2])

	assert.Equal(t, 3, len(allLines[2]))
	assert.Equal(t, "1-A3", allLines[2][0])
	assert.Equal(t, "1-B3", allLines[2][1])
	assert.Equal(t, "1-C3", allLines[2][2])

	assert.Equal(t, 3, len(allLines[3]))
	assert.Equal(t, "A1", allLines[3][0])
	assert.Equal(t, "B1", allLines[3][1])
	assert.Equal(t, "C1", allLines[3][2])

	assert.Equal(t, 2, len(allLines[4]))
	assert.Equal(t, "3-A2", allLines[4][0])
	assert.Equal(t, "3-B2", allLines[4][1])

	assert.Equal(t, 3, len(allLines[5]))
	assert.Equal(t, "A1", allLines[5][0])
	assert.Equal(t, "B1", allLines[5][1])
	assert.Equal(t, "C1", allLines[5][2])

	assert.Equal(t, 3, len(allLines[6]))
	assert.Equal(t, "A1", allLines[6][0])
	assert.Equal(t, "B1", allLines[6][1])
	assert.Equal(t, "C1", allLines[6][2])

	assert.Equal(t, 3, len(allLines[7]))
	assert.Equal(t, "5-A2", allLines[7][0])
	assert.Equal(t, "5-B2", allLines[7][1])
	assert.Equal(t, "5-C2", allLines[7][2])

	assert.Equal(t, 3, len(allLines[8]))
	assert.Equal(t, "5-A3", allLines[8][0])
	assert.Equal(t, "5-B3", allLines[8][1])
	assert.Equal(t, "5-C3", allLines[8][2])
}

func TestCustomTransactionDataParser_ParseOOXMLExcelDataLines_MultipleSheetWithDifferentColumnCount(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xlsx")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_with_different_header_row_excel_file.xlsx")
	assert.Nil(t, err)

	_, err = importer.ParseDataLines(context, testdata)
	assert.EqualError(t, err, errs.ErrFieldsInMultiTableAreDifferent.Message)
}

func TestCustomTransactionDataParser_ParseMSCFBExcelDataLines_EmptyData(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xls")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/empty_excel_file.xls")
	assert.Nil(t, err)

	allLines, err := importer.ParseDataLines(context, testdata)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(allLines))
}

func TestCustomTransactionDataParser_ParseMSCFBExcelDataLines_SingleSheet(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xls")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/simple_excel_file.xls")
	assert.Nil(t, err)

	allLines, err := importer.ParseDataLines(context, testdata)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(allLines))

	assert.Equal(t, 3, len(allLines[0]))
	assert.Equal(t, "A1", allLines[0][0])
	assert.Equal(t, "B1", allLines[0][1])
	assert.Equal(t, "C1", allLines[0][2])

	assert.Equal(t, 3, len(allLines[1]))
	assert.Equal(t, "A2", allLines[1][0])
	assert.Equal(t, "B2", allLines[1][1])
	assert.Equal(t, "C2", allLines[1][2])

	assert.Equal(t, 3, len(allLines[2]))
	assert.Equal(t, "A3", allLines[2][0])
	assert.Equal(t, "B3", allLines[2][1])
	assert.Equal(t, "C3", allLines[2][2])
}

func TestCustomTransactionDataParser_ParseMSCFBExcelDataLines_MultipleSheet(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xls")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_excel_file.xls")
	assert.Nil(t, err)

	allLines, err := importer.ParseDataLines(context, testdata)
	assert.Nil(t, err)

	assert.Equal(t, 9, len(allLines))

	assert.Equal(t, 3, len(allLines[0]))
	assert.Equal(t, "A1", allLines[0][0])
	assert.Equal(t, "B1", allLines[0][1])
	assert.Equal(t, "C1", allLines[0][2])

	assert.Equal(t, 3, len(allLines[1]))
	assert.Equal(t, "1-A2", allLines[1][0])
	assert.Equal(t, "1-B2", allLines[1][1])
	assert.Equal(t, "1-C2", allLines[1][2])

	assert.Equal(t, 3, len(allLines[2]))
	assert.Equal(t, "1-A3", allLines[2][0])
	assert.Equal(t, "1-B3", allLines[2][1])
	assert.Equal(t, "1-C3", allLines[2][2])

	assert.Equal(t, 3, len(allLines[3]))
	assert.Equal(t, "A1", allLines[3][0])
	assert.Equal(t, "B1", allLines[3][1])
	assert.Equal(t, "C1", allLines[3][2])

	assert.Equal(t, 3, len(allLines[4]))
	assert.Equal(t, "3-A2", allLines[4][0])
	assert.Equal(t, "3-B2", allLines[4][1])
	assert.Equal(t, "", allLines[4][2])

	assert.Equal(t, 3, len(allLines[5]))
	assert.Equal(t, "A1", allLines[5][0])
	assert.Equal(t, "B1", allLines[5][1])
	assert.Equal(t, "C1", allLines[5][2])

	assert.Equal(t, 3, len(allLines[6]))
	assert.Equal(t, "A1", allLines[6][0])
	assert.Equal(t, "B1", allLines[6][1])
	assert.Equal(t, "C1", allLines[6][2])

	assert.Equal(t, 3, len(allLines[7]))
	assert.Equal(t, "5-A2", allLines[7][0])
	assert.Equal(t, "5-B2", allLines[7][1])
	assert.Equal(t, "5-C2", allLines[7][2])

	assert.Equal(t, 3, len(allLines[8]))
	assert.Equal(t, "5-A3", allLines[8][0])
	assert.Equal(t, "5-B3", allLines[8][1])
	assert.Equal(t, "5-C3", allLines[8][2])
}

func TestCustomTransactionDataParser_ParseMSCFBExcelDataLines_MultipleSheetWithDifferentColumnCount(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataExcelFileParser("custom_xls")
	assert.Nil(t, err)

	context := core.NewNullContext()

	testdata, err := os.ReadFile("../../../testdata/multiple_sheets_with_different_header_row_excel_file.xls")
	assert.Nil(t, err)

	_, err = importer.ParseDataLines(context, testdata)
	assert.EqualError(t, err, errs.ErrFieldsInMultiTableAreDifferent.Message)
}
