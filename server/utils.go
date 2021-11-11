package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// M is a generic map
type M map[string]interface{}

type ErrorM map[string][]string

func (e ErrorM) Error() string {
	return ""
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(body io.Reader, input interface{}) error {
	return json.NewDecoder(body).Decode(input)
}

func validationError(w http.ResponseWriter, code int, _err error) {
	resp := ErrorM{}
	switch err := _err.(type) {
	case validator.ValidationErrors:
		for _, e := range err {
			field := e.Field()
			msg := checkTagRules(e.ActualTag(), field, e.Value())
			resp[field] = append(resp[field], msg)
		}
	default:
		resp["non_field_error"] = append(resp["non_field_error"], err.Error())
	}
	errorResponse(w, code, resp)
}

func badRequestError(w http.ResponseWriter, code int, _err error) {
	resp := ErrorM{
		"request_error": []string{_err.Error()},
	}
	errorResponse(w, code, resp)
}

func serverError(w http.ResponseWriter, code int, _err error) {
	resp := ErrorM{
		"server_error": []string{_err.Error()},
	}
	errorResponse(w, code, resp)
}

func errorResponse(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, M{"errors": err})
}

func checkTagRules(tag, field, value interface{}) (errMsg string) {
	if tag == "required" {
		errMsg = "this field is required"
	}
	if tag == "email" {
		errMsg = fmt.Sprintf("%q is not a valid email", value)
	}
	return
}
