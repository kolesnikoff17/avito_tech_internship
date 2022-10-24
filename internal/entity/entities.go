package entity

import (
	"errors"
	"fmt"
	"time"
)

// Balance -.
type Balance struct {
	ID     int    `json:"id" db:"user_id"`
	Amount string `json:"amount" db:"amount"`
}

// Order -.
type Order struct {
	ID          int    `json:"-" db:"order_id"`
	UserID      int    `json:"-" db:"user_id"`
	Sum         string `json:"sum" db:"order_sum"`
	ServiceID   int    `json:"-" db:"service_id"`
	ServiceName string `json:"service" db:"service_name"`
	StatusID    int    `json:"-" db:"status_id"`
	Status      string `json:"status" db:"status_name"`
	Time        MyTime `json:"time" db:"created"`
}

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
	switch src.(type) {
	case time.Time:
		t.Time = src.(time.Time)
	default:
		return errors.New("scan error: unknown type")
	}
	return nil
}

// History -.
type History struct {
	Orders  []Order `json:"orders"`
	UserID  int     `json:"-"`
	Limit   int     `json:"-"`
	OrderBy string  `json:"-"`
	Desc    bool    `json:"-"`
	Page    int     `json:"-"`
}

// SumByService -.
type SumByService struct {
	Sum  string `db:"sums"`
	Name string `db:"service_name"`
}

// Report -.
type Report struct {
	Sums []SumByService
}
