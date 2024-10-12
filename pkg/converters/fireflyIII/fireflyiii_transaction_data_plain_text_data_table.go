package fireflyIII

import (
	"encoding/csv"
	"io"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var fireflyIIITransactionSupportedColumns = []datatable.DataTableColumn{
	datatable.DATA_TABLE_TRANSACTION_TIME,
	datatable.DATA_TABLE_TRANSACTION_TYPE,
	datatable.DATA_TABLE_SUB_CATEGORY,
	datatable.DATA_TABLE_ACCOUNT_NAME,
	datatable.DATA_TABLE_ACCOUNT_CURRENCY,
	datatable.DATA_TABLE_AMOUNT,
	datatable.DATA_TABLE_RELATED_ACCOUNT_NAME,
	datatable.DATA_TABLE_RELATED_ACCOUNT_CURRENCY,
	datatable.DATA_TABLE_RELATED_AMOUNT,
	datatable.DATA_TABLE_TAGS,
	datatable.DATA_TABLE_DESCRIPTION,
}

var fireflyIIITransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_MODIFY_BALANCE: "Opening balance",
	models.TRANSACTION_TYPE_INCOME:         "Deposit",
	models.TRANSACTION_TYPE_EXPENSE:        "Withdrawal",
	models.TRANSACTION_TYPE_TRANSFER:       "Transfer",
}

// fireflyIIITransactionPlainTextDataTable defines the structure of firefly III transaction plain text data table
type fireflyIIITransactionPlainTextDataTable struct {
	allOriginalLines              [][]string
	originalHeaderLineColumnNames []string
	originalColumnIndex           map[datatable.DataTableColumn]int
}

// fireflyIIITransactionPlainTextDataRow defines the structure of firefly III transaction plain text data row
type fireflyIIITransactionPlainTextDataRow struct {
	dataTable     *fireflyIIITransactionPlainTextDataTable
	originalItems []string
	finalItems    map[datatable.DataTableColumn]string
}

// fireflyIIITransactionPlainTextDataRowIterator defines the structure of firefly III transaction plain text data row iterator
type fireflyIIITransactionPlainTextDataRowIterator struct {
	dataTable    *fireflyIIITransactionPlainTextDataTable
	currentIndex int
}

// DataRowCount returns the total count of data row
func (t *fireflyIIITransactionPlainTextDataTable) DataRowCount() int {
	if len(t.allOriginalLines) < 1 {
		return 0
	}

	return len(t.allOriginalLines) - 1
}

// GetDataColumnMapping returns data column map for data importer
func (t *fireflyIIITransactionPlainTextDataTable) GetDataColumnMapping() map[datatable.DataTableColumn]string {
	dataColumnMapping := make(map[datatable.DataTableColumn]string, len(fireflyIIITransactionSupportedColumns))

	for i := 0; i < len(fireflyIIITransactionSupportedColumns); i++ {
		column := fireflyIIITransactionSupportedColumns[i]
		dataColumnMapping[column] = utils.IntToString(int(column))
	}

	return dataColumnMapping
}

// HeaderLineColumnNames returns the header column name list
func (t *fireflyIIITransactionPlainTextDataTable) HeaderLineColumnNames() []string {
	columnIndexes := make([]string, len(fireflyIIITransactionSupportedColumns))

	for i := 0; i < len(fireflyIIITransactionSupportedColumns); i++ {
		column := fireflyIIITransactionSupportedColumns[i]

		if t.originalColumnIndex[column] >= 0 {
			columnIndexes[i] = utils.IntToString(int(column))
		} else {
			columnIndexes[i] = "-1"
		}
	}

	return columnIndexes
}

// DataRowIterator returns the iterator of data row
func (t *fireflyIIITransactionPlainTextDataTable) DataRowIterator() datatable.ImportedDataRowIterator {
	return &fireflyIIITransactionPlainTextDataRowIterator{
		dataTable:    t,
		currentIndex: 0,
	}
}

// IsValid returns whether this row contains valid data for importing
func (r *fireflyIIITransactionPlainTextDataRow) IsValid() bool {
	return true
}

// ColumnCount returns the total count of column in this data row
func (r *fireflyIIITransactionPlainTextDataRow) ColumnCount() int {
	return len(fireflyIIITransactionSupportedColumns)
}

// GetData returns the data in the specified column index
func (r *fireflyIIITransactionPlainTextDataRow) GetData(columnIndex int) string {
	if columnIndex >= len(fireflyIIITransactionSupportedColumns) {
		return ""
	}

	dataColumn := fireflyIIITransactionSupportedColumns[columnIndex]

	return r.finalItems[dataColumn]
}

// GetTime returns the time in the specified column index
func (r *fireflyIIITransactionPlainTextDataRow) GetTime(columnIndex int, timezoneOffset int16) (time.Time, error) {
	return utils.ParseFromLongDateTimeWithTimezone(r.GetData(columnIndex))
}

// GetTimezoneOffset returns the time zone offset in the specified column index
func (r *fireflyIIITransactionPlainTextDataRow) GetTimezoneOffset(columnIndex int) (*time.Location, error) {
	return nil, errs.ErrNotSupported
}

// HasNext returns whether the iterator does not reach the end
func (t *fireflyIIITransactionPlainTextDataRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allOriginalLines)
}

// Next returns the next imported data row
func (t *fireflyIIITransactionPlainTextDataRowIterator) Next(ctx core.Context, user *models.User) datatable.ImportedDataRow {
	if t.currentIndex+1 >= len(t.dataTable.allOriginalLines) {
		return nil
	}

	t.currentIndex++

	rowItems := t.dataTable.allOriginalLines[t.currentIndex]
	finalItems := t.dataTable.parseTransactionData(rowItems)

	return &fireflyIIITransactionPlainTextDataRow{
		dataTable:     t.dataTable,
		originalItems: rowItems,
		finalItems:    finalItems,
	}
}

func (t *fireflyIIITransactionPlainTextDataTable) parseTransactionData(items []string) map[datatable.DataTableColumn]string {
	data := make(map[datatable.DataTableColumn]string, 12)

	data[datatable.DATA_TABLE_SUB_CATEGORY] = ""

	for column, index := range t.originalColumnIndex {
		if index >= 0 && index < len(items) {
			data[column] = items[index]
		}
	}

	// trim trailing zero in decimal
	if data[datatable.DATA_TABLE_AMOUNT] != "" {
		data[datatable.DATA_TABLE_AMOUNT] = utils.TrimTrailingZerosInDecimal(data[datatable.DATA_TABLE_AMOUNT])
		amount, err := utils.ParseAmount(data[datatable.DATA_TABLE_AMOUNT])

		if err == nil {
			data[datatable.DATA_TABLE_AMOUNT] = utils.FormatAmount(-amount)
		}
	}

	if data[datatable.DATA_TABLE_RELATED_AMOUNT] != "" {
		data[datatable.DATA_TABLE_RELATED_AMOUNT] = utils.TrimTrailingZerosInDecimal(data[datatable.DATA_TABLE_RELATED_AMOUNT])
		amount, err := utils.ParseAmount(data[datatable.DATA_TABLE_RELATED_AMOUNT])

		if err == nil {
			data[datatable.DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(-amount)
		}
	} else {
		data[datatable.DATA_TABLE_RELATED_AMOUNT] = data[datatable.DATA_TABLE_AMOUNT]
	}

	// the related account currency field is foreign currency in firefly iii actually
	if data[datatable.DATA_TABLE_RELATED_ACCOUNT_CURRENCY] == "" {
		data[datatable.DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = data[datatable.DATA_TABLE_ACCOUNT_CURRENCY]
	}

	// the destination account of modify balance transaction in firefly iii is the asset account
	if data[datatable.DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE] {
		data[datatable.DATA_TABLE_ACCOUNT_NAME] = data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME]
	}

	// the destination account of income transaction in firefly iii is the asset account
	if data[datatable.DATA_TABLE_TRANSACTION_TYPE] == fireflyIIITransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME] {
		data[datatable.DATA_TABLE_ACCOUNT_NAME] = data[datatable.DATA_TABLE_RELATED_ACCOUNT_NAME]
	}

	return data
}

func createNewFireflyIIITransactionPlainTextDataTable(ctx core.Context, reader io.Reader) (*fireflyIIITransactionPlainTextDataTable, error) {
	allOriginalLines, err := parseAllLinesFromFireflyIIITransactionPlainText(ctx, reader)

	if err != nil {
		return nil, err
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[fireflyiii_transaction_data_plain_text_data_table.createNewFireflyIIITransactionPlainTextDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	originalHeaderItems := allOriginalLines[0]
	originalHeaderItemMap := make(map[string]int)

	for i := 0; i < len(originalHeaderItems); i++ {
		originalHeaderItemMap[originalHeaderItems[i]] = i
	}

	typeColumnIdx, typeColumnExists := originalHeaderItemMap["type"]
	amountColumnIdx, amountColumnExists := originalHeaderItemMap["amount"]
	foreignAmountColumnIdx, foreignAmountColumnExists := originalHeaderItemMap["foreign_amount"]
	currencyColumnIdx, currencyColumnExists := originalHeaderItemMap["currency_code"]
	foreignCurrencyColumnIdx, foreignCurrencyColumnExists := originalHeaderItemMap["foreign_currency_code"]
	descriptionColumnIdx, descriptionColumnExists := originalHeaderItemMap["description"]
	dateColumnIdx, dateColumnExists := originalHeaderItemMap["date"]
	sourceNameColumnIdx, sourceNameColumnExists := originalHeaderItemMap["source_name"]
	destinationNameColumnIdx, destinationNameColumnExists := originalHeaderItemMap["destination_name"]
	categoryColumnIdx, categoryColumnExists := originalHeaderItemMap["category"]
	tagsColumnIdx, tagsColumnExists := originalHeaderItemMap["tags"]

	if !typeColumnExists || !amountColumnExists || !dateColumnExists || !sourceNameColumnExists || !destinationNameColumnExists || !categoryColumnExists {
		log.Errorf(ctx, "[fireflyiii_transaction_data_plain_text_data_table.createNewFireflyIIITransactionPlainTextDataTable] cannot parse firefly III csv data, because missing essential columns in header row")
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if !foreignAmountColumnExists {
		foreignAmountColumnIdx = -1
	}

	if !currencyColumnExists {
		currencyColumnIdx = -1
	}

	if !foreignCurrencyColumnExists {
		foreignCurrencyColumnIdx = -1
	}

	if !descriptionColumnExists {
		descriptionColumnIdx = -1
	}

	if !tagsColumnExists {
		tagsColumnIdx = -1
	}

	return &fireflyIIITransactionPlainTextDataTable{
		allOriginalLines:              allOriginalLines,
		originalHeaderLineColumnNames: originalHeaderItems,
		originalColumnIndex: map[datatable.DataTableColumn]int{
			datatable.DATA_TABLE_TRANSACTION_TIME:         dateColumnIdx,
			datatable.DATA_TABLE_TRANSACTION_TYPE:         typeColumnIdx,
			datatable.DATA_TABLE_SUB_CATEGORY:             categoryColumnIdx,
			datatable.DATA_TABLE_ACCOUNT_NAME:             sourceNameColumnIdx,
			datatable.DATA_TABLE_ACCOUNT_CURRENCY:         currencyColumnIdx,
			datatable.DATA_TABLE_AMOUNT:                   amountColumnIdx,
			datatable.DATA_TABLE_RELATED_ACCOUNT_NAME:     destinationNameColumnIdx,
			datatable.DATA_TABLE_RELATED_ACCOUNT_CURRENCY: foreignCurrencyColumnIdx,
			datatable.DATA_TABLE_RELATED_AMOUNT:           foreignAmountColumnIdx,
			datatable.DATA_TABLE_TAGS:                     tagsColumnIdx,
			datatable.DATA_TABLE_DESCRIPTION:              descriptionColumnIdx,
		},
	}, nil
}

func parseAllLinesFromFireflyIIITransactionPlainText(ctx core.Context, reader io.Reader) ([][]string, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	allOriginalLines := make([][]string, 0)

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[fireflyiii_transaction_data_plain_text_data_table.parseAllLinesFromFireflyIIITransactionPlainText] cannot parse firefly III csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		allOriginalLines = append(allOriginalLines, items)
	}

	return allOriginalLines, nil
}
