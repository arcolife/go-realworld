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

// userResponse is a helper function used to return the User response in the format specified
// by the API spec.
func userResponse(user *conduit.User, _token ...string) M {
	if user == nil {
		return nil
	}
	var token string
	if len(_token) > 0 {
		token = _token[0]
	}
	return M{
		"email":    user.Email,
		"token":    token,
		"username": user.Username,
		"bio":      user.Bio,
		"image":    user.Image,
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
				serverError(w, err)
			}
			return
		}

		writeJSON(w, http.StatusCreated, M{"user": user})
	}
}

func (s *Server) loginUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		user, err := s.userService.Authenticate(r.Context(), input.User.Email, input.User.Password)

		if err != nil || user == nil {
			invalidUserCredentialsError(w)
			return
		}

		token, err := generateUserToken(user)

		if err != nil {
			serverError(w, err)
			return
		}

		user.Token = token

		writeJSON(w, http.StatusOK, M{"user": user})

	}
}

func (s *Server) getCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := userFromContext(ctx)
		user.Token = userTokenFromContext(ctx)

		writeJSON(w, http.StatusOK, M{"user": user})
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
			badRequestError(w)
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
			serverError(w, err)
			return
		}

		user.Token = userTokenFromContext(ctx)

		writeJSON(w, http.StatusOK, M{"user": user})
	}
}
