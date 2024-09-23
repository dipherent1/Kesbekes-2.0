package domains

import "time"

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`

	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
	Role     string `json:"role"`

	Email    string `json:"email" gorm:"uniqueindex;not null"`
	Username string `json:"username" gorm:"uniqueindex"`
	Password string `json:"password" gorm:"not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
