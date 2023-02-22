// Code generated by mockery v2.20.0. DO NOT EDIT.

package page

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// mockTemplate is an autogenerated mock type for the Template type
type mockTemplate struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1
func (_m *mockTemplate) Execute(_a0 io.Writer, _a1 interface{}) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer, interface{}) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockTemplate interface {
	mock.TestingT
	Cleanup(func())
}

// newMockTemplate creates a new instance of mockTemplate. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockTemplate(t mockConstructorTestingTnewMockTemplate) *mockTemplate {
	mock := &mockTemplate{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
