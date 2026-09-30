package request

type CreateOfficeRequest struct {
	Code    string  `json:"code" validate:"required,notblank,lowercase_hyphen,min=3,max=50"`
	Name    string  `json:"name" validate:"required,notblank,max=150"`
	Address *string `json:"address,omitempty" validate:"omitempty,notblank"`
}
