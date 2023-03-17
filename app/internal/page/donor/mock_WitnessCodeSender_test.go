// Code generated by mockery v2.20.0. DO NOT EDIT.

package donor

import (
	context "context"

	page "github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	mock "github.com/stretchr/testify/mock"
)

// mockWitnessCodeSender is an autogenerated mock type for the WitnessCodeSender type
type mockWitnessCodeSender struct {
	mock.Mock
}

// Send provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockWitnessCodeSender) Send(_a0 context.Context, _a1 *page.Lpa, _a2 page.Localizer) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *page.Lpa, page.Localizer) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockWitnessCodeSender interface {
	mock.TestingT
	Cleanup(func())
}

// newMockWitnessCodeSender creates a new instance of mockWitnessCodeSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockWitnessCodeSender(t mockConstructorTestingTnewMockWitnessCodeSender) *mockWitnessCodeSender {
	mock := &mockWitnessCodeSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
