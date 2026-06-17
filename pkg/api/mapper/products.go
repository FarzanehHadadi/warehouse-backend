package mapper

import (
	"warehouse/pkg/api/dto"
	"warehouse/pkg/models"
)

func ToProductDetailResponse(product *models.Product) dto.ProductSummary {
	if product == nil {

		return dto.ProductSummary{}
	}
	return dto.ProductSummary{
		ID:               product.ID,
		Name:             product.Name,
		CategoryID:       product.CategoryID,
		WarningThreshold: product.WarningThreshold,
		UnitID:           product.UnitID,
		Category:         ToCategorySummaries(*product.Category),
		Unit:             ToUnitSummaries(*product.Unit),
	}
}

func ToProductListResponse(products []*models.Product) []dto.ProductSummary {
	if products == nil {
		return []dto.ProductSummary{}

	}
	summaries := make([]dto.ProductSummary, len(products))
	for i, product := range products {
		summaries[i] = ToProductDetailResponse(product)
	}
	return summaries
}
func ToUnitSummaries(unit models.Unit) dto.Summary {
	return dto.Summary{
		Name: unit.Name,
		ID:   unit.ID,
	}
}
func ToCategorySummaries(category models.Category) dto.Summary {
	return dto.Summary{
		Name: category.Name,
		ID:   category.ID,
	}
}
