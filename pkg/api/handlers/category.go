package handlers

import (
	"errors"
	"warehouse/pkg/models"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
)

// HandleGetCategory godoc
//
//	@Summary		Get category by ID
//	@Description	Get a single category by its ID
//	@Tags			Categories
//	@Accept			json
//	 @Security     ApiKeyAuth
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	dto.SuccessCategoryResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/categories/{id} [get]
func (h *Handler) HandleGetCategory(c *gin.Context) {
	id := GetIDFromContext(c)

	cat, err := h.Repository.Category.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Response.NotFoundErr(c, "Category")

			return
		}
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.SuccessResponse(c, cat)
}

// HandlePostCategory godoc
//
//	@Summary		Create a category
//	@Description	Create a category
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	 @Security     ApiKeyAuth
//	@Param			category	body		dto.Category	true	"Category object with updated data"
//	@Success		200			"No Content - Category successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/categories [post]
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

// HandleDeleteCategory godoc
//
//		@Summary		Delete a category
//		@Description	Delete a category
//		@Tags			Categories
//		@Accept			json
//		@Produce		json
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Param			id			path		int				true	"Category ID"
//	@Success		200			"No Content"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/categories/{id} [delete]
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

// HandlePatchCategory godoc
//
//	@Summary		Modify a category
//	@Description	Modify a category
//	@Tags			Categories
//	 @Security     ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"Category ID"
//	@Param			category	body		dto.Category	true	"Category object with updated data"
//	@Success		200			"No Content - Category successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/categories/{id} [patch]
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

// HandleGetCategory godoc
//
//	@Summary		Get list of categories
//	@Description	Get list of categories
//	@Tags			Categories
//	 @Security     ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.SuccessCategoriesResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/categories [get]
func (h *Handler) HandleGetListCategories(c *gin.Context) {

	categories, err := h.Repository.Category.GetList()
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())

		return
	}

	h.Response.SuccessResponse(c, categories)
}
