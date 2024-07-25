package models

// TransactionTagIndex represents transaction and transaction tag relation stored in database
type TransactionTagIndex struct {
	TagIndexId      int64 `xorm:"PK"`
	Uid             int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_id)"`
	Deleted         bool  `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_id) NOT NULL"`
	TransactionTime int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id) NOT NULL"`
	TagId           int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id)"`
	TransactionId   int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_id)"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}
