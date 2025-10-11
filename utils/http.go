package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)

	for _, err := range errs {
		jsonField := strings.ToLower(string(err.Field()[0])) + err.Field()[1:]

		switch err.Tag() {
		case "required":
			errorMessages[jsonField] = fmt.Sprintf("O campo '%s' é obrigatório.", jsonField)
		case "email":
			errorMessages[jsonField] = "O formato do e-mail é inválido."
		case "numeric":
			errorMessages[jsonField] = fmt.Sprintf("O campo '%s' deve conter apenas números.", jsonField)
		default:
			errorMessages[jsonField] = fmt.Sprintf("O campo '%s' é inválido.", jsonField)
		}
	}

	return errorMessages
}
