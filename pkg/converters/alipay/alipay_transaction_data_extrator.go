package alipay

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func createNewAlipayTransactionBasicDataTable(ctx core.Context, originalDataTable datatable.BasicDataTable, fileHeaderLine string, dataHeaderStartContent []string, dataBottomEndLineRune rune) (datatable.BasicDataTable, error) {
	iterator := originalDataTable.DataRowIterator()
	allOriginalLines := make([][]string, 0)
	hasFileHeader := false
	foundContentBeforeDataHeaderLine := false

	for iterator.HasNext() {
		row := iterator.Next()

		if !hasFileHeader {
			if row.ColumnCount() <= 0 {
				continue
			} else if strings.Index(row.GetData(0), fileHeaderLine) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[alipay_transaction_data_extrator.createNewAlipayTransactionBasicDataTable] read unexpected line in row \"%s\" before read file header", iterator.CurrentRowId())
				continue
			}
		}

		if !foundContentBeforeDataHeaderLine {
			if row.ColumnCount() <= 0 {
				continue
			} else if utils.ContainsAnyString(row.GetData(0), dataHeaderStartContent) {
				foundContentBeforeDataHeaderLine = true
				continue
			} else {
				continue
			}
		}

		if foundContentBeforeDataHeaderLine {
			if row.ColumnCount() <= 0 {
				continue
			} else if row.ColumnCount() == 1 && dataBottomEndLineRune > 0 && utils.ContainsOnlyOneRune(row.GetData(0), dataBottomEndLineRune) {
				break
			}

			items := make([]string, row.ColumnCount())

			for i := 0; i < row.ColumnCount(); i++ {
				items[i] = strings.Trim(row.GetData(i), " ")
			}

			if len(allOriginalLines) > 0 && len(items) < len(allOriginalLines[0]) {
				log.Errorf(ctx, "[alipay_transaction_data_extrator.createNewAlipayTransactionBasicDataTable] cannot parse row \"%s\", because may missing some columns (column count %d in data row is less than header column count %d)", iterator.CurrentRowId(), len(items), len(allOriginalLines[0]))
				return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
			}

			allOriginalLines = append(allOriginalLines, items)
		}
	}

	if !hasFileHeader || !foundContentBeforeDataHeaderLine {
		return nil, errs.ErrInvalidFileHeader
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[alipay_transaction_data_extrator.createNewAlipayTransactionBasicDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	return csv.CreateNewCustomCsvBasicDataTable(allOriginalLines, true), nil
}
