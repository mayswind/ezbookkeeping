package models

// TransactionPictureInfo represents transaction picture file info stored in database
type TransactionPictureInfo struct {
	Uid              int64  `xorm:"INDEX(IDX_transaction_picture_uid_deleted_transaction_id_picture_id) INDEX(IDX_transaction_picture_uid_deleted_picture_id) NOT NULL"`
	Deleted          bool   `xorm:"INDEX(IDX_transaction_picture_uid_deleted_transaction_id_picture_id) INDEX(IDX_transaction_picture_uid_deleted_picture_id) NOT NULL"`
	TransactionId    int64  `xorm:"INDEX(IDX_transaction_picture_uid_deleted_transaction_id_picture_id) NOT NULL"`
	PictureId        int64  `xorm:"PK INDEX(IDX_transaction_picture_uid_deleted_transaction_id_picture_id) INDEX(IDX_transaction_picture_uid_deleted_picture_id)"`
	PictureExtension string `xorm:"VARCHAR(10) NOT NULL"`
	CreatedIp        string `xorm:"VARCHAR(39)"`
	CreatedUnixTime  int64
	UpdatedUnixTime  int64
	DeletedUnixTime  int64
}

// TransactionPictureInfoBasicResponse represents a view-object of transaction picture basic info
type TransactionPictureInfoBasicResponse struct {
	PictureId   int64  `json:"pictureId,string"`
	OriginalUrl string `json:"originalUrl"`
}

// ToTransactionPictureInfoBasicResponse returns a view-object according to database model
func (p *TransactionPictureInfo) ToTransactionPictureInfoBasicResponse(originalUrl string) *TransactionPictureInfoBasicResponse {
	return &TransactionPictureInfoBasicResponse{
		PictureId:   p.PictureId,
		OriginalUrl: originalUrl,
	}
}
