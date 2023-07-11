package auth

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type AccountType string

const (
	OAUTH       AccountType = "oauth"
	CREDENTIALS AccountType = "credentials"
)

type AccountProvider string

const (
	GOOGLE AccountProvider = "google"
)

type Account struct {
	Id                ulid.ULID       `json:"id"`
	UserId            ulid.ULID       `json:"user_id"`
	Type              AccountType     `json:"type"`
	PasswordHash      null.String     `json:"password_hash"`
	Provider          AccountProvider `json:"provider"`
	ProviderAccountId string          `json:"provider_account_id"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         null.Time       `json:"updated_at"`
	DeletedAt         null.Time       `json:"deleted_at"`
}

func NewAccount(
	userIdStr,
	typeStr string,
	passwordHash null.String,
	providerStr string,
	providerAccountId string,
) (Account, map[string]error) {
	errs := make(map[string]error)

	userId, err := validateUserId(userIdStr)
	if err != nil {
		errs["user_id"] = err
	}
	if err = validateAccountType(typeStr); err != nil {
		errs["type"] = err
	}
	if err = validateAccountProvider(providerStr); err != nil {
		errs["provider"] = err
	}
	if err = validateAccountProviderAccountId(providerAccountId); err != nil {
		errs["provider_account_id"] = err
	}
	if len(errs) != 0 {
		return Account{}, errs
	}

	id := ulid.Make()

	return Account{
		Id:                id,
		UserId:            userId,
		Type:              AccountType(typeStr),
		PasswordHash:      passwordHash,
		Provider:          AccountProvider(providerStr),
		ProviderAccountId: providerAccountId,
		CreatedAt:         time.Now(),
	}, nil
}
