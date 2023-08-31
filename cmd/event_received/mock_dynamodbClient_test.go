// Code generated by mockery v2.20.0. DO NOT EDIT.

package main

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockDynamodbClient is an autogenerated mock type for the dynamodbClient type
type mockDynamodbClient struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, pk, sk, v
func (_m *mockDynamodbClient) Get(ctx context.Context, pk string, sk string, v interface{}) error {
	ret := _m.Called(ctx, pk, sk, v)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, interface{}) error); ok {
		r0 = rf(ctx, pk, sk, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetOneByUID provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockDynamodbClient) GetOneByUID(_a0 context.Context, _a1 string, _a2 interface{}) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Put provides a mock function with given fields: _a0, _a1
func (_m *mockDynamodbClient) Put(_a0 context.Context, _a1 interface{}) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockDynamodbClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDynamodbClient creates a new instance of mockDynamodbClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDynamodbClient(t mockConstructorTestingTnewMockDynamodbClient) *mockDynamodbClient {
	mock := &mockDynamodbClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
