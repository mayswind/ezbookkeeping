package feidee

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

func createNewFeideeMymoneyAppTransactionBasicDataTable(ctx core.Context, originalDataTable datatable.BasicDataTable) (datatable.BasicDataTable, error) {
	iterator := originalDataTable.DataRowIterator()
	hasFileHeader := false
	fromIndex := 0
	toIndex := originalDataTable.DataRowCount()
	currentIndex := -1

	for iterator.HasNext() {
		row := iterator.Next()
		currentIndex++

		if !hasFileHeader {
			if row.ColumnCount() <= 0 {
				continue
			} else if strings.Index(row.GetData(0), feideeMymoneyAppTransactionDataCsvFileHeader) == 0 {
				hasFileHeader = true
				fromIndex = currentIndex + 1
				continue
			} else {
				log.Warnf(ctx, "[feidee_mymoney_app_transaction_data_extrator.createNewFeideeMymoneyAppTransactionBasicDataTable] read unexpected line in \"%s\" before read file header", iterator.CurrentRowId())
				continue
			}
		}
	}

	if !hasFileHeader {
		return nil, errs.ErrInvalidFileHeader
	}

	if toIndex-fromIndex < 2 {
		log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_extrator.createNewFeideeMymoneyAppTransactionBasicDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return datatable.CreateSubBasicTable(originalDataTable, fromIndex, toIndex, true), nil
}
