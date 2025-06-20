package mt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestMT940DataReaderParse(t *testing.T) {
	reader := &mt940DataReader{
		allLines: []string{
			"{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:",
			":20:MT940-2025001",
			":21:RELATEDREFERENCE",
			":25:123456789",
			":28C:123/1",
			":60F:C250601CNY1234,56",
			":61:2506010602DY123,45NTRFTEST//ABC123456",
			":86:First Transaction",
			"Additional Info",
			":61:2506020620CY234,56NSTFFOOBAR//DEF789012",
			":86:Second Transaction",
			"More Info",
			":62F:C250602CNY2345,67",
			":64:C250602CNY2345,67",
			"-}",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, "MT940-2025001", actualData.StatementReferenceNumber)
	assert.Equal(t, "RELATEDREFERENCE", actualData.RelatedReference)
	assert.Equal(t, "123456789", actualData.AccountId)
	assert.Equal(t, "123/1", actualData.SequentialNumber)

	assert.Equal(t, MT_MARK_CREDIT, actualData.OpeningBalance.DebitCreditMark)
	assert.Equal(t, "250601", actualData.OpeningBalance.Date)
	assert.Equal(t, "CNY", actualData.OpeningBalance.Currency)
	assert.Equal(t, "1234,56", actualData.OpeningBalance.Amount)

	assert.Equal(t, 2, len(actualData.Statements))

	assert.Equal(t, "250601", actualData.Statements[0].ValueDate)
	assert.Equal(t, "0602", actualData.Statements[0].EntryDate)
	assert.Equal(t, MT_MARK_DEBIT, actualData.Statements[0].CreditDebitMark)
	assert.Equal(t, "Y", actualData.Statements[0].FundsCode)
	assert.Equal(t, "123,45", actualData.Statements[0].Amount)
	assert.Equal(t, "NTRF", actualData.Statements[0].TransactionTypeIdentificationCode)
	assert.Equal(t, "TEST", actualData.Statements[0].ReferenceForAccountOwner)
	assert.Equal(t, "ABC123456", actualData.Statements[0].ReferenceOfAccountServicingInstitution)
	assert.Equal(t, "First Transaction", actualData.Statements[0].InformationToAccountOwner[0])
	assert.Equal(t, "Additional Info", actualData.Statements[0].InformationToAccountOwner[1])

	assert.Equal(t, "250602", actualData.Statements[1].ValueDate)
	assert.Equal(t, "0620", actualData.Statements[1].EntryDate)
	assert.Equal(t, MT_MARK_CREDIT, actualData.Statements[1].CreditDebitMark)
	assert.Equal(t, "Y", actualData.Statements[0].FundsCode)
	assert.Equal(t, "234,56", actualData.Statements[1].Amount)
	assert.Equal(t, "NSTF", actualData.Statements[1].TransactionTypeIdentificationCode)
	assert.Equal(t, "FOOBAR", actualData.Statements[1].ReferenceForAccountOwner)
	assert.Equal(t, "DEF789012", actualData.Statements[1].ReferenceOfAccountServicingInstitution)
	assert.Equal(t, "Second Transaction", actualData.Statements[1].InformationToAccountOwner[0])
	assert.Equal(t, "More Info", actualData.Statements[1].InformationToAccountOwner[1])

	assert.Equal(t, MT_MARK_CREDIT, actualData.ClosingBalance.DebitCreditMark)
	assert.Equal(t, "250602", actualData.ClosingBalance.Date)
	assert.Equal(t, "CNY", actualData.ClosingBalance.Currency)
	assert.Equal(t, "2345,67", actualData.ClosingBalance.Amount)

	assert.Equal(t, MT_MARK_CREDIT, actualData.ClosingAvailableBalance.DebitCreditMark)
	assert.Equal(t, "250602", actualData.ClosingAvailableBalance.Date)
	assert.Equal(t, "CNY", actualData.ClosingAvailableBalance.Currency)
	assert.Equal(t, "2345,67", actualData.ClosingAvailableBalance.Amount)
}

func TestMT940DataReaderParse_NoBlockHeaderFooter(t *testing.T) {
	reader := &mt940DataReader{
		allLines: []string{
			":20:MT940-2025001",
			":25:123456789",
			":28C:123/1",
			":60F:C250601CNY1234,56",
			":61:2506010602DY123,45NTRFTEST//ABC123456",
			":86:First Transaction",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, "MT940-2025001", actualData.StatementReferenceNumber)
	assert.Equal(t, "123456789", actualData.AccountId)
	assert.Equal(t, "123/1", actualData.SequentialNumber)

	assert.Equal(t, MT_MARK_CREDIT, actualData.OpeningBalance.DebitCreditMark)
	assert.Equal(t, "250601", actualData.OpeningBalance.Date)
	assert.Equal(t, "CNY", actualData.OpeningBalance.Currency)
	assert.Equal(t, "1234,56", actualData.OpeningBalance.Amount)

	assert.Equal(t, 1, len(actualData.Statements))

	assert.Equal(t, "250601", actualData.Statements[0].ValueDate)
	assert.Equal(t, "0602", actualData.Statements[0].EntryDate)
	assert.Equal(t, MT_MARK_DEBIT, actualData.Statements[0].CreditDebitMark)
	assert.Equal(t, "Y", actualData.Statements[0].FundsCode)
	assert.Equal(t, "123,45", actualData.Statements[0].Amount)
	assert.Equal(t, "NTRF", actualData.Statements[0].TransactionTypeIdentificationCode)
	assert.Equal(t, "TEST", actualData.Statements[0].ReferenceForAccountOwner)
	assert.Equal(t, "ABC123456", actualData.Statements[0].ReferenceOfAccountServicingInstitution)
	assert.Equal(t, "First Transaction", actualData.Statements[0].InformationToAccountOwner[0])
}

func TestMT940DataReaderParse_ReferenceForTheAccountOwnerTwoLine(t *testing.T) {
	reader := &mt940DataReader{
		allLines: []string{
			":61:250601D123,45NTRFABCDEFGHIJKLMNOP",
			"QRSTUVWXYZ",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(actualData.Statements))

	assert.Equal(t, "250601", actualData.Statements[0].ValueDate)
	assert.Equal(t, MT_MARK_DEBIT, actualData.Statements[0].CreditDebitMark)
	assert.Equal(t, "123,45", actualData.Statements[0].Amount)
	assert.Equal(t, "NTRF", actualData.Statements[0].TransactionTypeIdentificationCode)
	assert.Equal(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", actualData.Statements[0].ReferenceForAccountOwner)
}

func TestMT940DataReaderParse_InformationToAccountOwnerSixLine(t *testing.T) {
	reader := &mt940DataReader{
		allLines: []string{
			":61:250601D123,45NTRFTEST",
			":86:Additional Info Line 1",
			"Additional Info Line 2",
			"Additional Info Line 3",
			"Additional Info Line 4",
			"Additional Info Line 5",
			"Additional Info Line 6",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(actualData.Statements))

	assert.Equal(t, "250601", actualData.Statements[0].ValueDate)
	assert.Equal(t, MT_MARK_DEBIT, actualData.Statements[0].CreditDebitMark)
	assert.Equal(t, "123,45", actualData.Statements[0].Amount)
	assert.Equal(t, "NTRF", actualData.Statements[0].TransactionTypeIdentificationCode)
	assert.Equal(t, "TEST", actualData.Statements[0].ReferenceForAccountOwner)
	assert.Equal(t, 6, len(actualData.Statements[0].InformationToAccountOwner))
	assert.Equal(t, "Additional Info Line 1", actualData.Statements[0].InformationToAccountOwner[0])
	assert.Equal(t, "Additional Info Line 2", actualData.Statements[0].InformationToAccountOwner[1])
	assert.Equal(t, "Additional Info Line 3", actualData.Statements[0].InformationToAccountOwner[2])
	assert.Equal(t, "Additional Info Line 4", actualData.Statements[0].InformationToAccountOwner[3])
	assert.Equal(t, "Additional Info Line 5", actualData.Statements[0].InformationToAccountOwner[4])
	assert.Equal(t, "Additional Info Line 6", actualData.Statements[0].InformationToAccountOwner[5])
}

func TestMT940DataReaderParse_InformationToAccountOwnerMoreThanSixLine(t *testing.T) {
	reader := &mt940DataReader{
		allLines: []string{
			":61:250601D123,45NTRFTEST",
			":86:Additional Info Line 1",
			"Additional Info Line 2",
			"Additional Info Line 3",
			"Additional Info Line 4",
			"Additional Info Line 5",
			"Additional Info Line 6",
			"Additional Info Line 7",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(actualData.Statements))

	assert.Equal(t, "250601", actualData.Statements[0].ValueDate)
	assert.Equal(t, MT_MARK_DEBIT, actualData.Statements[0].CreditDebitMark)
	assert.Equal(t, "123,45", actualData.Statements[0].Amount)
	assert.Equal(t, "NTRF", actualData.Statements[0].TransactionTypeIdentificationCode)
	assert.Equal(t, "TEST", actualData.Statements[0].ReferenceForAccountOwner)
	assert.Equal(t, 6, len(actualData.Statements[0].InformationToAccountOwner))
	assert.Equal(t, "Additional Info Line 1", actualData.Statements[0].InformationToAccountOwner[0])
	assert.Equal(t, "Additional Info Line 2", actualData.Statements[0].InformationToAccountOwner[1])
	assert.Equal(t, "Additional Info Line 3", actualData.Statements[0].InformationToAccountOwner[2])
	assert.Equal(t, "Additional Info Line 4", actualData.Statements[0].InformationToAccountOwner[3])
	assert.Equal(t, "Additional Info Line 5", actualData.Statements[0].InformationToAccountOwner[4])
	assert.Equal(t, "Additional Info Line 6", actualData.Statements[0].InformationToAccountOwner[5])
}

func TestMT940DataReaderParse_DuplicateBlockHeader(t *testing.T) {
	reader := &mt940DataReader{
		allLines: []string{
			"{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:",
			":20:MT940-2025001",
			":25:123456789",
			":28C:123/1",
			":60F:C250601CNY1234,56",
			"{1:F01TESTBANK123456789}{2:I940TESTBANK}{4:",
			":61:2506010602DY123,45NTRFTEST//ABC123456",
			":86:First Transaction",
			"-}",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, "", actualData.StatementReferenceNumber)
	assert.Equal(t, "", actualData.AccountId)
	assert.Equal(t, "", actualData.SequentialNumber)

	assert.Nil(t, actualData.OpeningBalance)

	assert.Equal(t, 1, len(actualData.Statements))

	assert.Equal(t, "250601", actualData.Statements[0].ValueDate)
	assert.Equal(t, "0602", actualData.Statements[0].EntryDate)
	assert.Equal(t, MT_MARK_DEBIT, actualData.Statements[0].CreditDebitMark)
	assert.Equal(t, "Y", actualData.Statements[0].FundsCode)
	assert.Equal(t, "123,45", actualData.Statements[0].Amount)
	assert.Equal(t, "NTRF", actualData.Statements[0].TransactionTypeIdentificationCode)
	assert.Equal(t, "TEST", actualData.Statements[0].ReferenceForAccountOwner)
	assert.Equal(t, "ABC123456", actualData.Statements[0].ReferenceOfAccountServicingInstitution)
}

func TestMT940DataReaderParse_EmptyContent(t *testing.T) {
	reader := &mt940DataReader{
		allLines: []string{},
	}
	context := core.NewNullContext()

	_, err := reader.read(context)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestMT940DataReaderParseBalance_ValidBalance(t *testing.T) {
	reader := &mt940DataReader{}
	context := core.NewNullContext()

	balance, err := reader.parseBalance(context, "C250601CNY1234,56")
	assert.Nil(t, err)
	assert.Equal(t, MT_MARK_CREDIT, balance.DebitCreditMark)
	assert.Equal(t, "250601", balance.Date)
	assert.Equal(t, "CNY", balance.Currency)
	assert.Equal(t, "1234,56", balance.Amount)

	balance, err = reader.parseBalance(context, "D250602USD2345,67")
	assert.Nil(t, err)
	assert.Equal(t, MT_MARK_DEBIT, balance.DebitCreditMark)
	assert.Equal(t, "250602", balance.Date)
	assert.Equal(t, "USD", balance.Currency)
	assert.Equal(t, "2345,67", balance.Amount)
}

func TestMT940DataReaderParseBalance_InvalidBalance(t *testing.T) {
	reader := &mt940DataReader{}
	context := core.NewNullContext()

	_, err := reader.parseBalance(context, "X250601CNY1234,56")
	assert.EqualError(t, err, errs.ErrTransactionTypeInvalid.Message)

	_, err = reader.parseBalance(context, "C")
	assert.EqualError(t, err, errs.ErrInvalidMT940File.Message)
}

func TestMT940DataReaderParseStatement_ValidFields(t *testing.T) {
	reader := &mt940DataReader{}
	context := core.NewNullContext()

	statement, err := reader.parseStatement(context, "2506010602RDY123,45NTRFTEST//ABC123456")
	assert.Nil(t, err)
	assert.Equal(t, "250601", statement.ValueDate)
	assert.Equal(t, "0602", statement.EntryDate)
	assert.Equal(t, MT_MARK_REVERSAL_DEBIT, statement.CreditDebitMark)
	assert.Equal(t, "Y", statement.FundsCode)
	assert.Equal(t, "123,45", statement.Amount)
	assert.Equal(t, "NTRF", statement.TransactionTypeIdentificationCode)
	assert.Equal(t, "TEST", statement.ReferenceForAccountOwner)
	assert.Equal(t, "ABC123456", statement.ReferenceOfAccountServicingInstitution)

	statement, err = reader.parseStatement(context, "250601RC234,56NSTFFOOBAR")
	assert.Nil(t, err)
	assert.Equal(t, "250601", statement.ValueDate)
	assert.Equal(t, "", statement.EntryDate)
	assert.Equal(t, MT_MARK_REVERSAL_CREDIT, statement.CreditDebitMark)
	assert.Equal(t, "234,56", statement.Amount)
	assert.Equal(t, "NSTF", statement.TransactionTypeIdentificationCode)
	assert.Equal(t, "FOOBAR", statement.ReferenceForAccountOwner)
}

func TestMT940DataReaderParseStatement_InvalidField(t *testing.T) {
	reader := &mt940DataReader{}
	context := core.NewNullContext()

	_, err := reader.parseStatement(context, "250601X123,45NTRFTest")
	assert.EqualError(t, err, errs.ErrTransactionTypeInvalid.Message)
}

func TestMT940DataReaderParseStatement_MissingField(t *testing.T) {
	reader := &mt940DataReader{}
	context := core.NewNullContext()

	// Missing entry date
	_, err := reader.parseStatement(context, "2406")
	assert.EqualError(t, err, errs.ErrInvalidMT940File.Message)

	// Missing debit/credit mark
	_, err = reader.parseStatement(context, "250601060234,56NTRFTEST//ABC123456")
	assert.EqualError(t, err, errs.ErrTransactionTypeInvalid.Message)

	// Missing amount
	_, err = reader.parseStatement(context, "250601DNTRFTEST//ABC123456")
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	// Missing transaction type identification code
	_, err = reader.parseStatement(context, "250601D234,56TEST//ABC123456")
	assert.EqualError(t, err, errs.ErrInvalidMT940File.Message)

	// Missing reference for account owner
	_, err = reader.parseStatement(context, "250601D234,56NTRF//ABC123456")
	assert.EqualError(t, err, errs.ErrInvalidMT940File.Message)
}
