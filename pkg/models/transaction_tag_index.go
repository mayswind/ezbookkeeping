package models

// TransactionTagIndex represents transaction and transaction tag relation stored in database
type TransactionTagIndex struct {
	Uid             int64 `xorm:"PK INDEX(IDX_transaction_tag_index_uid_tag_id_transaction_time) INDEX(IDX_transaction_tag_index_uid_transaction_id)"`
	TagId           int64 `xorm:"PK INDEX(IDX_transaction_tag_index_uid_tag_id_transaction_time)"`
	TransactionId   int64 `xorm:"PK INDEX(IDX_transaction_tag_index_uid_transaction_id)"`
	TransactionTime int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_tag_id_transaction_time) NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}
