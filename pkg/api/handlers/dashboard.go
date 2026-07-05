package handlers

import (
	"warehouse/pkg/api/mapper"

	"github.com/gin-gonic/gin"
)

// HandleGetDashboard godoc
//
//	@Summary		Get dashboard
//	@Description	Get dashboard
//	@Tags			Dashboard
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Security		Bearer
//	@Produce		json
//	@Success		200	{object}	models.DashboardStats
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/dashboard [get]
func (h *Handler) HandleGetDashboard(c *gin.Context) {
	stats, err := h.Repository.Dashboard.GetStats()
	if err != nil {
		h.handleError(c, err, "Dashboard")
		return
	}
	h.Response.SuccessResponse(c, stats)
}

// HandleGetRecentActivities godoc
//
//	@Summary		Get recent activities
//	@Description	Get recent activities
//	@Tags			Dashboard
//	@Accept			json
//	@Security		ApiKeyAuth
//	@Security		Bearer
//	@Produce		json
//	@Success		200	{object}	dto.SuccessActivityResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/dashboard/activities [get]
func (h *Handler) HandleGetRecentActivities(c *gin.Context) {
	activities, err := h.Repository.Activity.GetRecent(10)
	if err != nil {
		h.Response.InternalServerErr(c, err.Error())
		return
	}
	h.Response.SuccessResponse(c, mapper.ToActivitySummaries(activities))
}
