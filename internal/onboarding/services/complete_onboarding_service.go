package services

import (
	"errors"
	"log"
	"time"

	"github.com/high-effort-low-stress/go-bank-api/internal/crypto"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/repositories"
	user_services "github.com/high-effort-low-stress/go-bank-api/internal/users/services"
	"github.com/high-effort-low-stress/go-bank-api/internal/validators"
	"gorm.io/gorm"
)

var (
	ErrRequestNotVerified  = errors.New("A solicitação de onboarding não foi verificada")
	ErrPasswordsDoNotMatch = errors.New("As senhas não coincidem")
	ErrWeakPassword        = errors.New("A senha não atende aos critérios de segurança")
)

type CompleteOnboardingService interface {
	Execute(token, password, confirmPassword string) error
}

type completeOnboardingService struct {
	onboardingRepo    repositories.OnboardingRequestRepository
	createUserService user_services.CreateUserService
}

func NewCompleteOnboardingService(
	onboardingRepo repositories.OnboardingRequestRepository,
	createUserService user_services.CreateUserService,
) CompleteOnboardingService {
	return &completeOnboardingService{onboardingRepo: onboardingRepo, createUserService: createUserService}
}

func (s *completeOnboardingService) Execute(token, password, confirmPassword string) error {
	if password != confirmPassword {
		return ErrPasswordsDoNotMatch
	}

	if !validators.ValidatePasswordPattern(password) {
		return ErrWeakPassword
	}

	hashedToken := crypto.HashTokenSHA256(token)
	onboardingRequest, err := s.onboardingRepo.FindByVerificationTokenHash(hashedToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidToken
		}
		log.Printf("Error finding onboarding request by token hash: %v", err)
		return ErrInternalServer
	}

	if time.Now().After(onboardingRequest.TokenExpiresAt) {
		return ErrExpiredToken
	}

	if onboardingRequest.Status == models.StatusCompleted {
		return ErrAlreadyVerified
	}

	if onboardingRequest.Status != models.StatusVerified {
		return ErrRequestNotVerified
	}

	_, _, err = s.createUserService.Execute(&user_services.CreateServiceRequest{
		FullName:       onboardingRequest.FullName,
		Email:          onboardingRequest.Email,
		DocumentNumber: onboardingRequest.DocumentNumber,
		Password:       password,
	})

	if err != nil {
		log.Printf("Error creating user and account: %v", err)
		return ErrInternalServer
	}

	onboardingRequest.Status = models.StatusCompleted
	err = s.onboardingRepo.Update(onboardingRequest)
	if err != nil {
		log.Printf("Error updating onboarding request: %v", err)
		return ErrInternalServer
	}

	return nil
}
