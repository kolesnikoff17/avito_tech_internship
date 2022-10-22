package entity

import "time"

// Balance -.
type Balance struct {
	ID     int    `json:"id" db:"user_id"`
	Amount string `json:"amount" db:"amount"`
}

// Order -.
type Order struct {
	ID          int       `json:"-" db:"order_id"`
	ServiceID   int       `json:"-" db:"service_id"`
	ServiceName string    `json:"service" db:"service_name"`
	UserID      int       `json:"-" db:"user_id"`
	StatusID    int       `json:"-" db:"status_id"`
	Status      string    `json:"status" db:"status_name"`
	Sum         string    `json:"sum" db:"order_sum"`
	Time        time.Time `json:"time" db:"created"`
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
