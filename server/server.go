package server

import (
	"log"
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
	log.Printf("server starting on %s", port)
	return s.server.ListenAndServe()
}

func healthCheck() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := RM{
			Status:  "success",
			Message: "healthy",
			Data:    M{"hello": "world"},
		}
		writeJSON(rw, http.StatusOK, resp)
	})
}

func errorCheck() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := Error{
			Code:    "internal",
			Message: "internal server error.",
		}

		errorResponse(rw, http.StatusInternalServerError, err)
	})
}
