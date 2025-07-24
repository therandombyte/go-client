package server

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
)

func (s *Server) loggerChain() alice.Chain {
	ac := alice.New(hlog.NewHandler(s.Logger),
			hlog.AccessHandler(func (r *http.Request, status, size int, duration time.Duration){
				hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Msg("request logged")
			}),
			hlog.RemoteAddrHandler("remove_ip"),
		)

	return ac
}
