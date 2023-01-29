package main

import (
	"fmt"

	"local.packages/token"
)

func main() {
	t := token.Token{Type: "type", Literal: "literal"}
	fmt.Println(t.Type)
	fmt.Println(t.Literal)
}
