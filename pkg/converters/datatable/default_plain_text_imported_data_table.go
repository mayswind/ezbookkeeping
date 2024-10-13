package datatable

import (
	"encoding/csv"
	"io"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// DefaultPlainTextImportedDataTable defines the structure of default plain text data table
type DefaultPlainTextImportedDataTable struct {
	allLines [][]string
}

// DefaultPlainTextImportedDataRow defines the structure of default plain text data table row
type DefaultPlainTextImportedDataRow struct {
	dataTable *DefaultPlainTextImportedDataTable
	allItems  []string
}

// DefaultPlainTextImportedDataRowIterator defines the structure of default plain text data table row iterator
type DefaultPlainTextImportedDataRowIterator struct {
	dataTable    *DefaultPlainTextImportedDataTable
	currentIndex int
}

// DataRowCount returns the total count of data row
func (t *DefaultPlainTextImportedDataTable) DataRowCount() int {
	if len(t.allLines) < 1 {
		return 0
	}

	return len(t.allLines) - 1
}

// HeaderColumnNames returns the header column name list
func (t *DefaultPlainTextImportedDataTable) HeaderColumnNames() []string {
	if len(t.allLines) < 1 {
		return nil
	}

	return t.allLines[0]
}

// DataRowIterator returns the iterator of data row
func (t *DefaultPlainTextImportedDataTable) DataRowIterator() ImportedDataRowIterator {
	return &DefaultPlainTextImportedDataRowIterator{
		dataTable:    t,
		currentIndex: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *DefaultPlainTextImportedDataRow) ColumnCount() int {
	return len(r.allItems)
}

// GetData returns the data in the specified column index
func (r *DefaultPlainTextImportedDataRow) GetData(columnIndex int) string {
	if columnIndex >= len(r.allItems) {
		return ""
	}

	return r.allItems[columnIndex]
}

// HasNext returns whether the iterator does not reach the end
func (t *DefaultPlainTextImportedDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allLines)
}

// Next returns the next imported data row
func (t *DefaultPlainTextImportedDataRowIterator) Next() ImportedDataRow {
	if t.currentIndex+1 >= len(t.dataTable.allLines) {
		return nil
	}

	t.currentIndex++

	rowItems := t.dataTable.allLines[t.currentIndex]

	return &DefaultPlainTextImportedDataRow{
		dataTable: t.dataTable,
		allItems:  rowItems,
	}
}

// CreateNewDefaultCsvDataTable returns default csv data table by io readers
func CreateNewDefaultCsvDataTable(ctx core.Context, reader io.Reader) (*DefaultPlainTextImportedDataTable, error) {
	return createNewDefaultPlainTextDataTable(ctx, reader, ',')
}

func createNewDefaultPlainTextDataTable(ctx core.Context, reader io.Reader, comma rune) (*DefaultPlainTextImportedDataTable, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = comma
	csvReader.FieldsPerRecord = -1

	allLines := make([][]string, 0)

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[Default_plain_text_imported_data_table.createNewDefaultPlainTextDataTable] cannot parse plain text data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if len(items) == 0 && items[0] == "" {
			continue
		}

		allLines = append(allLines, items)
	}

	return &DefaultPlainTextImportedDataTable{
		allLines: allLines,
	}, nil
}
