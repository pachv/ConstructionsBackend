package services

import (
	"context"
	"fmt"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
	"github.com/pachv/constructions/constructions/internal/repositories"
)

type SiteSectionsService struct {
	repo *repositories.SiteSectionsRepository
}

func NewSiteSectionsService(repo *repositories.SiteSectionsRepository) *SiteSectionsService {
	return &SiteSectionsService{repo: repo}
}

func (s *SiteSectionsService) GetAll(ctx context.Context) ([]entity.SiteSectionSummary, error) {
	items, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service site sections get all: %w", err)
	}
	return items, nil
}

func (s *SiteSectionsService) GetBySlugFull(ctx context.Context, slug string) (*entity.SiteSection, error) {
	sec, err := s.repo.GetBySlugFull(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("service site sections get by slug full: %w", err)
	}
	return sec, nil
}
