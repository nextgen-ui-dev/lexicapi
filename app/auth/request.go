package auth

type superadminSignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type onboardReq struct {
	Role           string   `json:"role"`
	EducationLevel string   `json:"education_level"`
	InterestIds    []string `json:"interest_ids"`
}
