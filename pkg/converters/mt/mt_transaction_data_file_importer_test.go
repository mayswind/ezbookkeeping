package mt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestMT940TransactionDataFileParseImportedData_MinimumValidData(t *testing.T) {
	importer := MT940TransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := importer.ParseImportedData(context, user, []byte(
		`{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:
		:20:123456789
		:25:12345678
		:28C:123/1
		:60F:C250601CNY123,45
		:61:2506010602C123,45NTRFTEST//ABC123456
		:86:Transaction 1
		:61:2506020603D234,56NTRFFOOBAR
		:86:Transaction 2
		:62F:C250601CNY123,45
		-}`), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 0, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, int64(1748736000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "12345678", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "CNY", allNewTransactions[0].OriginalSourceAccountCurrency)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, int64(1748822400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(23456), allNewTransactions[1].Amount)
	assert.Equal(t, "12345678", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "CNY", allNewTransactions[1].OriginalSourceAccountCurrency)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "12345678", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)
}

func TestMT940TransactionDataFileParseImportedData_ParseTransactionValidAmountAndCurrency(t *testing.T) {
	importer := MT940TransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		`{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:
		:20:123456789
		:25:12345678
		:28C:123/1
		:60F:C250601CNY123,45
		:61:250601C123,45NTRFTEST
		:86:Transaction 1
		:61:250602C0,12NTRFTEST
		:86:Transaction 2
		:61:250603C1,NTRFTEST
		:86:Transaction 3
		:62F:C250601CNY123,45
		-}`), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(allNewTransactions))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
}

func TestMT940TransactionDataFileParseImportedData_ParseTransactionInvalidAmountAndCurrency(t *testing.T) {
	importer := MT940TransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		`{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:
		:20:123456789
		:25:12345678
		:28C:123/1
		:60F:C250601CNY123,45
		:61:2506010602C123 45NTRFTEST
		:62F:C250601CNY123,45
		-}`), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidMT940File.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		`{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:
		:20:123456789
		:25:12345678
		:28C:123/1
		:60F:C250601CNY123,45
		:61:2506010602C12.45NTRFTEST
		:62F:C250601CNY123,45
		-}`), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidMT940File.Message)
}

func TestMT940TransactionDataFileParseImportedData_ParseTransactionType(t *testing.T) {
	importer := MT940TransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		`{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:
		:20:123456789
		:25:12345678
		:28C:123/1
		:60F:C250601CNY123,45
		:61:250601C123,45NTRFTEST
		:86:Transaction 1
		:61:250602D123,45NTRFTEST
		:86:Transaction 2
		:61:250603RC123,45NTRFTEST
		:86:Transaction 3
		:61:250604RD123,45NTRFTEST
		:86:Transaction 4
		:62F:C250601CNY123,45
		-}`), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[3].Type)
}

func TestMT940TransactionDataFileParseImportedData_ParseDescription(t *testing.T) {
	importer := MT940TransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		`{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:
		:20:123456789
		:25:12345678
		:28C:123/1
		:60F:C250601CNY123,45
		:61:2506010602C123,45NTRFTEST
		:86:Transaction 1
		Part 2
		Part 3
		:62F:C250601CNY123,45
		-}`), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Transaction 1\nPart 2\nPart 3", allNewTransactions[0].Comment)
}

func TestMT940TransactionDataFileParseImportedData_MissingRequiredField(t *testing.T) {
	importer := MT940TransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// Missing opening balance and closing balance
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		`{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:
		:20:123456789
		:28C:123/1
		:61:250601C123,45NTRFTEST
		:86:Transaction 1
		-}`), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}
