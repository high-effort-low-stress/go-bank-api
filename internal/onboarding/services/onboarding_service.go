package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/high-effort-low-stress/go-bank-api/internal/notification"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/repositories"
	"github.com/high-effort-low-stress/go-bank-api/internal/validators"

	"gorm.io/gorm"
)

var (
	ErrInvalidCPF     = errors.New("CPF inválido")
	ErrUserExists     = errors.New("O CPF ou E-mail já está cadastrado")
	ErrInternalServer = errors.New("Ocorreu um erro inesperado")
)

var websiteVerifyUrl = "verify"
var VerificationEmailTemplatePath = "templates/verification_email.html"
var subject = "Bem-vindo ao GoBank! Confirme seu e-mail."

type OnboardingService interface {
	StartOnboardingProcess(document, fullName, email string) error
}

type onboardingService struct {
	repo     repositories.OnboardingRequestRepository
	emailSvc notification.EmailService
	wg       *sync.WaitGroup
}

func NewOnboardingService(repo repositories.OnboardingRequestRepository, emailSvc notification.EmailService, wg *sync.WaitGroup) OnboardingService {
	return &onboardingService{repo: repo, emailSvc: emailSvc, wg: wg}
}

func (s *onboardingService) StartOnboardingProcess(document, fullName, email string) error {
	if !validators.IsValidCPF(document) {
		return ErrInvalidCPF
	}

	existingRequest, err := s.repo.FindByDocumentOrEmail(document, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking for existing onboarding request: %v", err)
		return ErrInternalServer
	}

	if existingRequest != nil {
		return ErrUserExists
	}

	rawToken, hashedToken, err := generateVerificationToken()
	if err != nil {
		log.Printf("Error generating verification token: %v", err)
		return ErrInternalServer
	}

	newRequest := &models.OnboardingRequest{
		FullName:              fullName,
		Email:                 email,
		DocumentNumber:        document,
		VerificationTokenHash: hashedToken,
		TokenExpiresAt:        time.Now().Add(1 * time.Hour),
		Status:                models.StatusPending,
	}

	if err := s.repo.Create(newRequest); err != nil {
		log.Printf("Error creating onboarding request: %v", err)
		return ErrInternalServer
	}

	if s.wg != nil {
		s.wg.Add(1)
	}

	go func() {
		if s.wg != nil {
			defer s.wg.Done()
		}
		if err := s.sendEmail(fullName, email, rawToken); err != nil {
			log.Printf("CRITICAL: Failed to send verification email to %s: %v", email, err)
		}
	}()

	log.Println("Onboarding process started successfully")
	return nil
}

func generateVerificationToken() (rawToken string, hashedToken string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	rawToken = hex.EncodeToString(bytes)

	hasher := sha256.New()
	hasher.Write([]byte(rawToken))
	hashedToken = hex.EncodeToString(hasher.Sum(nil))

	return rawToken, hashedToken, nil
}

func (s *onboardingService) sendEmail(fullName, email, rawToken string) error {
	verificationLink := fmt.Sprintf("%s/%s?token=%s", os.Getenv("WEBSITE_BASE_URL"), websiteVerifyUrl, rawToken)

	templateData := struct {
		FullName         string
		Subject          string
		VerificationLink string
	}{
		FullName:         fullName,
		Subject:          subject,
		VerificationLink: verificationLink,
	}

	emailRequest := &notification.EmailRequest{
		From:         os.Getenv("EMAIL_FROM"),
		To:           email,
		Subject:      subject,
		TemplatePath: VerificationEmailTemplatePath,
		TemplateData: templateData,
	}

	if err := s.emailSvc.SendEmail(emailRequest); err != nil {
		return err
	}

	return nil

}
