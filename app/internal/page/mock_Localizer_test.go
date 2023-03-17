// Code generated by mockery v2.20.0. DO NOT EDIT.

package page

import mock "github.com/stretchr/testify/mock"

// mockLocalizer is an autogenerated mock type for the Localizer type
type mockLocalizer struct {
	mock.Mock
}

// Count provides a mock function with given fields: messageID, count
func (_m *mockLocalizer) Count(messageID string, count int) string {
	ret := _m.Called(messageID, count)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, int) string); ok {
		r0 = rf(messageID, count)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Format provides a mock function with given fields: _a0, _a1
func (_m *mockLocalizer) Format(_a0 string, _a1 map[string]interface{}) string {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, map[string]interface{}) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// FormatCount provides a mock function with given fields: messageID, count, data
func (_m *mockLocalizer) FormatCount(messageID string, count int, data map[string]interface{}) string {
	ret := _m.Called(messageID, count, data)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, int, map[string]interface{}) string); ok {
		r0 = rf(messageID, count, data)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Possessive provides a mock function with given fields: s
func (_m *mockLocalizer) Possessive(s string) string {
	ret := _m.Called(s)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(s)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetShowTranslationKeys provides a mock function with given fields: s
func (_m *mockLocalizer) SetShowTranslationKeys(s bool) {
	_m.Called(s)
}

// ShowTranslationKeys provides a mock function with given fields:
func (_m *mockLocalizer) ShowTranslationKeys() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// T provides a mock function with given fields: _a0
func (_m *mockLocalizer) T(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTnewMockLocalizer interface {
	mock.TestingT
	Cleanup(func())
}

// newMockLocalizer creates a new instance of mockLocalizer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockLocalizer(t mockConstructorTestingTnewMockLocalizer) *mockLocalizer {
	mock := &mockLocalizer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
