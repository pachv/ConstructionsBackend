package services

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/pachv/constructions/constructions/internal/repositories"
)

type AdminEmailService struct {
	repo *repositories.AdminEmailRepository
}

func NewAdminEmailService(repo *repositories.AdminEmailRepository) *AdminEmailService {
	return &AdminEmailService{repo: repo}
}

var emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func (s *AdminEmailService) Get(ctx context.Context) (string, error) {
	item, err := s.repo.Get(ctx)
	if err != nil {
		return "", err
	}
	return item.Email, nil
}

func (s *AdminEmailService) Set(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is empty")
	}
	if !emailRe.MatchString(email) {
		return errors.New("invalid email")
	}
	return s.repo.Set(ctx, email)
}
