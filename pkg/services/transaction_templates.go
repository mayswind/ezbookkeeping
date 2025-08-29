package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// TransactionTemplateService represents transaction template service
type TransactionTemplateService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a transaction template service singleton instance
var (
	TransactionTemplates = &TransactionTemplateService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetTotalNormalTemplateCountByUid returns total normal template count of user
func (s *TransactionTemplateService) GetTotalNormalTemplateCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND template_type=?", uid, false, models.TRANSACTION_TEMPLATE_TYPE_NORMAL).Count(&models.TransactionTemplate{})

	return count, err
}

// GetTotalScheduledTemplateCountByUid returns total scheduled transaction count of user
func (s *TransactionTemplateService) GetTotalScheduledTemplateCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND template_type=?", uid, false, models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE).Count(&models.TransactionTemplate{})

	return count, err
}

// GetAllTemplatesByUid returns all transaction template models of user
func (s *TransactionTemplateService) GetAllTemplatesByUid(c core.Context, uid int64, templateType models.TransactionTemplateType) ([]*models.TransactionTemplate, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var templates []*models.TransactionTemplate
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND template_type=?", uid, false, templateType).Find(&templates)

	return templates, err
}

// GetTemplateByTemplateId returns a transaction template model according to transaction template id
func (s *TransactionTemplateService) GetTemplateByTemplateId(c core.Context, uid int64, templateId int64) (*models.TransactionTemplate, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if templateId <= 0 {
		return nil, errs.ErrTransactionTemplateIdInvalid
	}

	template := &models.TransactionTemplate{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(templateId).Where("uid=? AND deleted=?", uid, false).Get(template)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionTemplateNotFound
	}

	return template, nil
}

// GetMaxDisplayOrder returns the max display order
func (s *TransactionTemplateService) GetMaxDisplayOrder(c core.Context, uid int64, templateType models.TransactionTemplateType) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	template := &models.TransactionTemplate{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=? AND template_type=?", uid, false, templateType).OrderBy("display_order desc").Limit(1).Get(template)

	if err != nil {
		return 0, err
	}

	if has {
		return template.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreateTemplate saves a new transaction template model to database
func (s *TransactionTemplateService) CreateTemplate(c core.Context, template *models.TransactionTemplate) error {
	if template.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	template.TemplateId = s.GenerateUuid(uuid.UUID_TYPE_TEMPLATE)

	if template.TemplateId < 1 {
		return errs.ErrSystemIsBusy
	}

	template.Deleted = false
	template.CreatedUnixTime = time.Now().Unix()
	template.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(template.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		err := s.isTemplateValid(sess, template)

		if err != nil {
			return err
		}

		_, err = sess.Insert(template)
		return err
	})
}

// ModifyTemplate saves an existed transaction template model to database
func (s *TransactionTemplateService) ModifyTemplate(c core.Context, template *models.TransactionTemplate) error {
	if template.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	template.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(template.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		err := s.isTemplateValid(sess, template)

		if err != nil {
			return err
		}

		updatedRows, err := sess.ID(template.TemplateId).Cols("name", "type", "category_id", "account_id", "scheduled_frequency_type", "scheduled_frequency", "scheduled_start_time", "scheduled_end_time", "scheduled_at", "scheduled_timezone_utc_offset", "tag_ids", "amount", "related_account_id", "related_account_amount", "hide_amount", "comment", "updated_unix_time").Where("uid=? AND deleted=?", template.Uid, false).Update(template)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionTemplateNotFound
		}

		return err
	})
}

// HideTemplate updates hidden field of given transaction templates
func (s *TransactionTemplateService) HideTemplate(c core.Context, uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTemplate{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).In("template_id", ids).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionTemplateNotFound
		}

		return err
	})
}

// ModifyTemplateDisplayOrders updates display order of given transaction templates
func (s *TransactionTemplateService) ModifyTemplateDisplayOrders(c core.Context, uid int64, templates []*models.TransactionTemplate) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(templates); i++ {
		templates[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(templates); i++ {
			template := templates[i]
			updatedRows, err := sess.ID(template.TemplateId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(template)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrTransactionTemplateNotFound
			}
		}

		return nil
	})
}

// DeleteTemplate deletes an existed transaction template from database
func (s *TransactionTemplateService) DeleteTemplate(c core.Context, uid int64, templateId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTemplate{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.ID(templateId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionTemplateNotFound
		}

		return err
	})
}

// DeleteAllTemplates deletes all existed transaction templates from database
func (s *TransactionTemplateService) DeleteAllTemplates(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTemplate{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		return nil
	})
}

func (s *TransactionTemplateService) isTemplateValid(sess *xorm.Session, template *models.TransactionTemplate) error {
	// check accounts are valid
	sourceAccount := &models.Account{}
	destinationAccount := &models.Account{}
	has, err := sess.ID(template.AccountId).Where("uid=? AND deleted=?", template.Uid, false).Get(sourceAccount)

	if err != nil {
		return err
	} else if !has {
		return errs.ErrSourceAccountNotFound
	}

	if sourceAccount.Hidden {
		return errs.ErrCannotUseHiddenAccount
	}

	if template.Type == models.TRANSACTION_TYPE_TRANSFER {
		if template.RelatedAccountId <= 0 {
			return errs.ErrAccountIdInvalid
		} else {
			has, err = sess.ID(template.RelatedAccountId).Where("uid=? AND deleted=?", template.Uid, false).Get(destinationAccount)

			if err != nil {
				return err
			} else if !has {
				return errs.ErrDestinationAccountNotFound
			}

			if destinationAccount.Hidden {
				return errs.ErrCannotUseHiddenAccount
			}
		}
	}

	if sourceAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS || (destinationAccount != nil && destinationAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS) {
		return errs.ErrCannotAddTransactionToParentAccount
	}

	// check category is valid
	category := &models.TransactionCategory{}
	has, err = sess.ID(template.CategoryId).Where("uid=? AND deleted=?", template.Uid, false).Get(category)

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

	if (template.Type == models.TRANSACTION_TYPE_INCOME && category.Type != models.CATEGORY_TYPE_INCOME) ||
		(template.Type == models.TRANSACTION_TYPE_EXPENSE && category.Type != models.CATEGORY_TYPE_EXPENSE) ||
		(template.Type == models.TRANSACTION_TYPE_TRANSFER && category.Type != models.CATEGORY_TYPE_TRANSFER) {
		return errs.ErrTransactionCategoryTypeInvalid
	}

	// check tags are valid
	tagIds := template.GetTagIds()
	var tags []*models.TransactionTag
	err = sess.Where("uid=? AND deleted=?", template.Uid, false).In("tag_id", tagIds).Find(&tags)

	if err != nil {
		return err
	} else if len(tags) < len(tagIds) {
		return errs.ErrTransactionTagNotFound
	}

	for i := 0; i < len(tags); i++ {
		if tags[i].Hidden {
			return errs.ErrCannotUseHiddenTransactionTag
		}
	}

	return nil
}
