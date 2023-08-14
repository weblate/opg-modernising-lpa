// Code generated by mockery v2.20.0. DO NOT EDIT.

package event

import (
	context "context"

	eventbridge "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	mock "github.com/stretchr/testify/mock"
)

// mockEventbridgeClient is an autogenerated mock type for the eventbridgeClient type
type mockEventbridgeClient struct {
	mock.Mock
}

// PutEvents provides a mock function with given fields: ctx, params, optFns
func (_m *mockEventbridgeClient) PutEvents(ctx context.Context, params *eventbridge.PutEventsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *eventbridge.PutEventsOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *eventbridge.PutEventsInput, ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *eventbridge.PutEventsInput, ...func(*eventbridge.Options)) *eventbridge.PutEventsOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*eventbridge.PutEventsOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *eventbridge.PutEventsInput, ...func(*eventbridge.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockEventbridgeClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockEventbridgeClient creates a new instance of mockEventbridgeClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockEventbridgeClient(t mockConstructorTestingTnewMockEventbridgeClient) *mockEventbridgeClient {
	mock := &mockEventbridgeClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
