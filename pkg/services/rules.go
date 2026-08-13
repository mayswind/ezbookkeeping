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

// [PLUGIN:rules] RuleService implements CRUD for rules (and their triggers/actions as one
// unit). Self-contained service. See docs/PLUGIN_DESIGN.md.
type RuleService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a rule service singleton instance
var (
	Rules = &RuleService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetAllRulesByUid returns all rules of a user, optionally filtered by group, ordered by display order
func (s *RuleService) GetAllRulesByUid(c core.Context, uid int64, groupId int64) ([]*models.Rule, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	session := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false)

	if groupId > 0 {
		session = session.And("rule_group_id=?", groupId)
	}

	var rules []*models.Rule
	err := session.OrderBy("rule_group_id asc, display_order asc").Find(&rules)

	return rules, err
}

// GetActiveRulesByGroupId returns all active rules of a group (for engine evaluation)
func (s *RuleService) GetActiveRulesByGroupId(c core.Context, uid int64, groupId int64) ([]*models.Rule, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var rules []*models.Rule
	err := s.UserDataDB(uid).NewSession(c).
		Where("uid=? AND deleted=? AND active=? AND rule_group_id=? AND apply_on_create=?", uid, false, true, groupId, true).
		OrderBy("display_order asc").
		Find(&rules)

	return rules, err
}

// GetRuleByRuleId returns a single rule by its id
func (s *RuleService) GetRuleByRuleId(c core.Context, uid int64, ruleId int64) (*models.Rule, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if ruleId <= 0 {
		return nil, errs.ErrRuleIdInvalid
	}

	rule := &models.Rule{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(ruleId).Where("uid=? AND deleted=?", uid, false).Get(rule)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrRuleNotFound
	}

	return rule, nil
}

// GetTriggersByRuleId returns all triggers of a rule ordered by display order
func (s *RuleService) GetTriggersByRuleId(c core.Context, uid int64, ruleId int64) ([]*models.RuleTrigger, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var triggers []*models.RuleTrigger
	err := s.UserDataDB(uid).NewSession(c).
		Where("rule_id=?", ruleId).OrderBy("display_order asc").Find(&triggers)

	return triggers, err
}

// GetActionsByRuleId returns all actions of a rule ordered by display order
func (s *RuleService) GetActionsByRuleId(c core.Context, uid int64, ruleId int64) ([]*models.RuleAction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var actions []*models.RuleAction
	err := s.UserDataDB(uid).NewSession(c).
		Where("rule_id=?", ruleId).OrderBy("display_order asc").Find(&actions)

	return actions, err
}

// GetMaxDisplayOrder returns the max display order among the user's rules in a group
func (s *RuleService) GetMaxDisplayOrder(c core.Context, uid int64, groupId int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	rule := &models.Rule{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "rule_group_id", "display_order").
		Where("uid=? AND deleted=? AND rule_group_id=?", uid, false, groupId).OrderBy("display_order desc").Limit(1).Get(rule)

	if err != nil {
		return 0, err
	}

	if has {
		return rule.DisplayOrder, nil
	}
	return 0, nil
}

// CreateRuleWithDetails saves a new rule together with its triggers and actions in one transaction
func (s *RuleService) CreateRuleWithDetails(c core.Context, uid int64, rule *models.Rule, triggers []*models.RuleTrigger, actions []*models.RuleAction) (*models.Rule, []*models.RuleTrigger, []*models.RuleAction, error) {
	if uid <= 0 {
		return nil, nil, nil, errs.ErrUserIdInvalid
	}

	if rule.Name == "" {
		return nil, nil, nil, errs.ErrRuleNameEmpty
	}

	if rule.RuleGroupId <= 0 {
		return nil, nil, nil, errs.ErrRuleGroupNotSpecified
	}

	if len(triggers) < 1 {
		return nil, nil, nil, errs.ErrRuleNoTriggers
	}

	// Validate trigger/action types
	for i := 0; i < len(triggers); i++ {
		if !models.IsValidRuleTriggerType(triggers[i].TriggerType) {
			return nil, nil, nil, errs.ErrRuleTriggerTypeInvalid
		}
	}

	for i := 0; i < len(actions); i++ {
		if !models.IsValidRuleActionType(actions[i].ActionType) {
			return nil, nil, nil, errs.ErrRuleActionTypeInvalid
		}
	}

	// Generate the rule id and all trigger/action ids up front so the caller gets them back.
	rule.RuleId = s.GenerateUuid(uuid.UUID_TYPE_RULE)

	if rule.RuleId < 1 {
		return nil, nil, nil, errs.ErrSystemIsBusy
	}

	now := time.Now().Unix()
	rule.Uid = uid
	rule.Deleted = false
	rule.CreatedUnixTime = now
	rule.UpdatedUnixTime = now

	// Need N+M uuids for triggers+actions
	totalDetails := uint16(len(triggers) + len(actions))
	detailIds := s.GenerateUuids(uuid.UUID_TYPE_RULE, totalDetails)

	if len(detailIds) < int(totalDetails) {
		return nil, nil, nil, errs.ErrSystemIsBusy
	}

	idx := 0
	for i := 0; i < len(triggers); i++ {
		triggers[i].TriggerId = detailIds[idx]
		triggers[i].RuleId = rule.RuleId
		triggers[i].DisplayOrder = int32(i + 1)
		triggers[i].CreatedUnixTime = now
		triggers[i].UpdatedUnixTime = now
		idx++
	}

	for i := 0; i < len(actions); i++ {
		actions[i].ActionId = detailIds[idx]
		actions[i].RuleId = rule.RuleId
		actions[i].DisplayOrder = int32(i + 1)
		actions[i].CreatedUnixTime = now
		actions[i].UpdatedUnixTime = now
		idx++
	}

	err := s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		if _, e := sess.Insert(rule); e != nil {
			return e
		}

		if len(triggers) > 0 {
			if _, e := sess.Insert(triggers); e != nil {
				return e
			}
		}

		if len(actions) > 0 {
			if _, e := sess.Insert(actions); e != nil {
				return e
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return rule, triggers, actions, nil
}

// ModifyRuleWithDetails replaces an existed rule and all its triggers/actions
func (s *RuleService) ModifyRuleWithDetails(c core.Context, uid int64, rule *models.Rule, triggers []*models.RuleTrigger, actions []*models.RuleAction) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if rule.RuleId <= 0 {
		return errs.ErrRuleIdInvalid
	}

	if rule.Name == "" {
		return errs.ErrRuleNameEmpty
	}

	if rule.RuleGroupId <= 0 {
		return errs.ErrRuleGroupNotSpecified
	}

	if len(triggers) < 1 {
		return errs.ErrRuleNoTriggers
	}

	for i := 0; i < len(triggers); i++ {
		if !models.IsValidRuleTriggerType(triggers[i].TriggerType) {
			return errs.ErrRuleTriggerTypeInvalid
		}
	}

	for i := 0; i < len(actions); i++ {
		if !models.IsValidRuleActionType(actions[i].ActionType) {
			return errs.ErrRuleActionTypeInvalid
		}
	}

	now := time.Now().Unix()
	rule.Uid = uid
	rule.UpdatedUnixTime = now

	// Generate fresh ids for the new triggers/actions.
	totalDetails := uint16(len(triggers) + len(actions))
	detailIds := s.GenerateUuids(uuid.UUID_TYPE_RULE, totalDetails)

	if len(detailIds) < int(totalDetails) {
		return errs.ErrSystemIsBusy
	}

	idx := 0
	for i := 0; i < len(triggers); i++ {
		triggers[i].TriggerId = detailIds[idx]
		triggers[i].RuleId = rule.RuleId
		triggers[i].DisplayOrder = int32(i + 1)
		triggers[i].CreatedUnixTime = now
		triggers[i].UpdatedUnixTime = now
		idx++
	}

	for i := 0; i < len(actions); i++ {
		actions[i].ActionId = detailIds[idx]
		actions[i].RuleId = rule.RuleId
		actions[i].DisplayOrder = int32(i + 1)
		actions[i].CreatedUnixTime = now
		actions[i].UpdatedUnixTime = now
		idx++
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Update the rule row
		updatedRows, e := sess.ID(rule.RuleId).Cols("rule_group_id", "name", "comment", "active", "strict", "stop_processing", "apply_on_create", "apply_on_update", "updated_unix_time").
			Where("uid=? AND deleted=?", uid, false).Update(rule)

		if e != nil {
			return e
		} else if updatedRows < 1 {
			return errs.ErrRuleNotFound
		}

		// Delete old triggers/actions, insert new ones
		if _, e := sess.Where("rule_id=?", rule.RuleId).Delete(&models.RuleTrigger{}); e != nil {
			return e
		}

		if _, e := sess.Where("rule_id=?", rule.RuleId).Delete(&models.RuleAction{}); e != nil {
			return e
		}

		if len(triggers) > 0 {
			if _, e := sess.Insert(triggers); e != nil {
				return e
			}
		}

		if len(actions) > 0 {
			if _, e := sess.Insert(actions); e != nil {
				return e
			}
		}

		return nil
	})
}

// DeleteRule soft-deletes a rule. Its triggers/actions are hard-deleted since they have no
// Deleted flag (they are always owned by exactly one rule and replaced wholesale on edit).
func (s *RuleService) DeleteRule(c core.Context, uid int64, ruleId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if ruleId <= 0 {
		return errs.ErrRuleIdInvalid
	}

	now := time.Now().Unix()

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(ruleId).Cols("deleted", "deleted_unix_time").
			Where("uid=? AND deleted=?", uid, false).Update(&models.Rule{Deleted: true, DeletedUnixTime: now})

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrRuleNotFound
		}

		if _, err = sess.Where("rule_id=?", ruleId).Delete(&models.RuleTrigger{}); err != nil {
			return err
		}

		_, err = sess.Where("rule_id=?", ruleId).Delete(&models.RuleAction{})
		return err
	})
}
