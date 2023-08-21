package auth

import "github.com/oklog/ulid/v2"

type superadminSignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type onboardReq struct {
	Role           string `json:"role"`
	EducationLevel string `json:"education_level"`
	InterestIds    []ulid.ULID `json:"interest_ids"`
}
