package datatable

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

type testTransactionDataRowParser struct {
}

func (p *testTransactionDataRowParser) Parse(rowData map[TransactionDataTableColumn]string) (map[TransactionDataTableColumn]string, bool, error) {
	rowData[TRANSACTION_DATA_TABLE_DESCRIPTION] = "Test Description"
	return rowData, true, nil
}

func (p *testTransactionDataRowParser) GetAddedColumns() []TransactionDataTableColumn {
	return []TransactionDataTableColumn{TRANSACTION_DATA_TABLE_DESCRIPTION}
}

func TestBasicDataTableToTransactionDataTableWrapper_HasColumn(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          []*testBasicDataTableRow{},
	}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	transactionDataTable := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable, columnMapping)

	assert.True(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.True(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
	assert.True(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.False(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_DESCRIPTION))
	assert.False(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_CATEGORY))
}

func TestBasicDataTableToTransactionDataTableWrapper_TransactionRowCount(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"2024-01-01", "1", "100"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"2024-01-02", "2", "200"},
		},
		{
			rowId:      "3",
			rowColumns: []string{"2024-01-03", "1", "300"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	transactionDataTable := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable, columnMapping)
	assert.Equal(t, len(rows), transactionDataTable.TransactionRowCount())
}

func TestBasicDataTableToTransactionDataTableWrapper_TransactionRowIterator(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"2024-01-01", "1", "100"},
		},
		{
			rowId:      "2",
			rowColumns: []string{"2024-01-02", "2", "200"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	transactionDataTable := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable, columnMapping)
	iterator := transactionDataTable.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	firstRow, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, firstRow)
	assert.True(t, firstRow.IsValid())
	assert.Equal(t, "2024-01-01", firstRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "1", firstRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
	assert.Equal(t, "100", firstRow.GetData(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.True(t, iterator.HasNext())
	secondRow, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, secondRow)
	assert.True(t, secondRow.IsValid())
	assert.Equal(t, "2024-01-02", secondRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TIME))
	assert.Equal(t, "2", secondRow.GetData(TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))
	assert.Equal(t, "200", secondRow.GetData(TRANSACTION_DATA_TABLE_AMOUNT))

	assert.False(t, iterator.HasNext())
	emptyRow, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, emptyRow)
}

func TestBasicDataTableToTransactionDataTableWrapper_TransactionRowIterator_EmptyRow(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{""},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	transactionDataTable := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable, columnMapping)
	iterator := transactionDataTable.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.False(t, row.IsValid())
}

func TestBasicDataTableToTransactionDataTableWrapper_TransactionRowIterator_InvalidRow(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"2024-01-01", "1"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
	}

	transactionDataTable := CreateNewTransactionDataTableFromBasicDataTable(basicDataTable, columnMapping)
	iterator := transactionDataTable.TransactionRowIterator()

	assert.True(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, errs.ErrFewerFieldsInDataRowThanInHeaderRow, err)
	assert.Nil(t, row)
}

func TestBasicDataTableToTransactionDataTableWrapper_TransactionRowIterator_WithRowParserAddedColumn(t *testing.T) {
	columns := []string{"TransactionTime", "TransactionType", "Amount"}
	rows := []*testBasicDataTableRow{
		{
			rowId:      "1",
			rowColumns: []string{"2024-01-01", "1", "100"},
		},
	}

	basicDataTable := &testBasicDataTable{
		headerColumns: columns,
		rows:          rows,
	}

	columnMapping := map[TransactionDataTableColumn]string{
		TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "TransactionTime",
		TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "TransactionType",
		TRANSACTION_DATA_TABLE_AMOUNT:           "Amount",
		TRANSACTION_DATA_TABLE_DESCRIPTION:      "Description",
	}

	transactionDataTable := CreateNewTransactionDataTableFromBasicDataTableWithRowParser(basicDataTable, columnMapping, &testTransactionDataRowParser{})
	assert.True(t, transactionDataTable.HasColumn(TRANSACTION_DATA_TABLE_DESCRIPTION))

	iterator := transactionDataTable.TransactionRowIterator()
	assert.True(t, iterator.HasNext())
	row, err := iterator.Next(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, row)
	assert.True(t, row.IsValid())
	assert.Equal(t, "Test Description", row.GetData(TRANSACTION_DATA_TABLE_DESCRIPTION))
}
