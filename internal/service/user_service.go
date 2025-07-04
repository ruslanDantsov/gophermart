package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"time"
)

type UserRepository interface {
	Save(ctx context.Context, userData entity.UserData) error
	FindByLogin(ctx context.Context, login string) (*entity.UserData, error)
}

type PasswordManager interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, plainPassword string) error
}

type UserService struct {
	userRepository  UserRepository
	passwordManager PasswordManager
}

func NewUserService(userRepository UserRepository, passwordManager PasswordManager) *UserService {
	return &UserService{
		userRepository:  userRepository,
		passwordManager: passwordManager,
	}
}

func (s *UserService) AddUser(ctx context.Context, userCreateCommand command.UserCreateCommand) (*entity.UserData, error) {
	hashedPassword, err := s.passwordManager.Hash(userCreateCommand.Password)
	if err != nil {
		return nil, err
	}

	rawUserData := entity.UserData{
		ID:        uuid.New(),
		Login:     userCreateCommand.Login,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := s.userRepository.Save(ctx, rawUserData); err != nil {
		return nil, err
	}

	return &rawUserData, nil
}

func (s *UserService) FindByLoginAndPassword(ctx context.Context, login string, password string) (*entity.UserData, error) {
	userData, err := s.userRepository.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if err := s.passwordManager.Compare(userData.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return userData, nil
}
