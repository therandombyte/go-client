package app

import (
	"iv/cmd/login"
	"iv/pkg/logging"

	"github.com/spf13/cobra"
)

func init() {

}

func NewIVCommand(args []string) *cobra.Command {
	// fmt.Println("Init here")

	cmd := &cobra.Command{
		Use:   "iv",
		Short: "iv is a go client to make REST api calls to server",

		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(args)
		},
	}

	login := login.NewLoginCommand()

	cmd.AddCommand(login)
	return cmd
}

func Run(args []string) error {
	lgr := logging.InitLogger()
	lgr.Info().Msgf("Logging Initialized")
	return nil
}
