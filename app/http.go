package app

import (
	"encoding/json"
	"net/http"
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
