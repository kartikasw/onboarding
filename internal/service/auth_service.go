package service

import (
	"context"
	"fmt"
	"onboarding/common"
	"onboarding/internal/entity"
	"onboarding/internal/repository"
	pw "onboarding/pkg/password"
	"onboarding/pkg/token"
)

type AuthService interface {
	Register(ctx context.Context, email string, password string) (entity.UserViewModel, error)
	Login(ctx context.Context, email string, password string) (*token.JWTToken, error)
}

type IAuthService struct {
	userRepo repository.UserRepository
	jwtImpl  token.JWT
}

func NewAuthService(userRepo repository.UserRepository, jwtImpl token.JWT) AuthService {
	return &IAuthService{userRepo: userRepo, jwtImpl: jwtImpl}
}

func (s *IAuthService) Register(
	ctx context.Context,
	email string,
	password string,
) (entity.UserViewModel, error) {
	hashedPassword, err := pw.HashPassword(password)
	if err != nil {
		return entity.UserViewModel{}, err
	}

	arg := entity.User{
		Email:    email,
		Password: hashedPassword,
	}

	err = s.userRepo.CreateUser(ctx, arg)
	if err != nil {
		return entity.UserViewModel{}, err
	}

	return arg.ToViewModel(), nil
}

func (s *IAuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*token.JWTToken, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := pw.CheckPassword(password, user.Password); err != nil {
		return nil, fmt.Errorf("%d", common.ErrCredentiials)
	}

	jwtToken, err := s.jwtImpl.CreateAccessToken(user.UUID)
	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}
