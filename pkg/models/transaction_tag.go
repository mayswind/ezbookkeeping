package models

// TransactionTag represents transaction tag data stored in database
type TransactionTag struct {
	TagId           int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_tag_uid_deleted_group_order) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_tag_uid_deleted_group_order) NOT NULL"`
	TagGroupId      int64  `xorm:"INDEX(IDX_tag_uid_deleted_group_order) NOT NULL DEFAULT 0"`
	Name            string `xorm:"VARCHAR(64) NOT NULL"`
	DisplayOrder    int32  `xorm:"INDEX(IDX_tag_uid_deleted_group_order) NOT NULL"`
	Hidden          bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// TransactionTagGetRequest represents all parameters of transaction tag getting request
type TransactionTagGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// TransactionTagCreateRequest represents all parameters of transaction tag creation request
type TransactionTagCreateRequest struct {
	GroupId int64  `json:"groupId,string"`
	Name    string `json:"name" binding:"required,notBlank,max=64"`
}

// TransactionTagCreateBatchRequest represents all parameters of transaction tag batch creation request
type TransactionTagCreateBatchRequest struct {
	Tags       []*TransactionTagCreateRequest `json:"tags" binding:"required"`
	GroupId    int64                          `json:"groupId,string"`
	SkipExists bool                           `json:"skipExists"`
}

// TransactionTagModifyRequest represents all parameters of transaction tag modification request
type TransactionTagModifyRequest struct {
	Id      int64  `json:"id,string" binding:"required,min=1"`
	GroupId int64  `json:"groupId,string"`
	Name    string `json:"name" binding:"required,notBlank,max=64"`
}

// TransactionTagHideRequest represents all parameters of transaction tag hiding request
type TransactionTagHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

// TransactionTagMoveRequest represents all parameters of transaction tag moving request
type TransactionTagMoveRequest struct {
	NewDisplayOrders []*TransactionTagNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// TransactionTagNewDisplayOrderRequest represents a data pair of id and display order
type TransactionTagNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// TransactionTagDeleteRequest represents all parameters of transaction tag deleting request
type TransactionTagDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// TransactionTagInfoResponse represents a view-object of transaction tag
type TransactionTagInfoResponse struct {
	Id           int64  `json:"id,string"`
	Name         string `json:"name"`
	TagGroupId   int64  `json:"groupId,string"`
	DisplayOrder int32  `json:"displayOrder"`
	Hidden       bool   `json:"hidden"`
}

// FillFromOtherTag fills all the fields in this current tag from other transaction tag
func (t *TransactionTag) FillFromOtherTag(tag *TransactionTag) {
	t.TagId = tag.TagId
	t.Uid = tag.Uid
	t.Deleted = tag.Deleted
	t.Name = tag.Name
	t.TagGroupId = tag.TagGroupId
	t.DisplayOrder = tag.DisplayOrder
	t.Hidden = tag.Hidden
	t.CreatedUnixTime = tag.CreatedUnixTime
	t.UpdatedUnixTime = tag.UpdatedUnixTime
	t.DeletedUnixTime = tag.DeletedUnixTime
}

// ToTransactionTagInfoResponse returns a view-object according to database model
func (t *TransactionTag) ToTransactionTagInfoResponse() *TransactionTagInfoResponse {
	return &TransactionTagInfoResponse{
		Id:           t.TagId,
		Name:         t.Name,
		TagGroupId:   t.TagGroupId,
		DisplayOrder: t.DisplayOrder,
		Hidden:       t.Hidden,
	}
}

// TransactionTagInfoResponseSlice represents the slice data structure of TransactionTagInfoResponse
type TransactionTagInfoResponseSlice []*TransactionTagInfoResponse

// Len returns the count of items
func (s TransactionTagInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionTagInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionTagInfoResponseSlice) Less(i, j int) bool {
	if s[i].TagGroupId != s[j].TagGroupId {
		return s[i].TagGroupId < s[j].TagGroupId
	}

	return s[i].DisplayOrder < s[j].DisplayOrder
}
