package iif

// iifAccountDataset defines the structure of intuit interchange format (iif) account dataset
type iifAccountDataset struct {
	accountDataColumnIndexes map[string]int
	accounts                 []*iifAccountData
}

// iifAccountData defines the structure of intuit interchange format (iif) account data
type iifAccountData struct {
	dataItems []string
}

// iifTransactionDataset defines the structure of intuit interchange format (iif) transaction dataset
type iifTransactionDataset struct {
	transactionDataColumnIndexes map[string]int
	splitDataColumnIndexes       map[string]int
	transactions                 []*iifTransactionData
}

// iifTransactionData defines the structure of intuit interchange format (iif) transaction data
type iifTransactionData struct {
	dataItems []string
	splitData []*iifTransactionSplitData
}

// iifTransactionSplitData defines the structure of intuit interchange format (iif) transaction split data
type iifTransactionSplitData struct {
	dataItems []string
}

func (s *iifTransactionDataset) getTransactionDataItemValue(transactionData *iifTransactionData, columnName string) (string, bool) {
	if transactionData == nil {
		return "", false
	}

	index, exists := s.transactionDataColumnIndexes[columnName]

	if !exists || index < 0 || index >= len(transactionData.dataItems) {
		return "", false
	}

	return transactionData.dataItems[index], true
}

func (s *iifTransactionDataset) getSplitDataItemValue(splitData *iifTransactionSplitData, columnName string) (string, bool) {
	if splitData == nil {
		return "", false
	}

	index, exists := s.splitDataColumnIndexes[columnName]

	if !exists || index < 0 || index >= len(splitData.dataItems) {
		return "", false
	}

	return splitData.dataItems[index], true
}
