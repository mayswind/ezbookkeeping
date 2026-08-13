package models

// [PLUGIN:rules] Rule group — a named, ordered container for rules.
// Mirrors Firefly III's rule_groups concept. A group lets users organise rules and
// short-circuit evaluation (StopProcessing) once a rule in the group matches.
// This file is a self-contained extension; see docs/PLUGIN_DESIGN.md.

// RuleGroup represents a rule group stored in database
type RuleGroup struct {
	GroupId         int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_rule_group_uid_deleted) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_rule_group_uid_deleted) NOT NULL"`
	Name            string `xorm:"VARCHAR(64) NOT NULL"`
	Comment         string `xorm:"VARCHAR(255) NOT NULL"`
	DisplayOrder    int32  `xorm:"NOT NULL"`
	Active          bool   `xorm:"NOT NULL"`
	StopProcessing  bool   `xorm:"NOT NULL"` // stop all rule processing after a rule in this group matches
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// RuleGroupListRequest represents all parameters of rule group listing request
type RuleGroupListRequest struct{}

// RuleGroupGetRequest represents all parameters of a single rule group getting request
type RuleGroupGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// RuleGroupCreateRequest represents all parameters of a rule group creation request
type RuleGroupCreateRequest struct {
	Name           string `json:"name" binding:"required,notBlank,max=64"`
	Comment        string `json:"comment" binding:"max=255"`
	Active         bool   `json:"active"`
	StopProcessing bool   `json:"stopProcessing"`
	ClientSessionId string `json:"clientSessionId"`
}

// RuleGroupModifyRequest represents all parameters of a rule group modification request
type RuleGroupModifyRequest struct {
	Id             int64  `json:"id,string" binding:"required,min=1"`
	Name           string `json:"name" binding:"required,notBlank,max=64"`
	Comment        string `json:"comment" binding:"max=255"`
	Active         bool   `json:"active"`
	StopProcessing bool   `json:"stopProcessing"`
}

// RuleGroupDeleteRequest represents all parameters of a rule group deletion request
type RuleGroupDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// RuleGroupInfoResponse represents a view-object of a rule group
type RuleGroupInfoResponse struct {
	Id             int64  `json:"id,string"`
	Name           string `json:"name"`
	Comment        string `json:"comment"`
	DisplayOrder   int32  `json:"displayOrder"`
	Active         bool   `json:"active"`
	StopProcessing bool   `json:"stopProcessing"`
}

// ToRuleGroupInfoResponse returns a view-object according to the database model
func (g *RuleGroup) ToRuleGroupInfoResponse() *RuleGroupInfoResponse {
	return &RuleGroupInfoResponse{
		Id:             g.GroupId,
		Name:           g.Name,
		Comment:        g.Comment,
		DisplayOrder:   g.DisplayOrder,
		Active:         g.Active,
		StopProcessing: g.StopProcessing,
	}
}

// RuleGroupInfoResponseSlice represents the slice data structure of RuleGroupInfoResponse
type RuleGroupInfoResponseSlice []*RuleGroupInfoResponse

// Len returns the count of items
func (s RuleGroupInfoResponseSlice) Len() int { return len(s) }

// Swap swaps two items
func (s RuleGroupInfoResponseSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less reports whether the first item is less than the second one
func (s RuleGroupInfoResponseSlice) Less(i, j int) bool { return s[i].DisplayOrder < s[j].DisplayOrder }
