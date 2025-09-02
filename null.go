package date

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

var nullBytes = []byte("null")

type NullDate struct {
	Valid bool
	Date  Date
}

func NullDateFrom(date Date) NullDate {
	return NullDate{Valid: true, Date: date}
}

func (d NullDate) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return nullBytes, nil
	}
	return d.Date.MarshalJSON()
}

func (d *NullDate) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		d.Valid = false
		d.Date = Date{}
		return nil
	}

	if err := d.Date.UnmarshalJSON(data); err != nil {
		return err
	}

	d.Valid = true
	return nil
}

func (d NullDate) Value() (value driver.Value, err error) {
	if d.Valid {
		return d.Date.Value()
	}
	return nil, nil
}

func (d *NullDate) Scan(value any) error {
	if value == nil {
		d.Valid = false
		return nil
	}

	if t, ok := value.(time.Time); ok {
		if err := d.Date.Scan(t); err != nil {
			return err
		}
		d.Valid = true
		return nil
	}

	return fmt.Errorf("cannot scan type %T into NullTime", value)
}
