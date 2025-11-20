package response

import (
	entity "onboarding/internal/entity"
)

type UserResponse struct {
	Email string `json:"email"`
}

func NewUserResponse(user entity.User) UserResponse {
	return UserResponse{
		Email: user.Email,
	}
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func NewLoginResponse(token string, user entity.User) LoginResponse {
	return LoginResponse{
		Token: token,
		User:  NewUserResponse(user),
	}
}
