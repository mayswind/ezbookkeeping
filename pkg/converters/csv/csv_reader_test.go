package csv

import (
	"testing"

	"golang.org/x/text/encoding/unicode"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewCsvReaderFromData(t *testing.T) {
	context := core.NewNullContext()
	csvReader := CreateNewCsvReaderFromData([]byte("A1,B1,C1\n"+
		"A2,B2,C2\n"+
		"A3,B3,C3\n"), unicode.UTF8.NewDecoder(), ',', false)
	datatable, err := csvReader.ReadToCsvFileBasicDataTable(context, true)
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

func TestCreateNewCsvReaderFromData_SkipBlankLine(t *testing.T) {
	context := core.NewNullContext()
	csvReader := CreateNewCsvReaderFromData([]byte("\n"+
		"A1,B1,C1\n"+
		"A2,B2,C2\n"+
		"\n"+
		"A3,B3,C3\n"), unicode.UTF8.NewDecoder(), ',', false)
	datatable, err := csvReader.ReadToCsvFileBasicDataTable(context, true)
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
