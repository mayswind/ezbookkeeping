package beancount

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestBeancountDataReaderRead(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"; Test Beancount Data\n"+
		"option \"name_assets\" \"AssetsAccount\"\n"+
		"option \"name_liabilities\" \"LiabilitiesAccount\"\n"+
		"option \"name_equity\" \"EquityAccount\"\n"+
		"option \"name_income\" \"IncomeAccount\"\n"+
		"option \"name_expenses\" \"ExpensesAccount\"\n"+
		"\n"+
		"2024-01-01 open AssetsAccount:TestAccount\n"+
		"2024-01-02 open LiabilitiesAccount:TestAccount2\n"+
		"2024-01-03 open EquityAccount:Opening-Balances\n"+
		"\n"+
		"; The following transactions with tag1 and tag2\n"+
		"pushtag #tag1\n"+
		"pushtag #tag2\n"+
		"\n"+
		"2024-01-05 * \"Payee Name\" \"Foo Bar\" #tag3 #tag4 ^test-link\n"+
		"  IncomeAccount:TestCategory -123.45 CNY\n"+
		"  AssetsAccount:TestAccount 123.45 CNY\n"+
		"; The following transactions with tag2\n"+
		"poptag #tag1\n"+
		"2024-01-06 * \"test\n#test2\" #tag5 #tag6 ^test-link2\n"+
		"  LiabilitiesAccount:TestAccount2 -0.12 USD\n"+
		"  ExpensesAccount:TestCategory2 0.12 USD\n"+
		"2024-01-07 close AssetsAccount:TestAccount\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 5, len(actualData.accounts))
	assert.Equal(t, "AssetsAccount:TestAccount", actualData.accounts["AssetsAccount:TestAccount"].name)
	assert.Equal(t, beancountAssetsAccountType, actualData.accounts["AssetsAccount:TestAccount"].accountType)
	assert.Equal(t, "2024-01-01", actualData.accounts["AssetsAccount:TestAccount"].openDate)
	assert.Equal(t, "2024-01-07", actualData.accounts["AssetsAccount:TestAccount"].closeDate)

	assert.Equal(t, "LiabilitiesAccount:TestAccount2", actualData.accounts["LiabilitiesAccount:TestAccount2"].name)
	assert.Equal(t, beancountLiabilitiesAccountType, actualData.accounts["LiabilitiesAccount:TestAccount2"].accountType)
	assert.Equal(t, "2024-01-02", actualData.accounts["LiabilitiesAccount:TestAccount2"].openDate)

	assert.Equal(t, 2, len(actualData.transactions))

	assert.Equal(t, "2024-01-05", actualData.transactions[0].date)
	assert.Equal(t, "Payee Name", actualData.transactions[0].payee)
	assert.Equal(t, "Foo Bar", actualData.transactions[0].narration)
	assert.Equal(t, 2, len(actualData.transactions[0].postings))
	assert.Equal(t, "IncomeAccount:TestCategory", actualData.transactions[0].postings[0].account)
	assert.Equal(t, "-123.45", actualData.transactions[0].postings[0].amount)
	assert.Equal(t, "CNY", actualData.transactions[0].postings[0].commodity)
	assert.Equal(t, "AssetsAccount:TestAccount", actualData.transactions[0].postings[1].account)
	assert.Equal(t, "123.45", actualData.transactions[0].postings[1].amount)
	assert.Equal(t, "CNY", actualData.transactions[0].postings[1].commodity)

	assert.Equal(t, 4, len(actualData.transactions[0].tags))
	assert.Equal(t, actualData.transactions[0].tags[0], "tag1")
	assert.Equal(t, actualData.transactions[0].tags[1], "tag2")
	assert.Equal(t, actualData.transactions[0].tags[2], "tag3")
	assert.Equal(t, actualData.transactions[0].tags[3], "tag4")

	assert.Equal(t, 1, len(actualData.transactions[0].links))
	assert.Equal(t, actualData.transactions[0].links[0], "test-link")

	assert.Equal(t, "2024-01-06", actualData.transactions[1].date)
	assert.Equal(t, "", actualData.transactions[1].payee)
	assert.Equal(t, "test\n#test2", actualData.transactions[1].narration)
	assert.Equal(t, 2, len(actualData.transactions[1].postings))
	assert.Equal(t, "LiabilitiesAccount:TestAccount2", actualData.transactions[1].postings[0].account)
	assert.Equal(t, "-0.12", actualData.transactions[1].postings[0].amount)
	assert.Equal(t, "USD", actualData.transactions[1].postings[0].commodity)
	assert.Equal(t, "ExpensesAccount:TestCategory2", actualData.transactions[1].postings[1].account)
	assert.Equal(t, "0.12", actualData.transactions[1].postings[1].amount)
	assert.Equal(t, "USD", actualData.transactions[1].postings[1].commodity)

	assert.Equal(t, 3, len(actualData.transactions[1].tags))
	assert.Equal(t, actualData.transactions[1].tags[0], "tag2")
	assert.Equal(t, actualData.transactions[1].tags[1], "tag5")
	assert.Equal(t, actualData.transactions[1].tags[2], "tag6")

	assert.Equal(t, 1, len(actualData.transactions[1].links))
	assert.Equal(t, actualData.transactions[1].links[0], "test-link2")
}

func TestBeancountDataReaderRead_EmptyContent(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestBeancountDataReaderRead_UnsupportedInclude(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte("include \"other.beancount\""))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrBeancountFileNotSupportInclude.Message)
}

func TestBeancountDataReaderRead_SkipUnsupportedDirective(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"plugin \"beancount.plugins.plugin_name\"\n"+
		"unknown directive\n"+
		"2024-01-01 commodity USD\n"+
		"2024-01-01 price USD 1.08 CAD\n"+
		"2024-01-01 note Assets:Test \"some text\"\n"+
		"2024-01-01 document Assets:Test \"scheme://path\"\n"+
		"2024-01-01 event \"location\" \"address\"\n"+
		"2024-01-01 balance Assets:Test 100.00 USD\n"+
		"2024-01-01 pad Assets:Test Equity:Opening-Balances\n"+
		"2024-01-01 query \"Name\" \"\nSELECT FIELDS FROM TABLE\"\n"+
		"2024-01-01 custom \"Type\" \"Value\"\n"+
		"2024-01-01 unknown directive\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.Nil(t, err)
}

func TestBeancountDataReaderReadAndSetOption_AccountTypeName(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"option \"name_assets\" \"A\"\n"+
		"option \"name_liabilities\" \"L\"\n"+
		"option \"name_equity\" \"E\"\n"+
		"\n"+
		"2024-01-01 open A:TestAccount\n"+
		"2024-01-02 open L:TestAccount2\n"+
		"2024-01-03 open E:Opening-Balances\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(actualData.accounts))

	assert.Equal(t, "A:TestAccount", actualData.accounts["A:TestAccount"].name)
	assert.Equal(t, beancountAssetsAccountType, actualData.accounts["A:TestAccount"].accountType)

	assert.Equal(t, "L:TestAccount2", actualData.accounts["L:TestAccount2"].name)
	assert.Equal(t, beancountLiabilitiesAccountType, actualData.accounts["L:TestAccount2"].accountType)

	assert.Equal(t, "E:Opening-Balances", actualData.accounts["E:Opening-Balances"].name)
	assert.Equal(t, beancountEquityAccountType, actualData.accounts["E:Opening-Balances"].accountType)
	assert.True(t, actualData.accounts["E:Opening-Balances"].isOpeningBalanceEquityAccount())
}

func TestBeancountDataReaderReadAndSetOption_InvalidLineOrUnsupportedOption(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"option \"test\" \"Test\" \"Test2\"\n"+
		"option \"test\" \"Test\"\n"+
		"option \"test\"\n"+
		"option \n"+
		"option\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.Nil(t, err)
}

func TestBeancountDataReaderReadAndSetTags(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"pushtag #tag1\n"+
		"pushtag #tag2\n"+
		"pushtag #tag2\n"+
		"pushtag #tag1\n"+
		"\n"+
		"2024-01-01 * #tag3 #tag4\n"+
		"poptag #tag1\n"+
		"poptag #tag2\n"+
		"pushtag\n"+
		"pushtag    \n"+
		"pushtag tag\n"+
		"2024-01-02 * #tag5 #tag6\n"+
		"poptag #tag1\n"+
		"poptag #tag2\n"+
		"poptag\n"+
		"poptag \n"+
		"2024-01-03 * #tag5 #tag6\n"+
		"pushtag #tag3\n"+
		"pushtag #tag6\n"+
		"2024-01-04 * #tag5 #tag6\n"+
		"2024-01-05 * #tag5 #tag6 #tag6 #tag5\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 5, len(actualData.transactions))

	assert.Equal(t, 4, len(actualData.transactions[0].tags))
	assert.Equal(t, actualData.transactions[0].tags[0], "tag1")
	assert.Equal(t, actualData.transactions[0].tags[1], "tag2")
	assert.Equal(t, actualData.transactions[0].tags[2], "tag3")
	assert.Equal(t, actualData.transactions[0].tags[3], "tag4")

	assert.Equal(t, 2, len(actualData.transactions[1].tags))
	assert.Equal(t, actualData.transactions[1].tags[0], "tag5")
	assert.Equal(t, actualData.transactions[1].tags[1], "tag6")

	assert.Equal(t, 2, len(actualData.transactions[2].tags))
	assert.Equal(t, actualData.transactions[2].tags[0], "tag5")
	assert.Equal(t, actualData.transactions[2].tags[1], "tag6")

	assert.Equal(t, 3, len(actualData.transactions[3].tags))
	assert.Equal(t, actualData.transactions[3].tags[0], "tag3")
	assert.Equal(t, actualData.transactions[3].tags[1], "tag6")
	assert.Equal(t, actualData.transactions[3].tags[2], "tag5")

	assert.Equal(t, 3, len(actualData.transactions[4].tags))
	assert.Equal(t, actualData.transactions[4].tags[0], "tag3")
	assert.Equal(t, actualData.transactions[4].tags[1], "tag6")
	assert.Equal(t, actualData.transactions[4].tags[2], "tag5")
}

func TestBeancountDataReaderReadAccountLine_InvalidLine(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 open\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(actualData.accounts))
}

func TestBeancountDataReaderReadAccountLine_InvalidAccountType(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 open Test:TestAccount\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)

	reader, err = createNewBeancountDataReader(context, []byte(""+
		"option \"name_assets\" \"A\"\n"+
		"\n"+
		"2024-01-01 open Assets:TestAccount\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)
}

func TestBeancountDataReaderReadTransactionLine(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"2024-01-02 * \"test\ttest2\ntest3\" ; some comment\n"+
		"2024-01-03 ! \"test\" \"test2\"\n"+
		"2024-01-04 P \"test\" #tag #tag2 ; some comment\n"+
		"2024-01-05 txn \"test\" ^scheme://path/to/test/link ; some comment\n"+
		"2024-01-06 txn ; \"test\" \"test2\" #tag ^link\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 6, len(actualData.transactions))

	assert.Equal(t, "2024-01-01", actualData.transactions[0].date)
	assert.Equal(t, beancountDirectiveCompletedTransaction, actualData.transactions[0].directive)
	assert.Equal(t, "", actualData.transactions[0].payee)
	assert.Equal(t, "", actualData.transactions[0].narration)

	assert.Equal(t, "2024-01-02", actualData.transactions[1].date)
	assert.Equal(t, beancountDirectiveCompletedTransaction, actualData.transactions[1].directive)
	assert.Equal(t, "", actualData.transactions[1].payee)
	assert.Equal(t, "test\ttest2\ntest3", actualData.transactions[1].narration)

	assert.Equal(t, "2024-01-03", actualData.transactions[2].date)
	assert.Equal(t, beancountDirectiveInCompleteTransaction, actualData.transactions[2].directive)
	assert.Equal(t, "test", actualData.transactions[2].payee)
	assert.Equal(t, "test2", actualData.transactions[2].narration)

	assert.Equal(t, "2024-01-04", actualData.transactions[3].date)
	assert.Equal(t, beancountDirectivePaddingTransaction, actualData.transactions[3].directive)
	assert.Equal(t, "", actualData.transactions[3].payee)
	assert.Equal(t, "test", actualData.transactions[3].narration)

	assert.Equal(t, 2, len(actualData.transactions[3].tags))
	assert.Equal(t, actualData.transactions[3].tags[0], "tag")
	assert.Equal(t, actualData.transactions[3].tags[1], "tag2")

	assert.Equal(t, "2024-01-05", actualData.transactions[4].date)
	assert.Equal(t, beancountDirectiveTransaction, actualData.transactions[4].directive)
	assert.Equal(t, "", actualData.transactions[4].payee)
	assert.Equal(t, "test", actualData.transactions[4].narration)

	assert.Equal(t, 1, len(actualData.transactions[4].links))
	assert.Equal(t, actualData.transactions[4].links[0], "scheme://path/to/test/link")

	assert.Equal(t, "2024-01-06", actualData.transactions[5].date)
	assert.Equal(t, beancountDirectiveTransaction, actualData.transactions[5].directive)
	assert.Equal(t, "", actualData.transactions[5].payee)
	assert.Equal(t, "", actualData.transactions[5].narration)
}

func TestBeancountDataReaderReadTransactionPostingLine(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Income:TestCategory -123.45 CNY ; some comment\n"+
		"  Assets:TestAccount 123.45   CNY\n"+
		"2024-01-02 *\n"+
		"  Liabilities:TestAccount2 -0.23 USD ; some comment\n"+
		"  Expenses:TestCategory2    0.12 USD @@ 0.84 CNY\n"+
		"  Expenses:TestCategory3    0.11 USD @ 7.12 CNY\n"+
		"  ! Expenses:TestCategory4  0.00 USD {0.00 CNY}\n"+
		"  Expenses:TestCategory5 \n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(actualData.transactions))

	assert.Equal(t, "2024-01-01", actualData.transactions[0].date)
	assert.Equal(t, 2, len(actualData.transactions[0].postings))
	assert.Equal(t, "Income:TestCategory", actualData.transactions[0].postings[0].account)
	assert.Equal(t, "-123.45", actualData.transactions[0].postings[0].amount)
	assert.Equal(t, "CNY", actualData.transactions[0].postings[0].commodity)

	assert.Equal(t, "Assets:TestAccount", actualData.transactions[0].postings[1].account)
	assert.Equal(t, "123.45", actualData.transactions[0].postings[1].amount)
	assert.Equal(t, "CNY", actualData.transactions[0].postings[1].commodity)

	assert.Equal(t, "2024-01-02", actualData.transactions[1].date)
	assert.Equal(t, 4, len(actualData.transactions[1].postings))

	assert.Equal(t, "Liabilities:TestAccount2", actualData.transactions[1].postings[0].account)
	assert.Equal(t, "-0.23", actualData.transactions[1].postings[0].amount)
	assert.Equal(t, "USD", actualData.transactions[1].postings[0].commodity)
	assert.Equal(t, "Expenses:TestCategory2", actualData.transactions[1].postings[1].account)

	assert.Equal(t, "0.12", actualData.transactions[1].postings[1].amount)
	assert.Equal(t, "USD", actualData.transactions[1].postings[1].commodity)
	assert.Equal(t, "0.84", actualData.transactions[1].postings[1].totalCost)
	assert.Equal(t, "CNY", actualData.transactions[1].postings[1].totalCostCommodity)
	assert.Equal(t, "Expenses:TestCategory3", actualData.transactions[1].postings[2].account)

	assert.Equal(t, "0.11", actualData.transactions[1].postings[2].amount)
	assert.Equal(t, "USD", actualData.transactions[1].postings[2].commodity)
	assert.Equal(t, "7.12", actualData.transactions[1].postings[2].price)
	assert.Equal(t, "CNY", actualData.transactions[1].postings[2].priceCommodity)

	assert.Equal(t, "0.00", actualData.transactions[1].postings[3].amount)
	assert.Equal(t, "USD", actualData.transactions[1].postings[3].commodity)
}

func TestBeancountDataReaderReadTransactionPostingLine_AmountExpression(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Income:TestCategory (1.2-3.4) * 5.6 / 7.8 CNY\n"+
		"  Assets:TestAccount 1.2 * 3.4/-5.6 - 7.8 CNY\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(actualData.transactions))

	assert.Equal(t, "2024-01-01", actualData.transactions[0].date)
	assert.Equal(t, 2, len(actualData.transactions[0].postings))
	assert.Equal(t, "Income:TestCategory", actualData.transactions[0].postings[0].account)
	assert.Equal(t, "(1.2-3.4) * 5.6 / 7.8", actualData.transactions[0].postings[0].originalAmount)
	assert.Equal(t, "-1.58", actualData.transactions[0].postings[0].amount)
	assert.Equal(t, "CNY", actualData.transactions[0].postings[0].commodity)

	assert.Equal(t, "Assets:TestAccount", actualData.transactions[0].postings[1].account)
	assert.Equal(t, "1.2 * 3.4/-5.6 - 7.8", actualData.transactions[0].postings[1].originalAmount)
	assert.Equal(t, "-8.53", actualData.transactions[0].postings[1].amount)
	assert.Equal(t, "CNY", actualData.transactions[0].postings[1].commodity)
}

func TestBeancountDataReaderReadTransactionPostingLine_InvalidAmountExpression(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Income:TestCategory (1.2-3.4)*5.6/0 CNY\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)

	reader, err = createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Assets:TestAccount abc CNY\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestBeancountDataReaderReadTransactionPostingLine_InvalidAccountType(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Income:TestCategory -123.45 CNY\n"+
		"  Test:TestAccount 123.45 CNY\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)
}

func TestBeancountDataReaderReadTransactionPostingLine_InvalidCommodity(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Income:TestCategory -123.45 cny\n"+
		"  Assets:TestAccount 123.45 cny\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)
}

func TestBeancountDataReaderReadTransactionPostingLine_MissingAmount(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Assets:TestAccount\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(actualData.transactions))
	assert.Equal(t, 0, len(actualData.transactions[0].postings))

	reader, err = createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Assets:TestAccount \n"))
	assert.Nil(t, err)

	actualData, err = reader.read(context)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(actualData.transactions))
	assert.Equal(t, 0, len(actualData.transactions[0].postings))
}

func TestBeancountDataReaderReadTransactionPostingLine_MissingCommodity(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Assets:TestAccount 123.45\n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)

	reader, err = createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  Assets:TestAccount 123.45 \n"))
	assert.Nil(t, err)

	_, err = reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidBeancountFile.Message)
}

func TestBeancountDataReaderReadTransactionMetadataLine(t *testing.T) {
	context := core.NewNullContext()
	reader, err := createNewBeancountDataReader(context, []byte(""+
		"2024-01-01 *\n"+
		"  key: value\n"+
		"  key2: \"value 2\"\n"+
		"  key3: \n"+
		"  key4: \"\"\n"+
		"  key5 : \"\"\n"+
		"  key2: \"new value\"\n"+
		"  Income:TestCategory -123.45 CNY\n"+
		"  Assets:TestAccount 123.45 CNY\n"+
		"2024-01-02 *\n"+
		"  Liabilities:TestAccount2 -0.23 USD\n"+
		"    key6: value6\n"+
		"    key7: \"value 7\"\n"+
		"    key8: \n"+
		"    key9: \"\"\n"+
		"    key0 : \"\"\n"+
		"    key6: \"new value\"\n"+
		"  Expenses:TestCategory2 0.12 USD\n"))
	assert.Nil(t, err)

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(actualData.transactions))

	assert.Equal(t, "2024-01-01", actualData.transactions[0].date)
	assert.Equal(t, 2, len(actualData.transactions[0].postings))
	assert.Equal(t, 2, len(actualData.transactions[0].metadata))
	assert.Equal(t, "value", actualData.transactions[0].metadata["key"])
	assert.Equal(t, "value 2", actualData.transactions[0].metadata["key2"])

	assert.Equal(t, "2024-01-02", actualData.transactions[1].date)
	assert.Equal(t, 2, len(actualData.transactions[1].postings))
	assert.Equal(t, 2, len(actualData.transactions[1].postings[0].metadata))
	assert.Equal(t, "value6", actualData.transactions[1].postings[0].metadata["key6"])
	assert.Equal(t, "value 7", actualData.transactions[1].postings[0].metadata["key7"])
	assert.Equal(t, 0, len(actualData.transactions[1].postings[1].metadata))
}
