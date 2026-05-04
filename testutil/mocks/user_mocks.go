package mocks

import (
	"github.com/high-effort-low-stress/go-bank-api/internal/users/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/users/services"
	"github.com/stretchr/testify/mock"
)

type MockCreateUserService struct {
	mock.Mock
}

func (m *MockCreateUserService) Execute(req *services.CreateServiceRequest) (*models.User, *models.Account, error) {
	args := m.Called(req)

	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}

	var account *models.Account
	if args.Get(1) != nil {
		account = args.Get(1).(*models.Account)
	}

	return user, account, args.Error(2)
}
