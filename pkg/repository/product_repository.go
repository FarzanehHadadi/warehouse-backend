package repository

import (
	"context"
	"errors"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
func (pr *productRepository) Create(product *models.Product) (*models.Product, error) {
	result := pr.db.Model(&models.Product{}).Create(product)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, ErrDuplicateKey
		}
		return nil, result.Error
	}
	pr.db.Preload("Units").First(product, product.ID)
	pr.db.Preload("Categories").First(product, product.ID)

	return product, nil
}
func (pr *productRepository) FindByID(productId uint) (*models.Product, error) {

	var product *models.Product

	err := pr.db.Preload("Units").Preload("Categories").First(&product, productId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return product, nil
}
func (pr *productRepository) Delete(productId uint) error {

	result := pr.db.Delete(&models.Product{}, productId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
func (pr *productRepository) Update(productId uint, product *models.ProductUpdate) error {
	result := pr.db.Model(&models.Product{}).Where("id = ?", productId).Updates(product)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
func (pr *productRepository) GetList(req filter.Request) ([]*models.Product, *filter.CursorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query := pr.db.WithContext(ctx).Model(&models.Product{})
	query, err := filter.Apply(query, req, filter.ProductFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query, err = filter.ApplyCursor(query, req, filter.ProductFilterConfig)
	if err != nil {
		return nil, nil, err
	}
	query = query.Preload("Units")
	query = query.Preload("Categories")
	var products []*models.Product
	if err := query.Find(&products).Error; err != nil {

		return nil, nil, err
	}
	limit := req.Limit
	if limit < 0 {
		limit = 20
	}
	hasMore := len(products) > limit
	if hasMore {
		products = products[:limit]
	}
	var nextCursor string
	if len(products) > 0 {
		last := products[len(products)-1]
		if hasMore {
			nextCursor = filter.EncodeCursor(last.ID, last.CreatedAt)
		} else {
			nextCursor = ""
		}
	}

	return products, &filter.CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil

}
