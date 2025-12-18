package services

import (
	"context"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
	"github.com/pachv/constructions/constructions/internal/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllCategories(ctx context.Context) ([]entity.CatalogCategory, error) {
	return s.repo.GetAllCategories(ctx)
}

func (s *ProductService) GetAllSections(ctx context.Context) ([]entity.CatalogSection, error) {
	return s.repo.GetAllSections(ctx)
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]entity.CatalogProduct, error) {
	return s.repo.GetAllProducts(ctx)
}
