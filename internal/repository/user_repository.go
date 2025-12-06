package repository

import (
	"context"
	"onboarding/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	UpdateUserPassword(ctx context.Context, email, newPassword string) error
}

type IUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &IUserRepository{db: db}
}

func (r *IUserRepository) CreateUser(ctx context.Context, user entity.User) error {
	return r.db.WithContext(ctx).Omit("UUID").Create(user).Error
}

func (r *IUserRepository) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Take(&user, "uuid = ?", uuid).Error

	return user, err
}

func (r *IUserRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Order(nil).Take(&user, "email = ?", email).Error

	return user, err
}

func (r *IUserRepository) UpdateUserPassword(ctx context.Context, email, newPassword string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Take(&entity.User{}, "email = ?", email).Error; err != nil {
			return err
		}

		return tx.
			Model(&entity.User{}).
			Where("email = ?", email).
			Update("password", newPassword).Error
	})
}
