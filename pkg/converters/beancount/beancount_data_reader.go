package beancount

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const beancountDefaultAssetsAccountTypeName = "Assets"
const beancountDefaultLiabilitiesAccountTypeName = "Liabilities"
const beancountDefaultEquityAccountTypeName = "Equity"
const beancountDefaultIncomeAccountTypeName = "Income"
const beancountDefaultExpenseAccountTypeName = "Expenses"

const beancountOptionAssetsAccountTypeName = "name_assets"
const beancountOptionLiabilitiesAccountTypeName = "name_liabilities"
const beancountOptionEquityAccountTypeName = "name_equity"
const beancountOptionIncomeAccountTypeName = "name_income"
const beancountOptionExpenseAccountTypeName = "name_expenses"

const beancountCommentPrefix = ';'
const beancountAccountNameItemsSeparator = ":"
const beancountMetadataKeySuffix = ':'
const beancountPricePrefix = '@'
const beancountLinkPrefix = '^'
const beancountTagPrefix = '#'

// beancountDataReader defines the structure of Beancount data reader
type beancountDataReader struct {
	accountTypeNameMap         map[string]beancountAccountType
	accountTypeNameReversedMap map[beancountAccountType]string
	allData                    [][]string
}

// read returns the imported Beancount data
// Reference: https://beancount.github.io/docs/beancount_language_syntax.html
func (r *beancountDataReader) read(ctx core.Context) (*beancountData, error) {
	if len(r.allData) < 1 {
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	data := &beancountData{
		Accounts:     make(map[string]*beancountAccount),
		Transactions: make([]*beancountTransactionEntry, 0),
	}

	var err error
	var currentTransactionEntry *beancountTransactionEntry
	var currentTransactionPosting *beancountPosting
	var currentTags []string

	for i := 0; i < len(r.allData); i++ {
		items := r.allData[i]

		if len(items) == 0 || (len(items) == 1 && len(items[0]) == 0) || (len(r.getNotEmptyItemByIndex(items, 0)) > 0 && r.getNotEmptyItemByIndex(items, 0)[0] == beancountCommentPrefix) { // skip empty or comment lines
			continue
		}

		if r.getNotEmptyItemsCount(items) < 2 {
			log.Warnf(ctx, "[beancount_data_reader.read] cannot parse line#%d \"%s\", because not enough items in line", i, strings.Join(items, " "))
			continue
		}

		firstItem := items[0]

		if firstItem == "include" { // not support include directive
			return nil, errs.ErrBeancountFileNotSupportInclude
		} else if firstItem == "plugin" { // skip plugin directive lines
			currentTransactionEntry, currentTransactionPosting = r.updateCurrentState(data, currentTransactionEntry, currentTransactionPosting)
			continue
		} else if firstItem == "option" {
			currentTransactionEntry, currentTransactionPosting = r.updateCurrentState(data, currentTransactionEntry, currentTransactionPosting)
			r.readAndSetOption(ctx, i, items)
			continue
		} else if firstItem == "pushtag" {
			currentTransactionEntry, currentTransactionPosting = r.updateCurrentState(data, currentTransactionEntry, currentTransactionPosting)
			currentTags = r.readAndSetTags(ctx, i, items, currentTags, true)
			continue
		} else if firstItem == "poptag" {
			currentTransactionEntry, currentTransactionPosting = r.updateCurrentState(data, currentTransactionEntry, currentTransactionPosting)
			currentTags = r.readAndSetTags(ctx, i, items, currentTags, false)
			continue
		}

		if len(firstItem) == 0 { // original line has space prefix, maybe transaction posting or metadata line
			actualFirstItem := r.getNotEmptyItemByIndex(items, 0)

			if len(actualFirstItem) == 0 { // skip empty lines
				continue
			}

			if ('A' <= actualFirstItem[0] && actualFirstItem[0] <= 'Z') || actualFirstItem[0] == '!' { // transaction posting
				if currentTransactionEntry != nil && currentTransactionPosting != nil {
					currentTransactionEntry.Postings = append(currentTransactionEntry.Postings, currentTransactionPosting)
					currentTransactionPosting = nil
				}

				currentTransactionPosting, err = r.readTransactionPostingLine(ctx, i, items, data, actualFirstItem[0] == '!')

				if err != nil {
					return nil, err
				}
			} else if 'a' <= actualFirstItem[0] && actualFirstItem[0] <= 'z' { // metadata
				metadata := r.readTransactionMetadataLine(ctx, i, items)

				if metadata == nil {
					continue
				}

				metadataKey := metadata[0]
				metadataValue := metadata[1]

				if currentTransactionPosting != nil {
					if _, exists := currentTransactionPosting.Metadata[metadataKey]; !exists {
						currentTransactionPosting.Metadata[metadataKey] = metadataValue
					}
				} else if currentTransactionEntry != nil {
					if _, exists := currentTransactionEntry.Metadata[metadataKey]; !exists {
						currentTransactionEntry.Metadata[metadataKey] = metadataValue
					}
				}
			} else {
				log.Warnf(ctx, "[beancount_data_reader.read] cannot parse line#%d \"%s\", because line prefix is invalid", i, strings.Join(items, " "))
				currentTransactionEntry, currentTransactionPosting = r.updateCurrentState(data, currentTransactionEntry, currentTransactionPosting)
				continue
			}
		} else if _, err := utils.ParseFromLongDateFirstTime(firstItem, 0); err == nil { // original line has date as first item
			currentTransactionEntry, currentTransactionPosting = r.updateCurrentState(data, currentTransactionEntry, currentTransactionPosting)

			directive := r.getNotEmptyItemByIndex(items, 1)

			if directive == string(beancountDirectiveOpen) ||
				directive == string(beancountDirectiveClose) {
				_, err := r.readAccountLine(ctx, i, items, firstItem, beancountDirective(directive), data)

				if err != nil {
					return nil, err
				}
			} else if directive == string(beancountDirectiveTransaction) ||
				directive == string(beancountDirectiveCompletedTransaction) ||
				directive == string(beancountDirectiveInCompleteTransaction) ||
				directive == string(beancountDirectivePaddingTransaction) {
				currentTransactionEntry = r.readTransactionLine(ctx, i, items, firstItem, beancountDirective(directive), currentTags)
			} else if directive == string(beancountDirectiveCommodity) ||
				directive == string(beancountDirectivePrice) ||
				directive == string(beancountDirectiveNote) ||
				directive == string(beancountDirectiveDocument) ||
				directive == string(beancountDirectiveEvent) ||
				directive == string(beancountDirectiveBalance) ||
				directive == string(beancountDirectivePad) ||
				directive == string(beancountDirectiveQuery) ||
				directive == string(beancountDirectiveCustom) { // skip commodity / price / note / document / event / balance / pad / query / custom lines
				continue
			} else {
				log.Warnf(ctx, "[beancount_data_reader.read] cannot parse line#%d \"%s\", because directive is unknown", i, strings.Join(items, " "))
				continue
			}
		} else { // first item not start with date or space
			currentTransactionEntry, currentTransactionPosting = r.updateCurrentState(data, currentTransactionEntry, currentTransactionPosting)
			continue
		}
	}

	if currentTransactionEntry != nil {
		if currentTransactionPosting != nil {
			currentTransactionEntry.Postings = append(currentTransactionEntry.Postings, currentTransactionPosting)
			currentTransactionPosting = nil
		}

		data.Transactions = append(data.Transactions, currentTransactionEntry)
		currentTransactionEntry = nil
	}

	return data, nil
}

func (r *beancountDataReader) updateCurrentState(data *beancountData, currentTransactionEntry *beancountTransactionEntry, currentTransactionPosting *beancountPosting) (*beancountTransactionEntry, *beancountPosting) {
	if currentTransactionEntry != nil {
		if currentTransactionPosting != nil {
			currentTransactionEntry.Postings = append(currentTransactionEntry.Postings, currentTransactionPosting)
			currentTransactionPosting = nil
		}

		data.Transactions = append(data.Transactions, currentTransactionEntry)
		currentTransactionEntry = nil
		currentTransactionPosting = nil
	}

	return currentTransactionEntry, currentTransactionPosting
}

func (r *beancountDataReader) readAndSetOption(ctx core.Context, lineIndex int, items []string) {
	if r.getNotEmptyItemsCount(items) != 3 {
		log.Warnf(ctx, "[beancount_data_reader.readAndSetOption] cannot parse account type name option line#%d \"%s\", because items count in line not correct", lineIndex, strings.Join(items, " "))
		return
	}

	optionName := r.getNotEmptyItemByIndex(items, 1)
	optionValue := r.getNotEmptyItemByIndex(items, 2)

	switch optionName {
	case beancountOptionAssetsAccountTypeName:
		r.setAccountTypeNameMap(beancountAssetsAccountType, optionValue)
		break
	case beancountOptionLiabilitiesAccountTypeName:
		r.setAccountTypeNameMap(beancountLiabilitiesAccountType, optionValue)
		break
	case beancountOptionEquityAccountTypeName:
		r.setAccountTypeNameMap(beancountEquityAccountType, optionValue)
		break
	case beancountOptionIncomeAccountTypeName:
		r.setAccountTypeNameMap(beancountIncomeAccountType, optionValue)
		break
	case beancountOptionExpenseAccountTypeName:
		r.setAccountTypeNameMap(beancountExpensesAccountType, optionValue)
		break
	default:
		log.Warnf(ctx, "[beancount_data_reader.readAndSetOption] skip option line#%d \"%s\"", lineIndex, strings.Join(items, " "))
		break
	}
}

func (r *beancountDataReader) readAndSetTags(ctx core.Context, lineIndex int, items []string, currentTags []string, pushTag bool) []string {
	if r.getNotEmptyItemsCount(items) != 2 {
		log.Warnf(ctx, "[beancount_data_reader.readAndSetTags] cannot parse push/pop tag line#%d \"%s\", because items count in line not correct", lineIndex, strings.Join(items, " "))
		return currentTags
	}

	tag := r.getNotEmptyItemByIndex(items, 1)

	if len(tag) < 2 || tag[0] != beancountTagPrefix {
		log.Warnf(ctx, "[beancount_data_reader.readAndSetTags] cannot parse push/pop tag line#%d \"%s\", because tag is invalid", lineIndex, strings.Join(items, " "))
		return currentTags
	}

	tag = tag[1:]

	if pushTag {
		for i := 0; i < len(currentTags); i++ {
			if currentTags[i] == tag {
				return currentTags
			}
		}

		return append(currentTags, tag)
	} else { // pop tag
		for i := 0; i < len(currentTags); i++ {
			if currentTags[i] == tag {
				return append(currentTags[:i], currentTags[i+1:]...)
			}
		}

		return currentTags
	}
}

func (r *beancountDataReader) setAccountTypeNameMap(accountType beancountAccountType, accountTypeName string) {
	delete(r.accountTypeNameMap, r.accountTypeNameReversedMap[accountType])
	r.accountTypeNameMap[accountTypeName] = accountType
	r.accountTypeNameReversedMap[accountType] = accountTypeName
}

func (r *beancountDataReader) readAccountLine(ctx core.Context, lineIndex int, items []string, date string, directive beancountDirective, data *beancountData) (*beancountAccount, error) {
	if r.getNotEmptyItemsCount(items) < 3 {
		log.Warnf(ctx, "[beancount_data_reader.parseAccount] cannot parse account line#%d \"%s\", because items count in line not correct", lineIndex, strings.Join(items, " "))
		return nil, nil
	}

	var err error
	accountName := r.getNotEmptyItemByIndex(items, 2)
	account, exists := data.Accounts[accountName]

	if !exists {
		account, err = r.createAccount(ctx, data, accountName)

		if err != nil {
			return nil, err
		}
	}

	if directive == beancountDirectiveOpen {
		account.OpenDate = date
		return account, nil
	} else if directive == beancountDirectiveClose {
		account.CloseDate = date
		return account, nil
	} else {
		log.Warnf(ctx, "[beancount_data_reader.parseAccount] cannot parse account line#%d \"%s\", because directive is invalid", lineIndex, strings.Join(items, " "))
		return nil, nil
	}
}

func (r *beancountDataReader) createAccount(ctx core.Context, data *beancountData, accountName string) (*beancountAccount, error) {
	account := &beancountAccount{
		Name:        accountName,
		AccountType: beancountUnknownAccountType,
	}

	accountNameItems := strings.Split(accountName, beancountAccountNameItemsSeparator)

	if len(accountNameItems) > 1 {
		accountType, exists := r.accountTypeNameMap[accountNameItems[0]]

		if exists {
			account.AccountType = accountType
		} else {
			log.Warnf(ctx, "[beancount_data_reader.createAccount] cannot parse account \"%s\", because account type \"%s\" is invalid", accountName, accountNameItems[0])
			return nil, errs.ErrInvalidBeancountFile
		}
	}

	data.Accounts[accountName] = account
	return account, nil
}

func (r *beancountDataReader) readTransactionLine(ctx core.Context, lineIndex int, items []string, date string, directive beancountDirective, tags []string) *beancountTransactionEntry {
	transactionEntry := &beancountTransactionEntry{
		Date:      date,
		Directive: directive,
		Tags:      make([]string, 0),
		Links:     make([]string, 0),
		Metadata:  make(map[string]string),
	}

	transactionEntry.Tags = append(transactionEntry.Tags, tags...)

	allTags := make(map[string]bool, len(transactionEntry.Tags))

	for _, tag := range transactionEntry.Tags {
		allTags[tag] = true
	}

	// YYYY-MM-DD [txn|Flag] [[Payee] Narration] [#tag] [ˆlink]
	payeeNarrationFirstIndex := 2
	payeeNarrationLastIndex := len(items) - 1

	// parse remain items
	for i := payeeNarrationFirstIndex; i < len(items); i++ {
		item := items[i]

		if len(item) == 0 {
			continue
		}

		if item[0] == beancountCommentPrefix { // ; comment
			if i-1 < payeeNarrationLastIndex {
				payeeNarrationLastIndex = i - 1
			}

			break
		}

		if item[0] == beancountTagPrefix { // [#tag]
			tagName := item[1:]

			if _, exists := allTags[tagName]; !exists {
				transactionEntry.Tags = append(transactionEntry.Tags, tagName)
				allTags[tagName] = true
			}

			if i-1 < payeeNarrationLastIndex {
				payeeNarrationLastIndex = i - 1
			}
		} else if item[0] == beancountLinkPrefix { // [ˆlink]
			transactionEntry.Links = append(transactionEntry.Links, item[1:])

			if i-1 < payeeNarrationLastIndex {
				payeeNarrationLastIndex = i - 1
			}
		}
	}

	if payeeNarrationLastIndex-payeeNarrationFirstIndex >= 1 {
		transactionEntry.Payee = items[payeeNarrationFirstIndex]
		transactionEntry.Narration = items[payeeNarrationFirstIndex+1]
	} else if payeeNarrationLastIndex-payeeNarrationFirstIndex >= 0 {
		transactionEntry.Narration = items[payeeNarrationFirstIndex]
	}

	return transactionEntry
}

func (r *beancountDataReader) readTransactionPostingLine(ctx core.Context, lineIndex int, items []string, data *beancountData, hasFlag bool) (*beancountPosting, error) {
	// [Flag] Account Amount [{Cost}] [@ Price]
	accountNameExpectedIndex := 0

	if hasFlag {
		accountNameExpectedIndex = 1
	}

	if r.getNotEmptyItemsCount(items) <= accountNameExpectedIndex {
		log.Warnf(ctx, "[beancount_data_reader.readTransactionPostingLine] cannot parse transaction posting line#%d \"%s\", because items count in line not correct", lineIndex, strings.Join(items, " "))
		return nil, nil
	}

	accountName, accountNameActualIndex := r.getNotEmptyItemAndIndexByIndex(items, accountNameExpectedIndex)

	if accountName == "" || accountNameActualIndex < 0 {
		log.Warnf(ctx, "[beancount_data_reader.readTransactionPostingLine] cannot parse transaction posting line#%d \"%s\", because missing account name", lineIndex, strings.Join(items, " "))
		return nil, errs.ErrMissingAccountData
	}

	transactionPositing := &beancountPosting{
		Account:  accountName,
		Metadata: make(map[string]string),
	}

	amountActualLastIndex := -1
	transactionPositing.OriginalAmount, amountActualLastIndex = r.getOriginalAmountAndLastIndexFromIndex(items, accountNameActualIndex+1)

	if transactionPositing.OriginalAmount == "" || amountActualLastIndex < 0 {
		log.Warnf(ctx, "[beancount_data_reader.readTransactionPostingLine] cannot parse transaction posting line#%d \"%s\", because missing amount", lineIndex, strings.Join(items, " "))
		return nil, errs.ErrAmountInvalid
	}

	finalAmount, err := evaluateBeancountAmountExpression(ctx, transactionPositing.OriginalAmount)

	if err != nil {
		log.Warnf(ctx, "[beancount_data_reader.readTransactionPostingLine] cannot evaluate amount expression in line#%d \"%s\", because %s", lineIndex, strings.Join(items, " "), err.Error())
		return nil, errs.ErrAmountInvalid
	} else {
		transactionPositing.Amount = finalAmount
	}

	commodityActualIndex := -1
	transactionPositing.Commodity, commodityActualIndex = r.getNotEmptyItemAndIndexFromIndex(items, amountActualLastIndex+1)

	if transactionPositing.Commodity == "" || commodityActualIndex < 0 {
		log.Warnf(ctx, "[beancount_data_reader.readTransactionPostingLine] cannot parse transaction posting line#%d \"%s\", because missing commodity", lineIndex, strings.Join(items, " "))
		return nil, errs.ErrInvalidBeancountFile
	}

	if strings.ToUpper(transactionPositing.Commodity) != transactionPositing.Commodity { // The syntax for a currency is a word all in capital letters
		log.Warnf(ctx, "[beancount_data_reader.readTransactionPostingLine] cannot parse transaction posting line#%d \"%s\", because commodity name is not capital letters", lineIndex, strings.Join(items, " "))
		return nil, errs.ErrInvalidBeancountFile
	}

	// parse remain items
	if commodityActualIndex > 0 {
		for i := commodityActualIndex + 1; i < len(items); i++ {
			item := items[i]

			if len(item) == 0 {
				continue
			}

			if item[0] == beancountCommentPrefix { // ; comment
				break
			}

			if len(item) == 2 && item[0] == beancountPricePrefix && item[1] == beancountPricePrefix { // [@@ TotalCost]
				totalCost, totalCostActualIndex := r.getNotEmptyItemAndIndexFromIndex(items, i+1)

				if totalCostActualIndex > 0 {
					transactionPositing.TotalCost = totalCost
					i = totalCostActualIndex

					totalCostCommodity, totalCostCommodityActualIndex := r.getNotEmptyItemAndIndexFromIndex(items, totalCostActualIndex+1)

					if totalCostCommodityActualIndex > 0 {
						transactionPositing.TotalCostCommodity = totalCostCommodity
						i = totalCostCommodityActualIndex
					}
				}
			} else if len(item) == 1 && item[0] == beancountPricePrefix { // [@ Price]
				price, priceActualIndex := r.getNotEmptyItemAndIndexFromIndex(items, i+1)

				if priceActualIndex > 0 {
					transactionPositing.Price = price
					i = priceActualIndex

					priceCommodity, priceCommodityActualIndex := r.getNotEmptyItemAndIndexFromIndex(items, priceActualIndex+1)

					if priceCommodityActualIndex > 0 {
						transactionPositing.PriceCommodity = priceCommodity
						i = priceCommodityActualIndex
					}
				}
			}
		}
	}

	if transactionPositing.Account != "" {
		_, exists := data.Accounts[transactionPositing.Account]

		if !exists {
			_, err := r.createAccount(ctx, data, transactionPositing.Account)

			if err != nil {
				return nil, err
			}
		}
	}

	return transactionPositing, nil
}

func (r *beancountDataReader) readTransactionMetadataLine(ctx core.Context, lineIndex int, items []string) []string {
	key := r.getNotEmptyItemByIndex(items, 0)
	value := r.getNotEmptyItemByIndex(items, 1)

	if key == "" || value == "" {
		log.Warnf(ctx, "[beancount_data_reader.readTransactionMetadataLine] cannot parse metadata line#%d \"%s\", because key or value is empty", lineIndex, strings.Join(items, " "))
		return nil
	}

	if len(key) == 0 || key[len(key)-1] != beancountMetadataKeySuffix {
		log.Warnf(ctx, "[beancount_data_reader.readTransactionMetadataLine] cannot parse metadata line#%d \"%s\", because key is invalid correct", lineIndex, strings.Join(items, " "))
		return nil
	}

	key = key[:len(key)-1]

	return []string{key, value}
}

func (r *beancountDataReader) getNotEmptyItemByIndex(items []string, index int) string {
	item, _ := r.getNotEmptyItemAndIndexByIndex(items, index)
	return item
}

func (r *beancountDataReader) getNotEmptyItemAndIndexByIndex(items []string, index int) (string, int) {
	count := -1

	for i := 0; i < len(items); i++ {
		item := items[i]

		if len(item) == 0 {
			continue
		}

		count++

		if count == index {
			return items[i], i
		}
	}

	return "", -1
}

func (r *beancountDataReader) getNotEmptyItemAndIndexFromIndex(items []string, startIndex int) (string, int) {
	for i := startIndex; i < len(items); i++ {
		item := items[i]

		if len(item) == 0 {
			continue
		}

		return item, i
	}

	return "", -1
}

func (r *beancountDataReader) getNotEmptyItemsCount(items []string) int {
	count := 0

	for i := 0; i < len(items); i++ {
		if len(items[i]) > 0 {
			count++
		}
	}

	return count
}

func (r *beancountDataReader) getOriginalAmountAndLastIndexFromIndex(items []string, startIndex int) (string, int) {
	amountBuilder := strings.Builder{}
	lastIndex := -1

	for i := startIndex; i < len(items); i++ {
		item := items[i]

		if len(item) == 0 {
			continue
		}

		valid := true

		// The Amount in “Postings” can also be an arithmetic expression using ( ) * / - +
		for j := 0; j < len(item); j++ {
			if !(item[j] >= '0' && item[j] <= '9') && item[j] != '.' && item[j] != '(' && item[j] != ')' &&
				item[j] != '*' && item[j] != '/' && item[j] != '-' && item[j] != '+' {
				valid = false
				break
			}
		}

		if !valid {
			break
		}

		if amountBuilder.Len() > 0 {
			amountBuilder.WriteRune(' ')
		}

		amountBuilder.WriteString(item)
		lastIndex = i
	}

	return amountBuilder.String(), lastIndex
}

func createNewBeancountDataReader(ctx core.Context, data []byte) (*beancountDataReader, error) {
	fallback := unicode.UTF8.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(fallback))
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ' '
	csvReader.FieldsPerRecord = -1

	allData := make([][]string, 0)

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[beancount_data_reader.createNewBeancountDataReader] cannot parse data, because %s", err.Error())
			return nil, errs.ErrInvalidBeancountFile
		}

		allData = append(allData, items)
	}

	return &beancountDataReader{
		accountTypeNameMap: map[string]beancountAccountType{
			beancountDefaultAssetsAccountTypeName:      beancountAssetsAccountType,
			beancountDefaultLiabilitiesAccountTypeName: beancountLiabilitiesAccountType,
			beancountDefaultEquityAccountTypeName:      beancountEquityAccountType,
			beancountDefaultIncomeAccountTypeName:      beancountIncomeAccountType,
			beancountDefaultExpenseAccountTypeName:     beancountExpensesAccountType,
		},
		accountTypeNameReversedMap: map[beancountAccountType]string{
			beancountAssetsAccountType:      beancountDefaultAssetsAccountTypeName,
			beancountLiabilitiesAccountType: beancountDefaultLiabilitiesAccountTypeName,
			beancountEquityAccountType:      beancountDefaultEquityAccountTypeName,
			beancountIncomeAccountType:      beancountDefaultIncomeAccountTypeName,
			beancountExpensesAccountType:    beancountDefaultExpenseAccountTypeName,
		},
		allData: allData,
	}, nil
}
