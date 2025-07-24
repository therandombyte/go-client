package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
)

func (s *Server) loggerChain() alice.Chain {
	ac := alice.New(hlog.NewHandler(s.Logger),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Msg("request logged")
		}),
		hlog.RemoteAddrHandler("remove_ip"),
	)

	return ac
}

func (s *Server) authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth Middleware")
		h.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

func (s *Server) handleGetVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Version is 0.1")
}
