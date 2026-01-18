package request

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required,len=11,numeric" example:"12345678900"`
	FullName       string `json:"full_name" validate:"required,min=3,max=100" example:"João da Silva"`
	Email          string `json:"email" validate:"required,email,max=100" example:"joao.silva@example.com"`
	Phone          string `json:"phone" validate:"required,min=10,max=15,numeric" example:"11987654321"`
	AccountType    string `json:"account_type" validate:"required,oneof=checking savings" example:"checking"`
}
