package mapper

import (
	"warehouse/pkg/api/dto"
	"warehouse/pkg/models"
)

func ToOrderDetailResponse(order *models.Order) dto.OrderSummary {
	if order == nil {
		return dto.OrderSummary{}
	}
	return dto.OrderSummary{
		Quantity:      order.Quantity,
		ID:            order.ID,
		Price:         order.Price,
		ExpireDate:    order.ExpireDate,
		ProductStatus: order.ProductStatus,
		Description:   order.Description,
		Type:          order.Type,
		Product:       ToProductSummary(order.Product),
		Store:         ToStoreSummary(order.Store),
		Department:    ToDepartmentSummary(order.Department),
		ProductID:     order.ProductID,
		StoreID:       order.StoreID,
		DepartmentID:  order.DepartmentID,
	}
}

func ToOrderSummaries(orders []*models.Order) []dto.OrderSummary {
	if orders == nil {
		return []dto.OrderSummary{}
	}
	summaries := make([]dto.OrderSummary, len(orders))
	for i, order := range orders {
		summaries[i] = ToOrderDetailResponse(order)
	}
	return summaries
}

func ToOrderFromCreateRequest(req *dto.CreateOrderRequest) *models.Order {
	if req == nil {
		return nil
	}

	return &models.Order{
		ProductID:     req.ProductID,
		StoreID:       req.StoreID,
		DepartmentID:  req.DepartmentID,
		Type:          req.Type,
		Quantity:      req.Quantity,
		Price:         req.Price,
		ExpireDate:    req.ExpireDate,
		ProductStatus: req.ProductStatus,
		Description:   req.Description,
	}
}

func ToOrderUpdateFromRequest(req *dto.UpdateOrderRequest) *models.OrderUpdate {
	if req == nil {
		return nil
	}

	return &models.OrderUpdate{
		ProductID:     req.ProductID,
		StoreID:       req.StoreID,
		DepartmentID:  req.DepartmentID,
		Type:          req.Type,
		Quantity:      req.Quantity,
		Price:         req.Price,
		ExpireDate:    req.ExpireDate,
		ProductStatus: req.ProductStatus,
		Description:   req.Description,
	}
}
func ToProductSummary(product *models.Product) *dto.SimpleSummary {
	if product == nil {
		return nil
	}
	return &dto.SimpleSummary{
		ID:   product.ID,
		Name: product.Name,
	}
}
func ToProductSummaries(products []*models.Product) []dto.SimpleSummary {
	if products == nil {
		return []dto.SimpleSummary{}
	}
	summaries := make([]dto.SimpleSummary, len(products))
	for i, product := range products {
		summaries[i] = dto.SimpleSummary{
			ID:   product.ID,
			Name: product.Name,
		}
	}
	return summaries
}
func ToStoreSummary(store *models.Store) *dto.SimpleSummary {
	if store == nil {
		return nil
	}
	return &dto.SimpleSummary{
		ID:   store.ID,
		Name: store.Name,
	}
}
func ToDepartmentSummary(department *models.Department) *dto.SimpleSummary {
	if department == nil {
		return nil
	}
	return &dto.SimpleSummary{
		ID:   department.ID,
		Name: department.Name,
	}
}
