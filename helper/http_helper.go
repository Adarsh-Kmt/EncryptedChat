package helper

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Status int
	Error  any
}

type HTTPFunc func(w http.ResponseWriter, r *http.Request) *HTTPError

func MakeHTTPHandlerFunc(f HTTPFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if httpError := f(w, r); httpError != nil {

			WriteJSON(w, map[string]any{"error": httpError.Error}, httpError.Status)
		}
	}
}

func WriteJSON(w http.ResponseWriter, body any, status int) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
