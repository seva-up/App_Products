package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/seva-up/App_Products/internal/models"
)

type RedisRepository interface {
	SaveRefreshToken(tokenID string, userID int, email string, ttl time.Duration, metadata map[string]string) error
	GetRefreshToken(tokenID string) (*models.RefreshTokenData, error)
	DeleteRefreshToken(tokenID string) error
	BlockAccessToken(token string, ttl time.Duration) error
	IsAccessTokenBlocked(token string) (bool, error)

	SaveUserSession(userID int, sessionID string, sessionData map[string]interface{}, ttl time.Duration) error
	GetUserSessions(userID int) ([]map[string]interface{}, error)
	DeleteUserSession(userID int, sessionID string) error

	GenerateJWTToken(user *models.User, metadata *models.TokenMetadata) (*models.TokenPair, error)
	CreateAccessToken(user *models.User, sessionID string) (string, *models.Claims, error)
	CreateRefreshToken(user *models.User) (string, *jwt.RegisteredClaims, error)
	ParseToken(tokenString string) (*models.Claims, error)

	//SaveLoginAttempt(ctx context.Context) error
	//GetFailedLoginCount(ctx context.Context) error
}
