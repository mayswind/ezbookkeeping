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
	TransactionDataColumnIndexes map[string]int
	SplitDataColumnIndexes       map[string]int
	Transactions                 []*iifTransactionData
}

// iifTransactionData defines the structure of intuit interchange format (iif) transaction data
type iifTransactionData struct {
	DataItems []string
	SplitData []*iifTransactionSplitData
}

// iifTransactionSplitData defines the structure of intuit interchange format (iif) transaction split data
type iifTransactionSplitData struct {
	DataItems []string
}

func (s *iifTransactionDataset) getTransactionDataItemValue(transactionData *iifTransactionData, columnName string) (string, bool) {
	if transactionData == nil {
		return "", false
	}

	index, exists := s.TransactionDataColumnIndexes[columnName]

	if !exists || index < 0 || index >= len(transactionData.DataItems) {
		return "", false
	}

	return transactionData.DataItems[index], true
}

func (s *iifTransactionDataset) getSplitDataItemValue(splitData *iifTransactionSplitData, columnName string) (string, bool) {
	if splitData == nil {
		return "", false
	}

	index, exists := s.SplitDataColumnIndexes[columnName]

	if !exists || index < 0 || index >= len(splitData.DataItems) {
		return "", false
	}

	return splitData.DataItems[index], true
}
