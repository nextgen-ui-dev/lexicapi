package assistant

import (
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	// TODO: uncomment if client has supported user auth
	// r.Use(auth.UserAuthMiddleware)

	r.Post("/simplify", simplifyTextHandler)
	r.Post("/explain", explainTextHandler)

	return r
}
