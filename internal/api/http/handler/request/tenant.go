package request

type CreateTenantRequest struct {
	Name string `json:"name" validate:"required"`
}

type ProcessPayloadRequest struct {
	Payload interface{} `json:"payload" validate:"required"`
}
