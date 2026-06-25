package repository

import (
	"context"
	"warehouse/pkg/api/filter"
	"warehouse/pkg/models"

	"gorm.io/gorm"
)

type productRepository struct {
	*BaseRepository[models.Product]
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		BaseRepository: NewBaseRepository[models.Product](db),
	}
}
func (pr *productRepository) Create(product *models.Product) error {
	return pr.BaseRepository.Create(product)

}
func (pr *productRepository) FindByID(productId uint) (*models.Product, error) {

	return pr.BaseRepository.FindByID(productId, "Unit", "Category")

}
func (pr *productRepository) Delete(productId uint) error {

	return pr.BaseRepository.Delete(productId)

}
func (pr *productRepository) Update(productId uint, product *models.ProductUpdate) error {
	return pr.BaseRepository.Update(productId, product)

}
func (pr *productRepository) GetList(req filter.Request) ([]*models.Product, *filter.CursorResponse, error) {
	return pr.BaseRepository.GetList(req, filter.SimpleFilterConfig, "Unit", "Category")

}
func (pr *productRepository) Search(name string) ([]*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()
	query := pr.db.WithContext(ctx).Model(&models.Product{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	var products []*models.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	if len(products) > 100 {
		products = products[:100]
	}
	return products, nil
}
