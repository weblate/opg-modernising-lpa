// Code generated by mockery v2.20.0. DO NOT EDIT.

package attorney

import (
	http "net/http"

	sessions "github.com/gorilla/sessions"
	mock "github.com/stretchr/testify/mock"
)

// mockSessionStore is an autogenerated mock type for the SessionStore type
type mockSessionStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: r, name
func (_m *mockSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	ret := _m.Called(r, name)

	var r0 *sessions.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request, string) (*sessions.Session, error)); ok {
		return rf(r, name)
	}
	if rf, ok := ret.Get(0).(func(*http.Request, string) *sessions.Session); ok {
		r0 = rf(r, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sessions.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// New provides a mock function with given fields: r, name
func (_m *mockSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	ret := _m.Called(r, name)

	var r0 *sessions.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request, string) (*sessions.Session, error)); ok {
		return rf(r, name)
	}
	if rf, ok := ret.Get(0).(func(*http.Request, string) *sessions.Session); ok {
		r0 = rf(r, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sessions.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: r, w, s
func (_m *mockSessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	ret := _m.Called(r, w, s)

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request, http.ResponseWriter, *sessions.Session) error); ok {
		r0 = rf(r, w, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockSessionStore interface {
	mock.TestingT
	Cleanup(func())
}

// newMockSessionStore creates a new instance of mockSessionStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockSessionStore(t mockConstructorTestingTnewMockSessionStore) *mockSessionStore {
	mock := &mockSessionStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
