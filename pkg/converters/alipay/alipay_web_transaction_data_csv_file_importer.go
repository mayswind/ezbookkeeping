package alipay

// alipayWebTransactionDataCsvImporter defines the structure of alipay (web) csv importer for transaction data
type alipayWebTransactionDataCsvImporter struct {
	alipayTransactionDataCsvImporter
}

// Initialize a alipay (web) transaction data csv file importer singleton instance
var (
	AlipayWebTransactionDataCsvImporter = &alipayWebTransactionDataCsvImporter{
		alipayTransactionDataCsvImporter{
			fileHeaderLine:           "支付宝交易记录明细查询",
			dataHeaderStartContent:   "交易记录明细列表",
			dataBottomEndLineRune:    '-',
			timeColumnName:           "交易创建时间",
			categoryColumnName:       "",
			targetNameColumnName:     "交易对方",
			productNameColumnName:    "商品名称",
			amountColumnName:         "金额（元）",
			typeColumnName:           "收/支",
			relatedAccountColumnName: "",
			statusColumnName:         "交易状态",
			descriptionColumnName:    "备注",
		},
	}
)
