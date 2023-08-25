// Code generated by mockery v2.20.0. DO NOT EDIT.

package uid

import (
	context "context"

	aws "github.com/aws/aws-sdk-go-v2/aws"

	http "net/http"

	mock "github.com/stretchr/testify/mock"

	time "time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

// mockV4Signer is an autogenerated mock type for the v4Signer type
type mockV4Signer struct {
	mock.Mock
}

// SignHTTP provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7
func (_m *mockV4Signer) SignHTTP(_a0 context.Context, _a1 aws.Credentials, _a2 *http.Request, _a3 string, _a4 string, _a5 string, _a6 time.Time, _a7 ...func(*v4.SignerOptions)) error {
	_va := make([]interface{}, len(_a7))
	for _i := range _a7 {
		_va[_i] = _a7[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1, _a2, _a3, _a4, _a5, _a6)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, aws.Credentials, *http.Request, string, string, string, time.Time, ...func(*v4.SignerOptions)) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockV4Signer interface {
	mock.TestingT
	Cleanup(func())
}

// newMockV4Signer creates a new instance of mockV4Signer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockV4Signer(t mockConstructorTestingTnewMockV4Signer) *mockV4Signer {
	mock := &mockV4Signer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}