package service

import (
	"context"
	"onboarding/internal/entity"
)

type IAuthService struct {
}

type AuthService interface {
	Register(ctx context.Context, admin entity.User) (entity.User, error)
	Login(ctx context.Context, admin entity.User) (entity.User, string, error)
}

func NewAuthService() AuthService {
	return &IAuthService{}
}

func (s *IAuthService) Register(ctx context.Context, user entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (s *IAuthService) Login(ctx context.Context, user entity.User) (entity.User, string, error) {
	return entity.User{}, "", nil
}
