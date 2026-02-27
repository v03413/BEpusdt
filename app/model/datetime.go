package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var localLocation = time.Local

type Datetime time.Time

func (d *Datetime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*d).In(localLocation)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (d Datetime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(d)
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
	return time.Time(*d).In(localLocation).Format(time.DateTime)
}

func (d *Datetime) Year() string {
	return time.Time(*d).In(localLocation).Format("2006")
}

func (d *Datetime) Before(u time.Time) bool {
	return time.Time(*d).Before(u)
}

func (d *Datetime) Time() time.Time {
	return time.Time(*d).In(localLocation)
}

func (d *Datetime) Format(layout string) string {
	return time.Time(*d).In(localLocation).Format(layout)
}
