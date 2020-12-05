package services

import (
	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/uuid"
	"time"
	"xorm.io/xorm"
)

type TransactionTagService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

var (
	TransactionTags = &TransactionTagService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

func (s *TransactionTagService) GetAllTagsByUid(uid int64) ([]*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tags []*models.TransactionTag
	err := s.UserDataDB(uid).Where("uid=?", uid).OrderBy("display_order asc").Find(&tags)

	return tags, err
}

func (s *TransactionTagService) GetTagByTagId(uid int64, tagId int64) (*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if tagId <= 0 {
		return nil, errs.ErrTransactionTagIdInvalid
	}

	tag := &models.TransactionTag{}
	has, err := s.UserDataDB(uid).Where("uid=? AND tag_id=?", uid, tagId).Get(tag)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionTagNotFound
	}

	return tag, nil
}

func (s *TransactionTagService) GetMaxDisplayOrder(uid int64) (int, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	tag := &models.TransactionTag{}
	has, err := s.UserDataDB(uid).Cols("uid", "display_order").Where("uid=?", uid).OrderBy("display_order desc").Limit(1).Get(tag)

	if err != nil {
		return 0, err
	}

	if has {
		return tag.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

func (s *TransactionTagService) CreateTag(tag *models.TransactionTag) error {
	if tag.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsTagName(tag.Uid, tag.Name)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrTransactionTagNameAlreadyExists
	}

	tag.TagId = s.GenerateUuid(uuid.UUID_TYPE_TAG)

	tag.CreatedUnixTime = time.Now().Unix()
	tag.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tag.Uid).DoTransaction(func(sess *xorm.Session) error {
		_, err := sess.Insert(tag)
		return err
	})
}

func (s *TransactionTagService) ModifyTag(tag *models.TransactionTag) error {
	if tag.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsTagName(tag.Uid, tag.Name)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrTransactionTagNameAlreadyExists
	}

	tag.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tag.Uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("name", "updated_unix_time").Where("tag_id=? AND uid=?", tag.TagId, tag.Uid).Update(tag)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionTagNotFound
		}

		return err
	})
}

func (s *TransactionTagService) HideTag(uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTag{
		Hidden: hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").In("tag_id", ids).Where("uid=?", uid).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionTagNotFound
		}

		return err
	})
}

func (s *TransactionTagService) ModifyTagDisplayOrders(uid int64, tags []*models.TransactionTag) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(tags); i++ {
		tags[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(tags); i++ {
			tag := tags[i]
			updatedRows, err := sess.Cols("display_order", "updated_unix_time").Where("tag_id=? AND uid=?", tag.TagId, uid).Update(tag)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrTransactionTagNotFound
			}
		}

		return nil
	})
}

func (s *TransactionTagService) DeleteTags(uid int64, ids []int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	tag := &models.TransactionTag{}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		deletedRows, err := sess.In("tag_id", ids).Where("uid=?", uid).Delete(tag)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionTagNotFound
		}

		return err
	})
}

func (s *TransactionTagService) ExistsTagName(uid int64, name string) (bool, error) {
	if name == "" {
		return false, errs.ErrTransactionTagNameIsEmpty
	}

	return s.UserDB().Cols("name").Where("uid=? AND name=?", uid, name).Exist(&models.TransactionTag{})
}
