package response

type GetActiveOfficesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetOfficeResponse struct {
	ID       int64   `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Address  *string `json:"address"`
	IsActive bool    `json:"is_active"`
}
