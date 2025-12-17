package services

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
	"github.com/pachv/constructions/constructions/internal/repositories"
)

type UserService struct {
	logger          *slog.Logger
	userRepository  *repositories.UserRepository
	passwordService *PasswordService
}

func NewUserService(userRepository *repositories.UserRepository,
	logger *slog.Logger, passwordService *PasswordService) *UserService {
	return &UserService{
		userRepository:  userRepository,
		passwordService: passwordService,
		logger:          logger.With("component", "UserService"),
	}
}

func (s *UserService) RegisterUser(surname, name, login, fathername, email, phoneNumber, password string) (userId string, err error) {

	exist, err := s.userRepository.DoesUserExist(login)
	if err != nil {
		s.logger.Error("error while checking user exist : " + err.Error())
		return "", err
	}

	if exist {
		s.logger.Error("such user already exist : " + login)
		return "", fmt.Errorf("user already exist")
	}

	hashedPassword, err := s.passwordService.Hash(password)
	if err != nil {
		s.logger.Error("error while hashing password : " + err.Error())
		return "", err
	}

	userId = uuid.NewString()

	err = s.userRepository.RegisterUser(userId, surname, name, login, fathername, email, phoneNumber, hashedPassword)
	if err != nil {
		s.logger.Error("error while creating user : " + err.Error())
		return "", err
	}

	return userId, nil
}

var ErrUserAlreadyExists = errors.New("user already exist")
var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *UserService) Login(login, password string) (entity.User, error) {
	u, err := s.userRepository.GetUserByLogin(login)
	if err != nil {
		return entity.User{}, ErrInvalidCredentials
	}

	ok, err := s.passwordService.Compare(u.HashedPassword, password)
	if err != nil {
		return entity.User{}, err
	}
	if !ok {
		return entity.User{}, ErrInvalidCredentials
	}

	// никогда не отдаём пароль наружу
	u.HashedPassword = ""
	return u, nil
}

func (s *UserService) GetMe(userID string) (entity.User, error) {
	u, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return entity.User{}, err
	}
	u.HashedPassword = ""
	return u, nil
}

func (s *UserService) ChangePassword(userID, newPassword string) error {
	hashed, err := s.passwordService.Hash(newPassword)
	if err != nil {
		return err
	}

	return s.userRepository.UpdateUserPassword(userID, hashed)
}
