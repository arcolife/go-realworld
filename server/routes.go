package server

import (
	"os"
)

const MustAuth = true

func (s *Server) routes() {
	s.router.Use(Logger(os.Stdout))
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	{
		apiRouter.Handle("/health", healthCheck())
		apiRouter.Handle("/users", s.createUser())
		apiRouter.Handle("/users/login", s.loginUser())
	}

	authApiRoutes := apiRouter.PathPrefix("").Subrouter()
	authApiRoutes.Use(s.authenticate(MustAuth))

	{
		authApiRoutes.Handle("/user", s.getCurrentUser())

	}
}
