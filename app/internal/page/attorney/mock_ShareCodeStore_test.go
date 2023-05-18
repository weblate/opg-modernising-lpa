// Code generated by mockery v2.20.0. DO NOT EDIT.

package attorney

import (
	context "context"

	actor "github.com/ministryofjustice/opg-modernising-lpa/internal/actor"

	mock "github.com/stretchr/testify/mock"
)

// mockShareCodeStore is an autogenerated mock type for the ShareCodeStore type
type mockShareCodeStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockShareCodeStore) Get(_a0 context.Context, _a1 actor.Type, _a2 string) (actor.ShareCodeData, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 actor.ShareCodeData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, actor.Type, string) (actor.ShareCodeData, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, actor.Type, string) actor.ShareCodeData); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(actor.ShareCodeData)
	}

	if rf, ok := ret.Get(1).(func(context.Context, actor.Type, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockShareCodeStore interface {
	mock.TestingT
	Cleanup(func())
}

// newMockShareCodeStore creates a new instance of mockShareCodeStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockShareCodeStore(t mockConstructorTestingTnewMockShareCodeStore) *mockShareCodeStore {
	mock := &mockShareCodeStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
