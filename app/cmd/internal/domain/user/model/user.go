package model

import "time"

type User struct {
	ID         string    `json:"id" sql:"id"`
	Email      string    `json:"email" validate:"required" sql:"email"`
	Password   string    `json:"password" validate:"required" sql:"password"`
	Username   string    `json:"username" sql:"username"`
	TokenHash  string    `json:"token_hash" sql:"token_hash"`
	IsVerified bool      `json:"is_verified" sql:"is_verified"`
	CreatedAt  time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" sql:"updated_at"`
}
