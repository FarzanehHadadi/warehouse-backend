package handlers

import (
	"warehouse/pkg/api/auth"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/models"
	"warehouse/pkg/utils"

	"github.com/gin-gonic/gin"
)

// HandlePostCategory godoc
//
//	@Summary		Login
//	@Description	Login to warehouse
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			loginInfo	body	dto.UserDto	true	"enter your phone number and password to login"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.SuccessAuthResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/auth/login [post]
func (h *Handler) HandleLogin(c *gin.Context) {
	var user dto.UserDto

	if err := c.ShouldBindJSON(&user); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	//get user from db,
	findUser, err := h.Repository.User.FindByPhone(user.Mobile, user.Password)
	if err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	if !utils.CheckPasswordHash(user.Password, findUser.Password) {
		h.Response.BadRequestErr(c, "Invalid phone or password")
		return
	}
	// generate a token
	token, err := auth.GenerateToken(findUser.ID, auth.TokenTypeAccess)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	refreshToken, err := auth.GenerateToken(findUser.ID, auth.TokenTypeRefresh)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	userInfo := dto.SuccessAuthResponse{
		User: dto.UserLoginSummary{
			Phone:     findUser.Phone,
			Username:  findUser.Username,
			Email:     findUser.Email,
			FirstName: findUser.FirstName,
			LastName:  findUser.LastName,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}

	//return the token
	h.Response.SuccessResponse(c, userInfo)
}

// HandlePostRegister godoc
//
//	@Summary		Register
//	@Description	Register a new user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			registerInfo	body	models.CreateUserRequest	true	"enter your phone number and password to register"
//	@Security		ApiKeyAuth
//	@Security		AdminRegistrationKeyAuth
//	@Success		200	"No Content - User created successfully"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/auth/register [post]
func (h *Handler) HandlePostRegister(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}

	user := models.User{
		Phone:     req.Phone,
		Password:  req.Password,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := h.Repository.User.Create(&user); err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.SuccessResponse(c, "User created successfully")
}

// HandleRefreshToken godoc
//
//	@Summary		Refresh Token
//	@Description	Refresh the access token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			refreshToken	body	dto.RefreshTokenRequest	true	"enter your refresh token to refresh the access token"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.RefreshTokenResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/auth/refresh [post]
func (h *Handler) HandleRefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}

	// Parse refresh token
	token, err := auth.ValidateToken(req.RefreshToken, auth.TokenTypeRefresh)

	if err != nil {
		h.Response.UnauthorizedErr(c, "Invalid refresh token")
		return
	}

	newAccessToken, err := auth.GenerateToken(token.UserId, auth.TokenTypeAccess)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.SuccessResponse(c, &dto.RefreshTokenResponse{
		Token:        newAccessToken,
		RefreshToken: req.RefreshToken,
	})
}
