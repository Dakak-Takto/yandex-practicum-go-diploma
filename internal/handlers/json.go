package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// write json with Content-type header
func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// write json message ({key:value}) with content-type header
func JSONmsg(w http.ResponseWriter, status int, key, value string) {
	JSON(w, status, map[string]string{
		key: value,
	})
}

func decodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
