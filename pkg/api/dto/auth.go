package dto

type UserDto struct {
	Mobile   string `json:"mobile" binding:"required,min=11" example:"09123456789"`
	Password string `json:"password" binding:"required,min=3,max=72" example:"123"`
}
type SuccessAuthResponse struct {
	Token        string           `json:"token"`
	RefreshToken string           `json:"refresh_token"`
	User         UserLoginSummary `json:"user"`
}

type UserLoginSummary struct {
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
