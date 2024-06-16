package converters

import (
	"fmt"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// EzBookKeepingPlainFileExporter defines the structure of plain file exporter
type EzBookKeepingPlainFileExporter struct {
}

const lineSeparator = "\n"
const geoLocationSeparator = " "
const transactionTagSeparator = ";"
const headerLine = "Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Geographic Location,Tags,Description" + lineSeparator
const dataLineFormat = "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s" + lineSeparator

// toExportedContent returns the exported plain data
func (e *EzBookKeepingPlainFileExporter) toExportedContent(uid int64, separator string, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexs map[int64][]int64) ([]byte, error) {
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
		amount := e.getDisplayAmount(transaction.Amount)
		account2 := ""
		account2Currency := ""
		account2Amount := ""
		geoLocation := ""

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2 = e.replaceDelimiters(e.getAccountName(transaction.RelatedAccountId, accountMap), separator)
			account2Currency = e.getAccountCurrency(transaction.RelatedAccountId, accountMap)
			account2Amount = e.getDisplayAmount(transaction.RelatedAccountAmount)
		}

		if transaction.GeoLongitude != 0 || transaction.GeoLatitude != 0 {
			geoLocation = fmt.Sprintf("%f%s%f", transaction.GeoLongitude, geoLocationSeparator, transaction.GeoLatitude)
		}

		tags := e.replaceDelimiters(e.getTags(transaction.TransactionId, allTagIndexs, tagMap), separator)
		comment := e.replaceDelimiters(transaction.Comment, separator)

		ret.WriteString(fmt.Sprintf(actualDataLineFormat, transactionTime, transactionTimezone, transactionType, category, subCategory, account, accountCurrency, amount, account2, account2Currency, account2Amount, geoLocation, tags, comment))
	}

	return []byte(ret.String()), nil
}

func (e *EzBookKeepingPlainFileExporter) getTransactionTypeName(transactionDbType models.TransactionDbType) string {
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

func (e *EzBookKeepingPlainFileExporter) getTransactionCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
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

func (e *EzBookKeepingPlainFileExporter) getTransactionSubCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if exists {
		return category.Name
	} else {
		return ""
	}
}

func (e *EzBookKeepingPlainFileExporter) getAccountName(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return account.Name
	} else {
		return ""
	}
}

func (e *EzBookKeepingPlainFileExporter) getAccountCurrency(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return account.Currency
	} else {
		return ""
	}
}

func (e *EzBookKeepingPlainFileExporter) getDisplayAmount(amount int64) string {
	displayAmount := utils.Int64ToString(amount)
	integer := utils.SubString(displayAmount, 0, len(displayAmount)-2)
	decimals := utils.SubString(displayAmount, -2, 2)

	if integer == "" {
		integer = "0"
	} else if integer == "-" {
		integer = "-0"
	}

	if len(decimals) == 0 {
		decimals = "00"
	} else if len(decimals) == 1 {
		decimals = "0" + decimals
	}

	return integer + "." + decimals
}

func (e *EzBookKeepingPlainFileExporter) getTags(transactionId int64, allTagIndexs map[int64][]int64, tagMap map[int64]*models.TransactionTag) string {
	tagIndexs, exists := allTagIndexs[transactionId]

	if !exists {
		return ""
	}

	var ret strings.Builder

	for i := 0; i < len(tagIndexs); i++ {
		if i > 0 {
			ret.WriteString(transactionTagSeparator)
		}

		tagIndex := tagIndexs[i]
		tag, exists := tagMap[tagIndex]

		if !exists {
			continue
		}

		ret.WriteString(tag.Name)
	}

	return ret.String()
}

func (e *EzBookKeepingPlainFileExporter) replaceDelimiters(text string, separator string) string {
	text = strings.Replace(text, separator, " ", -1)
	text = strings.Replace(text, "\r\n", " ", -1)
	text = strings.Replace(text, "\n", " ", -1)

	return text
}
