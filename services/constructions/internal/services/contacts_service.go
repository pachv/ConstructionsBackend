package services

import (
	"context"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type ContactsRepo interface {
	GetEmail(ctx context.Context) (entity.ContactsEmailSetting, error)
	GetNumbers(ctx context.Context) ([]entity.ContactNumber, error)
	GetAddresses(ctx context.Context) ([]entity.ContactAddress, error)
}

type ContactsService struct {
	repo ContactsRepo
}

func NewContactsService(repo ContactsRepo) *ContactsService {
	return &ContactsService{repo: repo}
}

func (s *ContactsService) GetEmail(ctx context.Context) (entity.ContactsEmailSetting, error) {
	return s.repo.GetEmail(ctx)
}

func (s *ContactsService) GetNumbers(ctx context.Context) ([]entity.ContactNumber, error) {
	return s.repo.GetNumbers(ctx)
}

func (s *ContactsService) GetAddresses(ctx context.Context) ([]entity.ContactAddress, error) {
	return s.repo.GetAddresses(ctx)
}
