package models

// Project represents project data stored in database
type Project struct {
	ProjectId       int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_project_uid_deleted_order) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_project_uid_deleted_order) NOT NULL"`
	Name            string `xorm:"VARCHAR(64) NOT NULL"`
	Color           string `xorm:"VARCHAR(6) NOT NULL"`
	Comment         string `xorm:"VARCHAR(255) NOT NULL"`
	DisplayOrder    int32  `xorm:"INDEX(IDX_project_uid_deleted_order) NOT NULL"`
	Hidden          bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// ProjectGetRequest represents all parameters of project getting request
type ProjectGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// ProjectCreateRequest represents all parameters of project creation request
type ProjectCreateRequest struct {
	Name    string `json:"name" binding:"required,notBlank,max=64"`
	Color   string `json:"color" binding:"required,validHexRGBColor"`
	Comment string `json:"comment" binding:"max=255"`
}

// ProjectModifyRequest represents all parameters of project modification request
type ProjectModifyRequest struct {
	Id      int64  `json:"id,string" binding:"required,min=1"`
	Name    string `json:"name" binding:"required,notBlank,max=64"`
	Color   string `json:"color" binding:"required,validHexRGBColor"`
	Comment string `json:"comment" binding:"max=255"`
	Hidden  bool   `json:"hidden"`
}

// ProjectHideRequest represents all parameters of project hiding request
type ProjectHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

// ProjectMoveRequest represents all parameters of project moving request
type ProjectMoveRequest struct {
	NewDisplayOrders []*ProjectNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// ProjectNewDisplayOrderRequest represents a data pair of id and display order
type ProjectNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// ProjectDeleteRequest represents all parameters of project deleting request
type ProjectDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// ProjectInfoResponse represents a view-object of project
type ProjectInfoResponse struct {
	Id           int64  `json:"id,string"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	Comment      string `json:"comment"`
	DisplayOrder int32  `json:"displayOrder"`
	Hidden       bool   `json:"hidden"`
}

// FillFromOtherProject fills all the fields in this current project from other project
func (p *Project) FillFromOtherProject(project *Project) {
	p.ProjectId = project.ProjectId
	p.Uid = project.Uid
	p.Deleted = project.Deleted
	p.Name = project.Name
	p.Color = project.Color
	p.Comment = project.Comment
	p.DisplayOrder = project.DisplayOrder
	p.Hidden = project.Hidden
	p.CreatedUnixTime = project.CreatedUnixTime
	p.UpdatedUnixTime = project.UpdatedUnixTime
	p.DeletedUnixTime = project.DeletedUnixTime
}

// ToProjectInfoResponse returns a view-object according to database model
func (p *Project) ToProjectInfoResponse() *ProjectInfoResponse {
	return &ProjectInfoResponse{
		Id:           p.ProjectId,
		Name:         p.Name,
		Color:        p.Color,
		Comment:      p.Comment,
		DisplayOrder: p.DisplayOrder,
		Hidden:       p.Hidden,
	}
}

// ProjectInfoResponseSlice represents the slice data structure of ProjectInfoResponse
type ProjectInfoResponseSlice []*ProjectInfoResponse

// Len returns the count of items
func (s ProjectInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s ProjectInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s ProjectInfoResponseSlice) Less(i, j int) bool {
	return s[i].DisplayOrder < s[j].DisplayOrder
}
