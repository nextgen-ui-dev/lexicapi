package auth

import "github.com/go-chi/chi/v5"

func AdminRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sign-in", superadminSignInHandler)

	return r
}
