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

// update user preferenses by using user id
func (r *TelegramRepository) UpdateUserPreferenses(user *domains.User) error {
	return r.database.Model(&user).Where("user_id = ?", user.UserID).Update("preferenses", user.Preferenses).Error
}
