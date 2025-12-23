package jdcom

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestJDComFinanceCsvFileImporterParseImportedData_MinimumValidData(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支,交易分类\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,余额,交易成功,收入,其他\n" +
		"2025-09-01 12:34:56,xxx,xxx,123.45,银行卡,交易成功,支出,其他网购\n" +
		"2025-09-01 23:59:59,xxx,京东钱包余额充值,0.05,银行卡,交易成功,不计收支,余额\n" +
		"2025-09-02 23:59:59,xxx,京东余额提现,0.03,银行卡,交易成功,不计收支,余额\n"
	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := importer.ParseImportedData(context, user, []byte(data), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 3, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, "2025-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[0].Amount)
	assert.Equal(t, "余额", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "其他", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, "2025-09-01 12:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[1].Amount)
	assert.Equal(t, "银行卡", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "其他网购", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[2].Type)
	assert.Equal(t, "2025-09-01 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(5), allNewTransactions[2].Amount)
	assert.Equal(t, "银行卡", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "xxx", allNewTransactions[2].OriginalDestinationAccountName)
	assert.Equal(t, "余额", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, "2025-09-02 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(3), allNewTransactions[3].Amount)
	assert.Equal(t, "xxx", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "银行卡", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "余额", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "余额", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "银行卡", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[2].Uid)
	assert.Equal(t, "xxx", allNewAccounts[2].Name)
	assert.Equal(t, "CNY", allNewAccounts[2].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "其他网购", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "其他", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "余额", allNewSubTransferCategories[0].Name)
}

func TestJDComFinanceCsvFileImporterParseImportedData_ParseRefundTransaction(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,银行卡,退款成功,不计收支\n" +
		"2025-09-01 02:34:56,xxx,xxx,0.12(已全额退款),银行卡,交易成功,不计收支\n" +
		"2025-09-02 01:23:45,xxx,xxx,3.45,银行卡,退款成功,不计收支\n" +
		"2025-09-02 02:34:56,xxx,xxx,123.45(已退款3.45),银行卡,交易成功,支出\n"
	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data1), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2025-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12), allNewTransactions[0].Amount)
	assert.Equal(t, "银行卡", allNewTransactions[0].OriginalSourceAccountName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, "2025-09-01 02:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[1].Amount)
	assert.Equal(t, "银行卡", allNewTransactions[1].OriginalSourceAccountName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[2].Type)
	assert.Equal(t, "2025-09-02 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(-345), allNewTransactions[2].Amount)
	assert.Equal(t, "银行卡", allNewTransactions[2].OriginalSourceAccountName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[3].Type)
	assert.Equal(t, "2025-09-02 02:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[3].Amount)
	assert.Equal(t, "银行卡", allNewTransactions[3].OriginalSourceAccountName)
}

func TestJDComFinanceCsvFileImporterParseImportedData_ParseInvalidTime(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01T01:23:45,xxx,xxx,0.12,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data1), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	data2 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"09/01/2025 01:23:45,xxx,xxx,0.12,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data2), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestJDComFinanceCsvFileImporterParseImportedData_ParseInvalidType(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,银行卡,交易成功,转账\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestJDComFinanceCsvFileImporterParseImportedData_ParseInvalidAmount(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,￥0.12,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrAmountInvalid.Message)
}

func TestJDComFinanceCsvFileImporterParseImportedData_ParseAccountName(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// transfer to jd.com finance wallet
	data1 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支,交易分类\n" +
		"2025-09-01 01:23:45,xxx,京东钱包余额充值,0.05,银行卡,交易成功,不计收支,余额\n"
	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data1), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "银行卡", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "xxx", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)

	// transfer from jd.com finance wallet
	data2 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支,交易分类\n" +
		"2025-09-01 01:23:45,xxx,京东余额提现,0.05,银行卡,交易成功,不计收支,余额\n"
	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data2), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "xxx", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "银行卡", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)

	// transfer from other account
	data3 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支,交易分类\n" +
		"2025-09-01 01:23:45,xxx,京东小金库-转入,0.05,余额,交易成功,不计收支,小金库\n"
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data3), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "余额", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "xxx", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)

	// transfer to other account
	data4 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支,交易分类\n" +
		"2025-09-01 01:23:45,xxx,京东小金库-转出,0.05,余额,交易成功,不计收支,小金库\n"
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data4), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "xxx", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "余额", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)

	// refund
	data5 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支,交易分类\n" +
		"2025-09-01 01:23:45,xxx,价保退款,0.05,银行卡,交易成功,不计收支,其他\n"
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data5), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "银行卡", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)

	// repayment
	data6 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支,交易分类\n" +
		"2025-09-01 01:23:45,xxx,白条主动还款,0.05,银行卡,交易成功,不计收支,白条\n"
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data6), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "银行卡", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "xxx", allNewTransactions[0].OriginalDestinationAccountName)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[0].Type)
}

func TestJDComFinanceCsvFileImporterParseImportedData_ParseDescription(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,,0.12,银行卡,交易成功,支出\n"
	allNewTransactions, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data1), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].Comment)

	data2 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,交易说明,金额,收/付款方式,交易状态,收/支,备注\n" +
		"2025-09-01 01:23:45,xxx,xxx,Test,0.12,银行卡,交易成功,支出,\"foo\"\"bar,\ntest\"\n"
	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data2), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo\"bar,\ntest", allNewTransactions[0].Comment)

	data3 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,交易说明,金额,收/付款方式,交易状态,收/支,备注\n" +
		"2025-09-01 01:23:45,xxx,xxx,Test,0.12,银行卡,交易成功,支出,\n"
	allNewTransactions, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data3), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Test", allNewTransactions[0].Comment)
}

func TestJDComFinanceCsvFileImporterParseImportedData_SkipUnknownStatusTransaction(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,银行卡,xxxx,支出\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestJDComFinanceCsvFileImporterParseImportedData_SkipUnknownMemoTransferTransaction(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,银行卡,交易成功,不计收支\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestJDComFinanceCsvFileImporterParseImportedData_MissingFileHeader(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	data := "交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)

	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(""), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)
}

func TestJDComFinanceCsvFileImporterParseImportedData_MissingRequiredColumn(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Time Column
	data1 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"xxx,xxx,0.12,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data1), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)

	// Missing Merchant Name Column
	data2 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,交易说明,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,0.12,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data2), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Transaction Memo Column
	data3 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,金额,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,0.12,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data3), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	data4 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,收/付款方式,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,银行卡,交易成功,支出\n"
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data4), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Related Account Column
	data5 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,交易状态,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,交易成功,支出\n"
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data5), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Status Column
	data6 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,收/支\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,银行卡,支出\n"
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data6), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Type Column
	data7 := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态\n" +
		"2025-09-01 01:23:45,xxx,xxx,0.12,银行卡,交易成功\n"
	_, _, _, _, _, _, err = importer.ParseImportedData(context, user, []byte(data7), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}

func TestJDComFinanceCsvFileImporterParseImportedData_NoTransactionData(t *testing.T) {
	importer := JDComFinanceTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "导出信息：\n" +
		"京东账号名：xxxxxx\n" +
		"日期区间：2025-01-01 至 2025-09-01\n" +
		"\n" +
		"交易时间,商户名称,交易说明,金额,收/付款方式,交易状态,收/支\n"
	_, _, _, _, _, _, err := importer.ParseImportedData(context, user, []byte(data), time.UTC, converter.DefaultImporterOptions, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}
