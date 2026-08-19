package models

const UserCustomIconFileExtension = "png"

// UserCustomIcon represents a user uploaded icon stored in object storage.
type UserCustomIcon struct {
	IconId          int64 `xorm:"PK"`
	Uid             int64 `xorm:"INDEX(IDX_user_custom_icon_uid_deleted_order) NOT NULL"`
	Deleted         bool  `xorm:"INDEX(IDX_user_custom_icon_uid_deleted_order) NOT NULL"`
	DisplayOrder    int32 `xorm:"INDEX(IDX_user_custom_icon_uid_deleted_order) NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// UserCustomIconMoveRequest represents all parameters of user custom icon moving request
type UserCustomIconMoveRequest struct {
	NewDisplayOrders []*UserCustomIconNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// UserCustomIconNewDisplayOrderRequest represents a data pair of id and display order
type UserCustomIconNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// UserCustomIconDeleteRequest represents all parameters of user custom icon deletion request
type UserCustomIconDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// UserCustomIconInfoResponse represents the response data of a user custom icon
type UserCustomIconInfoResponse struct {
	Id           int64 `json:"id,string"`
	DisplayOrder int32 `json:"displayOrder"`
}

// ToUserCustomIconInfoResponse returns a view-object according to database model
func (i *UserCustomIcon) ToUserCustomIconInfoResponse() *UserCustomIconInfoResponse {
	return &UserCustomIconInfoResponse{
		Id:           i.IconId,
		DisplayOrder: i.DisplayOrder,
	}
}

// UserCustomIconInfoResponseSlice represents the slice data structure of UserCustomIconInfoResponse
type UserCustomIconInfoResponseSlice []*UserCustomIconInfoResponse

// Len returns the count of items
func (s UserCustomIconInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s UserCustomIconInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s UserCustomIconInfoResponseSlice) Less(i, j int) bool {
	return s[i].DisplayOrder < s[j].DisplayOrder
}
