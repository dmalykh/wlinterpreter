package interpreter

import (
	"github.com/dmalykh/wlinterpreter"
	"github.com/dmalykh/wlinterpreter/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestInterpreter_RegisterOperator(t *testing.T) {

	tests := []struct {
		name     string
		operator byte
		execute  wlinterpreter.ExecuteFunc
		want     int32
	}{
		{
			`Could add operator`,
			byte('M'),
			func(i wlinterpreter.Interpreter) error {
				return nil
			},
			8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i = &Interpreter{
				operators: make(map[byte]wlinterpreter.ExecuteFunc),
			}
			assert.NoError(t, i.RegisterOperator(tt.operator, tt.execute))
			_, exists := i.operators[tt.operator]
			assert.True(t, exists)
		})
	}
}

func TestNewInterpreter(t *testing.T) {
	var storage = mocks.NewStorage(t)
	t.Run(`use constructor`, func(t *testing.T) {
		var i = NewInterpreter(storage)
		assert.Equal(t, storage, i.InternalStorage())
		assert.NoError(t, i.RegisterOperator(byte('o'), func(i wlinterpreter.Interpreter) error { return nil }))
	})
	t.Run(`without constructor`, func(t *testing.T) {
		var i = &Interpreter{
			internalStorage: storage,
		}
		assert.Equal(t, storage, i.InternalStorage())
		assert.ErrorIs(t, ErrNilOperatorsMap, i.RegisterOperator(byte('o'), func(i wlinterpreter.Interpreter) error { return nil }))
	})
}

func TestInterpreter_Exec(t *testing.T) {

	tests := []struct {
		name     string
		operator byte
		call     byte
		wantErr  error
	}{
		{
			`Operator exists and its func called`,
			byte('*'),
			byte('*'),
			nil,
		},
		{
			`Operator doesn't exists`,
			byte('*'),
			byte('-'),
			wlinterpreter.ErrUnknownOperator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var execFunc = mocks.ExecuteFunc{}
			execFunc.On(`Execute`, mock.Anything).Return(func(i wlinterpreter.Interpreter) error { return nil })

			var i = &Interpreter{
				operators: make(map[byte]wlinterpreter.ExecuteFunc),
			}
			assert.NoError(t, i.RegisterOperator(tt.operator, execFunc.Execute))
			err := i.Exec(tt.call)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				execFunc.AssertNotCalled(t, `Execute`, mock.Anything)
			} else {
				execFunc.AssertNumberOfCalls(t, `Execute`, 1)
			}
		})
	}
}

func TestInterpreter_Run(t *testing.T) {
	tests := []struct {
		name     string
		register []byte
		check    []byte
		wantErr  error
	}{
		{
			`all ok`,
			[]byte{'z', '4', '-', '3', '$'},
			[]byte{'z', '4', '-'},
			nil,
		},
		{
			`unknown operator`,
			[]byte{'z', '4', '-'},
			[]byte{'z', ')', '-'},
			wlinterpreter.ErrUnknownOperator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i = NewInterpreter(nil)
			for _, op := range tt.register {
				assert.NoError(t, i.RegisterOperator(op, func(i wlinterpreter.Interpreter) error { return nil }))
			}
			assert.ErrorIs(t, i.Run(tt.check...), tt.wantErr)
		})
	}

}
