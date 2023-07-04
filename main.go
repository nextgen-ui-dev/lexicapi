package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lexica-app/lexicapi/adapters"
	"github.com/lexica-app/lexicapi/app"
	"github.com/lexica-app/lexicapi/app/article"
	"github.com/lexica-app/lexicapi/db"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		stdlog.Fatal("Failed to load config:", err)
	}

	app.ConfigureLogger(config)

	// Database connection injection
	pool := db.CreateConnPool(config.DbDsn)
	article.SetPool(pool)

	// Configure Adapters and Injection
	openaiAdapter := adapters.ConfigureOpenAIAdapter(config.OpenAIOrganizationId, config.OpenAIOrganizationId)
	article.SetOpenAIAdapter(openaiAdapter)

	// Router mounts
	r := chi.NewRouter()
	r.Use(app.ReqLoggerMiddleware)
	r.Use(middleware.Recoverer)

	r.NotFound(app.NotFound)
	r.MethodNotAllowed(app.MethodNotAllowed)

	r.Get("/", app.Heartbeat)

	// Admin Routes
	r.Group(func(r chi.Router) {
		r.Mount("/admin/article", article.AdminRouter())
	})

	log.Info().Msgf("Running server on port %s in %s mode...", config.Port, config.Env)
	http.ListenAndServe(":"+config.Port, r)
}
