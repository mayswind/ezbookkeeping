package dsv

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	csvconverter "github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var supportedFileTypeSeparators = map[string]rune{
	"custom_csv": ',',
	"custom_tsv": '\t',
}

var supportedFileEncodings = map[string]encoding.Encoding{
	"utf-8":        unicode.UTF8,                                           // UTF-8
	"utf-8-bom":    unicode.UTF8BOM,                                        // UTF-8 with BOM
	"utf-16le":     unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), // UTF-16 Little Endian
	"utf-16be":     unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),    // UTF-16 Big Endian
	"cp437":        charmap.CodePage437,                                    // OEM United States (CP-437)
	"cp863":        charmap.CodePage863,                                    // OEM Canadian French (CP-863)
	"cp037":        charmap.CodePage037,                                    // IBM EBCDIC US/Canada (CP-037)
	"cp1047":       charmap.CodePage1047,                                   // IBM EBCDIC Open Systems (CP-1047)
	"cp1140":       charmap.CodePage1140,                                   // IBM EBCDIC US/Canada with Euro (CP-1140)
	"iso-8859-1":   charmap.ISO8859_1,                                      // Western European (ISO-8859-1)
	"cp850":        charmap.CodePage850,                                    // Western European (CP-850)
	"cp858":        charmap.CodePage858,                                    // Western European with Euro (CP-858)
	"windows-1252": charmap.Windows1252,                                    // Western European (Windows-1252)
	"iso-8859-15":  charmap.ISO8859_15,                                     // Western European (ISO-8859-15)
	"iso-8859-4":   charmap.ISO8859_4,                                      // North European (ISO-8859-4)
	"iso-8859-10":  charmap.ISO8859_10,                                     // North European (ISO-8859-10)
	"cp865":        charmap.CodePage865,                                    // North European (CP-865)
	"iso-8859-2":   charmap.ISO8859_2,                                      // Central European (ISO-8859-2)
	"cp852":        charmap.CodePage852,                                    // Central European (CP-852)
	"windows-1250": charmap.Windows1250,                                    // Central European (Windows-1250)
	"iso-8859-14":  charmap.ISO8859_14,                                     // Celtic (ISO-8859-14)
	"iso-8859-3":   charmap.ISO8859_3,                                      // South European (ISO-8859-3)
	"cp860":        charmap.CodePage860,                                    // Portuguese (CP-860)
	"iso-8859-7":   charmap.ISO8859_7,                                      // Greek (ISO-8859-7)
	"windows-1253": charmap.Windows1253,                                    // Greek (Windows-1253)
	"iso-8859-9":   charmap.ISO8859_9,                                      // Turkish (ISO-8859-9)
	"windows-1254": charmap.Windows1254,                                    // Turkish (Windows-1254)
	"iso-8859-13":  charmap.ISO8859_13,                                     // Baltic (ISO-8859-13)
	"windows-1257": charmap.Windows1257,                                    // Baltic (Windows-1257)
	"iso-8859-16":  charmap.ISO8859_16,                                     // South-Eastern European (ISO-8859-16)
	"iso-8859-5":   charmap.ISO8859_5,                                      // Cyrillic (ISO-8859-5)
	"cp855":        charmap.CodePage855,                                    // Cyrillic (CP-855)
	"cp866":        charmap.CodePage866,                                    // Cyrillic (CP-866)
	"windows-1251": charmap.Windows1251,                                    // Cyrillic (Windows-1251)
	"koi8r":        charmap.KOI8R,                                          // Cyrillic (KOI8-R)
	"koi8u":        charmap.KOI8U,                                          // Cyrillic (KOI8-U)
	"iso-8859-6":   charmap.ISO8859_6,                                      // Arabic (ISO-8859-6)
	"windows-1256": charmap.Windows1256,                                    // Arabic (Windows-1256)
	"iso-8859-8":   charmap.ISO8859_8,                                      // Hebrew (ISO-8859-8)
	"cp862":        charmap.CodePage862,                                    // Hebrew (CP-862)
	"windows-1255": charmap.Windows1255,                                    // Hebrew (Windows-1255)
	"windows-874":  charmap.Windows874,                                     // Thai (Windows-874)
	"windows-1258": charmap.Windows1258,                                    // Vietnamese (Windows-1258)
	"gb18030":      simplifiedchinese.GB18030,                              // Chinese (Simplified, GB18030)
	"gbk":          simplifiedchinese.GBK,                                  // Chinese (Simplified, GBK)
	"big5":         traditionalchinese.Big5,                                // Chinese (Traditional, Big5)
	"euc-kr":       korean.EUCKR,                                           // Korean (EUC-KR)
	"euc-jp":       japanese.EUCJP,                                         // Japanese (EUC-JP)
	"iso-2022-jp":  japanese.ISO2022JP,                                     // Japanese (ISO-2022-JP)
	"shift_jis":    japanese.ShiftJIS,                                      // Japanese (Shift JIS)
}

var customTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: utils.IntToString(int(models.TRANSACTION_TYPE_MODIFY_BALANCE)),
	models.TRANSACTION_TYPE_INCOME:         utils.IntToString(int(models.TRANSACTION_TYPE_INCOME)),
	models.TRANSACTION_TYPE_EXPENSE:        utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE)),
	models.TRANSACTION_TYPE_TRANSFER:       utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER)),
}

type CustomTransactionDataDsvFileParser interface {
	ParseDsvFileLines(ctx core.Context, data []byte) ([][]string, error)
}

// customTransactionDataDsvFileImporter defines the structure of custom dsv importer for transaction data
type customTransactionDataDsvFileImporter struct {
	fileEncoding               encoding.Encoding
	separator                  rune
	columnIndexMapping         map[datatable.TransactionDataTableColumn]int
	transactionTypeNameMapping map[string]models.TransactionType
	hasHeaderLine              bool
	timeFormat                 string
	timezoneFormat             string
	amountDecimalSeparator     string
	amountDigitGroupingSymbol  string
	geoLocationSeparator       string
	transactionTagSeparator    string
}

// ParseDsvFileLines returns the parsed file lines for specified the dsv file data
func (c *customTransactionDataDsvFileImporter) ParseDsvFileLines(ctx core.Context, data []byte) ([][]string, error) {
	reader := transform.NewReader(bytes.NewReader(data), c.fileEncoding.NewDecoder())
	csvReader := csv.NewReader(reader)
	csvReader.Comma = c.separator
	csvReader.FieldsPerRecord = -1

	allLines := make([][]string, 0)

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[custom_transaction_data_dsv_file_importer.ParseDsvFileLines] cannot parse dsv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if len(items) == 1 && items[0] == "" {
			continue
		}

		for index := range items {
			items[index] = strings.Trim(items[index], " ")
		}

		allLines = append(allLines, items)
	}

	return allLines, nil
}

// ParseImportedData returns the imported data by parsing the custom transaction dsv data
func (c *customTransactionDataDsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	allLines, err := c.ParseDsvFileLines(ctx, data)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if !c.hasHeaderLine {
		allLines = append([][]string{{}}, allLines...)
	}

	dataTable := csvconverter.CreateNewCustomCsvImportedDataTable(allLines)
	transactionDataTable := CreateNewCustomPlainTextDataTable(dataTable, c.columnIndexMapping, c.transactionTypeNameMapping, c.timeFormat, c.timezoneFormat, c.amountDecimalSeparator, c.amountDigitGroupingSymbol)
	dataTableImporter := converter.CreateNewImporterWithTypeNameMapping(customTransactionTypeNameMapping, c.geoLocationSeparator, c.transactionTagSeparator)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

// IsDelimiterSeparatedValuesFileType returns whether the file type is the delimiter-separated values file type
func IsDelimiterSeparatedValuesFileType(fileType string) bool {
	_, exists := supportedFileTypeSeparators[fileType]
	return exists
}

// CreateNewCustomTransactionDataDsvFileParser returns a new custom dsv parser for transaction data
func CreateNewCustomTransactionDataDsvFileParser(fileType string, fileEncoding string) (CustomTransactionDataDsvFileParser, error) {
	separator, exists := supportedFileTypeSeparators[fileType]

	if !exists {
		return nil, errs.ErrImportFileTypeNotSupported
	}

	enc, exists := supportedFileEncodings[fileEncoding]

	if !exists {
		return nil, errs.ErrImportFileEncodingNotSupported
	}

	return &customTransactionDataDsvFileImporter{
		fileEncoding: enc,
		separator:    separator,
	}, nil
}

// CreateNewCustomTransactionDataDsvFileImporter returns a new custom dsv importer for transaction data
func CreateNewCustomTransactionDataDsvFileImporter(fileType string, fileEncoding string, columnIndexMapping map[datatable.TransactionDataTableColumn]int, transactionTypeNameMapping map[string]models.TransactionType, hasHeaderLine bool, timeFormat string, timezoneFormat string, amountDecimalSeparator string, amountDigitGroupingSymbol string, geoLocationSeparator string, transactionTagSeparator string) (converter.TransactionDataImporter, error) {
	separator, exists := supportedFileTypeSeparators[fileType]

	if !exists {
		return nil, errs.ErrImportFileTypeNotSupported
	}

	enc, exists := supportedFileEncodings[fileEncoding]

	if !exists {
		return nil, errs.ErrImportFileEncodingNotSupported
	}

	if _, exists = columnIndexMapping[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME]; !exists {
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if _, exists = columnIndexMapping[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE]; !exists {
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if _, exists = columnIndexMapping[datatable.TRANSACTION_DATA_TABLE_AMOUNT]; !exists {
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	return &customTransactionDataDsvFileImporter{
		fileEncoding:               enc,
		separator:                  separator,
		columnIndexMapping:         columnIndexMapping,
		transactionTypeNameMapping: transactionTypeNameMapping,
		hasHeaderLine:              hasHeaderLine,
		timeFormat:                 timeFormat,
		timezoneFormat:             timezoneFormat,
		amountDecimalSeparator:     amountDecimalSeparator,
		amountDigitGroupingSymbol:  amountDigitGroupingSymbol,
		geoLocationSeparator:       geoLocationSeparator,
		transactionTagSeparator:    transactionTagSeparator,
	}, nil
}
