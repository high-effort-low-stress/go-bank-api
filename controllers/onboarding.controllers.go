package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/high-effort-low-stress/go-bank-api/services"
	"github.com/high-effort-low-stress/go-bank-api/utils"
)

type StartOnboardingRequest struct {
	Document string `json:"document" binding:"required,numeric"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type OnboardingController struct {
	service services.OnboardingService
}

func NewOnboardingController(service services.OnboardingService) *OnboardingController {
	return &OnboardingController{service: service}
}

func (ctrl *OnboardingController) StartOnboarding(c *gin.Context) {
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

	err := ctrl.service.StartOnboardingProcess(req.Document, req.FullName, req.Email)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCPF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "O e-mail de verificação está sendo enviado."})
}
