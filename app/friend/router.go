package friend

import (
	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app/auth"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(auth.UserAuthMiddleware)

	r.Post("/{requesteeId}/add", sendFriendRequestHandler)
	r.Patch("/{friendId}/accept", acceptFriendRequestHandler)
	r.Delete("/{friendId}/reject", rejectFriendRequestHandler)

	return r
}
