package entity

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
