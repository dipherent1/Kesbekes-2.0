package domains

import "gorm.io/gorm"

type ChatInfo struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username"`
	ChatID   int64  `json:"chat_id" gorm:"uniqueindex;not null"`
}
