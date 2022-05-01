package interpreter

import (
	"errors"
	"fmt"
	"github.com/dmalykh/wlinterpreter"
)

var ErrNilOperatorsMap = errors.New(`can't register operator, probably you doesn't use NewInterpreter constructor`)

type Interpreter struct {
	operators       map[byte]wlinterpreter.ExecuteFunc
	internalStorage wlinterpreter.Storage
	history         []byte
	parent          *Interpreter
}

func (i *Interpreter) GetHistory(last int) []byte {
	var offset = len(i.history) - last
	if last < 0 || offset < 0 {
		offset = 0
	}
	return i.history[offset:]
}

// NewInterpreter is constructor for Interpreter
func NewInterpreter(storage wlinterpreter.Storage) wlinterpreter.Interpreter {
	return &Interpreter{
		operators:       make(map[byte]wlinterpreter.ExecuteFunc),
		internalStorage: storage,
		history:         make([]byte, 0),
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
	return i.Done()
}

func (i *Interpreter) Exec(operator byte) error {
	defer func() {
		//@TOD: make in goroutine, don't forget locks
		var fork = i
		for {
			fork.history = append(fork.history, operator)
			if fork.parent == nil {
				return
			}
			fork = fork.parent
		}
	}()
	f, exists := i.operators[operator]
	if !exists {
		return fmt.Errorf(`%w "%x (%s)"`, wlinterpreter.ErrUnknownOperator, operator, string(operator))
	}
	return f(i)
}

func (i *Interpreter) InternalStorage() wlinterpreter.Storage {
	return i.internalStorage
}

func (i *Interpreter) Fork() {
	var parent = *i
	var fork = Interpreter{
		operators:       i.operators,
		internalStorage: i.internalStorage,
		history:         make([]byte, 0),
		parent:          &parent,
	}
	*i = fork
}

func (i *Interpreter) Clone() wlinterpreter.Interpreter {
	return &Interpreter{
		operators:       i.operators,
		internalStorage: i.internalStorage,
		history:         make([]byte, 0),
	}
}

func (i *Interpreter) Done() error {
	var parent = i.parent
	if i.parent != nil {
		*i = *parent
	}
	return nil
}
