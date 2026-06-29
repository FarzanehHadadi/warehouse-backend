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
type Entity string

const (
	Department Entity = "Department"
	Manager    Entity = "Manager"
	Unit       Entity = "Unit"
	Product    Entity = "Product"
	Store      Entity = "Store"
	Order      Entity = "Order"
	Report     Entity = "Report"
	Activity   Entity = "Activity"
)

func NewHandler(repo *repository.Repository) *Handler {

	return &Handler{
		Repository: repo,
		Response:   appresponse.NewResponse(),
	}
}

// Helper functions
func (h *Handler) handleError(c *gin.Context, err error, resourceName string) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		h.Response.NotFoundErr(c, resourceName)
	case errors.Is(err, repository.ErrDuplicateKey):
		h.Response.ConflictErr(c, err.Error())
	default:
		h.Response.InternalServerErr(c, err.Error())
	}
}
