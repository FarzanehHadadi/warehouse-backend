package handlers

import (
	"log"
	"warehouse/pkg/api/auth"
	"warehouse/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

// HandlePostCategory godoc
//
//	@Summary		Login
//	@Description	Login to warehouse
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			loginInfo	body		dto.UserDto	true	"enter your phone number and password to login"
//	 @Security     ApiKeyAuth
//	@Success		200			{object}	dto.SuccessAuthResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		401			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/auth/login [post]
func (h *Handler) HandleLogin(c *gin.Context) {
	var user dto.UserDto

	if err := c.ShouldBindJSON(&user); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	//get user from db,
	findUser, err := h.Repository.User.FindByPhone(user.Mobile, user.Password)
	log.Println("findUser,", findUser)
	if err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}

	// generate a token
	token, err := auth.GenerateToken(findUser.ID)
	log.Println("token,", token)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return

	}
	//return the token
	h.Response.SuccessResponse(c, token)
}
