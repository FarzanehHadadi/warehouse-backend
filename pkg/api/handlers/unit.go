package handlers

import (
	"warehouse/pkg/models"
	"warehouse/pkg/repository"

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
	id := GetIDFromContext(c)

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
//	@Description	Get list of units
//	@Tags			Units
//	@Accept			json
//	@Produce		json
//	 @Security     ApiKeyAuth
//	@Success		200	{object}	dto.SuccessUnitListResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/units [get]
func (h *Handler) HandleGetUnitList(c *gin.Context) {
	unitsList, err := h.Repository.Unit.GetList()
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.SuccessResponse(c, unitsList)
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
	unitDb, err := h.Repository.Unit.Create(unit)
	if err != nil {
		switch err {
		case repository.ErrDuplicateKey:
			h.Response.ConflictErr(c, err.Error())
		default:
			h.Response.InternalServerErr(c, err.Error())

		}
		return
	}
	h.Response.CreatedResponse(c, unitDb)

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
	id := GetIDFromContext(c)
	if err := h.Repository.Unit.Delete(id); err != nil {
		//TODO : add switch for errors
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.NoContentResponse(c)
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
	id := GetIDFromContext(c)
	var unit *models.Unit

	if err := c.ShouldBindBodyWithJSON(&unit); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	err := h.Repository.Unit.Update(id, unit)
	//TODO : add switch for errors
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.NoContentResponse(c)
}
