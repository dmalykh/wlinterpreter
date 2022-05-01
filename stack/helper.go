// Package stack consists of methods for manipulations with stack.
// Methods wrap generic stack interface and get you ability for use generic stack interface with non-generic interfaces.
package stack

import "github.com/dmalykh/wlinterpreter"

// GetValue returns value in specified position of stack
func GetValue[S wlinterpreter.CellSize](stack wlinterpreter.Stack[S], position int) (S, error) {
	return stack.GetValue(position)
}

// SetValue set value to current caret position of stack
func SetValue[S wlinterpreter.CellSize](stack wlinterpreter.Stack[S], value S) error {
	return stack.SetValue(value)
}

// Move caret for delta steps. If positive right, if negative left.
func Move[S wlinterpreter.CellSize](stack wlinterpreter.Stack[S], delta int) error {
	return stack.Move(delta)
}

// GetPosition returns current caret position
func GetPosition[S wlinterpreter.CellSize](stack wlinterpreter.Stack[S]) int {
	return stack.GetPosition()
}

// SetPosition for caret in stack
func SetPosition[S wlinterpreter.CellSize](stack wlinterpreter.Stack[S], pos int) error {
	return stack.SetPosition(pos)
}
