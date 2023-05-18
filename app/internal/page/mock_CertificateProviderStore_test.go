// Code generated by mockery v2.20.0. DO NOT EDIT.

package page

import (
	context "context"

	actor "github.com/ministryofjustice/opg-modernising-lpa/internal/actor"

	mock "github.com/stretchr/testify/mock"
)

// mockCertificateProviderStore is an autogenerated mock type for the CertificateProviderStore type
type mockCertificateProviderStore struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *mockCertificateProviderStore) Create(_a0 context.Context) (*actor.CertificateProviderProvidedDetails, error) {
	ret := _m.Called(_a0)

	var r0 *actor.CertificateProviderProvidedDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*actor.CertificateProviderProvidedDetails, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *actor.CertificateProviderProvidedDetails); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*actor.CertificateProviderProvidedDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Put provides a mock function with given fields: _a0, _a1
func (_m *mockCertificateProviderStore) Put(_a0 context.Context, _a1 *actor.CertificateProviderProvidedDetails) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *actor.CertificateProviderProvidedDetails) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockCertificateProviderStore interface {
	mock.TestingT
	Cleanup(func())
}

// newMockCertificateProviderStore creates a new instance of mockCertificateProviderStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockCertificateProviderStore(t mockConstructorTestingTnewMockCertificateProviderStore) *mockCertificateProviderStore {
	mock := &mockCertificateProviderStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
