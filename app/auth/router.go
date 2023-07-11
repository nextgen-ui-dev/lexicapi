package auth

import "github.com/go-chi/chi/v5"

func AdminRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sign-in", superadminSignInHandler)

	return r
}

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sign-in/google", signInWithGoogleHandler)

	r.Group(func(r chi.Router) {
		r.Use(UserAuthMiddleware)

		r.Get("/userinfo", getUserInfoHandler)
	})

	return r
}
