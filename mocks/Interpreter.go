// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	wlinterpreter "github.com/dmalykh/wlinterpreter"
)

// Interpreter is an autogenerated mock type for the Interpreter type
type Interpreter struct {
	mock.Mock
}

// Clone provides a mock function with given fields:
func (_m *Interpreter) Clone() wlinterpreter.Interpreter {
	ret := _m.Called()

	var r0 wlinterpreter.Interpreter
	if rf, ok := ret.Get(0).(func() wlinterpreter.Interpreter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(wlinterpreter.Interpreter)
		}
	}

	return r0
}

// Done provides a mock function with given fields:
func (_m *Interpreter) Done() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exec provides a mock function with given fields: operator
func (_m *Interpreter) Exec(operator byte) error {
	ret := _m.Called(operator)

	var r0 error
	if rf, ok := ret.Get(0).(func(byte) error); ok {
		r0 = rf(operator)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fork provides a mock function with given fields:
func (_m *Interpreter) Frame() {
	_m.Called()
}

// GetHistory provides a mock function with given fields: last
func (_m *Interpreter) GetHistory(last int) []byte {
	ret := _m.Called(last)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(int) []byte); ok {
		r0 = rf(last)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// InternalStorage provides a mock function with given fields:
func (_m *Interpreter) InternalStorage() wlinterpreter.Storage {
	ret := _m.Called()

	var r0 wlinterpreter.Storage
	if rf, ok := ret.Get(0).(func() wlinterpreter.Storage); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(wlinterpreter.Storage)
		}
	}

	return r0
}

// RegisterOperator provides a mock function with given fields: operator, execute
func (_m *Interpreter) RegisterOperator(operator byte, execute wlinterpreter.ExecuteFunc) error {
	ret := _m.Called(operator, execute)

	var r0 error
	if rf, ok := ret.Get(0).(func(byte, wlinterpreter.ExecuteFunc) error); ok {
		r0 = rf(operator, execute)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: program
func (_m *Interpreter) Run(program ...byte) error {
	_va := make([]interface{}, len(program))
	for _i := range program {
		_va[_i] = program[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...byte) error); ok {
		r0 = rf(program...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewInterpreter creates a new instance of Interpreter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewInterpreter(t testing.TB) *Interpreter {
	mock := &Interpreter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}