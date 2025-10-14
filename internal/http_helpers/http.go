package http_helpers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateJsonRequest(c *gin.Context, req any) gin.H {
	err := c.ShouldBindJSON(&req)
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return gin.H{
			"error":   "Dados inválidos",
			"details": FormatValidationErrors(validationErrors),
		}
	}

	return gin.H{"error": "Corpo da requisição inválido"}

}

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
