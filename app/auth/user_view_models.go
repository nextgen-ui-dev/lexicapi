package auth

import "github.com/oklog/ulid/v2"

type UserSignIn struct {
	UserId       ulid.ULID `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
