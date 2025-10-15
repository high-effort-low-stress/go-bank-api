package repositories

import (
	"github.com/high-effort-low-stress/go-bank-api/internal/users/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUserWithAccount(user *models.User) (*models.User, *models.Account, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// todo criar conta
func (r *userRepository) CreateUserWithAccount(user *models.User) (*models.User, *models.Account, error) {
	var createdUser *models.User
	var createdAccount *models.Account

	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// account := &models.Account{
		// 	UserID:        user.ID,
		// 	AgencyNumber:  "0001", // Agência padrão
		// 	AccountNumber: generateAccountNumber(),
		// }

		// if err := tx.Create(account).Error; err != nil {
		// 	return err
		// }

		createdUser = user
		// createdAccount = account
		return nil
	})

	return createdUser, createdAccount, err
}
