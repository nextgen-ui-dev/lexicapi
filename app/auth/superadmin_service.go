package auth

import "context"

func superadminSignIn(ctx context.Context, body superadminSignInReq) (tokens SuperadminTokens, err error) {
	if err = superadmin.ValidateCredentials(body.Email, body.Password); err != nil {
		return
	}

	return
}
