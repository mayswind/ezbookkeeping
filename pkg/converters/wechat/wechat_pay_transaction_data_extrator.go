package wechat

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

func createNewWeChatPayTransactionBasicDataTable(ctx core.Context, originalDataTable datatable.BasicDataTable) (datatable.BasicDataTable, error) {
	iterator := originalDataTable.DataRowIterator()
	hasFileHeader := false
	foundContentBeforeDataHeaderLine := false
	headerColumnCount := 0
	fromIndex := 0
	toIndex := originalDataTable.DataRowCount()
	currentIndex := -1

	for iterator.HasNext() {
		row := iterator.Next()
		currentIndex++

		if !hasFileHeader {
			if row.ColumnCount() <= 0 {
				continue
			} else if strings.Index(row.GetData(0), wechatPayTransactionDataCsvFileHeader) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[wechat_pay_transaction_data_extrator.createNewWeChatPayTransactionBasicDataTable] read unexpected line in \"%s\" before read file header", iterator.CurrentRowId())
				continue
			}
		}

		if !foundContentBeforeDataHeaderLine {
			if row.ColumnCount() <= 0 {
				continue
			} else if strings.Index(row.GetData(0), wechatPayTransactionDataHeaderStartContentBeginning) == 0 {
				foundContentBeforeDataHeaderLine = true
				continue
			} else {
				continue
			}
		}

		if foundContentBeforeDataHeaderLine {
			if row.ColumnCount() <= 0 {
				continue
			}

			if headerColumnCount == 0 {
				headerColumnCount = row.ColumnCount()
				fromIndex = currentIndex
			} else if headerColumnCount > 0 && row.ColumnCount() < headerColumnCount {
				log.Errorf(ctx, "[wechat_pay_transaction_data_extrator.createNewWeChatPayTransactionBasicDataTable] cannot parse \"%s\", because may missing some columns (column count %d in data row is less than header column count %d)", iterator.CurrentRowId(), row.ColumnCount(), headerColumnCount)
				return nil, errs.NewImportError(errs.ErrFewerFieldsInDataRowThanInHeaderRow, iterator.CurrentRowId())
			}
		}
	}

	if !hasFileHeader || !foundContentBeforeDataHeaderLine {
		return nil, errs.ErrInvalidFileHeader
	}

	if toIndex-fromIndex < 2 {
		log.Errorf(ctx, "[wechat_pay_transaction_data_extrator.createNewWeChatPayTransactionBasicDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return datatable.CreateSubBasicTable(originalDataTable, fromIndex, toIndex, true), nil
}
