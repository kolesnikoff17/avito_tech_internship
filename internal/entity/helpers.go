package entity

import (
	"errors"
	"fmt"
	"time"
)

// MyTime is more fancy time format in history response
type MyTime struct {
	Time time.Time
}

// MarshalJSON -.
func (t *MyTime) MarshalJSON() ([]byte, error) {
	ts := fmt.Sprintf("\"%s\"", t.Time.Format("15:04 02 Jan 06 MST"))
	return []byte(ts), nil
}

// Scan -.
func (t *MyTime) Scan(src interface{}) error {
	switch src := src.(type) {
	case time.Time:
		t.Time = src
	default:
		return errors.New("scan error: unknown type")
	}
	return nil
}
