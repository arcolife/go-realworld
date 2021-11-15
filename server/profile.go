package server

import (
	"fmt"
	"net/http"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/gorilla/mux"
)

func (s *Server) getProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		filter := conduit.UserFilter{}
		filter.Username = &username
		fmt.Println(filter)
		users, err := s.userService.Users(r.Context(), filter)

		if err != nil {
			serverError(w, err)
			return
		}

		fmt.Println(users)

		if len(users) == 0 {
			notFoundError(w)
			return
		}

		writeJSON(w, http.StatusOK, M{"profile": users[0].Profile()})
	}
}
