package auth

type Superadmin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSuperadmin(email, password string) Superadmin {
	return Superadmin{
		Email:    email,
		Password: password,
	}
}

func (s Superadmin) ValidateCredentials(email, password string) (err error) {
	if err = validateSuperadminEmail(email); err != nil {
		return
	}
	if err = validateSuperadminPassword(password); err != nil {
		return
	}

	return nil
}
