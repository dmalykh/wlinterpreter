# # llinterpreter

Universal interpreter for simple esoteric languages that uses simple commands. Like brainfuck. Interpreter allows to work with any stack you want and with any cell size you need.
Because it uses [generics](https://go.dev/doc/tutorial/generics) for working with stack.


wlinterpreter consists of interfaces:
- **Interpreter** used only for interpretations simple commands.
- **Stack[Cell Size]** used for manipulations with stack. Also, for convenient working with it in non-generic interfaces helpers in stack implemented.
- **Storage** is simple additional FIFO storage that you may be need in development.


### Interpreter
Interpreter has methods for register and execute commands. 
```go
RegisterOperator('>', func(i wlinterpreter.Interpreter) error {
    ... do something ...
    return nil
})
```
and then...
```go
Run([]byte{'>', '>'})
```
Also it allows to Fork current interpreter to independent runtime and Clone it. It may be useful for loops implementations. 


## Brainfuck

For example, brainfuck dialect implemented for this interpreter. 
Finally, you can use it for brainfuck implementation.  

```go

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

	// Create waitgroup for reading output
	var wg = new(sync.WaitGroup)
	defer wg.Wait()
	go func() {
		for symbol := range output {
			wg.Add(1)
			func(s S) {
				defer wg.Done()
				fmt.Printf("%s", string(rune(s)))
			}(symbol)
		}
	}()

	// Run program
	var program = []byte(`++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.`)
	if err := bf.Run(program...); err != nil {
		log.Fatalln(err)
	}
}
```