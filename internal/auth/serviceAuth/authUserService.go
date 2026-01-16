package serviceAuth

import (
	"context"
	"fmt"

	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/models"
)

type authUS struct {
	userRepo auth.UserRepository
}

func NewUserService(userRepo auth.UserRepository) auth.UserService {
	return &authUS{userRepo: userRepo}
}

func (a *authUS) Register(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "AuthServiceUser.Register"

	result, err := a.userRepo.FindById(ctx, user.ID)
	if result != nil {
		return nil, fmt.Errorf("такой аккунт уже существует")
	}

	resultate, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать аккаунт %s", err)
	}

	return resultate, nil
}
