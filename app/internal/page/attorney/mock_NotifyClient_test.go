// Code generated by mockery v2.20.0. DO NOT EDIT.

package attorney

import (
	context "context"

	notify "github.com/ministryofjustice/opg-modernising-lpa/app/internal/notify"
	mock "github.com/stretchr/testify/mock"
)

// mockNotifyClient is an autogenerated mock type for the NotifyClient type
type mockNotifyClient struct {
	mock.Mock
}

// Email provides a mock function with given fields: ctx, email
func (_m *mockNotifyClient) Email(ctx context.Context, email notify.Email) (string, error) {
	ret := _m.Called(ctx, email)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, notify.Email) (string, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, notify.Email) string); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, notify.Email) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Sms provides a mock function with given fields: ctx, sms
func (_m *mockNotifyClient) Sms(ctx context.Context, sms notify.Sms) (string, error) {
	ret := _m.Called(ctx, sms)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, notify.Sms) (string, error)); ok {
		return rf(ctx, sms)
	}
	if rf, ok := ret.Get(0).(func(context.Context, notify.Sms) string); ok {
		r0 = rf(ctx, sms)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, notify.Sms) error); ok {
		r1 = rf(ctx, sms)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplateID provides a mock function with given fields: id
func (_m *mockNotifyClient) TemplateID(id notify.TemplateId) string {
	ret := _m.Called(id)

	var r0 string
	if rf, ok := ret.Get(0).(func(notify.TemplateId) string); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTnewMockNotifyClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockNotifyClient creates a new instance of mockNotifyClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockNotifyClient(t mockConstructorTestingTnewMockNotifyClient) *mockNotifyClient {
	mock := &mockNotifyClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
