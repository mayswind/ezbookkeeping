package converters

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// EzBookKeepingPlainFileConverter defines the structure of plain file converter
type EzBookKeepingPlainFileConverter struct {
}

const lineSeparator = "\n"
const geoLocationSeparator = " "
const transactionTagSeparator = ";"
const headerLine = "Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description" + lineSeparator
const dataLineFormat = "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s" + lineSeparator

// toExportedContent returns the exported plain data
func (e *EzBookKeepingPlainFileConverter) toExportedContent(uid int64, separator string, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error) {
	var ret strings.Builder

	ret.Grow(len(transactions) * 100)

	actualHeaderLine := headerLine
	actualDataLineFormat := dataLineFormat

	if separator != "," {
		actualHeaderLine = strings.Replace(headerLine, ",", separator, -1)
		actualDataLineFormat = strings.Replace(dataLineFormat, ",", separator, -1)
	}

	ret.WriteString(actualHeaderLine)

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			continue
		}

		transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
		transactionTime := utils.FormatUnixTimeToLongDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), transactionTimeZone)
		transactionTimezone := utils.FormatTimezoneOffset(transactionTimeZone)
		transactionType := e.getTransactionTypeName(transaction.Type)
		category := e.replaceDelimiters(e.getTransactionCategoryName(transaction.CategoryId, categoryMap), separator)
		subCategory := e.replaceDelimiters(e.getTransactionSubCategoryName(transaction.CategoryId, categoryMap), separator)
		account := e.replaceDelimiters(e.getAccountName(transaction.AccountId, accountMap), separator)
		accountCurrency := e.getAccountCurrency(transaction.AccountId, accountMap)
		amount := utils.FormatAmount(transaction.Amount)
		account2 := ""
		account2Currency := ""
		account2Amount := ""
		geoLocation := ""

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2 = e.replaceDelimiters(e.getAccountName(transaction.RelatedAccountId, accountMap), separator)
			account2Currency = e.getAccountCurrency(transaction.RelatedAccountId, accountMap)
			account2Amount = utils.FormatAmount(transaction.RelatedAccountAmount)
		}

		if transaction.GeoLongitude != 0 || transaction.GeoLatitude != 0 {
			geoLocation = fmt.Sprintf("%f%s%f", transaction.GeoLongitude, geoLocationSeparator, transaction.GeoLatitude)
		}

		tags := e.replaceDelimiters(e.getTags(transaction.TransactionId, allTagIndexes, tagMap), separator)
		comment := e.replaceDelimiters(transaction.Comment, separator)

		ret.WriteString(fmt.Sprintf(actualDataLineFormat, transactionTime, transactionTimezone, transactionType, category, subCategory, account, accountCurrency, amount, account2, account2Currency, account2Amount, geoLocation, tags, comment))
	}

	return []byte(ret.String()), nil
}

func (e *EzBookKeepingPlainFileConverter) parseImportedData(user *models.User, separator string, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) ([]*models.Transaction, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	lines := strings.Split(string(data), lineSeparator)

	if len(lines) < 2 {
		return nil, nil, nil, nil, errs.ErrOperationFailed
	}

	headerLineItems := strings.Split(lines[0], separator)
	headerItemMap := make(map[string]int)

	for i := 0; i < len(headerLineItems); i++ {
		headerItemMap[headerLineItems[i]] = i
	}

	timeColumnIdx, timeColumnExists := headerItemMap["Time"]
	timezoneColumnIdx, timezoneColumnExists := headerItemMap["Timezone"]
	typeColumnIdx, typeColumnExists := headerItemMap["Type"]
	subCategoryColumnIdx, subCategoryColumnExists := headerItemMap["Sub Category"]
	accountColumnIdx, accountColumnExists := headerItemMap["Account"]
	accountCurrencyColumnIdx, accountCurrencyColumnExists := headerItemMap["Account Currency"]
	amountColumnIdx, amountColumnExists := headerItemMap["Amount"]
	account2ColumnIdx, account2ColumnExists := headerItemMap["Account2"]
	account2CurrencyColumnIdx, account2CurrencyColumnExists := headerItemMap["Account2 Currency"]
	amount2ColumnIdx, amount2ColumnExists := headerItemMap["Account2 Amount"]
	geoLocationIdx, geoLocationExists := headerItemMap["Geographic Location"]
	tagsColumnIdx, tagsColumnExists := headerItemMap["Tags"]
	descriptionColumnIdx, descriptionColumnExists := headerItemMap["Description"]

	if !timeColumnExists || !typeColumnExists || !subCategoryColumnExists ||
		!accountColumnExists || !amountColumnExists || !account2ColumnExists || !amount2ColumnExists {
		return nil, nil, nil, nil, errs.ErrFormatInvalid
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

	allNewTransactions := make(ImportTransactionSlice, 0, len(lines))
	allNewAccounts := make([]*models.Account, 0)
	allNewSubCategories := make([]*models.TransactionCategory, 0)
	allNewTags := make([]*models.TransactionTag, 0)

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if len(line) < 1 {
			continue
		}

		lineItems := strings.Split(line, separator)

		if len(lineItems) < len(headerLineItems) {
			return nil, nil, nil, nil, errs.ErrFormatInvalid
		}

		timezoneOffset := defaultTimezoneOffset

		if timezoneColumnExists {
			transactionTimezone, err := utils.ParseFromTimezoneOffset(lineItems[timezoneColumnIdx])

			if err != nil {
				return nil, nil, nil, nil, err
			}

			timezoneOffset = utils.GetTimezoneOffsetMinutes(transactionTimezone)
		}

		transactionTime, err := utils.ParseFromLongDateTime(lineItems[timeColumnIdx], timezoneOffset)

		if err != nil {
			return nil, nil, nil, nil, err
		}

		transactionDbType, err := e.getTransactionDbType(lineItems[typeColumnIdx])

		if err != nil {
			return nil, nil, nil, nil, err
		}

		categoryId := int64(0)

		if transactionDbType != models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			transactionCategoryType, err := e.getTransactionCategoryType(transactionDbType)

			if err != nil {
				return nil, nil, nil, nil, err
			}

			subCategoryName := lineItems[subCategoryColumnIdx]

			if subCategoryName == "" {
				return nil, nil, nil, nil, errs.ErrFormatInvalid
			}

			subCategory, exists := categoryMap[subCategoryName]

			if !exists {
				subCategory = e.createNewTransactionCategoryModel(user.Uid, subCategoryName, transactionCategoryType)
				allNewSubCategories = append(allNewSubCategories, subCategory)
				categoryMap[subCategoryName] = subCategory
			}

			categoryId = subCategory.CategoryId
		}

		accountName := lineItems[accountColumnIdx]

		if accountName == "" {
			return nil, nil, nil, nil, errs.ErrFormatInvalid
		}

		account, exists := accountMap[accountName]

		if !exists {
			currency := user.DefaultCurrency

			if accountCurrencyColumnExists {
				currency = lineItems[accountCurrencyColumnIdx]

				if _, ok := validators.AllCurrencyNames[currency]; !ok {
					return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			}

			account = e.createNewAccountModel(user.Uid, accountName, currency)
			allNewAccounts = append(allNewAccounts, account)
			accountMap[accountName] = account
		}

		if accountCurrencyColumnExists {
			if account.Currency != lineItems[accountCurrencyColumnIdx] {
				return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
			}
		}

		amount, err := utils.ParseAmount(lineItems[amountColumnIdx])

		if err != nil {
			return nil, nil, nil, nil, err
		}

		relatedAccountId := int64(0)
		relatedAccountAmount := int64(0)

		if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2Name := lineItems[account2ColumnIdx]

			if account2Name == "" {
				return nil, nil, nil, nil, errs.ErrFormatInvalid
			}

			account2, exists := accountMap[account2Name]

			if !exists {
				currency := user.DefaultCurrency

				if accountCurrencyColumnExists {
					currency = lineItems[account2CurrencyColumnIdx]

					if _, ok := validators.AllCurrencyNames[currency]; !ok {
						return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
					}
				}

				account2 = e.createNewAccountModel(user.Uid, account2Name, currency)
				allNewAccounts = append(allNewAccounts, account2)
				accountMap[account2Name] = account2
			}

			if account2CurrencyColumnExists {
				if account2.Currency != lineItems[account2CurrencyColumnIdx] {
					return nil, nil, nil, nil, errs.ErrAccountCurrencyInvalid
				}
			}

			relatedAccountId = account2.AccountId
			relatedAccountAmount, err = utils.ParseAmount(lineItems[amount2ColumnIdx])

			if err != nil {
				return nil, nil, nil, nil, err
			}
		}

		geoLongitude := float64(0)
		geoLatitude := float64(0)

		if geoLocationExists {
			geoLocationItems := strings.Split(lineItems[geoLocationIdx], geoLocationSeparator)

			if len(geoLocationItems) == 2 {
				geoLongitude, err = utils.StringToFloat64(geoLocationItems[0])

				if err != nil {
					return nil, nil, nil, nil, err
				}

				geoLatitude, err = utils.StringToFloat64(geoLocationItems[1])

				if err != nil {
					return nil, nil, nil, nil, err
				}
			}
		}

		if tagsColumnExists {
			tagNames := strings.Split(lineItems[tagsColumnIdx], transactionTagSeparator)

			for i := 0; i < len(tagNames); i++ {
				tagName := tagNames[i]

				if tagName == "" {
					continue
				}

				tag, exists := tagMap[tagName]

				if !exists {
					tag = e.createNewTransactionTagModel(user.Uid, tagName)
					allNewTags = append(allNewTags, tag)
					tagMap[tagName] = tag
				}
			}
		}

		description := ""

		if descriptionColumnExists {
			description = lineItems[descriptionColumnIdx]
		}

		transaction := &models.Transaction{
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
		}

		allNewTransactions = append(allNewTransactions, transaction)
	}

	sort.Sort(allNewTransactions)

	return allNewTransactions, allNewAccounts, allNewSubCategories, allNewTags, nil
}

func (e *EzBookKeepingPlainFileConverter) getTransactionTypeName(transactionDbType models.TransactionDbType) string {
	if transactionDbType == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		return "Balance Modification"
	} else if transactionDbType == models.TRANSACTION_DB_TYPE_INCOME {
		return "Income"
	} else if transactionDbType == models.TRANSACTION_DB_TYPE_EXPENSE {
		return "Expense"
	} else if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return "Transfer"
	} else {
		return ""
	}
}

func (e *EzBookKeepingPlainFileConverter) getTransactionDbType(transactionTypeName string) (models.TransactionDbType, error) {
	if transactionTypeName == "Balance Modification" {
		return models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, nil
	} else if transactionTypeName == "Income" {
		return models.TRANSACTION_DB_TYPE_INCOME, nil
	} else if transactionTypeName == "Expense" {
		return models.TRANSACTION_DB_TYPE_EXPENSE, nil
	} else if transactionTypeName == "Transfer" {
		return models.TRANSACTION_DB_TYPE_TRANSFER_OUT, nil
	} else {
		return 0, errs.ErrTransactionTypeInvalid
	}
}

func (e *EzBookKeepingPlainFileConverter) getTransactionCategoryType(transactionType models.TransactionDbType) (models.TransactionCategoryType, error) {
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

func (e *EzBookKeepingPlainFileConverter) getTransactionCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if !exists {
		return ""
	}

	if category.ParentCategoryId == 0 {
		return category.Name
	}

	parentCategory, exists := categoryMap[category.ParentCategoryId]

	if !exists {
		return ""
	}

	return parentCategory.Name
}

func (e *EzBookKeepingPlainFileConverter) getTransactionSubCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if exists {
		return category.Name
	} else {
		return ""
	}
}

func (e *EzBookKeepingPlainFileConverter) getAccountName(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return account.Name
	} else {
		return ""
	}
}

func (e *EzBookKeepingPlainFileConverter) getAccountCurrency(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return account.Currency
	} else {
		return ""
	}
}

func (e *EzBookKeepingPlainFileConverter) getTags(transactionId int64, allTagIndexes map[int64][]int64, tagMap map[int64]*models.TransactionTag) string {
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
			ret.WriteString(transactionTagSeparator)
		}

		ret.WriteString(strings.Replace(tag.Name, transactionTagSeparator, " ", -1))
	}

	return ret.String()
}

func (e *EzBookKeepingPlainFileConverter) replaceDelimiters(text string, separator string) string {
	text = strings.Replace(text, separator, " ", -1)
	text = strings.Replace(text, "\r\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Replace(text, "\n", " ", -1)

	return text
}

func (e *EzBookKeepingPlainFileConverter) createNewAccountModel(uid int64, accountName string, currency string) *models.Account {
	return &models.Account{
		Uid:      uid,
		Name:     accountName,
		Currency: currency,
	}
}

func (e *EzBookKeepingPlainFileConverter) createNewTransactionCategoryModel(uid int64, categoryName string, transactionCategoryType models.TransactionCategoryType) *models.TransactionCategory {
	return &models.TransactionCategory{
		Uid:  uid,
		Name: categoryName,
		Type: transactionCategoryType,
	}
}

func (e *EzBookKeepingPlainFileConverter) createNewTransactionTagModel(uid int64, tagName string) *models.TransactionTag {
	return &models.TransactionTag{
		Uid:  uid,
		Name: tagName,
	}
}
