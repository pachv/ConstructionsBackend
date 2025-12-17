package services

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
	"github.com/pachv/constructions/constructions/internal/repositories"
)

var ErrInvalidAskQuestion = errors.New("invalid ask question")

type AskQuestionService struct {
	logger       *slog.Logger
	repo         *repositories.AskQuestionRepository
	mailService  *MailSendingService
	notifyEmails []string
	templatePath string
	emailSubject string
}

func NewAskQuestionService(
	logger *slog.Logger,
	repo *repositories.AskQuestionRepository,
	mailService *MailSendingService,
	notifyEmails []string,
	templatePath string,
) *AskQuestionService {
	return &AskQuestionService{
		logger:       logger.With("component", "AskQuestionService"),
		repo:         repo,
		mailService:  mailService,
		notifyEmails: notifyEmails,
		templatePath: templatePath,
		emailSubject: "Новый вопрос с сайта",
	}
}

func (s *AskQuestionService) Create(message, name, email, product string) (entity.AskQuestion, error) {
	// минимальная валидация (не усложняю)
	if message == "" || name == "" || email == "" || product == "" {
		return entity.AskQuestion{}, ErrInvalidAskQuestion
	}

	q := entity.AskQuestion{
		Id:      uuid.NewString(),
		Message: message,
		Name:    name,
		Email:   email,
		Product: product,
	}

	if err := s.repo.Save(q); err != nil {
		s.logger.Error("failed to save ask question", "err", err)
		return entity.AskQuestion{}, err
	}

	// Письмо: если не отправилось — логируем, но НЕ валим запрос (обычно так лучше для UX)
	if len(s.notifyEmails) > 0 && s.mailService != nil {
		data := struct {
			ID      string
			Message string
			Name    string
			Email   string
			Product string
		}{
			ID:      q.Id,
			Message: q.Message,
			Name:    q.Name,
			Email:   q.Email,
			Product: q.Product,
		}

		if err := s.mailService.SendHTMLFromTemplate(
			s.notifyEmails,
			s.emailSubject,
			s.templatePath,
			data,
		); err != nil {
			s.logger.Error("failed to send ask question email", "err", err, "ask_question_id", q.Id)
		}
	}

	return q, nil
}
