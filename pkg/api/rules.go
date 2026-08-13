package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// [PLUGIN:rules] RulesApi implements the rule group and rule HTTP API.
// Self-contained handler module; see docs/PLUGIN_DESIGN.md.
type RulesApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	ruleGroups *services.RuleGroupService
	rules      *services.RuleService
}

// Initialize a rules api singleton instance
var (
	Rules = &RulesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		ruleGroups: services.RuleGroups,
		rules:      services.Rules,
	}
)

// ---------- Rule Groups ----------

// RuleGroupListHandler returns all rule groups of the current user
func (a *RulesApi) RuleGroupListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	groups, e := a.ruleGroups.GetAllRuleGroupsByUid(c, uid)

	if e != nil {
		log.Errorf(c, "[rules.RuleGroupListHandler] failed to get rule groups for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	result := make([]*models.RuleGroupInfoResponse, 0, len(groups))
	for i := 0; i < len(groups); i++ {
		result = append(result, groups[i].ToRuleGroupInfoResponse())
	}

	return result, nil
}

// RuleGroupGetHandler returns one rule group
func (a *RulesApi) RuleGroupGetHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.RuleGroupGetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Warnf(c, "[rules.RuleGroupGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	group, e := a.ruleGroups.GetRuleGroupByGroupId(c, uid, req.Id)

	if e != nil {
		log.Errorf(c, "[rules.RuleGroupGetHandler] failed to get rule group \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	return group.ToRuleGroupInfoResponse(), nil
}

// RuleGroupCreateHandler creates a new rule group
func (a *RulesApi) RuleGroupCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var createReq models.RuleGroupCreateRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		log.Warnf(c, "[rules.RuleGroupCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	maxOrder, e := a.ruleGroups.GetMaxDisplayOrder(c, uid)
	if e != nil {
		log.Errorf(c, "[rules.RuleGroupCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	group := &models.RuleGroup{
		Uid:            uid,
		Name:           createReq.Name,
		Comment:        createReq.Comment,
		DisplayOrder:   maxOrder + 1,
		Active:         createReq.Active,
		StopProcessing: createReq.StopProcessing,
	}

	if e := a.ruleGroups.CreateRuleGroup(c, group); e != nil {
		log.Errorf(c, "[rules.RuleGroupCreateHandler] failed to create rule group for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	log.Infof(c, "[rules.RuleGroupCreateHandler] user \"uid:%d\" created rule group \"id:%d\"", uid, group.GroupId)
	return group.ToRuleGroupInfoResponse(), nil
}

// RuleGroupModifyHandler modifies an existed rule group
func (a *RulesApi) RuleGroupModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var modifyReq models.RuleGroupModifyRequest
	if err := c.ShouldBindJSON(&modifyReq); err != nil {
		log.Warnf(c, "[rules.RuleGroupModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	group := &models.RuleGroup{
		GroupId:        modifyReq.Id,
		Uid:            uid,
		Name:           modifyReq.Name,
		Comment:        modifyReq.Comment,
		Active:         modifyReq.Active,
		StopProcessing: modifyReq.StopProcessing,
	}

	if e := a.ruleGroups.ModifyRuleGroup(c, group); e != nil {
		log.Errorf(c, "[rules.RuleGroupModifyHandler] failed to modify rule group \"id:%d\" for user \"uid:%d\", because %s", modifyReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	updated, e := a.ruleGroups.GetRuleGroupByGroupId(c, uid, modifyReq.Id)
	if e != nil {
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	return updated.ToRuleGroupInfoResponse(), nil
}

// RuleGroupDeleteHandler deletes a rule group and its rules
func (a *RulesApi) RuleGroupDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var deleteReq models.RuleGroupDeleteRequest
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Warnf(c, "[rules.RuleGroupDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	if e := a.ruleGroups.DeleteRuleGroup(c, uid, deleteReq.Id); e != nil {
		log.Errorf(c, "[rules.RuleGroupDeleteHandler] failed to delete rule group \"id:%d\" for user \"uid:%d\", because %s", deleteReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	log.Infof(c, "[rules.RuleGroupDeleteHandler] user \"uid:%d\" deleted rule group \"id:%d\"", uid, deleteReq.Id)
	return true, nil
}

// ---------- Rules ----------

// RuleListHandler returns all rules (optionally filtered by group) with their triggers/actions
func (a *RulesApi) RuleListHandler(c *core.WebContext) (any, *errs.Error) {
	var listReq models.RuleListByGroupRequest
	_ = c.ShouldBindQuery(&listReq)

	uid := c.GetCurrentUid()
	rules, e := a.rules.GetAllRulesByUid(c, uid, listReq.GroupId)

	if e != nil {
		log.Errorf(c, "[rules.RuleListHandler] failed to get rules for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	result := make([]*models.RuleInfoResponse, 0, len(rules))
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		triggers, _ := a.rules.GetTriggersByRuleId(c, uid, rule.RuleId)
		actions, _ := a.rules.GetActionsByRuleId(c, uid, rule.RuleId)
		result = append(result, rule.ToRuleInfoResponse(triggers, actions))
	}

	return result, nil
}

// RuleGetHandler returns one rule with its triggers and actions
func (a *RulesApi) RuleGetHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.RuleGetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Warnf(c, "[rules.RuleGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	rule, e := a.rules.GetRuleByRuleId(c, uid, req.Id)

	if e != nil {
		log.Errorf(c, "[rules.RuleGetHandler] failed to get rule \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	triggers, _ := a.rules.GetTriggersByRuleId(c, uid, rule.RuleId)
	actions, _ := a.rules.GetActionsByRuleId(c, uid, rule.RuleId)

	return rule.ToRuleInfoResponse(triggers, actions), nil
}

// RuleCreateHandler creates a new rule with its triggers and actions
func (a *RulesApi) RuleCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var createReq models.RuleCreateRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		log.Warnf(c, "[rules.RuleCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	// Verify the rule group exists and belongs to the user.
	if _, e := a.ruleGroups.GetRuleGroupByGroupId(c, uid, createReq.RuleGroupId); e != nil {
		log.Warnf(c, "[rules.RuleCreateHandler] rule group \"id:%d\" not found for user \"uid:%d\"", createReq.RuleGroupId, uid)
		return nil, errs.ErrRuleGroupNotFound
	}

	maxOrder, e := a.rules.GetMaxDisplayOrder(c, uid, createReq.RuleGroupId)
	if e != nil {
		log.Errorf(c, "[rules.RuleCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	rule := &models.Rule{
		Uid:            uid,
		RuleGroupId:    createReq.RuleGroupId,
		Name:           createReq.Name,
		Comment:        createReq.Comment,
		DisplayOrder:   maxOrder + 1,
		Active:         createReq.Active,
		Strict:         createReq.Strict,
		StopProcessing: createReq.StopProcessing,
		ApplyOnCreate:  createReq.ApplyOnCreate,
		ApplyOnUpdate:  createReq.ApplyOnUpdate,
	}

	triggers := buildTriggersFromInput(createReq.Triggers)
	actions := buildActionsFromInput(createReq.Actions)

	rule, triggers, actions, e = a.rules.CreateRuleWithDetails(c, uid, rule, triggers, actions)

	if e != nil {
		log.Errorf(c, "[rules.RuleCreateHandler] failed to create rule for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	log.Infof(c, "[rules.RuleCreateHandler] user \"uid:%d\" created rule \"id:%d\"", uid, rule.RuleId)
	return rule.ToRuleInfoResponse(triggers, actions), nil
}

// RuleModifyHandler modifies an existed rule (replacing its triggers and actions)
func (a *RulesApi) RuleModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var modifyReq models.RuleModifyRequest
	if err := c.ShouldBindJSON(&modifyReq); err != nil {
		log.Warnf(c, "[rules.RuleModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	// Preserve the existing display order.
	existing, e := a.rules.GetRuleByRuleId(c, uid, modifyReq.Id)
	if e != nil {
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	rule := &models.Rule{
		RuleId:         modifyReq.Id,
		Uid:            uid,
		RuleGroupId:    modifyReq.RuleGroupId,
		Name:           modifyReq.Name,
		Comment:        modifyReq.Comment,
		DisplayOrder:   existing.DisplayOrder,
		Active:         modifyReq.Active,
		Strict:         modifyReq.Strict,
		StopProcessing: modifyReq.StopProcessing,
		ApplyOnCreate:  modifyReq.ApplyOnCreate,
		ApplyOnUpdate:  modifyReq.ApplyOnUpdate,
	}

	triggers := buildTriggersFromInput(modifyReq.Triggers)
	actions := buildActionsFromInput(modifyReq.Actions)

	if e := a.rules.ModifyRuleWithDetails(c, uid, rule, triggers, actions); e != nil {
		log.Errorf(c, "[rules.RuleModifyHandler] failed to modify rule \"id:%d\" for user \"uid:%d\", because %s", modifyReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	updated, _ := a.rules.GetRuleByRuleId(c, uid, modifyReq.Id)
	updatedTriggers, _ := a.rules.GetTriggersByRuleId(c, uid, modifyReq.Id)
	updatedActions, _ := a.rules.GetActionsByRuleId(c, uid, modifyReq.Id)

	log.Infof(c, "[rules.RuleModifyHandler] user \"uid:%d\" modified rule \"id:%d\"", uid, modifyReq.Id)
	return updated.ToRuleInfoResponse(updatedTriggers, updatedActions), nil
}

// RuleDeleteHandler deletes a rule
func (a *RulesApi) RuleDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var deleteReq models.RuleDeleteRequest
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Warnf(c, "[rules.RuleDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	if e := a.rules.DeleteRule(c, uid, deleteReq.Id); e != nil {
		log.Errorf(c, "[rules.RuleDeleteHandler] failed to delete rule \"id:%d\" for user \"uid:%d\", because %s", deleteReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	log.Infof(c, "[rules.RuleDeleteHandler] user \"uid:%d\" deleted rule \"id:%d\"", uid, deleteReq.Id)
	return true, nil
}

// ---------- helpers ----------

func buildTriggersFromInput(inputs []*models.RuleTriggerInput) []*models.RuleTrigger {
	triggers := make([]*models.RuleTrigger, len(inputs))
	for i := 0; i < len(inputs); i++ {
		triggers[i] = &models.RuleTrigger{
			TriggerType:    inputs[i].TriggerType,
			TriggerValue:   inputs[i].TriggerValue,
			Prohibited:     inputs[i].Prohibited,
			StopProcessing: inputs[i].StopProcessing,
		}
	}
	return triggers
}

func buildActionsFromInput(inputs []*models.RuleActionInput) []*models.RuleAction {
	actions := make([]*models.RuleAction, len(inputs))
	for i := 0; i < len(inputs); i++ {
		actions[i] = &models.RuleAction{
			ActionType:     inputs[i].ActionType,
			ActionValue:    inputs[i].ActionValue,
			StopProcessing: inputs[i].StopProcessing,
		}
	}
	return actions
}
