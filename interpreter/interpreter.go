package interpreter

import (
	"errors"
	"fmt"
	"github.com/dmalykh/wlinterpreter"
)

var ErrNilOperatorsMap = errors.New(`can't register operator, probably you doesn't use NewInterpreter constructor`)

type Interpreter struct {
	operators map[byte]wlinterpreter.ExecuteFunc
	//stack           wlinterpreter.Stack
	internalStorage wlinterpreter.Storage
}

// NewInterpreter is constructor for Interpreter
func NewInterpreter(stack wlinterpreter.Stack, storage wlinterpreter.Storage) wlinterpreter.Interpreter {
	return &Interpreter{
		operators: make(map[byte]wlinterpreter.ExecuteFunc),
		//stack:           stack,
		internalStorage: storage,
	}
}

func (i *Interpreter) RegisterOperator(operator byte, execute wlinterpreter.ExecuteFunc) error {
	if i.operators == nil {
		return ErrNilOperatorsMap
	}
	i.operators[operator] = execute
	return nil
}

func (i *Interpreter) Run(program ...byte) error {
	for s, operator := range program {
		if err := i.Exec(operator); err != nil {
			return fmt.Errorf(`error at %d symbol %x (%s): %w`, s, operator, string(operator), err)
		}
	}
	return nil
}

func (i *Interpreter) Exec(operator byte) error {
	f, exists := i.operators[operator]
	if !exists {
		return fmt.Errorf(`%w "%x (%s)"`, wlinterpreter.ErrUnknownOperator, operator, string(operator))
	}
	return f(i)
}

func (i *Interpreter) InternalStorage() wlinterpreter.Storage {
	return i.internalStorage
}

//
//func (i *Interpreter) Stack() wlinterpreter.Stack {
//	return i.stack
//}
