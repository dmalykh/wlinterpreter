package main

import (
	"fmt"
	"github.com/dmalykh/wlinterpreter/brainfuck"
	"github.com/dmalykh/wlinterpreter/interpreter"
	"github.com/dmalykh/wlinterpreter/stack"
	"github.com/dmalykh/wlinterpreter/storage/list"
	"log"
)

func main() {
	var store = list.New()
	var st = stack.NewStack[int32](3000)
	var wli = interpreter.NewInterpreter(st, store)
	bf, err := brainfuck.New(wli)
	if err != nil {
		panic(err)
	}

	go func() {
		for symbol := range bf.Output() {
			fmt.Print(string(symbol))
		}
	}()

	var program = []byte(`++++++++++[>+++++++>++++++++++>+++>+<<<<-]>++.>+.+++++++..+++.>++.<<+++++++++++++++.>.+++.------.--------.>+.>.`)
	if err := bf.Run(program...); err != nil {
		log.Fatalln(err)
	}
}
