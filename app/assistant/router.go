package assistant

import "github.com/go-chi/chi/v5"

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/explain", explainTextHandler)

	return r
}
