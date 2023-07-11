package auth

import "context"

func superadminSignIn(ctx context.Context, body superadminSignInReq) (tokens SuperadminTokens, err error) {
	if err = superadmin.ValidateCredentials(body.Email, body.Password); err != nil {
		return
	}

	accessToken, err := generateSuperadminAccessToken()
	if err != nil {
		return
	}

	return SuperadminTokens{
		AccessToken: accessToken,
	}, nil
}
