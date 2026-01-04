package model

type NotifyRecord struct {
	Id
	Txid string `Gorm:"type:varchar(64);uniqueIndex;not null;comment:交易哈希"`
	AutoTimeAt
}

func (nr NotifyRecord) TableName() string {

	return "bep_notify"
}

func IsNeedNotifyByTxid(txid string) bool {
	var record NotifyRecord
	var res = Db.Where("txid = ?", txid).Limit(1).Find(&record)
	if res.RowsAffected > 0 {

		return false
	}

	var order Order
	var res2 = Db.Where("ref_hash = ?", txid).Limit(1).Find(&order)
	if res2.RowsAffected > 0 {

		return false
	}

	return true
}
