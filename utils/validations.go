package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/uguremirmustafa/go_commerce/models"
)

var validate = validator.New()

func ValidateUserStruct(model models.User) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	err := validate.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateLoginRequestStruct(model models.LoginRequest) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	err := validate.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateSignupRequestStruct(model models.SignupRequest) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	err := validate.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
