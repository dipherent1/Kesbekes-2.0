package repositories

import (
	domains "kesbekes/Domains"

	"gorm.io/gorm"
)

type TelegramRepository struct {
	database *gorm.DB
}

func NewTelegramRepository(db *gorm.DB) *TelegramRepository {
	return &TelegramRepository{database: db}
}

func (r *TelegramRepository) StoreChat(chatInfo *domains.ChatInfo) error {
	return r.database.Create(&chatInfo).Error
}

func (r *TelegramRepository) StoreUser(user *domains.User) error {
	return r.database.Create(&user).Error
}

// update user preferences by adding to existing preferences using user id
func (r *TelegramRepository) UpdateUserPreferences(user *domains.User) error {
	var existingUser domains.User
	if err := r.database.Where("user_id = ?", user.UserID).First(&existingUser).Error; err != nil {
		return err
	}

	// Assuming preferences are stored as a slice of strings
	existingUser.Preferenses = append(existingUser.Preferenses, user.Preferenses...)

	return r.database.Save(&existingUser).Error
}

// get user preferences by user id
func (r *TelegramRepository) GetUserPreferences(userID int64) ([]string, error) {
	var user domains.User
	if err := r.database.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user.Preferenses, nil
}

// delete user preferences by user id
func (r *TelegramRepository) DeleteUserPreferences(userID int64, deleteWord string) error {
	var user domains.User
	if err := r.database.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	// Assuming preferences are stored as a slice of strings
	for i, pref := range user.Preferenses {
		if pref == deleteWord {
			user.Preferenses = append(user.Preferenses[:i], user.Preferenses[i+1:]...)
			break
		}
	}

	return r.database.Save(&user).Error
}
