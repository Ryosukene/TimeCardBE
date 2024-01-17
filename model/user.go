package model

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"unique"`
	Password   string    `json:"password"`
	Department string    `json:"department"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Email      string `json:"email" gorm:"unique"`
	Department string `json:"department"`
	Name       string `json:"name"`
}
