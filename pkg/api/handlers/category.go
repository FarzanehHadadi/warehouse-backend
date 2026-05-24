package handlers

import (
	"errors"
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
func (h *Handler) HandleDeleteCategory(c *gin.Context) {
	catId := c.Param("categoryId")
	if catId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
	}
	id, err := strconv.ParseInt(catId, 10, 64)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}
	if err := h.Repository.Category.Delete(uint(id)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusNoContent, nil)

}
func (h *Handler) HandlePatchCategory(c *gin.Context) {
	catId := c.Param("categoryId")
	if catId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id is required"})
	}
	id, err := strconv.ParseInt(catId, 10, 64)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repository.Category.Update(uint(id), &cat); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
	})

}
func (h *Handler) HandleGetListCategories(c *gin.Context) {

	categories, err := h.Repository.Category.GetList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}
