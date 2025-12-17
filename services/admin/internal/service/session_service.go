package service

import (
	"log/slog"

	"github.com/is_backend/services/admin/internal/domain/entity"
	"github.com/is_backend/services/admin/internal/repository"
)

type SessionService struct {
	sessionRepo *repository.SessionRepository
	logger      *slog.Logger
}

func NewSessionService(
	sessionRepo *repository.SessionRepository,
	logger *slog.Logger) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		logger:      logger.With("component", "SessionService"),
	}
}

func (s *SessionService) GetSession(sessionId string) (*entity.UserSession, error) {

	session, err := s.sessionRepo.GetSessionBySessionId(sessionId)
	if err != nil {
		s.logger.Error("cant get session : " + err.Error())
		return nil, err
	}

	return session, nil
}

func (s *SessionService) CreateSession(username string, userId string) (sessionId string, err error) {

	sessionId, err = s.sessionRepo.CreateSession(username, userId)
	if err != nil {
		s.logger.Error("cant create session : " + err.Error())
		return "", err
	}

	return sessionId, nil
}

func (s *SessionService) Delete(id string) error {

	err := s.sessionRepo.DeleteSession(id)
	if err != nil {
		s.logger.Error("cant delete sessions : " + err.Error())
		return err
	}

	return nil
}
