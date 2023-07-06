package auth

import (
	"net/http"

	"github.com/lexica-app/lexicapi/app"
)

func superadminSignInHandler(w http.ResponseWriter, r *http.Request) {
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Unimplemented"})
}
