package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/high-effort-low-stress/go-bank-api/utils"
)

type StartOnboardingRequest struct {
	Document string `json:"document" binding:"required,numeric"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// StartOnboarding é o handler para o endpoint POST /onboarding/start.
// Ele valida os dados de entrada e simula o início do processo de onboarding.
func StartOnboarding(c *gin.Context) {
	var req StartOnboardingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Dados inválidos",
				"details": utils.FormatValidationErrors(validationErrors),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido"})
		return
	}
	if !utils.IsValidCPF(req.Document) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CPF inválido"})
		return
	}

	if err := verifyConflict(req.Document, req.Email); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// 2. Criar o usuário no banco com status "pendente" e gerar um token de verificação.
	// 3. Disparar o e-mail de verificação (usando um serviço de e-mail).
	log.Printf("Iniciando onboarding para o e-mail: %s", req.Email)

	c.JSON(http.StatusAccepted, gin.H{"message": "O e-mail de verificação está sendo enviado."})
}

func verifyConflict(document string, email string) error {
	if document == "11122233344" {
		return fmt.Errorf("O CPF ou E-mail já está cadastrado.")
	}

	if email == "cadastrado@gmail.com" {
		return fmt.Errorf("O CPF ou E-mail já está cadastrado.")
	}

	return nil
}
