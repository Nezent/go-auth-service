package services

import (
	"github.com/Nezent/auth-service/internal/application/dto"
	"github.com/Nezent/auth-service/internal/domain/shared"
	"github.com/Nezent/auth-service/internal/domain/user"
	"github.com/Nezent/auth-service/internal/infrastructure/crypto"
	"github.com/Nezent/auth-service/internal/infrastructure/repository"
)

type UserService interface {
	CreateUser(email, password string) (*dto.CreateUserResponse, *shared.DomainError)
}

var _ UserService = (*UserServiceImpl)(nil)

type UserServiceImpl struct {
	userRepo *repository.UserRepositoryImpl
}

func NewUserService(userRepo *repository.UserRepositoryImpl) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// Add methods for user operations, e.g., CreateUser, GetUserByID, etc.

func (s *UserServiceImpl) CreateUser(email, password string) (*dto.CreateUserResponse, *shared.DomainError) {
	hash, salt, hashErr := crypto.HashPassword(password)
	if hashErr != nil {
		return nil, shared.NewDomainError("HASHING_FAILED", 500, hashErr.Error())
	}
	var user = &user.User{
		Email:    email,
		Password: hash,
		Salt:     salt,
	}

	id, repoErr := s.userRepo.Create(user)
	if repoErr != nil {
		return nil, repoErr
	}
	return &dto.CreateUserResponse{ID: id, Email: email}, nil
}
