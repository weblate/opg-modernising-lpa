// Code generated by mockery v2.20.0. DO NOT EDIT.

package secrets

import (
	context "context"

	secretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	mock "github.com/stretchr/testify/mock"
)

// mockSecretsManager is an autogenerated mock type for the secretsManager type
type mockSecretsManager struct {
	mock.Mock
}

// GetSecretValue provides a mock function with given fields: ctx, params, optFns
func (_m *mockSecretsManager) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *secretsmanager.GetSecretValueOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) *secretsmanager.GetSecretValueOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*secretsmanager.GetSecretValueOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockSecretsManager interface {
	mock.TestingT
	Cleanup(func())
}

// newMockSecretsManager creates a new instance of mockSecretsManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockSecretsManager(t mockConstructorTestingTnewMockSecretsManager) *mockSecretsManager {
	mock := &mockSecretsManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
