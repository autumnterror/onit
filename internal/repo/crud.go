package repo

import (
	"errors"
	"github.com/autumnterror/onit/internal/domain"
	"gorm.io/gorm"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mocks/mock_product_repo.go -package=mocks . ProductRepository
type ProductRepository interface {
	Create(product *domain.Product) error
	GetAll() ([]domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id string) error
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}
func (r *productRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Find(&products).Error
	return products, err
}
func (r *productRepository) GetByID(id string) (*domain.Product, error) {
	var product domain.Product

	err := r.db.First(&product, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}
func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}
func (r *productRepository) Delete(id string) error {
	return r.db.Delete(&domain.Product{}, "id = ?", id).Error
}
