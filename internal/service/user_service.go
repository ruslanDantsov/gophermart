package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IUserRepository interface {
	Save(ctx context.Context, userData model.UserData) error
}

type UserService struct {
	userRepository IUserRepository
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) AddUser(ctx context.Context, userCreateCommand command.UserCreateCommand) (*model.UserData, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreateCommand.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rawUserData := model.UserData{
		Id:        uuid.New(),
		Login:     userCreateCommand.Login,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := s.userRepository.Save(ctx, rawUserData); err != nil {
		return nil, err
	}

	return &rawUserData, nil
}
