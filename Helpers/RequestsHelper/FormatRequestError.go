package requestsHelper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func FormatValidationErrors(err error) []ValidationResponse {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {

		formattedErrors := make([]ValidationResponse, 0, len(validationErrors))

		for _, fieldError := range validationErrors {
			formattedErrors = append(formattedErrors, ValidationResponse{
				Field: fieldError.Field(),
				Tag:   fieldError.Tag(),

				Value: fmt.Sprintf("%v", fieldError.Value()),
			})
		}
		return formattedErrors
	}

	return nil
}
