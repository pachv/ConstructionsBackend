package services

import (
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

var ErrInvalidReview = errors.New("invalid review")

type ReviewRepo interface {
	Create(rv entity.Review) error
	GetAllPublished() ([]entity.Review, error)
}

type ReviewService struct {
	logger *slog.Logger
	repo   ReviewRepo
}

func NewReviewService(logger *slog.Logger, repo ReviewRepo) *ReviewService {
	return &ReviewService{
		logger: logger.With("component", "ReviewService"),
		repo:   repo,
	}
}

// imagePath может быть пустым ("") если фото не прикрепили
func (s *ReviewService) Create(name, position, text string, rating int, imagePath string, consent bool) (entity.Review, error) {
	if name == "" || text == "" {
		return entity.Review{}, ErrInvalidReview
	}
	if !consent {
		return entity.Review{}, ErrInvalidReview
	}
	if rating < 1 || rating > 5 {
		return entity.Review{}, ErrInvalidReview
	}

	rv := entity.Review{
		Id:         uuid.NewString(),
		Name:       name,
		Position:   position,
		Text:       text,
		Rating:     rating,
		ImagePath:  imagePath,
		Consent:    consent,
		CanPublish: true, // ! change to false in production
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(rv); err != nil {
		s.logger.Error("failed to create review", "err", err)
		return entity.Review{}, err
	}

	return rv, nil
}

func (s *ReviewService) GetAllPublished() ([]entity.Review, error) {
	items, err := s.repo.GetAllPublished()
	if err != nil {
		s.logger.Error("failed to get published reviews", "err", err)
		return nil, err
	}
	return items, nil
}
