package mocks

import (
	"github.com/high-effort-low-stress/go-bank-api/internal/notification"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/stretchr/testify/mock"
)

type MockOnboardingRepository struct {
	mock.Mock
}

func (m *MockOnboardingRepository) FindByDocumentOrEmail(document, email string) (*models.OnboardingRequest, error) {
	args := m.Called(document, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OnboardingRequest), args.Error(1)
}

func (m *MockOnboardingRepository) Create(req *models.OnboardingRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(req *notification.EmailRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockOnboardingRepository) FindByVerificationTokenHash(tokenHash string) (*models.OnboardingRequest, error) {
	args := m.Called(tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OnboardingRequest), args.Error(1)
}

func (m *MockOnboardingRepository) Update(req *models.OnboardingRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
