package brainfuck

import (
	"fmt"
	"github.com/dmalykh/wlinterpreter/mocks"
	"github.com/dmalykh/wlinterpreter/stack/slice"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBrainfuck_New(t *testing.T) {
	var interpreter = mocks.NewInterpreter(t)
	var stack = slice.NewStack[int32](30000) // Unfortunetly, mockery isn't working with generic types now...
	var inputChan = make(chan int32)
	var outputChan = make(chan int32)

	t.Run(`Error in allocation InternalStorage`, func(t *testing.T) {
		var storage = mocks.NewStorage(t)
		interpreter.On(`InternalStorage`).Return(storage)
		storage.On(`Alloc`, IGNORE).Return(fmt.Errorf(``))
		_, err := New[int32](stack, interpreter, inputChan, outputChan)
		assert.Error(t, err)
	})
}
