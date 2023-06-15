// Code generated by mockery v2.20.0. DO NOT EDIT.

package donor

import (
	identity "github.com/ministryofjustice/opg-modernising-lpa/app/internal/identity"
	mock "github.com/stretchr/testify/mock"
)

// mockYotiClient is an autogenerated mock type for the YotiClient type
type mockYotiClient struct {
	mock.Mock
}

// IsTest provides a mock function with given fields:
func (_m *mockYotiClient) IsTest() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ScenarioID provides a mock function with given fields:
func (_m *mockYotiClient) ScenarioID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SdkID provides a mock function with given fields:
func (_m *mockYotiClient) SdkID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// User provides a mock function with given fields: _a0
func (_m *mockYotiClient) User(_a0 string) (identity.UserData, error) {
	ret := _m.Called(_a0)

	var r0 identity.UserData
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (identity.UserData, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) identity.UserData); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(identity.UserData)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockYotiClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockYotiClient creates a new instance of mockYotiClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockYotiClient(t mockConstructorTestingTnewMockYotiClient) *mockYotiClient {
	mock := &mockYotiClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
