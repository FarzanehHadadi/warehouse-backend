package appresponse

import (
	"fmt"
	"net/http"
	"time"
	"warehouse/pkg/api/apperr"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/events"

	"github.com/gin-gonic/gin"
)

type Response struct{}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

type PaginatedList[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
	Limit      int    `json:"limit"`
}

func NewPaginatedList[T any](items []T, cursor filter.CursorResponse, limit int) PaginatedList[T] {
	return PaginatedList[T]{
		Items:      items,
		NextCursor: cursor.NextCursor,
		HasMore:    cursor.HasMore,
		Limit:      limit,
	}
}

func (r *Response) ListSuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (r *Response) CreatedResponse(c *gin.Context, data interface{}, entityType events.EntityType, entityID uint, description string, payload interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
	events.Log(c, events.Created, entityType, entityID, description, payload)
}

func (r *Response) NoContentResponse(c *gin.Context, action events.Action, entityType events.EntityType, entityID uint, description string, payload interface{}) {
	c.Status(http.StatusNoContent)
	events.Log(c, action, entityType, entityID, description, payload)
}

// Error response
func (r *Response) Error(c *gin.Context, err error) {
	if appErr, ok := err.(*apperr.AppError); ok {
		c.JSON(appErr.StatusCode, gin.H{
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}

	// Unknown error → Internal Server Error
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Something went wrong",
		},
	})
}
func (r *Response) NotFoundErr(c *gin.Context, resource string) {
	r.Error(c, apperr.NotFoundError(resource))
}

func (r *Response) BadRequestErr(c *gin.Context, message string) {
	r.Error(c, apperr.BadRequestError(message))
}
func (r *Response) ConflictErr(c *gin.Context, message string) {
	r.Error(c, apperr.ConflictError(message))
}
func (r *Response) InternalServerErr(c *gin.Context, message string) {
	r.Error(c, apperr.InternalServerError(message))
}
func (r *Response) UnauthorizedErr(c *gin.Context, message string) {
	r.Error(c, apperr.UnauthorizedError(message))
}
func (r *Response) ForbiddenErr(c *gin.Context, message string) {
	r.Error(c, apperr.ForbiddenError(message))
}

func (r *Response) ExcelResponse(c *gin.Context, filename string, data []byte) {
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func ExportFilename(prefix string) string {
	return fmt.Sprintf("%s_%s.xlsx", prefix, time.Now().Format("2006-01-02"))
}
