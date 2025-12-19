package services

import (
	"context"
	"fmt"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
	"github.com/pachv/constructions/constructions/internal/repositories"
)

type SiteSectionService struct {
	repo *repositories.SiteSectionRepository
}

func NewSiteSectionService(repo *repositories.SiteSectionRepository) *SiteSectionService {
	return &SiteSectionService{repo: repo}
}

func (s *SiteSectionService) GetAll(ctx context.Context) ([]entity.SiteSectionSummary, error) {
	items, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service get all site sections: %w", err)
	}
	return items, nil
}

func (s *SiteSectionService) GetBySlug(ctx context.Context, slug string) (*entity.SiteSection, error) {
	item, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("service get site section by slug: %w", err)
	}
	return item, nil
}
