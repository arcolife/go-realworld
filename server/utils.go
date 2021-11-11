package server

import (
	"encoding/json"
	"io"
	"net/http"
)

// M is a generic map
type M map[string]interface{}

func writeJSON(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(body io.Reader, input interface{}) error {
	return json.NewDecoder(body).Decode(input)
}
