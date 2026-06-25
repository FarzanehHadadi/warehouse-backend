package handlers

import (
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/api/mapper"
	"warehouse/pkg/models"

	"github.com/gin-gonic/gin"
)

// HandleGetStore godoc
//
//	@Summary		Get store by ID
//	@Description	Get a single store by its ID
//	@Tags			Stores
//	@Accept			json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Produce		json
//	@Param			id	path		int	true	"store ID"
//	@Success		200	{object}	dto.StoreSummary
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/stores/{id} [get]
func (h *Handler) HandleGetStore(c *gin.Context) {
	id := GetIDFromContext(c)
	store, err := h.Repository.Store.FindByID(id)
	if err != nil {
		h.handleError(c, err, "Store")
		return
	}
	h.Response.SuccessResponse(c, mapper.ToStoreDetailResponse(store))
}

// @Summary      Get stores List
// @Description  Retrieve stores with filtering, search, and cursor pagination
// @Tags         Stores
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        name            query    string    false  "Filter by name (partial match)"
// @Param        manager_name            query    string    false  "Filter by manager name (partial match)"
// @Param        created_after   query    string    false  "Created after date (YYYY-MM-DD)"
// @Param        created_before   query    string    false  "Created before date (YYYY-MM-DD)"
// @Param        sort_by         query    string    false  "Sort field" Enums(id,name,created_at)
// @Param        sort_order      query    string    false  "Sort direction" Enums(asc,desc)
// @Param        cursor          query    string    false  "Cursor for next page"
// @Param        limit           query    integer   false  "Number of items per page (max 100)" minimum(1) maximum(100)
// @Success      200  {object}  dto.storeListResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /v1/stores [get]
func (h *Handler) HandleGetStoreList(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.StoreFilterConfig)

	storesList, cursorResp, err := h.Repository.Store.GetList(*req)
	if err != nil {
		h.handleError(c, err, "Store")
		return
	}
	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		mapper.ToStoreSummaries(storesList), *cursorResp, req.Limit),
	)
}

// HandlePostStore godoc
//
//	@Summary		Create a store
//	@Description	Create a store
//	@Tags			Stores
//	@Accept			json
//	@Produce		json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Param			store	body		models.StoreUpdate	true	"store object with updated data"
//	@Success		204			{object}	dto.StoreSummary
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/stores [post]
func (h *Handler) HandlePostStore(c *gin.Context) {
	var store *models.Store

	if err := c.ShouldBindBodyWithJSON(&store); err != nil {
		h.handleError(c, err, "Store")
		return
	}
	err := h.Repository.Store.Create(store)
	if err != nil {
		h.handleError(c, err, "Store")
		return
	}
	h.Response.CreatedResponse(c, "")
}

// HandlePatchStore godoc
//
//	@Summary		Modify a store
//	@Description	Modify a store
//	@Tags			Stores
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"store ID"
//	@Param			store	body		models.StoreUpdate	true	"store object with updated data"
//	@Success		200			"No Content - store successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/stores/{id} [patch]
func (h *Handler) HandlePatchStore(c *gin.Context) {
	id := GetIDFromContext(c)
	var store *models.StoreUpdate

	if err := c.ShouldBindBodyWithJSON(&store); err != nil {
		h.handleError(c, err, "Store")
		return
	}
	if err := h.Repository.Store.Update(id, store); err != nil {
		h.handleError(c, err, "Store")
		return

	}
	h.Response.NoContentResponse(c)
}

// HandleDeleteStore godoc
//
//		@Summary		Delete a store
//		@Description	Delete a store
//		@Tags			Stores
//		@Accept			json
//		@Produce		json
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Param			id			path		int				true	"store ID"
//	@Success		200			"No Content"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/stores/{id} [delete]
func (h *Handler) HandleDeleteStore(c *gin.Context) {
	id := GetIDFromContext(c)

	if err := h.Repository.Store.Delete(id); err != nil {
		h.handleError(c, err, "Store")
		return

	}
	h.Response.NoContentResponse(c)
}
