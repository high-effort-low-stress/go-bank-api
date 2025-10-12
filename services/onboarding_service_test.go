package services_test

import (
	"sync"
	"testing"

	"github.com/high-effort-low-stress/go-bank-api/models"
	"github.com/high-effort-low-stress/go-bank-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// --- Mocks ---

// MockOnboardingRepository é um mock para a interface OnboardingRequestRepository.
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

// MockEmailService é um mock para a interface EmailService.
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendVerificationEmail(fullName, to, token string) error {
	args := m.Called(fullName, to, token)
	return args.Error(0)
}

// --- Testes ---

func TestStartOnboardingProcess_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOnboardingRepository)
	mockEmailSvc := new(MockEmailService)

	// Use a WaitGroup to synchronize the test with the goroutine
	var wg sync.WaitGroup
	service := services.NewOnboardingService(mockRepo, mockEmailSvc, &wg)
	fullName := "John Doe"
	email := "john.doe@example.com"
	validDocument := "68219090081" // CPF válido para o teste

	mockRepo.On("FindByDocumentOrEmail", validDocument, email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*models.OnboardingRequest")).Return(nil)
	mockEmailSvc.On("SendVerificationEmail", fullName, email, mock.AnythingOfType("string")).Return(nil)

	// Act
	err := service.StartOnboardingProcess(validDocument, fullName, email)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Wait for the goroutine to finish before asserting its mock calls
	wg.Wait()
	mockEmailSvc.AssertExpectations(t)
}

func TestStartOnboardingProcess_UserAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockOnboardingRepository)
	mockEmail := new(MockEmailService)

	var wg sync.WaitGroup
	service := services.NewOnboardingService(mockRepo, mockEmail, &wg)

	fullName := "Jane Doe"
	email := "jane.doe@example.com"
	validDocument := "68219090081"

	// Configura o mock para simular que o usuário já existe
	existingRequest := &models.OnboardingRequest{}
	mockRepo.On("FindByDocumentOrEmail", validDocument, email).Return(existingRequest, nil)

	// Act
	err := service.StartOnboardingProcess(validDocument, fullName, email)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrUserExists, err)
	mockRepo.AssertExpectations(t)

	// In failure cases, we can assert immediately that the mock was not called, as there's no goroutine.
	wg.Wait()
	mockEmail.AssertNotCalled(t, "SendVerificationEmail", mock.Anything, mock.Anything, mock.Anything)
}

func TestStartOnboardingProcess_InvalidCPF(t *testing.T) {
	// Arrange
	mockRepo := new(MockOnboardingRepository)
	mockEmail := new(MockEmailService)
	var wg sync.WaitGroup
	service := services.NewOnboardingService(mockRepo, mockEmail, &wg)

	fullName := "Invalid User"
	email := "invalid@example.com"
	document := "123" // CPF inválido

	// Act
	err := service.StartOnboardingProcess(document, fullName, email)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidCPF, err)
	mockRepo.AssertNotCalled(t, "FindByDocumentOrEmail")

	wg.Wait()
	mockEmail.AssertNotCalled(t, "SendVerificationEmail", mock.Anything, mock.Anything, mock.Anything)
}
