package auth

import (
	"net/http"

	"github.com/lexica-app/lexicapi/app"
)

func signInWithGoogleHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	_ = r.Header.Get("X-Google-Id-Token")

	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Unimplemented"})
}
