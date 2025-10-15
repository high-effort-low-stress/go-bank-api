package controllers

import (
	"errors"
	"log"
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

type CompleteOnboardingRequest struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,min=8"`
}

type OnboardingController struct {
	createOnboardingService   services.OnboardingService
	verifyEmailTokenService   services.VerifyEmailTokenService
	completeOnboardingService services.CompleteOnboardingService
}

func NewOnboardingController(
	createOnboardingService services.OnboardingService,
	verifyEmailTokenService services.VerifyEmailTokenService,
	completeOnboardingService services.CompleteOnboardingService,
) *OnboardingController {
	return &OnboardingController{
		createOnboardingService:   createOnboardingService,
		verifyEmailTokenService:   verifyEmailTokenService,
		completeOnboardingService: completeOnboardingService,
	}
}

func (ctrl *OnboardingController) StartOnboarding(c *gin.Context) {
	var req StartOnboardingRequest

	if response := http_helpers.ValidateJsonRequest(c, &req); response != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := ctrl.createOnboardingService.StartOnboardingProcess(req.Document, req.FullName, req.Email)
	if err == nil {
		c.JSON(http.StatusAccepted, gin.H{"message": "O e-mail de verificação está sendo enviado."})
		return
	}

	if errors.Is(err, services.ErrInvalidCPF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, services.ErrUserExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
}

func (ctrl *OnboardingController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": services.ErrInvalidToken.Error()})
		return
	}

	err := ctrl.verifyEmailTokenService.Execute(token)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Email verificado com sucesso!"})
		return
	}

	if errors.Is(err, services.ErrInvalidToken) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, services.ErrAlreadyVerified) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, services.ErrExpiredToken) {
		c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
}

func (ctrl *OnboardingController) CompleteOnboarding(c *gin.Context) {
	var req CompleteOnboardingRequest

	if response := http_helpers.ValidateJsonRequest(c, &req); response != nil {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := ctrl.completeOnboardingService.Execute(req.Token, req.Password, req.PasswordConfirmation)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Cadastro concluído com sucesso."})
		return
	}

	if errors.Is(err, services.ErrInvalidToken) || errors.Is(err, services.ErrRequestNotVerified) || errors.Is(err, services.ErrPasswordsDoNotMatch) || errors.Is(err, services.ErrWeakPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, services.ErrAlreadyVerified) {
		c.JSON(http.StatusConflict, gin.H{"error": services.ErrAlreadyVerified.Error()})
		return
	}

	if errors.Is(err, services.ErrExpiredToken) {
		c.JSON(http.StatusGone, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Error completing onboarding: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
}
