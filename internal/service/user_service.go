package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"time"
)

type IUserRepository interface {
	Save(ctx context.Context, userData model.UserData) error
	FindByLogin(ctx context.Context, login string) (*model.UserData, error)
}

type IPasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, plainPassword string) error
}

type UserService struct {
	UserRepository  IUserRepository
	PasswordService IPasswordService
}

func NewUserService(userRepository IUserRepository, passwordService IPasswordService) *UserService {
	return &UserService{
		UserRepository:  userRepository,
		PasswordService: passwordService,
	}
}

func (s *UserService) AddUser(ctx context.Context, userCreateCommand command.UserCreateCommand) (*model.UserData, error) {
	hashedPassword, err := s.PasswordService.Hash(userCreateCommand.Password)
	if err != nil {
		return nil, err
	}

	rawUserData := model.UserData{
		ID:        uuid.New(),
		Login:     userCreateCommand.Login,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := s.UserRepository.Save(ctx, rawUserData); err != nil {
		return nil, err
	}

	return &rawUserData, nil
}

func (s *UserService) FindByLoginAndPassword(ctx context.Context, login string, password string) (*model.UserData, error) {
	userData, err := s.UserRepository.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if err := s.PasswordService.Compare(userData.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return userData, nil
}
