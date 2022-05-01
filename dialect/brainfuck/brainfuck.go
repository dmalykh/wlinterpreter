// Package brainfuck implements functions to manipulate Interpreter for interpreting Brainfuck lang.
package brainfuck

import (
	"fmt"
	"github.com/dmalykh/wlinterpreter"
	sm "github.com/dmalykh/wlinterpreter/stack"
)

const (
	IGNORE = `ignore`
)

// Brainfuck implements wlinterpreter.Interpreter interface. It uses channels for inputs and outputs.
type Brainfuck[S wlinterpreter.CellSize] struct {
	interpreter wlinterpreter.Interpreter
	stack       wlinterpreter.Stack[S]
}

// New register operators for Brainfuck interpretation and returns created Brainfuck struct.
func New[S wlinterpreter.CellSize](stack wlinterpreter.Stack[S], interpreter wlinterpreter.Interpreter, input chan S, output chan S) (*Brainfuck[S], error) {

	var bf = &Brainfuck[S]{
		interpreter: interpreter,
		stack:       stack,
	}

	// Alloc blocks in internal storage
	if err := bf.Interpreter().InternalStorage().Alloc(IGNORE); err != nil {
		return nil, fmt.Errorf(`cann't alloc %s block; %w`, IGNORE, err)
	}

	// Add pointer to right operation
	if err := bf.registerRightOperator(); err != nil {
		return nil, fmt.Errorf(`cann't register right operator %w`, err)
	}
	// Add pointer to left operation
	if err := bf.registerLeftOperator(); err != nil {
		return nil, fmt.Errorf(`cann't register left operator %w`, err)
	}
	// Add incremental operation
	if err := bf.registerIncrementalOperator(); err != nil {
		return nil, fmt.Errorf(`cann't register incremental operator %w`, err)
	}
	// Add decremental operation
	if err := bf.registerDecrementalOperator(); err != nil {
		return nil, fmt.Errorf(`cann't register decremental operator %w`, err)
	}
	// Add jump operation
	if err := bf.registerStartJumpOperator(); err != nil {
		return nil, fmt.Errorf(`cann't register start jump operator %w`, err)
	}
	// Add end of jump operator
	if err := bf.registerEndJumpOperator(); err != nil {
		return nil, fmt.Errorf(`cann't register end jump operator %w`, err)
	}

	// Add output operation
	if err := bf.Interpreter().RegisterOperator('.', func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		value, err := sm.GetValue(stack, sm.GetPosition(stack))
		if err != nil {
			return fmt.Errorf(`cann't get value on %d: %w`, sm.GetPosition(stack), err)
		}
		output <- S(value)
		return nil
	}); err != nil {
		return nil, err
	}

	// Add input operation
	if err := bf.Interpreter().RegisterOperator(',', func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		var val = <-input
		if err := sm.SetValue[S](stack, val); err != nil {
			return fmt.Errorf(`cann't set value on %d: %w`, sm.GetPosition(stack), err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return bf, nil
}

// LoopMiddleware used for storing loop bracket positions for ignoring operators
func LoopMiddleware(f wlinterpreter.ExecuteFunc) wlinterpreter.ExecuteFunc {
	return func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		return f(i)
	}
}

// Use > to move stacks caret right
func (bf *Brainfuck[S]) registerRightOperator() error {
	return bf.Interpreter().RegisterOperator('>', LoopMiddleware(func(i wlinterpreter.Interpreter) error {
		if err := sm.Move(bf.stack, 1); err != nil {
			return fmt.Errorf(`cann't move 1: %w`, err)
		}
		return nil
	}))
}

// Use < to move stacks caret left
func (bf *Brainfuck[S]) registerLeftOperator() error {
	return bf.Interpreter().RegisterOperator('<', LoopMiddleware(func(i wlinterpreter.Interpreter) error {
		if err := sm.Move(bf.stack, -1); err != nil {
			return fmt.Errorf(`cann't move -1: %w`, err)
		}
		return nil
	}))
}

// Use + to increment stacks current value
func (bf *Brainfuck[S]) registerIncrementalOperator() error {
	return bf.Interpreter().RegisterOperator('+', LoopMiddleware(func(i wlinterpreter.Interpreter) error {
		value, err := sm.GetValue(bf.stack, sm.GetPosition(bf.stack))
		if err != nil {
			return fmt.Errorf(`cann't get value on %d: %w`, sm.GetPosition(bf.stack), err)
		}
		if err := sm.SetValue[S](bf.stack, S(value+1)); err != nil {
			return fmt.Errorf(`cann't set value %d on %d: %w`, value, sm.GetPosition(bf.stack), err)
		}
		return nil
	}))
}

// Use - to decrement stacks current value
func (bf *Brainfuck[S]) registerDecrementalOperator() error {
	return bf.Interpreter().RegisterOperator('-', LoopMiddleware(func(i wlinterpreter.Interpreter) error {
		value, err := sm.GetValue(bf.stack, sm.GetPosition(bf.stack))
		if err != nil {
			return fmt.Errorf(`cann't get value on %d: %w`, sm.GetPosition(bf.stack), err)
		}
		if err := sm.SetValue[S](bf.stack, S(value-1)); err != nil {
			return fmt.Errorf(`cann't set value %d on %d: %w`, value, sm.GetPosition(bf.stack), err)
		}
		return nil
	}))
}

// Use [ to store current position for a jump. If current value is 0, position added to IGNORE store.
func (bf *Brainfuck[S]) registerStartJumpOperator() error {
	return bf.Interpreter().RegisterOperator('[', func(i wlinterpreter.Interpreter) error {
		// If already ignored add another one
		if !i.InternalStorage().Empty(IGNORE) {
			if err := i.InternalStorage().Append(IGNORE, sm.GetPosition(bf.stack)); err != nil {
				return fmt.Errorf(`cann't append %s with bracket "[": %w`, IGNORE, err)
			}
			return nil
		}
		// Get value of current cell
		value, err := sm.GetValue(bf.stack, sm.GetPosition(bf.stack))
		if err != nil {
			return fmt.Errorf(`cann't get value on %d: %w`, sm.GetPosition(bf.stack), err)
		}
		// If value eq zero, ignore all operations while meet ]
		if value == 0 {
			if err := i.InternalStorage().Append(IGNORE, sm.GetPosition(bf.stack)); err != nil {
				return fmt.Errorf(`cann't append %s with bracket "[": %w`, IGNORE, err)
			}
			return nil
		}
		// Start new subprogram
		i.Fork()
		return nil
	})
}

// Use ] to remove position from offsets store and repeat last operations while value of current cell in nonzero
func (bf *Brainfuck[S]) registerEndJumpOperator() error {
	return bf.Interpreter().RegisterOperator(']', func(i wlinterpreter.Interpreter) error {
		// If ignored brackets, remove one go to next operator
		if !i.InternalStorage().Empty(IGNORE) {
			if _, err := i.InternalStorage().Pop(IGNORE); err != nil {
				return fmt.Errorf(`cann't pop %s with bracket "[": %w`, IGNORE, err)
			}
			return nil
		}

		// Run subprogram until value > 0
		for {
			// Make decision on value of current position
			value, err := sm.GetValue(bf.stack, sm.GetPosition(bf.stack))
			if err != nil {
				return fmt.Errorf(`cann't get value on %d: %w`, sm.GetPosition(bf.stack), err)
			}

			// If current value is 0, done current subprogram and go to next operator
			if value == 0 {
				if err := i.Done(); err != nil {
					return fmt.Errorf(`cann't close fork: %w`, err)
				}
				return nil
			}
			// Otherwise run current subprogram in independent runtime
			var subprogram = i.GetHistory(-1)
			if err := i.Clone().Run(subprogram[1:]...); err != nil { // Execute without [
				return fmt.Errorf(`error in subprogram on %s: %w`, string(subprogram), err)
			}
		}
	})
}

// Interpreter returns interpreter
func (b *Brainfuck[S]) Interpreter() wlinterpreter.Interpreter {
	return b.interpreter
}

// Run your program
func (b *Brainfuck[S]) Run(program ...byte) error {
	return b.Interpreter().Run(program...)
}
