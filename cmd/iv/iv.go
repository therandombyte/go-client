package main

import (
	"iv/cmd/iv/app"
)

func main() {
	command := app.NewIVCommand()
	command.Execute()
}
