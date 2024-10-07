package alipay

// alipayAppTransactionDataCsvImporter defines the structure of alipay app csv importer for transaction data
type alipayAppTransactionDataCsvImporter struct {
	alipayTransactionDataCsvImporter
}

// Initialize a alipay app transaction data csv file importer singleton instance
var (
	AlipayAppTransactionDataCsvImporter = &alipayAppTransactionDataCsvImporter{
		alipayTransactionDataCsvImporter{
			fileHeaderLine:           "------------------------------------------------------------------------------------",
			dataHeaderStartContent:   "支付宝（中国）网络技术有限公司  电子客户回单",
			timeColumnName:           "交易时间",
			categoryColumnName:       "交易分类",
			targetNameColumnName:     "交易对方",
			productNameColumnName:    "商品说明",
			amountColumnName:         "金额",
			typeColumnName:           "收/支",
			relatedAccountColumnName: "收/付款方式",
			statusColumnName:         "交易状态",
			descriptionColumnName:    "备注",
		},
	}
)
