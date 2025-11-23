package response

import (
	entity "onboarding/internal/entity"
)

type UserResponse struct {
	Email string `json:"email"`
}

func NewUserResponse(user entity.UserViewModel) UserResponse {
	return UserResponse{
		Email: user.Email,
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewLoginResponse(token string) LoginResponse {
	return LoginResponse{
		Token: token,
	}
}
