package main

import (
	"fmt"
	"github.com/dmalykh/wlinterpreter"
	"github.com/dmalykh/wlinterpreter/dialect/brainfuck"
	"github.com/dmalykh/wlinterpreter/interpreter"
	"github.com/dmalykh/wlinterpreter/stack/slice"
	"github.com/dmalykh/wlinterpreter/storage/list"
	"log"
	"sync"
)

func main() {
	Examlpe[int32]()
}

func Examlpe[S wlinterpreter.CellSize]() {
	// Get store for internal interpreter @TODO
	var store = list.New()
	// Create interpreter and stack
	var wli = interpreter.NewInterpreter(store)
	var st = slice.NewStack[S](3000)

	// Input and output channels
	var input = make(chan S)
	var output = make(chan S)
	defer func() {
		close(input)
		close(output)
	}()

	// Create brainfuck instance
	bf, err := brainfuck.New[S](st, wli, input, output)
	if err != nil {
		panic(err)
	}

	var wg = new(sync.WaitGroup)
	defer wg.Wait()
	go func() {
		for symbol := range output {
			wg.Add(1)
			func(s S) {
				defer wg.Done()
				r := rune(s)
				st := string(r)
				fmt.Printf("%s", st) //@TODO clean
			}(symbol)
		}
	}()

	//var program = []byte(`++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.`)
	var program = []byte(`++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.`)
	if err := bf.Run(program...); err != nil {
		log.Fatalln(err)
	}
}
