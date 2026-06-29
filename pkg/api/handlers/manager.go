package handlers

import (
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/api/mapper"
	"warehouse/pkg/events"
	"warehouse/pkg/models"
	"warehouse/pkg/utils"

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
	id := utils.GetIDFromContext(c)

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

	err := h.Repository.Manager.Create(manager)
	if err != nil {
		h.handleError(c, err, "Manager")

		return
	}
	h.Response.CreatedResponse(c, "", events.Manager, manager.ID, "Manager created successfully", manager)
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
	id := utils.GetIDFromContext(c)
	var dep *models.ManagerUpdate
	if err := c.ShouldBindBodyWithJSON(&dep); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	if err := h.Repository.Manager.Update(id, dep); err != nil {
		h.handleError(c, err, "Manager")
		return
	}
	h.Response.NoContentResponse(c, events.Updated, events.Manager, uint(id), "Manager updated successfully", dep)

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
	id := utils.GetIDFromContext(c)
	if err := h.Repository.Manager.Delete(id); err != nil {
		h.handleError(c, err, "Manager")
		return
	}
	h.Response.NoContentResponse(c, events.Deleted, events.Manager, uint(id), "Manager deleted successfully", nil)

}

// @Summary      Get Managers List
// @Description  Retrieve managers with filtering, search, and cursor pagination
// @Tags         Managers
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        search          query    string    false  "Global search in name, phone, email"
// @Param        name            query    string    false  "Filter by name (partial match)"
// @Param        phone            query    string    false  "Filter by phone (partial match)"
// @Param        created_after   query    string    false  "Created after date (YYYY-MM-DD)"
// @Param        created_before   query    string    false  "Created before date (YYYY-MM-DD)"
// @Param        sort_by         query    string    false  "Sort field" Enums(id,name,created_at)
// @Param        sort_order      query    string    false  "Sort direction" Enums(asc,desc)
// @Param        cursor          query    string    false  "Cursor for next page"
// @Param        limit           query    integer   false  "Number of items per page (max 100)" minimum(1) maximum(100)
// @Success      200  {object}  dto.ManagerListResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /v1/managers [get]
func (h *Handler) HandleGetManagerList(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.ManagerFilterConfig)
	managers, cursorResp, err := h.Repository.Manager.GetList(*req)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		mapper.ToManagerSummaries(managers),
		*cursorResp,
		req.Limit,
	))
}
