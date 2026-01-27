package auth

import (
	"context"

	"github.com/seva-up/App_Products/internal/auth/deliveryAuth/dtoAuth"
	"github.com/seva-up/App_Products/internal/models"
)

type UserService interface {
	Register(ctx context.Context, user *dtoAuth.InRegisters) (*models.User, error)
	Login(ctx context.Context, req *dtoAuth.LoginRequest, metadata *models.TokenMetadata) (*models.TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string, userID int) error
	RefreshTokens(ctx context.Context, refreshToken string, metadata *models.TokenMetadata) (*models.TokenPair, error)
	ValidateToken(ctx context.Context, token string) (*models.Claims, error)
	GetUserSessions(ctx context.Context, userID int) ([]map[string]interface{}, error)
}
