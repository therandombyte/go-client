package main

import (
	"fmt"
	"iv/cmd/iv/app"
)

func main() {
	fmt.Println("First Invocation")
	command := app.NewIVCommand()
	fmt.Println(command)
}
