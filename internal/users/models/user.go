package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserStatus string

const (
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
	StatusBlocked  UserStatus = "BLOCKED"
)

type User struct {
	ID             int64      `gorm:"primaryKey;autoIncrement;column:id"`
	PublicID       string     `gorm:"type:varchar(26);unique;not null"`
	FullName       string     `gorm:"type:varchar(255);not null"`
	Email          string     `gorm:"type:varchar(255);unique;not null"`
	DocumentNumber string     `gorm:"type:varchar(11);unique;not null;column:document_number"`
	PasswordHash   string     `gorm:"type:varchar(255);not null"`
	Status         UserStatus `gorm:"type:user_status;not null;default:'ACTIVE'"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	DeactivatedAt  *time.Time `gorm:"column:deactivated_at"`
}

func (User) TableName() string {
	return "user.users"
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.PublicID = ulid.Make().String()
	return
}
