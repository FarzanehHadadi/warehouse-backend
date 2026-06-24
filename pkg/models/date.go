package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

type Date struct {
	time.Time
}

func (d Date) IsZero() bool {
	return d.Time.IsZero()
}

func ParseDate(value string) (Date, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return Date{}, fmt.Errorf("date is required")
	}

	t, err := time.Parse(DateLayout, value)
	if err != nil {
		return Date{}, fmt.Errorf("invalid date %q, expected format YYYY-MM-DD", value)
	}

	return Date{Time: t}, nil
}

func (d Date) String() string {
	if d.Time.IsZero() {
		return ""
	}
	return d.Time.Format(DateLayout)
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "" {
		d.Time = time.Time{}
		return nil
	}

	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}

	d.Time = parsed.Time
	return nil
}

func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.String(), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)
	case []byte:
		return d.scanString(string(v))
	case string:
		return d.scanString(v)
	default:
		return fmt.Errorf("cannot scan %T into Date", value)
	}

	return nil
}

func (d *Date) scanString(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		d.Time = time.Time{}
		return nil
	}

	if len(value) >= len(DateLayout) {
		value = value[:len(DateLayout)]
	}

	parsed, err := ParseDate(value)
	if err != nil {
		return err
	}

	d.Time = parsed.Time
	return nil
}
