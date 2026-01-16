package auth

import (
	"context"

	"github.com/seva-up/App_Products/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindById(ctx context.Context, userId int) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userId int) error
}
