package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// OnboardingStatus define os possíveis status de um processo de onboarding.
type OnboardingStatus string

const (
	StatusPending   OnboardingStatus = "PENDING"
	StatusVerified  OnboardingStatus = "VERIFIED"
	StatusCompleted OnboardingStatus = "COMPLETED"
	StatusFailed    OnboardingStatus = "FAILED"
)

// OnboardingRequest representa a tabela onboarding_requests no banco de dados.
type OnboardingRequest struct {
	id                    int64            `gorm:"primaryKey;autoIncrement"`
	PublicID              string           `gorm:"type:varchar(26);unique;not null"`
	FullName              string           `gorm:"type:varchar(255);not null"`
	Email                 string           `gorm:"type:varchar(255);unique;not null"`
	DocumentNumber        string           `gorm:"type:varchar(11);unique;not null"`
	VerificationTokenHash string           `gorm:"type:varchar(255);unique;not null"`
	TokenExpiresAt        time.Time        `gorm:"not null"`
	Status                OnboardingStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`
	CreatedAt             time.Time        `gorm:"autoCreateTime"`
	UpdatedAt             time.Time        `gorm:"autoUpdateTime"`
}

// TableName define o nome da tabela para o GORM.
func (OnboardingRequest) TableName() string {
	return "onboarding.onboarding_requests"
}

func (or *OnboardingRequest) BeforeCreate(_ *gorm.DB) (err error) {
	or.PublicID = ulid.Make().String()
	return
}
