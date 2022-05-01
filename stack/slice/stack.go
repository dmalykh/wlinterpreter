// Package slice implements wlinterpreter.Stack interface using slice.
package slice

import (
	"fmt"
	"github.com/dmalykh/wlinterpreter"
)

type Stack[S wlinterpreter.CellSize] struct {
	stack   []S
	current int
}

func NewStack[S wlinterpreter.CellSize](size int) wlinterpreter.Stack[S] {
	var s = &Stack[S]{
		stack: make([]S, size),
	}
	return s
}

func (s *Stack[S]) GetValue(position int) (S, error) {
	if position >= len(s.stack) || position < 0 {
		return 0, fmt.Errorf(`%d: %w`, position, wlinterpreter.ErrStackOverflowed)
	}
	return s.stack[position], nil
}

func (s *Stack[S]) SetValue(value S) error {
	s.stack[s.current] = value
	return nil
}

func (s *Stack[S]) Move(delta int) error {
	if s.current+delta >= len(s.stack) || s.current+delta < 0 {
		return fmt.Errorf(`%d + %d: %w`, s.current, delta, wlinterpreter.ErrStackOverflowed)
	}
	s.current = s.current + delta
	return nil
}

func (s *Stack[S]) GetPosition() int {
	return s.current
}

func (s *Stack[S]) SetPosition(pos int) error {
	if pos >= len(s.stack) || pos < 0 {
		return fmt.Errorf(`%d: %w`, pos, wlinterpreter.ErrStackOverflowed)
	}
	s.current = pos
	return nil
}
