package app

import (
	"fmt"
	"iv/cmd/login"

	"github.com/spf13/cobra"
)

func init() {
	
}

func NewIVCommand() *cobra.Command {
	// fmt.Println("Init here")

	cmd := &cobra.Command {
		Use: "iv",
		Short: "iv is a go client to make REST api calls to server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Root iv command called")
		},
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	fmt.Println("RunE block")
		// 	return Run()
		// },
	}

	login := login.NewLoginCommand()

	cmd.AddCommand(login)
	return cmd
}

func Run() error{
	fmt.Println("In Run Block")
	return nil
}


