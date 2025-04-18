package converter

import (
	"sort"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// DataTableTransactionDataImporter defines the structure of plain text data table importer for transaction data
type DataTableTransactionDataImporter struct {
	transactionTypeMapping  map[string]models.TransactionType
	geoLocationSeparator    string
	transactionTagSeparator string
}

// ParseImportedData returns the imported transaction data
func (c *DataTableTransactionDataImporter) ParseImportedData(ctx core.Context, user *models.User, dataTable datatable.TransactionDataTable, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	if dataTable.TransactionRowCount() < 1 {
		log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse import data for user \"uid:%d\", because data table row count is less 1", user.Uid)
		return nil, nil, nil, nil, nil, nil, errs.ErrNotFoundTransactionDataInFile
	}

	nameDbTypeMap, err := c.buildTransactionTypeNameDbTypeMap()

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if !dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME) ||
		!dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE) ||
		!dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY) ||
		!dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME) ||
		!dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_AMOUNT) ||
		!dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME) {
		log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse import data for user \"uid:%d\", because missing essential columns in header row", user.Uid)
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if accountMap == nil {
		accountMap = make(map[string]*models.Account)
	}

	if expenseCategoryMap == nil {
		expenseCategoryMap = make(map[string]map[string]*models.TransactionCategory)
	}

	if incomeCategoryMap == nil {
		incomeCategoryMap = make(map[string]map[string]*models.TransactionCategory)
	}

	if transferCategoryMap == nil {
		transferCategoryMap = make(map[string]map[string]*models.TransactionCategory)
	}

	if tagMap == nil {
		tagMap = make(map[string]*models.TransactionTag)
	}

	allNewTransactions := make(models.ImportedTransactionSlice, 0, dataTable.TransactionRowCount())
	allNewAccounts := make([]*models.Account, 0)
	allNewSubExpenseCategories := make([]*models.TransactionCategory, 0)
	allNewSubIncomeCategories := make([]*models.TransactionCategory, 0)
	allNewSubTransferCategories := make([]*models.TransactionCategory, 0)
	allNewTags := make([]*models.TransactionTag, 0)

	dataRowIterator := dataTable.TransactionRowIterator()
	dataRowIndex := 0

	for dataRowIterator.HasNext() {
		dataRowIndex++
		dataRow, err := dataRowIterator.Next(ctx, user)

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse data row \"index:%d\" for user \"uid:%d\", because %s", dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, nil, nil, err
		}

		if !dataRow.IsValid() {
			continue
		}

		timezoneOffset := defaultTimezoneOffset

		if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE) {
			transactionTimezone, err := utils.ParseFromTimezoneOffset(dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE))

			if err != nil {
				log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse time zone \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE), dataRowIndex, user.Uid, err.Error())
				return nil, nil, nil, nil, nil, nil, errs.ErrTransactionTimeZoneInvalid
			}

			timezoneOffset = utils.GetTimezoneOffsetMinutes(transactionTimezone)
		}

		transactionTime, err := utils.ParseFromLongDateTime(dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME), timezoneOffset)

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse time \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, nil, nil, errs.ErrTransactionTimeInvalid
		}

		transactionDbType, err := c.getTransactionDbType(nameDbTypeMap, dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE))

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse transaction type \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, nil, nil, errs.Or(err, errs.ErrTransactionTypeInvalid)
		}

		categoryId := int64(0)
		categoryName := ""
		subCategoryName := ""

		if transactionDbType != models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			transactionCategoryType, err := c.getTransactionCategoryType(transactionDbType)

			if err != nil {
				log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse transaction category type in data row \"index:%d\" for user \"uid:%d\", because %s", dataRowIndex, user.Uid, err.Error())
				return nil, nil, nil, nil, nil, nil, errs.Or(err, errs.ErrTransactionTypeInvalid)
			}

			categoryName = dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_CATEGORY)
			subCategoryName = dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY)

			if transactionDbType == models.TRANSACTION_DB_TYPE_EXPENSE {
				subCategory, exists := c.getTransactionCategory(expenseCategoryMap, categoryName, subCategoryName)

				if !exists {
					subCategory = c.createNewTransactionCategoryModel(user.Uid, subCategoryName, transactionCategoryType)
					allNewSubExpenseCategories = append(allNewSubExpenseCategories, subCategory)

					if _, exists = expenseCategoryMap[subCategoryName]; !exists {
						expenseCategoryMap[subCategoryName] = make(map[string]*models.TransactionCategory)
					}

					expenseCategoryMap[subCategoryName][categoryName] = subCategory
				}

				categoryId = subCategory.CategoryId
			} else if transactionDbType == models.TRANSACTION_DB_TYPE_INCOME {
				subCategory, exists := c.getTransactionCategory(incomeCategoryMap, categoryName, subCategoryName)

				if !exists {
					subCategory = c.createNewTransactionCategoryModel(user.Uid, subCategoryName, transactionCategoryType)
					allNewSubIncomeCategories = append(allNewSubIncomeCategories, subCategory)

					if _, exists = incomeCategoryMap[subCategoryName]; !exists {
						incomeCategoryMap[subCategoryName] = make(map[string]*models.TransactionCategory)
					}

					incomeCategoryMap[subCategoryName][categoryName] = subCategory
				}

				categoryId = subCategory.CategoryId
			} else if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
				subCategory, exists := c.getTransactionCategory(transferCategoryMap, categoryName, subCategoryName)

				if !exists {
					subCategory = c.createNewTransactionCategoryModel(user.Uid, subCategoryName, transactionCategoryType)
					allNewSubTransferCategories = append(allNewSubTransferCategories, subCategory)

					if _, exists = transferCategoryMap[subCategoryName]; !exists {
						transferCategoryMap[subCategoryName] = make(map[string]*models.TransactionCategory)
					}

					transferCategoryMap[subCategoryName][categoryName] = subCategory
				}

				categoryId = subCategory.CategoryId
			}
		}

		accountName := dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME)
		accountCurrency := user.DefaultCurrency

		if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY) {
			accountCurrency = dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY)

			if _, ok := validators.AllCurrencyNames[accountCurrency]; !ok {
				log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] account currency \"%s\" is not supported in data row \"index:%d\" for user \"uid:%d\"", accountCurrency, dataRowIndex, user.Uid)
				return nil, nil, nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
			}
		}

		account, exists := accountMap[accountName]

		if !exists {
			account = c.createNewAccountModel(user.Uid, accountName, accountCurrency)
			allNewAccounts = append(allNewAccounts, account)
			accountMap[accountName] = account
		}

		if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY) {
			if account.Name != "" && account.Currency != accountCurrency {
				log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] currency \"%s\" in data row \"index:%d\" not equals currency \"%s\" of the account for user \"uid:%d\"", accountCurrency, dataRowIndex, account.Currency, user.Uid)
				return nil, nil, nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
			}
		} else if exists {
			accountCurrency = account.Currency
		}

		amount, err := utils.ParseAmount(dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_AMOUNT))

		if err != nil {
			log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse acmount \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_AMOUNT), dataRowIndex, user.Uid, err.Error())
			return nil, nil, nil, nil, nil, nil, errs.ErrAmountInvalid
		}

		relatedAccountId := int64(0)
		relatedAccountAmount := int64(0)
		account2Name := ""
		account2Currency := ""

		if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2Name = dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME)
			account2Currency = user.DefaultCurrency

			if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY) {
				account2Currency = dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY)

				if _, ok := validators.AllCurrencyNames[account2Currency]; !ok {
					log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] account2 currency \"%s\" is not supported in data row \"index:%d\" for user \"uid:%d\"", account2Currency, dataRowIndex, user.Uid)
					return nil, nil, nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			}

			account2, exists := accountMap[account2Name]

			if !exists {
				account2 = c.createNewAccountModel(user.Uid, account2Name, account2Currency)
				allNewAccounts = append(allNewAccounts, account2)
				accountMap[account2Name] = account2
			}

			if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY) {
				if account2.Name != "" && account2.Currency != account2Currency {
					log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] currency \"%s\" in data row \"index:%d\" not equals currency \"%s\" of the account2 for user \"uid:%d\"", account2Currency, dataRowIndex, account2.Currency, user.Uid)
					return nil, nil, nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			} else if exists {
				account2Currency = account2.Currency
			}

			relatedAccountId = account2.AccountId

			if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT) {
				relatedAccountAmount, err = utils.ParseAmount(dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT))

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse acmount2 \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, nil, nil, errs.ErrAmountInvalid
				}
			} else if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
				relatedAccountAmount = amount
			}
		}

		geoLongitude := float64(0)
		geoLatitude := float64(0)

		if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION) && c.geoLocationSeparator != "" {
			geoLocationItems := strings.Split(dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION), c.geoLocationSeparator)

			if len(geoLocationItems) == 2 {
				geoLongitude, err = utils.StringToFloat64(geoLocationItems[0])

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse geographic location \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, nil, nil, errs.ErrGeographicLocationInvalid
				}

				geoLatitude, err = utils.StringToFloat64(geoLocationItems[1])

				if err != nil {
					log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] cannot parse geographic location \"%s\" in data row \"index:%d\" for user \"uid:%d\", because %s", dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION), dataRowIndex, user.Uid, err.Error())
					return nil, nil, nil, nil, nil, nil, errs.ErrGeographicLocationInvalid
				}
			}
		}

		var tagIds []string
		var tagNames []string

		if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_TAGS) {
			var tagNameItems []string

			if c.transactionTagSeparator != "" {
				tagNameItems = strings.Split(dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TAGS), c.transactionTagSeparator)
			} else {
				tagNameItems = append(tagNameItems, dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_TAGS))
			}

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

		if dataTable.HasColumn(datatable.TRANSACTION_DATA_TABLE_DESCRIPTION) {
			description = dataRow.GetData(datatable.TRANSACTION_DATA_TABLE_DESCRIPTION)
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

	if len(allNewTransactions) < 1 {
		log.Errorf(ctx, "[data_table_transaction_data_exporter.ParseImportedData] no transaction data parsed for \"uid:%d\"", user.Uid)
		return nil, nil, nil, nil, nil, nil, errs.ErrNotFoundTransactionDataInFile
	}

	sort.Sort(allNewTransactions)

	return allNewTransactions, allNewAccounts, allNewSubExpenseCategories, allNewSubIncomeCategories, allNewSubTransferCategories, allNewTags, nil
}

func (c *DataTableTransactionDataImporter) buildTransactionTypeNameDbTypeMap() (map[string]models.TransactionDbType, error) {
	if c.transactionTypeMapping == nil {
		return nil, errs.ErrTransactionTypeInvalid
	}

	nameDbTypeMap := make(map[string]models.TransactionDbType, len(c.transactionTypeMapping))

	for name, transactionType := range c.transactionTypeMapping {
		if transactionType == models.TRANSACTION_TYPE_MODIFY_BALANCE {
			nameDbTypeMap[name] = models.TRANSACTION_DB_TYPE_MODIFY_BALANCE
		} else if transactionType == models.TRANSACTION_TYPE_INCOME {
			nameDbTypeMap[name] = models.TRANSACTION_DB_TYPE_INCOME
		} else if transactionType == models.TRANSACTION_TYPE_EXPENSE {
			nameDbTypeMap[name] = models.TRANSACTION_DB_TYPE_EXPENSE
		} else if transactionType == models.TRANSACTION_TYPE_TRANSFER {
			nameDbTypeMap[name] = models.TRANSACTION_DB_TYPE_TRANSFER_OUT
		} else {
			return nil, errs.ErrTransactionTypeInvalid
		}
	}

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

func (c *DataTableTransactionDataImporter) getTransactionCategory(categories map[string]map[string]*models.TransactionCategory, categoryName string, subCategoryName string) (*models.TransactionCategory, bool) {
	if len(categories) < 1 {
		return nil, false
	}

	subCategories, exists := categories[subCategoryName]

	if !exists || len(subCategories) < 1 {
		return nil, false
	}

	if categoryName == "" {
		for _, subCategory := range subCategories {
			if subCategory != nil {
				return subCategory, true
			}
		}
	}

	subCategory, exists := subCategories[categoryName]

	if !exists {
		for _, subCategory := range subCategories {
			if subCategory != nil {
				return subCategory, true
			}
		}
	}

	return subCategory, exists
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

// CreateNewImporterWithTypeNameMapping returns a new data table transaction data importer according to the specified arguments
func CreateNewImporterWithTypeNameMapping(transactionTypeMapping map[models.TransactionType]string, geoLocationSeparator string, transactionTagSeparator string) *DataTableTransactionDataImporter {
	return &DataTableTransactionDataImporter{
		transactionTypeMapping:  buildTransactionNameTypeMap(transactionTypeMapping),
		geoLocationSeparator:    geoLocationSeparator,
		transactionTagSeparator: transactionTagSeparator,
	}
}

// CreateNewSimpleImporter returns a new data table transaction data importer according to the specified arguments
func CreateNewSimpleImporter(transactionTypeMapping map[string]models.TransactionType) *DataTableTransactionDataImporter {
	return &DataTableTransactionDataImporter{
		transactionTypeMapping: transactionTypeMapping,
	}
}

// CreateNewSimpleImporterWithTypeNameMapping returns a new data table transaction data importer according to the specified arguments
func CreateNewSimpleImporterWithTypeNameMapping(transactionTypeMapping map[models.TransactionType]string) *DataTableTransactionDataImporter {
	return &DataTableTransactionDataImporter{
		transactionTypeMapping: buildTransactionNameTypeMap(transactionTypeMapping),
	}
}

func buildTransactionNameTypeMap(transactionTypeMapping map[models.TransactionType]string) map[string]models.TransactionType {
	if transactionTypeMapping == nil {
		return nil
	}

	typeNameMap := make(map[string]models.TransactionType, len(transactionTypeMapping))

	for transactionType, name := range transactionTypeMapping {
		typeNameMap[name] = transactionType
	}

	return typeNameMap
}
