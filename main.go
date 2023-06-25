package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "OK!"})
}

func main() {
	config, err := app.SetConfig()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Get("/", heartbeat)

	log.Println("Running server on port", config.Port, "in", config.Env, "mode...")
	http.ListenAndServe(":"+config.Port, r)
}
