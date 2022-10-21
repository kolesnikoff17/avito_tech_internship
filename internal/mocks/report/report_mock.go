// Code generated by mockery v2.14.0. DO NOT EDIT.

package reportmock

import (
	entity "balance_api/internal/entity"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ReportFile is an autogenerated mock type for the ReportFile type
type ReportFile struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, name, report
func (_m *ReportFile) Create(ctx context.Context, name string, report entity.Report) (string, error) {
	ret := _m.Called(ctx, name, report)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, entity.Report) string); ok {
		r0 = rf(ctx, name, report)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, entity.Report) error); ok {
		r1 = rf(ctx, name, report)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewReportFile interface {
	mock.TestingT
	Cleanup(func())
}

// NewReportFile creates a new instance of ReportFile. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReportFile(t mockConstructorTestingTNewReportFile) *ReportFile {
	mock := &ReportFile{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}