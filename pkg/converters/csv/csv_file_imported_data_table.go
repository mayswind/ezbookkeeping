package csv

import (
	"encoding/csv"
	"io"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// CsvFileImportedDataTable defines the structure of csv data table
type CsvFileImportedDataTable struct {
	allLines [][]string
}

// CsvFileImportedDataRow defines the structure of csv data table row
type CsvFileImportedDataRow struct {
	dataTable *CsvFileImportedDataTable
	allItems  []string
}

// CsvFileImportedDataRowIterator defines the structure of csv data table row iterator
type CsvFileImportedDataRowIterator struct {
	dataTable    *CsvFileImportedDataTable
	currentIndex int
}

// DataRowCount returns the total count of data row
func (t *CsvFileImportedDataTable) DataRowCount() int {
	if len(t.allLines) < 1 {
		return 0
	}

	return len(t.allLines) - 1
}

// HeaderColumnNames returns the header column name list
func (t *CsvFileImportedDataTable) HeaderColumnNames() []string {
	if len(t.allLines) < 1 {
		return nil
	}

	return t.allLines[0]
}

// DataRowIterator returns the iterator of data row
func (t *CsvFileImportedDataTable) DataRowIterator() datatable.ImportedDataRowIterator {
	return &CsvFileImportedDataRowIterator{
		dataTable:    t,
		currentIndex: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *CsvFileImportedDataRow) ColumnCount() int {
	return len(r.allItems)
}

// GetData returns the data in the specified column index
func (r *CsvFileImportedDataRow) GetData(columnIndex int) string {
	if columnIndex >= len(r.allItems) {
		return ""
	}

	return r.allItems[columnIndex]
}

// HasNext returns whether the iterator does not reach the end
func (t *CsvFileImportedDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allLines)
}

// Next returns the next imported data row
func (t *CsvFileImportedDataRowIterator) Next() datatable.ImportedDataRow {
	if t.currentIndex+1 >= len(t.dataTable.allLines) {
		return nil
	}

	t.currentIndex++

	rowItems := t.dataTable.allLines[t.currentIndex]

	return &CsvFileImportedDataRow{
		dataTable: t.dataTable,
		allItems:  rowItems,
	}
}

// CreateNewCsvDataTable returns comma separated values data table by io readers
func CreateNewCsvDataTable(ctx core.Context, reader io.Reader) (*CsvFileImportedDataTable, error) {
	return createNewCsvFileDataTable(ctx, reader, ',')
}

func createNewCsvFileDataTable(ctx core.Context, reader io.Reader, separator rune) (*CsvFileImportedDataTable, error) {
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
			log.Errorf(ctx, "[csv_file_imported_data_table.createNewCsvFileDataTable] cannot parse csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if len(items) == 0 && items[0] == "" {
			continue
		}

		allLines = append(allLines, items)
	}

	return &CsvFileImportedDataTable{
		allLines: allLines,
	}, nil
}
