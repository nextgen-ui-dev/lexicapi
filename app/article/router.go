package article

import "github.com/go-chi/chi/v5"

func AdminRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/category", createArticleCategoryHandler)
	r.Get("/category/{id}", getArticleCategoryByIdHandler)
	r.Delete("/category/{id}", deleteArticleCategoryHandler)

	return r
}
