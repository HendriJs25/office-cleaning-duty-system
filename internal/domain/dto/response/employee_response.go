package response

type GetActiveEmployeesResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}
