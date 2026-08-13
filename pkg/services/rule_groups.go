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

// [PLUGIN:rules] RuleGroupService implements CRUD for rule groups.
// Self-contained service following the existing TransactionCategoryService pattern.
// See docs/PLUGIN_DESIGN.md.
type RuleGroupService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a rule group service singleton instance
var (
	RuleGroups = &RuleGroupService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetAllRuleGroupsByUid returns all rule groups of a user ordered by display order
func (s *RuleGroupService) GetAllRuleGroupsByUid(c core.Context, uid int64) ([]*models.RuleGroup, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var groups []*models.RuleGroup
	err := s.UserDataDB(uid).NewSession(c).
		Where("uid=? AND deleted=?", uid, false).
		OrderBy("display_order asc").
		Find(&groups)

	return groups, err
}

// GetAllActiveRuleGroupsByUid returns all active rule groups of a user (for engine evaluation)
func (s *RuleGroupService) GetAllActiveRuleGroupsByUid(c core.Context, uid int64) ([]*models.RuleGroup, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var groups []*models.RuleGroup
	err := s.UserDataDB(uid).NewSession(c).
		Where("uid=? AND deleted=? AND active=?", uid, false, true).
		OrderBy("display_order asc").
		Find(&groups)

	return groups, err
}

// GetRuleGroupByGroupId returns a single rule group by its id
func (s *RuleGroupService) GetRuleGroupByGroupId(c core.Context, uid int64, groupId int64) (*models.RuleGroup, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if groupId <= 0 {
		return nil, errs.ErrRuleGroupIdInvalid
	}

	group := &models.RuleGroup{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(groupId).Where("uid=? AND deleted=?", uid, false).Get(group)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrRuleGroupNotFound
	}

	return group, nil
}

// GetMaxDisplayOrder returns the max display order among the user's rule groups
func (s *RuleGroupService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	group := &models.RuleGroup{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").
		Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(group)

	if err != nil {
		return 0, err
	}

	if has {
		return group.DisplayOrder, nil
	}
	return 0, nil
}

// CreateRuleGroup saves a new rule group
func (s *RuleGroupService) CreateRuleGroup(c core.Context, group *models.RuleGroup) error {
	if group.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if group.Name == "" {
		return errs.ErrRuleGroupNameEmpty
	}

	group.GroupId = s.GenerateUuid(uuid.UUID_TYPE_RULE)

	if group.GroupId < 1 {
		return errs.ErrSystemIsBusy
	}

	group.Deleted = false
	group.CreatedUnixTime = time.Now().Unix()
	group.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(group.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(group)
		return err
	})
}

// ModifyRuleGroup saves an existed rule group (mutable fields only)
func (s *RuleGroupService) ModifyRuleGroup(c core.Context, group *models.RuleGroup) error {
	if group.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if group.Name == "" {
		return errs.ErrRuleGroupNameEmpty
	}

	group.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(group.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(group.GroupId).Cols("name", "comment", "active", "stop_processing", "updated_unix_time").
			Where("uid=? AND deleted=?", group.Uid, false).Update(group)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrRuleGroupNotFound
		}

		return nil
	})
}

// DeleteRuleGroup soft-deletes a rule group and all rules/triggers/actions within it
func (s *RuleGroupService) DeleteRuleGroup(c core.Context, uid int64, groupId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if groupId <= 0 {
		return errs.ErrRuleGroupIdInvalid
	}

	now := time.Now().Unix()

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Verify ownership
		exists, err := sess.Cols("uid", "deleted", "group_id").Where("uid=? AND deleted=? AND group_id=?", uid, false, groupId).
			Limit(1).Exist(&models.RuleGroup{})

		if err != nil {
			return err
		} else if !exists {
			return errs.ErrRuleGroupNotFound
		}

		// Soft-delete the group
		_, err = sess.ID(groupId).Cols("deleted", "deleted_unix_time").
			Where("uid=? AND deleted=?", uid, false).Update(&models.RuleGroup{Deleted: true, DeletedUnixTime: now})

		if err != nil {
			return err
		}

		// Soft-delete all rules in this group
		_, err = sess.Cols("deleted", "deleted_unix_time").
			Where("uid=? AND deleted=? AND rule_group_id=?", uid, false, groupId).
			Update(&models.Rule{Deleted: true, DeletedUnixTime: now})

		return err
	})
}
