package handlers

import (
	"net/http"
	"warehouse/pkg/api/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleLogin(c *gin.Context) {
	var user dto.UserDto

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
