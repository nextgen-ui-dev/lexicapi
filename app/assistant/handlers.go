package assistant

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lexica-app/lexicapi/app"
)

func explainTextHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body ExplainTextReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	explained, err := explainText(ctx, body)
	if err != nil {
		switch {
		case
			errors.As(err, &ErrExplainedTextEmpty),
			errors.As(err, &ErrExplainedTextTooLong),
			errors.As(err, &ErrExplainedExplanationEmpty):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrInvalidOpenAIAPIKey):
			app.WriteHttpError(w, http.StatusUnauthorized, err)
		case errors.Is(err, ErrOpenAIRateLimited):
			app.WriteHttpError(w, http.StatusTooManyRequests, err)
		case errors.Is(err, ErrOpenAIServiceError):
			app.WriteHttpError(w, http.StatusServiceUnavailable, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, explained)
}
