package auth

type superadminSignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
