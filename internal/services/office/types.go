package office

type CreateOfficeInput struct {
	Code    string
	Name    string
	Address *string
}

type GetActiveOfficeResult struct {
	ID   int64
	Name string
}

type GetOfficeResult struct {
	ID       int64
	Code     string
	Name     string
	Address  *string
	IsActive bool
}

type UpdateOfficeInput struct {
	Name    string
	Address *string
}
