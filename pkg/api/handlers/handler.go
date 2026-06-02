package handlers

import (
	"errors"
	"warehouse/pkg/api/appresponse"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository // Direct access as you wanted
	Response   *appresponse.Response
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repository: repo,
		Response:   appresponse.NewResponse(),
	}
}
func GetIDFromContext(c *gin.Context) uint {
	if id, exists := c.Get("id"); exists {
		return id.(uint)
	}
	return 0
}

// Helper functions
func (h *Handler) handleManagerError(c *gin.Context, err error, title string) {
	if errors.Is(err, repository.ErrNotFound) {
		h.Response.NotFoundErr(c, title)
		return
	}
	h.Response.InternalServerErr(c, err.Error())
}
