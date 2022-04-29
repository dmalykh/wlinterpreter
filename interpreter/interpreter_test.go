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
	var stack = mocks.NewStack(t)
	t.Run(`use constructor`, func(t *testing.T) {
		var i = NewInterpreter(stack, storage)
		assert.Equal(t, storage, i.InternalStorage())
		assert.Equal(t, stack, i.Stack())
		assert.NoError(t, i.RegisterOperator(byte('o'), func(i wlinterpreter.Interpreter) error { return nil }))
	})
	t.Run(`without constructor`, func(t *testing.T) {
		var i = &Interpreter{
			stack:           stack,
			internalStorage: storage,
		}
		assert.Equal(t, storage, i.InternalStorage())
		assert.Equal(t, stack, i.Stack())
		assert.ErrorIs(t, ErrNilOperatorsMap, i.RegisterOperator(byte('o'), func(i wlinterpreter.Interpreter) error { return nil }))
	})
}

func TestInterpreter_Exec(t *testing.T) {
	type fields struct {
		operators map[byte]wlinterpreter.ExecuteFunc
		stack     wlinterpreter.Stack
	}
	type args struct {
		operator byte
		call     byte
		execute  wlinterpreter.ExecuteFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			`Operator exists and works`,
			args{
				operator: byte('*'),
				call:     byte('*'),
				execute: func(i wlinterpreter.Interpreter) error {
					i.Stack().Move(0)
					return nil
				},
			},
			nil,
		},
		{
			`Operator doesn't exists`,
			args{
				operator: byte('*'),
				call:     byte('-'),
				execute: func(i wlinterpreter.Interpreter) error {
					i.Stack().Move(0)
					return nil
				},
			},
			wlinterpreter.ErrUnknownOperator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stack = mocks.NewStack(t)
			if tt.wantErr == nil {
				stack.On("Move", mock.Anything).Return(tt.args.execute)
			}
			i := &Interpreter{
				operators: make(map[byte]wlinterpreter.ExecuteFunc),
				stack:     stack,
			}
			assert.NoError(t, i.RegisterOperator(tt.args.operator, tt.args.execute))
			err := i.Exec(tt.args.call)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				stack.AssertNotCalled(t, `Move`, mock.Anything)
			} else {
				stack.AssertCalled(t, `Move`, mock.Anything)
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
			var i = NewInterpreter(nil, nil)
			for _, op := range tt.register {
				assert.NoError(t, i.RegisterOperator(op, func(i wlinterpreter.Interpreter) error { return nil }))
			}
			assert.ErrorIs(t, i.Run(tt.check...), tt.wantErr)
		})
	}

}
