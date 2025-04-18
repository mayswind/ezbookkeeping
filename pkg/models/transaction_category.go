package models

// LevelOneTransactionCategoryParentId represents the parent id of level-one transaction category
const LevelOneTransactionCategoryParentId = 0

// TransactionCategoryType represents transaction category type
type TransactionCategoryType byte

// Transaction category types
const (
	CATEGORY_TYPE_INCOME   TransactionCategoryType = 1
	CATEGORY_TYPE_EXPENSE  TransactionCategoryType = 2
	CATEGORY_TYPE_TRANSFER TransactionCategoryType = 3
)

// TransactionCategory represents transaction category data stored in database
type TransactionCategory struct {
	CategoryId       int64                   `xorm:"PK"`
	Uid              int64                   `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
	Deleted          bool                    `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
	Type             TransactionCategoryType `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
	ParentCategoryId int64                   `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
	Name             string                  `xorm:"VARCHAR(64) NOT NULL"`
	DisplayOrder     int32                   `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
	Icon             int64                   `xorm:"NOT NULL"`
	Color            string                  `xorm:"VARCHAR(6) NOT NULL"`
	Hidden           bool                    `xorm:"NOT NULL"`
	Comment          string                  `xorm:"VARCHAR(255) NOT NULL"`
	CreatedUnixTime  int64
	UpdatedUnixTime  int64
	DeletedUnixTime  int64
}

// TransactionCategoryListRequest represents all parameters of transaction category listing request
type TransactionCategoryListRequest struct {
	Type     TransactionCategoryType `form:"type" binding:"min=0"`
	ParentId int64                   `form:"parent_id,string,default=-1" binding:"min=-1"`
}

// TransactionCategoryGetRequest represents all parameters of transaction category getting request
type TransactionCategoryGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// TransactionCategoryCreateRequest represents all parameters of single transaction category creation request
type TransactionCategoryCreateRequest struct {
	Name            string                  `json:"name" binding:"required,notBlank,max=64"`
	Type            TransactionCategoryType `json:"type" binding:"required"`
	ParentId        int64                   `json:"parentId,string" binding:"min=0"`
	Icon            int64                   `json:"icon,string" binding:"min=1"`
	Color           string                  `json:"color" binding:"required,len=6,validHexRGBColor"`
	Comment         string                  `json:"comment" binding:"max=255"`
	ClientSessionId string                  `json:"clientSessionId"`
}

// TransactionCategoryCreateBatchRequest represents all parameters of transaction category batch creation request
type TransactionCategoryCreateBatchRequest struct {
	Categories []*TransactionCategoryCreateWithSubCategories `json:"categories" binding:"required"`
}

// TransactionCategoryCreateWithSubCategories represents all parameters of multi transaction categories creation request
type TransactionCategoryCreateWithSubCategories struct {
	Name          string                              `json:"name" binding:"required,notBlank,max=64"`
	Type          TransactionCategoryType             `json:"type" binding:"required"`
	Icon          int64                               `json:"icon,string" binding:"min=1"`
	Color         string                              `json:"color" binding:"required,len=6,validHexRGBColor"`
	Comment       string                              `json:"comment" binding:"max=255"`
	SubCategories []*TransactionCategoryCreateRequest `json:"subCategories" binding:"required"`
}

// TransactionCategoryModifyRequest represents all parameters of transaction category modification request
type TransactionCategoryModifyRequest struct {
	Id       int64  `json:"id,string" binding:"required,min=1"`
	Name     string `json:"name" binding:"required,notBlank,max=64"`
	ParentId int64  `json:"parentId,string" binding:"min=0"`
	Icon     int64  `json:"icon,string" binding:"min=1"`
	Color    string `json:"color" binding:"required,len=6,validHexRGBColor"`
	Comment  string `json:"comment" binding:"max=255"`
	Hidden   bool   `json:"hidden"`
}

// TransactionCategoryHideRequest represents all parameters of transaction category hiding request
type TransactionCategoryHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

// TransactionCategoryMoveRequest represents all parameters of transaction category moving request
type TransactionCategoryMoveRequest struct {
	NewDisplayOrders []*TransactionCategoryNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// TransactionCategoryNewDisplayOrderRequest represents a data pair of id and display order
type TransactionCategoryNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// TransactionCategoryDeleteRequest represents all parameters of transaction category deleting request
type TransactionCategoryDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// TransactionCategoryInfoResponse represents a view-object of transaction category
type TransactionCategoryInfoResponse struct {
	Id            int64                                `json:"id,string"`
	Name          string                               `json:"name"`
	ParentId      int64                                `json:"parentId,string"`
	Type          TransactionCategoryType              `json:"type"`
	Icon          int64                                `json:"icon,string"`
	Color         string                               `json:"color"`
	Comment       string                               `json:"comment"`
	DisplayOrder  int32                                `json:"displayOrder"`
	Hidden        bool                                 `json:"hidden"`
	SubCategories TransactionCategoryInfoResponseSlice `json:"subCategories,omitempty"`
}

// ToTransactionCategoryInfoResponse returns a view-object according to database model
func (c *TransactionCategory) ToTransactionCategoryInfoResponse() *TransactionCategoryInfoResponse {
	return &TransactionCategoryInfoResponse{
		Id:           c.CategoryId,
		Name:         c.Name,
		ParentId:     c.ParentCategoryId,
		Type:         c.Type,
		Icon:         c.Icon,
		Color:        c.Color,
		Comment:      c.Comment,
		DisplayOrder: c.DisplayOrder,
		Hidden:       c.Hidden,
	}
}

// TransactionCategoryInfoResponseSlice represents the slice data structure of TransactionCategoryInfoResponse
type TransactionCategoryInfoResponseSlice []*TransactionCategoryInfoResponse

// Len returns the count of items
func (s TransactionCategoryInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionCategoryInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionCategoryInfoResponseSlice) Less(i, j int) bool {
	return s[i].DisplayOrder < s[j].DisplayOrder
}
