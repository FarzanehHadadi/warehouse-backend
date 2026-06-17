package handlers

import (
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/api/mapper"
	"warehouse/pkg/models"

	"github.com/gin-gonic/gin"
)

// HandleGetproduct godoc
//
//	@Summary		Get product by ID
//	@Description	Get a single product by its ID
//	@Tags			Products
//	@Accept			json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Produce		json
//	@Param			id	path		int	true	"product ID"
//	@Success		200	{object}	dto.ProductSummary
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/products/{id} [get]
func (h *Handler) HandleGetProduct(c *gin.Context) {
	id := GetIDFromContext(c)
	prd, err := h.Repository.Product.FindByID(id)
	if err != nil {
		h.handleError(c, err, "Product")
		return
	}
	h.Response.SuccessResponse(c, mapper.ToProductDetailResponse(prd))
}

// @Summary      Get products List
// @Description  Retrieve products with filtering, search, and cursor pagination
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        name            query    string    false  "Filter by name (partial match)"
// @Param        created_after   query    string    false  "Created after date (YYYY-MM-DD)"
// @Param        created_from   query    string    false  "Created before date (YYYY-MM-DD)"
// @Param        sort_by         query    string    false  "Sort field" Enums(id,name,created_at)
// @Param        sort_order      query    string    false  "Sort direction" Enums(asc,desc)
// @Param        cursor          query    string    false  "Cursor for next page"
// @Param        limit           query    integer   false  "Number of items per page (max 100)" minimum(1) maximum(100)
// @Success      200  {object}  dto.productListResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /v1/products [get]
func (h *Handler) HandleGetProductList(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.ProductFilterConfig)

	products, cursorResp, err := h.Repository.Product.GetList(*req)
	if err != nil {
		h.handleError(c, err, "Product")
		return
	}
	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		mapper.ToProductListResponse(products),
		*cursorResp,
		req.Limit,
	))
}

// HandlePostProduct godoc
//
//	@Summary		Create a product
//	@Description	Create a product
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Param			product	body		models.ProductUpdate	true	"product object with updated data"
//	@Success		204			{object}	dto.ProductSummary
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/products [post]
func (h *Handler) HandlePostProduct(c *gin.Context) {
	var product *models.Product
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		h.handleError(c, err, "Product")
		return
	}
	prd, err := h.Repository.Product.Create(product)
	if err != nil {
		h.handleError(c, err, "Product")
		return
	}
	h.Response.CreatedResponse(c, prd)
}

// HandlePatchProduct godoc
//
//	@Summary		Modify a product
//	@Description	Modify a product
//	@Tags			Products
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"product ID"
//	@Param			product	body		models.ProductUpdate	true	"product object with updated data"
//	@Success		200			"No Content - product successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/products/{id} [patch]
func (h *Handler) HandlePatchProduct(c *gin.Context) {
	id := GetIDFromContext(c)
	var product *models.ProductUpdate
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		h.handleError(c, err, "product")
		return
	}
	if err := h.Repository.Product.Update(id, product); err != nil {
		h.handleError(c, err, "product")
		return
	}
	h.Response.NoContentResponse(c)
}

// HandleDeleteProduct godoc
//
//		@Summary		Delete a product
//		@Description	Delete a product
//		@Tags			Products
//		@Accept			json
//		@Produce		json
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Param			id			path		int				true	"product ID"
//	@Success		200			"No Content"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/products/{id} [delete]
func (h *Handler) HandleDeleteProduct(c *gin.Context) {
	id := GetIDFromContext(c)
	if err := h.Repository.Product.Delete(id); err != nil {
		h.handleError(c, err, "Product")
		return
	}
	h.Response.NoContentResponse(c)
}
