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

	assert.Equal(t, 2, len(actualData.BankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.BankAccountTransactions[0].Date)
	assert.Equal(t, "-123.45", actualData.BankAccountTransactions[0].Amount)
	assert.Equal(t, "2024/10/12", actualData.BankAccountTransactions[1].Date)
	assert.Equal(t, "+234.56", actualData.BankAccountTransactions[1].Amount)

	assert.Equal(t, 1, len(actualData.CashAccountTransactions))
	assert.Equal(t, "2024/9/1", actualData.CashAccountTransactions[0].Date)
	assert.Equal(t, "100.00", actualData.CashAccountTransactions[0].Amount)
	assert.Equal(t, "Opening Balance", actualData.CashAccountTransactions[0].Payee)
	assert.Equal(t, "[Wallet]", actualData.CashAccountTransactions[0].Category)

	assert.Equal(t, 1, len(actualData.MemorizedTransactions))
	assert.Equal(t, qifCheckTransactionType, actualData.MemorizedTransactions[0].TransactionType)
	assert.Equal(t, "-123.45", actualData.MemorizedTransactions[0].Amount)
	assert.Equal(t, "2024/10/13", actualData.MemorizedTransactions[0].Amortization.FirstPaymentDate)
	assert.Equal(t, "3", actualData.MemorizedTransactions[0].Amortization.TotalYearsForLoan)
	assert.Equal(t, "1", actualData.MemorizedTransactions[0].Amortization.NumberOfPayments)
	assert.Equal(t, "2", actualData.MemorizedTransactions[0].Amortization.NumberOfPeriodsPerYear)
	assert.Equal(t, "12.34", actualData.MemorizedTransactions[0].Amortization.InterestRate)
	assert.Equal(t, "100.45", actualData.MemorizedTransactions[0].Amortization.CurrentLoanBalance)
	assert.Equal(t, "234.56", actualData.MemorizedTransactions[0].Amortization.OriginalLoanAmount)

	assert.Equal(t, 1, len(actualData.InvestmentAccountTransactions))
	assert.Equal(t, "2024/10/14", actualData.InvestmentAccountTransactions[0].Date)
	assert.Equal(t, "Buy", actualData.InvestmentAccountTransactions[0].Action)
	assert.Equal(t, "Test", actualData.InvestmentAccountTransactions[0].Security)
	assert.Equal(t, "12.34", actualData.InvestmentAccountTransactions[0].Price)
	assert.Equal(t, "10", actualData.InvestmentAccountTransactions[0].Quantity)
	assert.Equal(t, "-123.4", actualData.InvestmentAccountTransactions[0].Amount)

	assert.Equal(t, 2, len(actualData.Accounts))
	assert.Equal(t, "Test Account", actualData.Accounts[0].Name)
	assert.Equal(t, "Wallet", actualData.Accounts[1].Name)

	assert.Equal(t, 1, len(actualData.Categories))
	assert.Equal(t, "Test Category", actualData.Categories[0].Name)
	assert.Equal(t, qifIncomeTransaction, actualData.Categories[0].CategoryType)

	assert.Equal(t, 1, len(actualData.Classes))
	assert.Equal(t, "Test Class", actualData.Classes[0].Name)
	assert.Equal(t, "Foo Bar", actualData.Classes[0].Description)
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

	assert.Equal(t, 2, len(actualData.BankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.BankAccountTransactions[0].Date)
	assert.Equal(t, "-123.45", actualData.BankAccountTransactions[0].Amount)
	assert.Equal(t, "2024/10/12", actualData.BankAccountTransactions[1].Date)
	assert.Equal(t, "+234.56", actualData.BankAccountTransactions[1].Amount)

	assert.Equal(t, 1, len(actualData.CashAccountTransactions))
	assert.Equal(t, "2024/9/1", actualData.CashAccountTransactions[0].Date)
	assert.Equal(t, "100.00", actualData.CashAccountTransactions[0].Amount)
	assert.Equal(t, "Opening Balance", actualData.CashAccountTransactions[0].Payee)
	assert.Equal(t, "[Wallet]", actualData.CashAccountTransactions[0].Category)

	assert.Equal(t, 2, len(actualData.Accounts))
	assert.Equal(t, "Test Account", actualData.Accounts[0].Name)
	assert.Equal(t, "Wallet", actualData.Accounts[1].Name)
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

	assert.Equal(t, 0, len(actualData.BankAccountTransactions))
	assert.Equal(t, 0, len(actualData.CashAccountTransactions))
	assert.Equal(t, 0, len(actualData.MemorizedTransactions))
	assert.Equal(t, 0, len(actualData.InvestmentAccountTransactions))
	assert.Equal(t, 0, len(actualData.Accounts))
	assert.Equal(t, 0, len(actualData.Categories))
	assert.Equal(t, 0, len(actualData.Classes))
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

	assert.Equal(t, 1, len(actualData.BankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.BankAccountTransactions[0].Date)
	assert.Equal(t, "-123.45", actualData.BankAccountTransactions[0].Amount)
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

	assert.Equal(t, 2, len(actualData.BankAccountTransactions))
	assert.Equal(t, "2024/10/9", actualData.BankAccountTransactions[0].Date)
	assert.Equal(t, "-123.45", actualData.BankAccountTransactions[0].Amount)
	assert.Equal(t, "2024/10/11", actualData.BankAccountTransactions[1].Date)
	assert.Equal(t, "100.00", actualData.BankAccountTransactions[1].Amount)
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
	assert.Equal(t, "2024/10/12", actualData.Date)
	assert.Equal(t, "-123.45", actualData.Amount)
	assert.Equal(t, qifClearedStatusUnreconciled, actualData.ClearedStatus)
	assert.Equal(t, "100", actualData.Num)
	assert.Equal(t, "Foo", actualData.Payee)
	assert.Equal(t, "Bar", actualData.Memo)
	assert.Equal(t, 3, len(actualData.Addresses))
	assert.Equal(t, "Address 1", actualData.Addresses[0])
	assert.Equal(t, "Address 2", actualData.Addresses[1])
	assert.Equal(t, "Address 3", actualData.Addresses[2])
	assert.Equal(t, "Test Category", actualData.Category)
	assert.Equal(t, 2, len(actualData.SubTransactionCategory))
	assert.Equal(t, "Part1 Category", actualData.SubTransactionCategory[0])
	assert.Equal(t, "Part2 Category", actualData.SubTransactionCategory[1])
	assert.Equal(t, 2, len(actualData.SubTransactionMemo))
	assert.Equal(t, "Part1 Memo", actualData.SubTransactionMemo[0])
	assert.Equal(t, "Part2 Memo", actualData.SubTransactionMemo[1])
	assert.Equal(t, 2, len(actualData.SubTransactionAmount))
	assert.Equal(t, "-100.00", actualData.SubTransactionAmount[0])
	assert.Equal(t, "-23.45", actualData.SubTransactionAmount[1])
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
	assert.Equal(t, qifCheckTransactionType, actualData.TransactionType)
	assert.Equal(t, "2024/10/12", actualData.Date)
	assert.Equal(t, "-123.45", actualData.Amount)
	assert.Equal(t, qifClearedStatusCleared, actualData.ClearedStatus)
	assert.Equal(t, "100", actualData.Num)
	assert.Equal(t, "Foo", actualData.Payee)
	assert.Equal(t, "Bar", actualData.Memo)
	assert.Equal(t, "2024/10/13", actualData.Amortization.FirstPaymentDate)
	assert.Equal(t, "3", actualData.Amortization.TotalYearsForLoan)
	assert.Equal(t, "1", actualData.Amortization.NumberOfPayments)
	assert.Equal(t, "2", actualData.Amortization.NumberOfPeriodsPerYear)
	assert.Equal(t, "12.34", actualData.Amortization.InterestRate)
	assert.Equal(t, "100.45", actualData.Amortization.CurrentLoanBalance)
	assert.Equal(t, "234.56", actualData.Amortization.OriginalLoanAmount)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KD"})
	assert.Nil(t, err)
	assert.Equal(t, qifDepositTransactionType, actualData.TransactionType)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KP"})
	assert.Nil(t, err)
	assert.Equal(t, qifPaymentTransactionType, actualData.TransactionType)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KI"})
	assert.Nil(t, err)
	assert.Equal(t, qifInvestmentTransactionType, actualData.TransactionType)

	actualData, err = reader.parseMemorizedTransaction(context, []string{"KE"})
	assert.Nil(t, err)
	assert.Equal(t, qifElectronicPayeeTransactionType, actualData.TransactionType)
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
	assert.Equal(t, "2024/10/12", actualData.Date)
	assert.Equal(t, "Buy", actualData.Action)
	assert.Equal(t, "Test", actualData.Security)
	assert.Equal(t, "12.34", actualData.Price)
	assert.Equal(t, "10", actualData.Quantity)
	assert.Equal(t, "-123.4", actualData.Amount)
	assert.Equal(t, qifClearedStatusReconciled, actualData.ClearedStatus)
	assert.Equal(t, "Foo", actualData.Text)
	assert.Equal(t, "Bar", actualData.Memo)
	assert.Equal(t, "Test2", actualData.Commission)
	assert.Equal(t, "Account Name", actualData.AccountForTransfer)
	assert.Equal(t, "100", actualData.AmountTransferred)
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
	assert.Equal(t, "Account Name", actualData.Name)
	assert.Equal(t, "Account Type", actualData.AccountType)
	assert.Equal(t, "Some Text", actualData.Description)
	assert.Equal(t, "1234.56", actualData.CreditLimit)
	assert.Equal(t, "2024/10/12", actualData.StatementBalanceDate)
	assert.Equal(t, "123.45", actualData.StatementBalanceAmount)
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
	assert.Equal(t, "Category Name:Sub Category Name", actualData.Name)
	assert.Equal(t, "Some Text", actualData.Description)
	assert.Equal(t, true, actualData.TaxRelated)
	assert.Equal(t, qifIncomeTransaction, actualData.CategoryType)
	assert.Equal(t, "123.45", actualData.BudgetAmount)
	assert.Equal(t, "Test", actualData.TaxScheduleInformation)

	actualData2, err := reader.parseCategory(context, []string{
		"NCategory Name:Sub Category Name",
		"DSome Text",
		"E",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Category Name:Sub Category Name", actualData2.Name)
	assert.Equal(t, "Some Text", actualData2.Description)
	assert.Equal(t, false, actualData2.TaxRelated)
	assert.Equal(t, qifExpenseTransaction, actualData2.CategoryType)

	actualData3, err := reader.parseCategory(context, []string{
		"NCategory Name:Sub Category Name",
		"DSome Text",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Category Name:Sub Category Name", actualData3.Name)
	assert.Equal(t, "Some Text", actualData3.Description)
	assert.Equal(t, qifExpenseTransaction, actualData3.CategoryType)
}

func TestQifDataReaderParseClass_SupportedFields(t *testing.T) {
	reader := &qifDataReader{}
	context := core.NewNullContext()

	actualData, err := reader.parseClass(context, []string{
		"NClass Name",
		"DSome Text",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Class Name", actualData.Name)
	assert.Equal(t, "Some Text", actualData.Description)
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
	assert.Equal(t, qifClearedStatusUnreconciled, actualTransactionData.ClearedStatus)

	actualMemorizedTransactionData, err := reader.parseMemorizedTransaction(context, []string{
		"ZTest",
		"KZ",
		"",
	})
	assert.Nil(t, err)
	assert.Equal(t, qifInvalidTransactionType, actualMemorizedTransactionData.TransactionType)

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
