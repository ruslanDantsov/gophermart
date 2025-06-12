package service

import (
	"context"
	"errors"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, userData model.UserData) error {
	args := m.Called(ctx, userData)
	return args.Error(0)
}

func (m *MockUserRepository) FindByLogin(ctx context.Context, login string) (*model.UserData, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserData), args.Error(1)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) Compare(hashedPassword, plainPassword string) error {
	args := m.Called(hashedPassword, plainPassword)
	return args.Error(0)
}

// --- Tests ---

func TestUserService_AddUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	passwordService := new(MockPasswordService)

	passwordService.On("Hash", "password123").Return("hashed123", nil)

	repo.On("Save", mock.Anything, mock.MatchedBy(func(user model.UserData) bool {
		return user.Login == "testuser" && user.Password == "hashed123"
	})).Return(nil)

	userService := NewUserService(repo, passwordService)

	cmd := command.UserCreateCommand{
		Login:    "testuser",
		Password: "password123",
	}

	user, err := userService.AddUser(ctx, cmd)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, "hashed123", user.Password)

	repo.AssertExpectations(t)
	passwordService.AssertExpectations(t)
}

func TestUserService_AddUser_HashFails(t *testing.T) {
	ctx := context.Background()
	passwordService := new(MockPasswordService)
	passwordService.On("Hash", "password123").Return("", errors.New("hash error"))

	userService := NewUserService(nil, passwordService)

	cmd := command.UserCreateCommand{
		Login:    "testuser",
		Password: "password123",
	}

	user, err := userService.AddUser(ctx, cmd)

	assert.Nil(t, user)
	assert.EqualError(t, err, "hash error")

	passwordService.AssertExpectations(t)
}

func TestUserService_FindByLoginAndPassword_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	passwordService := new(MockPasswordService)

	user := &model.UserData{
		Id:        uuid.New(),
		Login:     "testuser",
		Password:  "hashed123",
		CreatedAt: time.Now(),
	}

	repo.On("FindByLogin", ctx, "testuser").Return(user, nil)
	passwordService.On("Compare", "hashed123", "password123").Return(nil)

	userService := NewUserService(repo, passwordService)

	result, err := userService.FindByLoginAndPassword(ctx, "testuser", "password123")

	assert.NoError(t, err)
	assert.Equal(t, user, result)

	repo.AssertExpectations(t)
	passwordService.AssertExpectations(t)
}

func TestUserService_FindByLoginAndPassword_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)

	repo.On("FindByLogin", ctx, "testuser").Return(nil, errors.New("not found"))

	userService := NewUserService(repo, nil)

	result, err := userService.FindByLoginAndPassword(ctx, "testuser", "any")

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")

	repo.AssertExpectations(t)
}

func TestUserService_FindByLoginAndPassword_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	passwordService := new(MockPasswordService)

	user := &model.UserData{
		Id:        uuid.New(),
		Login:     "testuser",
		Password:  "hashed123",
		CreatedAt: time.Now(),
	}

	repo.On("FindByLogin", ctx, "testuser").Return(user, nil)
	passwordService.On("Compare", "hashed123", "wrongpass").Return(errors.New("invalid"))

	userService := NewUserService(repo, passwordService)

	result, err := userService.FindByLoginAndPassword(ctx, "testuser", "wrongpass")

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid credentials")

	repo.AssertExpectations(t)
	passwordService.AssertExpectations(t)
}
