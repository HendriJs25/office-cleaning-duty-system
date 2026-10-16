package request

type CreateRoleRequest struct {
	Code            string   `json:"code" validate:"required,min=3,max=50,lowercase_hyphen"`
	Name            string   `json:"name" validate:"required,notblank,max=100"`
	PermissionCodes []string `json:"permission_codes" validate:"required,min=1,unique,dive,required"`
}
