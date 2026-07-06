package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// mergedTransactionDataTable defines the structure of merged transaction data table which merges multiple transaction data tables
type mergedTransactionDataTable struct {
	dataTables []TransactionDataTable
}

// mergedTransactionDataRowIterator defines the data row iterator structure of merged transaction data table
type mergedTransactionDataRowIterator struct {
	iterators    []TransactionDataRowIterator
	currentIndex int
}

// HasColumn returns whether any merged data table has specified column
func (t *mergedTransactionDataTable) HasColumn(column TransactionDataTableColumn) bool {
	for _, dt := range t.dataTables {
		if dt.HasColumn(column) {
			return true
		}
	}

	return false
}

// TransactionRowCount returns the total count of transaction data row across all merged data tables
func (t *mergedTransactionDataTable) TransactionRowCount() int {
	totalRowCount := 0

	for _, dt := range t.dataTables {
		totalRowCount += dt.TransactionRowCount()
	}

	return totalRowCount
}

// TransactionRowIterator returns the iterator of transaction data row which iterates through all merged data tables
func (t *mergedTransactionDataTable) TransactionRowIterator() TransactionDataRowIterator {
	iterators := make([]TransactionDataRowIterator, len(t.dataTables))

	for i, dt := range t.dataTables {
		iterators[i] = dt.TransactionRowIterator()
	}

	return &mergedTransactionDataRowIterator{
		iterators:    iterators,
		currentIndex: 0,
	}
}

// HasNext returns whether the iterator does not reach the end
func (it *mergedTransactionDataRowIterator) HasNext() bool {
	if it.currentIndex >= len(it.iterators) {
		return false
	}

	if it.iterators[it.currentIndex].HasNext() {
		return true
	}

	for i := it.currentIndex + 1; i < len(it.iterators); i++ {
		if it.iterators[i].HasNext() {
			return true
		}
	}

	return false
}

// Next returns the next transaction data row
func (it *mergedTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (dataRow TransactionDataRow, err error) {
	if it.currentIndex >= len(it.iterators) {
		return nil, nil
	}

	if it.iterators[it.currentIndex].HasNext() {
		return it.iterators[it.currentIndex].Next(ctx, user)
	}

	for it.currentIndex < len(it.iterators) {
		it.currentIndex++

		if it.currentIndex >= len(it.iterators) {
			break
		}

		if it.iterators[it.currentIndex].HasNext() {
			return it.iterators[it.currentIndex].Next(ctx, user)
		}
	}

	return nil, nil
}

// CreateNewMergedTransactionDataTable returns a merged transaction data table from multiple transaction data tables
func CreateNewMergedTransactionDataTable(dataTables []TransactionDataTable) TransactionDataTable {
	return &mergedTransactionDataTable{
		dataTables: dataTables,
	}
}
