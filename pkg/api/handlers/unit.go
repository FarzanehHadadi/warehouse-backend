package handlers

import (
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/events"
	"warehouse/pkg/models"
	"warehouse/pkg/repository"
	"warehouse/pkg/utils"

	"github.com/gin-gonic/gin"
)

// HandleGetUnit godoc
//
//	@Summary		Get unit by ID
//	@Description	Get a single unit by its ID
//	@Tags			Units
//	@Accept			json
//	 @Security     ApiKeyAuth
//
// @Security		Bearer
//
//	@Produce		json
//	@Param			id	path		int	true	"Unit ID"
//	@Success		200	{object}	dto.SuccessUnitResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/units/{id} [get]
func (h *Handler) HandleGetUnitById(c *gin.Context) {
	id := utils.GetIDFromContext(c)

	unit, err := h.Repository.Unit.FindByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			h.Response.NotFoundErr(c, err.Error())
		} else {
			h.Response.InternalServerErr(c, err.Error())

		}

		return
	}
	h.Response.SuccessResponse(c, unit)

}

// HandleGetUnitList godoc
//
//	@Summary		Get list of units
//	@Description	Retrieve units with filtering, search, and cursor pagination
//	@Tags			Units
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			search			query	string	false	"Global search in name"
//	@Param			name			query	string	false	"Filter by name (partial match)"
//	@Param			created_after	query	string	false	"Created after date (YYYY-MM-DD)"
//	@Param			created_before	query	string	false	"Created before date (YYYY-MM-DD)"
//	@Param			sort_by			query	string	false	"Sort field" Enums(id,name,created_at)
//	@Param			sort_order		query	string	false	"Sort direction" Enums(asc,desc)
//	@Param			cursor			query	string	false	"Cursor for next page"
//	@Param			limit			query	integer	false	"Number of items per page (max 100)" minimum(1) maximum(100)
//	@Success		200	{object}	dto.SuccessUnitListResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/units [get]
func (h *Handler) HandleGetUnitList(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.SimpleFilterConfig)
	unitsList, cursorResp, err := h.Repository.Unit.GetList(*req)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		unitsList,
		*cursorResp,
		req.Limit,
	))
}

// HandlePostUnit godoc
//
//	@Summary		Create a new unit
//	@Description	Create a new unit
//	@Tags			Units
//	@Accept			json
//	@Produce		json
//	@Param			unit	body	dto.Unit true	"Unit Item"
//	 @Security     ApiKeyAuth
//
// @Security		Bearer
//
//	@Success		200	{object}	dto.SuccessUnitResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/units [post]
func (h *Handler) HandlePostUnit(c *gin.Context) {
	var unit *models.Unit

	if err := c.ShouldBindBodyWithJSON(&unit); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	err := h.Repository.Unit.Create(unit)
	if err != nil {
		switch err {
		case repository.ErrDuplicateKey:
			h.Response.ConflictErr(c, err.Error())
		default:
			h.Response.InternalServerErr(c, err.Error())

		}
		return
	}
	h.Response.CreatedResponse(c, "", "Unit", unit.ID, "Unit created successfully", unit)

}

// HandleDeleteUnit godoc
//
//	@Summary		Delete unit
//	@Description	Delete an existing unit
//	@Tags			Units
//	@Accept			json
//	@Produce		json
//
// @Param			id		path	int		true	"Unit ID"
//
//	@Security     ApiKeyAuth
//
// @Security		Bearer
//
//	@Success		200	"No Content"
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/units/{id} [delete]
func (h *Handler) HandleDeleteUnit(c *gin.Context) {
	id := utils.GetIDFromContext(c)
	if err := h.Repository.Unit.Delete(id); err != nil {
		switch err {
		case repository.ErrNotFound:
			h.Response.NotFoundErr(c, err.Error())
		default:
			h.Response.InternalServerErr(c, err.Error())
		}
		return
	}
	h.Response.NoContentResponse(c, events.Deleted, "Unit", uint(id), "Unit deleted successfully", nil)
}

// HandlePatchUnit godoc
//
//	@Summary		Update an existing unit
//	@Description	Update an existing unit
//	@Tags			Units
//	@Accept			json
//	@Produce		json
//
// @Param			id		path	int		true	"Unit ID"
//
//	@Param			unit	body	dto.Unit true	"Unit Item"
//	 @Security     ApiKeyAuth
//
// @Security		Bearer
//
//	@Success		200 "No Content"
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/units/{id} [patch]
func (h *Handler) HandlePatchUnit(c *gin.Context) {
	id := utils.GetIDFromContext(c)
	var unit *models.Unit

	if err := c.ShouldBindBodyWithJSON(&unit); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	err := h.Repository.Unit.Update(id, unit)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			h.Response.NotFoundErr(c, err.Error())
		default:
			h.Response.InternalServerErr(c, err.Error())
		}
		return
	}
	h.Response.NoContentResponse(c, events.Updated, "Unit", uint(id), "Unit updated successfully", unit)
}
