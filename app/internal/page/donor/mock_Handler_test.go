// Code generated by mockery v2.20.0. DO NOT EDIT.

package donor

import (
	http "net/http"

	page "github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	mock "github.com/stretchr/testify/mock"
)

// mockHandler is an autogenerated mock type for the Handler type
type mockHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: data, w, r
func (_m *mockHandler) Execute(data page.AppData, w http.ResponseWriter, r *http.Request) error {
	ret := _m.Called(data, w, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(page.AppData, http.ResponseWriter, *http.Request) error); ok {
		r0 = rf(data, w, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockHandler interface {
	mock.TestingT
	Cleanup(func())
}

// newMockHandler creates a new instance of mockHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockHandler(t mockConstructorTestingTnewMockHandler) *mockHandler {
	mock := &mockHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
