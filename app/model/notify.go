package model

type NotifyRecord struct {
	Id
	Txid string `gorm:"type:varchar(64);uniqueIndex;not null;comment:交易哈希"`
	AutoTimeAt
}

func (nr NotifyRecord) TableName() string {

	return "bep_notify"
}

func IsNeedNotifyByTxid(txid string) bool {
	var row NotifyRecord
	var res = Db.Where("txid = ?", txid).Limit(1).Find(&row)
	if res.RowsAffected > 0 {

		return false
	}

	var row2 Order
	var res2 = Db.Where("trade_hash = ?", txid).Limit(1).Find(&row2)
	if res2.RowsAffected > 0 {

		return false
	}

	return true
}
