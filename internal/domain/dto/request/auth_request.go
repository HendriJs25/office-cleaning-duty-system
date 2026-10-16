package request

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyRegistrationTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
type RegisterRequest struct {
	Token           string `json:"token" validate:"required"`
	UserName        string `json:"user_name" validate:"required,notblank,max=200"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,validpassword"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}
