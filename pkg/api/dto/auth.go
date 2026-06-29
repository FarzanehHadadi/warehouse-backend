package dto

type UserDto struct {
	Mobile   string `json:"mobile" binding:"required,min=11"`
	Password string `json:"password" binding:"required,min=3,max=72"`
}
type SuccessAuthResponse struct {
	Token string           `json:"token"`
	User  UserLoginSummary `json:"user"`
}

type UserLoginSummary struct {
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
