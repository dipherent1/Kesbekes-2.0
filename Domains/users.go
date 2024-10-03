package domains

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`

	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
	Role     string `json:"role"`
	UserID   int64  `json:"user_id"`

	Email       string         `json:"email" gorm:"uniqueindex"`
	Username    string         `json:"username" gorm:"uniqueindex"`
	Password    string         `json:"password"`
	Preferenses pq.StringArray `json:"preferenses" gorm:"type:text[]"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
