package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	router *mux.Router
}

func NewServer() *Server {
	s := Server{
		server: &http.Server{
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
		router: mux.NewRouter().StrictSlash(true),
	}

	s.routes()
	s.server.Handler = s.router

	return &s
}

func (s *Server) Run(port string) error {
	s.server.Addr = port
	return s.server.ListenAndServe()
}
