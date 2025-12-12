package ofx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestOFXTransactionDataFileParseImportedData_MinimumValidData(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>CHECK</TRNTYPE>\n"+
			"            <DTPOSTED>20240901123456.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>-0.12</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>XFER</TRNTYPE>\n"+
			"            <DTPOSTED>20240901225959.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>-1.00</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>XFER</TRNTYPE>\n"+
			"            <DTPOSTED>20240901235959.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>2.00</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"  <CREDITCARDMSGSRSV1>\n"+
			"    <CCSTMTTRNRS>\n"+
			"      <CCSTMTRS>\n"+
			"        <CURDEF>USD</CURDEF>\n"+
			"        <CCACCTFROM>\n"+
			"          <ACCTID>456</ACCTID>\n"+
			"        </CCACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>ATM</TRNTYPE>\n"+
			"            <DTPOSTED>20240902012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>1.23</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>POS</TRNTYPE>\n"+
			"            <DTPOSTED>20240902123456.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>-0.01</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </CCSTMTRS>\n"+
			"    </CCSTMTTRNRS>\n"+
			"  </CREDITCARDMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 6, len(allNewTransactions))
	assert.Equal(t, 3, len(allNewAccounts))
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
	assert.Equal(t, int64(1725202799), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(100), allNewTransactions[2].Amount)
	assert.Equal(t, "123", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, int64(1725206399), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(200), allNewTransactions[3].Amount)
	assert.Equal(t, "", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "123", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[4].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[4].Type)
	assert.Equal(t, int64(1725211425), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
	assert.Equal(t, int64(123), allNewTransactions[4].Amount)
	assert.Equal(t, "", allNewTransactions[4].OriginalSourceAccountName)
	assert.Equal(t, "USD", allNewTransactions[4].OriginalSourceAccountCurrency)
	assert.Equal(t, "456", allNewTransactions[4].OriginalDestinationAccountName)
	assert.Equal(t, "USD", allNewTransactions[4].OriginalDestinationAccountCurrency)
	assert.Equal(t, "", allNewTransactions[4].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[5].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[5].Type)
	assert.Equal(t, int64(1725251696), utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime))
	assert.Equal(t, int64(1), allNewTransactions[5].Amount)
	assert.Equal(t, "456", allNewTransactions[5].OriginalSourceAccountName)
	assert.Equal(t, "USD", allNewTransactions[5].OriginalSourceAccountCurrency)
	assert.Equal(t, "", allNewTransactions[5].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "123", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[2].Uid)
	assert.Equal(t, "456", allNewAccounts[2].Name)
	assert.Equal(t, "USD", allNewAccounts[2].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "", allNewSubTransferCategories[0].Name)
}

func TestOFXTransactionDataFileParseImportedData_ParseAccountTo(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, allNewAccounts, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>XFER</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>-123.45</TRNAMT>\n"+
			"            <BANKACCTTO>\n"+
			"              <ACCTID>456</ACCTID>\n"+
			"            </BANKACCTTO>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"  <CREDITCARDMSGSRSV1>\n"+
			"    <CCSTMTTRNRS>\n"+
			"      <CCSTMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <CCACCTFROM>\n"+
			"          <ACCTID>456</ACCTID>\n"+
			"        </CCACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>XFER</TRNTYPE>\n"+
			"            <DTPOSTED>20240902012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>-1.23</TRNAMT>\n"+
			"            <CCACCTTO>\n"+
			"              <ACCTID>789</ACCTID>\n"+
			"            </CCACCTTO>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </CCSTMTRS>\n"+
			"    </CCSTMTTRNRS>\n"+
			"  </CREDITCARDMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, 3, len(allNewAccounts))

	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
	assert.Equal(t, "123", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "CNY", allNewTransactions[0].OriginalSourceAccountCurrency)
	assert.Equal(t, "456", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, "CNY", allNewTransactions[0].OriginalDestinationAccountCurrency)

	assert.Equal(t, int64(123), allNewTransactions[1].Amount)
	assert.Equal(t, "456", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "CNY", allNewTransactions[1].OriginalSourceAccountCurrency)
	assert.Equal(t, "789", allNewTransactions[1].OriginalDestinationAccountName)
	assert.Equal(t, "CNY", allNewTransactions[1].OriginalDestinationAccountCurrency)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "123", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "456", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[2].Uid)
	assert.Equal(t, "789", allNewAccounts[2].Name)
	assert.Equal(t, "CNY", allNewAccounts[2].Currency)
}

func TestOFXTransactionDataFileParseImportedData_ParseValidTransactionTime(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901123456</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901123456.789</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901125959.000[-3]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901122959.000[-3.5]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240902030405.000[0]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 6, len(allNewTransactions))

	assert.Equal(t, int64(1725148800), utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime))
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime))
	assert.Equal(t, int64(1725194096), utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime))
	assert.Equal(t, int64(1725206399), utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime))
	assert.Equal(t, int64(1725206399), utils.GetUnixTimeFromTransactionTime(allNewTransactions[4].TransactionTime))
	assert.Equal(t, int64(1725246245), utils.GetUnixTimeFromTransactionTime(allNewTransactions[5].TransactionTime))
}

func TestOFXTransactionDataFileParseImportedData_ParseInvalidTransactionTime(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>2024</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>2024-09-01</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>202491</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	_, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901 12:34:56</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestOFXTransactionDataFileParseImportedData_ParseAmount_CommaAsDecimalPoint(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123,45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, int64(12345), allNewTransactions[0].Amount)
}

func TestOFXTransactionDataFileParseImportedData_ParseInvalidAmount(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123 45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestOFXTransactionDataFileParseImportedData_ParseTransactionCurrency(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"        <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"            <CURRENCY>USD</CURRENCY>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "USD", allNewTransactions[0].OriginalSourceAccountCurrency)
}

func TestOFXTransactionDataFileParseImportedData_ParseDescription(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	allNewTransactions, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"            <NAME>Test</NAME>\n"+
			"            <MEMO>foo    bar\t#test</MEMO>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo    bar\t#test", allNewTransactions[0].Comment)

	allNewTransactions, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"            <NAME>Test</NAME>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Test", allNewTransactions[0].Comment)

	allNewTransactions, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"            <PAYEE>\n"+
			"              <NAME>Test</NAME>\n"+
			"            </PAYEE>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Test", allNewTransactions[0].Comment)
}

func TestOFXTransactionDataFileParseImportedData_MissingAccountFromNode(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Posted Date Node
	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"  <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingAccountData.Message)
}

func TestOFXTransactionDataFileParseImportedData_MissingCurrencyNode(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Default Currency Node
	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"    <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAccountCurrencyInvalid.Message)
}

func TestOFXTransactionDataFileParseImportedData_MissingTransactionRequiredNode(t *testing.T) {
	importer := OFXTransactionDataImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Posted Date Node
	_, _, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"    <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingTransactionTime.Message)

	// Missing Transaction Type Node
	_, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"    <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"            <TRNAMT>123.45</TRNAMT>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTypeInvalid.Message)

	// Missing Amount Node
	_, _, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(
		"<OFX>\n"+
			"    <BANKMSGSRSV1>\n"+
			"    <STMTTRNRS>\n"+
			"      <STMTRS>\n"+
			"        <CURDEF>CNY</CURDEF>\n"+
			"        <BANKACCTFROM>\n"+
			"          <ACCTID>123</ACCTID>\n"+
			"        </BANKACCTFROM>\n"+
			"        <BANKTRANLIST>\n"+
			"          <STMTTRN>\n"+
			"            <TRNTYPE>DEP</TRNTYPE>\n"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"          </STMTTRN>\n"+
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"), 0, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}
