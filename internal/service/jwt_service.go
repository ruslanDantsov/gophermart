package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTService struct {
	JWTSecret string
}

type TokenResult struct {
	AccessToken string
	ExpiresIn   int64
}

func NewAuthService(secret string) *JWTService {
	return &JWTService{JWTSecret: secret}
}

func (s *JWTService) GenerateJWT(id uuid.UUID, username string) (*TokenResult, error) {
	expirationTime := time.Now().Add(time.Hour * 1).Unix()
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResult{
		AccessToken: signed,
		ExpiresIn:   expirationTime - time.Now().Unix(),
	}, nil
}
