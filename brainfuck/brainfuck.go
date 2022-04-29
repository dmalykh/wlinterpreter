// Package brainfuck implements functions to manipulate Interpreter for interpreting Brainfuck lang.
package brainfuck

import (
	"fmt"
	"github.com/dmalykh/wlinterpreter"
)

const (
	IGNORE  = `ignore`
	SQUARES = `squares`
)

// Brainfuck implements wlinterpreter.Interpreter interface. It uses channels for inputs and outputs.
type Brainfuck[S wlinterpreter.CellSize] struct {
	stack       wlinterpreter.Stack[S]
	interpreter wlinterpreter.Interpreter
	inputChan   chan rune
	outputChan  chan rune
}

func New(interpreter wlinterpreter.Interpreter) (*Brainfuck, error) {

	var bf = &Brainfuck{
		interpreter: interpreter,
		inputChan:   make(chan rune),
		outputChan:  make(chan rune),
	}

	// Alloc blocks in internal storage
	if err := bf.Interpreter().InternalStorage().Alloc(SQUARES); err != nil {
		return nil, fmt.Errorf(`can't alloc %s block; %w`, SQUARES, err)
	}
	if err := bf.Interpreter().InternalStorage().Alloc(IGNORE); err != nil {
		return nil, fmt.Errorf(`can't alloc %s block; %w`, IGNORE, err)
	}

	// Add pointer to right operation
	if err := bf.Interpreter().RegisterOperator('>', func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		if err := i.Stack().Move(1); err != nil {
			return fmt.Errorf(`cann't move 1: %w`, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	// Add pointer to left operation
	if err := bf.Interpreter().RegisterOperator('<', func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		if err := i.Stack().Move(-1); err != nil {
			return fmt.Errorf(`cann't move -1: %w`, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	// Add incrementation operation
	if err := bf.Interpreter().RegisterOperator('+', func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		value, err := i.Stack().GetValue(i.Stack().GetPosition())
		if err != nil {
			return fmt.Errorf(`cann't get value on %d: %w`, i.Stack().GetPosition(), err)
		}
		if err := i.Stack().SetValue(value + 1); err != nil {
			return fmt.Errorf(`cann't set value %d on %d: %w`, value, i.Stack().GetPosition(), err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	// Add output operation
	if err := bf.Interpreter().RegisterOperator('.', func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		value, err := i.Stack().GetValue(i.Stack().GetPosition())
		if err != nil {
			return fmt.Errorf(`cann't get value on %d: %w`, i.Stack().GetPosition(), err)
		}
		bf.outputChan <- rune(value)
		return nil
	}); err != nil {
		return nil, err
	}
	// Add input operation
	if err := bf.Interpreter().RegisterOperator(',', func(i wlinterpreter.Interpreter) error {
		if !i.InternalStorage().Empty(IGNORE) {
			return nil
		}
		i.Stack().SetValue(<-bf.inputChan)
		return nil
	}); err != nil {
		return nil, err
	}

	// Add jump operation
	if err := bf.Interpreter().RegisterOperator('[', func(i wlinterpreter.Interpreter) error {
		// If already ignored, add to ignored another one
		if !i.InternalStorage().Empty(IGNORE) {
			if err := i.InternalStorage().Append(IGNORE, int32(i.Stack().GetPosition())); err != nil {
				return fmt.Errorf(`can't append %s with bracket "[": %w`, IGNORE, err)
			}
			return nil
		}

		// Add to ignored if value in current stack position equals 0
		if i.Stack().GetValue(i.Stack().GetPosition()) == 0 {
			if err := i.InternalStorage().Append(IGNORE, i.Stack().GetPosition()); err != nil {
				return fmt.Errorf(`can't append %s with bracket "[": %w`, IGNORE, err)
			}
			return nil
		}

		// Add to squares if value in current stack position not equals 0
		if err := i.InternalStorage().Append(SQUARES, i.Stack().GetPosition()); err != nil {
			return fmt.Errorf(`can't append %s with bracket "[": %w`, SQUARES, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Add end of jump operator
	if err := bf.Interpreter().RegisterOperator(']', func(i wlinterpreter.Interpreter) error {
		// If in ignored brackets, remove it
		if !i.InternalStorage().Empty(IGNORE) {
			if _, err := i.InternalStorage().Pop(IGNORE); err != nil {
				return fmt.Errorf(`can't pop element for %s: %w`, IGNORE, err)
			}
			return nil
		}
		// Close squares, go to start if it greater then 0
		start, err := i.InternalStorage().Pop(SQUARES)
		if err != nil {
			return fmt.Errorf(`can't pop element for %s: %w`, SQUARES, err)
		}
		if i.Stack().GetValue(start) > 0 {
			i.Stack().SetPosition(start)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return bf, nil
}

// Interpreter returns interpreter
func (b *Brainfuck) Interpreter() wlinterpreter.Interpreter {
	return b.interpreter
}

// Run your program
func (b *Brainfuck) Run(program ...byte) error {
	defer b.shutdown()
	return b.Interpreter().Run(program...)
}

func (b *Brainfuck) Output() chan rune {
	return b.outputChan
}

func (b *Brainfuck) shutdown() {
	close(b.outputChan)
	close(b.inputChan)
}
