package appresponse

import (
	"net/http"
	"warehouse/pkg/api/apperr"

	"github.com/gin-gonic/gin"
)

type Response struct{}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (r *Response) CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func (r *Response) NoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
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
