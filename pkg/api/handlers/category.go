package handlers

import (
	"net/http"
	"warehouse/pkg/models"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleGetCategory(c *gin.Context) {
}
func (h *Handler) HandlePostCategory(c *gin.Context) {
	var cat *models.Category

	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cat, err := h.Repository.Category.Create(cat)
	if err != nil {
		switch err {
		case repository.ErrDuplicateKey:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, cat)

}
func (h *Handler) HandleDeleteCategory(c *gin.Context)    {}
func (h *Handler) HandlePatchCategory(c *gin.Context)     {}
func (h *Handler) HandleGetListCategories(c *gin.Context) {}
