package alipay

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestAlipayCsvFileImporterParseImportedData_MinimumValidData(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 01:23:45 ,xxxx            ,0.12   ,收入      ,交易成功    ,\n" +
		"2024-09-01 12:34:56 ,xxxx            ,123.45  ,支出      ,交易成功    ,\n" +
		"2024-09-01 23:59:59 ,充值-普通充值             ,0.05   ,不计收支    ,交易成功    ,\n" +
		"2024-09-02 23:59:59 ,提现-普通提现             ,0.03   ,不计收支    ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, err := converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 4, len(allNewTransactions))
	assert.Equal(t, 2, len(allNewAccounts))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))
	assert.Equal(t, 0, len(allNewTags))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(12), allNewTransactions[0].Amount)
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[1].Type)
	assert.Equal(t, "2024-09-01 12:34:56", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[1].TransactionTime), time.UTC))
	assert.Equal(t, int64(12345), allNewTransactions[1].Amount)
	assert.Equal(t, "", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[2].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[2].Type)
	assert.Equal(t, "2024-09-01 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[2].TransactionTime), time.UTC))
	assert.Equal(t, int64(5), allNewTransactions[2].Amount)
	assert.Equal(t, "", allNewTransactions[2].OriginalSourceAccountName)
	assert.Equal(t, "Alipay", allNewTransactions[2].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[2].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewTransactions[3].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, allNewTransactions[3].Type)
	assert.Equal(t, "2024-09-02 23:59:59", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[3].TransactionTime), time.UTC))
	assert.Equal(t, int64(3), allNewTransactions[3].Amount)
	assert.Equal(t, "Alipay", allNewTransactions[3].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalDestinationAccountName)
	assert.Equal(t, "", allNewTransactions[3].OriginalCategoryName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Alipay", allNewAccounts[0].Name)
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

func TestAlipayCsvFileImporterParseImportedData_ParseRefundTransaction(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 01:23:45 ,0.12   ,不计收支    ,退款成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12), allNewTransactions[0].Amount)
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)

	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 01:23:45 ,0.12   ,收入      ,退税成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, allNewTransactions[0].Type)
	assert.Equal(t, "2024-09-01 01:23:45", utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(allNewTransactions[0].TransactionTime), time.UTC))
	assert.Equal(t, int64(-12), allNewTransactions[0].Amount)
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalCategoryName)
}

func TestAlipayCsvFileImporterParseImportedData_ParseInvalidTime(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01T12:34:56 ,0.12   ,收入      ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)

	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,收/支     ,交易状态    ,\n" +
		"09/01/2024 12:34:56 ,0.12   ,收入      ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrTransactionTimeInvalid.Message)
}

func TestAlipayCsvFileImporterParseImportedData_ParseInvalidType(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,0.12   ,        ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestAlipayCsvFileImporterParseImportedData_ParseAccountName(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	// income to alipay wallet
	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,收入      ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalSourceAccountName)

	// refund to other account
	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,不计收支    ,退款成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)

	// transfer to alipay wallet
	data3, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,test                ,充值-普通充值             ,0.12   ,不计收支    ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data3), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalDestinationAccountName)

	// transfer from alipay wallet
	data4, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,test                ,提现-实时提现             ,0.12   ,不计收支    ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data4), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "Alipay", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)

	// transfer in
	data5, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,test                ,xx-转入             ,0.12   ,不计收支    ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data5), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)

	// transfer out
	data6, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,test                ,xx-转出             ,0.12   ,不计收支    ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data6), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)

	// repayment
	data7, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,交易对方            ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,test                ,xx还款             ,0.12   ,不计收支    ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data7), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "test", allNewTransactions[0].OriginalDestinationAccountName)
}

func TestAlipayCsvFileImporterParseImportedData_ParseCategory(t *testing.T) {
	converter := AlipayAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("------------------------------------------------------------------------------------\n" +
		"导出信息：\n" +
		"姓名：xxx\n" +
		"支付宝账户：xxx@xxx.xxx\n" +
		"起始时间：[2024-01-01 00:00:00]    终止时间：[2024-09-01 23:59:59]\n" +
		"导出交易类型：[全部]\n" +
		"------------------------支付宝（中国）网络技术有限公司  电子客户回单------------------------\n" +
		"交易时间,交易分类,商品说明,收/支,金额,交易状态,\n" +
		"2024-09-01 01:23:45,Test Category,xxxx,收入,0.12,交易成功,\n" +
		"2024-09-01 12:34:56,Test Category2,xxxx,支出,123.45,交易成功,\n" +
		"2024-09-01 23:59:59,Test Category3,充值-普通充值,不计收支,0.05,交易成功,\n")
	assert.Nil(t, err)

	allNewTransactions, _, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(allNewTransactions))
	assert.Equal(t, 1, len(allNewSubExpenseCategories))
	assert.Equal(t, 1, len(allNewSubIncomeCategories))
	assert.Equal(t, 1, len(allNewSubTransferCategories))

	assert.Equal(t, int64(1234567890), allNewSubExpenseCategories[0].Uid)
	assert.Equal(t, "Test Category2", allNewSubExpenseCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubIncomeCategories[0].Uid)
	assert.Equal(t, "Test Category", allNewSubIncomeCategories[0].Name)

	assert.Equal(t, int64(1234567890), allNewSubTransferCategories[0].Uid)
	assert.Equal(t, "Test Category3", allNewSubTransferCategories[0].Name)
}

func TestAlipayCsvFileImporterParseImportedData_ParseRelatedAccount(t *testing.T) {
	converter := AlipayAppTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("------------------------------------------------------------------------------------\n" +
		"导出信息：\n" +
		"姓名：xxx\n" +
		"支付宝账户：xxx@xxx.xxx\n" +
		"起始时间：[2024-01-01 00:00:00]    终止时间：[2024-09-01 23:59:59]\n" +
		"导出交易类型：[全部]\n" +
		"------------------------支付宝（中国）网络技术有限公司  电子客户回单------------------------\n" +
		"交易时间,商品说明,收/支,金额,收/付款方式,交易状态,\n" +
		"2024-09-01 03:45:07,余额宝-单次转入,不计收支,0.01,Test Account,交易成功,\n" +
		"2024-09-01 05:07:29,信用卡还款,不计收支,0.02,Test Account2,交易成功,\n")
	assert.Nil(t, err)

	allNewTransactions, allNewAccounts, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(allNewTransactions))
	assert.Equal(t, 3, len(allNewAccounts))

	assert.Equal(t, int64(1234567890), allNewTransactions[0].Uid)
	assert.Equal(t, int64(1), allNewTransactions[0].Amount)
	assert.Equal(t, "Test Account", allNewTransactions[0].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[0].OriginalDestinationAccountName)

	assert.Equal(t, int64(1234567890), allNewTransactions[1].Uid)
	assert.Equal(t, int64(2), allNewTransactions[1].Amount)
	assert.Equal(t, "Test Account2", allNewTransactions[1].OriginalSourceAccountName)
	assert.Equal(t, "", allNewTransactions[1].OriginalDestinationAccountName)

	assert.Equal(t, int64(1234567890), allNewAccounts[0].Uid)
	assert.Equal(t, "Test Account", allNewAccounts[0].Name)
	assert.Equal(t, "CNY", allNewAccounts[0].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[1].Uid)
	assert.Equal(t, "", allNewAccounts[1].Name)
	assert.Equal(t, "CNY", allNewAccounts[1].Currency)

	assert.Equal(t, int64(1234567890), allNewAccounts[2].Uid)
	assert.Equal(t, "Test Account2", allNewAccounts[2].Name)
	assert.Equal(t, "CNY", allNewAccounts[2].Currency)
}

func TestAlipayCsvFileImporterParseImportedData_ParseDescription(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,收/支     ,交易状态    ,备注                  ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,收入      ,交易成功    ,test2               ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err := converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "test2", allNewTransactions[0].Comment)

	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,收/支     ,交易状态    ,备注                  ,\n" +
		"2024-09-01 12:34:56 ,test                ,0.12   ,收入      ,交易成功    ,                    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	allNewTransactions, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(allNewTransactions))
	assert.Equal(t, "test", allNewTransactions[0].Comment)
}

func TestAlipayCsvFileImporterParseImportedData_SkipClosedIncomeOrTransferTransaction(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 01:23:45 ,xxxx            ,0.12   ,收入      ,交易关闭    ,\n" +
		"2024-09-01 23:59:59 ,充值-普通充值             ,0.05   ,不计收支    ,交易关闭    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestAlipayCsvFileImporterParseImportedData_SkipUnknownProductTransferTransaction(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 23:59:59 ,xxxx                ,0.05   ,不计收支    ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestAlipayCsvFileImporterParseImportedData_SkipUnknownStatusTransaction(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1234567890,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,商品名称                ,金额（元）,收/支     ,交易状态    ,\n" +
		"2024-09-01 01:23:45 ,xxxx            ,0.12   ,收入      ,xxxx    ,\n" +
		"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}

func TestAlipayCsvFileImporterParseImportedData_MissingFileHeader(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	data, err := simplifiedchinese.GB18030.NewEncoder().String(
		"交易创建时间              ,金额（元）,收/支     ,交易状态    ,\n" +
			"2024-09-01 12:34:56 ,0.12   ,收入      ,交易成功    ,\n" +
			"------------------------------------------------------------------------------------\n")
	assert.Nil(t, err)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)

	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(""), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrInvalidFileHeader.Message)
}

func TestAlipayCsvFileImporterParseImportedData_MissingRequiredColumn(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	// Missing Time Column
	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"金额（元）,收/支     ,交易状态    ,\n" +
		"0.12   ,收入      ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Amount Column
	data2, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,收/支     ,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,收入      ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data2), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Status Column
	data3, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,收/支     ,\n" +
		"2024-09-01 12:34:56 ,0.12   ,收入      ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data3), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)

	// Missing Type Column
	data4, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,交易状态    ,\n" +
		"2024-09-01 12:34:56 ,0.12   ,交易成功    ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data4), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrMissingRequiredFieldInHeaderRow.Message)
}

func TestAlipayCsvFileImporterParseImportedData_NoTransactionData(t *testing.T) {
	converter := AlipayWebTransactionDataCsvFileImporter
	context := core.NewNullContext()

	user := &models.User{
		Uid:             1,
		DefaultCurrency: "CNY",
	}

	data1, err := simplifiedchinese.GB18030.NewEncoder().String("支付宝交易记录明细查询\n" +
		"账号:[xxx@xxx.xxx]\n" +
		"起始日期:[2024-01-01 00:00:00]    终止日期:[2024-09-01 23:59:59]\n" +
		"---------------------------------交易记录明细列表------------------------------------\n" +
		"交易创建时间              ,金额（元）,收/支     ,交易状态    ,\n" +
		"------------------------------------------------------------------------------------\n")
	_, _, _, _, _, _, err = converter.ParseImportedData(context, user, []byte(data1), 0, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, errs.ErrNotFoundTransactionDataInFile.Message)
}
