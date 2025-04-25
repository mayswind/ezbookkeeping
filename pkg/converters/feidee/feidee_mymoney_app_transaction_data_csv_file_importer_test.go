package feidee

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestFeideeMymoneyCsvFileImporterParseImportedData_MinimumValidData(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 00:00:00\",\"\",\"Test Account\",\"123.45\",\"\",\"\"\n"+
		"\"余额变更\",\"2024-09-01 01:00:00\",\"\",\"Test Account2\",\"-0.12\",\"\",\"\"\n"+
		"\"收入\",\"2024-09-01 01:23:45\",\"Test Category\",\"Test Account\",\"0.12\",\"\",\"\"\n"+
		"\"支出\",\"2024-09-01 12:34:56\",\"Test Category2\",\"Test Account\",\"1.00\",\"\",\"\"\n"+
		"\"转出\",\"2024-09-01 23:59:59\",\"Test Category3\",\"Test Account\",\"0.05\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 23:59:59\",\"Test Category3\",\"Test Account2\",\"0.05\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-02 23:59:59\",\"Test Category3\",\"Test Account\",\"0.5\",\"\",\"00000000-0000-0000-0000-000000000002\"\n"+
		"\"转出\",\"2024-09-02 23:59:59\",\"Test Category3\",\"Test Account2\",\"0.5\",\"\",\"00000000-0000-0000-0000-000000000002\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 6, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 2, len(allNewSubExpenseCategories))
	assert.Equal(t, 2, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, "2024-09-01 01:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[2].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[3].Type)
	assert.Equal(t, "2024-09-01 12:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(100), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[4].Type)
	assert.Equal(t, "2024-09-01 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime), time.UTC))
	assert.Equal(t, int64(5), allNewTransactions[4].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[4].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category3", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[5].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[5].Type)
	assert.Equal(t, "2024-09-02 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime), time.UTC))
	assert.Equal(t, int64(50), allNewTransactions[5].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[5].OriginalSourceAccountName)
	assert.Equal(t, "Test Account", allNewTransactions[5].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category3", allNewTransactions[5].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[1].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[1].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "Test Category3", allNewSubTransferCategories[0].Name)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseOutstandingBalanceModification(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"负债变更\",\"2024-09-01 00:00:00\",\"\",\"Test Account\",\"123.45\",\"\",\"\"\n"+
		"\"负债变更\",\"2024-09-01 01:00:00\",\"\",\"Test Account2\",\"-0.12\",\"\",\"\"\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, "2024-09-01 01:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"收入\",\"2024-09-01T12:34:56\",\"Test Category\",\"Test Account\",\"0.12\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"收入\",\"09/01/2024 12:34:56\",\"Test Category\",\"Test Account\",\"0.12\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseInvalidType(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"Type\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account\",\"0.12\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTypeInvalid.Message)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseValidAccountCurrency(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"账户币种\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 01:23:45\",\"\",\"Test Account\",\"USD\",\"123.45\",\"\",\"\"\n"+
		"\"转出\",\"2024-09-01 12:34:56\",\"Test Category3\",\"Test Account\",\"USD\",\"1.23\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 12:34:56\",\"Test Category3\",\"Test Account2\",\"EUR\",\"1.10\",\"\",\"00000000-0000-0000-0000-000000000001\""), 0, nil, nil, nil, nil, nil)

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

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseInvalidAccountCurrency(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"账户币种\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 01:23:45\",\"\",\"Test Account\",\"USD\",\"123.45\",\"\",\"\"\n"+
		"\"转出\",\"2024-09-01 12:34:56\",\"Test Category3\",\"Test Account\",\"CNY\",\"1.23\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 12:34:56\",\"Test Category3\",\"Test Account2\",\"EUR\",\"1.10\",\"\",\"00000000-0000-0000-0000-000000000001\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"账户币种\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 01:23:45\",\"\",\"Test Account\",\"USD\",\"123.45\",\"\",\"\"\n"+
		"\"转出\",\"2024-09-01 12:34:56\",\"Test Category3\",\"Test Account2\",\"CNY\",\"1.23\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 12:34:56\",\"Test Category3\",\"Test Account\",\"EUR\",\"1.10\",\"\",\"00000000-0000-0000-0000-000000000001\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseNotSupportedCurrency(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"账户币种\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 01:23:45\",\"\",\"Test Account\",\"XXX\",\"123.45\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"账户币种\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"转出\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account\",\"USD\",\"123.45\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account2\",\"XXX\",\"123.45\",\"\",\"00000000-0000-0000-0000-000000000001\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"账户币种\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"转出\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account\",\"XXX\",\"123.45\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account2\",\"USD\",\"123.45\",\"\",\"00000000-0000-0000-0000-000000000001\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseInvalidAmount(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 01:23:45\",\"\",\"Test Account\",\"123 45\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"负债变更\",\"2024-09-01 01:23:45\",\"\",\"Test Account\",\"123 45\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"转出\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account\",\"123 45\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account2\",\"123.45\",\"\",\"00000000-0000-0000-0000-000000000001\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"转出\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account\",\"123.45\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account2\",\"123 45\",\"\",\"00000000-0000-0000-0000-000000000001\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_ParseDescription(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"支出\",\"2024-09-01 12:34:56\",\"Test Category\",\"Test Account\",\"123.45\",\"Test\n"+
		"A new line break\",\"\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Test\nA new line break", allNewTransactions[0].Comment)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_InvalidRelatedId(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"转出\",\"2024-09-01 23:59:59\",\"Test Category3\",\"Test Account\",\"0.05\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrRelatedIdCannotBeBlank.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"转入\",\"2024-09-01 23:59:59\",\"Test Category3\",\"Test Account\",\"0.05\",\"\",\"\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrRelatedIdCannotBeBlank.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"转出\",\"2024-09-01 23:59:59\",\"Test Category3\",\"Test Account\",\"0.05\",\"\",\"00000000-0000-0000-0000-000000000001\"\n"+
		"\"转入\",\"2024-09-02 23:59:59\",\"Test Category3\",\"Test Account\",\"0.5\",\"\",\"00000000-0000-0000-0000-000000000002\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrFoundRecordNotHasRelatedRecord.Message)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_MissingFileHeader(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 00:00:00\",\"\",\"Test Account\",\"123.45\",\"\",\"\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)
}

func TestFeideeMymoneyCsvFileImporterParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := FeideeMymoneyAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Time Column
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"\",\"Test Account\",\"123.45\",\"\",\"\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Type Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"2024-09-01 00:00:00\",\"Test Category\",\"Test Account\",\"123.45\",\"\",\"\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Sub Category Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"账户\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 00:00:00\",\"Test Account\",\"123.45\",\"\",\"\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Account Name Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"金额\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 00:00:00\",\"\",\"123.45\",\"\",\"\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"备注\",\"关联Id\"\n"+
		"\"余额变更\",\"2024-09-01 00:00:00\",\"\",\"Test Account\",\"\",\"\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Related ID Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("随手记导出文件(headers:v5;xxxxx)\n"+
		"\"交易类型\",\"日期\",\"子类别\",\"账户\",\"金额\",\"备注\"\n"+
		"\"余额变更\",\"2024-09-01 00:00:00\",\"\",\"Test Account\",\"123.45\",\"\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}
