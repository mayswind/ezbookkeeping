package feidee

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

func createNewFeideeMymoneyAppTransactionBasicDataTable(ctx core.Context, originalDataTable datatable.BasicDataTable) (datatable.BasicDataTable, error) {
	iterator := originalDataTable.DataRowIterator()
	allOriginalLines := make([][]string, 0)
	hasFileHeader := false

	for iterator.HasNext() {
		row := iterator.Next()

		if !hasFileHeader {
			if row.ColumnCount() <= 0 {
				continue
			} else if strings.Index(row.GetData(0), feideeMymoneyAppTransactionDataCsvFileHeader) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[feidee_mymoney_app_transaction_data_extrator.createNewFeideeMymoneyAppTransactionBasicDataTable] read unexpected line in row \"%s\" before read file header", iterator.CurrentRowId())
				continue
			}
		}

		items := make([]string, row.ColumnCount())

		for i := 0; i < row.ColumnCount(); i++ {
			items[i] = strings.Trim(row.GetData(i), " ")
		}

		allOriginalLines = append(allOriginalLines, items)
	}

	if !hasFileHeader {
		return nil, errs.ErrInvalidFileHeader
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_extrator.createNewFeideeMymoneyAppTransactionBasicDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return csv.CreateNewCustomCsvBasicDataTable(allOriginalLines, true), nil
}
