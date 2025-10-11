package repositories

import (
	"github.com/high-effort-low-stress/go-bank-api/models"
	"gorm.io/gorm"
)

type OnboardingRequestRepository interface {
	FindByDocumentOrEmail(document, email string) (*models.OnboardingRequest, error)
	Create(onboardingRequest *models.OnboardingRequest) error
}

type onboardingRequestRepository struct {
	db *gorm.DB
}

func NewOnboardingRequestRepository(db *gorm.DB) OnboardingRequestRepository {
	return &onboardingRequestRepository{db: db}
}

func (r *onboardingRequestRepository) FindByDocumentOrEmail(document, email string) (*models.OnboardingRequest, error) {
	var onboardingRequest models.OnboardingRequest
	result := r.db.Where("document_number = ? OR email = ?", document, email).First(&onboardingRequest)
	if result.Error != nil {
		return nil, result.Error
	}
	return &onboardingRequest, nil
}

// Create cria uma nova solicitação de onboarding no banco de dados.
func (r *onboardingRequestRepository) Create(onboardingRequest *models.OnboardingRequest) error {
	result := r.db.Create(onboardingRequest)
	return result.Error
}
