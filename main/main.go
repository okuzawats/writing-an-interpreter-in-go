package main

import (
	"fmt"

	"local.packages/animals"
)

func main() {
	fmt.Println(animals.ElephantFeed())
	fmt.Println(animals.MonkeyFeed())
	fmt.Println(animals.RabbitFeed())
}
