package auth

import (
	"context"
	"fmt"
)

func signInWithGoogle(ctx context.Context, idToken string) (signIn UserSignIn, errs map[string]error, err error) {
	payload, err := validateUserGoogleIdToken(ctx, idToken)
	if err != nil {
		return
	}

	_, name, email, imageUrl := extractProfileFromGoogleIdTokenPayload(payload)
	user, errs := NewUserWithOAuth(name, email, imageUrl)
	if errs != nil {
		return
	}

	fmt.Println(user)

	return signIn, nil, nil
}
