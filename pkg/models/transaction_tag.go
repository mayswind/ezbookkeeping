package models

type TransactionTag struct {
	TagId           int64  `xorm:"PK"`
	Uid             int64  `xorm:"UNIQUE(UQE_tag_uid_name) NOT NULL"`
	Name            string `xorm:"UNIQUE(UQE_tag_uid_name) VARCHAR(32) NOT NULL"`
	DisplayOrder    int    `xorm:"NOT NULL"`
	Hidden          bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}

type TransactionTagGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

type TransactionTagCreateRequest struct {
	Name string `json:"name" binding:"required,notBlank,max=32"`
}

type TransactionTagModifyRequest struct {
	Id   int64  `json:"id,string" binding:"required,min=1"`
	Name string `json:"name" binding:"required,notBlank,max=32"`
}

type TransactionTagHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

type TransactionTagMoveRequest struct {
	NewDisplayOrders []*TransactionTagNewDisplayOrderRequest `json:"newDisplayOrders"`
}

type TransactionTagNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int   `json:"displayOrder"`
}

type TransactionTagDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

type TransactionTagInfoResponse struct {
	Id           int64  `json:"id,string"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"displayOrder"`
	Hidden       bool   `json:"hidden"`
}

func (c *TransactionTag) ToTransactionTagInfoResponse() *TransactionTagInfoResponse {
	return &TransactionTagInfoResponse{
		Id:           c.TagId,
		Name:         c.Name,
		DisplayOrder: c.DisplayOrder,
		Hidden:       c.Hidden,
	}
}

type TransactionTagInfoResponseSlice []*TransactionTagInfoResponse

func (c TransactionTagInfoResponseSlice) Len() int {
	return len(c)
}

func (c TransactionTagInfoResponseSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c TransactionTagInfoResponseSlice) Less(i, j int) bool {
	return c[i].DisplayOrder < c[j].DisplayOrder
}
