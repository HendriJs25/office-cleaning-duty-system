package request

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyRegistrationTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
