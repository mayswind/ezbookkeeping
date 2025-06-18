package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// CsvFileBasicDataTable defines the structure of csv data table
type CsvFileBasicDataTable struct {
	allLines [][]string
}

// CsvFileBasicDataTableRow defines the structure of csv data table row
type CsvFileBasicDataTableRow struct {
	dataTable *CsvFileBasicDataTable
	allItems  []string
}

// CsvFileBasicDataTableRowIterator defines the structure of csv data table row iterator
type CsvFileBasicDataTableRowIterator struct {
	dataTable    *CsvFileBasicDataTable
	currentIndex int
}

// DataRowCount returns the total count of data row
func (t *CsvFileBasicDataTable) DataRowCount() int {
	if len(t.allLines) < 1 {
		return 0
	}

	return len(t.allLines) - 1
}

// HeaderColumnNames returns the header column name list
func (t *CsvFileBasicDataTable) HeaderColumnNames() []string {
	if len(t.allLines) < 1 {
		return nil
	}

	return t.allLines[0]
}

// DataRowIterator returns the iterator of data row
func (t *CsvFileBasicDataTable) DataRowIterator() datatable.BasicDataTableRowIterator {
	return &CsvFileBasicDataTableRowIterator{
		dataTable:    t,
		currentIndex: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *CsvFileBasicDataTableRow) ColumnCount() int {
	return len(r.allItems)
}

// GetData returns the data in the specified column index
func (r *CsvFileBasicDataTableRow) GetData(columnIndex int) string {
	if columnIndex >= len(r.allItems) {
		return ""
	}

	return r.allItems[columnIndex]
}

// HasNext returns whether the iterator does not reach the end
func (t *CsvFileBasicDataTableRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allLines)
}

// CurrentRowId returns current index
func (t *CsvFileBasicDataTableRowIterator) CurrentRowId() string {
	return fmt.Sprintf("line#%d", t.currentIndex)
}

// Next returns the next basic data row
func (t *CsvFileBasicDataTableRowIterator) Next() datatable.BasicDataTableRow {
	if t.currentIndex+1 >= len(t.dataTable.allLines) {
		return nil
	}

	t.currentIndex++

	rowItems := t.dataTable.allLines[t.currentIndex]

	return &CsvFileBasicDataTableRow{
		dataTable: t.dataTable,
		allItems:  rowItems,
	}
}

// CreateNewCsvBasicDataTable returns comma separated values data table by io readers
func CreateNewCsvBasicDataTable(ctx core.Context, reader io.Reader) (datatable.BasicDataTable, error) {
	return createNewCsvFileBasicDataTable(ctx, reader, ',')
}

// CreateNewCustomCsvBasicDataTable returns character separated values data table by io readers
func CreateNewCustomCsvBasicDataTable(allLines [][]string) datatable.BasicDataTable {
	return &CsvFileBasicDataTable{
		allLines: allLines,
	}
}

func createNewCsvFileBasicDataTable(ctx core.Context, reader io.Reader, separator rune) (*CsvFileBasicDataTable, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = separator
	csvReader.FieldsPerRecord = -1

	allLines := make([][]string, 0)

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[csv_file_basic_data_table.createNewCsvFileDataTable] cannot parse csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if len(items) == 1 && items[0] == "" {
			continue
		}

		allLines = append(allLines, items)
	}

	return &CsvFileBasicDataTable{
		allLines: allLines,
	}, nil
}
