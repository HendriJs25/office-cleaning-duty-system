package role

type RoleResult struct {
	ID       int64
	Code     string
	Name     string
	IsActive bool
}

type GetAllActiveRolesResult struct {
	ID   int64
	Name string
}
