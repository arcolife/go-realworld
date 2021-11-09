package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// M is a generic map
type M map[string]interface{}

// R represents a standard JSON API response format containing the status, message and data.
type R struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RM is an alternate form of R that uses a generic map to envelope the actual payload.
type RM struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    M      `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(body io.Reader, input interface{}) error {
	return json.NewDecoder(body).Decode(input)
}

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e Error) ErrorMessage() string {
	return e.Message
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%s message=%s", e.Code, e.Message)
}

func errorResponse(w http.ResponseWriter, code int, errs ...error) {
	_errors := []string{}

	for _, err := range errs {
		switch e := err.(type) {
		case Error:
			_errors = append(_errors, e.Message)
		default:
		}
	}

	writeJSON(w, code, M{"errors": _errors})

}
