package wechat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestWeChatPayCsvFileImporterParseImportedData_MinimumValidData(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),当前状态\n" +
		"2024-09-01 01:23:45,二维码收款,收入,￥0.12,已收钱\n" +
		"2024-09-01 12:34:56,商户消费,支出,￥123.45,支付成功\n" +
		"2024-09-01 23:59:59,零钱充值,/,￥0.05,充值完成\n" +
		"2024-09-02 23:59:59,零钱提现,/,￥0.03,提现已到账\n"
	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 2, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[0].Amount)
	assert.Equal(t, "Wallet", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "二维码收款", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, "2024-09-01 12:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[1].Amount)
	assert.Equal(t, "", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "商户消费", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[2].Type)
	assert.Equal(t, "2024-09-01 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(5), allNewTransactions[2].Amount)
	assert.Equal(t, "", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Wallet", allNewTransactions[2].OriginalDestinationAccountName)
	assert.Equal(t, "零钱充值", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, "2024-09-02 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(3), allNewTransactions[3].Amount)
	assert.Equal(t, "Wallet", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "零钱提现", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Wallet", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "商户消费", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "二维码收款", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "零钱充值", allNewSubTransferCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[1].Uid)
	assert.Equal(t, "零钱提现", allNewSubTransferCategories[1].Name)
}

func TestWeChatPayCsvFileImporterParseImportedData_ParseRefundTransaction(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),当前状态\n" +
		"2024-09-01 01:23:45,xxx-退款,收入,￥0.12,已全额退款\n"
	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12), allNewTransactions[0].Amount)
	assert.Equal(t, "Wallet", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "xxx-退款", allNewTransactions[0].OriginalCategoryName)
}

func TestWeChatPayCsvFileImporterParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),当前状态\n" +
		"2024-09-01T01:23:45,二维码收款,收入,￥0.12,已收钱\n"
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	data2 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),当前状态\n" +
		"09/01/2024 12:34:56,二维码收款,收入,￥0.12,已收钱\n"
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestWeChatPayCsvFileImporterParseImportedData_ParseInvalidType(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),当前状态\n" +
		"2024-09-01T01:23:45,xxx,,￥0.12,支付成功\n"
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestWeChatPayCsvFileImporterParseImportedData_ParseAccountName(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// income to wechat wallet without related account name
	data1 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),支付方式,当前状态\n" +
		"2024-09-01 01:23:45,二维码收款,收入,￥0.12,/,已收钱\n"
	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Wallet", allNewTransactions[0].OriginalSourceAccountName)

	// refund to other account
	data2 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),支付方式,当前状态\n" +
		"2024-09-01 01:23:45,xxx-退款,收入,￥0.12,test,已全额退款\n"
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "test", allNewTransactions[0].OriginalSourceAccountName)

	// transfer to wechat wallet
	data3 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),支付方式,当前状态\n" +
		"2024-09-01 23:59:59,零钱充值,/,￥0.05,test,充值完成\n"
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data3), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "test", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Wallet", allNewTransactions[0].OriginalDestinationAccountName)

	// transfer from wechat wallet
	data4 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),支付方式,当前状态\n" +
		"2024-09-02 23:59:59,零钱提现,/,￥0.03,test,提现已到账\n"
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data4), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Wallet", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)
}

func TestWeChatPayCsvFileImporterParseImportedData_ParseDescription(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),当前状态,备注\n" +
		"2024-09-01 01:23:45,二维码收款,收入,￥0.12,已收钱,\"/\"\n"
	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].Comment)

	data2 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元),当前状态,备注\n" +
		"2024-09-01 01:23:45,二维码收款,收入,￥0.12,已收钱,\"foo\"\"bar,\ntest\"\n"
	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "foo\"bar,\ntest", allNewTransactions[0].Comment)
}

func TestWeChatPayCsvFileImporterParseImportedData_MissingFileHeader(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	data := "交易时间,交易类型,收/支,金额(元),当前状态\n" +
		"2024-09-01 01:23:45,二维码收款,收入,￥0.12,已收钱\n"
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)
}

func TestWeChatPayCsvFileImporterParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := WeChatPayTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Time Column
	data1 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易类型,收/支,金额(元),当前状态\n" +
		"二维码收款,收入,￥0.12,已收钱\n"
	_, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Category Column
	data2 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,收/支,金额(元),当前状态\n" +
		"2024-09-01 01:23:45,收入,￥0.12,已收钱\n"
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Type Column
	data3 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,金额(元),当前状态\n" +
		"2024-09-01 01:23:45,二维码收款,￥0.12,已收钱\n"
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data3), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	data4 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,当前状态\n" +
		"2024-09-01 01:23:45,二维码收款,收入,已收钱\n"
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data4), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Status Column
	data5 := "微信支付账单明细,,,,\n" +
		"微信昵称：[xxx],,,,\n" +
		"起始时间：[2024-01-01 00:00:00] 终止时间：[2024-09-01 23:59:59],,,,\n" +
		",,,,\n" +
		"----------------------微信支付账单明细列表--------------------,,,,\n" +
		"交易时间,交易类型,收/支,金额(元)\n" +
		"2024-09-01 01:23:45,二维码收款,收入,￥0.12\n"
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data5), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}
