package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "OK!"})
}

func main() {
	r := chi.NewRouter()
	r.Get("/", heartbeat)

	http.ListenAndServe(":8080", r)
}
