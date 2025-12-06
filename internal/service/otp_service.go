package service

import (
	"context"
	"onboarding/internal/repository"
	otp "onboarding/internal/repository/otp"
)

type OtpService interface {
	SendOtpForgotPassword(ctx context.Context, email string) error
	VerifyOtpForgotPassword(ctx context.Context, email string, otpCode string) error
}

type IOtpService struct {
	userRepo repository.UserRepository
	otpRepo  otp.OtpRepository
}

func NewOtpService(
	userRepo repository.UserRepository,
	otpRepo otp.OtpRepository,
) OtpService {
	return &IOtpService{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (s *IOtpService) SendOtpForgotPassword(ctx context.Context, email string) error {
	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	return s.otpRepo.SendOtp(ctx, email, otp.ServiceForgotPassword)
}

func (s *IOtpService) VerifyOtpForgotPassword(ctx context.Context, email string, otpCode string) error {
	return s.otpRepo.VerifyOtp(ctx, email, otpCode, otp.ServiceForgotPassword)
}
