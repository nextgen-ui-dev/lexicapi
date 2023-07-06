package main

import (
	stdlog "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lexica-app/lexicapi/adapters"
	"github.com/lexica-app/lexicapi/app"
	"github.com/lexica-app/lexicapi/app/article"
	"github.com/lexica-app/lexicapi/app/auth"
	"github.com/lexica-app/lexicapi/db"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		stdlog.Fatal("Failed to load config:", err)
	}

	// App Configurations
	app.ConfigureLogger(config)
	app.ConfigureCors(config)

	// Configure Adapters and Dependency Injection
	pool := db.CreateConnPool(config.DbDsn)
	openaiAdapter := adapters.ConfigureOpenAIAdapter(config.OpenAIOrganizationId, config.OpenAIAPIKey)

	article.SetPool(pool)
	article.SetOpenAIAdapter(openaiAdapter)

	auth.ConfigureLexicaAPIKey(config.LexicaApiKey)
	auth.ConfigureSuperadminJWTSecret(config.LexicaSuperadminJwtSecret)
	auth.ConfigureSuperadmin(config.LexicaSuperadminEmail, config.LexicaSuperadminPassword)

	r := chi.NewRouter()

	// Global middlewares
	r.Use(app.ReqLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(app.CorsMiddleware)
	r.Use(auth.LexicaAPIKeyMiddleware)

	// Default route handlers
	r.NotFound(app.NotFound)
	r.MethodNotAllowed(app.MethodNotAllowed)
	r.Get("/", app.Heartbeat)

	// Admin Routes
	r.Group(func(r chi.Router) {
		r.Mount("/admin/auth", auth.AdminRouter())
		r.Mount("/admin/article", article.AdminRouter())
	})

	// Normal Routes
	r.Group(func(r chi.Router) {
		r.Mount("/article", article.Router())
	})

	log.Info().Msgf("Running server on port %s in %s mode...", config.Port, config.Env)
	http.ListenAndServe(":"+config.Port, r)
}
