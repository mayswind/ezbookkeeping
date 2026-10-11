package services

import (
	"strings"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// QueuedTransactionService represents the queued transaction (draft queue) service
type QueuedTransactionService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a queued transaction service singleton instance
var (
	QueuedTransactions = &QueuedTransactionService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetAllQueuedTransactionsByUid returns the queue of the user, optionally including the archived items
func (s *QueuedTransactionService) GetAllQueuedTransactionsByUid(c core.Context, uid int64, includeArchived bool) ([]*models.QueuedTransaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var queuedTransactions []*models.QueuedTransaction
	sess := s.UserDataDB(uid).NewSession(c)

	if includeArchived {
		// An archived item past its retention is already gone as far as the user is concerned, even
		// when the purge job has not removed the row yet.
		expirationUnixTime := time.Now().Unix() - models.QueuedTransactionArchivedRetentionHours*60*60
		sess = sess.Where("uid=? AND deleted=? AND (status=? OR archived_unix_time>?)", uid, false, models.QUEUED_TRANSACTION_STATUS_PENDING, expirationUnixTime)
	} else {
		sess = sess.Where("uid=? AND deleted=? AND status=?", uid, false, models.QUEUED_TRANSACTION_STATUS_PENDING)
	}

	err := sess.OrderBy("status asc, transaction_time desc, created_unix_time desc").Find(&queuedTransactions)

	return queuedTransactions, err
}

// GetQueuedTransactionById returns a queued transaction model of the user by the queued transaction id
func (s *QueuedTransactionService) GetQueuedTransactionById(c core.Context, uid int64, queuedTransactionId int64) (*models.QueuedTransaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if queuedTransactionId <= 0 {
		return nil, errs.ErrQueuedTransactionIdInvalid
	}

	queuedTransaction := &models.QueuedTransaction{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(queuedTransactionId).Where("uid=? AND deleted=?", uid, false).Get(queuedTransaction)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrQueuedTransactionNotFound
	}

	return queuedTransaction, nil
}

// GetQueuedTransactionCountsByUid returns how many pending and archived items the queue of the user holds
func (s *QueuedTransactionService) GetQueuedTransactionCountsByUid(c core.Context, uid int64) (pendingCount int64, archivedCount int64, err error) {
	if uid <= 0 {
		return 0, 0, errs.ErrUserIdInvalid
	}

	pendingCount, err = s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND status=?", uid, false, models.QUEUED_TRANSACTION_STATUS_PENDING).Count(&models.QueuedTransaction{})

	if err != nil {
		return 0, 0, err
	}

	expirationUnixTime := time.Now().Unix() - models.QueuedTransactionArchivedRetentionHours*60*60
	archivedCount, err = s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND status=? AND archived_unix_time>?", uid, false, models.QUEUED_TRANSACTION_STATUS_ARCHIVED, expirationUnixTime).Count(&models.QueuedTransaction{})

	if err != nil {
		return 0, 0, err
	}

	return pendingCount, archivedCount, nil
}

// CreateQueuedTransaction saves a new queued transaction to database
func (s *QueuedTransactionService) CreateQueuedTransaction(c core.Context, queuedTransaction *models.QueuedTransaction, tagIds []int64) error {
	if queuedTransaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	queuedTransaction.QueuedTransactionId = s.GenerateUuid(uuid.UUID_TYPE_QUEUED_TRANSACTION)

	if queuedTransaction.QueuedTransactionId < 1 {
		return errs.ErrSystemIsBusy
	}

	now := time.Now().Unix()

	queuedTransaction.Deleted = false
	queuedTransaction.Status = models.QUEUED_TRANSACTION_STATUS_PENDING
	queuedTransaction.ArchivedUnixTime = 0
	queuedTransaction.TagIds = strings.Join(utils.Int64ArrayToStringArray(utils.ToUniqueInt64Slice(tagIds)), ",")
	queuedTransaction.CreatedUnixTime = now
	queuedTransaction.UpdatedUnixTime = now

	return s.UserDataDB(queuedTransaction.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(queuedTransaction)
		return err
	})
}

// ArchiveQueuedTransaction marks a queued transaction of the user as archived (rejected), so it is
// purged once the retention period elapses
func (s *QueuedTransactionService) ArchiveQueuedTransaction(c core.Context, uid int64, queuedTransactionId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if queuedTransactionId <= 0 {
		return errs.ErrQueuedTransactionIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.QueuedTransaction{
		Status:           models.QUEUED_TRANSACTION_STATUS_ARCHIVED,
		ArchivedUnixTime: now,
		UpdatedUnixTime:  now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(queuedTransactionId).Cols("status", "archived_unix_time", "updated_unix_time").Where("uid=? AND deleted=? AND status=?", uid, false, models.QUEUED_TRANSACTION_STATUS_PENDING).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrQueuedTransactionNotFound
		}

		return nil
	})
}

// ArchiveAllPendingQueuedTransactions marks every pending item of the queue of the user as archived
// and returns how many were affected
func (s *QueuedTransactionService) ArchiveAllPendingQueuedTransactions(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.QueuedTransaction{
		Status:           models.QUEUED_TRANSACTION_STATUS_ARCHIVED,
		ArchivedUnixTime: now,
		UpdatedUnixTime:  now,
	}

	var updatedRows int64

	err := s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		rows, err := sess.Cols("status", "archived_unix_time", "updated_unix_time").Where("uid=? AND deleted=? AND status=?", uid, false, models.QUEUED_TRANSACTION_STATUS_PENDING).Update(updateModel)

		if err != nil {
			return err
		}

		updatedRows = rows
		return nil
	})

	if err != nil {
		return 0, err
	}

	return updatedRows, nil
}

// RestoreQueuedTransaction moves an archived queued transaction of the user back to the pending queue
func (s *QueuedTransactionService) RestoreQueuedTransaction(c core.Context, uid int64, queuedTransactionId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if queuedTransactionId <= 0 {
		return errs.ErrQueuedTransactionIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.QueuedTransaction{
		Status:           models.QUEUED_TRANSACTION_STATUS_PENDING,
		ArchivedUnixTime: 0,
		UpdatedUnixTime:  now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(queuedTransactionId).Cols("status", "archived_unix_time", "updated_unix_time").Where("uid=? AND deleted=? AND status=?", uid, false, models.QUEUED_TRANSACTION_STATUS_ARCHIVED).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrQueuedTransactionNotFound
		}

		return nil
	})
}

// DeleteQueuedTransaction removes a queued transaction of the user from the queue, which is what
// accepting one does after the real transaction has been created out of it
func (s *QueuedTransactionService) DeleteQueuedTransaction(c core.Context, uid int64, queuedTransactionId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if queuedTransactionId <= 0 {
		return errs.ErrQueuedTransactionIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.QueuedTransaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.ID(queuedTransactionId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrQueuedTransactionNotFound
		}

		return nil
	})
}

// DeleteAllQueuedTransactions removes the whole queue of the user
func (s *QueuedTransactionService) DeleteAllQueuedTransactions(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.QueuedTransaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)
		return err
	})
}

// DeleteAllExpiredArchivedQueuedTransactions removes every archived queued transaction whose
// retention period has elapsed, across all user data databases
func (s *QueuedTransactionService) DeleteAllExpiredArchivedQueuedTransactions(c core.Context, currentUnixTime int64) error {
	expirationUnixTime := currentUnixTime - models.QueuedTransactionArchivedRetentionHours*60*60

	updateModel := &models.QueuedTransaction{
		Deleted:         true,
		DeletedUnixTime: currentUnixTime,
	}

	totalDeletedRows := int64(0)

	for i := 0; i < s.UserDataDBCount(); i++ {
		deletedRows, err := s.UserDataDBByIndex(i).NewSession(c).Cols("deleted", "deleted_unix_time").Where("deleted=? AND status=? AND archived_unix_time>? AND archived_unix_time<=?", false, models.QUEUED_TRANSACTION_STATUS_ARCHIVED, 0, expirationUnixTime).Update(updateModel)

		if err != nil {
			return err
		}

		totalDeletedRows += deletedRows
	}

	if totalDeletedRows > 0 {
		log.Infof(c, "[queued_transactions.DeleteAllExpiredArchivedQueuedTransactions] %d archived queued transactions have been removed", totalDeletedRows)
	}

	return nil
}
