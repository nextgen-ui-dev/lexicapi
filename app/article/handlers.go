package article

import (
	"encoding/json"
	"net/http"

	"github.com/lexica-app/lexicapi/app"
)

func createArticleCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body createArticleCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	category, err := createArticleCategory(ctx, body.Name)
	if err != nil {
		switch err {
		case ErrArticleCategoryNameExists, ErrArticleCategoryNameTooLong, ErrArticleCategoryNameEmpty:
			app.WriteHttpError(w, http.StatusBadRequest, err)
		default:
			app.WriteHttpError(w, http.StatusInternalServerError, err)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, category)
}
