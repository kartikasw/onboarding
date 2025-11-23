package service

import (
	"context"
	"onboarding/internal/entity"
	"onboarding/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (entity.UserViewModel, error)
}

type IUserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &IUserService{
		userRepo: userRepo,
	}
}

func (s *IUserService) GetUser(ctx context.Context, uuid uuid.UUID) (entity.UserViewModel, error) {
	result := entity.UserViewModel{}

	user, err := s.userRepo.GetUserByUUID(ctx, uuid)
	if err != nil {
		return result, err
	}

	return user.ToViewModel(), err
}
