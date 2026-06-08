package models

type User struct {
	Basic
	Phone    string `json:"phone" gorm:"not null; uniqueIndex" `
	Password string `json:"password"`
	Username string `json:"user_name"`
	Email    string `json:"email"`
}
