// Code generated by mockery v2.20.0. DO NOT EDIT.

package attorney

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// mockErrorHandler is an autogenerated mock type for the ErrorHandler type
type mockErrorHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockErrorHandler) Execute(_a0 http.ResponseWriter, _a1 *http.Request, _a2 error) {
	_m.Called(_a0, _a1, _a2)
}

type mockConstructorTestingTnewMockErrorHandler interface {
	mock.TestingT
	Cleanup(func())
}

// newMockErrorHandler creates a new instance of mockErrorHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockErrorHandler(t mockConstructorTestingTnewMockErrorHandler) *mockErrorHandler {
	mock := &mockErrorHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
