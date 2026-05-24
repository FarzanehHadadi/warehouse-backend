package dto

type UserDto struct {
	Mobile   string `json:"mobile" binding:"required,min=11"`
	Password string `json:"password" binding:"required,min=3,max=72"`
}
