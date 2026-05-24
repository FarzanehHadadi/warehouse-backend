package handlers

import (
	"net/http"
	"strconv"
	"warehouse/pkg/models"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleGetCategory(c *gin.Context) {
	catId := c.Param("categoryId")
	if catId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
	}
	id, err := strconv.ParseInt(catId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	cat, err := h.Repository.Category.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": cat,
	})
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

	c.JSON(http.StatusCreated, cat)

}
func (h *Handler) HandleDeleteCategory(c *gin.Context)    {}
func (h *Handler) HandlePatchCategory(c *gin.Context)     {}
func (h *Handler) HandleGetListCategories(c *gin.Context) {}
