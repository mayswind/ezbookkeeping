package converters

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// DataTableColumn represents the data column type of data table
type DataTableColumn byte

// Data table columns
const (
	DATA_TABLE_TRANSACTION_TIME         DataTableColumn = 1
	DATA_TABLE_TRANSACTION_TIMEZONE     DataTableColumn = 2
	DATA_TABLE_TRANSACTION_TYPE         DataTableColumn = 3
	DATA_TABLE_CATEGORY                 DataTableColumn = 4
	DATA_TABLE_SUB_CATEGORY             DataTableColumn = 5
	DATA_TABLE_ACCOUNT_NAME             DataTableColumn = 6
	DATA_TABLE_ACCOUNT_CURRENCY         DataTableColumn = 7
	DATA_TABLE_AMOUNT                   DataTableColumn = 8
	DATA_TABLE_RELATED_ACCOUNT_NAME     DataTableColumn = 9
	DATA_TABLE_RELATED_ACCOUNT_CURRENCY DataTableColumn = 10
	DATA_TABLE_RELATED_AMOUNT           DataTableColumn = 11
	DATA_TABLE_GEOGRAPHIC_LOCATION      DataTableColumn = 12
	DATA_TABLE_TAGS                     DataTableColumn = 13
	DATA_TABLE_DESCRIPTION              DataTableColumn = 14
)

// DataTableTransactionDataExporter defines the structure of plain text data table exporter for transaction data
type DataTableTransactionDataExporter struct {
	dataColumnMapping       map[DataTableColumn]string
	transactionTypeMapping  map[models.TransactionType]string
	geoLocationSeparator    string
	transactionTagSeparator string
}

// DataTableTransactionDataImporter defines the structure of plain text data table importer for transaction data
type DataTableTransactionDataImporter struct {
	dataColumnMapping       map[DataTableColumn]string
	transactionTypeMapping  map[models.TransactionType]string
	geoLocationSeparator    string
	transactionTagSeparator string
}

func (c *DataTableTransactionDataExporter) buildExportedContent(ctx core.Context, dataTableBuilder DataTableBuilder, uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) error {
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			continue
		}

		dataRowMap := make(map[DataTableColumn]string, 15)
		transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)

		dataRowMap[DATA_TABLE_TRANSACTION_TIME] = utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), transactionTimeZone)
		dataRowMap[DATA_TABLE_TRANSACTION_TIMEZONE] = utils.FormatTimezoneOffset(transactionTimeZone)
		dataRowMap[DATA_TABLE_TRANSACTION_TYPE] = dataTableBuilder.ReplaceDelimiters(c.getDisplayTransactionTypeName(transaction.Type))
		dataRowMap[DATA_TABLE_CATEGORY] = c.getExportedTransactionCategoryName(dataTableBuilder, transaction.CategoryId, categoryMap)
		dataRowMap[DATA_TABLE_SUB_CATEGORY] = c.getExportedTransactionSubCategoryName(dataTableBuilder, transaction.CategoryId, categoryMap)
		dataRowMap[DATA_TABLE_ACCOUNT_NAME] = c.getExportedAccountName(dataTableBuilder, transaction.AccountId, accountMap)
		dataRowMap[DATA_TABLE_ACCOUNT_CURRENCY] = c.getAccountCurrency(dataTableBuilder, transaction.AccountId, accountMap)
		dataRowMap[DATA_TABLE_AMOUNT] = utils.FormatAmount(transaction.Amount)

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			dataRowMap[DATA_TABLE_RELATED_ACCOUNT_NAME] = c.getExportedAccountName(dataTableBuilder, transaction.RelatedAccountId, accountMap)
			dataRowMap[DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = c.getAccountCurrency(dataTableBuilder, transaction.RelatedAccountId, accountMap)
			dataRowMap[DATA_TABLE_RELATED_AMOUNT] = utils.FormatAmount(transaction.RelatedAccountAmount)
		}

		dataRowMap[DATA_TABLE_GEOGRAPHIC_LOCATION] = c.getExportedGeographicLocation(transaction)
		dataRowMap[DATA_TABLE_TAGS] = c.getExportedTags(dataTableBuilder, transaction.TransactionId, allTagIndexes, tagMap)
		dataRowMap[DATA_TABLE_DESCRIPTION] = dataTableBuilder.ReplaceDelimiters(transaction.Comment)

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

func (c *DataTableTransactionDataExporter) getExportedTransactionCategoryName(dataTableBuilder DataTableBuilder, categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if !exists {
		return ""
	}

	if category.ParentCategoryId == 0 {
		return dataTableBuilder.ReplaceDelimiters(category.Name)
	}

	parentCategory, exists := categoryMap[category.ParentCategoryId]

	if !exists {
		return ""
	}

	return dataTableBuilder.ReplaceDelimiters(parentCategory.Name)
}

func (c *DataTableTransactionDataExporter) getExportedTransactionSubCategoryName(dataTableBuilder DataTableBuilder, categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if exists {
		return dataTableBuilder.ReplaceDelimiters(category.Name)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataExporter) getExportedAccountName(dataTableBuilder DataTableBuilder, accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return dataTableBuilder.ReplaceDelimiters(account.Name)
	} else {
		return ""
	}
}

func (c *DataTableTransactionDataExporter) getAccountCurrency(dataTableBuilder DataTableBuilder, accountId int64, accountMap map[int64]*models.Account) string {
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

func (c *DataTableTransactionDataExporter) getExportedTags(dataTableBuilder DataTableBuilder, transactionId int64, allTagIndexes map[int64][]int64, tagMap map[int64]*models.TransactionTag) string {
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

func (c *DataTableTransactionDataImporter) parseImportedData(ctx core.Context, user *models.User, dataTable ImportedDataTable, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	if dataTable.DataRowCount() < 1 {
		log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse import data for user \"uid:%d\", because data table row count is less 1", user.Uid)
		return nil, nil, nil, nil, errs.ErrNotFoundTransactionDataInFile
	}

	nameDbTypeMap, err := c.buildTransactionTypeNameDbTypeMap()

	if err != nil {
		return nil, nil, nil, nil, err
	}

	headerLineItems := dataTable.HeaderLineColumnNames()
	headerItemMap := make(map[string]int)

	for i := 0; i < len(headerLineItems); i++ {
		headerItemMap[headerLineItems[i]] = i
	}

	timeColumnIdx, timeColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TRANSACTION_TIME]]
	timezoneColumnIdx, timezoneColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TRANSACTION_TIMEZONE]]
	typeColumnIdx, typeColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TRANSACTION_TYPE]]
	subCategoryColumnIdx, subCategoryColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_SUB_CATEGORY]]
	accountColumnIdx, accountColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_ACCOUNT_NAME]]
	accountCurrencyColumnIdx, accountCurrencyColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_ACCOUNT_CURRENCY]]
	amountColumnIdx, amountColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_AMOUNT]]
	account2ColumnIdx, account2ColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_RELATED_ACCOUNT_NAME]]
	account2CurrencyColumnIdx, account2CurrencyColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_RELATED_ACCOUNT_CURRENCY]]
	amount2ColumnIdx, amount2ColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_RELATED_AMOUNT]]
	geoLocationIdx, geoLocationExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_GEOGRAPHIC_LOCATION]]
	tagsColumnIdx, tagsColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_TAGS]]
	descriptionColumnIdx, descriptionColumnExists := headerItemMap[c.dataColumnMapping[DATA_TABLE_DESCRIPTION]]

	if !timeColumnExists || !typeColumnExists || !subCategoryColumnExists ||
		!accountColumnExists || !amountColumnExists || !account2ColumnExists {
		log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse import data for user \"uid:%d\", because missing essential columns in header row", user.Uid)
		return nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if accountMap == nil {
		accountMap = make(map[string]*models.Account)
	}

	if categoryMap == nil {
		categoryMap = make(map[string]*models.TransactionCategory)
	}

	if tagMap == nil {
		tagMap = make(map[string]*models.TransactionTag)
	}

	allNewTransactions := make(models.ImportedTransactionSlice, 0, dataTable.DataRowCount())
	allNewAccounts := make([]*models.Account, 0)
	allNewSubCategories := make([]*models.TransactionCategory, 0)
	allNewTags := make([]*models.TransactionTag, 0)

	dataRowIterator := dataTable.DataRowIterator()
	dataRowIndex := 0

	for dataRowIterator.HasNext() {
		dataRowIndex++
		dataRow := dataRowIterator.Next()
		columnCount := dataRow.ColumnCount()

		if columnCount < 1 || (columnCount == 1 && dataRow.GetData(0) == "") {
			continue
		}

		if columnCount < len(headerLineItems) {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse data row \"index:%d\" for user \"uid:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", dataRowIndex, user.Uid, columnCount, len(headerLineItems))
			return nil, nil, nil, nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
		}

		timezoneOffset := defaultTimezoneOffset

		if timezoneColumnExists {
			transactionTimezone, err := dataRow.GetTimezoneOffset(timezoneColumnIdx)

			if err != nil {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse time zone \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(timezoneColumnIdx), dataRowIndex, user.Uid, err.Error())
				return nil, nil, nil, nil, errs.ErrTransactionTimeZoneInvalid
			}

			timezoneOffset = utils.GetTimezoneOffsetMinutes(transactionTimezone)
		}

		transactionTime, err := dataRow.GetTime(timeColumnIdx, timezoneOffset)

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse time \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(timeColumnIdx), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, errs.ErrTransactionTimeInvalid
		}

		transactionDbType, err := c.getTransactionDbType(nameDbTypeMap, dataRow.GetData(typeColumnIdx))

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse transaction type \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(typeColumnIdx), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, errs.Or(err, errs.ErrTransactionTypeInvalid)
		}

		categoryId := int64(0)
		subCategoryName := ""

		if transactionDbType != models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			transactionCategoryType, err := c.getTransactionCategoryType(transactionDbType)

			if err != nil {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse transaction category type in data row \"index:%d\" for user \"uid:%d\", because %s", dataRowIndex, user.Uid, err.Error())
				return nil, nil, nil, nil, errs.Or(err, errs.ErrTransactionTypeInvalid)
			}

			subCategoryName = dataRow.GetData(subCategoryColumnIdx)
			subCategory, exists := categoryMap[subCategoryName]

			if !exists {
				subCategory = c.createNewTransactionCategoryModel(user.Uid, subCategoryName, transactionCategoryType)
				allNewSubCategories = append(allNewSubCategories, subCategory)
				categoryMap[subCategoryName] = subCategory
			}

			categoryId = subCategory.CategoryId
		}

		accountName := dataRow.GetData(accountColumnIdx)

		if accountName == "" {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account name is empty in data row \"index:%d\" for user \"uid:%d\"", dataRowIndex, user.Uid)
			return nil, nil, nil, nil, errs.ErrAccountNameCannotBeBlank
		}

		accountCurrency := user.DefaultCurrency

		if accountCurrencyColumnExists {
			accountCurrency = dataRow.GetData(accountCurrencyColumnIdx)

			if _, ok := validators.AllCurrencyNames[accountCurrency]; !ok {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account currency \"%s\" is not supported in data row \"index:%d\" for user \"uid:%d\"", accountCurrency, dataRowIndex, user.Uid)
				return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
			}
		}

		account, exists := accountMap[accountName]

		if !exists {
			account = c.createNewAccountModel(user.Uid, accountName, accountCurrency)
			allNewAccounts = append(allNewAccounts, account)
			accountMap[accountName] = account
		}

		if accountCurrencyColumnExists {
			if account.Currency != accountCurrency {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] currency \"%s\" in data row \"index:%d\" not equals currency \"%s\" of the account for user \"uid:%d\"", accountCurrency, dataRowIndex, account.Currency, user.Uid)
				return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
			}
		} else if exists {
			accountCurrency = account.Currency
		}

		amount, err := utils.ParseAmount(dataRow.GetData(amountColumnIdx))

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse acmount \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(amountColumnIdx), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, errs.ErrAmountInvalid
		}

		relatedAccountId := int64(0)
		relatedAccountAmount := int64(0)
		account2Name := ""
		account2Currency := ""

		if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2Name = dataRow.GetData(account2ColumnIdx)

			if account2Name == "" {
				log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account2 name is empty in data row \"index:%d\" for user \"uid:%d\"", dataRowIndex, user.Uid)
				return nil, nil, nil, nil, errs.ErrDestinationAccountNameCannotBeBlank
			}

			account2Currency = user.DefaultCurrency

			if account2CurrencyColumnExists {
				account2Currency = dataRow.GetData(account2CurrencyColumnIdx)

				if _, ok := validators.AllCurrencyNames[account2Currency]; !ok {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] account2 currency \"%s\" is not supported in data row \"index:%d\" for user \"uid:%d\"", account2Currency, dataRowIndex, user.Uid)
					return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			}

			account2, exists := accountMap[account2Name]

			if !exists {
				account2 = c.createNewAccountModel(user.Uid, account2Name, account2Currency)
				allNewAccounts = append(allNewAccounts, account2)
				accountMap[account2Name] = account2
			}

			if account2CurrencyColumnExists {
				if account2.Currency != account2Currency {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] currency \"%s\" in data row \"index:%d\" not equals currency \"%s\" of the account2 for user \"uid:%d\"", account2Currency, dataRowIndex, account2.Currency, user.Uid)
					return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			} else if exists {
				account2Currency = account2.Currency
			}

			relatedAccountId = account2.AccountId

			if amount2ColumnExists {
				relatedAccountAmount, err = utils.ParseAmount(dataRow.GetData(amount2ColumnIdx))

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse acmount2 \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(amount2ColumnIdx), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, errs.ErrAmountInvalid
				}
			} else if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
				relatedAccountAmount = amount
			}
		}

		geoLongitude := float64(0)
		geoLatitude := float64(0)

		if geoLocationExists {
			geoLocationItems := strings.Split(dataRow.GetData(geoLocationIdx), c.geoLocationSeparator)

			if len(geoLocationItems) == 2 {
				geoLongitude, err = utils.StringToFloat64(geoLocationItems[0])

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse geographic location \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(geoLocationIdx), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, errs.ErrGeographicLocationInvalid
				}

				geoLatitude, err = utils.StringToFloat64(geoLocationItems[1])

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_converter.parseImportedData] cannot parse geographic location \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(geoLocationIdx), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, errs.ErrGeographicLocationInvalid
				}
			}
		}

		var tagIds []string
		var tagNames []string

		if tagsColumnExists {
			tagNameItems := strings.Split(dataRow.GetData(tagsColumnIdx), c.transactionTagSeparator)

			for i := 0; i < len(tagNameItems); i++ {
				tagName := tagNameItems[i]

				if tagName == "" {
					continue
				}

				tag, exists := tagMap[tagName]

				if !exists {
					tag = c.createNewTransactionTagModel(user.Uid, tagName)
					allNewTags = append(allNewTags, tag)
					tagMap[tagName] = tag
				}

				if tag != nil {
					tagIds = append(tagIds, utils.Int64ToString(tag.TagId))
				}

				tagNames = append(tagNames, tagName)
			}
		}

		description := ""

		if descriptionColumnExists {
			description = dataRow.GetData(descriptionColumnIdx)
		}

		transaction := &models.ImportTransaction{
			Transaction: &models.Transaction{
				Uid:                  user.Uid,
				Type:                 transactionDbType,
				CategoryId:           categoryId,
				TransactionTime:      utils.GetMinTransactionTimeFromUnixTime(transactionTime.Unix()),
				TimezoneUtcOffset:    timezoneOffset,
				AccountId:            account.AccountId,
				Amount:               amount,
				HideAmount:           false,
				RelatedAccountId:     relatedAccountId,
				RelatedAccountAmount: relatedAccountAmount,
				Comment:              description,
				GeoLongitude:         geoLongitude,
				GeoLatitude:          geoLatitude,
				CreatedIp:            "127.0.0.1",
			},
			TagIds:                             tagIds,
			OriginalCategoryName:               subCategoryName,
			OriginalSourceAccountName:          accountName,
			OriginalSourceAccountCurrency:      accountCurrency,
			OriginalDestinationAccountName:     account2Name,
			OriginalDestinationAccountCurrency: account2Currency,
			OriginalTagNames:                   tagNames,
		}

		allNewTransactions = append(allNewTransactions, transaction)
	}

	sort.Sort(allNewTransactions)

	return allNewTransactions, allNewAccounts, allNewSubCategories, allNewTags, nil
}

func (c *DataTableTransactionDataImporter) buildTransactionTypeNameDbTypeMap() (map[string]models.TransactionDbType, error) {
	if c.transactionTypeMapping == nil {
		return nil, errs.ErrTransactionTypeInvalid
	}

	nameDbTypeMap := make(map[string]models.TransactionDbType, len(c.transactionTypeMapping))
	nameDbTypeMap[c.transactionTypeMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]] = models.TRANSACTION_DB_TYPE_MODIFY_BALANCE
	nameDbTypeMap[c.transactionTypeMapping[models.TRANSACTION_TYPE_INCOME]] = models.TRANSACTION_DB_TYPE_INCOME
	nameDbTypeMap[c.transactionTypeMapping[models.TRANSACTION_TYPE_EXPENSE]] = models.TRANSACTION_DB_TYPE_EXPENSE
	nameDbTypeMap[c.transactionTypeMapping[models.TRANSACTION_TYPE_TRANSFER]] = models.TRANSACTION_DB_TYPE_TRANSFER_OUT

	return nameDbTypeMap, nil
}

func (c *DataTableTransactionDataImporter) getTransactionDbType(nameDbTypeMap map[string]models.TransactionDbType, transactionTypeName string) (models.TransactionDbType, error) {
	transactionType, exists := nameDbTypeMap[transactionTypeName]

	if !exists {
		return 0, errs.ErrTransactionTypeInvalid
	}

	return transactionType, nil
}

func (c *DataTableTransactionDataImporter) getTransactionCategoryType(transactionType models.TransactionDbType) (models.TransactionCategoryType, error) {
	if transactionType == models.TRANSACTION_DB_TYPE_INCOME {
		return models.CATEGORY_TYPE_INCOME, nil
	} else if transactionType == models.TRANSACTION_DB_TYPE_EXPENSE {
		return models.CATEGORY_TYPE_EXPENSE, nil
	} else if transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		return models.CATEGORY_TYPE_TRANSFER, nil
	} else {
		return 0, errs.ErrTransactionTypeInvalid
	}
}

func (c *DataTableTransactionDataImporter) createNewAccountModel(uid int64, accountName string, currency string) *models.Account {
	return &models.Account{
		Uid:      uid,
		Name:     accountName,
		Currency: currency,
	}
}

func (c *DataTableTransactionDataImporter) createNewTransactionCategoryModel(uid int64, categoryName string, transactionCategoryType models.TransactionCategoryType) *models.TransactionCategory {
	return &models.TransactionCategory{
		Uid:  uid,
		Name: categoryName,
		Type: transactionCategoryType,
	}
}

func (c *DataTableTransactionDataImporter) createNewTransactionTagModel(uid int64, tagName string) *models.TransactionTag {
	return &models.TransactionTag{
		Uid:  uid,
		Name: tagName,
	}
}
