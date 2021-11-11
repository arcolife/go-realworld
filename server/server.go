package server

import (
	"log"
	"net/http"
	"time"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/0xdod/go-realworld/inmem"
	"github.com/gorilla/mux"
)

type Server struct {
	server      *http.Server
	router      *mux.Router
	userService conduit.UserService
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
	s.userService = &inmem.UserService{}

	s.server.Handler = s.router

	return &s
}

func (s *Server) Run(port string) error {
	s.server.Addr = port
	log.Printf("server starting on %s", port)
	return s.server.ListenAndServe()
}

func healthCheck() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := M{
			"status":  "available",
			"message": "healthy",
			"data":    M{"hello": "beautiful"},
		}
		writeJSON(rw, http.StatusOK, resp)
	})
}
