// Code generated by mockery v2.20.0. DO NOT EDIT.

package attorney

import (
	context "context"

	page "github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	mock "github.com/stretchr/testify/mock"
)

// mockLpaStore is an autogenerated mock type for the LpaStore type
type mockLpaStore struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *mockLpaStore) Create(_a0 context.Context) (*page.Lpa, error) {
	ret := _m.Called(_a0)

	var r0 *page.Lpa
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*page.Lpa, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *page.Lpa); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*page.Lpa)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: _a0
func (_m *mockLpaStore) Get(_a0 context.Context) (*page.Lpa, error) {
	ret := _m.Called(_a0)

	var r0 *page.Lpa
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*page.Lpa, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *page.Lpa); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*page.Lpa)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: _a0
func (_m *mockLpaStore) GetAll(_a0 context.Context) ([]*page.Lpa, error) {
	ret := _m.Called(_a0)

	var r0 []*page.Lpa
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*page.Lpa, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*page.Lpa); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*page.Lpa)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Put provides a mock function with given fields: _a0, _a1
func (_m *mockLpaStore) Put(_a0 context.Context, _a1 *page.Lpa) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *page.Lpa) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockLpaStore interface {
	mock.TestingT
	Cleanup(func())
}

// newMockLpaStore creates a new instance of mockLpaStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockLpaStore(t mockConstructorTestingTnewMockLpaStore) *mockLpaStore {
	mock := &mockLpaStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
