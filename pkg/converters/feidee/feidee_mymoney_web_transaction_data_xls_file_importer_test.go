package feidee

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestGetFeideeMymoneyWebColumnNameMapping_LegacyFormat(t *testing.T) {
	importer := FeideeMymoneyWebTransactionDataXlsFileImporter
	headers := []string{"交易类型", "分类", "子分类", "账户1", "账户2", "金额", "日期", "成员", "项目", "商家", "备注"}

	mapping := importer.getFeideeMymoneyWebColumnNameMapping(headers)
	assert.Equal(t, feideeMymoneyWebDataLegacyColumnNameMapping, mapping)
}

func TestGetFeideeMymoneyWebColumnNameMapping_NewFormatExpense(t *testing.T) {
	importer := FeideeMymoneyWebTransactionDataXlsFileImporter
	headers := []string{"交易类型", "日期", "一级分类", "二级分类", "支出账户", "金额", "成员", "商家", "项目", "备注"}

	mapping := importer.getFeideeMymoneyWebColumnNameMapping(headers)
	assert.Equal(t, feideeMymoneyWebDataExpenseColumnNameMapping, mapping)
	assert.Equal(t, "支出账户", mapping[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME])
	assert.Equal(t, "一级分类", mapping[datatable.TRANSACTION_DATA_TABLE_CATEGORY])
}

func TestGetFeideeMymoneyWebColumnNameMapping_NewFormatIncome(t *testing.T) {
	importer := FeideeMymoneyWebTransactionDataXlsFileImporter
	headers := []string{"交易类型", "日期", "一级分类", "二级分类", "收入账户", "金额", "成员", "商家", "项目", "备注"}

	mapping := importer.getFeideeMymoneyWebColumnNameMapping(headers)
	assert.Equal(t, feideeMymoneyWebDataIncomeColumnNameMapping, mapping)
	assert.Equal(t, "收入账户", mapping[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME])
	assert.Equal(t, "一级分类", mapping[datatable.TRANSACTION_DATA_TABLE_CATEGORY])
}

func TestGetFeideeMymoneyWebColumnNameMapping_NewFormatBalanceModification(t *testing.T) {
	importer := FeideeMymoneyWebTransactionDataXlsFileImporter
	headers := []string{"交易类型", "日期", "一级分类", "二级分类", "账户1", "账户2", "金额", "成员", "商家", "项目", "备注"}

	mapping := importer.getFeideeMymoneyWebColumnNameMapping(headers)
	assert.Equal(t, feideeMymoneyWebDataBalanceModificationColumnNameMapping, mapping)
	assert.Equal(t, "账户1", mapping[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME])
	assert.Equal(t, "账户2", mapping[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME])
	assert.Equal(t, "一级分类", mapping[datatable.TRANSACTION_DATA_TABLE_CATEGORY])
}

func TestGetFeideeMymoneyWebColumnNameMapping_NewFormatTransfer(t *testing.T) {
	importer := FeideeMymoneyWebTransactionDataXlsFileImporter
	headers := []string{"交易类型", "日期", "转出账户", "转入账户", "金额", "成员", "商家", "项目", "备注"}

	mapping := importer.getFeideeMymoneyWebColumnNameMapping(headers)
	assert.Equal(t, feideeMymoneyWebDataTransferColumnNameMapping, mapping)
	assert.Equal(t, "转出账户", mapping[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME])
	assert.Equal(t, "转入账户", mapping[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME])
}

func TestFeideeMymoneyTransactionDataXlsImporterParseImportedData_MinimumValidData(t *testing.T) {
	importer := FeideeMymoneyWebTransactionDataXlsFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	testdata, err := os.ReadFile("../../../testdata/feidee_mymoney_test_file.xls")
	assert.Nil(t, err)

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := importer.ParseImportedData(context, user, testdata, time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 8, len(allNewTransactions))
	assert.Equal(t, 3, len(allNewAccounts))
	assert.Equal(t, 3, len(allNewSubExpenseCategories))
	assert.Equal(t, 3, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2026-07-07 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(-1230), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[1].Type)
	assert.Equal(t, "2026-07-07 00:01:02", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(1), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[1].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[2].Type)
	assert.Equal(t, "2026-07-07 00:29:10", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[3].Type)
	assert.Equal(t, "2026-07-07 00:30:05", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[4].Type)
	assert.Equal(t, "2026-07-07 00:30:15", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime), time.UTC))
	assert.Equal(t, int64(50), allNewTransactions[4].Amount)
	assert.Equal(t, "Test Account3", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[5].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[5].Type)
	assert.Equal(t, "2026-07-07 00:31:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime), time.UTC))
	assert.Equal(t, int64(98), allNewTransactions[5].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[5].OriginalSourceAccountName)
	assert.Equal(t, "Test Category3", allNewTransactions[5].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[6].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[6].Type)
	assert.Equal(t, "2026-07-07 00:31:33", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[6].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12300), allNewTransactions[6].Amount)
	assert.Equal(t, "Test Description", allNewTransactions[6].Comment)
	assert.Equal(t, "Test Account2", allNewTransactions[6].OriginalSourceAccountName)
	assert.Equal(t, "Test Category4", allNewTransactions[6].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[7].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[7].Type)
	assert.Equal(t, "2026-07-07 00:32:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[7].TransactionTime), time.UTC))
	assert.Equal(t, int64(100), allNewTransactions[7].Amount)
	assert.Equal(t, "Test Account3", allNewTransactions[7].OriginalSourceAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[7].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account3", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[2].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[2].Name)
	assert.Equal(t, "CNY", allNewAccounts[2].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[1].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[2].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[2].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Test Category4", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[1].Uid)
	assert.Equal(t, "Test Category3", allNewSubIncomeCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[2].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[2].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}

func TestFeideeMymoneyTransactionDataXlsImporterParseImportedData_MinimumLegacyValidData(t *testing.T) {
	importer := FeideeMymoneyWebTransactionDataXlsFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	testdata, err := os.ReadFile("../../../testdata/feidee_mymoney_legacy_test_file.xls")
	assert.Nil(t, err)

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := importer.ParseImportedData(context, user, testdata, time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
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
	assert.Equal(t, "Test Category3", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[5].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[5].Type)
	assert.Equal(t, "2024-09-10 00:00:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime), time.UTC))
	assert.Equal(t, int64(-54300), allNewTransactions[5].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[5].OriginalSourceAccountName)
	assert.Equal(t, "Test Category5", allNewTransactions[5].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[6].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[6].Type)
	assert.Equal(t, "2024-09-11 05:06:00", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[6].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12340), allNewTransactions[6].Amount)
	assert.Equal(t, "Line1\nLine2", allNewTransactions[6].Comment)
	assert.Equal(t, "Test Account", allNewTransactions[6].OriginalSourceAccountName)
	assert.Equal(t, "Test Category4", allNewTransactions[6].OriginalCategoryName)

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

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[2].Uid)
	assert.Equal(t, "Test Category4", allNewSubExpenseCategories[2].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[1].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[2].Uid)
	assert.Equal(t, "Test Category5", allNewSubIncomeCategories[2].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "Test Category3", allNewSubTransferCategories[0].Name)
}
