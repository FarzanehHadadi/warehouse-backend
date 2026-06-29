package handlers

import (
	"errors"
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/events"
	"warehouse/pkg/models"
	"warehouse/pkg/repository"
	"warehouse/pkg/utils"

	"github.com/gin-gonic/gin"
)

// HandleGetCategory godoc
//
//	@Summary		Get category by ID
//	@Description	Get a single category by its ID
//	@Tags			Categories
//	@Accept			json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	dto.SuccessCategoryResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/categories/{id} [get]
func (h *Handler) HandleGetCategory(c *gin.Context) {
	id := utils.GetIDFromContext(c)

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
//	 @Security     Bearer
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
	err := h.Repository.Category.Create(cat)
	if err != nil {
		switch err {
		case repository.ErrDuplicateKey:
			h.Response.BadRequestErr(c, err.Error())
		default:
			h.Response.InternalServerErr(c, err.Error())
		}
		return
	}

	h.Response.CreatedResponse(c, "", "Category", cat.ID, "Category created successfully", cat)

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
	id := utils.GetIDFromContext(c)

	if err := h.Repository.Category.Delete(uint(id)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Response.NotFoundErr(c, "Category")
			return
		}
		h.Response.InternalServerErr(c, err.Error())

		return

	}
	h.Response.NoContentResponse(c, events.Deleted, "Category", uint(id), "Category deleted successfully", nil)

}

// HandlePatchCategory godoc
//
//	@Summary		Modify a category
//	@Description	Modify a category
//	@Tags			Categories
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
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
	id := utils.GetIDFromContext(c)

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

	h.Response.NoContentResponse(c, events.Updated, "Category", uint(id), "Category updated successfully", cat)

}

// HandleGetListCategories godoc
//
//	@Summary		Get list of categories
//	@Description	Retrieve categories with filtering, search, and cursor pagination
//	@Tags			Categories
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			search			query	string	false	"Global search in name"
//	@Param			name			query	string	false	"Filter by name (partial match)"
//	@Param			created_after	query	string	false	"Created after date (YYYY-MM-DD)"
//	@Param			created_before	query	string	false	"Created before date (YYYY-MM-DD)"
//	@Param			sort_by			query	string	false	"Sort field" Enums(id,name,created_at)
//	@Param			sort_order		query	string	false	"Sort direction" Enums(asc,desc)
//	@Param			cursor			query	string	false	"Cursor for next page"
//	@Param			limit			query	integer	false	"Number of items per page (max 100)" minimum(1) maximum(100)
//	@Success		200	{object}	dto.SuccessCategoriesResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/categories [get]
func (h *Handler) HandleGetListCategories(c *gin.Context) {

	req := dto.NewPaginationRequestFromConfig(c, filter.SimpleFilterConfig)
	categories, cursorResp, err := h.Repository.Category.GetList(*req)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())

		return
	}

	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		categories,
		*cursorResp,
		req.Limit,
	))
}
