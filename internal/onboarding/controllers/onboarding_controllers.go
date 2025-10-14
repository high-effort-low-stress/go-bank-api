package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/high-effort-low-stress/go-bank-api/internal/http_helpers"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/services"
)

type StartOnboardingRequest struct {
	Document string `json:"document" binding:"required,numeric"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type OnboardingController struct {
	createOnboardingService services.OnboardingService
	verifyEmailTokenService services.VerifyEmailTokenService
}

func NewOnboardingController(
	createOnboardingService services.OnboardingService,
	verifyEmailTokenService services.VerifyEmailTokenService,
) *OnboardingController {
	return &OnboardingController{
		createOnboardingService: createOnboardingService,
		verifyEmailTokenService: verifyEmailTokenService,
	}
}

func (ctrl *OnboardingController) StartOnboarding(c *gin.Context) {
	var req StartOnboardingRequest

	if response := http_helpers.ValidateJsonRequest(c, &req); response != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := ctrl.createOnboardingService.StartOnboardingProcess(req.Document, req.FullName, req.Email)
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

func (ctrl *OnboardingController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
	}

	err := ctrl.verifyEmailTokenService.Execute(token)

	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrAlreadyVerified) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verificado com sucesso!"})

}
