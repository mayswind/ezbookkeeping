package models

// TransactionTag represents transaction tag data stored in database
type TransactionTag struct {
	TagId           int64  `xorm:"PK"`
	Uid             int64  `xorm:"UNIQUE(UQE_tag_uid_name) NOT NULL"`
	Name            string `xorm:"UNIQUE(UQE_tag_uid_name) VARCHAR(32) NOT NULL"`
	DisplayOrder    int    `xorm:"NOT NULL"`
	Hidden          bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}

// TransactionTagGetRequest represents all parameters of transaction tag getting request
type TransactionTagGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// TransactionTagCreateRequest represents all parameters of transaction tag creation request
type TransactionTagCreateRequest struct {
	Name string `json:"name" binding:"required,notBlank,max=32"`
}

// TransactionTagModifyRequest represents all parameters of transaction tag modification request
type TransactionTagModifyRequest struct {
	Id   int64  `json:"id,string" binding:"required,min=1"`
	Name string `json:"name" binding:"required,notBlank,max=32"`
}

// TransactionTagHideRequest represents all parameters of transaction tag hiding request
type TransactionTagHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

// TransactionTagMoveRequest represents all parameters of transaction tag moving request
type TransactionTagMoveRequest struct {
	NewDisplayOrders []*TransactionTagNewDisplayOrderRequest `json:"newDisplayOrders"`
}

// TransactionTagNewDisplayOrderRequest represents a data pair of id and display order
type TransactionTagNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int   `json:"displayOrder"`
}

// TransactionTagDeleteRequest represents all parameters of transaction tag deleting request
type TransactionTagDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// TransactionTagInfoResponse represents a view-object of transaction tag
type TransactionTagInfoResponse struct {
	Id           int64  `json:"id,string"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"displayOrder"`
	Hidden       bool   `json:"hidden"`
}

// ToTransactionTagInfoResponse returns a view-object according to database model
func (c *TransactionTag) ToTransactionTagInfoResponse() *TransactionTagInfoResponse {
	return &TransactionTagInfoResponse{
		Id:           c.TagId,
		Name:         c.Name,
		DisplayOrder: c.DisplayOrder,
		Hidden:       c.Hidden,
	}
}

// TransactionTagInfoResponseSlice represents the slice data structure of TransactionTagInfoResponse
type TransactionTagInfoResponseSlice []*TransactionTagInfoResponse

// Len returns the count of items
func (c TransactionTagInfoResponseSlice) Len() int {
	return len(c)
}

// Swap swaps two items
func (c TransactionTagInfoResponseSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less reports whether the first item is less than the second one
func (c TransactionTagInfoResponseSlice) Less(i, j int) bool {
	return c[i].DisplayOrder < c[j].DisplayOrder
}
