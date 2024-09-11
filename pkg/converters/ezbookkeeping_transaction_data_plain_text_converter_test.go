package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestEzBookKeepingPlainFileConverterToExportedContent(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	transactions := make([]*models.Transaction, 3)
	transactions[0] = &models.Transaction{
		TransactionId:     1,
		TransactionTime:   1725165296000,
		Type:              models.TRANSACTION_DB_TYPE_INCOME,
		TimezoneUtcOffset: 480,
		CategoryId:        2,
		AccountId:         1,
		Amount:            12345,
		GeoLongitude:      123.45,
		GeoLatitude:       45.67,
		Comment:           "Hello,World",
	}
	transactions[1] = &models.Transaction{
		TransactionId:     2,
		TransactionTime:   1725194096000,
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TimezoneUtcOffset: 0,
		CategoryId:        4,
		AccountId:         1,
		Amount:            -10,
		GeoLongitude:      0,
		GeoLatitude:       0,
		Comment:           "Foo#Bar",
	}
	transactions[2] = &models.Transaction{
		TransactionId:        3,
		TransactionTime:      1725212096000,
		Type:                 models.TRANSACTION_DB_TYPE_TRANSFER_OUT,
		TimezoneUtcOffset:    -300,
		CategoryId:           6,
		AccountId:            1,
		Amount:               12345,
		RelatedAccountId:     2,
		RelatedAccountAmount: 1735,
		Comment:              "T\te\rs\nt\r\ntest",
	}

	accountMap := make(map[int64]*models.Account, 2)
	accountMap[1] = &models.Account{
		AccountId: 1,
		Name:      "Test Account",
		Currency:  "CNY",
	}
	accountMap[2] = &models.Account{
		AccountId: 2,
		Name:      "Test Account2",
		Currency:  "USD",
	}

	categoryMap := make(map[int64]*models.TransactionCategory, 6)
	categoryMap[1] = &models.TransactionCategory{
		CategoryId: 1,
		Type:       models.CATEGORY_TYPE_INCOME,
		Name:       "Test Category",
	}
	categoryMap[2] = &models.TransactionCategory{
		CategoryId:       2,
		Type:             models.CATEGORY_TYPE_INCOME,
		ParentCategoryId: 1,
		Name:             "Test Sub Category",
	}
	categoryMap[3] = &models.TransactionCategory{
		CategoryId: 3,
		Type:       models.CATEGORY_TYPE_EXPENSE,
		Name:       "Test Category2",
	}
	categoryMap[4] = &models.TransactionCategory{
		CategoryId:       4,
		Type:             models.CATEGORY_TYPE_EXPENSE,
		ParentCategoryId: 3,
		Name:             "Test Sub Category2",
	}
	categoryMap[5] = &models.TransactionCategory{
		CategoryId: 5,
		Type:       models.CATEGORY_TYPE_TRANSFER,
		Name:       "Test Category3",
	}
	categoryMap[6] = &models.TransactionCategory{
		CategoryId:       6,
		Type:             models.CATEGORY_TYPE_TRANSFER,
		ParentCategoryId: 5,
		Name:             "Test Sub Category3",
	}

	tagMap := make(map[int64]*models.TransactionTag, 2)
	tagMap[1] = &models.TransactionTag{
		TagId: 1,
		Name:  "Test,Tag",
	}
	tagMap[2] = &models.TransactionTag{
		TagId: 2,
		Name:  "Test;Tag2",
	}

	allTagIndexes := make(map[int64][]int64, 2)
	allTagIndexes[1] = []int64{1, 2}
	allTagIndexes[2] = []int64{3, 1, 4}
	allTagIndexes[3] = []int64{2, 3}

	expectedContent := "Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n" +
		"2024-09-01 12:34:56,+08:00,Income,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,123.450000 45.670000,Test Tag;Test Tag2,Hello World\n" +
		"2024-09-01 12:34:56,+00:00,Expense,Test Category2,Test Sub Category2,Test Account,CNY,-0.10,,,,,Test Tag,Foo#Bar\n" +
		"2024-09-01 12:34:56,-05:00,Transfer,Test Category3,Test Sub Category3,Test Account,CNY,123.45,Test Account2,USD,17.35,,Test Tag2,T\te s t test\n"
	actualContent, err := converter.ToExportedContent(context, 123, transactions, accountMap, categoryMap, tagMap, allTagIndexes)

	assert.Nil(t, err)
	assert.Equal(t, expectedContent, string(actualContent))
}

func TestEzBookKeepingPlainFileConverterParseImportedData_MinimumValidData(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 00:00:00,Balance Modification,,Test Account,123.45,,\n"+
		"2024-09-01 01:23:45,Income,Test Category,Test Account,0.12,,\n"+
		"2024-09-01 12:34:56,Expense,Test Category2,Test Account,1.00,,\n"+
		"2024-09-01 23:59:59,Transfer,Test Category3,Test Account,0.05,Test Account2,0.05"), 0, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 3, len(allNewSubCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725153825), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725235199), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
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

	assert.Equal(t, int64(1234567890), allNewSubCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubCategories[1].Uid)
	assert.Equal(t, "Test Category2", allNewSubCategories[1].Name)

	assert.Equal(t, int64(1234567890), allNewSubCategories[2].Uid)
	assert.Equal(t, "Test Category3", allNewSubCategories[2].Name)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01T12:34:56,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"09/01/2024 12:34:56,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidType(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Type,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseValidTimezone(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Timezone,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,+08:00,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725165296), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidTimezone(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Timezone,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Asia/Shanghai,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidAccountName(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Expense,Test Category,,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Transfer,Test Category,Test Account,123.45,,123.45"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseValidAccountCurrency(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Balance Modification,Test Category,Test Account,USD,123.45,,,\n"+
		"2024-09-01 12:34:56,Transfer,Test Category2,Test Account,USD,1.23,Test Account2,EUR,1.10"), 0, nil, nil, nil)

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

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidAccountCurrency(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Balance Modification,,Test Account,USD,123.45,,,\n"+
		"2024-09-01 12:34:56,Transfer,Test Category3,Test Account,CNY,1.23,Test Account2,EUR,1.10"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Balance Modification,,Test Account,USD,123.45,,,\n"+
		"2024-09-01 12:34:56,Transfer,Test Category3,Test Account2,CNY,1.23,Test Account,EUR,1.10"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseNotSupportedCurrency(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Balance Modification,Test Category,Test Account,XXX,123.45,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Transfer,Test Category,Test Account,USD,123.45,Test Account2,XXX,123.45"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidAmount(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123 45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Transfer,Test Category,Test Account,123.45,Test Account2,123 45"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseNoAmount2(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,"), 0, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, int64(0), allNewTransactions[0].RelatedAccountAmount)

	allNewTransactions, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2\n"+
		"2024-09-01 12:34:56,Transfer,Test Category,Test Account,123.45,Test Account2"), 0, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, int64(12345), allNewTransactions[0].RelatedAccountAmount)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseValidGeographicLocation(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,123.45 45.56"), 0, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 123.45, allNewTransactions[0].GeoLongitude)
	assert.Equal(t, 45.56, allNewTransactions[0].GeoLatitude)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidGeographicLocation(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,1"), 0, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, float64(0), allNewTransactions[0].GeoLongitude)
	assert.Equal(t, float64(0), allNewTransactions[0].GeoLatitude)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,a b"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,1 "), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseTag(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, allNewTags, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Tags\n"+
		"2024-09-01 00:00:00,Balance Modification,,Test Account,123.45,,,foo;;bar.;#test;hello\tworld;;"), 0, nil, nil, nil)

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

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseDescription(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.ParseImportedData(context, user, []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Description\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,foo    bar\t#test"), 0, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := EzBookKeepingTransactionDataCSVFileConverter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.ParseImportedData(context, user, []byte(""), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Type,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Account,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Timezone,Type,Category,Sub Category,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.ParseImportedData(context, user, []byte("Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}
