package errs

import "net/http"

// [PLUGIN:rules] Error codes related to rules engine (groups, rules, triggers, actions)
var (
	ErrRuleGroupIdInvalid      = NewNormalError(NormalSubcategoryRule, 0, http.StatusBadRequest, "rule group id is invalid")
	ErrRuleGroupNotFound       = NewNormalError(NormalSubcategoryRule, 1, http.StatusBadRequest, "rule group not found")
	ErrRuleGroupNameEmpty      = NewNormalError(NormalSubcategoryRule, 2, http.StatusBadRequest, "rule group name is empty")
	ErrRuleIdInvalid           = NewNormalError(NormalSubcategoryRule, 3, http.StatusBadRequest, "rule id is invalid")
	ErrRuleNotFound            = NewNormalError(NormalSubcategoryRule, 4, http.StatusBadRequest, "rule not found")
	ErrRuleNameEmpty           = NewNormalError(NormalSubcategoryRule, 5, http.StatusBadRequest, "rule name is empty")
	ErrRuleGroupNotSpecified   = NewNormalError(NormalSubcategoryRule, 6, http.StatusBadRequest, "rule group is not specified")
	ErrRuleTriggerTypeInvalid  = NewNormalError(NormalSubcategoryRule, 7, http.StatusBadRequest, "rule trigger type is invalid")
	ErrRuleActionTypeInvalid   = NewNormalError(NormalSubcategoryRule, 8, http.StatusBadRequest, "rule action type is invalid")
	ErrRuleNoTriggers          = NewNormalError(NormalSubcategoryRule, 9, http.StatusBadRequest, "rule must have at least one trigger")
	ErrRuleTriggerValueMissing = NewNormalError(NormalSubcategoryRule, 10, http.StatusBadRequest, "rule trigger value is missing")
)
