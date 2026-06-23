package handlers

import (
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/api/mapper"
	"warehouse/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandleGetOrder godoc
//
//	@Summary		Get order by ID
//	@Description	Get a single order by its ID
//	@Tags			Orders
//	@Accept			json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Produce		json
//	@Param			id	path		int	true	"order ID"
//	@Success		200	{object}	dto.OrderSummary
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/orders/{id} [get]
func (h *Handler) HandleGetOrder(c *gin.Context) {
	id := GetIDFromContext(c)
	order, err := h.Repository.Order.FindByID(id)
	if err != nil {
		h.handleError(c, err, "Order")
		return
	}
	h.Response.SuccessResponse(c, mapper.ToOrderDetailResponse(order))
}

// @Summary      Get orders List
// @Description  Retrieve orders with filtering, search, and cursor pagination
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        product_status      query    string    false  "Filter by product status" Enums(good,defective,unknown)
// @Param        type      query    string    false  "Filter by type" Enums(inbound,outbound)
// @Param        created_after   query    string    false  "Created after date (YYYY-MM-DD)"
// @Param        created_before   query    string    false  "Created before date (YYYY-MM-DD)"
// @Param        product_id       query    integer   false  "Filter by product ID"
// @Param        store_id         query    integer   false  "Filter by store ID"
// @Param        department_id    query    integer   false  "Filter by department ID"
// @Param        sort_by         query    string    false  "Sort field" Enums(id,name,created_at)
// @Param        sort_order      query    string    false  "Sort direction" Enums(asc,desc)
// @Param        cursor          query    string    false  "Cursor for next page"
// @Param        limit           query    integer   false  "Number of items per page (max 100)" minimum(1) maximum(100)
// @Success      200  {object}  dto.orderListResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /v1/orders [get]
func (h *Handler) HandleGetOrderList(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.OrderFilterConfig)
	ordersList, cursorResp, err := h.Repository.Order.GetList(*req)
	if err != nil {
		logger.Log.Error("error getting orders list", zap.Error(err))
		h.handleError(c, err, "Order")
		return
	}
	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		mapper.ToOrderSummaries(ordersList), *cursorResp, req.Limit),
	)
}

// HandlePostOrder godoc
//
//	@Summary		Create a order
//	@Description	Create a order
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Param			order	body		dto.CreateOrderRequest	true	"Order object"
//	@Success		201			{object}	dto.OrderSummary
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/orders [post]
func (h *Handler) HandlePostOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		logger.Log.Error("error binding body", zap.Error(err))
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	ord, err := h.Repository.Order.Create(mapper.ToOrderFromCreateRequest(&req))
	if err != nil {
		logger.Log.Error("error creating order", zap.Error(err))
		h.handleError(c, err, "Order")
		return
	}
	h.Response.CreatedResponse(c, mapper.ToOrderDetailResponse(ord))
}

// HandlePatchOrder godoc
//
//	@Summary		Modify a order
//	@Description	Modify a order
//	@Tags			Orders
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"order ID"
//	@Param			order	body		dto.UpdateOrderRequest	true	"Order object with updated data"
//	@Success		200			"No Content - order successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/orders/{id} [patch]
func (h *Handler) HandlePatchOrder(c *gin.Context) {
	id := GetIDFromContext(c)
	var req dto.UpdateOrderRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	if err := h.Repository.Order.Update(id, mapper.ToOrderUpdateFromRequest(&req)); err != nil {
		h.handleError(c, err, "Order")
		return

	}
	h.Response.NoContentResponse(c)
}

// HandleDeleteOrder godoc
//
//		@Summary		Delete a order
//		@Description	Delete a order
//		@Tags			Orders
//		@Accept			json
//		@Produce		json
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Param			id			path		int				true	"order ID"
//	@Success		200			"No Content"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/orders/{id} [delete]
func (h *Handler) HandleDeleteOrder(c *gin.Context) {
	id := GetIDFromContext(c)

	if err := h.Repository.Order.Delete(id); err != nil {
		h.handleError(c, err, "Order")
		return

	}
	h.Response.NoContentResponse(c)
}
