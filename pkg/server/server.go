package server

import (
	"context"
	"fmt"
	"iv/pkg/logging"
	"iv/pkg/server/driver"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

// Server is a HTTP server
type Server struct {
	mux    *http.ServeMux
	Driver driver.Server
	Logger zerolog.Logger // to be passed as generics?
	Addr   string
	Services
}

type Services struct {

}

func New(sm *http.ServeMux, ds driver.Server, lgr zerolog.Logger) *Server {
	return &Server{
		mux:    sm,
		Driver: ds,
		Logger: lgr,
	}
}

func (s *Server) ListenAndServe() error {
	return s.Driver.ListenAndServe(s.Addr, s.mux)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return nil
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
	fmt.Println("Launching server in driver")
	d.Server.Addr = addr
	d.Server.Handler = h
	return d.Server.ListenAndServe()
}

func (d *Driver) Shutdown(ctx context.Context) error {
	return nil
}

func RunServer() error {
	lgr := logging.InitLogger()
	lgr.Info().Msgf("Logging Initialized")
	// server multiplexer is often called router that routes incoming
	// requests to its handler
	s := New(http.NewServeMux(), NewDriver(), lgr)
	s.Addr = ":8081"
	errCh := make(chan error, 1)
	fmt.Println("Starting to serve... ")
	go func() {
		errCh <- s.ListenAndServe()
	}()

	// channel to listen for interrupt signal
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		lgr.Info().Msgf("server start error: %v", err)
	case <- sigInt:
		lgr.Info().Msgf("shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			lgr.Info().Msgf("graceful shutdown error: %v", err)
		}
	}
	return nil
}
