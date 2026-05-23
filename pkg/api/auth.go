package api

import (
	"net/http"
	"warehouse/pkg/models"

	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
