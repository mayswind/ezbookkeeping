package models

// [PLUGIN:rules] Rule — a named set of triggers (conditions) and actions (mutations).
// When a TransactionCreateRequest is processed, the engine evaluates rule groups in order,
// then rules within each group. If a rule's triggers match (ALL if Strict, else ANY),
// its actions are applied to the request DTO BEFORE it is turned into a Transaction.
// See docs/PLUGIN_DESIGN.md.

// Rule represents a rule stored in database
type Rule struct {
	RuleId          int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_rule_uid_deleted_group_order) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_rule_uid_deleted_group_order) NOT NULL"`
	RuleGroupId     int64  `xorm:"INDEX(IDX_rule_uid_deleted_group_order) NOT NULL"`
	Name            string `xorm:"VARCHAR(64) NOT NULL"`
	Comment         string `xorm:"VARCHAR(255) NOT NULL"`
	DisplayOrder    int32  `xorm:"INDEX(IDX_rule_uid_deleted_group_order) NOT NULL"`
	Active          bool   `xorm:"NOT NULL"`
	Strict          bool   `xorm:"NOT NULL"` // true = ALL triggers must match, false = ANY
	StopProcessing  bool   `xorm:"NOT NULL"` // stop the whole group once this rule matches
	ApplyOnCreate   bool   `xorm:"NOT NULL"` // evaluate on transaction creation/import
	ApplyOnUpdate   bool   `xorm:"NOT NULL"` // evaluate on transaction modification (v1: not wired)
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// RuleTriggerInput represents a trigger in a rule create/modify request
type RuleTriggerInput struct {
	TriggerType    RuleTriggerType `json:"triggerType" binding:"required"`
	TriggerValue   string          `json:"triggerValue" binding:"max=255"`
	Prohibited     bool            `json:"prohibited"`
	StopProcessing bool            `json:"stopProcessing"`
}

// RuleActionInput represents an action in a rule create/modify request
type RuleActionInput struct {
	ActionType     RuleActionType `json:"actionType" binding:"required"`
	ActionValue    string         `json:"actionValue" binding:"max=255"`
	StopProcessing bool           `json:"stopProcessing"`
}

// RuleListByGroupRequest represents all parameters of rule listing request
type RuleListByGroupRequest struct {
	GroupId int64 `form:"groupId,string,default=0" binding:"min=0"`
}

// RuleGetRequest represents all parameters of a single rule getting request
type RuleGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// RuleCreateRequest represents all parameters of a rule creation request (rule + triggers + actions)
type RuleCreateRequest struct {
	RuleGroupId     int64              `json:"ruleGroupId,string" binding:"required,min=1"`
	Name            string             `json:"name" binding:"required,notBlank,max=64"`
	Comment         string             `json:"comment" binding:"max=255"`
	Active          bool               `json:"active"`
	Strict          bool               `json:"strict"`
	StopProcessing  bool               `json:"stopProcessing"`
	ApplyOnCreate   bool               `json:"applyOnCreate"`
	ApplyOnUpdate   bool               `json:"applyOnUpdate"`
	Triggers        []*RuleTriggerInput `json:"triggers" binding:"required,min=1,dive"`
	Actions         []*RuleActionInput  `json:"actions" binding:"required,min=1,dive"`
	ClientSessionId string             `json:"clientSessionId"`
}

// RuleModifyRequest represents all parameters of a rule modification request
type RuleModifyRequest struct {
	Id             int64               `json:"id,string" binding:"required,min=1"`
	RuleGroupId    int64               `json:"ruleGroupId,string" binding:"required,min=1"`
	Name           string              `json:"name" binding:"required,notBlank,max=64"`
	Comment        string              `json:"comment" binding:"max=255"`
	Active         bool                `json:"active"`
	Strict         bool                `json:"strict"`
	StopProcessing bool                `json:"stopProcessing"`
	ApplyOnCreate  bool                `json:"applyOnCreate"`
	ApplyOnUpdate  bool                `json:"applyOnUpdate"`
	Triggers       []*RuleTriggerInput `json:"triggers" binding:"required,min=1,dive"`
	Actions        []*RuleActionInput  `json:"actions" binding:"required,min=1,dive"`
}

// RuleDeleteRequest represents all parameters of a rule deletion request
type RuleDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// RuleInfoResponse represents a view-object of a rule (with its triggers and actions)
type RuleInfoResponse struct {
	Id             int64                      `json:"id,string"`
	RuleGroupId    int64                      `json:"ruleGroupId,string"`
	Name           string                     `json:"name"`
	Comment        string                     `json:"comment"`
	DisplayOrder   int32                      `json:"displayOrder"`
	Active         bool                       `json:"active"`
	Strict         bool                       `json:"strict"`
	StopProcessing bool                       `json:"stopProcessing"`
	ApplyOnCreate  bool                       `json:"applyOnCreate"`
	ApplyOnUpdate  bool                       `json:"applyOnUpdate"`
	Triggers       []*RuleTriggerInfoResponse `json:"triggers"`
	Actions        []*RuleActionInfoResponse  `json:"actions"`
}

// ToRuleInfoResponse returns a view-object according to the database model and its
// triggers / actions. The triggers and actions are ordered by DisplayOrder by the caller.
func (r *Rule) ToRuleInfoResponse(triggers []*RuleTrigger, actions []*RuleAction) *RuleInfoResponse {
	resp := &RuleInfoResponse{
		Id:             r.RuleId,
		RuleGroupId:    r.RuleGroupId,
		Name:           r.Name,
		Comment:        r.Comment,
		DisplayOrder:   r.DisplayOrder,
		Active:         r.Active,
		Strict:         r.Strict,
		StopProcessing: r.StopProcessing,
		ApplyOnCreate:  r.ApplyOnCreate,
		ApplyOnUpdate:  r.ApplyOnUpdate,
	}

	resp.Triggers = make([]*RuleTriggerInfoResponse, 0, len(triggers))
	for i := 0; i < len(triggers); i++ {
		resp.Triggers = append(resp.Triggers, triggers[i].ToRuleTriggerInfoResponse())
	}

	resp.Actions = make([]*RuleActionInfoResponse, 0, len(actions))
	for i := 0; i < len(actions); i++ {
		resp.Actions = append(resp.Actions, actions[i].ToRuleActionInfoResponse())
	}

	return resp
}
