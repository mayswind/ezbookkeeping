package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"xorm.io/builder"
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

const pageCountForLoadTransactionAmounts = 1000

// TransactionService represents transaction service
type TransactionService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a transaction service singleton instance
var (
	Transactions = &TransactionService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetTotalTransactionCountByUid returns total transaction count of user
func (s *TransactionService) GetTotalTransactionCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.Transaction{})

	return count, err
}

// GetAllTransactions returns all transactions
func (s *TransactionService) GetAllTransactions(c core.Context, uid int64, pageCount int32, noDuplicated bool) ([]*models.Transaction, error) {
	maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(time.Now().Unix())
	var allTransactions []*models.Transaction

	for maxTransactionTime > 0 {
		transactions, err := s.GetAllTransactionsByMaxTime(c, uid, maxTransactionTime, pageCount, noDuplicated)

		if err != nil {
			return nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < int(pageCount) {
			maxTransactionTime = 0
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	return allTransactions, nil
}

// GetAllTransactionsByMaxTime returns all transactions before given time
func (s *TransactionService) GetAllTransactionsByMaxTime(c core.Context, uid int64, maxTransactionTime int64, count int32, noDuplicated bool) ([]*models.Transaction, error) {
	return s.GetTransactionsByMaxTime(c, uid, maxTransactionTime, 0, 0, nil, nil, nil, false, models.TRANSACTION_TAG_FILTER_HAS_ANY, "", "", 1, count, false, noDuplicated)
}

// GetTransactionsByMaxTime returns transactions before given time
func (s *TransactionService) GetTransactionsByMaxTime(c core.Context, uid int64, maxTransactionTime int64, minTransactionTime int64, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, tagIds []int64, noTags bool, tagFilterType models.TransactionTagFilterType, amountFilter string, keyword string, page int32, count int32, needOneMoreItem bool, noDuplicated bool) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if page < 0 {
		return nil, errs.ErrPageIndexInvalid
	} else if page == 0 {
		page = 1
	}

	if count < 1 {
		return nil, errs.ErrPageCountInvalid
	}

	var transactions []*models.Transaction
	var err error

	actualCount := count

	if needOneMoreItem {
		actualCount++
	}

	condition, conditionParams := s.buildTransactionQueryCondition(uid, maxTransactionTime, minTransactionTime, transactionType, categoryIds, accountIds, tagIds, amountFilter, keyword, noDuplicated)
	sess := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...)
	sess = s.appendFilterTagIdsConditionToQuery(sess, uid, maxTransactionTime, minTransactionTime, tagIds, noTags, tagFilterType)

	err = sess.Limit(int(actualCount), int(count*(page-1))).OrderBy("transaction_time desc").Find(&transactions)

	return transactions, err
}

// GetTransactionsInMonthByPage returns all transactions in given year and month
func (s *TransactionService) GetTransactionsInMonthByPage(c core.Context, uid int64, year int32, month int32, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, tagIds []int64, noTags bool, tagFilterType models.TransactionTagFilterType, amountFilter string, keyword string) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	minTransactionTime, maxTransactionTime, err := utils.GetTransactionTimeRangeByYearMonth(year, month)

	if err != nil {
		return nil, errs.ErrSystemError
	}

	var transactions []*models.Transaction

	condition, conditionParams := s.buildTransactionQueryCondition(uid, maxTransactionTime, minTransactionTime, transactionType, categoryIds, accountIds, tagIds, amountFilter, keyword, true)
	sess := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...)
	sess = s.appendFilterTagIdsConditionToQuery(sess, uid, maxTransactionTime, minTransactionTime, tagIds, noTags, tagFilterType)

	err = sess.OrderBy("transaction_time desc").Find(&transactions)

	transactionsInMonth := make([]*models.Transaction, 0, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		transactionUnixTime := utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime)
		transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)

		if utils.IsUnixTimeEqualsYearAndMonth(transactionUnixTime, transactionTimeZone, year, month) {
			transactionsInMonth = append(transactionsInMonth, transaction)
		}
	}

	return transactionsInMonth, err
}

// GetTransactionByTransactionId returns a transaction model according to transaction id
func (s *TransactionService) GetTransactionByTransactionId(c core.Context, uid int64, transactionId int64) (*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if transactionId <= 0 {
		return nil, errs.ErrTransactionIdInvalid
	}

	transaction := &models.Transaction{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(transactionId).Where("uid=? AND deleted=?", uid, false).Get(transaction)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionNotFound
	}

	return transaction, nil
}

// GetAllTransactionCount returns total count of transactions
func (s *TransactionService) GetAllTransactionCount(c core.Context, uid int64) (int64, error) {
	return s.GetTransactionCount(c, uid, 0, 0, 0, nil, nil, nil, false, models.TRANSACTION_TAG_FILTER_HAS_ANY, "", "")
}

// GetTransactionCount returns count of transactions
func (s *TransactionService) GetTransactionCount(c core.Context, uid int64, maxTransactionTime int64, minTransactionTime int64, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, tagIds []int64, noTags bool, tagFilterType models.TransactionTagFilterType, amountFilter string, keyword string) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	condition, conditionParams := s.buildTransactionQueryCondition(uid, maxTransactionTime, minTransactionTime, transactionType, categoryIds, accountIds, tagIds, amountFilter, keyword, true)
	sess := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...)
	sess = s.appendFilterTagIdsConditionToQuery(sess, uid, maxTransactionTime, minTransactionTime, tagIds, noTags, tagFilterType)

	return sess.Count(&models.Transaction{})
}

// CreateTransaction saves a new transaction to database
func (s *TransactionService) CreateTransaction(c core.Context, transaction *models.Transaction, tagIds []int64, pictureIds []int64) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	// Check whether account id is valid
	err := s.isAccountIdValid(transaction)

	if err != nil {
		return err
	}

	now := time.Now().Unix()

	needTransactionUuidCount := 1

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		needTransactionUuidCount = 2
	}

	transactionUuids := s.GenerateUuids(uuid.UUID_TYPE_TRANSACTION, uint16(needTransactionUuidCount))

	if len(transactionUuids) < needTransactionUuidCount {
		return errs.ErrSystemIsBusy
	}

	tagIds = utils.ToUniqueInt64Slice(tagIds)
	needTagIndexUuidCount := uint16(len(tagIds))
	tagIndexUuids := s.GenerateUuids(uuid.UUID_TYPE_TAG_INDEX, needTagIndexUuidCount)

	if len(tagIndexUuids) < int(needTagIndexUuidCount) {
		return errs.ErrSystemIsBusy
	}

	transaction.TransactionId = transactionUuids[0]

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		transaction.RelatedId = transactionUuids[1]
	}

	transaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

	transaction.CreatedUnixTime = now
	transaction.UpdatedUnixTime = now

	transactionTagIndexes := make([]*models.TransactionTagIndex, len(tagIds))

	for i := 0; i < len(tagIds); i++ {
		transactionTagIndexes[i] = &models.TransactionTagIndex{
			TagIndexId:      tagIndexUuids[i],
			Uid:             transaction.Uid,
			Deleted:         false,
			TagId:           tagIds[i],
			TransactionId:   transaction.TransactionId,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		}
	}

	pictureUpdateModel := &models.TransactionPictureInfo{
		TransactionId:   transaction.TransactionId,
		UpdatedUnixTime: now,
	}

	userDataDb := s.UserDataDB(transaction.Uid)

	return userDataDb.DoTransaction(c, func(sess *xorm.Session) error {
		return s.doCreateTransaction(c, userDataDb, sess, transaction, transactionTagIndexes, tagIds, pictureIds, pictureUpdateModel)
	})
}

// BatchCreateTransactions saves new transactions to database
func (s *TransactionService) BatchCreateTransactions(c core.Context, uid int64, transactions []*models.Transaction, allTagIds map[int][]int64, processHandler core.TaskProcessUpdateHandler) error {
	now := time.Now().Unix()
	currentProcess := float64(0)
	processUpdateStep := int(math.Max(100.0, float64(len(transactions)/100.0)))

	needTransactionUuidCount := uint16(0)
	needTagIndexUuidCount := uint16(0)

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Uid != uid {
			return errs.ErrUserIdInvalid
		}

		// Check whether account id is valid
		err := s.isAccountIdValid(transaction)

		if err != nil {
			return err
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			needTransactionUuidCount += 2
		} else {
			needTransactionUuidCount++
		}

		transaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

		transaction.CreatedUnixTime = now
		transaction.UpdatedUnixTime = now
	}

	for index, tagIds := range allTagIds {
		if index < 0 || index >= len(transactions) {
			return errs.ErrOperationFailed
		}

		uniqueTagIds := utils.ToUniqueInt64Slice(tagIds)
		needTagIndexUuidCount += uint16(len(uniqueTagIds))
	}

	if needTransactionUuidCount > uint16(65535) || needTagIndexUuidCount > uint16(65535) {
		return errs.ErrImportTooManyTransaction
	}

	transactionUuids := s.GenerateUuids(uuid.UUID_TYPE_TRANSACTION, needTransactionUuidCount)
	transactionUuidIndex := 0

	if len(transactionUuids) < int(needTransactionUuidCount) {
		return errs.ErrSystemIsBusy
	}

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		transaction.TransactionId = transactionUuids[transactionUuidIndex]
		transactionUuidIndex++

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			transaction.RelatedId = transactionUuids[transactionUuidIndex]
			transactionUuidIndex++
		}
	}

	tagIndexUuids := s.GenerateUuids(uuid.UUID_TYPE_TAG_INDEX, needTagIndexUuidCount)
	tagIndexUuidIndex := 0

	if len(tagIndexUuids) < int(needTagIndexUuidCount) {
		return errs.ErrSystemIsBusy
	}

	allTransactionTagIndexes := make(map[int64][]*models.TransactionTagIndex)
	allTransactionTagIds := make(map[int64][]int64)

	for index, tagIds := range allTagIds {
		transaction := transactions[index]
		uniqueTagIds := utils.ToUniqueInt64Slice(tagIds)

		transactionTagIndexes := make([]*models.TransactionTagIndex, len(uniqueTagIds))

		for i := 0; i < len(uniqueTagIds); i++ {
			transactionTagIndexes[i] = &models.TransactionTagIndex{
				TagIndexId:      tagIndexUuids[tagIndexUuidIndex],
				Uid:             transaction.Uid,
				Deleted:         false,
				TagId:           uniqueTagIds[i],
				TransactionId:   transaction.TransactionId,
				CreatedUnixTime: now,
				UpdatedUnixTime: now,
			}

			tagIndexUuidIndex++
		}

		allTransactionTagIndexes[transaction.TransactionId] = transactionTagIndexes
		allTransactionTagIds[transaction.TransactionId] = uniqueTagIds
	}

	userDataDb := s.UserDataDB(uid)

	return userDataDb.DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(transactions); i++ {
			transaction := transactions[i]
			transactionTagIndexes := allTransactionTagIndexes[transaction.TransactionId]
			transactionTagIds := allTransactionTagIds[transaction.TransactionId]
			err := s.doCreateTransaction(c, userDataDb, sess, transaction, transactionTagIndexes, transactionTagIds, nil, nil)

			currentProcess = float64(i) / float64(len(transactions)) * 100

			if processHandler != nil && i%processUpdateStep == 0 {
				processHandler(currentProcess)
			}

			if err != nil {
				transactionUnixTime := utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime)
				transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
				log.Errorf(c, "[transactions.BatchCreateTransactions] failed to create trasaction (datetime: %s, type: %s, amount: %d)", utils.FormatUnixTimeToLongDateTime(transactionUnixTime, transactionTimeZone), transaction.Type, transaction.Amount)
				return err
			}
		}

		return nil
	})
}

// CreateScheduledTransactions saves all scheduled transactions that should be created now
func (s *TransactionService) CreateScheduledTransactions(c core.Context, currentUnixTime int64, interval time.Duration) error {
	var allTemplates []*models.TransactionTemplate
	intervalMinute := int(interval / time.Minute)
	currentTime := time.Unix(currentUnixTime, 0)
	currentMinute := (currentTime.Minute() / intervalMinute) * intervalMinute

	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentMinute, 0, 0, time.Local)
	startTimeInUTC := startTime.In(time.UTC)

	minutesElapsedOfDayInUtc := startTimeInUTC.Hour()*60 + startTimeInUTC.Minute()
	secondsElapsedOfDayInUtc := minutesElapsedOfDayInUtc * 60
	todayFirstTimeInUTC := startTimeInUTC.Add(time.Duration(-secondsElapsedOfDayInUtc) * time.Second)
	todayFirstUnixTimeInUTC := todayFirstTimeInUTC.Unix()

	minScheduledAt := minutesElapsedOfDayInUtc
	maxScheduledAt := minScheduledAt + intervalMinute

	for i := 0; i < s.UserDataDBCount(); i++ {
		var templates []*models.TransactionTemplate
		err := s.UserDataDBByIndex(i).NewSession(c).Where("deleted=? AND template_type=? AND (scheduled_frequency_type=? OR scheduled_frequency_type=?) AND (scheduled_start_time IS NULL OR scheduled_start_time<=?) AND (scheduled_end_time IS NULL OR scheduled_end_time>=?) AND scheduled_at>=? AND scheduled_at<?", false, models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, startTime.Unix(), startTime.Unix(), minScheduledAt, maxScheduledAt).Find(&templates)

		if err != nil {
			return err
		}

		allTemplates = append(allTemplates, templates...)
	}

	if len(allTemplates) < 1 {
		return nil
	}

	log.Infof(c, "[transactions.CreateScheduledTransactions] should process %d scheduled transaction templates now (scheduled at from %d to %d)", len(allTemplates), minScheduledAt, maxScheduledAt)

	successCount := 0
	skipCount := 0
	failedCount := 0

	for i := 0; i < len(allTemplates); i++ {
		template := allTemplates[i]

		if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED {
			skipCount++
			log.Warnf(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" disabled scheduled transaction frequency", template.TemplateId)
			continue
		}

		if (template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY &&
			template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY) ||
			template.ScheduledFrequency == "" {
			skipCount++
			log.Warnf(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" has invalid scheduled transaction frequency", template.TemplateId)
			continue
		}

		frequencyValues, err := utils.StringArrayToInt64Array(strings.Split(template.ScheduledFrequency, ","))

		if err != nil {
			skipCount++
			log.Warnf(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" has invalid scheduled transaction frequency, because %s", template.TemplateId, err.Error())
			continue
		}

		frequencyValueSet := utils.ToSet(frequencyValues)
		templateTimeZone := time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
		transactionUnixTime := todayFirstUnixTimeInUTC + int64(template.ScheduledAt)*60
		transactionTime := time.Unix(transactionUnixTime, 0).In(templateTimeZone)

		if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY && !frequencyValueSet[int64(transactionTime.Weekday())] {
			skipCount++
			log.Infof(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" does not need to create transaction, today is %s", template.TemplateId, startTimeInUTC.Weekday())
			continue
		} else if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY && !frequencyValueSet[int64(transactionTime.Day())] {
			skipCount++
			log.Infof(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" does not need to create transaction, today is %d of month", template.TemplateId, startTimeInUTC.Day())
			continue
		}

		if template.ScheduledStartTime != nil && *template.ScheduledStartTime > transactionUnixTime {
			skipCount++
			log.Infof(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" does not need to create transaction, now is earlier than the start time %d", template.TemplateId, *template.ScheduledStartTime)
			continue
		}

		if template.ScheduledEndTime != nil && *template.ScheduledEndTime < transactionUnixTime {
			skipCount++
			log.Infof(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" does not need to create transaction, now is later than the end time %d", template.TemplateId, *template.ScheduledEndTime)
			continue
		}

		var transactionDbType models.TransactionDbType

		if template.Type == models.TRANSACTION_TYPE_EXPENSE {
			transactionDbType = models.TRANSACTION_DB_TYPE_EXPENSE
		} else if template.Type == models.TRANSACTION_TYPE_INCOME {
			transactionDbType = models.TRANSACTION_DB_TYPE_INCOME
		} else if template.Type == models.TRANSACTION_TYPE_TRANSFER {
			transactionDbType = models.TRANSACTION_DB_TYPE_TRANSFER_OUT
		} else {
			skipCount++
			log.Warnf(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" has invalid transaction type", template.TemplateId)
			continue
		}

		transaction := &models.Transaction{
			Uid:               template.Uid,
			Type:              transactionDbType,
			CategoryId:        template.CategoryId,
			TransactionTime:   utils.GetMinTransactionTimeFromUnixTime(transactionTime.Unix()),
			TimezoneUtcOffset: template.ScheduledTimezoneUtcOffset,
			AccountId:         template.AccountId,
			Amount:            template.Amount,
			HideAmount:        template.HideAmount,
			Comment:           template.Comment,
			CreatedIp:         "127.0.0.1",
			ScheduledCreated:  true,
		}

		if template.Type == models.TRANSACTION_TYPE_TRANSFER {
			transaction.RelatedAccountId = template.RelatedAccountId
			transaction.RelatedAccountAmount = template.RelatedAccountAmount
		}

		tagIds := template.GetTagIds()
		err = s.CreateTransaction(c, transaction, tagIds, nil)

		if err == nil {
			successCount++
			log.Infof(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" has created a new trasaction \"id:%d\"", template.TemplateId, transaction.TransactionId)
		} else {
			failedCount++
			log.Errorf(c, "[transactions.CreateScheduledTransactions] transaction template \"id:%d\" failed to create new trasaction", template.TemplateId)
		}
	}

	log.Infof(c, "[transactions.CreateScheduledTransactions] %d transactions has been created successfully, %d templates does not need to create transactions and %d transactions failed to create", successCount, skipCount, failedCount)

	return nil
}

// ModifyTransaction saves an existed transaction to database
func (s *TransactionService) ModifyTransaction(c core.Context, transaction *models.Transaction, currentTagIdsCount int, addTagIds []int64, removeTagIds []int64, addPictureIds []int64, removePictureIds []int64) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	needTagIndexUuidCount := uint16(len(addTagIds))
	tagIndexUuids := s.GenerateUuids(uuid.UUID_TYPE_TAG_INDEX, needTagIndexUuidCount)

	if len(tagIndexUuids) < int(needTagIndexUuidCount) {
		return errs.ErrSystemIsBusy
	}

	updateCols := make([]string, 0, 16)

	now := time.Now().Unix()

	transaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))
	transaction.UpdatedUnixTime = now
	updateCols = append(updateCols, "updated_unix_time")

	addTagIds = utils.ToUniqueInt64Slice(addTagIds)
	removeTagIds = utils.ToUniqueInt64Slice(removeTagIds)

	transactionTagIndexes := make([]*models.TransactionTagIndex, len(addTagIds))

	for i := 0; i < len(addTagIds); i++ {
		transactionTagIndexes[i] = &models.TransactionTagIndex{
			TagIndexId:      tagIndexUuids[i],
			Uid:             transaction.Uid,
			Deleted:         false,
			TagId:           addTagIds[i],
			TransactionId:   transaction.TransactionId,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		}
	}

	err := s.UserDataDB(transaction.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Get and verify current transaction
		oldTransaction := &models.Transaction{}
		has, err := sess.ID(transaction.TransactionId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(oldTransaction)

		if err != nil {
			log.Errorf(c, "[transactions.ModifyTransaction] failed to get current transaction, because %s", err.Error())
			return err
		} else if !has {
			return errs.ErrTransactionNotFound
		}

		transaction.Type = oldTransaction.Type

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			transaction.RelatedId = oldTransaction.RelatedId
		}

		// Check whether account id is valid
		err = s.isAccountIdValid(transaction)

		if err != nil {
			return err
		}

		// Get and verify source and destination account (if necessary)
		sourceAccount, destinationAccount, err := s.getAccountModels(sess, transaction)

		if err != nil {
			log.Errorf(c, "[transactions.ModifyTransaction] failed to get account, because %s", err.Error())
			return err
		}

		if sourceAccount.Hidden || (destinationAccount != nil && destinationAccount.Hidden) {
			return errs.ErrCannotModifyTransactionInHiddenAccount
		}

		if sourceAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS || (destinationAccount != nil && destinationAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS) {
			return errs.ErrCannotModifyTransactionInParentAccount
		}

		if (transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) &&
			sourceAccount.Currency == destinationAccount.Currency && transaction.Amount != transaction.RelatedAccountAmount {
			return errs.ErrTransactionSourceAndDestinationAmountNotEqual
		}

		if (transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) &&
			(transaction.Amount < 0 || transaction.RelatedAccountAmount < 0) {
			return errs.ErrTransferTransactionAmountCannotBeLessThanZero
		}

		oldSourceAccount, oldDestinationAccount, err := s.getOldAccountModels(sess, transaction, oldTransaction, sourceAccount, destinationAccount)

		if err != nil {
			log.Errorf(c, "[transactions.ModifyTransaction] failed to get old account, because %s", err.Error())
			return err
		}

		if oldSourceAccount.Hidden || (oldDestinationAccount != nil && oldDestinationAccount.Hidden) {
			return errs.ErrCannotAddTransactionToHiddenAccount
		}

		// Append modified columns and verify
		if transaction.CategoryId != oldTransaction.CategoryId {
			// Get and verify category
			err = s.isCategoryValid(sess, transaction)

			if err != nil {
				return err
			}

			updateCols = append(updateCols, "category_id")
		}

		modifyTransactionTime := false

		if utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime) != utils.GetUnixTimeFromTransactionTime(oldTransaction.TransactionTime) {
			if oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
				return errs.ErrBalanceModificationTransactionCannotModifyTime
			}

			sameSecondLatestTransaction := &models.Transaction{}
			minTransactionTime := utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))
			maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

			has, err = sess.Where("uid=? AND deleted=? AND transaction_time>=? AND transaction_time<=?", transaction.Uid, false, minTransactionTime, maxTransactionTime).OrderBy("transaction_time desc").Limit(1).Get(sameSecondLatestTransaction)

			if err != nil {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to get trasaction time, because %s", err.Error())
				return err
			}

			if has && sameSecondLatestTransaction.TransactionTime < maxTransactionTime-1 {
				transaction.TransactionTime = sameSecondLatestTransaction.TransactionTime + 1
			} else if has && sameSecondLatestTransaction.TransactionTime == maxTransactionTime-1 {
				return errs.ErrTooMuchTransactionInOneSecond
			}

			updateCols = append(updateCols, "transaction_time")
			modifyTransactionTime = true
		}

		if transaction.TimezoneUtcOffset != oldTransaction.TimezoneUtcOffset {
			updateCols = append(updateCols, "timezone_utc_offset")
		}

		if transaction.AccountId != oldTransaction.AccountId {
			updateCols = append(updateCols, "account_id")
		}

		if transaction.Amount != oldTransaction.Amount {
			if oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
				transaction.RelatedAccountAmount = oldTransaction.RelatedAccountAmount + transaction.Amount - oldTransaction.Amount
				updateCols = append(updateCols, "related_account_amount")
			}

			updateCols = append(updateCols, "amount")
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			if transaction.RelatedAccountId != oldTransaction.RelatedAccountId {
				updateCols = append(updateCols, "related_account_id")
			}

			if transaction.RelatedAccountAmount != oldTransaction.RelatedAccountAmount {
				updateCols = append(updateCols, "related_account_amount")
			}
		}

		if transaction.HideAmount != oldTransaction.HideAmount {
			updateCols = append(updateCols, "hide_amount")
		}

		if transaction.Comment != oldTransaction.Comment {
			updateCols = append(updateCols, "comment")
		}

		if transaction.GeoLongitude != oldTransaction.GeoLongitude {
			updateCols = append(updateCols, "geo_longitude")
		}

		if transaction.GeoLatitude != oldTransaction.GeoLatitude {
			updateCols = append(updateCols, "geo_latitude")
		}

		// Get and verify tags
		err = s.isTagsValid(sess, transaction, transactionTagIndexes, addTagIds)

		if err != nil {
			return err
		}

		// Get and verify pictures
		err = s.isPicturesValid(sess, transaction, addPictureIds)

		if err != nil {
			return err
		}

		// Not allow to add transaction before balance modification transaction
		if transaction.Type != models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			otherTransactionExists := false

			if destinationAccount != nil && sourceAccount.AccountId != destinationAccount.AccountId {
				otherTransactionExists, err = sess.Cols("uid", "deleted", "account_id").Where("uid=? AND deleted=? AND type=? AND (account_id=? OR account_id=?) AND transaction_time>=?", transaction.Uid, false, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, sourceAccount.AccountId, destinationAccount.AccountId, transaction.TransactionTime).Limit(1).Exist(&models.Transaction{})
			} else {
				otherTransactionExists, err = sess.Cols("uid", "deleted", "account_id").Where("uid=? AND deleted=? AND type=? AND account_id=? AND transaction_time>=?", transaction.Uid, false, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, sourceAccount.AccountId, transaction.TransactionTime).Limit(1).Exist(&models.Transaction{})
			}

			if err != nil {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to get whether other transactions exist, because %s", err.Error())
				return err
			} else if otherTransactionExists {
				return errs.ErrCannotAddTransactionBeforeBalanceModificationTransaction
			}
		}

		// Update transaction row
		updatedRows, err := sess.ID(transaction.TransactionId).Cols(updateCols...).Where("uid=? AND deleted=?", transaction.Uid, false).Update(transaction)

		if err != nil {
			log.Errorf(c, "[transactions.ModifyTransaction] failed to update transaction, because %s", err.Error())
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionNotFound
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			relatedTransaction := s.GetRelatedTransferTransaction(transaction)

			if utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime) != utils.GetUnixTimeFromTransactionTime(relatedTransaction.TransactionTime) {
				return errs.ErrTooMuchTransactionInOneSecond
			}

			relatedUpdateCols := s.getRelatedUpdateColumns(updateCols)
			updatedRows, err := sess.ID(relatedTransaction.TransactionId).Cols(relatedUpdateCols...).Where("uid=? AND deleted=?", relatedTransaction.Uid, false).Update(relatedTransaction)

			if err != nil {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to update related transaction, because %s", err.Error())
				return err
			} else if updatedRows < 1 {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to update related transaction")
				return errs.ErrDatabaseOperationFailed
			}
		}

		// Update transaction tag index
		if len(removeTagIds) > 0 {
			tagIndexUpdateModel := &models.TransactionTagIndex{
				Deleted:         true,
				DeletedUnixTime: now,
			}

			deletedRows, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", transaction.Uid, false, transaction.TransactionId).In("tag_id", removeTagIds).Update(tagIndexUpdateModel)

			if err != nil {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to remove old transaction tag index, because %s", err.Error())
				return err
			} else if deletedRows < 1 {
				return errs.ErrTransactionTagNotFound
			}
		}

		if len(transactionTagIndexes) > 0 {
			for i := 0; i < len(transactionTagIndexes); i++ {
				transactionTagIndex := transactionTagIndexes[i]
				transactionTagIndex.TransactionTime = transaction.TransactionTime

				_, err := sess.Insert(transactionTagIndex)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to add new transaction tag index, because %s", err.Error())
					return err
				}
			}
		} else if len(transactionTagIndexes) == 0 && currentTagIdsCount > 0 && modifyTransactionTime {
			tagIndexUpdateModel := &models.TransactionTagIndex{
				TransactionTime: transaction.TransactionTime,
			}

			_, err := sess.Where("uid=? AND deleted=? AND transaction_id=?", transaction.Uid, false, transaction.TransactionId).Update(tagIndexUpdateModel)

			if err != nil {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to update transaction tag index, because %s", err.Error())
				return err
			}
		}

		// Update transaction picture
		if len(removePictureIds) > 0 {
			pictureUpdateModel := &models.TransactionPictureInfo{
				Deleted:         true,
				DeletedUnixTime: now,
			}

			deletedRows, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", transaction.Uid, false, transaction.TransactionId).In("picture_id", removePictureIds).Update(pictureUpdateModel)

			if err != nil {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to remove old transaction picture info, because %s", err.Error())
				return err
			} else if deletedRows < 1 {
				return errs.ErrTransactionPictureNotFound
			}
		}

		if len(addPictureIds) > 0 {
			pictureUpdateModel := &models.TransactionPictureInfo{
				TransactionId:   transaction.TransactionId,
				UpdatedUnixTime: now,
			}

			_, err = sess.Cols("transaction_id", "updated_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", transaction.Uid, false, models.TransactionPictureNewPictureTransactionId).In("picture_id", addPictureIds).Update(pictureUpdateModel)

			if err != nil {
				log.Errorf(c, "[transactions.ModifyTransaction] failed to update new transaction picture info, because %s", err.Error())
				return err
			}
		}

		// Update account table
		if oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			if transaction.AccountId != oldTransaction.AccountId {
				return errs.ErrBalanceModificationTransactionCannotChangeAccountId
			}

			if transaction.RelatedAccountAmount != oldTransaction.RelatedAccountAmount {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)+(%d)", oldTransaction.RelatedAccountAmount, transaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			var oldAccountNewAmount int64 = 0
			var newAccountNewAmount int64 = 0

			if transaction.AccountId == oldTransaction.AccountId {
				oldAccountNewAmount = transaction.Amount
			} else if transaction.AccountId != oldTransaction.AccountId {
				newAccountNewAmount = transaction.Amount
			}

			if oldAccountNewAmount != oldTransaction.Amount {
				oldSourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldSourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)+(%d)", oldTransaction.Amount, oldAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldSourceAccount.Uid, false).Update(oldSourceAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newAccountNewAmount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", newAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			var oldAccountNewAmount int64 = 0
			var newAccountNewAmount int64 = 0

			if transaction.AccountId == oldTransaction.AccountId {
				oldAccountNewAmount = transaction.Amount
			} else if transaction.AccountId != oldTransaction.AccountId {
				newAccountNewAmount = transaction.Amount
			}

			if oldAccountNewAmount != oldTransaction.Amount {
				oldSourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldSourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)-(%d)", oldTransaction.Amount, oldAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldSourceAccount.Uid, false).Update(oldSourceAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newAccountNewAmount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", newAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			var oldSourceAccountNewAmount int64 = 0
			var newSourceAccountNewAmount int64 = 0

			if transaction.AccountId == oldTransaction.AccountId {
				oldSourceAccountNewAmount = transaction.Amount
			} else if transaction.AccountId != oldTransaction.AccountId {
				newSourceAccountNewAmount = transaction.Amount
			}

			if oldSourceAccountNewAmount != oldTransaction.Amount {
				oldSourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldSourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)-(%d)", oldTransaction.Amount, oldSourceAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldSourceAccount.Uid, false).Update(oldSourceAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newSourceAccountNewAmount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", newSourceAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}

			var oldDestinationAccountNewAmount int64 = 0
			var newDestinationAccountNewAmount int64 = 0

			if transaction.RelatedAccountId == oldTransaction.RelatedAccountId {
				oldDestinationAccountNewAmount = transaction.RelatedAccountAmount
			} else if transaction.RelatedAccountId != oldTransaction.RelatedAccountId {
				newDestinationAccountNewAmount = transaction.RelatedAccountAmount
			}

			if oldDestinationAccountNewAmount != oldTransaction.RelatedAccountAmount {
				oldDestinationAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldDestinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)+(%d)", oldTransaction.RelatedAccountAmount, oldDestinationAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldDestinationAccount.Uid, false).Update(oldDestinationAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newDestinationAccountNewAmount != 0 {
				destinationAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", newDestinationAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

				if err != nil {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance, because %s", err.Error())
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.ModifyTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			return errs.ErrTransactionTypeInvalid
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// DeleteTransaction deletes an existed transaction from database
func (s *TransactionService) DeleteTransaction(c core.Context, uid int64, transactionId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Transaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	tagIndexUpdateModel := &models.TransactionTagIndex{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	pictureUpdateModel := &models.TransactionPictureInfo{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Get and verify current transaction
		oldTransaction := &models.Transaction{}
		has, err := sess.ID(transactionId).Where("uid=? AND deleted=?", uid, false).Get(oldTransaction)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrTransactionNotFound
		}

		// Get and verify source and destination account
		sourceAccount, destinationAccount, err := s.getAccountModels(sess, oldTransaction)

		if err != nil {
			return err
		}

		if sourceAccount.Hidden || (destinationAccount != nil && destinationAccount.Hidden) {
			return errs.ErrCannotDeleteTransactionInHiddenAccount
		}

		if sourceAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS || (destinationAccount != nil && destinationAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS) {
			return errs.ErrCannotDeleteTransactionInParentAccount
		}

		// Update transaction row to deleted
		deletedRows, err := sess.ID(oldTransaction.TransactionId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionNotFound
		}

		if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			deletedRows, err = sess.ID(oldTransaction.RelatedId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

			if err != nil {
				return err
			} else if deletedRows < 1 {
				return errs.ErrTransactionNotFound
			}
		}

		// Update transaction tag index
		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", uid, false, oldTransaction.TransactionId).Update(tagIndexUpdateModel)

		if err != nil {
			return err
		}

		// Update transaction picture
		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", uid, false, oldTransaction.TransactionId).Update(pictureUpdateModel)

		if err != nil {
			return err
		}

		// Update account table
		if oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			if oldTransaction.RelatedAccountAmount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.DeleteTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			if oldTransaction.Amount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.DeleteTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			if oldTransaction.Amount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", oldTransaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					log.Errorf(c, "[transactions.DeleteTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			if oldTransaction.Amount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedSourceRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", oldTransaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedSourceRows < 1 {
					log.Errorf(c, "[transactions.DeleteTransaction] failed to update account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}

			if oldTransaction.RelatedAccountAmount != 0 {
				destinationAccount.UpdatedUnixTime = time.Now().Unix()
				updatedDestinationRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

				if err != nil {
					return err
				} else if updatedDestinationRows < 1 {
					log.Errorf(c, "[transactions.DeleteTransaction] failed to update related account balance")
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			return errs.ErrTransactionTypeInvalid
		}

		return err
	})
}

// DeleteAllTransactions deletes all existed transactions from database
func (s *TransactionService) DeleteAllTransactions(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Transaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	tagIndexUpdateModel := &models.TransactionTagIndex{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	pictureUpdateModel := &models.TransactionPictureInfo{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	accountUpdateModel := &models.Account{
		Balance:         0,
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Update all transaction to deleted
		_, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		// Update all transaction tag index to deleted
		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(tagIndexUpdateModel)

		if err != nil {
			return err
		}

		// Update all transaction picture to deleted
		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(pictureUpdateModel)

		if err != nil {
			return err
		}

		// Update all account table to deleted
		_, err = sess.Cols("balance", "deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(accountUpdateModel)

		if err != nil {
			return err
		}

		return nil
	})
}

// GetRelatedTransferTransaction returns the related transaction for transfer transaction
func (s *TransactionService) GetRelatedTransferTransaction(originalTransaction *models.Transaction) *models.Transaction {
	var relatedType models.TransactionDbType
	var relatedTransactionTime int64

	if originalTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		relatedType = models.TRANSACTION_DB_TYPE_TRANSFER_IN
		relatedTransactionTime = originalTransaction.TransactionTime + 1
	} else if originalTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		relatedType = models.TRANSACTION_DB_TYPE_TRANSFER_OUT
		relatedTransactionTime = originalTransaction.TransactionTime - 1
	} else {
		return nil
	}

	relatedTransaction := &models.Transaction{
		TransactionId:        originalTransaction.RelatedId,
		Uid:                  originalTransaction.Uid,
		Deleted:              originalTransaction.Deleted,
		Type:                 relatedType,
		CategoryId:           originalTransaction.CategoryId,
		TransactionTime:      relatedTransactionTime,
		TimezoneUtcOffset:    originalTransaction.TimezoneUtcOffset,
		AccountId:            originalTransaction.RelatedAccountId,
		Amount:               originalTransaction.RelatedAccountAmount,
		RelatedId:            originalTransaction.TransactionId,
		RelatedAccountId:     originalTransaction.AccountId,
		RelatedAccountAmount: originalTransaction.Amount,
		Comment:              originalTransaction.Comment,
		GeoLongitude:         originalTransaction.GeoLongitude,
		GeoLatitude:          originalTransaction.GeoLatitude,
		CreatedIp:            originalTransaction.CreatedIp,
		CreatedUnixTime:      originalTransaction.CreatedUnixTime,
		UpdatedUnixTime:      originalTransaction.UpdatedUnixTime,
		DeletedUnixTime:      originalTransaction.DeletedUnixTime,
	}

	return relatedTransaction
}

// GetAccountsTotalIncomeAndExpense returns the every accounts total income and expense amount by specific date range
func (s *TransactionService) GetAccountsTotalIncomeAndExpense(c core.Context, uid int64, startUnixTime int64, endUnixTime int64, utcOffset int16, useTransactionTimezone bool) (map[int64]int64, map[int64]int64, error) {
	if uid <= 0 {
		return nil, nil, errs.ErrUserIdInvalid
	}

	clientLocation := time.FixedZone("Client Timezone", int(utcOffset)*60)
	startLocalDateTime := utils.FormatUnixTimeToNumericLocalDateTime(startUnixTime, clientLocation)
	endLocalDateTime := utils.FormatUnixTimeToNumericLocalDateTime(endUnixTime, clientLocation)

	startUnixTime = utils.GetMinUnixTimeWithSameLocalDateTime(startUnixTime, utcOffset)
	endUnixTime = utils.GetMaxUnixTimeWithSameLocalDateTime(endUnixTime, utcOffset)

	startTransactionTime := utils.GetMinTransactionTimeFromUnixTime(startUnixTime)
	endTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(endUnixTime)

	condition := "uid=? AND deleted=? AND (type=? OR type=?) AND transaction_time>=? AND transaction_time<=?"
	conditionParams := make([]any, 0, 4)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)

	minTransactionTime := startTransactionTime
	maxTransactionTime := endTransactionTime
	var allTransactions []*models.Transaction

	for maxTransactionTime > 0 {
		var transactions []*models.Transaction

		finalConditionParams := make([]any, 0, 6)
		finalConditionParams = append(finalConditionParams, conditionParams...)
		finalConditionParams = append(finalConditionParams, minTransactionTime)
		finalConditionParams = append(finalConditionParams, maxTransactionTime)

		err := s.UserDataDB(uid).NewSession(c).Select("type, account_id, transaction_time, timezone_utc_offset, amount").Where(condition, finalConditionParams...).Limit(pageCountForLoadTransactionAmounts, 0).OrderBy("transaction_time desc").Find(&transactions)

		if err != nil {
			return nil, nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < pageCountForLoadTransactionAmounts {
			maxTransactionTime = 0
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	incomeAmounts := make(map[int64]int64)
	expenseAmounts := make(map[int64]int64)

	for i := 0; i < len(allTransactions); i++ {
		transaction := allTransactions[i]
		timeZone := clientLocation

		if useTransactionTimezone {
			timeZone = time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
		}

		localDateTime := utils.FormatUnixTimeToNumericLocalDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), timeZone)

		if localDateTime < startLocalDateTime || localDateTime > endLocalDateTime {
			continue
		}

		var amountsMap map[int64]int64

		if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			amountsMap = incomeAmounts
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			amountsMap = expenseAmounts
		}

		totalAmounts, exists := amountsMap[transaction.AccountId]

		if !exists {
			totalAmounts = 0
		}

		totalAmounts += transaction.Amount
		amountsMap[transaction.AccountId] = totalAmounts
	}

	return incomeAmounts, expenseAmounts, nil
}

// GetAccountsAndCategoriesTotalIncomeAndExpense returns the every accounts and categories total income and expense amount by specific date range
func (s *TransactionService) GetAccountsAndCategoriesTotalIncomeAndExpense(c core.Context, uid int64, startUnixTime int64, endUnixTime int64, tagIds []int64, noTags bool, tagFilterType models.TransactionTagFilterType, utcOffset int16, useTransactionTimezone bool) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	clientLocation := time.FixedZone("Client Timezone", int(utcOffset)*60)
	var startLocalDateTime, endLocalDateTime, startTransactionTime, endTransactionTime int64

	if startUnixTime > 0 {
		startLocalDateTime = utils.FormatUnixTimeToNumericLocalDateTime(startUnixTime, clientLocation)
		startUnixTime = utils.GetMinUnixTimeWithSameLocalDateTime(startUnixTime, utcOffset)
		startTransactionTime = utils.GetMinTransactionTimeFromUnixTime(startUnixTime)
	}

	if endUnixTime > 0 {
		endLocalDateTime = utils.FormatUnixTimeToNumericLocalDateTime(endUnixTime, clientLocation)
		endUnixTime = utils.GetMaxUnixTimeWithSameLocalDateTime(endUnixTime, utcOffset)
		endTransactionTime = utils.GetMaxTransactionTimeFromUnixTime(endUnixTime)
	}

	condition := "uid=? AND deleted=? AND (type=? OR type=?)"
	conditionParams := make([]any, 0, 4)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)

	minTransactionTime := startTransactionTime
	maxTransactionTime := endTransactionTime
	var allTransactions []*models.Transaction

	for maxTransactionTime >= 0 {
		var transactions []*models.Transaction

		finalCondition := condition
		finalConditionParams := make([]any, 0, 6)
		finalConditionParams = append(finalConditionParams, conditionParams...)

		if minTransactionTime > 0 {
			finalCondition = finalCondition + " AND transaction_time>=?"
			finalConditionParams = append(finalConditionParams, minTransactionTime)
		}

		if maxTransactionTime > 0 {
			finalCondition = finalCondition + " AND transaction_time<=?"
			finalConditionParams = append(finalConditionParams, maxTransactionTime)
		}

		sess := s.UserDataDB(uid).NewSession(c).Select("category_id, account_id, transaction_time, timezone_utc_offset, amount").Where(finalCondition, finalConditionParams...)
		sess = s.appendFilterTagIdsConditionToQuery(sess, uid, maxTransactionTime, minTransactionTime, tagIds, noTags, tagFilterType)

		err := sess.Limit(pageCountForLoadTransactionAmounts, 0).OrderBy("transaction_time desc").Find(&transactions)

		if err != nil {
			return nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < pageCountForLoadTransactionAmounts {
			maxTransactionTime = -1
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	transactionTotalAmountsMap := make(map[string]*models.Transaction)

	for i := 0; i < len(allTransactions); i++ {
		transaction := allTransactions[i]
		timeZone := clientLocation

		if useTransactionTimezone {
			timeZone = time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
		}

		localDateTime := utils.FormatUnixTimeToNumericLocalDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), timeZone)

		if (startLocalDateTime > 0 && localDateTime < startLocalDateTime) || (endLocalDateTime > 0 && localDateTime > endLocalDateTime) {
			continue
		}

		groupKey := fmt.Sprintf("%d_%d", transaction.CategoryId, transaction.AccountId)
		totalAmounts, exists := transactionTotalAmountsMap[groupKey]

		if !exists {
			totalAmounts = &models.Transaction{
				CategoryId: transaction.CategoryId,
				AccountId:  transaction.AccountId,
				Amount:     0,
			}

			transactionTotalAmountsMap[groupKey] = totalAmounts
		}

		totalAmounts.Amount += transaction.Amount
	}

	transactionTotalAmounts := make([]*models.Transaction, 0, len(transactionTotalAmountsMap))

	for _, totalAmounts := range transactionTotalAmountsMap {
		transactionTotalAmounts = append(transactionTotalAmounts, totalAmounts)
	}

	return transactionTotalAmounts, nil
}

// GetAccountsAndCategoriesMonthlyIncomeAndExpense returns the every accounts monthly income and expense amount by specific date range
func (s *TransactionService) GetAccountsAndCategoriesMonthlyIncomeAndExpense(c core.Context, uid int64, startYear int32, startMonth int32, endYear int32, endMonth int32, tagIds []int64, noTags bool, tagFilterType models.TransactionTagFilterType, utcOffset int16, useTransactionTimezone bool) (map[int32][]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	clientLocation := time.FixedZone("Client Timezone", int(utcOffset)*60)
	var startTransactionTime, endTransactionTime int64
	var err error

	if startYear > 0 && startMonth > 0 {
		startTransactionTime, _, err = utils.GetTransactionTimeRangeByYearMonth(startYear, startMonth)

		if err != nil {
			return nil, errs.ErrSystemError
		}
	}

	if endYear > 0 && endMonth > 0 {
		_, endTransactionTime, err = utils.GetTransactionTimeRangeByYearMonth(endYear, endMonth)

		if err != nil {
			return nil, errs.ErrSystemError
		}
	}

	condition := "uid=? AND deleted=? AND (type=? OR type=?)"
	conditionParams := make([]any, 0, 4)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)

	minTransactionTime := startTransactionTime
	maxTransactionTime := endTransactionTime
	var allTransactions []*models.Transaction

	for maxTransactionTime >= 0 {
		var transactions []*models.Transaction

		finalCondition := condition
		finalConditionParams := make([]any, 0, 6)
		finalConditionParams = append(finalConditionParams, conditionParams...)

		if minTransactionTime > 0 {
			finalCondition = finalCondition + " AND transaction_time>=?"
			finalConditionParams = append(finalConditionParams, minTransactionTime)
		}

		if maxTransactionTime > 0 {
			finalCondition = finalCondition + " AND transaction_time<=?"
			finalConditionParams = append(finalConditionParams, maxTransactionTime)
		}

		sess := s.UserDataDB(uid).NewSession(c).Select("category_id, account_id, transaction_time, timezone_utc_offset, amount").Where(finalCondition, finalConditionParams...)
		sess = s.appendFilterTagIdsConditionToQuery(sess, uid, maxTransactionTime, minTransactionTime, tagIds, noTags, tagFilterType)

		err := sess.Limit(pageCountForLoadTransactionAmounts, 0).OrderBy("transaction_time desc").Find(&transactions)

		if err != nil {
			return nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < pageCountForLoadTransactionAmounts {
			maxTransactionTime = -1
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	startYearMonth := startYear*100 + startMonth
	endYearMonth := endYear*100 + endMonth
	transactionsMonthlyAmountsMap := make(map[string]*models.Transaction)
	transactionsMonthlyAmounts := make(map[int32][]*models.Transaction)

	for i := 0; i < len(allTransactions); i++ {
		transaction := allTransactions[i]
		timeZone := clientLocation

		if useTransactionTimezone {
			timeZone = time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
		}

		yearMonth := utils.FormatUnixTimeToNumericYearMonth(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), timeZone)

		if (startYearMonth > 0 && yearMonth < startYearMonth) || (endYearMonth > 0 && yearMonth > endYearMonth) {
			continue
		}

		groupKey := fmt.Sprintf("%d_%d_%d", yearMonth, transaction.CategoryId, transaction.AccountId)
		transactionAmounts, exists := transactionsMonthlyAmountsMap[groupKey]

		if !exists {
			transactionAmounts = &models.Transaction{
				CategoryId: transaction.CategoryId,
				AccountId:  transaction.AccountId,
			}
			transactionsMonthlyAmountsMap[groupKey] = transactionAmounts
		}

		transactionAmounts.Amount += transaction.Amount
	}

	for groupKey, transaction := range transactionsMonthlyAmountsMap {
		groupKeyParts := strings.Split(groupKey, "_")
		yearMonth, _ := utils.StringToInt32(groupKeyParts[0])
		monthlyAmounts, exists := transactionsMonthlyAmounts[yearMonth]

		if !exists {
			monthlyAmounts = make([]*models.Transaction, 0, 0)
		}

		monthlyAmounts = append(monthlyAmounts, transaction)
		transactionsMonthlyAmounts[yearMonth] = monthlyAmounts
	}

	return transactionsMonthlyAmounts, nil
}

// GetTransactionMapByList returns a transaction map by a list
func (s *TransactionService) GetTransactionMapByList(transactions []*models.Transaction) map[int64]*models.Transaction {
	transactionMap := make(map[int64]*models.Transaction)

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		transactionMap[transaction.TransactionId] = transaction
	}

	return transactionMap
}

// GetTransactionIds returns transaction ids list
func (s *TransactionService) GetTransactionIds(transactions []*models.Transaction) []int64 {
	transactionIds := make([]int64, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transactionIds[i] = transactions[i].TransactionId
	}

	return transactionIds
}

func (s *TransactionService) doCreateTransaction(c core.Context, database *datastore.Database, sess *xorm.Session, transaction *models.Transaction, transactionTagIndexes []*models.TransactionTagIndex, tagIds []int64, pictureIds []int64, pictureUpdateModel *models.TransactionPictureInfo) error {
	// Get and verify source and destination account
	sourceAccount, destinationAccount, err := s.getAccountModels(sess, transaction)

	if err != nil {
		return err
	}

	if sourceAccount.Hidden || (destinationAccount != nil && destinationAccount.Hidden) {
		return errs.ErrCannotAddTransactionToHiddenAccount
	}

	if sourceAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS || (destinationAccount != nil && destinationAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS) {
		return errs.ErrCannotAddTransactionToParentAccount
	}

	if (transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) &&
		sourceAccount.Currency == destinationAccount.Currency && transaction.Amount != transaction.RelatedAccountAmount {
		return errs.ErrTransactionSourceAndDestinationAmountNotEqual
	}

	if (transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) &&
		(transaction.Amount < 0 || transaction.RelatedAccountAmount < 0) {
		return errs.ErrTransferTransactionAmountCannotBeLessThanZero
	}

	// Get and verify category
	err = s.isCategoryValid(sess, transaction)

	if err != nil {
		return err
	}

	// Get and verify tags
	err = s.isTagsValid(sess, transaction, transactionTagIndexes, tagIds)

	if err != nil {
		return err
	}

	// Get and verify pictures
	err = s.isPicturesValid(sess, transaction, pictureIds)

	if err != nil {
		return err
	}

	// Verify balance modification transaction and calculate real amount
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		otherTransactionExists, err := sess.Cols("uid", "deleted", "account_id").Where("uid=? AND deleted=? AND account_id=?", transaction.Uid, false, sourceAccount.AccountId).Limit(1).Exist(&models.Transaction{})

		if err != nil {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to get whether other transactions exist, because %s", err.Error())
			return err
		} else if otherTransactionExists {
			return errs.ErrBalanceModificationTransactionCannotAddWhenNotEmpty
		}

		transaction.RelatedAccountId = transaction.AccountId
		transaction.RelatedAccountAmount = transaction.Amount - sourceAccount.Balance
	} else { // Not allow to add transaction before balance modification transaction
		otherTransactionExists := false

		if destinationAccount != nil && sourceAccount.AccountId != destinationAccount.AccountId {
			otherTransactionExists, err = sess.Cols("uid", "deleted", "account_id").Where("uid=? AND deleted=? AND type=? AND (account_id=? OR account_id=?) AND transaction_time>=?", transaction.Uid, false, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, sourceAccount.AccountId, destinationAccount.AccountId, transaction.TransactionTime).Limit(1).Exist(&models.Transaction{})
		} else {
			otherTransactionExists, err = sess.Cols("uid", "deleted", "account_id").Where("uid=? AND deleted=? AND type=? AND account_id=? AND transaction_time>=?", transaction.Uid, false, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE, sourceAccount.AccountId, transaction.TransactionTime).Limit(1).Exist(&models.Transaction{})
		}

		if err != nil {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to get whether other transactions exist, because %s", err.Error())
			return err
		} else if otherTransactionExists {
			return errs.ErrCannotAddTransactionBeforeBalanceModificationTransaction
		}
	}

	// Insert transaction row
	var relatedTransaction *models.Transaction

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		relatedTransaction = s.GetRelatedTransferTransaction(transaction)
	}

	insertTransactionSavePointName := "insert_transaction"
	err = database.SetSavePoint(sess, insertTransactionSavePointName)

	if err != nil {
		log.Errorf(c, "[transactions.doCreateTransaction] failed to set save point \"%s\", because %s", insertTransactionSavePointName, err.Error())
		return err
	}

	createdRows, err := sess.Insert(transaction)

	if err != nil || createdRows < 1 { // maybe another transaction has same time
		if err != nil {
			log.Warnf(c, "[transactions.doCreateTransaction] cannot create trasaction, because %s, regenerate transaction time value", err.Error())
		} else {
			log.Warnf(c, "[transactions.doCreateTransaction] cannot create trasaction, regenerate transaction time value")
		}

		err = database.RollbackToSavePoint(sess, insertTransactionSavePointName)

		if err != nil {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to rollback to save point \"%s\", because %s", insertTransactionSavePointName, err.Error())
			return err
		}

		sameSecondLatestTransaction := &models.Transaction{}
		minTransactionTime := utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))
		maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

		has, err := sess.Where("uid=? AND transaction_time>=? AND transaction_time<=?", transaction.Uid, minTransactionTime, maxTransactionTime).OrderBy("transaction_time desc").Limit(1).Get(sameSecondLatestTransaction)

		if err != nil {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to get trasaction time, because %s", err.Error())
			return err
		} else if !has {
			log.Errorf(c, "[transactions.doCreateTransaction] it should have transactions in %d - %d, but result is empty", minTransactionTime, maxTransactionTime)
			return errs.ErrDatabaseOperationFailed
		} else if sameSecondLatestTransaction.TransactionTime == maxTransactionTime-1 {
			return errs.ErrTooMuchTransactionInOneSecond
		}

		transaction.TransactionTime = sameSecondLatestTransaction.TransactionTime + 1
		createdRows, err := sess.Insert(transaction)

		if err != nil {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to add transaction again, because %s", err.Error())
			return err
		} else if createdRows < 1 {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to add transaction again")
			return errs.ErrDatabaseOperationFailed
		}
	}

	if relatedTransaction != nil {
		relatedTransaction.TransactionTime = transaction.TransactionTime + 1

		if utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime) != utils.GetUnixTimeFromTransactionTime(relatedTransaction.TransactionTime) {
			return errs.ErrTooMuchTransactionInOneSecond
		}

		createdRows, err := sess.Insert(relatedTransaction)

		if err != nil {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to add related transaction, because %s", err.Error())
			return err
		} else if createdRows < 1 {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to add related transaction")
			return errs.ErrDatabaseOperationFailed
		}
	}

	err = nil

	// Insert transaction tag index
	if len(transactionTagIndexes) > 0 {
		for i := 0; i < len(transactionTagIndexes); i++ {
			transactionTagIndex := transactionTagIndexes[i]
			transactionTagIndex.TransactionTime = transaction.TransactionTime

			_, err := sess.Insert(transactionTagIndex)

			if err != nil {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to add transaction tag index, because %s", err.Error())
				return err
			}
		}
	}

	// Update transaction picture
	if len(pictureIds) > 0 {
		_, err = sess.Cols("transaction_id", "updated_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", transaction.Uid, false, models.TransactionPictureNewPictureTransactionId).In("picture_id", pictureIds).Update(pictureUpdateModel)

		if err != nil {
			log.Errorf(c, "[transactions.doCreateTransaction] failed to update transaction picture info, because %s", err.Error())
			return err
		}
	}

	// Update account table
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.RelatedAccountAmount != 0 {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance, because %s", err.Error())
				return err
			} else if updatedRows < 1 {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance")
				return errs.ErrDatabaseOperationFailed
			}
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
		if transaction.Amount != 0 {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance, because %s", err.Error())
				return err
			} else if updatedRows < 1 {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance")
				return errs.ErrDatabaseOperationFailed
			}
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
		if transaction.Amount != 0 {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", transaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance, because %s", err.Error())
				return err
			} else if updatedRows < 1 {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance")
				return errs.ErrDatabaseOperationFailed
			}
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if transaction.Amount != 0 {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedSourceRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", transaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance, because %s", err.Error())
				return err
			} else if updatedSourceRows < 1 {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance")
				return errs.ErrDatabaseOperationFailed
			}
		}

		if transaction.RelatedAccountAmount != 0 {
			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedDestinationRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance, because %s", err.Error())
				return err
			} else if updatedDestinationRows < 1 {
				log.Errorf(c, "[transactions.doCreateTransaction] failed to update account balance")
				return errs.ErrDatabaseOperationFailed
			}
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return errs.ErrTransactionTypeInvalid
	}

	return err
}

func (s *TransactionService) buildTransactionQueryCondition(uid int64, maxTransactionTime int64, minTransactionTime int64, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, tagIds []int64, amountFilter string, keyword string, noDuplicated bool) (string, []any) {
	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, 16)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	if maxTransactionTime > 0 {
		condition = condition + " AND transaction_time<=?"
		conditionParams = append(conditionParams, maxTransactionTime)
	}

	if minTransactionTime > 0 {
		condition = condition + " AND transaction_time>=?"
		conditionParams = append(conditionParams, minTransactionTime)
	}

	var accountIdsCondition strings.Builder
	accountIdConditionParams := make([]any, 0, len(accountIds))

	for i := 0; i < len(accountIds); i++ {
		if i > 0 {
			accountIdsCondition.WriteString(",")
		}

		accountIdsCondition.WriteString("?")
		accountIdConditionParams = append(accountIdConditionParams, accountIds[i])
	}

	if models.TRANSACTION_DB_TYPE_MODIFY_BALANCE <= transactionType && transactionType <= models.TRANSACTION_DB_TYPE_EXPENSE {
		condition = condition + " AND type=?"
		conditionParams = append(conditionParams, transactionType)
	} else if transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		if len(accountIds) == 0 {
			condition = condition + " AND type=?"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
		} else if len(accountIds) == 1 {
			condition = condition + " AND (type=? OR type=?)"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_IN)
		} else { // len(accountsIds) > 1
			condition = condition + " AND (type=? OR (type=? AND related_account_id NOT IN (" + accountIdsCondition.String() + ")))"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_IN)
			conditionParams = append(conditionParams, accountIdConditionParams...)
		}
	} else {
		if noDuplicated {
			if len(accountIds) == 0 {
				condition = condition + " AND (type=? OR type=? OR type=? OR type=?)"
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE)
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
			} else if len(accountIds) == 1 {
				// Do Nothing
			} else { // len(accountsIds) > 1
				condition = condition + " AND (type=? OR type=? OR type=? OR type=? OR (type=? AND related_account_id NOT IN (" + accountIdsCondition.String() + ")))"
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE)
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
				conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_IN)
				conditionParams = append(conditionParams, accountIdConditionParams...)
			}
		}
	}

	if len(categoryIds) > 0 {
		var conditions strings.Builder

		for i := 0; i < len(categoryIds); i++ {
			if i > 0 {
				conditions.WriteString(",")
			}

			conditions.WriteString("?")
			conditionParams = append(conditionParams, categoryIds[i])
		}

		if conditions.Len() > 1 {
			condition = condition + " AND category_id IN (" + conditions.String() + ")"
		} else {
			condition = condition + " AND category_id = " + conditions.String()
		}
	}

	if len(accountIds) > 0 {
		if accountIdsCondition.Len() > 1 {
			condition = condition + " AND account_id IN (" + accountIdsCondition.String() + ")"
		} else {
			condition = condition + " AND account_id = " + accountIdsCondition.String()
		}

		conditionParams = append(conditionParams, accountIdConditionParams...)
	}

	if amountFilter != "" {
		amountFilterItems := strings.Split(amountFilter, ":")

		if len(amountFilterItems) == 2 && amountFilterItems[0] == "gt" {
			value, err := utils.StringToInt64(amountFilterItems[1])

			if err == nil {
				condition = condition + " AND amount > ?"
				conditionParams = append(conditionParams, value)
			}
		} else if len(amountFilterItems) == 2 && amountFilterItems[0] == "lt" {
			value, err := utils.StringToInt64(amountFilterItems[1])

			if err == nil {
				condition = condition + " AND amount < ?"
				conditionParams = append(conditionParams, value)
			}
		} else if len(amountFilterItems) == 2 && amountFilterItems[0] == "eq" {
			value, err := utils.StringToInt64(amountFilterItems[1])

			if err == nil {
				condition = condition + " AND amount = ?"
				conditionParams = append(conditionParams, value)
			}
		} else if len(amountFilterItems) == 2 && amountFilterItems[0] == "ne" {
			value, err := utils.StringToInt64(amountFilterItems[1])

			if err == nil {
				condition = condition + " AND amount <> ?"
				conditionParams = append(conditionParams, value)
			}
		} else if len(amountFilterItems) == 3 && amountFilterItems[0] == "bt" {
			value1, err := utils.StringToInt64(amountFilterItems[1])
			value2, err := utils.StringToInt64(amountFilterItems[2])

			if err == nil {
				condition = condition + " AND amount >= ? AND amount <= ?"
				conditionParams = append(conditionParams, value1)
				conditionParams = append(conditionParams, value2)
			}
		} else if len(amountFilterItems) == 3 && amountFilterItems[0] == "nb" {
			value1, err := utils.StringToInt64(amountFilterItems[1])
			value2, err := utils.StringToInt64(amountFilterItems[2])

			if err == nil {
				condition = condition + " AND (amount < ? OR amount > ?)"
				conditionParams = append(conditionParams, value1)
				conditionParams = append(conditionParams, value2)
			}
		}
	}

	if keyword != "" {
		condition = condition + " AND comment LIKE ?"
		conditionParams = append(conditionParams, "%%"+keyword+"%%")
	}

	return condition, conditionParams
}

func (s *TransactionService) appendFilterTagIdsConditionToQuery(sess *xorm.Session, uid int64, maxTransactionTime int64, minTransactionTime int64, tagIds []int64, noTags bool, tagFilterType models.TransactionTagFilterType) *xorm.Session {
	subQueryCondition := builder.And(builder.Eq{"uid": uid}, builder.Eq{"deleted": false})

	if maxTransactionTime > 0 {
		subQueryCondition = subQueryCondition.And(builder.Lte{"transaction_time": maxTransactionTime})
	}

	if minTransactionTime > 0 {
		subQueryCondition = subQueryCondition.And(builder.Gte{"transaction_time": minTransactionTime})
	}

	if noTags {
		subQuery := builder.Select("transaction_id").From("transaction_tag_index").Where(subQueryCondition)
		sess.NotIn("transaction_id", subQuery).NotIn("related_id", subQuery)
		return sess
	}

	if len(tagIds) < 1 {
		return sess
	}

	subQueryCondition = subQueryCondition.And(builder.In("tag_id", tagIds))
	subQuery := builder.Select("transaction_id").From("transaction_tag_index").Where(subQueryCondition)

	if tagFilterType == models.TRANSACTION_TAG_FILTER_HAS_ALL || tagFilterType == models.TRANSACTION_TAG_FILTER_NOT_HAS_ALL {
		subQuery = subQuery.GroupBy("transaction_id").Having(fmt.Sprintf("COUNT(DISTINCT tag_id) >= %d", len(tagIds)))
	}

	if tagFilterType == models.TRANSACTION_TAG_FILTER_HAS_ANY || tagFilterType == models.TRANSACTION_TAG_FILTER_HAS_ALL {
		sess.And(builder.Or(builder.In("transaction_id", subQuery), builder.In("related_id", subQuery)))
	} else if tagFilterType == models.TRANSACTION_TAG_FILTER_NOT_HAS_ANY || tagFilterType == models.TRANSACTION_TAG_FILTER_NOT_HAS_ALL {
		sess.NotIn("transaction_id", subQuery).NotIn("related_id", subQuery)
	}

	return sess
}

func (s *TransactionService) isAccountIdValid(transaction *models.Transaction) error {
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.RelatedAccountId != 0 && transaction.RelatedAccountId != transaction.AccountId {
			return errs.ErrTransactionDestinationAccountCannotBeSet
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME ||
		transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
		if transaction.RelatedAccountId != 0 {
			return errs.ErrTransactionDestinationAccountCannotBeSet
		} else if transaction.RelatedAccountAmount != 0 {
			return errs.ErrTransactionDestinationAmountCannotBeSet
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if transaction.AccountId == transaction.RelatedAccountId {
			return errs.ErrTransactionSourceAndDestinationIdCannotBeEqual
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return errs.ErrTransactionTypeInvalid
	} else {
		return errs.ErrTransactionTypeInvalid
	}

	return nil
}

func (s *TransactionService) getAccountModels(sess *xorm.Session, transaction *models.Transaction) (sourceAccount *models.Account, destinationAccount *models.Account, err error) {
	sourceAccount = &models.Account{}
	destinationAccount = &models.Account{}

	has, err := sess.ID(transaction.AccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(sourceAccount)

	if err != nil {
		return nil, nil, err
	} else if !has {
		return nil, nil, errs.ErrSourceAccountNotFound
	}

	// check whether the related account is valid
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.RelatedAccountId != 0 && transaction.RelatedAccountId != transaction.AccountId {
			return nil, nil, errs.ErrAccountIdInvalid
		} else {
			destinationAccount = sourceAccount
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME || transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
		if transaction.RelatedAccountId != 0 {
			return nil, nil, errs.ErrAccountIdInvalid
		}

		destinationAccount = nil
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		if transaction.RelatedAccountId <= 0 {
			return nil, nil, errs.ErrAccountIdInvalid
		} else {
			has, err = sess.ID(transaction.RelatedAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(destinationAccount)

			if err != nil {
				return nil, nil, err
			} else if !has {
				return nil, nil, errs.ErrDestinationAccountNotFound
			}
		}
	}

	// check whether the parent accounts are valid
	if sourceAccount.ParentAccountId > 0 && destinationAccount != nil && sourceAccount.ParentAccountId != destinationAccount.ParentAccountId && destinationAccount.ParentAccountId > 0 {
		var accounts []*models.Account
		err := sess.Where("uid=? AND deleted=? and (account_id=? or account_id=?)", transaction.Uid, false, sourceAccount.ParentAccountId, destinationAccount.ParentAccountId).Find(&accounts)

		if err != nil {
			return nil, nil, err
		}

		if len(accounts) < 2 {
			return nil, nil, errs.ErrAccountNotFound
		}

		for i := 0; i < len(accounts); i++ {
			account := accounts[i]

			if account.Hidden {
				return nil, nil, errs.ErrCannotUseHiddenAccount
			}
		}
	} else if sourceAccount.ParentAccountId > 0 && (destinationAccount == nil || sourceAccount.ParentAccountId == destinationAccount.ParentAccountId || destinationAccount.ParentAccountId == 0) {
		sourceParentAccount := &models.Account{}
		has, err = sess.ID(sourceAccount.ParentAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(sourceParentAccount)

		if err != nil {
			return nil, nil, err
		} else if !has {
			return nil, nil, errs.ErrSourceAccountNotFound
		}

		if sourceParentAccount.Hidden {
			return nil, nil, errs.ErrCannotUseHiddenAccount
		}
	} else if sourceAccount.ParentAccountId == 0 && destinationAccount != nil && destinationAccount.ParentAccountId > 0 {
		destinationParentAccount := &models.Account{}
		has, err = sess.ID(destinationAccount.ParentAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(destinationParentAccount)

		if err != nil {
			return nil, nil, err
		} else if !has {
			return nil, nil, errs.ErrDestinationAccountNotFound
		}

		if destinationParentAccount.Hidden {
			return nil, nil, errs.ErrCannotUseHiddenAccount
		}
	}

	return sourceAccount, destinationAccount, nil
}

func (s *TransactionService) getOldAccountModels(sess *xorm.Session, transaction *models.Transaction, oldTransaction *models.Transaction, sourceAccount *models.Account, destinationAccount *models.Account) (oldSourceAccount *models.Account, oldDestinationAccount *models.Account, err error) {
	oldSourceAccount = &models.Account{}
	oldDestinationAccount = &models.Account{}

	if transaction.AccountId == oldTransaction.AccountId {
		oldSourceAccount = sourceAccount
	} else {
		has, err := sess.ID(oldTransaction.AccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(oldSourceAccount)

		if err != nil {
			return nil, nil, err
		} else if !has {
			return nil, nil, errs.ErrSourceAccountNotFound
		}
	}

	if transaction.RelatedAccountId == oldTransaction.RelatedAccountId {
		oldDestinationAccount = destinationAccount
	} else {
		has, err := sess.ID(oldTransaction.RelatedAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(oldDestinationAccount)

		if err != nil {
			return nil, nil, err
		} else if !has {
			return nil, nil, errs.ErrDestinationAccountNotFound
		}
	}
	return oldSourceAccount, oldDestinationAccount, nil
}

func (s *TransactionService) getRelatedUpdateColumns(updateCols []string) []string {
	relatedUpdateCols := make([]string, len(updateCols))

	for i := 0; i < len(updateCols); i++ {
		if updateCols[i] == "account_id" {
			relatedUpdateCols[i] = "related_account_id"
		} else if updateCols[i] == "related_account_id" {
			relatedUpdateCols[i] = "account_id"
		} else if updateCols[i] == "amount" {
			relatedUpdateCols[i] = "related_account_amount"
		} else if updateCols[i] == "related_account_amount" {
			relatedUpdateCols[i] = "amount"
		} else {
			relatedUpdateCols[i] = updateCols[i]
		}
	}

	return relatedUpdateCols
}

func (s *TransactionService) isCategoryValid(sess *xorm.Session, transaction *models.Transaction) error {
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.CategoryId != 0 {
			return errs.ErrBalanceModificationTransactionCannotSetCategory
		}
	} else {
		category := &models.TransactionCategory{}
		has, err := sess.ID(transaction.CategoryId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(category)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrTransactionCategoryNotFound
		}

		if category.Hidden {
			return errs.ErrCannotUseHiddenTransactionCategory
		}

		if category.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
			return errs.ErrCannotUsePrimaryCategoryForTransaction
		}

		if (transaction.Type == models.TRANSACTION_DB_TYPE_INCOME && category.Type != models.CATEGORY_TYPE_INCOME) ||
			(transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE && category.Type != models.CATEGORY_TYPE_EXPENSE) ||
			((transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) && category.Type != models.CATEGORY_TYPE_TRANSFER) {
			return errs.ErrTransactionCategoryTypeInvalid
		}

		parentCategory := &models.TransactionCategory{}
		has, err = sess.ID(category.ParentCategoryId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(parentCategory)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrTransactionCategoryNotFound
		}

		if parentCategory.Hidden {
			return errs.ErrCannotUseHiddenTransactionCategory
		}
	}

	return nil
}

func (s *TransactionService) isTagsValid(sess *xorm.Session, transaction *models.Transaction, transactionTagIndexes []*models.TransactionTagIndex, tagIds []int64) error {
	if len(transactionTagIndexes) > 0 {
		var tags []*models.TransactionTag
		err := sess.Where("uid=? AND deleted=?", transaction.Uid, false).In("tag_id", tagIds).Find(&tags)

		if err != nil {
			return err
		}

		tagMap := make(map[int64]*models.TransactionTag)

		for i := 0; i < len(tags); i++ {
			if tags[i].Hidden {
				return errs.ErrCannotUseHiddenTransactionTag
			}

			tagMap[tags[i].TagId] = tags[i]
		}

		for i := 0; i < len(transactionTagIndexes); i++ {
			if _, exists := tagMap[transactionTagIndexes[i].TagId]; !exists {
				return errs.ErrTransactionTagNotFound
			}
		}
	}

	return nil
}

func (s *TransactionService) isPicturesValid(sess *xorm.Session, transaction *models.Transaction, pictureIds []int64) error {
	if len(pictureIds) > 0 {
		var pictureInfos []*models.TransactionPictureInfo
		err := sess.Where("uid=? AND deleted=?", transaction.Uid, false).In("picture_id", pictureIds).Find(&pictureInfos)

		if err != nil {
			return err
		}

		pictureInfoMap := make(map[int64]*models.TransactionPictureInfo)

		for i := 0; i < len(pictureInfos); i++ {
			if pictureInfos[i].TransactionId != models.TransactionPictureNewPictureTransactionId && pictureInfos[i].TransactionId != transaction.TransactionId {
				return errs.ErrTransactionPictureIdInvalid
			}

			pictureInfoMap[pictureInfos[i].PictureId] = pictureInfos[i]
		}

		for i := 0; i < len(pictureIds); i++ {
			if _, exists := pictureInfoMap[pictureIds[i]]; !exists {
				return errs.ErrTransactionPictureNotFound
			}
		}
	}

	return nil
}
