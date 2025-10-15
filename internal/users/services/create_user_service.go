package services

import (
	"github.com/high-effort-low-stress/go-bank-api/internal/crypto"
	"github.com/high-effort-low-stress/go-bank-api/internal/users/models"
	"github.com/high-effort-low-stress/go-bank-api/internal/users/repositories"
)

type CreateServiceRequest struct {
	FullName       string
	Email          string
	DocumentNumber string
	Password       string
}

type CreateUserService interface {
	Execute(request *CreateServiceRequest) (*models.User, *models.Account, error)
}

type createUserService struct {
	userRepo repositories.UserRepository
}

func NewCreateUserService(userRepo repositories.UserRepository) CreateUserService {
	return &createUserService{userRepo: userRepo}
}

func (s *createUserService) Execute(request *CreateServiceRequest) (*models.User, *models.Account, error) {
	passwordHash, err := crypto.HashPassword(request.Password)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		FullName:       request.FullName,
		Email:          request.Email,
		DocumentNumber: request.DocumentNumber,
		PasswordHash:   passwordHash,
	}

	return s.userRepo.CreateUserWithAccount(user)
}
