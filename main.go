package main

import (
	"encoding/json"
	stdlog "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lexica-app/lexicapi/app"
	"github.com/rs/zerolog/log"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "OK!"})
}

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		stdlog.Fatal("Failed to load config:", err)
	}

	app.ConfigureLogger(config)

	r := chi.NewRouter()
	r.Use(app.ReqLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Get("/", heartbeat)

	log.Info().Msgf("Running server on port %s in %s mode...", config.Port, config.Env)
	http.ListenAndServe(":"+config.Port, r)
}
