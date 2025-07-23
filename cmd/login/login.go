package login

import (
	"fmt"

	"github.com/spf13/cobra"
)


func NewLoginCommand() *cobra.Command {
	fmt.Println("Login here")

	cmd := &cobra.Command {
		Use: "iv login",
		Short: "iv login will log into the REST server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Login subcommand called")
		},
	}
	
	return cmd
}
