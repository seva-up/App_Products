package auth

import (
	"context"

	"github.com/seva-up/App_Products/internal/models"
)

type UserService interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	//Login(ctx context.Context)error
	//Logout(ctx context.Context)error
}
