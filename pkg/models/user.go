package models

import "time"

type User struct {
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"user_name"`
	Email     string    `json:"email"`
}
