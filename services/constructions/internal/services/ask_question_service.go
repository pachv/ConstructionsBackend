package services

import (
	"errors"
	"fmt"
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
func (s *AskQuestionService) Create(message, name, phone string, email, product *string, consent bool) (entity.AskQuestion, error) {
	if message == "" || name == "" || phone == "" {
		return entity.AskQuestion{}, ErrInvalidAskQuestion
	}
	if !consent {
		return entity.AskQuestion{}, ErrInvalidAskQuestion
	}

	q := entity.AskQuestion{
		Id:      uuid.NewString(),
		Message: message,
		Name:    name,
		Phone:   phone,
		Email:   email,
		Product: product,
		Consent: consent,
	}

	if err := s.repo.Save(q); err != nil {
		s.logger.Error("failed to save ask question", "err", err)
		return entity.AskQuestion{}, err
	}

	if s.mailService != nil && s.notifyEmails[0] != "" {
		fmt.Println("here")
		data := struct {
			ID      string
			Message string
			Name    string
			Phone   string
			Email   string
			Product string
			Consent bool
		}{
			ID:      q.Id,
			Message: q.Message,
			Name:    q.Name,
			Phone:   q.Phone,
			Email:   derefStr(q.Email),
			Product: derefStr(q.Product),
			Consent: q.Consent,
		}

		err := s.mailService.SendHTMLFromTemplate(
			s.notifyEmails,
			"Новый вопрос с сайта",
			s.templatePath,
			data,
		)
		if err != nil {
			s.logger.Error("cant send email : " + err.Error())
		}
	}

	return q, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
