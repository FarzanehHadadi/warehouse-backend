package handlers

import (
	"warehouse/pkg/repository"
)

type Handler struct {
	Repository *repository.Repository // Direct access as you wanted
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repository: repo,
	}
}
