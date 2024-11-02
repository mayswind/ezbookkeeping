package ofx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestCreateNewOFXFileReader_OFX1(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<BANKMSGSRSV1>\n"+
			"<STMTTRNRS>\n"+
			"<STMTRS>\n"+
			"<CURDEF>CNY\n"+
			"<BANKACCTFROM>\n"+
			"<ACCTID>123\n"+
			"</BANKACCTFROM>\n"+
			"<BANKTRANLIST>\n"+
			"<STMTTRN>\n"+
			"<TRNTYPE>DEP\n"+
			"<DTPOSTED>20240901012345.000[+8:CST]\n"+
			"<TRNAMT>123.45\n"+
			"</STMTTRN>\n"+
			"</BANKTRANLIST>\n"+
			"</STMTRS>\n"+
			"</STMTTRNRS>\n"+
			"</BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion1, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "103", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)

	assert.Equal(t, "CNY", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.DefaultCurrency)

	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)
	assert.Equal(t, "123", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom.AccountId)

	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.Equal(t, ofxDepositTransaction, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].PostedDate)
	assert.Equal(t, "123.45", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Amount)
}

func TestCreateNewOFXFileReader_OFX1WithoutBreakLine(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>"+
			"<BANKMSGSRSV1>"+
			"<STMTTRNRS>"+
			"<STMTRS>"+
			"<CURDEF>CNY"+
			"<BANKACCTFROM>"+
			"<ACCTID>123"+
			"</BANKACCTFROM>"+
			"<BANKTRANLIST>"+
			"<STMTTRN>"+
			"<TRNTYPE>DEP"+
			"<DTPOSTED>20240901012345.000[+8:CST]"+
			"<TRNAMT>123.45"+
			"</STMTTRN>"+
			"</BANKTRANLIST>"+
			"</STMTRS>"+
			"</STMTTRNRS>"+
			"</BANKMSGSRSV1>"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion1, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "103", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)

	assert.Equal(t, "CNY", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.DefaultCurrency)

	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)
	assert.Equal(t, "123", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom.AccountId)

	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.Equal(t, ofxDepositTransaction, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].PostedDate)
	assert.Equal(t, "123.45", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Amount)
}

func TestCreateNewOFXFileReader_OFX1ParseBankAccountFrom(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<BANKMSGSRSV1>\n"+
			"<STMTTRNRS>\n"+
			"<STMTRS>\n"+
			"<BANKACCTFROM>\n"+
			"<BANKID>1234567890\n"+
			"<BRANCHID>2345678901\n"+
			"<ACCTID>3456789012\n"+
			"<ACCTTYPE>CHECKING\n"+
			"<ACCTKEY>4567890123\n"+
			"</BANKACCTFROM>\n"+
			"</STMTRS>\n"+
			"</STMTTRNRS>\n"+
			"</BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)

	account := ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom
	assert.Equal(t, "1234567890", account.BankId)
	assert.Equal(t, "2345678901", account.BranchId)
	assert.Equal(t, "3456789012", account.AccountId)
	assert.Equal(t, ofxCheckingAccount, account.AccountType)
	assert.Equal(t, "4567890123", account.AccountKey)
}

func TestCreateNewOFXFileReader_OFX1ParseCreditCardAccountFrom(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<CREDITCARDMSGSRSV1>\n"+
			"<CCSTMTTRNRS>\n"+
			"<CCSTMTRS>\n"+
			"<CCACCTFROM>\n"+
			"<ACCTID>3456789012\n"+
			"<ACCTKEY>4567890123\n"+
			"</CCACCTFROM>\n"+
			"</CCSTMTRS>\n"+
			"</CCSTMTTRNRS>\n"+
			"</CREDITCARDMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1)
	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse.StatementResponse)
	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)

	account := ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom
	assert.Equal(t, "3456789012", account.AccountId)
	assert.Equal(t, "4567890123", account.AccountKey)
}

func TestCreateNewOFXFileReader_OFX1ParseBankTransactionList(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<BANKMSGSRSV1>\n"+
			"<STMTTRNRS>\n"+
			"<STMTRS>\n"+
			"<BANKTRANLIST>\n"+
			"<DTSTART>20240901012345.000[+8:CST]\n"+
			"<DTEND>20240901235959.000[+8:CST]\n"+
			"</BANKTRANLIST>\n"+
			"</STMTRS>\n"+
			"</STMTTRNRS>\n"+
			"</BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList)

	transactionList := ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList
	assert.Equal(t, "20240901012345.000[+8:CST]", transactionList.StartDate)
	assert.Equal(t, "20240901235959.000[+8:CST]", transactionList.EndDate)
}

func TestCreateNewOFXFileReader_OFX1ParseCreditTransactionList(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<CREDITCARDMSGSRSV1>\n"+
			"<CCSTMTTRNRS>\n"+
			"<CCSTMTRS>\n"+
			"<BANKTRANLIST>\n"+
			"<DTSTART>20240901012345.000[+8:CST]\n"+
			"<DTEND>20240901235959.000[+8:CST]\n"+
			"</BANKTRANLIST>\n"+
			"</CCSTMTRS>\n"+
			"</CCSTMTTRNRS>\n"+
			"</CREDITCARDMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1)
	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse.StatementResponse)
	assert.NotNil(t, ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList)

	transactionList := ofxFile.CreditCardMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList
	assert.Equal(t, "20240901012345.000[+8:CST]", transactionList.StartDate)
	assert.Equal(t, "20240901235959.000[+8:CST]", transactionList.EndDate)
}

func TestCreateNewOFXFileReader_OFX1ParseTransaction(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<BANKMSGSRSV1>\n"+
			"<STMTTRNRS>\n"+
			"<STMTRS>\n"+
			"<BANKACCTFROM>\n"+
			"<ACCTID>123\n"+
			"</BANKACCTFROM>\n"+
			"<BANKTRANLIST>\n"+
			"<STMTTRN>\n"+
			"<FITID>1234567890\n"+
			"<TRNTYPE>CASH\n"+
			"<DTPOSTED>20240901012345.000[+8:CST]\n"+
			"<TRNAMT>123.45\n"+
			"<NAME>Test Name\n"+
			"<MEMO>Some Text\n"+
			"<CURRENCY>CNY\n"+
			"<ORIGCURRENCY>USD\n"+
			"</STMTTRN>\n"+
			"</BANKTRANLIST>\n"+
			"</STMTRS>\n"+
			"</STMTTRNRS>\n"+
			"</BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList)
	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0])

	transaction := ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0]
	assert.Equal(t, "1234567890", transaction.TransactionId)
	assert.Equal(t, ofxCashWithdrawalTransaction, transaction.TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", transaction.PostedDate)
	assert.Equal(t, "123.45", transaction.Amount)
	assert.Equal(t, "Test Name", transaction.Name)
	assert.Equal(t, "Some Text", transaction.Memo)
	assert.Equal(t, "CNY", transaction.Currency)
	assert.Equal(t, "USD", transaction.OriginalCurrency)
}

func TestCreateNewOFXFileReader_OFX1ParseTransactionPayee(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<BANKMSGSRSV1>\n"+
			"<STMTTRNRS>\n"+
			"<STMTRS>\n"+
			"<BANKACCTFROM>\n"+
			"<ACCTID>123\n"+
			"</BANKACCTFROM>\n"+
			"<BANKTRANLIST>\n"+
			"<STMTTRN>\n"+
			"<TRNTYPE>DEP\n"+
			"<DTPOSTED>20240901012345.000[+8:CST]\n"+
			"<TRNAMT>123.45\n"+
			"<PAYEE>\n"+
			"<NAME>Test Name\n"+
			"<ADDR1>Address 1\n"+
			"<ADDR2>Address 2\n"+
			"<ADDR3>Address 3\n"+
			"<CITY>City Name\n"+
			"<STATE>State Name\n"+
			"<POSTALCODE>10000000\n"+
			"<COUNTRY>Country Name\n"+
			"<PHONE>11111111111\n"+
			"</PAYEE>\n"+
			"</STMTTRN>\n"+
			"</BANKTRANLIST>\n"+
			"</STMTRS>\n"+
			"</STMTTRNRS>\n"+
			"</BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList)
	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0])
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Payee)

	payee := ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Payee
	assert.Equal(t, "Test Name", payee.Name)
	assert.Equal(t, "Address 1", payee.Address1)
	assert.Equal(t, "Address 2", payee.Address2)
	assert.Equal(t, "Address 3", payee.Address3)
	assert.Equal(t, "City Name", payee.City)
	assert.Equal(t, "State Name", payee.State)
	assert.Equal(t, "10000000", payee.PostalCode)
	assert.Equal(t, "Country Name", payee.Country)
	assert.Equal(t, "11111111111", payee.Phone)
}

func TestCreateNewOFXFileReader_OFX1WithEndElement(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"<BANKMSGSRSV1>\n"+
			"<STMTTRNRS>\n"+
			"<STMTRS>\n"+
			"<CURDEF>CNY</CURDEF>\n"+
			"<BANKACCTFROM>\n"+
			"<ACCTID>123</ACCTID>\n"+
			"</BANKACCTFROM>\n"+
			"<BANKTRANLIST>\n"+
			"<STMTTRN>\n"+
			"<TRNTYPE>DEP</TRNTYPE>\n"+
			"<DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>\n"+
			"<TRNAMT>123.45</TRNAMT>\n"+
			"</STMTTRN>\n"+
			"</BANKTRANLIST>\n"+
			"</STMTRS>\n"+
			"</STMTTRNRS>\n"+
			"</BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion1, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "103", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)

	assert.Equal(t, "CNY", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.DefaultCurrency)

	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)
	assert.Equal(t, "123", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom.AccountId)

	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.Equal(t, ofxDepositTransaction, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].PostedDate)
	assert.Equal(t, "123.45", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Amount)
}

func TestCreateNewOFXFileReader_OFX1WithBlanklinesInHeader(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"\n"+
			"\n"+
			"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>"+
			"<BANKMSGSRSV1>"+
			"<STMTTRNRS>"+
			"<STMTRS>"+
			"<CURDEF>CNY"+
			"<BANKACCTFROM>"+
			"<ACCTID>123"+
			"</BANKACCTFROM>"+
			"<BANKTRANLIST>"+
			"<STMTTRN>"+
			"<TRNTYPE>DEP"+
			"<DTPOSTED>20240901012345.000[+8:CST]"+
			"<TRNAMT>123.45"+
			"</STMTTRN>"+
			"</BANKTRANLIST>"+
			"</STMTRS>"+
			"</STMTTRNRS>"+
			"</BANKMSGSRSV1>"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion1, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "103", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)

	assert.Equal(t, "CNY", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.DefaultCurrency)

	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)
	assert.Equal(t, "123", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom.AccountId)

	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.Equal(t, ofxDepositTransaction, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].PostedDate)
	assert.Equal(t, "123.45", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Amount)
}

func TestCreateNewOFXFileReader_OFX1WithoutCharset(t *testing.T) {
	context := core.NewNullContext()

	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"FOO:BAR\n"+
			"\n"+
			"<OFX>\n"+
			"</OFX>"))
	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion1, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "103", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)
}

func TestCreateNewOFXFileReader_OFX1WithInvalidHeaderVersion(t *testing.T) {
	context := core.NewNullContext()
	_, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:200\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"</OFX>"))

	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)
}

func TestCreateNewOFXFileReader_OFX1WithInvalidHeader(t *testing.T) {
	context := core.NewNullContext()
	_, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:XML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"\n"+
			"<OFX>\n"+
			"</OFX>"))

	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)
}

func TestCreateNewOFXFileReader_OFX1WithUnknownHeader(t *testing.T) {
	context := core.NewNullContext()

	reader, err := createNewOFXFileReader(context, []byte(
		"OFXHEADER:100\n"+
			"DATA:OFXSGML\n"+
			"VERSION:103\n"+
			"SECURITY:NONE\n"+
			"ENCODING:USASCII\n"+
			"CHARSET:1252\n"+
			"COMPRESSION:NONE\n"+
			"OLDFILEUID:NONE\n"+
			"NEWFILEUID:NONE\n"+
			"FOO:BAR\n"+
			"\n"+
			"<OFX>\n"+
			"</OFX>"))
	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion1, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "103", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)
}

func TestCreateNewOFXFileReader_OFX2(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>\n"+
			"<?OFX OFXHEADER=\"200\" VERSION=\"211\" SECURITY=\"NONE\" OLDFILEUID=\"NONE\" NEWFILEUID=\"NONE\"?>\n"+
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
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion2, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "211", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)

	assert.Equal(t, "CNY", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.DefaultCurrency)

	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)
	assert.Equal(t, "123", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom.AccountId)

	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.Equal(t, ofxDepositTransaction, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].PostedDate)
	assert.Equal(t, "123.45", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Amount)
}

func TestCreateNewOFXFileReader_OFX2WithoutBreakLine(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>"+
			"<?OFX OFXHEADER=\"200\" VERSION=\"211\" SECURITY=\"NONE\" OLDFILEUID=\"NONE\" NEWFILEUID=\"NONE\"?>"+
			"<OFX>"+
			"  <BANKMSGSRSV1>"+
			"    <STMTTRNRS>"+
			"      <STMTRS>"+
			"        <CURDEF>CNY</CURDEF>"+
			"        <BANKACCTFROM>"+
			"          <ACCTID>123</ACCTID>"+
			"        </BANKACCTFROM>"+
			"        <BANKTRANLIST>"+
			"          <STMTTRN>"+
			"            <TRNTYPE>DEP</TRNTYPE>"+
			"            <DTPOSTED>20240901012345.000[+8:CST]</DTPOSTED>"+
			"            <TRNAMT>123.45</TRNAMT>"+
			"          </STMTTRN>"+
			"        </BANKTRANLIST>"+
			"      </STMTRS>"+
			"    </STMTTRNRS>"+
			"  </BANKMSGSRSV1>"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion2, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "211", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)

	assert.Equal(t, "CNY", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.DefaultCurrency)

	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)
	assert.Equal(t, "123", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom.AccountId)

	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.Equal(t, ofxDepositTransaction, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].PostedDate)
	assert.Equal(t, "123.45", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Amount)
}

func TestCreateNewOFXFileReader_OFX2WithoutOFXHeader(t *testing.T) {
	context := core.NewNullContext()
	_, err := createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>"+
			"<OFX>"+
			"</OFX>"))

	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)
}

func TestCreateNewOFXFileReader_OFX2WithInvalidHeaderVersion(t *testing.T) {
	context := core.NewNullContext()
	_, err := createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>"+
			"<?OFX OFXHEADER=\"100\" VERSION=\"211\" SECURITY=\"NONE\" OLDFILEUID=\"NONE\" NEWFILEUID=\"NONE\"?>"+
			"<OFX>"+
			"</OFX>"))

	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)
}

func TestCreateNewOFXFileReader_OFX2WithInvalidHeader(t *testing.T) {
	context := core.NewNullContext()
	_, err := createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>"+
			"<?OFX?>"+
			"<OFX>"+
			"</OFX>"))
	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)

	_, err = createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>"+
			"<?OFX OFXHEADER=200?>"+
			"<OFX>"+
			"</OFX>"))
	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)

	_, err = createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>"+
			"<?OFX OFXHEADER=\"200\" VERSION=\"211\" SECURITY=\"NONE\" OLDFILEUID=\"NONE\" NEWFILEUID=\"NONE\" test=\"\"?>"+
			"<OFX>"+
			"</OFX>"))
	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)
}

func TestCreateNewOFXFileReader_OFX2WithUnknownHeader(t *testing.T) {
	context := core.NewNullContext()

	reader, err := createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>"+
			"<?OFX OFXHEADER=\"200\" VERSION=\"211\" SECURITY=\"NONE\" OLDFILEUID=\"NONE\" NEWFILEUID=\"NONE\" FOO=\"BAR\"?>"+
			"<OFX>"+
			"</OFX>"))
	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.NotNil(t, ofxFile)

	assert.NotNil(t, ofxFile.FileHeader)
	assert.Equal(t, ofxVersion2, ofxFile.FileHeader.OFXDeclarationVersion)
	assert.Equal(t, "211", ofxFile.FileHeader.OFXDataVersion)
	assert.Equal(t, "NONE", ofxFile.FileHeader.Security)
	assert.Equal(t, "NONE", ofxFile.FileHeader.OldFileUid)
	assert.Equal(t, "NONE", ofxFile.FileHeader.NewFileUid)
}

func TestCreateNewOFXFileReader_OFX2WithSGML(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
		"<?xml version=\"1.0\" encoding=\"US-ASCII\"?>\n"+
			"<?OFX OFXHEADER=\"200\" VERSION=\"211\" SECURITY=\"NONE\" OLDFILEUID=\"NONE\" NEWFILEUID=\"NONE\"?>\n"+
			"<BANKMSGSRSV1>\n"+
			"<STMTTRNRS>\n"+
			"<STMTRS>\n"+
			"<CURDEF>CNY\n"+
			"<BANKACCTFROM>\n"+
			"<ACCTID>123\n"+
			"</BANKACCTFROM>\n"+
			"<BANKTRANLIST>\n"+
			"<STMTTRN>\n"+
			"<TRNTYPE>DEP\n"+
			"<DTPOSTED>20240901012345.000[+8:CST]\n"+
			"<TRNAMT>123.45\n"+
			"</STMTTRN>\n"+
			"</BANKTRANLIST>\n"+
			"</STMTRS>\n"+
			"</STMTTRNRS>\n"+
			"</BANKMSGSRSV1>\n"+
			"</OFX>"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidOFXFile.Message)
}

func TestCreateNewOFXFileReader_OFX2WithoutAnyHeader(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewOFXFileReader(context, []byte(
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
			"        </BANKTRANLIST>\n"+
			"      </STMTRS>\n"+
			"    </STMTTRNRS>\n"+
			"  </BANKMSGSRSV1>\n"+
			"</OFX>"))

	assert.Nil(t, err)

	ofxFile, err := reader.read(context)
	assert.Nil(t, err)
	assert.NotNil(t, ofxFile)
	assert.Nil(t, ofxFile.FileHeader)

	assert.NotNil(t, ofxFile.BankMessageResponseV1)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse)
	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse)

	assert.Equal(t, "CNY", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.DefaultCurrency)

	assert.NotNil(t, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom)
	assert.Equal(t, "123", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.AccountFrom.AccountId)

	assert.Equal(t, 1, len(ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions))
	assert.Equal(t, ofxDepositTransaction, ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].TransactionType)
	assert.Equal(t, "20240901012345.000[+8:CST]", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].PostedDate)
	assert.Equal(t, "123.45", ofxFile.BankMessageResponseV1.StatementTransactionResponse.StatementResponse.TransactionList.StatementTransactions[0].Amount)
}
