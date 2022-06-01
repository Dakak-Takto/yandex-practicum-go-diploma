package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, i any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(i)
}

func JSONmsg(w http.ResponseWriter, status int, key, value string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		key: value,
	})
}

func decodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
