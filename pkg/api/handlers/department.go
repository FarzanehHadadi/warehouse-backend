package handlers

import (
	"errors"
	"warehouse/pkg/models"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
)

// HandleGetDepartment godoc
//
//	@Summary		Get department by ID
//	@Description	Get a single department by its ID
//	@Tags			Departments
//	@Accept			json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Produce		json
//	@Param			id	path		int	true	"Department ID"
//	@Success		200	{object}	dto.DepartmentListResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/departments/{id} [get]
func (h *Handler) HandleGetDepartment(c *gin.Context) {
	id := GetIDFromContext(c)
	department, err := h.Repository.Department.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.Response.NotFoundErr(c, "Department")

		} else {
			h.Response.InternalServerErr(c, err.Error())
		}
		return

	}

	h.Response.SuccessResponse(c, department)
}

// HandlePostDepartment godoc
//
//	@Summary		Create a department
//	@Description	Create a department
//	@Tags			Departments
//	@Accept			json
//	@Produce		json
//	 @Security     ApiKeyAuth
//	 @Security     Bearer
//	@Param			department	body		dto.Department	true	"Department object with updated data"
//	@Success		200			"No Content - Department successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/departments [post]
func (h *Handler) HandlePostDepartment(c *gin.Context) {
	var department *models.Department

	if err := c.ShouldBindJSON(&department); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}

	dep, err := h.Repository.Department.Create(department)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			h.Response.ConflictErr(c, err.Error())
		} else {
			h.Response.InternalServerErr(c, err.Error())

		}
		return
	}
	h.Response.CreatedResponse(c, dep)
}

// HandlePatchDepartment godoc
//
//	@Summary		Modify a department
//	@Description	Modify a department
//	@Tags			Departments
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int				true	"Department ID"
//	@Param			Department	body		dto.Department	true	"Department object with updated data"
//	@Success		200			"No Content - Department successfully updated"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/departments/{id} [patch]
func (h *Handler) HandlePatchDepartment(c *gin.Context) {
	id := GetIDFromContext(c)
	var dep *models.DepartmentUpdate
	if err := c.ShouldBindBodyWithJSON(&dep); err != nil {
		h.Response.BadRequestErr(c, err.Error())
		return
	}
	if err := h.Repository.Department.Update(id, dep); err != nil {
		switch err {
		case repository.ErrNotFound:
			h.Response.NotFoundErr(c, "Department")
		default:
			h.Response.InternalServerErr(c, err.Error())

		}
		return
	}
	h.Response.NoContentResponse(c)

}

// HandleDeleteDepartment godoc
//
//		@Summary		Delete a department
//		@Description	Delete a department
//		@Tags			Departments
//		@Accept			json
//		@Produce		json
//	 @Security     ApiKeyAuth
//
// @Security     Bearer
//
//	@Param			id			path		int				true	"Department ID"
//	@Success		200			"No Content"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v1/departments/{id} [delete]
func (h *Handler) HandleDeleteDepartment(c *gin.Context) {
	id := GetIDFromContext(c)
	if err := h.Repository.Department.Delete(id); err != nil {
		switch err {
		case repository.ErrNotFound:
			h.Response.NotFoundErr(c, "Department")
		default:
			h.Response.InternalServerErr(c, err.Error())

		}
		return
	}
	h.Response.NoContentResponse(c)

}

// HandleGetDepartmentList godoc
//
//	@Summary		Get list of departments
//	@Description	Get list of departments
//	@Tags			Departments
//	 @Security     ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.DepartmentListResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/v1/departments [get]
func (h *Handler) HandleGetDepartmentList(c *gin.Context) {
	depList, err := h.Repository.Department.GetList()
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			h.Response.NotFoundErr(c, "Department")
		default:
			h.Response.InternalServerErr(c, err.Error())

		}
		return
	}
	h.Response.SuccessResponse(c, depList)
}
