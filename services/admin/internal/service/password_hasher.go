package service

import "golang.org/x/crypto/bcrypt"

type PasswordHasherService struct {
}

func NewPasswordHasherService() *PasswordHasherService {
	return &PasswordHasherService{}
}

func (s *PasswordHasherService) HashPassword(password string) (hashedPassword string, err error) {

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func (s *PasswordHasherService) IsPasswordCorrect(password, hashedPassword string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, nil
}
