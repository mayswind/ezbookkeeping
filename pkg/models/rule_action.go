package models

// [PLUGIN:rules] Rule action — a mutation applied to a transaction create request when
// its rule's triggers match. A rule has many actions, executed in order. See
// docs/PLUGIN_DESIGN.md.

// RuleActionType represents the kind of mutation an action performs.
type RuleActionType string

// Supported rule action types for v1.
const (
	RULE_ACTION_SET_CATEGORY     RuleActionType = "set_category"
	RULE_ACTION_CLEAR_CATEGORY   RuleActionType = "clear_category"
	RULE_ACTION_ADD_TAG          RuleActionType = "add_tag"
	RULE_ACTION_REMOVE_TAG       RuleActionType = "remove_tag"
	RULE_ACTION_REMOVE_ALL_TAGS  RuleActionType = "remove_all_tags"
	RULE_ACTION_SET_DESCRIPTION  RuleActionType = "set_description"
	RULE_ACTION_APPEND_TO_DESC   RuleActionType = "append_to_description"
	RULE_ACTION_PREPEND_TO_DESC  RuleActionType = "prepend_to_description"
	RULE_ACTION_SET_AMOUNT       RuleActionType = "set_amount"
	RULE_ACTION_SET_SOURCE_ACCT  RuleActionType = "set_source_account"
)

// AllRuleActionTypes returns all supported action types (used for validation + UI).
func AllRuleActionTypes() []RuleActionType {
	return []RuleActionType{
		RULE_ACTION_SET_CATEGORY,
		RULE_ACTION_CLEAR_CATEGORY,
		RULE_ACTION_ADD_TAG,
		RULE_ACTION_REMOVE_TAG,
		RULE_ACTION_REMOVE_ALL_TAGS,
		RULE_ACTION_SET_DESCRIPTION,
		RULE_ACTION_APPEND_TO_DESC,
		RULE_ACTION_PREPEND_TO_DESC,
		RULE_ACTION_SET_AMOUNT,
		RULE_ACTION_SET_SOURCE_ACCT,
	}
}

// IsValidRuleActionType reports whether t is a supported action type.
func IsValidRuleActionType(t RuleActionType) bool {
	for _, v := range AllRuleActionTypes() {
		if v == t {
			return true
		}
	}
	return false
}

// RuleAction represents a rule action stored in database
type RuleAction struct {
	ActionId        int64          `xorm:"PK"`
	RuleId          int64          `xorm:"INDEX(IDX_rule_action_rule_id_order) NOT NULL"`
	ActionType      RuleActionType `xorm:"VARCHAR(50) NOT NULL"`
	ActionValue     string         `xorm:"VARCHAR(255) NOT NULL"`
	DisplayOrder    int32          `xorm:"INDEX(IDX_rule_action_rule_id_order) NOT NULL"`
	StopProcessing  bool           `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}

// RuleActionInfoResponse represents a view-object of a rule action
type RuleActionInfoResponse struct {
	Id             int64          `json:"id,string"`
	ActionType     RuleActionType `json:"actionType"`
	ActionValue    string         `json:"actionValue"`
	DisplayOrder   int32          `json:"displayOrder"`
	StopProcessing bool           `json:"stopProcessing"`
}

// ToRuleActionInfoResponse returns a view-object according to the database model
func (a *RuleAction) ToRuleActionInfoResponse() *RuleActionInfoResponse {
	return &RuleActionInfoResponse{
		Id:             a.ActionId,
		ActionType:     a.ActionType,
		ActionValue:    a.ActionValue,
		DisplayOrder:   a.DisplayOrder,
		StopProcessing: a.StopProcessing,
	}
}
