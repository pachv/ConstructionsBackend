package services

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

var ErrInvalidCallbackRequest = errors.New("invalid callback request")

type CallbackRepo interface {
	Save(cb entity.CallbackRequest) error
}

type CallbackService struct {
	logger       *slog.Logger
	repo         CallbackRepo
	mailService  *MailSendingService
	notifyEmail  string
	templatePath string
	subject      string
}

func NewCallbackService(
	logger *slog.Logger,
	repo CallbackRepo,
	mailService *MailSendingService,
	notifyEmail string,
	templatePath string,
) *CallbackService {
	return &CallbackService{
		logger:       logger.With("component", "CallbackService"),
		repo:         repo,
		mailService:  mailService,
		notifyEmail:  notifyEmail,
		templatePath: templatePath,
		subject:      "Заявка на обратный звонок",
	}
}

func (s *CallbackService) Create(name, phone string, consent bool) (entity.CallbackRequest, error) {
	if name == "" || phone == "" || !consent {
		return entity.CallbackRequest{}, ErrInvalidCallbackRequest
	}

	cb := entity.CallbackRequest{
		Id:      uuid.NewString(),
		Name:    name,
		Phone:   phone,
		Consent: consent,
	}

	if err := s.repo.Save(cb); err != nil {
		s.logger.Error("failed to save callback request", "err", err)
		return entity.CallbackRequest{}, err
	}

	// письмо админу (не валим запрос, если письмо не ушло)
	if s.mailService != nil && s.notifyEmail != "" {
		data := struct {
			ID      string
			Name    string
			Phone   string
			Consent bool
		}{
			ID:      cb.Id,
			Name:    cb.Name,
			Phone:   cb.Phone,
			Consent: cb.Consent,
		}

		if err := s.mailService.SendHTMLFromTemplate(
			[]string{s.notifyEmail},
			s.subject,
			s.templatePath,
			data,
		); err != nil {
			s.logger.Error("failed to send callback email", "err", err, "callback_id", cb.Id)
		}
	}

	return cb, nil
}
