package response

type ErrorResponse struct {
	Message string `json:"message" example:"validation error: invalid input"`
}
