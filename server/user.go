package server

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fid reflect.StructField) string {
		name := strings.SplitN(fid.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})
}

func userResponse(_user *conduit.User) M {
	return M{
		"email":    _user.Email,
		"token":    _user.Token,
		"username": _user.Username,
		"bio":      _user.Bio,
		"image":    _user.Image,
	}
}

func (s *Server) createUser() http.HandlerFunc {

	type Input struct {
		User struct {
			Email    string `json:"email" validate:"required,email"`
			Username string `json:"username" validate:"required,min=2"`
			Password string `json:"password" validate:"required,min=8,max=72"`
		} `json:"user" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, http.StatusUnprocessableEntity, err)
			return
		}

		user := conduit.User{
			Email:    input.User.Email,
			Username: input.User.Username,
		}
		user.SetPassword(input.User.Password)

		if err := s.userService.CreateUser(r.Context(), &user); err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		if err := writeJSON(w, http.StatusCreated, M{"user": userResponse(&user)}); err != nil {
			log.Printf("json endcoding error: %v", err)
		}
	}
}

func (s *Server) loginUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8,max=72"`
		} `json:"user" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, http.StatusUnprocessableEntity, err)
			return
		}

		user, err := s.userService.Authenticate(r.Context(), input.User.Email, input.User.Password)

		if err != nil {
			// unauthorized request
			return
		}

		if err := writeJSON(w, http.StatusOK, M{"user": userResponse(user)}); err != nil {
			log.Printf("json endcoding error: %v", err)
		}
	}
}
