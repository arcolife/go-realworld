package server

func (s *Server) routes() {
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	{
		apiRouter.Handle("/health", healthCheck())
		apiRouter.Handle("/users", s.createUser())
	}
}
