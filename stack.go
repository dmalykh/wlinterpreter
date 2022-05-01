package wlinterpreter

import "errors"

var ErrStackOverflowed = errors.New(`stack overflowed`)

// CellSize is interface for cell sizes
type CellSize interface {
	~int8 | ~int32 | ~int64
}

//Stack interface consist of methods for manipulations with stack.
type Stack[S CellSize] interface {
	// GetValue returns value in specified position
	GetValue(position int) (S, error)
	// SetValue set value to current caret position
	SetValue(value S) error
	// Move caret for delta steps. If positive right, if negative left.
	Move(delta int) error
	// GetPosition returns current caret position
	GetPosition() int
	// SetPosition for caret in stack
	SetPosition(pos int) error
}
