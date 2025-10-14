package services_test

import (
	"sync"
	"testing"

	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/services"
	"github.com/high-effort-low-stress/go-bank-api/testutil/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestStartOnboardingProcess_Success(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	mockEmailSvc := new(mocks.MockEmailService)

	var wg sync.WaitGroup
	service := services.NewOnboardingService(mockRepo, mockEmailSvc, &wg)
	fullName := "John Doe"
	email := "john.doe@example.com"
	validDocument := "68219090081"

	mockRepo.On("FindByDocumentOrEmail", validDocument, email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*models.OnboardingRequest")).Return(nil)
	mockEmailSvc.On("SendEmail", mock.AnythingOfType("*notification.EmailRequest")).Return(nil)

	err := service.StartOnboardingProcess(validDocument, fullName, email)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	wg.Wait()
	mockEmailSvc.AssertExpectations(t)
}

func TestStartOnboardingProcess_UserAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	mockEmail := new(mocks.MockEmailService)

	var wg sync.WaitGroup
	service := services.NewOnboardingService(mockRepo, mockEmail, &wg)

	fullName := "Jane Doe"
	email := "jane.doe@example.com"
	validDocument := "68219090081"

	existingRequest := &models.OnboardingRequest{}
	mockRepo.On("FindByDocumentOrEmail", validDocument, email).Return(existingRequest, nil)

	err := service.StartOnboardingProcess(validDocument, fullName, email)

	assert.Error(t, err)
	assert.Equal(t, services.ErrUserExists, err)
	mockRepo.AssertExpectations(t)

	wg.Wait()
	mockEmail.AssertNotCalled(t, "SendEmail", mock.Anything, mock.Anything, mock.Anything)
}

func TestStartOnboardingProcess_InvalidCPF(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	mockEmail := new(mocks.MockEmailService)
	var wg sync.WaitGroup
	service := services.NewOnboardingService(mockRepo, mockEmail, &wg)

	fullName := "Invalid User"
	email := "invalid@example.com"
	invalidDocument := "123"

	err := service.StartOnboardingProcess(invalidDocument, fullName, email)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidCPF, err)
	mockRepo.AssertNotCalled(t, "FindByDocumentOrEmail")

	wg.Wait()
	mockEmail.AssertNotCalled(t, "SendEmail", mock.Anything, mock.Anything, mock.Anything)
}
