package request

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequest struct {
	RoleID          int64  `json:"role_id" validate:"required"`
	EmployeeID      *int64 `json:"employee_id,omitempty"`
	UserName        string `json:"user_name" validate:"required,notblank,max=200"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,validpassword"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}
