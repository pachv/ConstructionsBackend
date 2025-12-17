package services

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword = errors.New("password is empty")

	ErrInvalidHash = errors.New("invalid password hash")
)

type PasswordService struct {
	cost int
}

func NewPasswordService(cost int) *PasswordService {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordService{cost: cost}
}

func (s *PasswordService) Hash(password string) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", ErrEmptyPassword
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func (s *PasswordService) Compare(hash string, password string) (bool, error) {
	if hash == "" {
		return false, ErrInvalidHash
	}
	if strings.TrimSpace(password) == "" {
		return false, ErrEmptyPassword
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, ErrInvalidHash
}
