package mapper

import (
	"warehouse/pkg/api/dto"
	"warehouse/pkg/models"
)

func ToOrderDetailResponse(order *models.Order) dto.OrderSummary {

	return dto.OrderSummary{}

}

func ToOrderSummaries(stores []*models.Order) []dto.OrderSummary {

	return []dto.OrderSummary{}

}
