package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/onboarding/services"
	"github.com/high-effort-low-stress/go-bank-api/internal/utils/crypto"
	"github.com/high-effort-low-stress/go-bank-api/testutil/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestVerifyEmailTokenService_Execute_Success(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewVerifyEmailTokenService(mockRepo)

	token := "valid-token"
	hashedToken := crypto.HashTokenSHA256(token)
	request := &models.OnboardingRequest{
		Status:         models.StatusPending,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.OnboardingRequest")).Return(nil)

	err := service.Execute(token)

	assert.NoError(t, err)
	assert.Equal(t, models.StatusVerified, request.Status)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailTokenService_Execute_TokenNotFound(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewVerifyEmailTokenService(mockRepo)

	token := "non-existent-token"
	hashedToken := crypto.HashTokenSHA256(token)

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(nil, gorm.ErrRecordNotFound)

	err := service.Execute(token)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidToken, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailTokenService_Execute_TokenExpired(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewVerifyEmailTokenService(mockRepo)

	token := "expired-token"
	hashedToken := crypto.HashTokenSHA256(token)
	request := &models.OnboardingRequest{
		Status:         models.StatusPending,
		TokenExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)

	err := service.Execute(token)

	assert.Error(t, err)
	assert.Equal(t, services.ErrExpiredToken, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailTokenService_Execute_AlreadyCompleted(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewVerifyEmailTokenService(mockRepo)

	token := "completed-token"
	hashedToken := crypto.HashTokenSHA256(token)
	request := &models.OnboardingRequest{
		Status:         models.StatusCompleted,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)

	err := service.Execute(token)

	assert.Error(t, err)
	assert.Equal(t, services.ErrAlreadyVerified, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailTokenService_Execute_AlreadyVerifiedIdempotency(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewVerifyEmailTokenService(mockRepo)

	token := "verified-token"
	hashedToken := crypto.HashTokenSHA256(token)
	request := &models.OnboardingRequest{
		Status:         models.StatusVerified,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)

	err := service.Execute(token)

	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailTokenService_Execute_DatabaseErrorOnFind(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewVerifyEmailTokenService(mockRepo)

	token := "token"
	hashedToken := crypto.HashTokenSHA256(token)

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(nil, errors.New("db error"))

	err := service.Execute(token)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyEmailTokenService_Execute_DatabaseErrorOnUpdate(t *testing.T) {
	mockRepo := new(mocks.MockOnboardingRepository)
	service := services.NewVerifyEmailTokenService(mockRepo)

	token := "token"
	hashedToken := crypto.HashTokenSHA256(token)
	request := &models.OnboardingRequest{
		Status:         models.StatusPending,
		TokenExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindByVerificationTokenHash", hashedToken).Return(request, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.OnboardingRequest")).Return(errors.New("db error"))

	err := service.Execute(token)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInternalServer, err)
	mockRepo.AssertExpectations(t)
}
