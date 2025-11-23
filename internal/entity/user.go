package entity

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserViewModel struct {
	UUID  uuid.UUID
	Email string `json:"email"`
}

func (e User) ToViewModel() UserViewModel {
	return UserViewModel{
		UUID:  e.UUID,
		Email: e.Email,
	}
}
