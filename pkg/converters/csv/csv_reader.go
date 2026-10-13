package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// CsvReader defines the structure of csv file reader
type CsvReader struct {
	reader    io.Reader
	separator rune
	trimSpace bool
}

// ReadAllLines reads all lines from the csv file
func (r *CsvReader) ReadAllLines(ctx core.Context) ([]*datatable.PlainTextDataTableLineData, error) {
	return r.read(ctx)
}

// ReadToCsvFileBasicDataTable returns a CsvFileBasicDataTable instance with the data read from the csv file
func (r *CsvReader) ReadToCsvFileBasicDataTable(ctx core.Context, hasTitleLine bool) (datatable.BasicDataTable, error) {
	allLines, err := r.read(ctx)

	if err != nil {
		return nil, err
	}

	return datatable.CreateNewPlainTextDataTableWithLineNumbers(allLines, hasTitleLine), nil
}

func (r *CsvReader) read(ctx core.Context) ([]*datatable.PlainTextDataTableLineData, error) {
	csvReader := csv.NewReader(r.reader)
	csvReader.Comma = r.separator
	csvReader.FieldsPerRecord = -1

	allLines := make([]*datatable.PlainTextDataTableLineData, 0)

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[csv_reader.read] cannot parse csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if len(items) == 1 && items[0] == "" {
			continue
		}

		if r.trimSpace {
			for index := range items {
				items[index] = strings.Trim(items[index], " ")
			}
		}

		lineNumber, _ := csvReader.FieldPos(0)
		allLines = append(allLines, &datatable.PlainTextDataTableLineData{
			LineNumber: lineNumber,
			Items:      items,
		})
	}

	return allLines, nil
}

// CreateNewCsvReader creates a new CsvReader with io reader
func CreateNewCsvReader(reader io.Reader, separator rune, trimSpace bool) *CsvReader {
	return &CsvReader{
		reader:    reader,
		separator: separator,
		trimSpace: trimSpace,
	}
}

// CreateNewCsvReaderFromData creates a new CsvReader with byte data and encoding
func CreateNewCsvReaderFromData(data []byte, transformer transform.Transformer, separator rune, trimSpace bool) *CsvReader {
	reader := transform.NewReader(bytes.NewReader(data), transformer)
	return CreateNewCsvReader(reader, separator, trimSpace)
}
