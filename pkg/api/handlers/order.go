package handlers

import (
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/export"
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

// HandleExportOrder godoc
//
//	@Summary		Export orders to Excel
//	@Description	Export filtered orders as an Excel (.xlsx) file. Supports the same query filters as the orders list endpoint.
//	@Tags			Orders
//	@Accept			json
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Security		ApiKeyAuth
//	@Security		Bearer
//	@Param			product_status	query	string	false	"Filter by product status"	Enums(good,defective,unknown)
//	@Param			type			query	string	false	"Filter by order type"		Enums(inbound,outbound)
//	@Param			created_after	query	string	false	"Created after date (YYYY-MM-DD)"
//	@Param			created_before	query	string	false	"Created before date (YYYY-MM-DD)"
//	@Param			product_id		query	integer	false	"Filter by product ID"
//	@Param			store_id		query	integer	false	"Filter by store ID"
//	@Param			department_id	query	integer	false	"Filter by department ID"
//	@Param			sort_by			query	string	false	"Sort field"					Enums(id,created_at)
//	@Param			sort_order		query	string	false	"Sort direction"				Enums(asc,desc)
//	@Success		200	{file}	binary	"Excel file download"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/orders/export [get]
func (h *Handler) HandleExportOrder(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.OrderFilterConfig)
	ordersList, err := h.Repository.Order.GetListNoPagination(*req)
	if err != nil {
		h.handleError(c, err, "Order")
		return
	}
	customHeaders := map[string]string{
		"Basic":         "ID",
		"Product":       "Product Name",
		"Store":         "Store Name",
		"Department":    "Department Name",
		"ExpireDate":    "Expire Date",
		"ProductStatus": "Product Status",
		"Type":          "Type",
		"ProductID":     "Product ID",
		"StoreID":       "Store ID",
		"DepartmentID":  "Department ID",
		"Quantity":      "Quantity",
		"Price":         "Price",
		"Description":   "Description",
	}

	err = export.ExportToExcel(c, ordersList, "Orders", "orders_export", customHeaders)
	if err != nil {
		h.handleError(c, err, "Excel Export")
	}

}

// func (h *Handler) HandleExportOrder(c *gin.Context) {
// 	req := dto.NewPaginationRequestFromConfig(c, filter.OrderFilterConfig)
// 	ordersList, err := h.Repository.Order.GetListNoPagination(*req)
// 	if err != nil {
// 		h.handleError(c, err, "Order")
// 		return
// 	}

// 	f := excelize.NewFile()
// 	defer f.Close()

// 	sheetName := "Orders"
// 	index, err := f.NewSheet(sheetName)
// 	if err != nil {
// 		h.handleError(c, err, "Excel")
// 		return
// 	}
// 	f.SetActiveSheet(index)

// 	// === Define Headers ===
// 	headers := []string{
// 		"ID", "Product", "Store", "Department", "Description",
// 		"Quantity", "Price", "Expire Date", "Product Status",
// 		"Type", "Product ID", "Store ID", "Department ID",
// 	}

// 	// Set headers
// 	for col, header := range headers {
// 		cell := fmt.Sprintf("%s1", string(rune('A'+col)))
// 		f.SetCellValue(sheetName, cell, header)
// 	}

// 	// === Create Styles ===
// 	headerStyle, _ := f.NewStyle(&excelize.Style{
// 		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
// 		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
// 		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
// 	})

// 	dateStyle, _ := f.NewStyle(&excelize.Style{
// 		NumFmt:    14, // Built-in date format (yyyy-mm-dd)
// 		Alignment: &excelize.Alignment{Horizontal: "center"},
// 	})

// 	numberStyle, _ := f.NewStyle(&excelize.Style{
// 		NumFmt:    1, // General number format
// 		Alignment: &excelize.Alignment{Horizontal: "right"},
// 	})

// 	// Apply header style and height
// 	f.SetRowHeight(sheetName, 1, 25)
// 	f.SetCellStyle(sheetName, "A1", "M1", headerStyle)

// 	// === Fill Data ===
// 	row := 2
// 	for _, order := range ordersList {
// 		colA := fmt.Sprintf("A%d", row)
// 		colB := fmt.Sprintf("B%d", row)
// 		colC := fmt.Sprintf("C%d", row)
// 		colD := fmt.Sprintf("D%d", row)
// 		colE := fmt.Sprintf("E%d", row)
// 		colF := fmt.Sprintf("F%d", row)
// 		colG := fmt.Sprintf("G%d", row)
// 		colH := fmt.Sprintf("H%d", row)
// 		colI := fmt.Sprintf("I%d", row)
// 		colJ := fmt.Sprintf("J%d", row)
// 		colK := fmt.Sprintf("K%d", row)
// 		colL := fmt.Sprintf("L%d", row)
// 		colM := fmt.Sprintf("M%d", row)

// 		f.SetCellValue(sheetName, colA, order.ID)

// 		if order.Product != nil {
// 			f.SetCellValue(sheetName, colB, order.Product.Name) // Change .Name if your field is different
// 		}
// 		if order.Store != nil {
// 			f.SetCellValue(sheetName, colC, order.Store.Name)
// 		}
// 		if order.Department != nil {
// 			f.SetCellValue(sheetName, colD, order.Department.Name)
// 		}

// 		f.SetCellValue(sheetName, colE, order.Description)
// 		f.SetCellValue(sheetName, colF, order.Quantity)
// 		f.SetCellValue(sheetName, colG, order.Price)

// 		// Expire Date
// 		if !order.ExpireDate.IsZero() {
// 			f.SetCellValue(sheetName, colH, order.ExpireDate.Time)
// 			f.SetCellStyle(sheetName, colH, colH, dateStyle)
// 		}

// 		f.SetCellValue(sheetName, colI, string(order.ProductStatus))
// 		f.SetCellValue(sheetName, colJ, string(order.Type))
// 		f.SetCellValue(sheetName, colK, order.ProductID)
// 		f.SetCellValue(sheetName, colL, order.StoreID)
// 		f.SetCellValue(sheetName, colM, order.DepartmentID)

// 		// Apply number style to quantity and price
// 		f.SetCellStyle(sheetName, colF, colG, numberStyle)

// 		row++
// 	}

// 	// Auto-fit column width
// 	for col := 'A'; col <= 'M'; col++ {
// 		f.SetColWidth(sheetName, string(col), string(col), 18)
// 	}

// 	// Save to buffer and send as download
// 	buffer, err := f.WriteToBuffer()
// 	if err != nil {
// 		h.handleError(c, err, "Excel")
// 		return
// 	}

// 	filename := fmt.Sprintf("orders_export_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))

// 	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// 	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
// 	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
// }

// func (h *Handler) HandleExportOrder(c *gin.Context) {
// 	req := dto.NewPaginationRequestFromConfig(c, filter.OrderFilterConfig)
// 	ordersList, err := h.Repository.Order.GetListNoPagination(*req)
// 	if err != nil {
// 		h.handleError(c, err, "Order")
// 		return
// 	}

// 	file, err := export.OrdersToExcel(ordersList)
// 	if err != nil {
// 		logger.Log.Error("error generating orders excel", zap.Error(err))
// 		h.Response.InternalServerErr(c, "failed to generate excel file")
// 		return
// 	}

//		h.Response.ExcelResponse(c, appresponse.ExportFilename("orders"), file)
//	}
