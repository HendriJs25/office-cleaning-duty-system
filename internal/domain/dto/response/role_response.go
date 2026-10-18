package response

type RoleResponse struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type GetActiveRolesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetRoleDetailResponse struct {
	ID          int64                 `json:"id"`
	Code        string                `json:"code"`
	Name        string                `json:"name"`
	IsActive    bool                  `json:"is_active"`
	Permissions []PermissionsResponse `json:"permissions"`
}

type PermissionsResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
