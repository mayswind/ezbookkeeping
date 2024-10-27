package iif

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestIIFTransactionDataFileParseImportedData_MinimumValidData(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte(
		"!ACCNT\tNAME\tACCNTTYPE\n"+
			"ACCNT\tTest Account\tBANK\n"+
			"ACCNT\tTest Account2\tBANK\n"+
			"ACCNT\tTest Category\tINC\n"+
			"ACCNT\tTest Category2\tEXP\n"+
			"!TRNS\tTRNSTYPE\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tTRNSTYPE\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\t\n"+
			"TRNS\tBEGINBALCHECK\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\tBEGINBALCHECK\t09/01/2024\txxx\t-123.45\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tDEPOSIT\t09/02/2024\tTest Account\t0.12\n"+
			"SPL\tDEPOSIT\t09/02/2024\tTest Category\t-0.12\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tCREDIT CARD\t09/03/2024\tTest Account\t-1.00\n"+
			"SPL\tCREDIT CARD\t09/03/2024\tTest Category2\t1.00\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tTRANSFER\t09/04/2024\tTest Account\t-0.05\n"+
			"SPL\tTRANSFER\t09/04/2024\tTest Account2\t0.05\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tGENERAL JOURNAL\t09/05/2024\tTest Account\t0.06\n"+
			"SPL\tGENERAL JOURNAL\t09/05/2024\tTest Account2\t-0.06\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tDEPOSIT\t09/06/2024\tTest Category\t-23.45\n"+
			"SPL\tDEPOSIT\t09/06/2024\tTest Account2\t23.45\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tCREDIT CARD\t09/07/2024\tTest Category2\t34.56\n"+
			"SPL\tCREDIT CARD\t09/07/2024\tTest Account2\t-34.56\n"+
			"ENDTRNS\t\t\t\t\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 7, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
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
	assert.Equal(t, "Test Account", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[4].Type)
	assert.Equal(t, int64(1725494400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
	assert.Equal(t, int64(6), allNewTransactions[4].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "Test Account", allNewTransactions[4].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[5].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[5].Type)
	assert.Equal(t, int64(1725580800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime))
	assert.Equal(t, int64(2345), allNewTransactions[5].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[5].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[5].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[5].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[6].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[6].Type)
	assert.Equal(t, int64(1725667200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[6].TransactionTime))
	assert.Equal(t, int64(3456), allNewTransactions[6].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[6].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[6].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[6].OriginalCategoryName)

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
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}

func TestIIFTransactionDataFileParseImportedData_MinimumValidDataWithoutAccountData(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Category\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Category", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)
}

func TestIIFTransactionDataFileParseImportedData_MultipleDataset(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!ACCNT\tNAME\tACCNTTYPE\n"+
			"ACCNT\tTest Account3\tBANK\n"+
			"ACCNT\tTest Account4\tBANK\n"+
			"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/05/2024\tTest Account\t-0.05\n"+
			"SPL\t09/05/2024\tTest Account2\t0.05\n"+
			"ENDTRNS\t\t\t\n"+
			"!TRNS\tTRNSID\tTRNSTYPE\tDATE\tACCNT\tNAME\tCLASS\tAMOUNT\tDOCNUM\tMEMO\tCLEAR\tTOPRINT\tADDR5\tDUEDATE\tTERMS\n"+
			"!SPL\tSPLID\tTRNSTYPE\tDATE\tACCNT\tNAME\tCLASS\tAMOUNT\tDOCNUM\tMEMO\tCLEAR\tQNTY\tREIMBEXP\tSERVICEDATE\tOTHER2\n"+
			"!ENDTRNS\t\t\t\t\t\t\t\t\t\t\t\t\t\t\n"+
			"TRNS\t\tTRANSFER\t09/04/2024\tTest Account3\tTest Category\tTest Class\t123.45\t\t\t\t\t\t\t\n"+
			"SPL\t\tTRANSFER\t09/04/2024\tTest Account4\t\t\t-123.45\t\t\t\t\t\t\t\n"+
			"ENDTRNS\t\t\t\t\t\t\t\t\t\t\t\t\t\t\n"+
			"!CLASS\tNAME\tHIDDEN\n"+
			"CLASS\tTest Class\tN\n"+
			"!TRNS\tTRNSTYPE\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tTRNSTYPE\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\t\n"+
			"TRNS\tBEGINBALCHECK\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\tBEGINBALCHECK\t09/01/2024\txxx\t-123.45\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tDEPOSIT\t09/02/2024\tTest Account\t0.12\n"+
			"SPL\tDEPOSIT\t09/02/2024\tTest Category\t-0.12\n"+
			"ENDTRNS\t\t\t\t\n"+
			"TRNS\tCREDIT CARD\t09/03/2024\tTest Account\t-1.00\n"+
			"SPL\tCREDIT CARD\t09/03/2024\tTest Category2\t1.00\n"+
			"ENDTRNS\t\t\t\t\n"+
			"!ACCNT\tTEST\tNAME\tACCNTTYPE\n"+
			"ACCNT\t\tTest Category\tINC\n"+
			"ACCNT\t\tTest Category2\tEXP\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 5, len(allNewTransactions))

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
	assert.Equal(t, "Test Account", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725408000), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[3].Amount)
	assert.Equal(t, "Test Account4", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "Test Account3", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[4].Type)
	assert.Equal(t, int64(1725494400), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
	assert.Equal(t, int64(5), allNewTransactions[4].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[4].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[4].OriginalCategoryName)
}

func TestIIFTransactionDataFileParseImportedData_ParseCategoryAndSubCategory(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, allNewSubExpenseCategories, allNewSubIncomeCategories, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!ACCNT\tNAME\tACCNTTYPE\n"+
			"ACCNT\tTest Parent Category:Test Category\tINC\n"+
			"ACCNT\tTest Parent Category2:Test Category2\tEXP\n"+
			"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Parent Category:Test Category\t-123.45\n"+
			"ENDTRNS\t\t\t\n"+
			"TRNS\t09/02/2024\tTest Account2\t-123.45\n"+
			"SPL\t09/02/2024\tTest Parent Category2:Test Category2\t123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category2", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[0].Name)
}

func TestIIFTransactionDataFileParseImportedData_ParseNameAsTransferCategory(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, allNewSubTransferCategories, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tNAME\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tNAME\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\tTest Category\t-123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t\t123.45\n"+
			"ENDTRNS\t\t\t\t\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewSubTransferCategories))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Test Account2", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, "Test Category", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubTransferCategories[0].Name)
}

func TestIIFTransactionDataFileParseImportedData_ParseShortMonthDayFormatTime(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t9/01/2024\tTest Account\t123.45\n"+
			"SPL\t9/01/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"+
			"TRNS\t09/2/2024\tTest Account\t123.45\n"+
			"SPL\t09/2/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"+
			"TRNS\t9/3/2024\tTest Account\t123.45\n"+
			"SPL\t9/3/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 3, len(allNewTransactions))
	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(1725235200), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(1725321600), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
}

func TestIIFTransactionDataFileParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t2024/09/01\tTest Account\t123.45\n"+
			"SPL\t2024/09/01\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t9/1/24\tTest Account\t123.45\n"+
			"SPL\t9/1/24\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t2024-09-01\tTest Account\t123.45\n"+
			"SPL\t2024-09-01\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t9/24\tTest Account\t123.45\n"+
			"SPL\t9/24\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestIIFTransactionDataFileParseImportedData_ParseInvalidAmount(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123 45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123 45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestIIFTransactionDataFileParseImportedData_ParseDescription(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\tMEMO\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\tMEMO\n"+
			"!ENDTRNS\t\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\t\"foo    bar\t#test\"\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\t\n"+
			"ENDTRNS\t\t\t\t\n"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)
}

func TestIIFTransactionDataFileParseImportedData_NotSupportedToParseSplitTransaction(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-100.00\n"+
			"SPL\t09/01/2024\tTest Account3\t-23.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotSupportedSplitTransactions.Message)
}

func TestIIFTransactionDataFileParseImportedData_InvalidDataLines(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Transaction Line
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Missing Split Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Missing Transaction End Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Missing Transaction End Line (following is another header)
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"!ACCNT\tNAME\tACCNTTYPE\n"+
			"ACCNT\tTest Account\tBANK\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Invalid Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"TEST\t\t\t\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Repeat Transaction Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Repeat Transaction End Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\t\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\t\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)
}

func TestIIFTransactionDataFileParseImportedData_InvalidHeaderLines(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing All Sample Lines
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Missing Transaction Sample Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Missing Split Sample Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Missing Transaction End Sample Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Missing Transaction End Sample Line (following is data line)
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!SPL\tDATE\tACCNT\tAMOUNT\n"+
			"TRNS\t09/01/2024\tTest Account\t123.45\n"+
			"SPL\t09/01/2024\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)

	// Invalid Sample Line
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\tAMOUNT\n"+
			"!TEST\tDATE\tACCNT\tAMOUNT\n"+
			"!ENDTRNS\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidIIFFile.Message)
}

func TestIIFTransactionDataFileParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := IifTransactionDataFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Date Column
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(
		"!TRNS\tACCNT\tAMOUNT\t\n"+
			"!SPL\tACCNT\tAMOUNT\t\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\tTest Account\t123.45\n"+
			"SPL\tTest Account2\t-123.45\n"+
			"ENDTRNS\t\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Account Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tAMOUNT\t\n"+
			"!SPL\tDATE\tAMOUNT\t\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\t123.45\n"+
			"SPL\t09/01/2024\t-123.45\n"+
			"ENDTRNS\t\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(
		"!TRNS\tDATE\tACCNT\t\n"+
			"!SPL\tDATE\tACCNT\t\n"+
			"!ENDTRNS\t\t\t\n"+
			"TRNS\t09/01/2024\tTest Account\n"+
			"SPL\t09/01/2024\tTest Account2\n"+
			"ENDTRNS\t\t\t\t\n"), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}
