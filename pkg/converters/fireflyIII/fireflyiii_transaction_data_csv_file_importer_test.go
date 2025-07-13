package fireflyIII

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestFireFlyIIICsvFileConverterParseImportedData_MinimumValidData(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",123.45,2024-09-01T00:00:00+08:00,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\"Asset account\",\n"+
		"Deposit,0.12,2024-09-01T01:23:45+08:00,\"A revenue account\",\"Revenue account\",\"Test Account\",\"Asset account\",\"Test Category\"\n"+
		"Withdrawal,-1.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category2\"\n"+
		"Transfer,0.05,2024-09-01T23:59:59+08:00,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category3\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

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
	assert.Equal(t, "Test Category", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725165296), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725206399), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category3", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "Test Category3", allNewSubTransferCategories[0].Name)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-1.00,2024-09-01T12:34:56,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-1.00,2024-09-01 12:34:56+08:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseValidTransactionType(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// income transactions
	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Deposit,10.00,2024-09-01T12:34:56+08:00,\"A revenue account\",\"Revenue account\",\"Test Account\",\"Asset account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Deposit,10.00,2024-09-01T12:34:56+08:00,\"A revenue account\",\"Revenue account\",\"Test Account\",\"Debt\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)

	// expense transactions
	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-10.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-10.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Debt\",\"A expense account\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)

	// opening balance transactions
	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",10.00,2024-09-01T12:34:56+08:00,\"Initial balance\",\"Initial balance account\",\"Test Account\",\"Asset account\",\"\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)

	// transfer transactions
	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,10.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[0].OriginalDestinationAccountName)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-10.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"Test Account2\",\"Debt\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[0].OriginalDestinationAccountName)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Deposit,10.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Debt\",\"Test Account2\",\"Asset account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[0].OriginalDestinationAccountName)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,10.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Debt\",\"Test Account2\",\"Debt\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[0].OriginalDestinationAccountName)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseInvalidTransactionType(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,10.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Revenue account\",\"Test Account2\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTypeInvalid.Message)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseAccountNameAsCategoryName(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-1.00,2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "A expense account", allNewTransactions[0].OriginalCategoryName)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Deposit,10.00,2024-09-01T12:34:56+08:00,\"A revenue account\",\"Revenue account\",\"Test Account\",\"Asset account\",\"\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "A revenue account", allNewTransactions[0].OriginalCategoryName)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseValidTimezone(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-1.00,2024-09-01T12:34:56-10:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725230096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-1.00,2024-09-01T12:34:56+00:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-1.00,2024-09-01T12:34:56+12:45,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725148196), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseValidAccountCurrency(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,foreign_amount,date,currency_code,foreign_currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",123.45,,2024-09-01T00:00:00+08:00,USD,,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\"Asset account\",\n"+
		"Transfer,1.23,-1.10,2024-09-01T23:59:59+08:00,USD,EUR,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category2\""), 0, nil, nil, nil, nil, nil)

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

func TestFireFlyIIICsvFileConverterParseImportedData_ParseValidForeignAmountAndCurrency(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,foreign_amount,date,currency_code,foreign_currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,10.00,15.00,2024-09-01T12:34:56+08:00,USD,EUR,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, int64(1500), allNewTransactions[0].RelatedAccountAmount)
	assert.Equal(t, "USD", allNewTransactions[0].OriginalSourceAccountCurrency)
	assert.Equal(t, "EUR", allNewTransactions[0].OriginalDestinationAccountCurrency)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,currency_code,foreign_currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,10.00,2024-09-01T12:34:56+08:00,USD,EUR,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1000), allNewTransactions[0].Amount)
	assert.Equal(t, int64(1000), allNewTransactions[0].RelatedAccountAmount)
	assert.Equal(t, "USD", allNewTransactions[0].OriginalSourceAccountCurrency)
	assert.Equal(t, "EUR", allNewTransactions[0].OriginalDestinationAccountCurrency)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,10.00,2024-09-01T12:34:56+08:00,USD,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category\""), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "USD", allNewTransactions[0].OriginalSourceAccountCurrency)
	assert.Equal(t, "USD", allNewTransactions[0].OriginalDestinationAccountCurrency)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseInvalidAccountCurrency(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,foreign_amount,date,currency_code,foreign_currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",123.45,,2024-09-01T00:00:00+08:00,USD,,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\"Asset account\",\n"+
		"Transfer,1.23,1.10,2024-09-01T23:59:59+08:00,CNY,EUR,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category3\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,foreign_amount,date,currency_code,foreign_currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",123.45,,2024-09-01T00:00:00+08:00,USD,,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\n"+
		"Transfer,1.23,1.10,2024-09-01T23:59:59+08:00,CNY,EUR,\"Test Account2\",\"Asset account\",\"Test Account\",\"Asset account\",\"Test Category3\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseNotSupportedCurrency(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,foreign_amount,date,currency_code,foreign_currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",123.45,,2024-09-01T00:00:00+08:00,XXX,,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\"Asset account\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,foreign_amount,date,currency_code,foreign_currency_code,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,123.45,123.45,2024-09-01T23:59:59+08:00,USD,XXX,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category2\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseInvalidAmount(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-123 45,2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\"\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,foreign_amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Transfer,123.45,123 45,2024-09-01T23:59:59+08:00,\"Test Account\",\"Asset account\",\"Test Account2\",\"Asset account\",\"Test Category2\""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseDescription(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,description,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-123.45,\"foo    bar\t#test\",2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\"\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)
}

func TestFireFlyIIICsvFileConverterParseImportedData_ParseTags(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, allNewTags, err := converter.ParseImportedData(context, user, []byte("type,amount,tags,date,source_name,source_type,destination_name,destination_type,category\n"+
		"Withdrawal,-123.45,\"tag1,tag2,tag3\",2024-09-01T12:34:56+08:00,\"Test Account\",\"Asset account\",\"A expense account\",\"Expense account\",\"Test Category\"\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 3, len(allNewTags))
	assert.Equal(t, int64(1234567890), allNewTags[0].Uid)
	assert.Equal(t, "tag1", allNewTags[0].Name)
	assert.Equal(t, int64(1234567890), allNewTags[1].Uid)
	assert.Equal(t, "tag2", allNewTags[1].Name)
	assert.Equal(t, int64(1234567890), allNewTags[2].Uid)
	assert.Equal(t, "tag3", allNewTags[2].Name)
}

func TestFireFlyIIICsvFileConverterParseImportedData_MissingFileHeader(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}

func TestFireFlyIIICsvFileConverterParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := FireflyIIITransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Time Column
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte("type,amount,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",123.45,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\"Asset account\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Type Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("amount,date,source_name,source_type,destination_name,destination_type,category\n"+
		"123.45,2024-09-01T00:00:00+08:00,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\"Asset account\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Account Name Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,destination_name,category\n"+
		"\"Opening balance\",123.45,2024-09-01T00:00:00+08:00,\"Test Account\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,date,source_name,source_type,destination_name,destination_type,category\n"+
		"\"Opening balance\",2024-09-01T00:00:00+08:00,\"Initial balance for \"\"Test Account\"\"\",\"Initial balance account\",\"Test Account\",\"Asset account\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Account2 Name Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,category\n"+
		"\"Opening balance\",123.45,2024-09-01T00:00:00+08:00,\"Initial balance for \"\"Test Account\"\"\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Source Account Type Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,destination_name,destination_type,category\n"+
		"\"Opening balance\",123.45,2024-09-01T00:00:00+08:00,\"Initial balance for \"\"Test Account\"\"\",\"Test Account\",\"Asset account\",\"Asset account\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Destination Account Type Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte("type,amount,date,source_name,source_type,destination_name,category\n"+
		"\"Opening balance\",123.45,2024-09-01T00:00:00+08:00,\"Initial balance for \"\"Test Account\"\"\",\"Asset account\",\"Test Account\",\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}
