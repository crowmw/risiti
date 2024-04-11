// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	model "github.com/crowmw/risiti/model"
	mock "github.com/stretchr/testify/mock"
)

// IReceiptService is an autogenerated mock type for the IReceiptService type
type IReceiptService struct {
	mock.Mock
}

// Create provides a mock function with given fields: receipt
func (_m *IReceiptService) Create(receipt model.Receipt) (model.Receipt, error) {
	ret := _m.Called(receipt)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 model.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func(model.Receipt) (model.Receipt, error)); ok {
		return rf(receipt)
	}
	if rf, ok := ret.Get(0).(func(model.Receipt) model.Receipt); ok {
		r0 = rf(receipt)
	} else {
		r0 = ret.Get(0).(model.Receipt)
	}

	if rf, ok := ret.Get(1).(func(model.Receipt) error); ok {
		r1 = rf(receipt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadAll provides a mock function with given fields:
func (_m *IReceiptService) ReadAll() ([]model.Receipt, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadAll")
	}

	var r0 []model.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.Receipt, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.Receipt); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Receipt)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadByName provides a mock function with given fields: name
func (_m *IReceiptService) ReadByName(name string) (model.Receipt, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for ReadByName")
	}

	var r0 model.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (model.Receipt, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) model.Receipt); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(model.Receipt)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadByText provides a mock function with given fields: text
func (_m *IReceiptService) ReadByText(text string) ([]model.Receipt, error) {
	ret := _m.Called(text)

	if len(ret) == 0 {
		panic("no return value specified for ReadByText")
	}

	var r0 []model.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]model.Receipt, error)); ok {
		return rf(text)
	}
	if rf, ok := ret.Get(0).(func(string) []model.Receipt); ok {
		r0 = rf(text)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Receipt)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIReceiptService creates a new instance of IReceiptService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIReceiptService(t interface {
	mock.TestingT
	Cleanup(func())
}) *IReceiptService {
	mock := &IReceiptService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}