package beancount

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestBeancountTransactionDataFileParseImportedData_MinimumValidData(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte(
		"2024-09-01 *\n"+
			"  Equity:Opening-Balances -123.45 CNY\n"+
			"  Assets:TestAccount 123.45 CNY\n"+
			"2024-09-02 *\n"+
			"  Income:TestCategory -0.12 CNY\n"+
			"  Assets:TestAccount 0.12 CNY\n"+
			"2024-09-03 *\n"+
			"  Assets:TestAccount -1.00 CNY\n"+
			"  Expenses:TestCategory2 1.00 CNY\n"+
			"2024-09-04 *\n"+
			"  Assets:TestAccount -0.05 CNY\n"+
			"  Assets:TestAccount2 0.05 CNY\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Income:TestCategory", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Expenses:TestCategory2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Assets:TestAccount2", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Assets:TestAccount", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Assets:TestAccount2", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Expenses:TestCategory2", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Income:TestCategory", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}

func TestBeancountTransactionDataFileParseImportedData_MinimumValidData2(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte(
		"2024-09-01 *\n"+
			"  Assets:TestAccount 123.45 CNY\n"+
			"  Equity:Opening-Balances -123.45 CNY\n"+
			"2024-09-02 *\n"+
			"  Assets:TestAccount 0.12 CNY\n"+
			"  Income:TestCategory -0.12 CNY\n"+
			"2024-09-03 *\n"+
			"  Expenses:TestCategory2 1.00 CNY\n"+
			"  Assets:TestAccount -1.00 CNY\n"+
			"2024-09-04 *\n"+
			"  Assets:TestAccount2 0.05 CNY\n"+
			"  Assets:TestAccount -0.05 CNY\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Income:TestCategory", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Expenses:TestCategory2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Assets:TestAccount2", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Assets:TestAccount", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Assets:TestAccount2", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Expenses:TestCategory2", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Income:TestCategory", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}

func TestBeancountTransactionDataFileParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"2024/09/01 *\n"+
			"  Equity:Opening-Balances -123.45 CNY\n"+
			"  Assets:TestAccount 123.45 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestBeancountTransactionDataFileParseImportedData_ParseValidCurrency(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"2024-09-01 * \"Payee Name\" \"Hello\nWorld\"\n"+
			"  Assets:TestAccount -0.12 USD\n"+
			"  Assets:TestAccount2 0.84 CNY\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, int64(12), allNewTransactions[0].Amount)
	assert.Equal(t, int64(84), allNewTransactions[0].RelatedAccountAmount)
	assert.Equal(t, "Assets:TestAccount", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "USD", allNewTransactions[0].OriginalSourceAccountCurrency)
	assert.Equal(t, "Assets:TestAccount2", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, "CNY", allNewTransactions[0].OriginalDestinationAccountCurrency)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Assets:TestAccount", allNewAccounts[0].Name)
	assert.Equal(t, "USD", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Assets:TestAccount2", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)
}

func TestBeancountTransactionDataFileParseImportedData_ParseInvalidAmount(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"2024-09-01 *\n"+
			"  Equity:Opening-Balances -abc CNY\n"+
			"  Assets:TestAccount abc CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"2024-09-01 *\n"+
			"  Equity:Opening-Balances -1/0 CNY\n"+
			"  Assets:TestAccount 1/0 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestBeancountTransactionDataFileParseImportedData_ParseDescription(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"2024-09-01 * \"foo    bar\t#test\n\"\n"+
			"  Equity:Opening-Balances -123.45 CNY\n"+
			"  Assets:TestAccount 123.45 CNY\n"+
			"2024-09-02 * \"Payee Name\" \"Hello\nWorld\"\n"+
			"  Income:TestCategory -0.12 CNY\n"+
			"  Assets:TestAccount 0.12 CNY\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))

	assert.Equal(t, "foo    bar\t#test\n", allNewTransactions[0].Comment)
	assert.Equal(t, "Hello\nWorld", allNewTransactions[1].Comment)
}

func TestBeancountTransactionDataFileParseImportedData_InvalidTransaction(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"2024-09-02 * \"Payee Name\" \"Hello\nWorld\"\n"+
			"  Assets:TestAccount 0.11 CNY\n"+
			"  Assets:TestAccount2 0.11 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"2024-09-02 * \"Payee Name\" \"Hello\nWorld\"\n"+
			"  Expenses:TestCategory -0.11 CNY\n"+
			"  Expenses:TestCategory2 0.11 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrThereAreNotSupportedTransactionType.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"2024-09-02 * \"Payee Name\" \"Hello\nWorld\"\n"+
			"  Income:TestCategory -0.11 CNY\n"+
			"  Income:TestCategory2 0.11 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrThereAreNotSupportedTransactionType.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"2024-09-02 * \"Payee Name\" \"Hello\nWorld\"\n"+
			"  Equity:TestCategory -0.11 CNY\n"+
			"  Equity:TestCategory2 0.11 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrThereAreNotSupportedTransactionType.Message)
}

func TestBeancountTransactionDataFileParseImportedData_NotSupportedToParseSplitTransaction(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"2024-09-02 * \"Payee Name\" \"Hello\nWorld\"\n"+
			"  Assets:TestAccount -0.23 CNY\n"+
			"  Assets:TestAccount2 0.11 CNY\n"+
			"  Assets:TestAccount3 0.12 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotSupportedSplitTransactions.Message)
}

func TestBeancountTransactionDataFileParseImportedData_MissingTransactionRequiredData(t *testing.T) {
	converter := BeancountTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// Missing Transaction Time
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"* \"narration\"\n"+
			"  Equity:Opening-Balances -123.45 CNY\n"+
			"  Assets:TestAccount 123.45 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)

	// Missing Account Name
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"2024-09-01 * \"narration\"\n"+
			"  Equity:Opening-Balances -123.45 CNY\n"+
			"   123.45 CNY\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)

	// Missing Amount
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"2024-09-01 * \"narration\"\n"+
			"  Equity:Opening-Balances\n"+
			"  Assets:TestAccount\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)

	// Missing Commodity
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"2024-09-01 * \"narration\"\n"+
			"  Equity:Opening-Balances -123.45\n"+
			"  Assets:TestAccount 123.45\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)
}
