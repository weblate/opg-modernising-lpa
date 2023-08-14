// Code generated by mockery v2.20.0. DO NOT EDIT.

package app

import (
	context "context"

	uid "github.com/ministryofjustice/opg-modernising-lpa/app/internal/uid"
	mock "github.com/stretchr/testify/mock"
)

// mockUidClient is an autogenerated mock type for the UidClient type
type mockUidClient struct {
	mock.Mock
}

// CreateCase provides a mock function with given fields: _a0, _a1
func (_m *mockUidClient) CreateCase(_a0 context.Context, _a1 *uid.CreateCaseRequestBody) (uid.CreateCaseResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 uid.CreateCaseResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *uid.CreateCaseRequestBody) (uid.CreateCaseResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *uid.CreateCaseRequestBody) uid.CreateCaseResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(uid.CreateCaseResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *uid.CreateCaseRequestBody) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockUidClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockUidClient creates a new instance of mockUidClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockUidClient(t mockConstructorTestingTnewMockUidClient) *mockUidClient {
	mock := &mockUidClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
