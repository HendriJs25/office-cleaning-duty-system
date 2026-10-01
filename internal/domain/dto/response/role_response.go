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
