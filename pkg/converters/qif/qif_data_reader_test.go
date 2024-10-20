package qif

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestQifDataReaderParse(t *testing.T) {
	reader := &qifDataReader{
		allLines: []string{
			"!Type:Bank",
			"D2024/10/9",
			"T-123.45",
			"^",
			"D2024/10/12",
			"T+234.56",
			"^",
			"!Type:Cash",
			"D2024/9/1",
			"T100.00",
			"POpening Balance",
			"L[Wallet]",
			"^",
			"!Type:Memorized",
			"KC",
			"T-123.45",
			"12024/10/13",
			"23",
			"31",
			"42",
			"512.34",
			"6100.45",
			"7234.56",
			"^",
			"!Type:Invst",
			"D2024/10/14",
			"NBuy",
			"YTest",
			"I12.34",
			"Q10",
			"T-123.4",
			"^",
			"!Account",
			"NTest Account",
			"^",
			"NWallet",
			"^",
			"!Type:Cat",
			"NTest Category",
			"I",
			"^",
			"!Type:Class",
			"NTest Class",
			"DFoo Bar",
			"^",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(actualData.bankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.bankAccountTransactions[0].date)
	assert.Equal(t, "-123.45", actualData.bankAccountTransactions[0].amount)
	assert.Equal(t, "2024/10/12", actualData.bankAccountTransactions[1].date)
	assert.Equal(t, "+234.56", actualData.bankAccountTransactions[1].amount)

	assert.Equal(t, 1, len(actualData.cashAccountTransactions))
	assert.Equal(t, "2024/9/1", actualData.cashAccountTransactions[0].date)
	assert.Equal(t, "100.00", actualData.cashAccountTransactions[0].amount)
	assert.Equal(t, "Opening Balance", actualData.cashAccountTransactions[0].payee)
	assert.Equal(t, "[Wallet]", actualData.cashAccountTransactions[0].category)

	assert.Equal(t, 1, len(actualData.memorizedTransactions))
	assert.Equal(t, qifCheckTransactionType, actualData.memorizedTransactions[0].transactionType)
	assert.Equal(t, "-123.45", actualData.memorizedTransactions[0].amount)
	assert.Equal(t, "2024/10/13", actualData.memorizedTransactions[0].amortization.firstPaymentDate)
	assert.Equal(t, "3", actualData.memorizedTransactions[0].amortization.totalYearsForLoan)
	assert.Equal(t, "1", actualData.memorizedTransactions[0].amortization.numberOfPayments)
	assert.Equal(t, "2", actualData.memorizedTransactions[0].amortization.numberOfPeriodsPerYear)
	assert.Equal(t, "12.34", actualData.memorizedTransactions[0].amortization.interestRate)
	assert.Equal(t, "100.45", actualData.memorizedTransactions[0].amortization.currentLoanBalance)
	assert.Equal(t, "234.56", actualData.memorizedTransactions[0].amortization.originalLoanAmount)

	assert.Equal(t, 1, len(actualData.investmentAccountTransactions))
	assert.Equal(t, "2024/10/14", actualData.investmentAccountTransactions[0].date)
	assert.Equal(t, "Buy", actualData.investmentAccountTransactions[0].action)
	assert.Equal(t, "Test", actualData.investmentAccountTransactions[0].security)
	assert.Equal(t, "12.34", actualData.investmentAccountTransactions[0].price)
	assert.Equal(t, "10", actualData.investmentAccountTransactions[0].quantity)
	assert.Equal(t, "-123.4", actualData.investmentAccountTransactions[0].amount)

	assert.Equal(t, 2, len(actualData.accounts))
	assert.Equal(t, "Test Account", actualData.accounts[0].name)
	assert.Equal(t, "Wallet", actualData.accounts[1].name)

	assert.Equal(t, 1, len(actualData.categories))
	assert.Equal(t, "Test Category", actualData.categories[0].name)
	assert.Equal(t, qifIncomeTransaction, actualData.categories[0].categoryType)

	assert.Equal(t, 1, len(actualData.classes))
	assert.Equal(t, "Test Class", actualData.classes[0].name)
	assert.Equal(t, "Foo Bar", actualData.classes[0].description)
}

func TestQifDataReaderParse_AccountEntryBeforeTransaction(t *testing.T) {
	reader := &qifDataReader{
		allLines: []string{
			"!Account",
			"NTest Account",
			"^",
			"!Type:Bank",
			"D2024/10/9",
			"T-123.45",
			"^",
			"D2024/10/12",
			"T+234.56",
			"^",
			"!Account",
			"NWallet",
			"^",
			"!Type:Cash",
			"D2024/9/1",
			"T100.00",
			"POpening Balance",
			"L[Wallet]",
			"^",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(actualData.bankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.bankAccountTransactions[0].date)
	assert.Equal(t, "-123.45", actualData.bankAccountTransactions[0].amount)
	assert.Equal(t, "2024/10/12", actualData.bankAccountTransactions[1].date)
	assert.Equal(t, "+234.56", actualData.bankAccountTransactions[1].amount)

	assert.Equal(t, 1, len(actualData.cashAccountTransactions))
	assert.Equal(t, "2024/9/1", actualData.cashAccountTransactions[0].date)
	assert.Equal(t, "100.00", actualData.cashAccountTransactions[0].amount)
	assert.Equal(t, "Opening Balance", actualData.cashAccountTransactions[0].payee)
	assert.Equal(t, "[Wallet]", actualData.cashAccountTransactions[0].category)

	assert.Equal(t, 2, len(actualData.accounts))
	assert.Equal(t, "Test Account", actualData.accounts[0].name)
	assert.Equal(t, "Wallet", actualData.accounts[1].name)
}

func TestQifDataReaderParse_EmptyContent(t *testing.T) {
	reader := &qifDataReader{
		allLines: []string{},
	}
	context := core.NewNullContext()

	_, err := reader.read(context)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestQifDataReaderParse_EmptyEntry(t *testing.T) {
	reader := &qifDataReader{
		allLines: []string{
			"!Type:Bank",
			"^",
			"!Type:Cash",
			"^",
			"!Type:Memorized",
			"^",
			"!Type:Invst",
			"^",
			"!Account",
			"^",
			"!Type:Cat",
			"^",
			"!Type:Class",
			"^",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(actualData.bankAccountTransactions))
	assert.Equal(t, 0, len(actualData.cashAccountTransactions))
	assert.Equal(t, 0, len(actualData.memorizedTransactions))
	assert.Equal(t, 0, len(actualData.investmentAccountTransactions))
	assert.Equal(t, 0, len(actualData.accounts))
	assert.Equal(t, 0, len(actualData.categories))
	assert.Equal(t, 0, len(actualData.classes))
}

func TestQifDataReaderParse_UnsupportedEntryHeader(t *testing.T) {
	reader := &qifDataReader{
		allLines: []string{
			"!Type:Bank",
			"D2024/10/9",
			"T-123.45",
			"^",
			"!Type:Unknown",
			"D2024/10/11",
			"T100.00",
			"^",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(actualData.bankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.bankAccountTransactions[0].date)
	assert.Equal(t, "-123.45", actualData.bankAccountTransactions[0].amount)
}

func TestQifDataReaderParse_UnsupportedLine(t *testing.T) {
	reader := &qifDataReader{
		allLines: []string{
			"!Type:Bank",
			"D2024/10/9",
			"T-123.45",
			"^",
			"!Option:Unknown",
			"D2024/10/11",
			"T100.00",
			"^",
		},
	}
	context := core.NewNullContext()

	actualData, err := reader.read(context)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(actualData.bankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.bankAccountTransactions[0].date)
	assert.Equal(t, "-123.45", actualData.bankAccountTransactions[0].amount)
	assert.Equal(t, "2024/10/11", actualData.bankAccountTransactions[1].date)
	assert.Equal(t, "100.00", actualData.bankAccountTransactions[1].amount)
}

func TestQifDataReaderParse_NewEntryHeaderAfterUnclosedEntry(t *testing.T) {
	reader := &qifDataReader{
		allLines: []string{
			"!Type:Bank",
			"D2024/10/9",
			"T-123.45",
			"!Type:Cash",
			"D2024/9/1",
			"T100.00",
			"POpening Balance",
			"L[Wallet]",
			"^",
		},
	}
	context := core.NewNullContext()

	_, err := reader.read(context)
	assert.EqualError(t, err, errs.ErrInvalidQIFFile.Message)
}

func TestQifDataReaderParseTransaction_SupportedFields(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualData, err := reader.parseTransaction(context, []string{
		"D2024/10/12",
		"T-123.45",
		"C",
		"N100",
		"PFoo",
		"MBar",
		"AAddress 1",
		"AAddress 2",
		"AAddress 3",
		"LTest Category",
		"SPart1 Category",
		"EPart1 Memo",
		"$-100.00",
		"SPart2 Category",
		"EPart2 Memo",
		"$-23.45",
	}, false)

	assert.Nil(t, err)
	assert.Equal(t, "2024/10/12", actualData.date)
	assert.Equal(t, "-123.45", actualData.amount)
	assert.Equal(t, qifClearedStatusUnreconciled, actualData.clearedStatus)
	assert.Equal(t, "100", actualData.num)
	assert.Equal(t, "Foo", actualData.payee)
	assert.Equal(t, "Bar", actualData.memo)
	assert.Equal(t, 3, len(actualData.addresses))
	assert.Equal(t, "Address 1", actualData.addresses[0])
	assert.Equal(t, "Address 2", actualData.addresses[1])
	assert.Equal(t, "Address 3", actualData.addresses[2])
	assert.Equal(t, "Test Category", actualData.category)
	assert.Equal(t, 2, len(actualData.subTransactionCategory))
	assert.Equal(t, "Part1 Category", actualData.subTransactionCategory[0])
	assert.Equal(t, "Part2 Category", actualData.subTransactionCategory[1])
	assert.Equal(t, 2, len(actualData.subTransactionMemo))
	assert.Equal(t, "Part1 Memo", actualData.subTransactionMemo[0])
	assert.Equal(t, "Part2 Memo", actualData.subTransactionMemo[1])
	assert.Equal(t, 2, len(actualData.subTransactionAmount))
	assert.Equal(t, "-100.00", actualData.subTransactionAmount[0])
	assert.Equal(t, "-23.45", actualData.subTransactionAmount[1])
}

func TestQifDataReaderParseMemorizedTransaction_SupportedFields(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualData, err := reader.parseMemorizedTransaction(context, []string{
		"KC",
		"D2024/10/12",
		"T-123.45",
		"C*",
		"N100",
		"PFoo",
		"MBar",
		"12024/10/13",
		"23",
		"31",
		"42",
		"512.34",
		"6100.45",
		"7234.56",
	})

	assert.Nil(t, err)
	assert.Equal(t, qifCheckTransactionType, actualData.transactionType)
	assert.Equal(t, "2024/10/12", actualData.date)
	assert.Equal(t, "-123.45", actualData.amount)
	assert.Equal(t, qifClearedStatusCleared, actualData.clearedStatus)
	assert.Equal(t, "100", actualData.num)
	assert.Equal(t, "Foo", actualData.payee)
	assert.Equal(t, "Bar", actualData.memo)
	assert.Equal(t, "2024/10/13", actualData.amortization.firstPaymentDate)
	assert.Equal(t, "3", actualData.amortization.totalYearsForLoan)
	assert.Equal(t, "1", actualData.amortization.numberOfPayments)
	assert.Equal(t, "2", actualData.amortization.numberOfPeriodsPerYear)
	assert.Equal(t, "12.34", actualData.amortization.interestRate)
	assert.Equal(t, "100.45", actualData.amortization.currentLoanBalance)
	assert.Equal(t, "234.56", actualData.amortization.originalLoanAmount)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KD"})
	assert.Nil(t, err)
	assert.Equal(t, qifDepositTransactionType, actualData.transactionType)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KP"})
	assert.Nil(t, err)
	assert.Equal(t, qifPaymentTransactionType, actualData.transactionType)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KI"})
	assert.Nil(t, err)
	assert.Equal(t, qifInvestmentTransactionType, actualData.transactionType)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KE"})
	assert.Nil(t, err)
	assert.Equal(t, qifElectronicPayeeTransactionType, actualData.transactionType)
}

func TestQifDataReaderParseInvestmentTransaction_SupportedFields(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualData, err := reader.parseInvestmentTransaction(context, []string{
		"D2024/10/12",
		"NBuy",
		"YTest",
		"I12.34",
		"Q10",
		"T-123.4",
		"CR",
		"PFoo",
		"MBar",
		"OTest2",
		"LAccount Name",
		"$100",
	})

	assert.Nil(t, err)
	assert.Equal(t, "2024/10/12", actualData.date)
	assert.Equal(t, "Buy", actualData.action)
	assert.Equal(t, "Test", actualData.security)
	assert.Equal(t, "12.34", actualData.price)
	assert.Equal(t, "10", actualData.quantity)
	assert.Equal(t, "-123.4", actualData.amount)
	assert.Equal(t, qifClearedStatusReconciled, actualData.clearedStatus)
	assert.Equal(t, "Foo", actualData.text)
	assert.Equal(t, "Bar", actualData.memo)
	assert.Equal(t, "Test2", actualData.commission)
	assert.Equal(t, "Account Name", actualData.accountForTransfer)
	assert.Equal(t, "100", actualData.amountTransferred)
}

func TestQifDataReaderParseAccount_SupportedFields(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualData, err := reader.parseAccount(context, []string{
		"NAccount Name",
		"TAccount Type",
		"DSome Text",
		"L1234.56",
		"/2024/10/12",
		"$123.45",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Account Name", actualData.name)
	assert.Equal(t, "Account Type", actualData.accountType)
	assert.Equal(t, "Some Text", actualData.description)
	assert.Equal(t, "1234.56", actualData.creditLimit)
	assert.Equal(t, "2024/10/12", actualData.statementBalanceDate)
	assert.Equal(t, "123.45", actualData.statementBalanceAmount)
}

func TestQifDataReaderParseCategory_SupportedFields(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualData, err := reader.parseCategory(context, []string{
		"NCategory Name:Sub Category Name",
		"DSome Text",
		"T",
		"I",
		"B123.45",
		"RTest",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Category Name:Sub Category Name", actualData.name)
	assert.Equal(t, "Some Text", actualData.description)
	assert.Equal(t, true, actualData.taxRelated)
	assert.Equal(t, qifIncomeTransaction, actualData.categoryType)
	assert.Equal(t, "123.45", actualData.budgetAmount)
	assert.Equal(t, "Test", actualData.taxScheduleInformation)

	actualData2, err := reader.parseCategory(context, []string{
		"NCategory Name:Sub Category Name",
		"DSome Text",
		"E",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Category Name:Sub Category Name", actualData2.name)
	assert.Equal(t, "Some Text", actualData2.description)
	assert.Equal(t, false, actualData2.taxRelated)
	assert.Equal(t, qifExpenseTransaction, actualData2.categoryType)

	actualData3, err := reader.parseCategory(context, []string{
		"NCategory Name:Sub Category Name",
		"DSome Text",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Category Name:Sub Category Name", actualData3.name)
	assert.Equal(t, "Some Text", actualData3.description)
	assert.Equal(t, qifExpenseTransaction, actualData3.categoryType)
}

func TestQifDataReaderParseClass_SupportedFields(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualData, err := reader.parseClass(context, []string{
		"NClass Name",
		"DSome Text",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Class Name", actualData.name)
	assert.Equal(t, "Some Text", actualData.description)
}

func TestQifDataReaderParse_UnsupportedFieldsOrValues(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualTransactionData, err := reader.parseTransaction(context, []string{
		"ZTest",
		"CZ",
		"",
	}, false)
	assert.Nil(t, err)
	assert.Equal(t, qifClearedStatusUnreconciled, actualTransactionData.clearedStatus)

	actualMemorizedTransactionData, err := reader.parseMemorizedTransaction(context, []string{
		"ZTest",
		"KZ",
		"",
	})
	assert.Nil(t, err)
	assert.Equal(t, qifInvalidTransactionType, actualMemorizedTransactionData.transactionType)

	_, err = reader.parseInvestmentTransaction(context, []string{
		"ZTest",
		"",
	})
	assert.Nil(t, err)

	_, err = reader.parseAccount(context, []string{
		"ZTest",
		"",
	})
	assert.Nil(t, err)

	_, err = reader.parseCategory(context, []string{
		"ZTest",
		"",
	})
	assert.Nil(t, err)

	_, err = reader.parseClass(context, []string{
		"ZTest",
		"",
	})
	assert.Nil(t, err)
}
