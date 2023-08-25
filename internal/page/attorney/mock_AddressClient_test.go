// Code generated by mockery v2.20.0. DO NOT EDIT.

package attorney

import (
	context "context"

	place "github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	mock "github.com/stretchr/testify/mock"
)

// mockAddressClient is an autogenerated mock type for the AddressClient type
type mockAddressClient struct {
	mock.Mock
}

// LookupPostcode provides a mock function with given fields: ctx, postcode
func (_m *mockAddressClient) LookupPostcode(ctx context.Context, postcode string) ([]place.Address, error) {
	ret := _m.Called(ctx, postcode)

	var r0 []place.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]place.Address, error)); ok {
		return rf(ctx, postcode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []place.Address); ok {
		r0 = rf(ctx, postcode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]place.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, postcode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockAddressClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockAddressClient creates a new instance of mockAddressClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockAddressClient(t mockConstructorTestingTnewMockAddressClient) *mockAddressClient {
	mock := &mockAddressClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}