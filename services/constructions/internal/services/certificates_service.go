package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
	"github.com/pachv/constructions/constructions/internal/repositories"
)

type CertificateService struct {
	repo   *repositories.CertificateRepository
	logger *slog.Logger
}

func NewCertificateService(repo *repositories.CertificateRepository, logger *slog.Logger) *CertificateService {
	return &CertificateService{repo: repo, logger: logger}
}

func (s *CertificateService) GetAll(ctx context.Context) ([]entity.Certificate, error) {
	s.logger.Debug("CertificateService.GetAll: start")
	items, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("CertificateService.GetAll: repo failed", "err", err)
		return nil, fmt.Errorf("get certificates: %w", err)
	}
	return items, nil
}
