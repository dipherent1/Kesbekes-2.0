package domains

import "gorm.io/gorm"

type ChatInfo struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username" gorm: uniqueindex`
	ChatID   int64  `json:"chat_id" gorm:"uniqueindex;not null"`
	Users    []User `json:"users" gorm:"many2many:user_chats;"`
}
