package main

import (
	"iv/cmd/iv/app"
	"os"
)

func main() {
	command := app.NewIVCommand(os.Args)
	command.Execute()
	// os.Exit(1)
}
