package ofx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestOFXTransactionDataFileParseImportedData_MinimumValidData(t *testing.T) {
	converter := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte(
		"<OFX>"+
			"  <BANKMSGSRSV1>"+
			"  <STMTTRNRS>"+
			"    <STMTRS>"+
			"      <CURDEF>CNY</CURDEF>"+
			"      <BANKACCTFROM>"+
			"        <ACCTID>123</ACCTID>"+
			"      </BANKACCTFROM>"+
			"      <BANKTRANLIST>"+
			"        <STMTTRN>"+
			"          <TRNTYPE>DEP</TRNTYPE>"+
			"          <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>"+
			"          <TRNAMT>123.45</TRNAMT>"+
			"        </STMTTRN>"+
			"        <STMTTRN>"+
			"          <TRNTYPE>CHECK</TRNTYPE>"+
			"          <DTPOSTED>20240901123456.000[+8:CST]</DTPOSTED>"+
			"          <TRNAMT>-0.12</TRNAMT>"+
			"        </STMTTRN>"+
			"        <STMTTRN>"+
			"          <TRNTYPE>XFER</TRNTYPE>"+
			"          <DTPOSTED>20240901235959.000[+8:CST]</DTPOSTED>"+
			"          <TRNAMT>-1.00</TRNAMT>"+
			"        </STMTTRN>"+
			"      </BANKTRANLIST>"+
			"    </STMTRS>"+
			"  </STMTTRNRS>"+
			"  </BANKMSGSRSV1>"+
			"</OFX>"), 0, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 3, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, int64(1725125025), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "123", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "CNY", allNewTransactions[0].OriginalSourceAccountCurrency)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, int64(1725165296), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "123", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "CNY", allNewTransactions[1].OriginalSourceAccountCurrency)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[2].Type)
	assert.Equal(t, int64(1725206399), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "123", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "123", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}
