package app

import (
	"fmt"
	"iv/cmd/login"

	"github.com/spf13/cobra"
)

func NewIVCommand() *cobra.Command {
	fmt.Println("Init here")

	cmd := &cobra.Command {
		Use: "iv",
		Short: "iv is a go client to make REST api calls to server",
	}

	login := login.NewLoginCommand()

	cmd.AddCommand(login)
	return cmd
}


