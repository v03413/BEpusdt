package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Datetime time.Time

func (d *Datetime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*d)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (d Datetime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(d)
	// 判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (d *Datetime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*d = Datetime(value)

		return nil
	}

	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (d *Datetime) String() string {

	return time.Time(*d).Format(time.DateTime)
}

func (d *Datetime) Year() string {

	return time.Time(*d).Format("2006")
}

func (d *Datetime) Before(u time.Time) bool {

	return time.Time(*d).Before(u)
}

func (d *Datetime) Time() time.Time {

	return time.Time(*d)
}

func (d *Datetime) Format(layout string) string {

	return time.Time(*d).Format(layout)
}
