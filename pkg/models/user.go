package models

type User struct {
	Phone    string `json:"phone" binding:"required,min=11"`
	Password string `json:"password" binding:"required,min=3,max=72"`
}
