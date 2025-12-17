package service

import (
	"fmt"
	"log/slog"

	"github.com/is_backend/services/admin/internal/domain/entity"
	"github.com/is_backend/services/admin/internal/repository"
)

type AuthService struct {
	UserRepository        *repository.UserRepository
	passwordHasherService *PasswordHasherService
	logger                *slog.Logger
}

func NewAuthService(
	UserRepository *repository.UserRepository,
	passwordHasherService *PasswordHasherService,
	logger *slog.Logger) *AuthService {
	return &AuthService{
		UserRepository: UserRepository,
		logger:         logger.With("component", "AuthService"),
	}
}

func (s *AuthService) LoginUser(username string, password string) (*entity.User, error) {
	user, err := s.UserRepository.GetUserByUsername(username)
	if err != nil {
		s.logger.Error("cant get users : " + err.Error())
		return nil, err
	}

	correct, err := s.passwordHasherService.IsPasswordCorrect(password, user.HashedPassword)
	if err != nil {
		s.logger.Error("password is wrong : " + err.Error())
		return nil, err
	}

	if !correct {
		s.logger.Error("password is wrong : " + err.Error())
		return nil, err
	}

	s.logger.Debug("auth successful")

	return user, nil
}

func (s *AuthService) GetUsers(page int, search, orderBy string) ([]*entity.User, int, error) {

	users, pageAmount, err := s.UserRepository.GetUsers(page, search, orderBy)
	if err != nil {
		s.logger.Error("cant get users : " + err.Error())
		return nil, 0, err
	}

	return users, pageAmount, nil
}

func (s *AuthService) GetUsersAmount() (int, error) {
	amount, err := s.UserRepository.GetUsersAmount()
	if err != nil {
		s.logger.Error("cant get users amount : " + err.Error())
		return 0, err
	}

	return amount, nil
}

func (s *AuthService) CreateUser(username, password string) error {

	fmt.Println("CreateUser AuthService")
	fmt.Println("username is " + username)
	fmt.Println("password is " + password)

	hashedPassword, err := s.passwordHasherService.HashPassword(password)
	if err != nil {
		s.logger.Error("cant hash password : " + err.Error())
		return fmt.Errorf("cant hash password ")
	}

	err = s.UserRepository.CreateUser(username, hashedPassword)
	if err != nil {
		s.logger.Error("cant create user: " + err.Error())
		return fmt.Errorf("cant create user ")
	}

	return nil
}

func (s *AuthService) DeleteUser(id string) error {

	err := s.UserRepository.DeleteUser(id)
	if err != nil {
		s.logger.Error("cant delete user : " + err.Error())
		return err
	}

	return nil
}

func (s *AuthService) GetUser(id string) (*entity.User, error) {

	user, err := s.UserRepository.GetUserById(id)
	if err != nil {
		s.logger.Error("cant get user : " + err.Error())
		return nil, err
	}

	return user, nil

}

func (s *AuthService) UpdateUser(id, username, password string) error {

	err := s.UserRepository.UpdateUsername(id, username)
	if err != nil {
		s.logger.Error("cant update username : " + err.Error())
		return fmt.Errorf("ant update username")
	}

	if password != "" {
		hashedPassword, err := s.passwordHasherService.HashPassword(password)
		if err != nil {
			s.logger.Error("cant hash password : " + err.Error())
			return fmt.Errorf("cant hash password ")
		}

		err = s.UserRepository.UpdatePassword(id, hashedPassword)
		if err != nil {
			s.logger.Error("cant update username : " + err.Error())
			return fmt.Errorf("cant update password")
		}
	}

	return nil
}
