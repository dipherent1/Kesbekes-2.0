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
