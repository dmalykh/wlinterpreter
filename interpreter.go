package wlinterpreter

import (
	"errors"
)

type ExecuteFunc func(i Interpreter) error

var ErrUnknownOperator = errors.New(`got unknown operator`)

// Interpreter implements functions to register and execute operator for manipulations with stack.
// The RegisterOperator method saves relation between operator and ExecuteFunc that should be calls when operator
// would be detected in execution:
//		RegisterOperator('^', func(i Interpreter) error {
//			i.Stack()....do something...
//			return nil
//		})
//
//		Exec('^') // something will done
//
// For manipulations with stack while operator executes use Stack interface.
//
// Also, if you need to store some data while Interpreter works you should use InternalStorage() with Storage interface.
//

type Interpreter interface {
	// RegisterOperator creates relation between operator and execute. When operator will got at execution, the execute will called.
	RegisterOperator(operator byte, execute ExecuteFunc) error
	// Run full program. Calls the Exec method every byte in program
	Run(program ...byte) error
	// Exec executes operators registered function or return error if operator didn't registered
	Exec(operator byte) error
	// InternalStorage return Storage interface for storing data while Interpreter works
	InternalStorage() Storage
	// Stack return Stack interface for manipulations with stack while Interpreter works
	Stack() Stack
}
