package dsv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestIsDelimiterSeparatedValuesFileType(t *testing.T) {
	assert.True(t, IsDelimiterSeparatedValuesFileType("custom_csv"))
	assert.True(t, IsDelimiterSeparatedValuesFileType("custom_tsv"))
	assert.True(t, IsDelimiterSeparatedValuesFileType("custom_ssv"))

	assert.False(t, IsDelimiterSeparatedValuesFileType("dsv"))
	assert.False(t, IsDelimiterSeparatedValuesFileType("csv"))
	assert.False(t, IsDelimiterSeparatedValuesFileType("tsv"))
	assert.False(t, IsDelimiterSeparatedValuesFileType("ssv"))
}

func TestCustomTransactionDataDsvFileParser_ParseDsvFileLines(t *testing.T) {
	importer, err := CreateNewCustomTransactionDataDsvFileParser("custom_csv", "utf-8")
	assert.Nil(t, err)

	context := core.NewNullContext()

	allLines, err := importer.ParseDsvFileLines(context, []byte(
		"2024-09-01 00:00:00,B,123.45\n"+
			"2024-09-01 01:23:45,I,0.12\n"))
	assert.Nil(t, err)

	assert.Equal(t, 2, len(allLines))

	assert.Equal(t, 3, len(allLines[0]))
	assert.Equal(t, "2024-09-01 00:00:00", allLines[0][0])
	assert.Equal(t, "B", allLines[0][1])
	assert.Equal(t, "123.45", allLines[0][2])

	assert.Equal(t, 3, len(allLines[1]))
	assert.Equal(t, "2024-09-01 01:23:45", allLines[1][0])
	assert.Equal(t, "I", allLines[1][1])
	assert.Equal(t, "0.12", allLines[1][2])

	importer, err = CreateNewCustomTransactionDataDsvFileParser("custom_tsv", "utf-8")
	assert.Nil(t, err)

	allLines, err = importer.ParseDsvFileLines(context, []byte(
		"2024-09-01 12:34:56\tE\t1.00\n"+
			"2024-09-01 23:59:59\tT\t0.05"))
	assert.Nil(t, err)

	assert.Equal(t, 2, len(allLines))

	assert.Equal(t, 3, len(allLines[0]))
	assert.Equal(t, "2024-09-01 12:34:56", allLines[0][0])
	assert.Equal(t, "E", allLines[0][1])
	assert.Equal(t, "1.00", allLines[0][2])

	assert.Equal(t, 3, len(allLines[1]))
	assert.Equal(t, "2024-09-01 23:59:59", allLines[1][0])
	assert.Equal(t, "T", allLines[1][1])
	assert.Equal(t, "0.05", allLines[1][2])

	importer, err = CreateNewCustomTransactionDataDsvFileParser("custom_ssv", "utf-8")
	assert.Nil(t, err)

	allLines, err = importer.ParseDsvFileLines(context, []byte(
		"2024-09-01 12:34:56;E;1.00\n"+
			"2024-09-01 23:59:59;T;0.05"))
	assert.Nil(t, err)

	assert.Equal(t, 2, len(allLines))

	assert.Equal(t, 3, len(allLines[0]))
	assert.Equal(t, "2024-09-01 12:34:56", allLines[0][0])
	assert.Equal(t, "E", allLines[0][1])
	assert.Equal(t, "1.00", allLines[0][2])

	assert.Equal(t, 3, len(allLines[1]))
	assert.Equal(t, "2024-09-01 23:59:59", allLines[1][0])
	assert.Equal(t, "T", allLines[1][1])
	assert.Equal(t, "0.05", allLines[1][2])
}

func TestCustomTransactionDataDsvFileImporter_MinimumValidData(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"I": models.TRANSACTION_TYPE_INCOME,
		"E": models.TRANSACTION_TYPE_EXPENSE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", ".", "", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,B,123.45\n"+
			"2024-09-01 01:23:45,I,0.12\n"+
			"2024-09-01 12:34:56,E,1.00\n"+
			"2024-09-01 23:59:59,T,0.05"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725153825), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725235199), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}

func TestCustomTransactionDataDsvFileImporter_WithAllSupportedColumns(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE:     1,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         2,
		datatable.TRANSACTION_DATA_TABLE_CATEGORY:                 3,
		datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:             4,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             5,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         6,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   7,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     8,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: 9,
		datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           10,
		datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION:      11,
		datatable.TRANSACTION_DATA_TABLE_TAGS:                     12,
		datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:              13,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"Balance Modification": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"Income":               models.TRANSACTION_TYPE_INCOME,
		"Expense":              models.TRANSACTION_TYPE_EXPENSE,
		"Transfer":             models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, true, "YYYY-MM-DD HH:mm:ss", "", ".", "", " ", "lonlat", ";")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := importer.ParseImportedData(context, user, []byte(
		"\"Time\",\"Timezone\",\"Type\",\"Category\",\"Sub Category\",\"Account\",\"Account Currency\",\"Amount\",\"Account2\",\"Account2 Currency\",\"Account2 Amount\",\"Geographic Location\",\"Tags\",\"Description\"\n"+
			"\"2024-09-01 00:00:00\",\"+08:00\",\"Balance Modification\",\"\",\"\",\"Test Account\",\"CNY\",\"123.45\",\"\",\"\",\"\",\"\",\"\",\"\"\n"+
			"\"2024-09-01 01:23:45\",\"+08:00\",\"Income\",\"Test Category\",\"Test Sub Category\",\"Test Account\",\"CNY\",\"0.12\",\"\",\"\",\"\",\"123.450000 45.670000\",\"Test Tag;Test Tag2\",\"Hello World\"\n"+
			"\"2024-09-01 12:34:56\",\"+00:00\",\"Expense\",\"Test Category2\",\"Test Sub Category2\",\"Test Account\",\"CNY\",\"1.00\",\"\",\"\",\"\",\"\",\"Test Tag\",\"Foo#Bar\"\n"+
			"\"2024-09-01 23:59:59\",\"-05:00\",\"Transfer\",\"Test Category3\",\"Test Sub Category3\",\"Test Account\",\"CNY\",\"0.05\",\"Test Account2\",\"USD\",\"0.35\",\"\",\"Test Tag2\",\"foo\tbar\""), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 2, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725120000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725125025), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Test Sub Category", allNewTransactions[1].OriginalCategoryName)
	assert.Equal(t, 123.45, allNewTransactions[1].GeoLongitude)
	assert.Equal(t, 45.67, allNewTransactions[1].GeoLatitude)
	assert.Equal(t, "Hello World", allNewTransactions[1].Comment)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Sub Category2", allNewTransactions[2].OriginalCategoryName)
	assert.Equal(t, "Foo#Bar", allNewTransactions[2].Comment)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725253199), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "Test Sub Category3", allNewTransactions[3].OriginalCategoryName)
	assert.Equal(t, "foo\tbar", allNewTransactions[3].Comment)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[1].Name)
	assert.Equal(t, "USD", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Test Sub Category2", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Test Sub Category", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "Test Sub Category3", allNewSubTransferCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewTags[0].Uid)
	assert.Equal(t, "Test Tag", allNewTags[0].Name)

	assert.Equal(t, int64(1234567890), allNewTags[1].Uid)
	assert.Equal(t, "Test Tag2", allNewTags[1].Name)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidTime(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01T12:34:56,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"09/01/2024 12:34:56,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseTransactionWithoutType(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"I": models.TRANSACTION_TYPE_INCOME,
		"E": models.TRANSACTION_TYPE_EXPENSE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,A,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidType(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": 0,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,B,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTypeInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseTimeWithTimezone(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ssZ", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56-10:00,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725230096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56+00:00,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56+12:45,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725148196), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
}

func TestCustomTransactionDataDsvFileImporter_ParseTimeWithTimezone2(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ssZZ", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56-1000,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725230096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56+0000,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56+1245,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725148196), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
}

func TestCustomTransactionDataDsvFileImporter_ParseValidTimezone(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE: 1,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:               3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,-10:00,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725230096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,+00:00,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,+12:45,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725148196), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
}

func TestCustomTransactionDataDsvFileImporter_ParseValidTimezone2(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE: 1,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:               3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "ZZ", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,-1000,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725230096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,+0000,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,+1245,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725148196), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidTimezoneFormat(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE: 1,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:               3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "z", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,CST,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrImportFileTransactionTimezoneFormatInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidTimezone(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE: 1,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:               3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,Asia/Shanghai,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeZoneInvalid.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,-0700,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeZoneInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidTimezone2(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE: 1,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:               3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "ZZ", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,Asia/Shanghai,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeZoneInvalid.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,0700,E,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeZoneInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseAmountWithCustomFormat(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_tsv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ",", ".", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56\tE\t1.234,56"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(123456), allNewTransactions[0].Amount)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidAmountWithCustomFormat(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_tsv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", ",", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56\tE\t1.234,56"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidAmountWithCustomFormat2(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_tsv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ",", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56\tE\t1.234,56"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParsePrimaryCategory(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_CATEGORY:         2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"I": models.TRANSACTION_TYPE_INCOME,
		"E": models.TRANSACTION_TYPE_EXPENSE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,B,,123.45\n"+
			"2024-09-01 01:23:45,I,Test Category,0.12\n"+
			"2024-09-01 12:34:56,E,Test Category2,1.00\n"+
			"2024-09-01 23:59:59,T,Test Category3,0.05"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725153825), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, "Test Category", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, "Test Category2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725235199), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Category3", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "Test Category3", allNewSubTransferCategories[0].Name)
}

func TestCustomTransactionDataDsvFileImporter_ParseValidAccountCurrency(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         1,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             2,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         3,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   4,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     5,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: 6,
		datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           7,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 01:23:45,B,Test Account,USD,123.45,,,\n"+
			"2024-09-01 12:34:56,T,Test Account,USD,1.23,Test Account2,EUR,1.10"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "USD", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[1].Name)
	assert.Equal(t, "EUR", allNewAccounts[1].Currency)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidAccountCurrency(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         1,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             2,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         3,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   4,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     5,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: 6,
		datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           7,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 01:23:45,B,Test Account,USD,123.45,,,\n"+
			"2024-09-01 12:34:56,T,Test Account,CNY,1.23,Test Account2,EUR,1.10"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 01:23:45,B,Test Account,USD,123.45,,,\n"+
			"2024-09-01 12:34:56,T,Test Account2,CNY,1.23,Test Account,EUR,1.10"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseNotSupportedCurrency(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:         0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:         1,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:             2,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY:         3,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:                   4,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME:     5,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY: 6,
		datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:           7,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 01:23:45,B,Test Account,XXX,123.45,,,"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 01:23:45,T,Test Account,USD,123.45,Test Account2,XXX,123.45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseValidAmount(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
		datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:   3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
		"I": models.TRANSACTION_TYPE_INCOME,
		"E": models.TRANSACTION_TYPE_EXPENSE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,B,123.45000000,\n"+
			"2024-09-01 01:23:45,I,0.12000000,\n"+
			"2024-09-01 12:34:56,E,1.00000000,\n"+
			"2024-09-01 23:59:59,T,0.05000000,0.35000000"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, int64(1725153825), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, int64(1725235199), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, int64(35), allNewTransactions[3].RelatedAccountAmount)
}

func TestCustomTransactionDataDsvFileImporter_ParseAmountWithSpaceDigitGroupingSymbol(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", " ", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// normal space
	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,E,1 234,\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(123400), allNewTransactions[0].Amount)

	// no-break space (NBSP)
	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,E,1 234,\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(123400), allNewTransactions[0].Amount)

	// narrow no-break space (NNBSP)
	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,E,1 234,\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(123400), allNewTransactions[0].Amount)

	// figure space
	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,E,1 234,\n"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(123400), allNewTransactions[0].Amount)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidAmount(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     1,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:               3,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: 4,
		datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT:       5,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,E,Test Account,123 45,,"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,T,Test Account,123.45,Test Account2,123 45"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseNoAmount2(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     1,
		datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         2,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:               3,
		datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: 4,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,E,Test Account,123.45,"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, int64(0), allNewTransactions[0].RelatedAccountAmount)

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,T,Test Account,123.45,Test Account2"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, int64(12345), allNewTransactions[0].RelatedAccountAmount)
}

func TestCustomTransactionDataDsvFileImporter_ParseValidGeographicLocation(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:    0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:    1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:              2,
		datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION: 3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", ";", "lonlat", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,E,123.45,123.45;45.56"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 123.45, allNewTransactions[0].GeoLongitude)
	assert.Equal(t, 45.56, allNewTransactions[0].GeoLatitude)
}

func TestCustomTransactionDataDsvFileImporter_ParseInvalidGeographicLocation(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:    0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:    1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:              2,
		datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION: 3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", " ", "lonlat", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,E,123.45,,,1"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, float64(0), allNewTransactions[0].GeoLongitude)
	assert.Equal(t, float64(0), allNewTransactions[0].GeoLatitude)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,E,123.45,a b"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrGeographicLocationInvalid.Message)
}

func TestCustomTransactionDataDsvFileImporter_ParseTag(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
		datatable.TRANSACTION_DATA_TABLE_TAGS:             3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", ";")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, allNewTags, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,E,123.45,foo;;bar.;#test;hello\tworld;;"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTags[0].Uid)
	assert.Equal(t, "foo", allNewTags[0].Name)

	assert.Equal(t, int64(1234567890), allNewTags[1].Uid)
	assert.Equal(t, "bar.", allNewTags[1].Name)

	assert.Equal(t, int64(1234567890), allNewTags[2].Uid)
	assert.Equal(t, "#test", allNewTags[2].Name)

	assert.Equal(t, int64(1234567890), allNewTags[3].Uid)
	assert.Equal(t, "hello\tworld", allNewTags[3].Name)
}

func TestCustomTransactionDataDsvFileImporter_ParseTagWithoutSeparator(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
		datatable.TRANSACTION_DATA_TABLE_TAGS:             3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"E": models.TRANSACTION_TYPE_EXPENSE,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, allNewTags, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 00:00:00,E,123.45,foo;;bar.;#test;hello\tworld;;"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTags[0].Uid)
	assert.Equal(t, "foo;;bar.;#test;hello	world;;", allNewTags[0].Name)
}

func TestCustomTransactionDataDsvFileImporter_ParseDescription(t *testing.T) {
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
		datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:      3,
	}
	transactionTypeMapping := map[string]models.TransactionType{
		"T": models.TRANSACTION_TYPE_TRANSFER,
	}
	importer, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.Nil(t, err)

	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"2024-09-01 12:34:56,T,123.45,foo    bar\t#test"), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)
}

func TestCustomTransactionDataDsvFileImporter_InvalidSeparator(t *testing.T) {
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
	}
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	_, err := CreateNewCustomTransactionDataDsvFileImporter("test", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.EqualError(t, err, errs.ErrImportFileTypeNotSupported.Message)
}

func TestCustomTransactionDataDsvFileImporter_InvalidFileEncoding(t *testing.T) {
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
	}
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           2,
	}
	_, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "ascii", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.EqualError(t, err, errs.ErrImportFileEncodingNotSupported.Message)
}

func TestCustomTransactionDataDsvFileImporter_MissingRequiredColumn(t *testing.T) {
	transactionTypeMapping := map[string]models.TransactionType{
		"B": models.TRANSACTION_TYPE_MODIFY_BALANCE,
	}

	// Missing Time Column
	columnIndexMapping := map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 0,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           1,
	}
	_, err := CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Type Column
	columnIndexMapping = map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_AMOUNT:           1,
	}
	_, err = CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	columnIndexMapping = map[datatable.TransactionDataTableColumn]int{
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: 0,
		datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: 1,
	}
	_, err = CreateNewCustomTransactionDataDsvFileImporter("custom_csv", "utf-8", columnIndexMapping, transactionTypeMapping, false, "YYYY-MM-DD HH:mm:ss", "", ".", "", "", "", "")
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}
