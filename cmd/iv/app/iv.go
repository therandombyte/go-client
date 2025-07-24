package app

import (
	"iv/cmd/login"
	"iv/pkg/server"

	"github.com/spf13/cobra"
)

func init() {

}

func NewIVCommand(args []string) *cobra.Command {
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

	return server.RunServer() 
	// lgr := logging.InitLogger()
	// lgr.Info().Msgf("Logging Initialized")

	// // server multiplexer is often called router that routes incoming
	// // requests to its handler
	// s := server.New(http.NewServeMux(), server.NewDriver(), lgr)
	// s.Addr = ":8081"
	// errCh := make(chan error, 1)
	// fmt.Println("Starting to serve... ")
	// go func() {
	// 	errCh <- s.ListenAndServe()
	// }()

	// // channel to listen for interrupt signal
	// sigInt := make(chan os.Signal, 1)
	// signal.Notify(sigInt, os.Interrupt, syscall.SIGTERM)



	// return s.ListenAndServe()
}
