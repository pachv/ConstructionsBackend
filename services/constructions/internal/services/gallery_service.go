package services

import (
	"context"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type GalleryRepo interface {
	GetAllCategories(ctx context.Context) ([]entity.GalleryCategory, error)
	GetPhotosByCategorySlug(ctx context.Context, categorySlug string) ([]entity.GalleryPhoto, error)
}

type GalleryService struct {
	repo GalleryRepo
}

func NewGalleryService(repo GalleryRepo) *GalleryService {
	return &GalleryService{repo: repo}
}

func (s *GalleryService) GetAllCategories(ctx context.Context) ([]entity.GalleryCategory, error) {
	return s.repo.GetAllCategories(ctx)
}

func (s *GalleryService) GetPhotosByCategorySlug(ctx context.Context, categorySlug string) ([]entity.GalleryPhoto, error) {
	return s.repo.GetPhotosByCategorySlug(ctx, categorySlug)
}
