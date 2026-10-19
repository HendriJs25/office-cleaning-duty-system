package role

type RoleResult struct {
	ID       int64
	Code     string
	Name     string
	IsActive bool
}

type GetRoleDetailResult struct {
	ID          int64
	Code        string
	Name        string
	IsActive    bool
	Permissions []PermissionInfo
}

type PermissionInfo struct {
	ID   int64
	Code string
	Name string
}

type GetAllActiveRolesResult struct {
	ID   int64
	Name string
}

type CreateRoleInput struct {
	Code            string
	Name            string
	PermissionsCode []string
}

type UpdateRoleInput struct {
	Name            string
	Description     *string
	PermissionsCode []string
}
