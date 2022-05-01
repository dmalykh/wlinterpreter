// Interpreter implements functions to register and execute operator for manipulations with stack.
// The RegisterOperator method saves relation between operator and ExecuteFunc that should be calls when operator
// would be detected in execution:
//		RegisterOperator('^', func(i Interpreter) error {
//			i.Stack()....do something...
//			return nil
//		})
//
//		Exec('^') // something will be done
//
// For manipulations with stack while operator executes use Stack interface.
//
// Also, if you need to store some data while Interpreter works you should use InternalStorage() with Storage interface.
//

package wlinterpreter

import (
	"errors"
)

type ExecuteFunc func(i Interpreter) error

var ErrUnknownOperator = errors.New(`got unknown operator`)

type Interpreter interface {
	// RegisterOperator creates relation between operator and execute. When operator will be gotten at execution the ExecuteFunc will be called.
	RegisterOperator(operator byte, execute ExecuteFunc) error
	// Run full program. Calls the Exec method every byte in program
	Run(program ...byte) error
	// Exec executes operators registered function or return error if operator didn't register
	Exec(operator byte) error
	// InternalStorage return Storage interface for storing data while Interpreter works
	InternalStorage() Storage
	// GetHistory returns last executed operators slice. If last < 0 returns full history slice.
	GetHistory(last int) []byte
	// Frame current interpreter and use it. Forks allows to use independent history, that copied for all parents. You can use frames for loops.
	Frame()
	// Done closes current frame. If parent exists, current interpreter is parent.
	Done() error
	// Clone current interpreter with its settings except history and return it.
	Clone() Interpreter
}
