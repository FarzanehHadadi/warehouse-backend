package handlers

import (
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/mapper"
	"warehouse/pkg/models"

	"github.com/gin-gonic/gin"
)

// HandleGetManager godoc
//
//	@Summary		Get manager by ID
//	@Description	Get a single manager by its ID
//	@Tags			Managers
//	@Accept			json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Produce		json
//	@Param			id	path		int	true	"Manager ID"
//	@Success		200	{object}	dto.ManagerSummary
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/managers/{id} [get]
func (h *Handler) HandleGetManager(c *gin.Context) {
	id := GetIDFromContext(c)

	manager, err := h.Repository.Manager.FindByID(id)
	if err != nil {
		h.handleError(c, err, "Manager")
		return
	}

	h.Response.SuccessResponse(c, mapper.ToManagerDetailResponse(manager))
}

// HandlePostManager godoc
//
//	@Summary		Create a manager
//	@Description	Create a manager
//	@Tags			Managers
//	@Accept			json
//	@Produce		json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Param			Manager	body		dto.CreateManagerRequest	true	"Manager object with updated data"
//	@Success		204			{object}	dto.ManagerSummary
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/managers [post]
func (h *Handler) HandlePostManager(c *gin.Context) {
	var manager *models.Manager

	if err := c.ShouldBindJSON(&manager); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}

	dep, err := h.Repository.Manager.Create(manager)
	if err != nil {
		h.handleError(c, err, "Manager")

		return
	}
	h.Response.CreatedResponse(c, dep)
}

// HandlePatchManager godoc
//
//	@Summary		Modify a Manager
//	@Description	Modify a Manager
//	@Tags			Managers
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"Manager ID"
//	@Param			Manager	body		dto.CreateManagerRequest	true	"Manager object with updated data"
//	@Success		200			"No Content - Manager successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/managers/{id} [patch]
func (h *Handler) HandlePatchManager(c *gin.Context) {
	id := GetIDFromContext(c)
	var dep *models.ManagerUpdate
	if err := c.ShouldBindBodyWithJSON(&dep); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	if err := h.Repository.Manager.Update(id, dep); err != nil {
		h.handleError(c, err, "Manager")
		return
	}
	h.Response.NoContentResponse(c)

}

// HandleDeleteManager godoc
//
//		@Summary		Delete a manager
//		@Description	Delete a manager
//		@Tags			Managers
//		@Accept			json
//		@Produce		json
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Param			id			path		int				true	"Manager ID"
//	@Success		200			"No Content"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/managers/{id} [delete]
func (h *Handler) HandleDeleteManager(c *gin.Context) {
	id := GetIDFromContext(c)
	if err := h.Repository.Manager.Delete(id); err != nil {
		h.handleError(c, err, "Manager")
		return
	}
	h.Response.NoContentResponse(c)

}

// HandleGetManagerList godoc
//
//	@Summary		Get list of managers
//	@Description	Get list of managers
//	@Tags			Managers
//	 @Security     ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.ManagerListResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/managers [get]
func (h *Handler) HandleGetManagerList(c *gin.Context) {
	managers, total, err := h.Repository.Manager.GetList()
	if err != nil {
		h.handleError(c, err, "Manager")
		return
	}
	response := &dto.ManagerListResponse{
		Items: mapper.ToManagerSummaries(managers),
		Total: total,
	}

	h.Response.SuccessResponse(c, response)
}
