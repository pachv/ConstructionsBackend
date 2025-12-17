package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	refreshTokenSecret string
}

func NewTokenService(refreshTokenSecret string) *TokenService {
	return &TokenService{
		refreshTokenSecret: refreshTokenSecret,
	}
}

func (s *TokenService) CreateRefreshToken(id, login string) (token string, err error) {
	refreshClaims := jwt.MapClaims{
		"id":    id,
		"login": login,
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	return refresh.SignedString([]byte(s.refreshTokenSecret))

}

func (s *TokenService) IsRefreshTokenCorrect(refreshToken string) (bool, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.refreshTokenSecret), nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, nil
	}

	// Проверяем что токен не истек
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("invalid claims format")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return false, nil
		}
	} else {
		return false, errors.New("expiration time not found")
	}

	return true, nil
}

func (s *TokenService) GetRefreshTokenClaims(refreshToken string) (id, login string, err error) {
	claims, err := s.parseAndValidateToken(refreshToken, s.refreshTokenSecret)
	if err != nil {
		return "", "", err
	}

	id, err = s.extractStringClaim(claims, "id")
	if err != nil {
		return "", "", err
	}

	login, err = s.extractStringClaim(claims, "login")
	if err != nil {
		return "", "", err
	}

	return id, login, nil
}

func (s *TokenService) parseAndValidateToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	return claims, nil
}

func (s *TokenService) extractStringClaim(claims jwt.MapClaims, key string) (string, error) {
	valueInterface, exists := claims[key]
	if !exists {
		return "", errors.New(key + " not found in token claims")
	}

	value, ok := valueInterface.(string)
	if !ok {
		return "", errors.New(key + " has invalid type")
	}

	return value, nil
}

func (s *TokenService) extractInt64Claim(claims jwt.MapClaims, key string) (int64, error) {
	valueInterface, exists := claims[key]
	if !exists {
		return 0, errors.New(key + " not found in token claims")
	}

	valueFloat, ok := valueInterface.(float64)
	if !ok {
		return 0, errors.New(key + " has invalid type")
	}

	return int64(valueFloat), nil
}
