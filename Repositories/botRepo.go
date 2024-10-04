package repositories

import (
	"errors"
	"fmt"
	domains "kesbekes/Domains"

	"gorm.io/gorm"
)

type TelegramRepository struct {
	database *gorm.DB
}

func NewTelegramRepository(db *gorm.DB) *TelegramRepository {
	return &TelegramRepository{database: db}
}

func (r *TelegramRepository) StoreChat(chatInfo *domains.ChatInfo, userID int64) error {
	// Check if chat with the same chat_id already exists
	var existingChat domains.ChatInfo
	err := r.database.Where("chat_id = ?", chatInfo.ChatID).First(&existingChat).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Retrieve the existing user
	existingUser, err := r.GetUser(userID)
	if err != nil {
		return err
	}

	// If chat exists, update the Users field
	if existingChat.ID != 0 {
		existingChat.Users = append(existingChat.Users, *existingUser)
		// Save the updated chat with the new user
		return r.database.Save(&existingChat).Error
	}

	// If chat doesn't exist, add the user to the new chat and create it
	chatInfo.Users = append(chatInfo.Users, *existingUser)
	err = r.database.Create(&chatInfo).Error
	if err != nil {
		return err
	}
	fmt.Println("---Chat saved successfully---")

	// Save the user with the new chat added
	return r.database.Save(&existingUser).Error
}

func (r *TelegramRepository) GetUser(UserID int64) (*domains.User, error) {
	var user domains.User
	if err := r.database.Where("user_id = ?", UserID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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
