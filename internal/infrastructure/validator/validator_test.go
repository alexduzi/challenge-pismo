package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required,min=3,max=10"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,gte=18,lte=100"`
	CPF   string `validate:"required,cpf"`
}

func TestValidateStruct_Success(t *testing.T) {
	valid := TestStruct{
		Name:  "João",
		Email: "joao@example.com",
		Age:   25,
		CPF:   "12345678909",
	}

	err := ValidateStruct(&valid)
	assert.NoError(t, err)
}

func TestValidateStruct_RequiredField(t *testing.T) {
	invalid := TestStruct{
		Email: "joao@example.com",
		Age:   25,
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	invalid := TestStruct{
		Name:  "João",
		Email: "invalid-email",
		Age:   25,
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email must be a valid email address")
}

func TestValidateCPF_Valid(t *testing.T) {
	validCPFs := []string{
		"12345678909",
		"11144477735",
	}

	for _, cpf := range validCPFs {
		valid := TestStruct{
			Name:  "Test",
			Email: "test@test.com",
			Age:   25,
			CPF:   cpf,
		}
		err := ValidateStruct(&valid)
		assert.NoError(t, err, "CPF %s should be valid", cpf)
	}
}

func TestValidateCPF_Invalid(t *testing.T) {
	invalidCPFs := []string{
		"00000000000",
		"11111111111",
		"12345678900",
		"123",
	}

	for _, cpf := range invalidCPFs {
		invalid := TestStruct{
			Name:  "Test",
			Email: "test@test.com",
			Age:   25,
			CPF:   cpf,
		}
		err := ValidateStruct(&invalid)
		assert.Error(t, err, "CPF %s should be invalid", cpf)
	}
}
