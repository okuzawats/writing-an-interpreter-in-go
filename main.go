package main

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"okuzawats.com/go/evaluator"
	"okuzawats.com/go/lexer"
	"okuzawats.com/go/object"
	"okuzawats.com/go/parser"

	"okuzawats.com/go/repl"
)

func main() {
	// 引数を受け取る。1つ目の引数は捨てる。
	args := os.Args

	if len(args) == 1 {
		// REPLの起動
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to tyep in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	} else if len(args) == 2 {
		// プログラムファイルの実行
		file, err := os.Open(args[1])
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		out := os.Stdout

		l := lexer.New(string(content))
		p := parser.New(l)

		evaluated := evaluator.Eval(p.ParseProgram(), object.NewEnvironment())

		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
