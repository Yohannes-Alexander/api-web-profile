package models

import "time"

type User struct {
	ID        string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Name      string    `json:"name" gorm:"column:name"`
	Email     string    `json:"email" gorm:"column:email;uniqueIndex"`
	Password  string    `json:"-" gorm:"column:password"`
	Role      string    `json:"role" gorm:"column:role"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// optional refresh token model for migrations if you want
type RefreshToken struct {
	ID        string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	UserID    string    `json:"user_id" gorm:"column:user_id;type:uuid"`
	Token     string    `json:"token" gorm:"column:token;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}
