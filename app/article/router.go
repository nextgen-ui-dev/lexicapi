package article

import "github.com/go-chi/chi/v5"

func AdminRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/category", getArticleCategoriesHandler)
	r.Post("/category", createArticleCategoryHandler)
	r.Get("/category/{id}", getArticleCategoryByIdHandler)
	r.Delete("/category/{id}", deleteArticleCategoryHandler)
	r.Patch("/category/{id}", updateArticleCategoryHandler)

	return r
}
