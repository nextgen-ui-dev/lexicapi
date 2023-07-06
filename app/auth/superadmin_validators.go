package auth

import (
	"strings"

	"github.com/jellydator/validation"
)

var ErrIncorrectCredentials = validation.NewError("auth:incorrect_credentials", "Incorrect email or password")

func validateSuperadminEmail(email string) (err error) {
	email = strings.TrimSpace(email)
	return validation.Validate(
		&email,
		validation.NewStringRuleWithError(func(str string) bool {
			return str == superadmin.Email
		}, ErrIncorrectCredentials),
	)
}

func validateSuperadminPassword(password string) (err error) {
	password = strings.TrimSpace(password)
	return validation.Validate(
		&password,
		validation.NewStringRuleWithError(func(str string) bool {
			return str == superadmin.Password
		}, ErrIncorrectCredentials),
	)
}
