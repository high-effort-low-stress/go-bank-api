package http_helpers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/high-effort-low-stress/go-bank-api/internal/utils/http_helpers"
	"github.com/stretchr/testify/assert"
)

type testRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      string `json:"age" binding:"numeric"`
}

func TestValidateJsonRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should return nil when request is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]string{
			"fullName": "John Doe",
			"email":    "john@example.com",
			"age":      "30",
		}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		var req testRequest
		errResponse := http_helpers.ValidateJsonRequest(c, &req)

		assert.Nil(t, errResponse)
		assert.Equal(t, "John Doe", req.FullName)
	})

	t.Run("Should return error when JSON is malformed", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid-json}"))
		c.Request.Header.Set("Content-Type", "application/json")

		var req testRequest
		errResponse := http_helpers.ValidateJsonRequest(c, &req)

		assert.NotNil(t, errResponse)
		assert.Equal(t, "Corpo da requisição inválido", errResponse["error"])
	})

	t.Run("Should return validation details when fields are invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Email inválido e campo 'age' não numérico
		body := map[string]string{
			"fullName": "John Doe",
			"email":    "invalid-email",
			"age":      "abc",
		}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		var req testRequest
		errResponse := http_helpers.ValidateJsonRequest(c, &req)

		assert.NotNil(t, errResponse)
		assert.Equal(t, "Dados inválidos", errResponse["error"])

		details := errResponse["details"].(map[string]string)
		assert.Equal(t, "O formato do e-mail é inválido.", details["email"])
		assert.Equal(t, "O campo 'age' deve conter apenas números.", details["age"])
	})
}

func TestFormatValidationErrors(t *testing.T) {
	validate := validator.New()

	t.Run("Should format different validation tags correctly", func(t *testing.T) {
		type sample struct {
			Name  string `validate:"required"`
			Email string `validate:"email"`
			Code  string `validate:"numeric"`
			Other string `validate:"min=5"`
		}

		s := sample{
			Name:  "",
			Email: "not-an-email",
			Code:  "abc",
			Other: "123",
		}

		err := validate.Struct(s)
		verrs := err.(validator.ValidationErrors)

		formatted := http_helpers.FormatValidationErrors(verrs)

		assert.Equal(t, "O campo 'name' é obrigatório.", formatted["name"])
		assert.Equal(t, "O formato do e-mail é inválido.", formatted["email"])
		assert.Equal(t, "O campo 'code' deve conter apenas números.", formatted["code"])
		// Teste do caso default
		assert.Equal(t, "O campo 'other' é inválido.", formatted["other"])
	})
}
