package models

// TransactionTag represents transaction tag data stored in database
type TransactionTag struct {
	TagId           int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_tag_uid_deleted_order) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_tag_uid_deleted_order) NOT NULL"`
	Name            string `xorm:"VARCHAR(64) NOT NULL"`
	DisplayOrder    int32  `xorm:"INDEX(IDX_tag_uid_deleted_order) NOT NULL"`
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
	Name string `json:"name" binding:"required,notBlank,max=64"`
}

// TransactionTagModifyRequest represents all parameters of transaction tag modification request
type TransactionTagModifyRequest struct {
	Id   int64  `json:"id,string" binding:"required,min=1"`
	Name string `json:"name" binding:"required,notBlank,max=64"`
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
	DisplayOrder int32  `json:"displayOrder"`
	Hidden       bool   `json:"hidden"`
}

// ToTransactionTagInfoResponse returns a view-object according to database model
func (t *TransactionTag) ToTransactionTagInfoResponse() *TransactionTagInfoResponse {
	return &TransactionTagInfoResponse{
		Id:           t.TagId,
		Name:         t.Name,
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
	return s[i].DisplayOrder < s[j].DisplayOrder
}
