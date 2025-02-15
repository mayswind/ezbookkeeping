package iif

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

const iifAccountSampleLineSignColumnName = "!ACCNT"
const iifTransactionSampleLineSignColumnName = "!TRNS"
const iifTransactionSplitSampleLineSignColumnName = "!SPL"
const iifTransactionEndSampleLineSignColumnName = "!ENDTRNS"

const iifAccountLineSignColumnName = "ACCNT"
const iifTransactionLineSignColumnName = "TRNS"
const iifTransactionSplitLineSignColumnName = "SPL"
const iifTransactionEndLineSignColumnName = "ENDTRNS"

// iifDataReader defines the structure of intuit interchange format (iif) data reader
type iifDataReader struct {
	reader *csv.Reader
}

// read returns the iif transaction dataset
func (r *iifDataReader) read(ctx core.Context) ([]*iifAccountDataset, []*iifTransactionDataset, error) {
	allAccountDatasets := make([]*iifAccountDataset, 0)
	allTransactionDatasets := make([]*iifTransactionDataset, 0)

	currentDatasetType := ""
	lastLineSign := ""

	var currentAccountDataset *iifAccountDataset
	var currentTransactionDataset *iifTransactionDataset
	var currentTransactionData *iifTransactionData

	for {
		items, err := r.reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[iif_data_reader.read] cannot parse tsv data, because %s", err.Error())
			return nil, nil, errs.ErrInvalidIIFFile
		}

		if len(items) == 1 && items[0] == "" {
			continue
		}

		if len(items[0]) < 1 {
			log.Errorf(ctx, "[iif_data_reader.read] line first column is empty")
			return nil, nil, errs.ErrInvalidIIFFile
		}

		// read sample line
		if items[0][0] == '!' {
			if lastLineSign != "" {
				log.Errorf(ctx, "[iif_data_reader.read] iif missing transaction end line")
				return nil, nil, errs.ErrInvalidIIFFile
			}

			if currentAccountDataset != nil {
				allAccountDatasets = append(allAccountDatasets, currentAccountDataset)
				currentAccountDataset = nil
			}

			if currentTransactionDataset != nil {
				allTransactionDatasets = append(allTransactionDatasets, currentTransactionDataset)
				currentTransactionDataset = nil
			}

			if items[0] == iifTransactionSplitSampleLineSignColumnName || items[0] == iifTransactionEndSampleLineSignColumnName {
				log.Errorf(ctx, "[iif_data_reader.read] read transaction split sample line or transaction end sample line sign before transaction sample line sign")
				return nil, nil, errs.ErrInvalidIIFFile
			} else {
				currentDatasetType = items[0]
				lastLineSign = ""
			}

			if currentDatasetType == iifAccountSampleLineSignColumnName {
				currentAccountDataset, err = r.readAccountSampleLine(ctx, items)

				if err != nil {
					return nil, nil, err
				}
			} else if currentDatasetType == iifTransactionSampleLineSignColumnName {
				currentTransactionDataset, err = r.readTransactionSampleLines(ctx, items)

				if err != nil {
					return nil, nil, err
				}
			} // not process (read sample line) for other dataset type

			continue
		}

		// read data lines
		if currentDatasetType == "" {
			log.Errorf(ctx, "[iif_data_reader.read] cannot read data line before sample line")
			return nil, nil, errs.ErrInvalidIIFFile
		} else if currentDatasetType == iifAccountSampleLineSignColumnName && currentAccountDataset != nil {
			if items[0] == iifAccountLineSignColumnName {
				accountData := &iifAccountData{
					dataItems: items,
				}
				currentAccountDataset.accounts = append(currentAccountDataset.accounts, accountData)
			} else {
				log.Errorf(ctx, "[iif_data_reader.read] iif line expected reading account sign, but actual is \"%s\"", items[0])
				return nil, nil, errs.ErrInvalidIIFFile
			}
		} else if currentDatasetType == iifTransactionSampleLineSignColumnName && currentTransactionDataset != nil {
			if lastLineSign == "" {
				if items[0] == iifTransactionLineSignColumnName {
					currentTransactionData = &iifTransactionData{
						dataItems: items,
						splitData: make([]*iifTransactionSplitData, 0),
					}
					lastLineSign = items[0]
				} else {
					log.Errorf(ctx, "[iif_data_reader.read] iif line expected reading transaction sign, but actual is \"%s\"", items[0])
					return nil, nil, errs.ErrInvalidIIFFile
				}
			} else if lastLineSign == iifTransactionLineSignColumnName || lastLineSign == iifTransactionSplitLineSignColumnName {
				if items[0] == iifTransactionSplitLineSignColumnName {
					if currentTransactionData == nil {
						log.Errorf(ctx, "[iif_data_reader.read] expected current transaction data is not nil, but read \"%s\"", items[0])
						return nil, nil, errs.ErrInvalidIIFFile
					}

					currentTransactionData.splitData = append(currentTransactionData.splitData, &iifTransactionSplitData{
						dataItems: items,
					})
					lastLineSign = items[0]
				} else if items[0] == iifTransactionEndLineSignColumnName {
					if currentTransactionData == nil {
						log.Errorf(ctx, "[iif_data_reader.read] expected current transaction data is not nil, but read \"%s\"", items[0])
						return nil, nil, errs.ErrInvalidIIFFile
					}

					if len(currentTransactionData.splitData) < 1 {
						log.Errorf(ctx, "[iif_data_reader.read] expected reading transaction split line, but read \"%s\"", items[0])
						return nil, nil, errs.ErrInvalidIIFFile
					}

					currentTransactionDataset.transactions = append(currentTransactionDataset.transactions, currentTransactionData)
					lastLineSign = ""
				} else {
					log.Errorf(ctx, "[iif_data_reader.read] iif line expected reading split sign or transaction end sign, but actual is \"%s\"", items[0])
					return nil, nil, errs.ErrInvalidIIFFile
				}
			} else {
				log.Errorf(ctx, "[iif_data_reader.read] iif missing transaction sample end line")
				return nil, nil, errs.ErrInvalidIIFFile
			}
		} // not process (read data line) for other dataset type
	}

	if lastLineSign != "" {
		log.Errorf(ctx, "[iif_data_reader.read] iif missing transaction end line")
		return nil, nil, errs.ErrInvalidIIFFile
	}

	if currentAccountDataset != nil {
		allAccountDatasets = append(allAccountDatasets, currentAccountDataset)
	}

	if currentTransactionDataset != nil {
		allTransactionDatasets = append(allTransactionDatasets, currentTransactionDataset)
	}

	return allAccountDatasets, allTransactionDatasets, nil
}

func (r *iifDataReader) readAccountSampleLine(ctx core.Context, items []string) (*iifAccountDataset, error) {
	accountSampleItems := items
	accountDataColumnIndexes := make(map[string]int, len(accountSampleItems))

	for i := 1; i < len(accountSampleItems); i++ {
		columnName := accountSampleItems[i]
		accountDataColumnIndexes[columnName] = i
	}

	return &iifAccountDataset{
		accountDataColumnIndexes: accountDataColumnIndexes,
		accounts:                 make([]*iifAccountData, 0),
	}, nil
}

func (r *iifDataReader) readTransactionSampleLines(ctx core.Context, items []string) (*iifTransactionDataset, error) {
	transactionSampleItems := items
	transactionDataColumnIndexes := make(map[string]int, len(transactionSampleItems))

	for i := 1; i < len(transactionSampleItems); i++ {
		columnName := transactionSampleItems[i]
		transactionDataColumnIndexes[columnName] = i
	}

	splitSampleItems, err := r.reader.Read()

	if err == io.EOF {
		log.Errorf(ctx, "[iif_data_reader.readTransactionSampleLines] expected reading transaction split sample line, but read eof")
		return nil, errs.ErrInvalidIIFFile
	}

	if len(splitSampleItems) < 1 || splitSampleItems[0] != iifTransactionSplitSampleLineSignColumnName {
		log.Errorf(ctx, "[iif_data_reader.readTransactionSampleLines] expected reading transaction split sample line, but read \"%s\"", strings.Join(splitSampleItems, "\t"))
		return nil, errs.ErrInvalidIIFFile
	}

	splitDataColumnIndexes := make(map[string]int, len(splitSampleItems))

	for i := 1; i < len(splitSampleItems); i++ {
		columnName := splitSampleItems[i]
		splitDataColumnIndexes[columnName] = i
	}

	transactionEndSampleItems, err := r.reader.Read()

	if err == io.EOF {
		log.Errorf(ctx, "[iif_data_reader.readTransactionSampleLines] expected reading transaction end sample line, but read eof")
		return nil, errs.ErrInvalidIIFFile
	}

	if len(transactionEndSampleItems) < 1 || transactionEndSampleItems[0] != iifTransactionEndSampleLineSignColumnName {
		log.Errorf(ctx, "[iif_data_reader.readTransactionSampleLines] expected reading transaction end sample line, but read \"%s\"", strings.Join(transactionEndSampleItems, "\t"))
		return nil, errs.ErrInvalidIIFFile
	}

	return &iifTransactionDataset{
		transactionDataColumnIndexes: transactionDataColumnIndexes,
		splitDataColumnIndexes:       splitDataColumnIndexes,
		transactions:                 make([]*iifTransactionData, 0),
	}, nil
}

func createNewIifDataReader(data []byte) *iifDataReader {
	reader := bytes.NewReader(data)
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = -1

	return &iifDataReader{
		reader: csvReader,
	}
}
