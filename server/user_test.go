package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/0xdod/go-realworld/conduit"
	"github.com/gorilla/mux"
)

type MockUserService struct {
	CreateUserFn  func(*conduit.User) error
	UserByEmailFn func(string) *conduit.User
}

func (m *MockUserService) CreateUser(_ context.Context, user *conduit.User) error {
	return m.CreateUserFn(user)
}

func (m *MockUserService) UserByEmail(_ context.Context, email string) (*conduit.User, error) {
	return m.UserByEmailFn(email), nil
}

func (m *MockUserService) Authenticate(_ context.Context, email, password string) (*conduit.User, error) {
	return nil, nil
}

func Test_createUser(t *testing.T) {
	userStore := &MockUserService{}
	srv := &Server{
		router:      mux.NewRouter(),
		userService: userStore,
	}
	srv.routes()

	input := `{
		"user": {
			"email":    "e@mail.com",
			"username": "username",
			"password": "passwerd"
		}
	}`

	body := strings.NewReader(input)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", body)
	w := httptest.NewRecorder()

	var user conduit.User
	userStore.CreateUserFn = func(u *conduit.User) error {
		user = *u
		return nil
	}

	srv.router.ServeHTTP(w, req)

	if code := w.Code; code != http.StatusCreated {
		t.Errorf("expected status code of 201, but got %d", code)
	}
	expectedResp := userResponse(&user)
	gotResp := M{}
	extractResponseBody(w.Body, &gotResp)

	if !reflect.DeepEqual(expectedResp, gotResp) {
		t.Errorf("expected response %v, but got %v", expectedResp, gotResp)
	}
}

func Test_loginUser(t *testing.T) {

}

func extractResponseBody(body io.Reader, v interface{}) {
	mm := M{}
	_ = readJSON(body, &mm)
	byt, err := json.Marshal(mm["user"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byt, v)
}
