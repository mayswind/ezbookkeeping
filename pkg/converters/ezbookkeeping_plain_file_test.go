package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestEzBookKeepingPlainFileConverterParseImportedData_MinimumValidData(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubCategories, allNewTags, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
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

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725153825), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725235199), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)

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
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01T12:34:56,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"09/01/2024 12:34:56,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidType(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Type,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseValidTimezone(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Timezone,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,+08:00,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(1725165296), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidTimezone(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Timezone,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Asia/Shanghai,Expense,Test Category,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidSubCategoryName(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Expense,,Test Account,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidAccountName(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Expense,Test Category,,123.45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Transfer,Test Category,Test Account,123.45,,123.45"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseValidAccountCurrency(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
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
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Balance Modification,,Test Account,USD,123.45,,,\n"+
		"2024-09-01 12:34:56,Transfer,Test Category3,Test Account,CNY,1.23,Test Account2,EUR,1.10"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Balance Modification,,Test Account,USD,123.45,,,\n"+
		"2024-09-01 12:34:56,Transfer,Test Category3,Test Account2,CNY,1.23,Test Account,EUR,1.10"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseNotSupportedCurrency(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Balance Modification,Test Category,Test Account,XXX,123.45,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount\n"+
		"2024-09-01 01:23:45,Transfer,Test Category,Test Account,USD,123.45,Test Account2,XXX,123.45"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidAmount(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123 45,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount\n"+
		"2024-09-01 12:34:56,Transfer,Test Category,Test Account,123.45,Test Account2,123 45"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseValidGeographicLocation(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,123.45 45.56"), 0, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 123.45, allNewTransactions[0].GeoLongitude)
	assert.Equal(t, 45.56, allNewTransactions[0].GeoLatitude)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseInvalidGeographicLocation(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,1"), 0, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, float64(0), allNewTransactions[0].GeoLongitude)
	assert.Equal(t, float64(0), allNewTransactions[0].GeoLatitude)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,a b"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Geographic Location\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,1 "), 0, nil, nil, nil)
	assert.NotNil(t, err)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_ParseTag(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, allNewTags, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Tags\n"+
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
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, err := converter.parseImportedData(user, ",", []byte("Time,Type,Sub Category,Account,Amount,Account2,Account2 Amount,Description\n"+
		"2024-09-01 12:34:56,Expense,Test Category,Test Account,123.45,,,foo    bar\t#test"), 0, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)
}

func TestEzBookKeepingPlainFileConverterParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := &EzBookKeepingPlainFileConverter{}

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, err := converter.parseImportedData(user, ",", []byte(""), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Type,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Account,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Timezone,Type,Category,Sub Category,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,CNY,123.45,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)

	_, _, _, _, err = converter.parseImportedData(user, ",", []byte("Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description\n"+
		"2024-09-01 00:00:00,+08:00,Balance Modification,Test Category,Test Sub Category,Test Account,CNY,123.45,,,,,"), 0, nil, nil, nil)
	assert.NotNil(t, err)
}
