package models

type User struct {
	Basic
	Phone     string `json:"phone" gorm:"not null; uniqueIndex" `
	Password  string `json:"-"`
	Username  string `json:"user_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
}
type CreateUserRequest struct {
	Phone     string  `json:"phone" gorm:"not null; uniqueIndex" `
	Password  string  `json:"password" gorm:"not null; size:255"`
	Username  string  `json:"user_name"`
	Email     *string `json:"email"`
	FirstName string  `json:"first_name" gorm:"not null"`
	LastName  string  `json:"last_name" gorm:"not null"`
}
