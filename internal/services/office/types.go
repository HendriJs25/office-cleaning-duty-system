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
