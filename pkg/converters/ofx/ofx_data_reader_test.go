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

func TestCreateNewOFXFileReader_OFX1WithBlanklines(t *testing.T) {
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
