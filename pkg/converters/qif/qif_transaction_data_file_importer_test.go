package qif

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestQIFTransactionDataFileParseImportedData_MinimumValidData(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T123.45\n"+
			"POpening Balance\n"+
			"L[Test Account]\n"+
			"^\n"+
			"D2024-09-02\n"+
			"T0.12\n"+
			"LTest Category\n"+
			"^\n"+
			"D2024-09-03\n"+
			"T-1.00\n"+
			"LTest Category2\n"+
			"^\n"+
			"D2024-09-04\n"+
			"T-0.05\n"+
			"L[Test Account2]\n"+
			"^\n"+
			"D2024-09-05\n"+
			"T0.06\n"+
			"L[Test Account2]\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 5, len(allNewTransactions))
	assert.Equal(t, 3, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[4].Type)
	assert.Equal(t, int64(1725494400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
	assert.Equal(t, int64(6), allNewTransactions[4].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[4].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[2].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[2].Name)
	assert.Equal(t, "CNY", allNewAccounts[2].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}

func TestQIFTransactionDataFileParseImportedData_ParseYearMonthDayDateFormatTime(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"^\n"+
			"D2024-9-2\n"+
			"T-123.45\n"+
			"^\n"+
			"D2024/9/3\n"+
			"T-123.45\n"+
			"^\n"+
			"D2024.9.4\n"+
			"T-123.45\n"+
			"^\n"+
			"D2024'9.5\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 5, len(allNewTransactions))

	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(1725494400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
}

func TestQIFTransactionDataFileParseImportedData_ParseMonthDayYearDateFormatTime(t *testing.T) {
	importer := QifMonthDayYearTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D09-01-2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D9-2-2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D9/3/2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D9.4.2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D9.5'2024\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 5, len(allNewTransactions))

	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(1725494400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
}

func TestQIFTransactionDataFileParseImportedData_ParseDayYearMonthDateFormatTime(t *testing.T) {
	importer := QifDayMonthYearTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D01-09-2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D2-9-2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D3/9/2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D4.9.2024\n"+
			"T-123.45\n"+
			"^\n"+
			"D5'9.2024\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 5, len(allNewTransactions))

	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(1725494400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
}

func TestQIFTransactionDataFileParseImportedData_ParseShortYearMonthDayDateFormatTime(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D24-09-01\n"+
			"T-123.45\n"+
			"^\n"+
			"D24-9-2\n"+
			"T-123.45\n"+
			"^\n"+
			"D24/9/3\n"+
			"T-123.45\n"+
			"^\n"+
			"D24.9.4\n"+
			"T-123.45\n"+
			"^\n"+
			"D24'9.5\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 5, len(allNewTransactions))

	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(1725494400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
}

func TestQIFTransactionDataFileParseImportedData_ParseInvalidTime(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024 09 01\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestQIFTransactionDataFileParseImportedData_ParseAmountWithThousandsSeparator(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123,456.78\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))

	assert.Equal(t, int64(12345678), allNewTransactions[0].Amount)
}

func TestQIFTransactionDataFileParseImportedData_ParseInvalidAmount(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123 45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestQIFTransactionDataFileParseImportedData_ParseAccountType(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Cash\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)

	allNewTransactions, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"!Type:CCard\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)

	allNewTransactions, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"!Type:Oth A\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)

	allNewTransactions, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"!Type:Oth L\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
}

func TestQIFTransactionDataFileParseImportedData_ParseAccount(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Account\n"+
			"NTest Account\n"+
			"^\n"+
			"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewAccounts))

	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)
}

func TestQIFTransactionDataFileParseImportedData_ParseAmountWithLeadingPlusSign(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T+123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
}

func TestQIFTransactionDataFileParseImportedData_ParseSubCategory(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, allNewSubExpenseCategories, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"LTest Category:Sub Category\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))

	assert.Equal(t, "Sub Category", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Sub Category", allNewSubExpenseCategories[0].Name)
}

func TestQIFTransactionDataFileParseImportedData_ParseDescription(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"PTest\n"+
			"Mfoo    bar\t#test\n"+
			"^\n"+
			"D2024-09-02\n"+
			"T-234.56\n"+
			"PTest2\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)
	assert.Equal(t, "", allNewTransactions[1].Comment)

	allNewTransactions, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"PTest\n"+
			"Mfoo    bar\t#test\n"+
			"^\n"+
			"D2024-09-02\n"+
			"T-234.56\n"+
			"PTest2\n"+
			"^\n"), 0, converter.DefaultImporterOptions.WithPayeeAsDescription(), nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)
	assert.Equal(t, "Test2", allNewTransactions[1].Comment)
}

func TestQIFTransactionDataFileParseImportedData_WithAdditionalOptions(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"PTest2\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 0, len(allNewTransactions[0].OriginalTagNames))

	allNewTransactions, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"T-123.45\n"+
			"PTest2\n"+
			"^\n"), 0, converter.DefaultImporterOptions.WithPayeeAsTag(), nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewTransactions[0].OriginalTagNames))
	assert.Equal(t, "Test2", allNewTransactions[0].OriginalTagNames[0])
}

func TestQIFTransactionDataFileParseImportedData_MissingRequiredFields(t *testing.T) {
	importer := QifYearMonthDayTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Time Field
	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"T-123.45\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingTransactionTime.Message)

	// Missing Amount Field
	_, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"!Type:Bank\n"+
			"D2024-09-01\n"+
			"^\n"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}
