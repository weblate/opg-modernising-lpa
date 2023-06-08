// Code generated by mockery v2.20.0. DO NOT EDIT.

package page

import (
	context "context"

	onelogin "github.com/ministryofjustice/opg-modernising-lpa/internal/onelogin"
	mock "github.com/stretchr/testify/mock"
)

// mockOneLoginClient is an autogenerated mock type for the OneLoginClient type
type mockOneLoginClient struct {
	mock.Mock
}

// AuthCodeURL provides a mock function with given fields: state, nonce, locale, identity
func (_m *mockOneLoginClient) AuthCodeURL(state string, nonce string, locale string, identity bool) string {
	ret := _m.Called(state, nonce, locale, identity)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, bool) string); ok {
		r0 = rf(state, nonce, locale, identity)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// EndSessionURL provides a mock function with given fields: idToken, postLogoutURL
func (_m *mockOneLoginClient) EndSessionURL(idToken string, postLogoutURL string) string {
	ret := _m.Called(idToken, postLogoutURL)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(idToken, postLogoutURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Exchange provides a mock function with given fields: ctx, code, nonce
func (_m *mockOneLoginClient) Exchange(ctx context.Context, code string, nonce string) (string, string, error) {
	ret := _m.Called(ctx, code, nonce)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, string, error)); ok {
		return rf(ctx, code, nonce)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, code, nonce)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, code, nonce)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, code, nonce)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UserInfo provides a mock function with given fields: ctx, accessToken
func (_m *mockOneLoginClient) UserInfo(ctx context.Context, accessToken string) (onelogin.UserInfo, error) {
	ret := _m.Called(ctx, accessToken)

	var r0 onelogin.UserInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (onelogin.UserInfo, error)); ok {
		return rf(ctx, accessToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) onelogin.UserInfo); ok {
		r0 = rf(ctx, accessToken)
	} else {
		r0 = ret.Get(0).(onelogin.UserInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accessToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockOneLoginClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockOneLoginClient creates a new instance of mockOneLoginClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockOneLoginClient(t mockConstructorTestingTnewMockOneLoginClient) *mockOneLoginClient {
	mock := &mockOneLoginClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
