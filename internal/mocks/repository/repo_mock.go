// Code generated by mockery v2.14.0. DO NOT EDIT.

package repomock

import (
  entity "balance_api/internal/entity"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// BalanceRepo is an autogenerated mock type for the BalanceRepo type
type BalanceRepo struct {
	mock.Mock
}

// ChangeOrderStatus provides a mock function with given fields: ctx, order
func (_m *BalanceRepo) ChangeOrderStatus(ctx context.Context, order entity.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateOrder provides a mock function with given fields: ctx, order
func (_m *BalanceRepo) CreateOrder(ctx context.Context, order entity.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *BalanceRepo) GetByID(ctx context.Context, id int) (entity.Balance, error) {
	ret := _m.Called(ctx, id)

	var r0 entity.Balance
	if rf, ok := ret.Get(0).(func(context.Context, int) entity.Balance); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(entity.Balance)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistory provides a mock function with given fields: ctx, history
func (_m *BalanceRepo) GetHistory(ctx context.Context, history entity.History) (entity.History, error) {
	ret := _m.Called(ctx, history)

	var r0 entity.History
	if rf, ok := ret.Get(0).(func(context.Context, entity.History) entity.History); ok {
		r0 = rf(ctx, history)
	} else {
		r0 = ret.Get(0).(entity.History)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, entity.History) error); ok {
		r1 = rf(ctx, history)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetReport provides a mock function with given fields: ctx, year, month
func (_m *BalanceRepo) GetReport(ctx context.Context, year int, month int) (entity.Report, error) {
	ret := _m.Called(ctx, year, month)

	var r0 entity.Report
	if rf, ok := ret.Get(0).(func(context.Context, int, int) entity.Report); ok {
		r0 = rf(ctx, year, month)
	} else {
		r0 = ret.Get(0).(entity.Report)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, year, month)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Increase provides a mock function with given fields: ctx, balance
func (_m *BalanceRepo) Increase(ctx context.Context, balance entity.Balance) error {
	ret := _m.Called(ctx, balance)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Balance) error); ok {
		r0 = rf(ctx, balance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyOrder provides a mock function with given fields: ctx, order
func (_m *BalanceRepo) VerifyOrder(ctx context.Context, order entity.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewBalanceRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewBalanceRepo creates a new instance of BalanceRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBalanceRepo(t mockConstructorTestingTNewBalanceRepo) *BalanceRepo {
	mock := &BalanceRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
