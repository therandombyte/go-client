package driver

import (
	"context"
	"net/http"
)

// package driver defines interface for custom HTTP listeners
// Application code should use package server

// Server dispaches request to http.Handler
type Server interface {
	ListenAndServe(addr string, h http.Handler) error
	Shutdown(ctx context.Context) error
}
