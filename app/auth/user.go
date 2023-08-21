package auth

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrUnverifiedUser = errors.New("User must be verified first")
)

type UserStatus uint
type UserRole null.String
type UserEducationLevel null.String

const (
	NOT_VERIFIED UserStatus = iota
	NOT_ONBOARDED
	ACTIVE
)

var (
	STUDENT  = UserRole(null.StringFrom("pelajar"))
	EDUCATOR = UserRole(null.StringFrom("pengajar"))
	CIVILIAN = UserRole(null.StringFrom("umum"))
)

var (
	SMP     UserEducationLevel = UserEducationLevel(null.StringFrom("smp"))
	SMA     UserEducationLevel = UserEducationLevel(null.StringFrom("sma"))
	SARJANA UserEducationLevel = UserEducationLevel(null.StringFrom("sarjana"))
	LAINNYA UserEducationLevel = UserEducationLevel(null.StringFrom("lainnya"))
)

type User struct {
	Id             ulid.ULID          `json:"id"`
	Name           null.String        `json:"name"`
	Email          null.String        `json:"email"`
	ImageUrl       null.String        `json:"image_url"`
	Status         UserStatus         `json:"status"`
	Role           UserRole           `json:"role"`
	EducationLevel UserEducationLevel `json:"education_level"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      null.Time          `json:"updated_at"`
	DeletedAt      null.Time          `json:"deleted_at"`
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

func (u *User) Onboard(roleStr, educationLevelStr string) (errs map[string]error) {
	errs = make(map[string]error)

	if u.Status == ACTIVE {
		return nil
	}

	if u.Status == NOT_VERIFIED {
		errs["status"] = ErrUnverifiedUser
		return errs
	}

	role, err := validateUserRole(null.StringFrom(roleStr))
	if err != nil {
		errs["role"] = err
	}

	educationLevel, err := validateUserEducationLevel(null.StringFrom(educationLevelStr))
	if err != nil {
		errs["education_level"] = err
	}

	if len(errs) != 0 {
		return errs
	}

	u.Role = role
	u.EducationLevel = educationLevel
	u.Status = ACTIVE
	u.UpdatedAt = null.TimeFrom(time.Now())

	return nil
}
