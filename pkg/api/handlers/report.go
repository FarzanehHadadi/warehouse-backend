package handlers

import (
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/api/dto"
	"warehouse/pkg/api/export"
	"warehouse/pkg/api/filter"

	"github.com/gin-gonic/gin"
)

// HandleGetThresholdProximityReport godoc
//
//	@Summary		Threshold proximity report
//	@Description	Returns products whose current stock quantity is at or below the warning threshold, including category, unit, and store count.
//	@Tags			Reports
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Security		Bearer
//	@Param			product_id	query		integer	false	"Filter by product ID"
//	@Param			store_id	query		integer	false	"Filter by store ID"
//	@Param			category_id	query		integer	false	"Filter by category ID"
//	@Param			sort_by		query		string	false	"Sort field"		Enums(current_quantity,product_id)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc,desc)
//	@Param			cursor		query		string	false	"Cursor for next page"
//	@Param			limit		query		integer	false	"Number of items per page (max 100)"	minimum(1)	maximum(100)
//	@Success		200			{object}	dto.ThresholdProximityListResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/reports/threshold-proximity [get]
func (h *Handler) HandleGetThresholdProximityReport(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.ThresholdProximityFilterConfig)
	reports, cursorResp, err := h.Repository.Report.GetThresholdProximity(c.Request.Context(), *req, true)
	if err != nil {
		h.handleError(c, err, "Report")
		return
	}

	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		reports,
		*cursorResp,
		req.Limit,
	))
}

// HandleExportThresholdProximityReport godoc
//
//	@Summary		Export threshold proximity report to Excel
//	@Description	Export filtered threshold proximity report as an Excel (.xlsx) file. Supports the same query filters as the threshold proximity report list endpoint.
//	@Tags			Reports
//	@Accept			json
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Security		ApiKeyAuth
//	@Security		Bearer
//	@Param			product_id	query		integer	false	"Filter by product ID"
//	@Param			store_id	query		integer	false	"Filter by store ID"
//	@Param			category_id	query		integer	false	"Filter by category ID"
//	@Param			sort_by		query		string	false	"Sort field"		Enums(current_quantity,product_id)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc,desc)
//	@Success		200			{file}		binary	"Excel file download"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/reports/threshold-proximity/export [get]
func (h *Handler) HandleExportThresholdProximityReport(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.ThresholdProximityFilterConfig)
	reports, _, err := h.Repository.Report.GetThresholdProximity(c.Request.Context(), *req, false)
	if err != nil {
		h.handleError(c, err, "Report")
		return
	}
	customHeaders := map[string]string{
		"ProductID":        "Product ID",
		"ProductName":      "Product Name",
		"CategoryID":       "Category ID",
		"CategoryName":     "Category Name",
		"UnitID":           "Unit ID",
		"UnitName":         "Unit Name",
		"WarningThreshold": "Warning Threshold",
		"CurrentQuantity":  "Current Quantity",
		"StoresCount":      "Stores Count",
	}

	err = export.ExportToExcel(c, reports, "Threshold Proximity Report", "threshold_proximity_report_export", customHeaders)
	if err != nil {
		h.handleError(c, err, "Threshold Proximity Report Excel Export")
	}
}

// HandleGetStoreProductQuantitiesReport godoc
//
//	@Summary		Store product quantities report
//	@Description	Returns total product quantity per store, with optional filters by product and category.
//	@Tags			Reports
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Security		Bearer
//	@Param			product_id	query		integer	false	"Filter by product ID"
//	@Param			category_id	query		integer	false	"Filter by category ID"
//	@Param			sort_by		query		string	false	"Sort field"		Enums(total_quantity,product_id,store_id)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc,desc)
//	@Param			cursor		query		string	false	"Cursor for next page"
//	@Param			limit		query		integer	false	"Number of items per page (max 100)"	minimum(1)	maximum(100)
//	@Success		200			{object}	dto.StoreProductQuantityListResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/reports/store-product-quantities [get]
func (h *Handler) HandleGetStoreProductQuantitiesReport(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.StoreProductQuantityFilterConfig)
	reports, cursorResp, err := h.Repository.Report.GetStoreProductQuantities(c.Request.Context(), *req, true)
	if err != nil {
		h.handleError(c, err, "Report")
		return
	}

	h.Response.ListSuccessResponse(c, appresponse.NewPaginatedList(
		reports,
		*cursorResp,
		req.Limit,
	))
}

// HandleExportStoreProductQuantitiesReport godoc
//
//	@Summary		Export store product quantities report to Excel
//	@Description	Export filtered store product quantities report as an Excel (.xlsx) file. Supports the same query filters as the store product quantities report list endpoint.
//	@Tags			Reports
//	@Accept			json
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Security		ApiKeyAuth
//	@Security		Bearer
//	@Param			product_id	query		integer	false	"Filter by product ID"
//	@Param			category_id	query		integer	false	"Filter by category ID"
//	@Param			sort_by		query		string	false	"Sort field"		Enums(total_quantity,product_id,store_id)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc,desc)
//	@Success		200			{file}		binary	"Excel file download"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/reports/store-product-quantities/export [get]
func (h *Handler) HandleExportStoreProductQuantitiesReport(c *gin.Context) {
	req := dto.NewPaginationRequestFromConfig(c, filter.StoreProductQuantityFilterConfig)
	reports, _, err := h.Repository.Report.GetStoreProductQuantities(c.Request.Context(), *req, false)
	if err != nil {
		h.handleError(c, err, "Report")
		return
	}

	customHeaders := map[string]string{
		"ProductID":     "Product ID",
		"CategoryID":    "Category ID",
		"StoreID":       "Store ID",
		"TotalQuantity": "Total Quantity",
	}

	err = export.ExportToExcel(c, reports, "Store Product Quantities", "store_product_quantities_report_export", customHeaders)
	if err != nil {
		h.handleError(c, err, "Store Product Quantities Excel Export")
	}
}
