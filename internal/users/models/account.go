package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Account struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;column:id"`
	PublicID      string    `gorm:"type:varchar(26);unique;not null"`
	UserID        int64     `gorm:"not null"`
	User          User      `gorm:"foreignKey:UserID"`
	AccountNumber string    `gorm:"type:varchar(10);unique;not null"`
	AgencyNumber  string    `gorm:"type:varchar(6);not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (Account) TableName() string {
	return "user.accounts"
}

func (a *Account) BeforeCreate(_ *gorm.DB) (err error) {
	a.PublicID = ulid.Make().String()
	return
}
