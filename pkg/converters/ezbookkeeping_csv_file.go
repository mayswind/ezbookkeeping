package converters

import (
	"fmt"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// EzBookKeepingCSVFileExporter defines the structure of csv file exporter
type EzBookKeepingCSVFileExporter struct {
	DataConverter
}

const csvHeaderLine = "Time,Timezone,Type,Category,Sub Category,Account,Account Currency,Amount,Account2,Account2 Currency,Account2 Amount,Tags,Comment\n"
const csvDataLineFormat = "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n"

// ToExportedContent returns the exported csv data
func (e *EzBookKeepingCSVFileExporter) ToExportedContent(uid int64, timezone *time.Location, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexs map[int64][]int64) ([]byte, error) {
	var ret strings.Builder

	ret.Grow(len(transactions) * 100)
	ret.WriteString(csvHeaderLine)

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			continue
		}

		transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
		transactionTime := utils.FormatUnixTimeToLongDateTimeWithoutSecond(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), transactionTimeZone)
		transactionTimezone := utils.FormatTimezoneOffset(transactionTimeZone)
		transactionType := e.getTransactionTypeName(transaction.Type)
		category := e.getTransactionCategoryName(transaction.CategoryId, categoryMap)
		subCategory := e.getTransactionSubCategoryName(transaction.CategoryId, categoryMap)
		account := e.getAccountName(transaction.AccountId, accountMap)
		accountCurrency := e.getAccountCurrency(transaction.AccountId, accountMap)
		amount := e.getDisplayAmount(transaction.Amount)
		account2 := ""
		account2Currency := ""
		account2Amount := ""

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2 = e.getAccountName(transaction.RelatedAccountId, accountMap)
			account2Currency = e.getAccountCurrency(transaction.RelatedAccountId, accountMap)
			account2Amount = e.getDisplayAmount(transaction.RelatedAccountAmount)
		}

		tags := e.getTags(transaction.TransactionId, allTagIndexs, tagMap)
		comment := e.getComment(transaction.Comment)

		ret.WriteString(fmt.Sprintf(csvDataLineFormat, transactionTime, transactionTimezone, transactionType, category, subCategory, account, accountCurrency, amount, account2, account2Currency, account2Amount, tags, comment))
	}

	return []byte(ret.String()), nil
}

func (e *EzBookKeepingCSVFileExporter) getTransactionTypeName(transactionDbType models.TransactionDbType) string {
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

func (e *EzBookKeepingCSVFileExporter) getTransactionCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
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

func (e *EzBookKeepingCSVFileExporter) getTransactionSubCategoryName(categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if exists {
		return category.Name
	} else {
		return ""
	}
}

func (e *EzBookKeepingCSVFileExporter) getAccountName(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return account.Name
	} else {
		return ""
	}
}

func (e *EzBookKeepingCSVFileExporter) getAccountCurrency(accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return account.Currency
	} else {
		return ""
	}
}

func (e *EzBookKeepingCSVFileExporter) getDisplayAmount(amount int64) string {
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

func (e *EzBookKeepingCSVFileExporter) getTags(transactionId int64, allTagIndexs map[int64][]int64, tagMap map[int64]*models.TransactionTag) string {
	tagIndexs, exists := allTagIndexs[transactionId]

	if !exists {
		return ""
	}

	var ret strings.Builder

	for i := 0; i < len(tagIndexs); i++ {
		if i > 0 {
			ret.WriteString(";")
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

func (e *EzBookKeepingCSVFileExporter) getComment(comment string) string {
	comment = strings.Replace(comment, ",", " ", -1)
	comment = strings.Replace(comment, "\r\n", " ", -1)
	comment = strings.Replace(comment, "\n", " ", -1)

	return comment
}
