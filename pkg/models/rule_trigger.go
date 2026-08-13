package models

// [PLUGIN:rules] Rule trigger — a single condition evaluated against a transaction create
// request. A rule has many triggers; whether ALL or ANY must match is controlled by
// Rule.Strict. See docs/PLUGIN_DESIGN.md.

// RuleTriggerType represents the kind of comparison a trigger performs.
type RuleTriggerType string

// Supported rule trigger types for v1. Each has a comparison function in the rules engine.
const (
	RULE_TRIGGER_DESCRIPTION_IS        RuleTriggerType = "description_is"
	RULE_TRIGGER_DESCRIPTION_CONTAINS  RuleTriggerType = "description_contains"
	RULE_TRIGGER_AMOUNT_IS             RuleTriggerType = "amount_is"
	RULE_TRIGGER_AMOUNT_LESS           RuleTriggerType = "amount_less"
	RULE_TRIGGER_AMOUNT_MORE           RuleTriggerType = "amount_more"
	RULE_TRIGGER_SOURCE_ACCOUNT_IS     RuleTriggerType = "source_account_is"
	RULE_TRIGGER_DESTINATION_ACCOUNT_IS RuleTriggerType = "destination_account_is"
	RULE_TRIGGER_CATEGORY_IS           RuleTriggerType = "category_is"
	RULE_TRIGGER_HAS_NO_CATEGORY       RuleTriggerType = "has_no_category"
	RULE_TRIGGER_HAS_ANY_CATEGORY      RuleTriggerType = "has_any_category"
)

// AllRuleTriggerTypes returns all supported trigger types (used for validation + UI).
func AllRuleTriggerTypes() []RuleTriggerType {
	return []RuleTriggerType{
		RULE_TRIGGER_DESCRIPTION_IS,
		RULE_TRIGGER_DESCRIPTION_CONTAINS,
		RULE_TRIGGER_AMOUNT_IS,
		RULE_TRIGGER_AMOUNT_LESS,
		RULE_TRIGGER_AMOUNT_MORE,
		RULE_TRIGGER_SOURCE_ACCOUNT_IS,
		RULE_TRIGGER_DESTINATION_ACCOUNT_IS,
		RULE_TRIGGER_CATEGORY_IS,
		RULE_TRIGGER_HAS_NO_CATEGORY,
		RULE_TRIGGER_HAS_ANY_CATEGORY,
	}
}

// IsValidRuleTriggerType reports whether t is a supported trigger type.
func IsValidRuleTriggerType(t RuleTriggerType) bool {
	for _, v := range AllRuleTriggerTypes() {
		if v == t {
			return true
		}
	}
	return false
}

// RuleTrigger represents a rule trigger stored in database
type RuleTrigger struct {
	TriggerId       int64           `xorm:"PK"`
	RuleId          int64           `xorm:"INDEX(IDX_rule_trigger_rule_id_order) NOT NULL"`
	TriggerType     RuleTriggerType `xorm:"VARCHAR(50) NOT NULL"`
	TriggerValue    string          `xorm:"VARCHAR(255) NOT NULL"`
	DisplayOrder    int32           `xorm:"INDEX(IDX_rule_trigger_rule_id_order) NOT NULL"`
	Prohibited      bool            `xorm:"NOT NULL"` // negate the match
	StopProcessing  bool            `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}

// RuleTriggerInfoResponse represents a view-object of a rule trigger
type RuleTriggerInfoResponse struct {
	Id             int64           `json:"id,string"`
	TriggerType    RuleTriggerType `json:"triggerType"`
	TriggerValue   string          `json:"triggerValue"`
	DisplayOrder   int32           `json:"displayOrder"`
	Prohibited     bool            `json:"prohibited"`
	StopProcessing bool            `json:"stopProcessing"`
}

// ToRuleTriggerInfoResponse returns a view-object according to the database model
func (t *RuleTrigger) ToRuleTriggerInfoResponse() *RuleTriggerInfoResponse {
	return &RuleTriggerInfoResponse{
		Id:             t.TriggerId,
		TriggerType:    t.TriggerType,
		TriggerValue:   t.TriggerValue,
		DisplayOrder:   t.DisplayOrder,
		Prohibited:     t.Prohibited,
		StopProcessing: t.StopProcessing,
	}
}
