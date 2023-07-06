package article

import (
	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app/auth"
)

func AdminRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(auth.SuperadminAuthMiddleware)

	r.Get("/category", getArticleCategoriesHandler)
	r.Post("/category", createArticleCategoryHandler)
	r.Get("/category/{id}", getArticleCategoryByIdHandler)
	r.Delete("/category/{id}", deleteArticleCategoryHandler)
	r.Patch("/category/{id}", updateArticleCategoryHandler)

	r.Get("/", getArticlesHandler)
	r.Post("/", createArticleHandler)
	r.Get("/{id}", getArticleByIdHandler)
	r.Put("/{id}", updateArticleHandler)
	r.Delete("/{id}", removeArticleHandler)

	r.Post("/{articleId}/text", createArticleTextHandler)
	r.Patch("/{articleId}/text/{id}", updateArticleTextHandler)
	r.Delete("/{articleId}/text/{id}", removeArticleTextHandler)
	r.Post("/{articleId}/text/generate", generateOpenAIArticleTextHandler)
	r.Patch("/{articleId}/text/{id}/regenerate", regenerateOpenAIArticleTextHandler)

	return r
}

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/category", getArticleCategoriesHandler)
	r.Get("/category/{id}", getArticleCategoryByIdHandler)

	r.Get("/", getArticlesHandler)
	r.Get("/{id}", getArticleByIdHandler)

	return r
}
