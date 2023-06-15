// Code generated by mockery v2.20.0. DO NOT EDIT.

package certificateprovider

import (
	context "context"

	page "github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	mock "github.com/stretchr/testify/mock"
)

// mockDonorStore is an autogenerated mock type for the DonorStore type
type mockDonorStore struct {
	mock.Mock
}

// GetAny provides a mock function with given fields: _a0
func (_m *mockDonorStore) GetAny(_a0 context.Context) (*page.Lpa, error) {
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

type mockConstructorTestingTnewMockDonorStore interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDonorStore creates a new instance of mockDonorStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDonorStore(t mockConstructorTestingTnewMockDonorStore) *mockDonorStore {
	mock := &mockDonorStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
