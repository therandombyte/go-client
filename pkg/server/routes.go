package server

const (
	vraHostPath string = "http://localhost:8080/api"
)

func (s *Server) registerRoutes() {
	s.mux.Handle("GET /api",
		s.loggerChain().
			Append(s.authHandler).
			ThenFunc(s.handleGetVersion))
}
