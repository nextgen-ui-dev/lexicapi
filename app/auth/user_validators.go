package auth

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/jellydator/validation/is"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidUserId       = validation.NewError("auth:invalid_user_id", "Invalid user id")
	ErrUserNameTooLong     = validation.NewError("auth:user_name_too_long", "Name can't be longer than 255 characters")
	ErrInvalidUserEmail    = validation.NewError("auth:invalid_user_email", "Invalid user email")
	ErrInvalidUserImageUrl = validation.NewError("auth:invalid_user_image_url", "Invalid user image url")
	ErrInvalidUserStatus   = validation.NewError("auth:invalid_user_status", "Invalid user status")
)

func validateUserId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidUserId
	}

	return id, nil
}

func validateUserName(name string) (err error) {
	name = strings.TrimSpace(name)
	return validation.Validate(
		&name,
		validation.When(
			!validation.IsEmpty(name),
			validation.Length(1, 255).ErrorObject(ErrUserNameTooLong),
		),
	)
}

func validateUserEmail(email string) (err error) {
	email = strings.TrimSpace(email)
	return validation.Validate(
		&email,
		validation.When(
			!validation.IsEmpty(email),
			is.Email.ErrorObject(ErrInvalidUserEmail),
		),
	)
}

func validateUserImageUrl(imageUrl string) (err error) {
	imageUrl = strings.TrimSpace(imageUrl)
	return validation.Validate(
		&imageUrl,
		validation.When(
			!validation.IsEmpty(imageUrl),
			is.URL.ErrorObject(ErrInvalidUserImageUrl),
		),
	)
}

func validateUserStatus(status int) (err error) {
	if status < int(NOT_VERIFIED) || status > int(ACTIVE) {
		return ErrInvalidUserStatus
	}

	return nil
}
