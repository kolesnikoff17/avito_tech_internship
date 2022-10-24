package usecase

import (
	"balance_api/internal/entity"
	reportmock "balance_api/internal/mocks/report"
	repomock "balance_api/internal/mocks/repository"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	r := repomock.NewBalanceRepo(t)
	uc := New(r, reportmock.NewReportFile(t))

	r.On("GetByID", ctx, 1).Return(entity.Balance{ID: 1, Amount: "200"}, nil)
	r.On("GetByID", ctx, 2).Return(entity.Balance{}, entity.ErrNoID)

	type TestCase struct {
		name        string
		id          int
		expectedVal entity.Balance
		expectedErr error
	}

	cases := []TestCase{{
		name: "valid",
		id:   1,
		expectedVal: entity.Balance{
			ID:     1,
			Amount: "200",
		},
		expectedErr: nil,
	}, {
		name:        "no such id",
		id:          2,
		expectedVal: entity.Balance{},
		expectedErr: entity.ErrNoID,
	},
	}

	for _, tc := range cases {
		val, err := uc.GetByID(ctx, tc.id)
		assert.Equal(t, tc.expectedVal, val)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	r := repomock.NewBalanceRepo(t)
	uc := New(r, reportmock.NewReportFile(t))

	r.On("GetByID", ctx, 1).Return(entity.Balance{ID: 1, Amount: "300"}, nil)
	r.On("CheckServiceID", ctx, 1).Return(nil)
	r.On("GetOrderByID", ctx, 1).Return(entity.Order{}, entity.ErrOrderNoExists)
	r.On("CreateOrder", ctx, entity.Order{ID: 1, ServiceID: 1, UserID: 1, Sum: "200"}).
		Return(nil)

	r.On("GetByID", ctx, 2).Return(entity.Balance{ID: 2, Amount: "300"}, nil)
	r.On("CheckServiceID", ctx, 2).Return(nil)
	r.On("GetOrderByID", ctx, 2).
		Return(entity.Order{ID: 2, ServiceID: 2, UserID: 2, Sum: "200"}, nil)

	r.On("GetByID", ctx, 3).Return(entity.Balance{}, entity.ErrNoID)

	r.On("GetByID", ctx, 4).Return(entity.Balance{ID: 4, Amount: "100"}, nil)

	r.On("GetByID", ctx, 5).Return(entity.Balance{ID: 5, Amount: "300"}, nil)
	r.On("CheckServiceID", ctx, 5).Return(entity.ErrNoService)

	type TestCase struct {
		name        string
		val         entity.Order
		expectedErr error
	}

	cases := []TestCase{{
		name:        "valid",
		val:         entity.Order{ID: 1, ServiceID: 1, UserID: 1, Sum: "200"},
		expectedErr: nil,
	}, {
		name:        "already exist",
		val:         entity.Order{ID: 2, ServiceID: 2, UserID: 2, Sum: "200"},
		expectedErr: entity.ErrOrderExists,
	}, {
		name:        "no such user",
		val:         entity.Order{ID: 3, ServiceID: 3, UserID: 3, Sum: "200"},
		expectedErr: entity.ErrNoID,
	}, {
		name:        "not enough money",
		val:         entity.Order{ID: 4, ServiceID: 4, UserID: 4, Sum: "200"},
		expectedErr: entity.ErrNotEnoughMoney,
	}, {
		name:        "no such service",
		val:         entity.Order{ID: 5, ServiceID: 5, UserID: 5, Sum: "200"},
		expectedErr: entity.ErrNoService,
	},
	}

	for _, tc := range cases {
		err := uc.CreateOrder(ctx, tc.val)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestChangeOrderStatus(t *testing.T) {
	ctx := context.Background()
	r := repomock.NewBalanceRepo(t)
	uc := New(r, reportmock.NewReportFile(t))

	r.On("GetOrderByID", ctx, 1).
		Return(entity.Order{ID: 1, ServiceID: 1, UserID: 1, Sum: "200", StatusID: 1}, nil)
	r.On("CommitOrder", ctx, entity.Order{ID: 1, ServiceID: 1, UserID: 1, Sum: "200", StatusID: 2}).
		Return(nil)

	r.On("GetOrderByID", ctx, 2).
		Return(entity.Order{ID: 2, ServiceID: 2, UserID: 2, Sum: "200", StatusID: 1}, nil)
	r.On("RollbackOrder", ctx, entity.Order{ID: 2, ServiceID: 2, UserID: 2, Sum: "200", StatusID: 3}).
		Return(nil)

	r.On("GetOrderByID", ctx, 3).
		Return(entity.Order{}, entity.ErrOrderNoExists)
	r.On("GetOrderByID", ctx, 4).
		Return(entity.Order{ID: 4, ServiceID: 4, UserID: 4, Sum: "200", StatusID: 1}, nil)
	r.On("GetOrderByID", ctx, 5).
		Return(entity.Order{ID: 5, ServiceID: 5, UserID: 5, Sum: "200", StatusID: 2}, nil)

	type TestCase struct {
		name        string
		val         entity.Order
		expectedErr error
	}

	cases := []TestCase{{
		name:        "valid",
		val:         entity.Order{ID: 1, ServiceID: 1, UserID: 1, Sum: "200", StatusID: 2},
		expectedErr: nil,
	}, {
		name:        "valid rollback",
		val:         entity.Order{ID: 2, ServiceID: 2, UserID: 2, Sum: "200", StatusID: 3},
		expectedErr: nil,
	}, {
		name:        "no such order id",
		val:         entity.Order{ID: 3, ServiceID: 3, UserID: 3, Sum: "200", StatusID: 2},
		expectedErr: entity.ErrOrderNoExists,
	}, {
		name:        "wrong order data",
		val:         entity.Order{ID: 4, ServiceID: 4, UserID: 4, Sum: "300", StatusID: 3},
		expectedErr: entity.ErrOrderMismatch,
	}, {
		name:        "order already committed",
		val:         entity.Order{ID: 5, ServiceID: 5, UserID: 5, Sum: "200", StatusID: 2},
		expectedErr: entity.ErrCantChangeStatus,
	},
	}

	for _, tc := range cases {
		err := uc.ChangeOrderStatus(ctx, tc.val)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestIncrease(t *testing.T) {
	ctx := context.Background()
	r := repomock.NewBalanceRepo(t)
	uc := New(r, reportmock.NewReportFile(t))

	r.On("GetByID", ctx, 1).Return(entity.Balance{}, entity.ErrNoID)
	r.On("CreateUser", ctx, entity.Balance{ID: 1, Amount: "200"}).Return(nil)
	r.On("GetByID", ctx, 2).Return(entity.Balance{ID: 2, Amount: "1"}, nil)
	r.On("Increase", ctx, entity.Balance{ID: 2, Amount: "200"}).Return(nil)

	type TestCase struct {
		name        string
		val         entity.Balance
		expectedErr error
	}

	cases := []TestCase{{
		name:        "valid new user",
		val:         entity.Balance{ID: 1, Amount: "200"},
		expectedErr: nil,
	}, {
		name:        "valid",
		val:         entity.Balance{ID: 2, Amount: "200"},
		expectedErr: nil,
	},
	}

	for _, tc := range cases {
		err := uc.Increase(ctx, tc.val)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestGetHistory(t *testing.T) {
	ctx := context.Background()
	r := repomock.NewBalanceRepo(t)
	uc := New(r, reportmock.NewReportFile(t))

	r.On("GetByID", ctx, 1).Return(entity.Balance{ID: 1, Amount: "200"}, nil)
	r.On("GetHistory", ctx, entity.History{UserID: 1, Limit: 10, OrderBy: "date", Desc: true, Page: 1}).
		Return(entity.History{Orders: []entity.Order{
			{ID: 1, ServiceName: "aboba", Status: "approved", Sum: "10", Time: entity.MyTime{Time: time.Unix(10, 0)}},
			{ID: 2, ServiceName: "aboba", Status: "canceled", Sum: "20", Time: entity.MyTime{Time: time.Unix(10, 0)}},
		}, UserID: 1, Limit: 10, OrderBy: "date", Desc: true, Page: 1}, nil)

	r.On("GetByID", ctx, 2).Return(entity.Balance{}, entity.ErrNoID)

	r.On("GetByID", ctx, 3).Return(entity.Balance{ID: 3, Amount: "200"}, nil)
	r.On("GetHistory", ctx, entity.History{UserID: 3, Limit: 2, OrderBy: "sum", Desc: false, Page: 10}).
		Return(entity.History{}, entity.ErrEmptyPage)

	type TestCase struct {
		name        string
		val         entity.History
		expectedVal entity.History
		expectedErr error
	}

	cases := []TestCase{{
		name: "valid",
		val:  entity.History{UserID: 1, Limit: 10, OrderBy: "date", Desc: true, Page: 1},
		expectedVal: entity.History{Orders: []entity.Order{
			{ID: 1, ServiceName: "aboba", Status: "approved", Sum: "10", Time: entity.MyTime{Time: time.Unix(10, 0)}},
			{ID: 2, ServiceName: "aboba", Status: "canceled", Sum: "20", Time: entity.MyTime{Time: time.Unix(10, 0)}},
		}, UserID: 1, Limit: 10, OrderBy: "date", Desc: true, Page: 1},
		expectedErr: nil,
	}, {
		name:        "no such user",
		val:         entity.History{UserID: 2, Limit: 2, OrderBy: "sum", Desc: false, Page: 1},
		expectedVal: entity.History{},
		expectedErr: entity.ErrNoID,
	}, {
		name:        "empty page",
		val:         entity.History{UserID: 3, Limit: 2, OrderBy: "sum", Desc: false, Page: 10},
		expectedVal: entity.History{},
		expectedErr: entity.ErrEmptyPage,
	},
	}

	for _, tc := range cases {
		history, err := uc.GetHistory(ctx, tc.val)
		assert.Equal(t, tc.expectedVal, history)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestUpdateReport(t *testing.T) {
	ctx := context.Background()
	r := repomock.NewBalanceRepo(t)
	f := reportmock.NewReportFile(t)
	uc := New(r, f)

	r.On("GetReport", ctx, 2022, 9).
		Return(entity.Report{Sums: []entity.SumByService{{Sum: "1", Name: "a"}}}, nil)
	f.On("Create", ctx, "2022-09", entity.Report{Sums: []entity.SumByService{{Sum: "1", Name: "a"}}}).
		Return("2022-09.csv", nil)

	r.On("GetReport", ctx, 1980, 1).Return(entity.Report{Sums: nil}, entity.ErrEmptyReport)

	type TestCase struct {
		name        string
		date        []int
		expectedVal string
		expectedErr error
	}

	cases := []TestCase{{
		name:        "valid",
		date:        []int{2022, 9},
		expectedVal: "2022-09.csv",
		expectedErr: nil,
	}, {
		name:        "empty report",
		date:        []int{1980, 1},
		expectedVal: "",
		expectedErr: entity.ErrEmptyReport,
	},
	}

	for _, tc := range cases {
		name, err := uc.UpdateReport(ctx, tc.date[0], tc.date[1])
		assert.Equal(t, tc.expectedVal, name)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestGetDir(t *testing.T) {
	r := repomock.NewBalanceRepo(t)
	f := reportmock.NewReportFile(t)
	uc := New(r, f)

	f.On("GetDir").Return("reports/")

	assert.Equal(t, "reports/", uc.GetReportDir())
}
