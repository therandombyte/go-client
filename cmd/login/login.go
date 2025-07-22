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
	}
	
	return cmd
}
