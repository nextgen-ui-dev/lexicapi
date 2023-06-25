package app

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrNotFound         = errors.New("Page not found")
	ErrMethodNotAllowed = errors.New("Method not allowed")
)

func WriteHttpBodyJson(w http.ResponseWriter, status int, body any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		WriteHttpError(w, http.StatusInternalServerError, err)
	}
}

func WriteHttpError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "OK!"})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteHttpError(w, http.StatusNotFound, ErrNotFound)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	WriteHttpError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
}
