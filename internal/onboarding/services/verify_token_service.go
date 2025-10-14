package services

import (
	"errors"
	"log"
	"time"

	"github.com/high-effort-low-stress/go-bank-api/internal/crypto"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/repositories"
	"gorm.io/gorm"
)

var (
	ErrInvalidToken    = errors.New("Token inválido ou expirado")
	ErrExpiredToken    = errors.New("Token expirado")
	ErrAlreadyVerified = errors.New("Usuário já verificado")
)

type VerifyEmailTokenService interface {
	Execute(token string) error
}

type verifyEmailTokenService struct {
	repo repositories.OnboardingRequestRepository
}

func NewVerifyEmailTokenService(repo repositories.OnboardingRequestRepository) VerifyEmailTokenService {
	return &verifyEmailTokenService{repo: repo}
}

func (s *verifyEmailTokenService) Execute(token string) error {
	hashedToken := crypto.HashTokenSHA256(token)

	onboardingRequest, err := s.repo.FindByVerificationTokenHash(hashedToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidToken
		}
		log.Printf("Error finding onboarding request by token hash: %v", err)
		return ErrInternalServer
	}

	if time.Now().After(onboardingRequest.TokenExpiresAt) {
		return ErrInvalidToken
	}

	if onboardingRequest.Status == models.StatusCompleted {
		return ErrAlreadyVerified
	}

	if onboardingRequest.Status == models.StatusVerified {
		return nil // Operational idempotency
	}

	onboardingRequest.Status = models.StatusVerified
	if err := s.repo.Update(onboardingRequest); err != nil {
		log.Printf("Error updating onboarding request status: %v", err)
		return ErrInternalServer
	}

	return nil
}
