package server

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func Logger(w io.Writer) func(h http.Handler) http.Handler {
	return (func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(w, h)
	})
}

func (s *Server) routes() {
	s.router.Use(Logger(os.Stdout))
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	{
		apiRouter.Handle("/health", healthCheck())
		apiRouter.Handle("/users", s.createUser())
		apiRouter.Handle("/users/login", s.loginUser())
	}
}
