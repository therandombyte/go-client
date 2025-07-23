package server

import (
	"iv/pkg/server/driver"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// Server is a HTTP server
type Server struct {
	mux    *http.ServeMux
	Driver driver.Server
	Logger zerolog.Logger // to be passed as generics?
	Addr   string
	// Services
}

// Driver implements the driver.Server Interface
type Driver struct {
	Server http.Server
}

func NewDriver() *Driver {
	return &Driver{
		Server: http.Server{
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}

func (d *Driver) ListenAndServe(addr string, h http.Handler) error {
	return nil
}

func (d *Driver) Shutdown() error {
	return nil
}
