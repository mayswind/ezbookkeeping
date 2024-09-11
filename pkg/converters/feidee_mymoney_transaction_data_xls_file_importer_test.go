package converters

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestFeideeMymoneyTransactionDataXlsImporterParseImportedData_MinimumValidData(t *testing.T) {
	converter := FeideeMymoneyTransactionDataXlsImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	testdata, err := os.ReadFile("../../testdata/feidee_mymoney_test_file.xls")
	assert.Nil(t, err)

	allNewTransactions, allNewAccounts, allNewSubCategories, allNewTags, err := converter.ParseImportedData(context, user, testdata, 0, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 6, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 5, len(allNewSubCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, "2024-09-01 12:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, "2024-09-01 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Comment5", allNewTransactions[3].Comment)
	assert.Equal(t, "Test Account", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category3", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[4].Type)
	assert.Equal(t, "2024-09-10 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime), time.UTC))
	assert.Equal(t, int64(-54300), allNewTransactions[4].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "Test Category5", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[5].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[5].Type)
	assert.Equal(t, "2024-09-11 05:06:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12340), allNewTransactions[5].Amount)
	assert.Equal(t, "Line1\nLine2", allNewTransactions[5].Comment)
	assert.Equal(t, "Test Account", allNewTransactions[5].OriginalSourceAccountName)
	assert.Equal(t, "Test Category4", allNewTransactions[5].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubCategories[1].Uid)
	assert.Equal(t, "Test Category5", allNewSubCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubCategories[2].Uid)
	assert.Equal(t, "Test Category2", allNewSubCategories[2].Name)

	assert.Equal(t, int64(1234567890), allNewSubCategories[3].Uid)
	assert.Equal(t, "Test Category4", allNewSubCategories[3].Name)

	assert.Equal(t, int64(1234567890), allNewSubCategories[4].Uid)
	assert.Equal(t, "Test Category3", allNewSubCategories[4].Name)
}
