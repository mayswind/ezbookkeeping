package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// DataTableTransactionDataExporter defines the structure of plain text data table exporter for transaction data
type DataTableTransactionDataExporter struct {
	transactionTypeMapping  map[models.TransactionType]string
	geoLocationSeparator    string
	transactionTagSeparator string
}

// BuildExportedContent writes the exported transaction data to the data table builder
func (c *DataTableTransactionDataExporter) BuildExportedContent(ctx core.Context, dataTableBuilder datatable.TransactionDataTableBuilder, uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) error {
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			continue
		}

		dataRowMap := make(map[datatable.TransactionDataTableColumn]string, 15)
		transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)

		dataRowMap[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), transactionTimeZone)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(transactionTimeZone)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = dataTableBuilder.ReplaceDelimiters(c.getDisplayTransactionTypeName(transaction.Type))
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = c.getExportedTransactionCategoryName(dataTableBuilder, transaction.CategoryId, categoryMap)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = c.getExportedTransactionSubCategoryName(dataTableBuilder, transaction.CategoryId, categoryMap)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = c.getExportedAccountName(dataTableBuilder, transaction.AccountId, accountMap)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = c.getAccountCurrency(dataTableBuilder, transaction.AccountId, accountMap)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = utils.FormatAmount(transaction.Amount)

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			dataRowMap[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = c.getExportedAccountName(dataTableBuilder, transaction.RelatedAccountId, accountMap)
			dataRowMap[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = c.getAccountCurrency(dataTableBuilder, transaction.RelatedAccountId, accountMap)
			dataRowMap[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(transaction.RelatedAccountAmount)
		}

		dataRowMap[datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION] = c.getExportedGeographicLocation(transaction)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_TAGS] = c.getExportedTags(dataTableBuilder, transaction.TransactionId, allTagIndexes, tagMap)
		dataRowMap[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = dataTableBuilder.ReplaceDelimiters(transaction.Comment)

		dataTableBuilder.AppendTransaction(dataRowMap)
	}

	return nil
}

func (c *DataTableTransactionDataExporter) getDisplayTransactionTypeName(transactionDbType models.TransactionDbType) string {
	transactionType, err := transactionDbType.ToTransactionType()

	if err != nil {
		return ""
	}

	transactionTypeName, exists := c.transactionTypeMapping[transactionType]

	if !exists {
		return ""
	}

	return transactionTypeName
}

func (c *DataTableTransactionDataExporter) getExportedTransactionCategoryName(dataTableBuilder datatable.TransactionDataTableBuilder, categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if !exists {
		return ""
	}

	if category.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
		return dataTableBuilder.ReplaceDelimiters(category.Name)
	}

	parentCategory, exists := categoryMap[category.ParentCategoryId]

	if !exists {
		return ""
	}

	return dataTableBuilder.ReplaceDelimiters(parentCategory.Name)
}

func (c *DataTableTransactionDataExporter) getExportedTransactionSubCategoryName(dataTableBuilder datatable.TransactionDataTableBuilder, categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if exists {
		return dataTableBuilder.ReplaceDelimiters(category.Name)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataExporter) getExportedAccountName(dataTableBuilder datatable.TransactionDataTableBuilder, accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return dataTableBuilder.ReplaceDelimiters(account.Name)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataExporter) getAccountCurrency(dataTableBuilder datatable.TransactionDataTableBuilder, accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return dataTableBuilder.ReplaceDelimiters(account.Currency)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataExporter) getExportedGeographicLocation(transaction *models.Transaction) string {
	if transaction.GeoLongitude != 0 || transaction.GeoLatitude != 0 {
		return fmt.Sprintf("%f%s%f", transaction.GeoLongitude, c.geoLocationSeparator, transaction.GeoLatitude)
	}

	return ""
}

func (c *DataTableTransactionDataExporter) getExportedTags(dataTableBuilder datatable.TransactionDataTableBuilder, transactionId int64, allTagIndexes map[int64][]int64, tagMap map[int64]*models.TransactionTag) string {
	tagIndexes, exists := allTagIndexes[transactionId]

	if !exists {
		return ""
	}

	var ret strings.Builder

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]
		tag, exists := tagMap[tagIndex]

		if !exists {
			continue
		}

		if ret.Len() > 0 {
			ret.WriteString(c.transactionTagSeparator)
		}

		ret.WriteString(strings.Replace(tag.Name, c.transactionTagSeparator, " ", -1))
	}

	return dataTableBuilder.ReplaceDelimiters(ret.String())
}

// CreateNewExporter returns a new data table transaction data exporter according to the specified arguments
func CreateNewExporter(transactionTypeMapping map[models.TransactionType]string, geoLocationSeparator string, transactionTagSeparator string) *DataTableTransactionDataExporter {
	return &DataTableTransactionDataExporter{
		transactionTypeMapping:  transactionTypeMapping,
		geoLocationSeparator:    geoLocationSeparator,
		transactionTagSeparator: transactionTagSeparator,
	}
}
