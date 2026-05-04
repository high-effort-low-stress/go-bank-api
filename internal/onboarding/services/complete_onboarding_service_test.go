package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/services"
	user_models "github.com/high-effort-low-stress/go-bank-api/internal/users/models"
	user_services "github.com/high-effort-low-stress/go-bank-api/internal/users/services"
	"github.com/high-effort-low-stress/go-bank-api/internal/utils/crypto"
	"github.com/high-effort-low-stress/go-bank-api/testutil/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCompleteOnboardingService_Execute_Success(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	mockCreateUserSvc := new(mocks.MockCreateUserService)
	service := services.NewCompleteOnboardingService(mockRepo, mockCreateUserSvc)

	token := "valid-token"
	password := "StrongPassword123!"
	hashedToken := crypto.HashTokenSHA256(token)

	request := &models.OnboardingRequest{
		FullName:       "John Doe",
		Email:          "john@example.com",
		DocumentNumber: "12345678900",
		Status:         models.StatusVerified,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)
	mockCreateUserSvc.On("Execute", mock.MatchedBy(func(req *user_services.CreateServiceRequest) bool {
		return req.Email == request.Email && req.Password == password
	})).Return(&user_models.User{}, &user_models.Account{}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.OnboardingRequest")).Return(nil)

	err := service.Execute(token, password, password)

	assert.NoError(t, err)
	assert.Equal(t, models.StatusCompleted, request.Status)
	mockRepo.AssertExpectations(t)
	mockCreateUserSvc.AssertExpectations(t)
}

func TestCompleteOnboardingService_Execute_PasswordsDoNotMatch(t *testing.T) {
	service := services.NewCompleteOnboardingService(nil, nil)

	err := service.Execute("token", "pass1", "pass2")

	assert.Error(t, err)
	assert.Equal(t, services.ErrPasswordsDoNotMatch, err)
}

func TestCompleteOnboardingService_Execute_WeakPassword(t *testing.T) {
	service := services.NewCompleteOnboardingService(nil, nil)

	// Senha curta e sem caracteres especiais
	err := service.Execute("token", "123", "123")

	assert.Error(t, err)
	assert.Equal(t, services.ErrWeakPassword, err)
}

func TestCompleteOnboardingService_Execute_TokenNotFound(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewCompleteOnboardingService(mockRepo, nil)

	token := "non-existent"
	password := "StrongPassword123!"
	hashedToken := crypto.HashTokenSHA256(token)

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(nil, gorm.ErrRecordNotFound)

	err := service.Execute(token, password, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidToken, err)
}

func TestCompleteOnboardingService_Execute_TokenExpired(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewCompleteOnboardingService(mockRepo, nil)

	token := "expired"
	password := "StrongPassword123!"
	hashedToken := crypto.HashTokenSHA256(token)

	request := &models.OnboardingRequest{
		TokenExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)

	err := service.Execute(token, password, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrExpiredToken, err)
}

func TestCompleteOnboardingService_Execute_AlreadyCompleted(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewCompleteOnboardingService(mockRepo, nil)

	token := "completed"
	password := "StrongPassword123!"
	hashedToken := crypto.HashTokenSHA256(token)

	request := &models.OnboardingRequest{
		Status:         models.StatusCompleted,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}
	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)

	err := service.Execute(token, password, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrAlreadyVerified, err)
}

func TestCompleteOnboardingService_Execute_RequestNotVerified(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewCompleteOnboardingService(mockRepo, nil)

	token := "pending"
	password := "StrongPassword123!"
	hashedToken := crypto.HashTokenSHA256(token)

	// Status ainda é Pending, não Verified
	request := &models.OnboardingRequest{
		Status:         models.StatusPending,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}
	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)

	err := service.Execute(token, password, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrRequestNotVerified, err)
}

func TestCompleteOnboardingService_Execute_CreateUserFailure(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	mockCreateUserSvc := new(mocks.MockCreateUserService)
	service := services.NewCompleteOnboardingService(mockRepo, mockCreateUserSvc)

	token := "token"
	password := "StrongPassword123!"
	hashedToken := crypto.HashTokenSHA256(token)

	request := &models.OnboardingRequest{
		Status:         models.StatusVerified,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)
	mockCreateUserSvc.On("Execute", mock.Anything).Return(nil, nil, errors.New("creation error"))

	err := service.Execute(token, password, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInternalServer, err)
}

func TestCompleteOnboardingService_Execute_UpdateRepoFailure(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	mockCreateUserSvc := new(mocks.MockCreateUserService)
	service := services.NewCompleteOnboardingService(mockRepo, mockCreateUserSvc)

	token := "token"
	password := "StrongPassword123!"
	hashedToken := crypto.HashTokenSHA256(token)
	request := &models.OnboardingRequest{
		Status:         models.StatusVerified,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)
	mockCreateUserSvc.On("Execute", mock.Anything).Return(&user_models.User{}, &user_models.Account{}, nil)
	mockRepo.On("Update", mock.Anything).Return(errors.New("db error"))

	err := service.Execute(token, password, password)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInternalServer, err)
}
