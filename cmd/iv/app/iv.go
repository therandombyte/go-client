package app

import (
	"iv/cmd/login"
	"os"

	"github.com/rs/zerolog"
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
	lgr := initLogger()
	lgr.Info().Msgf("Logging Initialized")
	return nil
}

func initLogger() zerolog.Logger {
	lvl, err := zerolog.ParseLevel("Debug")
	if err != nil {
		panic("Error Initializing Logger")
	}
	lgr := zerolog.New(os.Stdout).Level(lvl)
	lgr = lgr.With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return lgr
}
