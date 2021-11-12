package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorM map[string][]string

func (e ErrorM) Error() string {
	return ""
}

func validationError(w http.ResponseWriter, _err error) {
	resp := ErrorM{}
	switch err := _err.(type) {
	case validator.ValidationErrors:
		for _, e := range err {
			field := e.Field()
			msg := checkTagRules(e)
			resp[field] = append(resp[field], msg)
		}
	default:
		resp["non_field_error"] = append(resp["non_field_error"], err.Error())
	}
	errorResponse(w, http.StatusUnprocessableEntity, resp)
}

func invalidUserCredentialsError(w http.ResponseWriter) {
	msg := "invalid authentication credentials"
	writeJSON(w, http.StatusUnauthorized, M{"errors": msg})
}

func invalidAuthTokenError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Token")
	msg := "invalid or missing authentication token"
	writeJSON(w, http.StatusUnauthorized, M{"errors": msg})
}

func serverError(w http.ResponseWriter) {
	writeJSON(w, http.StatusInternalServerError, M{"error": "internal error"})
}

func errorResponse(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, M{"errors": err})
}

func checkTagRules(e validator.FieldError) (errMsg string) {
	tag, field, param, value := e.ActualTag(), e.Field(), e.Param(), e.Value()

	if tag == "required" {
		errMsg = "this field is required"
	}

	if tag == "email" {
		errMsg = fmt.Sprintf("%q is not a valid email", value)
	}

	if tag == "min" {
		errMsg = fmt.Sprintf("%s must be greater than %v", field, param)
	}

	if tag == "max" {
		errMsg = fmt.Sprintf("%s must be less than %v", field, param)
	}
	return
}

func logError(err error) {
	log.Println(err)
}
