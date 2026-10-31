package models

import "encoding/json/v2"

const MaxImportPresetDataSize = 1<<24 - 1

// ImportPreset represents a saved import configuration stored in database
type ImportPreset struct {
	PresetId        int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_import_preset_uid_deleted_order) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_import_preset_uid_deleted_order) NOT NULL"`
	Name            string `xorm:"VARCHAR(64) NOT NULL"`
	DisplayOrder    int32  `xorm:"INDEX(IDX_import_preset_uid_deleted_order) NOT NULL"`
	Data            string `xorm:"MEDIUMBLOB"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// ImportPresetData represents a reusable import configuration
type ImportPresetData struct {
	FileType          string          `json:"fileType" binding:"required,max=64"`
	FileSubType       string          `json:"fileSubType"`
	FileEncoding      string          `json:"fileEncoding"`
	HandlingMethod    int             `json:"handlingMethod"`
	AdditionalOptions map[string]bool `json:"additionalOptions"`
	AdditionalPrompt  string          `json:"additionalPrompt"`
	DataMapping       string          `json:"dataMapping"`
	CustomScript      string          `json:"customScript"`
	AutoExecuteScript bool            `json:"autoExecuteScript"`
	ReplaceRules      string          `json:"replaceRules"`
}

// ImportPresetCreateRequest represents all parameters of import preset creation request
type ImportPresetCreateRequest struct {
	Name            string `json:"name" binding:"required,notBlank,max=64"`
	ClientSessionId string `json:"clientSessionId"`
	ImportPresetData
}

// ImportPresetModifyRequest represents all parameters of import preset modification request
type ImportPresetModifyRequest struct {
	Id   int64  `json:"id,string" binding:"required,min=1"`
	Name string `json:"name" binding:"required,notBlank,max=64"`
	ImportPresetData
}

// ImportPresetGetRequest represents all parameters of import preset getting request
type ImportPresetGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// ImportPresetMoveRequest represents all parameters of import preset moving request
type ImportPresetMoveRequest struct {
	NewDisplayOrders []*ImportPresetNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// ImportPresetNewDisplayOrderRequest represents a data pair of id and display order
type ImportPresetNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// ImportPresetDeleteRequest represents all parameters of import preset deleting request
type ImportPresetDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// ImportPresetInfoResponse represents a view-object of import preset info
type ImportPresetInfoResponse struct {
	Id           int64  `json:"id,string"`
	Name         string `json:"name"`
	DisplayOrder int32  `json:"displayOrder"`
	ImportPresetData
}

// ToImportPresetInfoResponse returns the complete preset configuration
func (p *ImportPreset) ToImportPresetInfoResponse() (*ImportPresetInfoResponse, error) {
	response := &ImportPresetInfoResponse{
		Id:           p.PresetId,
		Name:         p.Name,
		DisplayOrder: p.DisplayOrder,
	}

	if len(p.Data) > 0 {
		err := json.Unmarshal([]byte(p.Data), &response.ImportPresetData)

		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// ImportPresetInfoResponseSlice represents the slice data structure of ImportPresetInfoResponse
type ImportPresetInfoResponseSlice []*ImportPresetInfoResponse

// Len returns the count of items
func (s ImportPresetInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s ImportPresetInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s ImportPresetInfoResponseSlice) Less(i, j int) bool {
	return s[i].DisplayOrder < s[j].DisplayOrder
}
