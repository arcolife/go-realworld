package server

import "net/http"

func (s *Server) createUser() http.HandlerFunc {
	type input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}
