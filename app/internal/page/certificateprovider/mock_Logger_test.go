// Code generated by mockery v2.20.0. DO NOT EDIT.

package certificateprovider

import mock "github.com/stretchr/testify/mock"

// mockLogger is an autogenerated mock type for the Logger type
type mockLogger struct {
	mock.Mock
}

// Print provides a mock function with given fields: v
func (_m *mockLogger) Print(v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

type mockConstructorTestingTnewMockLogger interface {
	mock.TestingT
	Cleanup(func())
}

// newMockLogger creates a new instance of mockLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockLogger(t mockConstructorTestingTnewMockLogger) *mockLogger {
	mock := &mockLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
