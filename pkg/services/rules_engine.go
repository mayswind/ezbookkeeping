package services

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// [PLUGIN:rules] RulesEngine evaluates the user's rule groups/rules against a
// TransactionCreateRequest and mutates the request DTO in place BEFORE the transaction is
// created. It is invoked from exactly two call-sites in pkg/api/transactions.go
// (TransactionCreateHandler and TransactionImportHandler), gated by config.EnableRulesEngine.
//
// Design: in-memory evaluation (simpler than Firefly III's search-composition approach).
// For each request we build a context map of category/tag/account names once, then evaluate
// every active rule group in order. This touches ZERO core transaction code: it only reads
// the request DTO and writes back to it. See docs/PLUGIN_DESIGN.md.

// ApplyToCreateRequest evaluates all active rules (apply_on_create=true) for the user and
// mutates req in place. Returns nil on success; rule-matching errors for a single rule are
// logged and skipped so one bad rule cannot block transaction creation.
func (s *RuleService) ApplyToCreateRequest(c core.Context, uid int64, req *models.TransactionCreateRequest) (bool, *errs.Error) {
	if uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	if req == nil {
		return false, nil
	}

	// 1. Load active rule groups (ordered) — none means nothing to do.
	groups, err := RuleGroups.GetAllActiveRuleGroupsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[rules_engine.ApplyToCreateRequest] failed to load rule groups for user \"uid:%d\", because %s", uid, err.Error())
		return false, errs.Or(err, errs.ErrOperationFailed)
	}

	if len(groups) == 0 {
		return false, nil
	}

	// 2. Build the evaluation context once (category/tag/account name resolution).
	ctx, err := s.buildEvaluationContext(c, uid)

	if err != nil {
		log.Errorf(c, "[rules_engine.ApplyToCreateRequest] failed to build evaluation context for user \"uid:%d\", because %s", uid, err.Error())
		return false, errs.Or(err, errs.ErrOperationFailed)
	}

	changed := false

	// 3. For each group, load its active rules and evaluate.
	for _, group := range groups {
		rules, e := s.GetActiveRulesByGroupId(c, uid, group.GroupId)

		if e != nil {
			log.Warnf(c, "[rules_engine.ApplyToCreateRequest] failed to load rules for group \"id:%d\", because %s", group.GroupId, e.Error())
			continue
		}

		groupHadMatch := false

		for _, rule := range rules {
			matched, matchErr := s.evaluateRuleTriggers(c, uid, rule, req, ctx)

			if matchErr != nil {
				log.Warnf(c, "[rules_engine.ApplyToCreateRequest] failed to evaluate rule \"id:%d\", because %s", rule.RuleId, matchErr.Error())
				continue
			}

			if !matched {
				continue
			}

			groupHadMatch = true

			// Apply the rule's actions in order.
			actions, _ := s.GetActionsByRuleId(c, uid, rule.RuleId)

			for _, action := range actions {
				if applied := s.applyAction(c, uid, action, req, ctx); applied {
					changed = true
				}

				if action.StopProcessing {
					break
				}
			}

			// Rule-level short-circuit: stop processing further rules in this group.
			if rule.StopProcessing {
				break
			}
		}

		// Group-level short-circuit: stop all subsequent groups if a rule matched.
		if groupHadMatch && group.StopProcessing {
			break
		}
	}

	return changed, nil
}

// ruleEvaluationContext holds the name->id / id->name maps needed to evaluate triggers and
// resolve action targets by name (find-or-create).
type ruleEvaluationContext struct {
	categoriesByName map[string]int64                 // lowercase name -> category id
	categoriesById   map[int64]*models.TransactionCategory
	tagsByName       map[string]int64                 // lowercase name -> tag id
	tagIdStrings     map[string]bool                  // set of tag ids already present on the request (as strings)
	accountsById     map[int64]*models.Account
}

func (s *RuleService) buildEvaluationContext(c core.Context, uid int64) (*ruleEvaluationContext, error) {
	ctx := &ruleEvaluationContext{
		categoriesByName: make(map[string]int64),
		categoriesById:   make(map[int64]*models.TransactionCategory),
		tagsByName:       make(map[string]int64),
		tagIdStrings:     make(map[string]bool),
		accountsById:     make(map[int64]*models.Account),
	}

	// Categories (all types, all levels)
	categories, err := TransactionCategories.GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(categories); i++ {
		cat := categories[i]
		ctx.categoriesById[cat.CategoryId] = cat
		ctx.categoriesByName[strings.ToLower(cat.Name)] = cat.CategoryId
	}

	// Tags
	tags, err := TransactionTags.GetAllTagsByUid(c, uid)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		ctx.tagsByName[strings.ToLower(tag.Name)] = tag.TagId
	}

	// Accounts
	accounts, err := Accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(accounts); i++ {
		acc := accounts[i]
		ctx.accountsById[acc.AccountId] = acc
	}

	return ctx, nil
}

// evaluateRuleTriggers returns whether the rule's triggers match the request.
// Strict=true -> ALL triggers must match; Strict=false -> ANY trigger matches.
// Each trigger's match result can be negated by its Prohibited flag.
func (s *RuleService) evaluateRuleTriggers(c core.Context, uid int64, rule *models.Rule, req *models.TransactionCreateRequest, ctx *ruleEvaluationContext) (bool, error) {
	triggers, err := s.GetTriggersByRuleId(c, uid, rule.RuleId)

	if err != nil {
		return false, err
	}

	if len(triggers) == 0 {
		return false, nil
	}

	if rule.Strict {
		// ALL must match (after negation)
		for _, trigger := range triggers {
			matched := s.triggerMatches(trigger, req, ctx)
			matched = matched != trigger.Prohibited // XOR with prohibited
			if !matched {
				return false, nil
			}
		}
		return true, nil
	}

	// ANY must match (after negation)
	for _, trigger := range triggers {
		matched := s.triggerMatches(trigger, req, ctx)
		matched = matched != trigger.Prohibited
		if matched {
			return true, nil
		}
	}
	return false, nil
}

// triggerMatches evaluates a single trigger against the request (without the prohibited flag)
func (s *RuleService) triggerMatches(trigger *models.RuleTrigger, req *models.TransactionCreateRequest, ctx *ruleEvaluationContext) bool {
	switch trigger.TriggerType {
	case models.RULE_TRIGGER_DESCRIPTION_IS:
		return strings.EqualFold(strings.TrimSpace(req.Comment), strings.TrimSpace(trigger.TriggerValue))

	case models.RULE_TRIGGER_DESCRIPTION_CONTAINS:
		return strings.Contains(strings.ToLower(req.Comment), strings.ToLower(trigger.TriggerValue))

	case models.RULE_TRIGGER_AMOUNT_IS:
		return absInt64(req.SourceAmount) == parseAmountToCents(trigger.TriggerValue)

	case models.RULE_TRIGGER_AMOUNT_LESS:
		return absInt64(req.SourceAmount) < parseAmountToCents(trigger.TriggerValue)

	case models.RULE_TRIGGER_AMOUNT_MORE:
		return absInt64(req.SourceAmount) > parseAmountToCents(trigger.TriggerValue)

	case models.RULE_TRIGGER_SOURCE_ACCOUNT_IS:
		name := accountNameById(ctx, req.SourceAccountId)
		return strings.EqualFold(strings.TrimSpace(name), strings.TrimSpace(trigger.TriggerValue))

	case models.RULE_TRIGGER_DESTINATION_ACCOUNT_IS:
		name := accountNameById(ctx, req.DestinationAccountId)
		if name == "" {
			return false
		}
		return strings.EqualFold(strings.TrimSpace(name), strings.TrimSpace(trigger.TriggerValue))

	case models.RULE_TRIGGER_CATEGORY_IS:
		cat := ctx.categoriesById[req.CategoryId]
		if cat == nil {
			return false
		}
		return strings.EqualFold(cat.Name, trigger.TriggerValue)

	case models.RULE_TRIGGER_HAS_NO_CATEGORY:
		return req.CategoryId <= 0

	case models.RULE_TRIGGER_HAS_ANY_CATEGORY:
		return req.CategoryId > 0
	}

	return false
}

// applyAction mutates req according to the action. Returns true if a change was made.
func (s *RuleService) applyAction(c core.Context, uid int64, action *models.RuleAction, req *models.TransactionCreateRequest, ctx *ruleEvaluationContext) bool {
	switch action.ActionType {
	case models.RULE_ACTION_SET_CATEGORY:
		catId, ok := ctx.categoriesByName[strings.ToLower(strings.TrimSpace(action.ActionValue))]
		if !ok {
			// find-or-create (mirrors Firefly III SetCategory)
			catId = s.findOrCreateCategory(c, uid, action.ActionValue, ctx)
		}
		if catId > 0 {
			req.CategoryId = catId
			return true
		}
		return false

	case models.RULE_ACTION_CLEAR_CATEGORY:
		req.CategoryId = 0
		return true

	case models.RULE_ACTION_ADD_TAG:
		tagName := strings.TrimSpace(action.ActionValue)
		if tagName == "" {
			return false
		}
		tagId, ok := ctx.tagsByName[strings.ToLower(tagName)]
		if !ok {
			tagId = s.findOrCreateTag(c, uid, tagName, ctx)
		}
		if tagId > 0 {
			tagIdStr := utils.Int64ToString(tagId)
			if !containsString(req.TagIds, tagIdStr) {
				req.TagIds = append(req.TagIds, tagIdStr)
			}
			return true
		}
		return false

	case models.RULE_ACTION_REMOVE_TAG:
		tagName := strings.TrimSpace(action.ActionValue)
		tagId, ok := ctx.tagsByName[strings.ToLower(tagName)]
		if !ok {
			return false
		}
		tagIdStr := utils.Int64ToString(tagId)
		newTagIds := make([]string, 0, len(req.TagIds))
		changed := false
		for i := 0; i < len(req.TagIds); i++ {
			if req.TagIds[i] == tagIdStr {
				changed = true
				continue
			}
			newTagIds = append(newTagIds, req.TagIds[i])
		}
		if changed {
			req.TagIds = newTagIds
			return true
		}
		return false

	case models.RULE_ACTION_REMOVE_ALL_TAGS:
		if len(req.TagIds) > 0 {
			req.TagIds = []string{}
			return true
		}
		return false

	case models.RULE_ACTION_SET_DESCRIPTION:
		req.Comment = action.ActionValue
		return true

	case models.RULE_ACTION_APPEND_TO_DESC:
		req.Comment = req.Comment + action.ActionValue
		return true

	case models.RULE_ACTION_PREPEND_TO_DESC:
		req.Comment = action.ActionValue + req.Comment
		return true

	case models.RULE_ACTION_SET_AMOUNT:
		cents := parseAmountToCents(action.ActionValue)
		if cents > 0 {
			// Preserve the sign of the original amount
			if req.SourceAmount < 0 {
				req.SourceAmount = -cents
			} else {
				req.SourceAmount = cents
			}
			return true
		}
		return false

	case models.RULE_ACTION_SET_SOURCE_ACCT:
		// Resolve account by name; if found, set the source account id
		for id, acc := range ctx.accountsById {
			if strings.EqualFold(strings.TrimSpace(acc.Name), strings.TrimSpace(action.ActionValue)) {
				req.SourceAccountId = id
				return true
			}
		}
		return false
	}

	return false
}

// findOrCreateCategory resolves a category by name, creating it if missing (expense type by
// default). Returns 0 on failure. The created id is cached in the context.
func (s *RuleService) findOrCreateCategory(c core.Context, uid int64, name string, ctx *ruleEvaluationContext) int64 {
	category := &models.TransactionCategory{
		Uid:              uid,
		Type:             models.CATEGORY_TYPE_EXPENSE,
		ParentCategoryId: models.LevelOneTransactionCategoryParentId,
		Name:             name,
		DisplayOrder:     0,
		Icon:             1,
		Color:            "BDBDBD",
	}

	if err := TransactionCategories.CreateCategory(c, category); err != nil {
		log.Warnf(c, "[rules_engine.findOrCreateCategory] failed to create category %q, because %s", name, err.Error())
		return 0
	}

	ctx.categoriesById[category.CategoryId] = category
	ctx.categoriesByName[strings.ToLower(name)] = category.CategoryId
	return category.CategoryId
}

// findOrCreateTag resolves a tag by name, creating it if missing. Returns 0 on failure.
func (s *RuleService) findOrCreateTag(c core.Context, uid int64, name string, ctx *ruleEvaluationContext) int64 {
	tag := &models.TransactionTag{
		Uid:           uid,
		TagGroupId:    0,
		Name:          name,
		DisplayOrder:  0,
	}

	if err := TransactionTags.CreateTag(c, tag); err != nil {
		log.Warnf(c, "[rules_engine.findOrCreateTag] failed to create tag %q, because %s", name, err.Error())
		return 0
	}

	ctx.tagsByName[strings.ToLower(name)] = tag.TagId
	return tag.TagId
}

// --- helpers ---

func accountNameById(ctx *ruleEvaluationContext, accountId int64) string {
	if acc, ok := ctx.accountsById[accountId]; ok {
		return acc.Name
	}
	return ""
}

// absInt64 returns the absolute value (amounts are stored in minor units; expense amounts
// may be negative depending on transaction type, so we compare magnitudes).
func absInt64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// parseAmountToCents parses a decimal amount string like "12.50" into minor units (cents).
// Returns 0 on parse failure or negative input.
func parseAmountToCents(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	bigVal, err := utils.ParseAmount(s)
	if err != nil {
		return 0
	}

	return bigVal
}

// containsString reports whether the slice contains the value
func containsString(slice []string, value string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			return true
		}
	}
	return false
}
