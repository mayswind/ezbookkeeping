package _default

import (
	"fmt"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// ezBookKeepingTransactionPlainTextDataTable defines the structure of ezbookkeeping transaction plain text data table
type ezBookKeepingTransactionPlainTextDataTable struct {
	columnSeparator       string
	lineSeparator         string
	allLines              []string
	headerLineColumnNames []string
}

// ezBookKeepingTransactionPlainTextDataRow defines the structure of ezbookkeeping transaction plain text data row
type ezBookKeepingTransactionPlainTextDataRow struct {
	allItems []string
}

// ezBookKeepingTransactionPlainTextDataRowIterator defines the structure of ezbookkeeping transaction plain text data row iterator
type ezBookKeepingTransactionPlainTextDataRowIterator struct {
	dataTable    *ezBookKeepingTransactionPlainTextDataTable
	currentIndex int
}

// ezBookKeepingTransactionPlainTextDataTableBuilder defines the structure of ezbookkeeping transaction plain text data table builder
type ezBookKeepingTransactionPlainTextDataTableBuilder struct {
	columnSeparator       string
	lineSeparator         string
	columns               []datatable.DataTableColumn
	dataColumnNameMapping map[datatable.DataTableColumn]string
	dataLineFormat        string
	builder               *strings.Builder
}

// DataRowCount returns the total count of data row
func (t *ezBookKeepingTransactionPlainTextDataTable) DataRowCount() int {
	if len(t.allLines) < 1 {
		return 0
	}

	return len(t.allLines) - 1
}

// HeaderLineColumnNames returns the header column name list
func (t *ezBookKeepingTransactionPlainTextDataTable) HeaderLineColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *ezBookKeepingTransactionPlainTextDataTable) DataRowIterator() datatable.ImportedDataRowIterator {
	return &ezBookKeepingTransactionPlainTextDataRowIterator{
		dataTable:    t,
		currentIndex: 0,
	}
}

// IsValid returns whether this row contains valid data for importing
func (r *ezBookKeepingTransactionPlainTextDataRow) IsValid() bool {
	return true
}

// ColumnCount returns the total count of column in this data row
func (r *ezBookKeepingTransactionPlainTextDataRow) ColumnCount() int {
	return len(r.allItems)
}

// GetData returns the data in the specified column index
func (r *ezBookKeepingTransactionPlainTextDataRow) GetData(columnIndex int) string {
	if columnIndex >= len(r.allItems) {
		return ""
	}

	return r.allItems[columnIndex]
}

// GetTime returns the time in the specified column index
func (r *ezBookKeepingTransactionPlainTextDataRow) GetTime(columnIndex int, timezoneOffset int16) (time.Time, error) {
	return utils.ParseFromLongDateTime(r.GetData(columnIndex), timezoneOffset)
}

// GetTimezoneOffset returns the time zone offset in the specified column index
func (r *ezBookKeepingTransactionPlainTextDataRow) GetTimezoneOffset(columnIndex int) (*time.Location, error) {
	return utils.ParseFromTimezoneOffset(r.GetData(columnIndex))
}

// HasNext returns whether the iterator does not reach the end
func (t *ezBookKeepingTransactionPlainTextDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allLines)
}

// Next returns the next imported data row
func (t *ezBookKeepingTransactionPlainTextDataRowIterator) Next(ctx core.Context, user *models.User) datatable.ImportedDataRow {
	if t.currentIndex+1 >= len(t.dataTable.allLines) {
		return nil
	}

	t.currentIndex++

	rowContent := t.dataTable.allLines[t.currentIndex]
	rowItems := strings.Split(rowContent, t.dataTable.columnSeparator)

	return &ezBookKeepingTransactionPlainTextDataRow{
		allItems: rowItems,
	}
}

// AppendTransaction appends the specified transaction to data builder
func (b *ezBookKeepingTransactionPlainTextDataTableBuilder) AppendTransaction(data map[datatable.DataTableColumn]string) {
	dataRowParams := make([]any, len(b.columns))

	for i := 0; i < len(b.columns); i++ {
		dataRowParams[i] = data[b.columns[i]]
	}

	b.builder.WriteString(fmt.Sprintf(b.dataLineFormat, dataRowParams...))
}

// ReplaceDelimiters returns the text after removing the delimiters
func (b *ezBookKeepingTransactionPlainTextDataTableBuilder) ReplaceDelimiters(text string) string {
	text = strings.Replace(text, "\r\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, b.columnSeparator, " ", -1)
	text = strings.Replace(text, b.lineSeparator, " ", -1)

	return text
}

// String returns the textual representation of this data
func (b *ezBookKeepingTransactionPlainTextDataTableBuilder) String() string {
	return b.builder.String()
}

func (b *ezBookKeepingTransactionPlainTextDataTableBuilder) generateHeaderLine() string {
	var ret strings.Builder

	for i := 0; i < len(b.columns); i++ {
		if ret.Len() > 0 {
			ret.WriteString(b.columnSeparator)
		}

		dataColumn := b.columns[i]
		columnName := b.dataColumnNameMapping[dataColumn]

		ret.WriteString(columnName)
	}

	ret.WriteString(b.lineSeparator)

	return ret.String()
}

func (b *ezBookKeepingTransactionPlainTextDataTableBuilder) generateDataLineFormat() string {
	var ret strings.Builder

	for i := 0; i < len(b.columns); i++ {
		if ret.Len() > 0 {
			ret.WriteString(b.columnSeparator)
		}

		ret.WriteString("%s")
	}

	ret.WriteString(b.lineSeparator)

	return ret.String()
}

func createNewezbookkeepingTransactionPlainTextDataTable(content string, columnSeparator string, lineSeparator string) (*ezBookKeepingTransactionPlainTextDataTable, error) {
	allLines := strings.Split(content, lineSeparator)

	if len(allLines) < 2 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	headerLine := allLines[0]
	headerLine = strings.ReplaceAll(headerLine, "\r", "")
	headerLineItems := strings.Split(headerLine, columnSeparator)

	return &ezBookKeepingTransactionPlainTextDataTable{
		columnSeparator:       columnSeparator,
		lineSeparator:         lineSeparator,
		allLines:              allLines,
		headerLineColumnNames: headerLineItems,
	}, nil
}

func createNewezbookkeepingTransactionPlainTextDataTableBuilder(transactionCount int, columns []datatable.DataTableColumn, dataColumnNameMapping map[datatable.DataTableColumn]string, columnSeparator string, lineSeparator string) *ezBookKeepingTransactionPlainTextDataTableBuilder {
	var builder strings.Builder
	builder.Grow(transactionCount * 100)

	dataTableBuilder := &ezBookKeepingTransactionPlainTextDataTableBuilder{
		columnSeparator:       columnSeparator,
		lineSeparator:         lineSeparator,
		columns:               columns,
		dataColumnNameMapping: dataColumnNameMapping,
		builder:               &builder,
	}

	headerLine := dataTableBuilder.generateHeaderLine()
	dataLineFormat := dataTableBuilder.generateDataLineFormat()

	dataTableBuilder.builder.WriteString(headerLine)
	dataTableBuilder.dataLineFormat = dataLineFormat

	return dataTableBuilder
}
