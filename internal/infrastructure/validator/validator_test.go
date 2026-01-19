package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestAccount struct {
	DocumentNumber string `json:"document_number" validate:"required,len=11,numeric"`
	FullName       string `json:"full_name" validate:"required,min=3,max=100"`
	Email          string `json:"email" validate:"required,email,max=100"`
	Phone          string `json:"phone" validate:"required,min=10,max=15,numeric"`
	AccountType    string `json:"account_type" validate:"required,oneof=checking savings"`
}

type TestTransaction struct {
	AccountID       int64   `json:"account_id" validate:"required"`
	OperationTypeID int     `json:"operation_type_id" validate:"gte=1,lte=4"`
	Amount          float64 `json:"amount" validate:"gt=0"`
}

type TestGtStruct struct {
	Value float64 `json:"value" validate:"gt=0"`
}

type TestGteStruct struct {
	Value int `json:"value" validate:"gte=1"`
}

type TestAlphaStruct struct {
	Name string `json:"name" validate:"required,alpha"`
}

type TestAlphanumStruct struct {
	Code string `json:"code" validate:"required,alphanum"`
}

type TestLtStruct struct {
	Value int `json:"value" validate:"required,lt=10"`
}

type TestLteStruct struct {
	Value int `json:"value" validate:"required,lte=10"`
}

func TestValidateStruct_Success(t *testing.T) {
	valid := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	err := ValidateStruct(&valid)
	assert.NoError(t, err)
}

func TestValidateStruct_RequiredField(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "documentnumber is required")
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "invalid-email",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email must be a valid email address")
}

func TestValidateStruct_MinLength(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "Jo", // Less than 3 characters
		Email:          "joao@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fullname must be at least 3 characters")
}

func TestValidateStruct_MaxLength(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "1198765432112345678", // More than 15 characters
		AccountType:    "checking",
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "phone must be at most 15 characters")
}

func TestValidateStruct_ExactLength(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "1234567890", // 10 characters instead of 11
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "documentnumber must be exactly 11 characters")
}

func TestValidateStruct_Numeric(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "1234567890A", // Contains a letter
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "documentnumber must contain only numbers")
}

func TestValidateStruct_OneOf(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "11987654321",
		AccountType:    "investment", // Not in "checking savings"
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accounttype must be one of [checking savings]")
}

func TestValidateStruct_Gt(t *testing.T) {
	invalid := TestGtStruct{
		Value: 0, // Not greater than 0
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must be greater than 0")
}

func TestValidateStruct_GtNegative(t *testing.T) {
	invalid := TestGtStruct{
		Value: -5.5, // Negative value, not greater than 0
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must be greater than 0")
}

func TestValidateStruct_Gte(t *testing.T) {
	invalid := TestGteStruct{
		Value: 0, // Less than 1
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must be greater than or equal to 1")
}

func TestValidateStruct_GteNegative(t *testing.T) {
	invalid := TestGteStruct{
		Value: -1, // Negative value, less than 1
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must be greater than or equal to 1")
}

func TestValidateStruct_Lte(t *testing.T) {
	invalid := TestTransaction{
		AccountID:       1,
		OperationTypeID: 5, // Greater than 4
		Amount:          100,
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operationtypeid must be less than or equal to 4")
}

func TestValidateStruct_Lt(t *testing.T) {
	invalid := TestLtStruct{
		Value: 10, // Not less than 10
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must be less than 10")
}

func TestValidateStruct_Alpha(t *testing.T) {
	invalid := TestAlphaStruct{
		Name: "João123", // Contains numbers
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name must contain only letters")
}

func TestValidateStruct_Alphanum(t *testing.T) {
	invalid := TestAlphanumStruct{
		Code: "ABC@123", // Contains special character
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code must contain only letters and numbers")
}

func TestValidateStruct_NilInput(t *testing.T) {
	err := ValidateStruct(nil)
	assert.Error(t, err)
}

func TestValidateStructAll_Success(t *testing.T) {
	valid := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	errors := ValidateStructAll(&valid)
	assert.Nil(t, errors)
}

func TestValidateStructAll_MultipleErrors(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "",        // required error
		FullName:       "Jo",      // min length error
		Email:          "invalid", // email error
		Phone:          "11987654321",
		AccountType:    "investment", // oneof error
	}

	errors := ValidateStructAll(&invalid)
	assert.NotNil(t, errors)
	assert.GreaterOrEqual(t, len(errors), 3)
}

func TestValidateStructAll_SingleError(t *testing.T) {
	invalid := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "invalid-email",
		Phone:          "11987654321",
		AccountType:    "checking",
	}

	errors := ValidateStructAll(&invalid)
	assert.NotNil(t, errors)
	assert.Len(t, errors, 1)
	assert.Contains(t, errors[0], "email must be a valid email address")
}

func TestValidateStructAll_NilInput(t *testing.T) {
	errors := ValidateStructAll(nil)
	assert.NotNil(t, errors)
	assert.Contains(t, errors[0], "validation error")
}

func TestFormatValidationError_DefaultCase(t *testing.T) {
	type TestUnknownTag struct {
		Field string `validate:"uuid"` // uuid tag is not in switch case
	}

	invalid := TestUnknownTag{
		Field: "not-a-uuid",
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "field is invalid")
}

func TestGetValidator_ReturnsInstance(t *testing.T) {
	v := GetValidator()
	assert.NotNil(t, v)
}

func TestGetValidator_SameInstance(t *testing.T) {
	v1 := GetValidator()
	v2 := GetValidator()
	assert.Same(t, v1, v2)
}

func TestValidateTransaction_Success(t *testing.T) {
	valid := TestTransaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          100.50,
	}

	err := ValidateStruct(&valid)
	assert.NoError(t, err)
}

func TestValidateTransaction_AllOperationTypes(t *testing.T) {
	operationTypes := []int{1, 2, 3, 4}

	for _, opType := range operationTypes {
		valid := TestTransaction{
			AccountID:       1,
			OperationTypeID: opType,
			Amount:          100.50,
		}

		err := ValidateStruct(&valid)
		assert.NoError(t, err, "Operation type %d should be valid", opType)
	}
}

func TestValidateTransaction_InvalidOperationType(t *testing.T) {
	invalidOps := []int{0, 5, 10, -1}

	for _, opType := range invalidOps {
		invalid := TestTransaction{
			AccountID:       1,
			OperationTypeID: opType,
			Amount:          100.50,
		}

		err := ValidateStruct(&invalid)
		assert.Error(t, err, "Operation type %d should be invalid", opType)
	}
}

func TestValidateTransaction_ZeroAmount(t *testing.T) {
	invalid := TestTransaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          0,
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount must be greater than 0")
}

func TestValidateTransaction_SmallPositiveAmount(t *testing.T) {
	valid := TestTransaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          0.01,
	}

	err := ValidateStruct(&valid)
	assert.NoError(t, err)
}

func TestValidateTransaction_NegativeAmount(t *testing.T) {
	invalid := TestTransaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          -100.50,
	}

	err := ValidateStruct(&invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount must be greater than 0")
}

func TestValidateAccount_AllAccountTypes(t *testing.T) {
	accountTypes := []string{"checking", "savings"}

	for _, accType := range accountTypes {
		valid := TestAccount{
			DocumentNumber: "12345678901",
			FullName:       "João da Silva",
			Email:          "joao@example.com",
			Phone:          "11987654321",
			AccountType:    accType,
		}

		err := ValidateStruct(&valid)
		assert.NoError(t, err, "Account type %s should be valid", accType)
	}
}

func TestValidateAccount_ValidEmails(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.com",
		"user+tag@example.org",
		"user@subdomain.domain.com",
	}

	for _, email := range validEmails {
		valid := TestAccount{
			DocumentNumber: "12345678901",
			FullName:       "João da Silva",
			Email:          email,
			Phone:          "11987654321",
			AccountType:    "checking",
		}

		err := ValidateStruct(&valid)
		assert.NoError(t, err, "Email %s should be valid", email)
	}
}

func TestValidateAccount_InvalidEmails(t *testing.T) {
	invalidEmails := []string{
		"invalid",
		"@domain.com",
		"user@",
		"user@.com",
	}

	for _, email := range invalidEmails {
		invalid := TestAccount{
			DocumentNumber: "12345678901",
			FullName:       "João da Silva",
			Email:          email,
			Phone:          "11987654321",
			AccountType:    "checking",
		}

		err := ValidateStruct(&invalid)
		assert.Error(t, err, "Email %s should be invalid", email)
	}
}

func TestValidateAccount_PhoneBoundaries(t *testing.T) {
	validMin := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "1198765432", // 10 digits
		AccountType:    "checking",
	}
	err := ValidateStruct(&validMin)
	assert.NoError(t, err)

	validMax := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "551198765432100", // 15 digits
		AccountType:    "checking",
	}
	err = ValidateStruct(&validMax)
	assert.NoError(t, err)

	invalidMin := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "João da Silva",
		Email:          "joao@example.com",
		Phone:          "119876543", // 9 digits
		AccountType:    "checking",
	}
	err = ValidateStruct(&invalidMin)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "phone must be at least 10 characters")
}

func TestValidateAccount_FullNameBoundaries(t *testing.T) {
	validMin := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "Ana",
		Email:          "ana@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}
	err := ValidateStruct(&validMin)
	assert.NoError(t, err)

	invalidMin := TestAccount{
		DocumentNumber: "12345678901",
		FullName:       "An",
		Email:          "an@example.com",
		Phone:          "11987654321",
		AccountType:    "checking",
	}
	err = ValidateStruct(&invalidMin)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fullname must be at least 3 characters")
}
