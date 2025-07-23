package logging

import (
	"os"

	"github.com/rs/zerolog"
)

// currently the logging level and time format is fixed
// dont want the caller to be bothered by this
// {TODO}: have to make it configurable
func InitLogger() zerolog.Logger {
	lvl, err := zerolog.ParseLevel("Debug")
	if err != nil {
		panic("Error Initializing Logger")
	}
	lgr := zerolog.New(os.Stdout).Level(lvl)
	lgr = lgr.With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return lgr
}
