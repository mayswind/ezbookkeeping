package feidee

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestFeideeMymoneyElecloudTransactionDataXlsxImporterParseImportedData_MinimumValidData(t *testing.T) {
	importer := FeideeMymoneyElecloudTransactionDataXlsxFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "USD",
	}

	testdata, err := os.ReadFile("../../../testdata/feidee_mymoney_elecloud_test_file.xlsx")
	assert.Nil(t, err)

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, _, err := importer.ParseImportedData(context, user, testdata, 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 7, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 3, len(allNewSubExpenseCategories))
	assert.Equal(t, 3, len(allNewSubIncomeCategories))
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
	assert.Equal(t, "Test Account2", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[4].Type)
	assert.Equal(t, "2024-09-01 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime), time.UTC))
	assert.Equal(t, int64(5), allNewTransactions[4].Amount)
	assert.Equal(t, "Test Comment5", allNewTransactions[4].Comment)
	assert.Equal(t, "Test Account", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[4].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[5].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[5].Type)
	assert.Equal(t, "2024-09-10 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime), time.UTC))
	assert.Equal(t, int64(-654300), allNewTransactions[5].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[5].OriginalSourceAccountName)
	assert.Equal(t, "Test Category5", allNewTransactions[5].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[6].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[6].Type)
	assert.Equal(t, "2024-09-11 05:06:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[6].TransactionTime), time.UTC))
	assert.Equal(t, int64(-112340), allNewTransactions[6].Amount)
	assert.Equal(t, "Foo#\\r\\nBar", allNewTransactions[6].Comment)
	assert.Equal(t, "Test Account", allNewTransactions[6].OriginalSourceAccountName)
	assert.Equal(t, "Test Category4", allNewTransactions[6].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[1].Uid)
	assert.Equal(t, "Test Category4", allNewSubExpenseCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[2].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[2].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[1].Uid)
	assert.Equal(t, "Test Category5", allNewSubIncomeCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[2].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[2].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}
