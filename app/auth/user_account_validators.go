package auth

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidAccountId              = validation.NewError("auth:invalid_account_id", "Invalid account id")
	ErrInvalidAccountType            = validation.NewError("auth:invalid_account_type", "Invalid account type")
	ErrUnsupportedAccountProvider    = validation.NewError("auth:unsupported_account_provider", "Provider is not supported")
	ErrAccountProviderAccountIdEmpty = validation.NewError("auth:provider_account_id_empty", "Provider account id can't be empty")
)

func validateAccountId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidAccountId
	}

	return id, nil
}

func validateAccountType(typeStr string) (err error) {
	switch typeStr {
	case string(OAUTH), string(CREDENTIALS):
		return nil
	default:
		return ErrInvalidAccountType
	}
}

func validateAccountProvider(provider string) (err error) {
	switch provider {
	case string(GOOGLE):
		return nil
	default:
		return ErrUnsupportedAccountProvider
	}
}

func validateAccountProviderAccountId(providerAccountId string) (err error) {
	providerAccountId = strings.TrimSpace(providerAccountId)
	return validation.Validate(
		&providerAccountId,
		validation.Required.ErrorObject(ErrAccountProviderAccountIdEmpty),
	)
}
