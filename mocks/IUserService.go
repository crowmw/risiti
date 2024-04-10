// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	model "github.com/crowmw/risiti/model"
	mock "github.com/stretchr/testify/mock"
)

// IUserService is an autogenerated mock type for the IUserService type
type IUserService struct {
	mock.Mock
}

// AnyExists provides a mock function with given fields:
func (_m *IUserService) AnyExists() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AnyExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: user
func (_m *IUserService) Create(user model.User) (model.User, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(model.User) (model.User, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(model.User) model.User); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	if rf, ok := ret.Get(1).(func(model.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Read provides a mock function with given fields: email
func (_m *IUserService) Read(email string) (model.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (model.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) model.User); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIUserService creates a new instance of IUserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUserService {
	mock := &IUserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
