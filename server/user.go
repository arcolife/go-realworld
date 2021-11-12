package server

import (
	"errors"
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

func userResponse(_user *conduit.User, _token ...string) M {
	if _user == nil {
		return nil
	}
	var token string
	if len(_token) > 0 {
		token = _token[0]
	}
	return M{
		"email":    _user.Email,
		"token":    token,
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
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		user := conduit.User{
			Email:    input.User.Email,
			Username: input.User.Username,
		}
		user.SetPassword(input.User.Password)

		if err := s.userService.CreateUser(r.Context(), &user); err != nil {
			switch {
			case errors.Is(err, conduit.ErrDuplicateEmail):
				err = ErrorM{"email": []string{"this email is already in use"}}
				errorResponse(w, http.StatusConflict, err)
			case errors.Is(err, conduit.ErrDuplicateUsername):
				err = ErrorM{"username": []string{"this username is already in use"}}
				errorResponse(w, http.StatusConflict, err)
			default:
				serverError(w)
			}
			logError(err)
			return
		}

		if err := writeJSON(w, http.StatusCreated, M{"user": userResponse(&user)}); err != nil {
			logError(err)
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
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		user, err := s.userService.Authenticate(r.Context(), input.User.Email, input.User.Password)

		if err != nil || user == nil {
			invalidUserCredentialsError(w)
			return
		}

		token, err := generateUserToken(user)

		if err != nil {
			logError(err)
			serverError(w)
			return
		}

		if err := writeJSON(w, http.StatusOK, M{"user": userResponse(user, token)}); err != nil {
			logError(err)
		}
	}
}

func (s *Server) getCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := userFromContext(ctx)
		token := userTokenFromContext(ctx)
		err := writeJSON(w, http.StatusOK, M{"user": userResponse(user, token)})

		if err != nil {
			logError(err)
		}

	}
}

func (s *Server) updateUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Email    *string `json:"email,omitempty"`
			Username *string `json:"username,omitempty"`
			Bio      *string `json:"bio,omitempty"`
			Image    *string `json:"image,omitempty"`
			Password *string `json:"password,omitempty"`
		} `json:"user,omitempty" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		ctx := r.Context()
		user := userFromContext(ctx)
		patch := conduit.UserPatch{
			Username: input.User.Username,
			Bio:      input.User.Bio,
			Email:    input.User.Email,
			Image:    input.User.Image,
		}

		if v := input.User.Password; v != nil {
			user.SetPassword(*v)
		}

		err := s.userService.UpdateUser(ctx, user, patch)
		if err != nil {
			logError(err)
			serverError(w)
			return
		}

		token := userTokenFromContext(ctx)

		if err := writeJSON(w, http.StatusOK, M{"user": userResponse(user, token)}); err != nil {
			logError(err)
		}
	}
}
