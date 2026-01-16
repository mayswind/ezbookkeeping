package models

// TransactionTagGroup represents transaction tag group data stored in database
type TransactionTagGroup struct {
	TagGroupId      int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_tag_group_uid_deleted_order) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_tag_group_uid_deleted_order) NOT NULL"`
	Name            string `xorm:"VARCHAR(64) NOT NULL"`
	DisplayOrder    int32  `xorm:"INDEX(IDX_tag_group_uid_deleted_order) NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// TransactionTagGroupGetRequest represents all parameters of transaction tag group getting request
type TransactionTagGroupGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// TransactionTagGroupCreateRequest represents all parameters of transaction tag group creation request
type TransactionTagGroupCreateRequest struct {
	Name string `json:"name" binding:"required,notBlank,max=64"`
}

// TransactionTagGroupModifyRequest represents all parameters of transaction tag group modification request
type TransactionTagGroupModifyRequest struct {
	Id   int64  `json:"id,string" binding:"required,min=1"`
	Name string `json:"name" binding:"required,notBlank,max=64"`
}

// TransactionTagGroupMoveRequest represents all parameters of transaction tag group moving request
type TransactionTagGroupMoveRequest struct {
	NewDisplayOrders []*TransactionTagGroupNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// TransactionTagGroupNewDisplayOrderRequest represents a data pair of id and display order
type TransactionTagGroupNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// TransactionTagGroupDeleteRequest represents all parameters of transaction tag group deleting request
type TransactionTagGroupDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// TransactionTagGroupInfoResponse represents a view-object of transaction tag group
type TransactionTagGroupInfoResponse struct {
	Id           int64  `json:"id,string"`
	Name         string `json:"name"`
	DisplayOrder int32  `json:"displayOrder"`
}

// ToTransactionTagGroupInfoResponse returns a view-object according to database model
func (t *TransactionTagGroup) ToTransactionTagGroupInfoResponse() *TransactionTagGroupInfoResponse {
	return &TransactionTagGroupInfoResponse{
		Id:           t.TagGroupId,
		Name:         t.Name,
		DisplayOrder: t.DisplayOrder,
	}
}

// TransactionTagGroupInfoResponseSlice represents the slice data structure of TransactionTagGroupInfoResponse
type TransactionTagGroupInfoResponseSlice []*TransactionTagGroupInfoResponse

// Len returns the count of items
func (s TransactionTagGroupInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionTagGroupInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionTagGroupInfoResponseSlice) Less(i, j int) bool {
	return s[i].DisplayOrder < s[j].DisplayOrder
}
