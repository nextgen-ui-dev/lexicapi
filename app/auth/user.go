package auth

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type UserStatus uint

const (
	NOT_VERIFIED UserStatus = iota
	NOT_ONBOARDED
	ACTIVE
)

type User struct {
	Id        ulid.ULID   `json:"id"`
	Name      null.String `json:"name"`
	Email     null.String `json:"email"`
	ImageUrl  null.String `json:"image_url"`
	Status    UserStatus  `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}

func NewUserWithOAuth(name, email, imageUrl null.String) (User, map[string]error) {
	errs := make(map[string]error)

	if err := validateUserName(name.String); err != nil {
		errs["name"] = err
	}
	if err := validateUserEmail(email.String); err != nil {
		errs["email"] = err
	}
	if err := validateUserImageUrl(imageUrl.String); err != nil {
		errs["image_url"] = err
	}
	if len(errs) != 0 {
		return User{}, errs
	}

	id := ulid.Make()

	return User{
		Id:        id,
		Name:      name,
		Email:     email,
		ImageUrl:  imageUrl,
		Status:    NOT_ONBOARDED,
		CreatedAt: time.Now(),
	}, nil
}
