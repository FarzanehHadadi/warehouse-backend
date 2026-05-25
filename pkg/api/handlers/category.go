package handlers

import (
	"errors"
	"warehouse/pkg/models"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleGetCategory(c *gin.Context) {
	id := GetIDFromContext(c)

	cat, err := h.Repository.Category.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Response.NotFoundErr(c, "Category")

			return
		}
		h.Response.InternalServerErr(c, err.Error())
	}
	h.Response.SuccessResponse(c, cat)
}
func (h *Handler) HandlePostCategory(c *gin.Context) {
	var cat *models.Category

	if err := c.ShouldBindJSON(&cat); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	cat, err := h.Repository.Category.Create(cat)
	if err != nil {
		switch err {
		case repository.ErrDuplicateKey:
			h.Response.BadRequestErr(c, err.Error())
		default:
			h.Response.InternalServerErr(c, err.Error())
		}
		return
	}

	h.Response.CreatedResponse(c, cat)

}
func (h *Handler) HandleDeleteCategory(c *gin.Context) {
	id := GetIDFromContext(c)

	if err := h.Repository.Category.Delete(uint(id)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Response.NotFoundErr(c, "Category")
			return
		}
		h.Response.InternalServerErr(c, err.Error())

		return

	}
	h.Response.NoContentResponse(c)

}
func (h *Handler) HandlePatchCategory(c *gin.Context) {
	id := GetIDFromContext(c)

	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	if err := h.Repository.Category.Update(uint(id), &cat); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Response.NotFoundErr(c, "Category")

			return
		}
		h.Response.InternalServerErr(c, err.Error())

		return
	}

	h.Response.SuccessResponse(c, "")

}
func (h *Handler) HandleGetListCategories(c *gin.Context) {

	categories, err := h.Repository.Category.GetList()
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())

		return
	}

	h.Response.SuccessResponse(c, categories)
}
