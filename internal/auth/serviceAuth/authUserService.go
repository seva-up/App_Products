package serviceAuth

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/seva-up/App_Products/config"
	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/auth/deliveryAuth/dtoAuth"
	"github.com/seva-up/App_Products/internal/models"
)

type authUS struct {
	userRepo  auth.UserRepository
	redisRepo auth.RedisRepository
	cfg       *config.Config
}

func NewUserService(userRepo auth.UserRepository, redisRepo auth.RedisRepository, cfg *config.Config) auth.UserService {
	return &authUS{userRepo: userRepo, redisRepo: redisRepo, cfg: cfg}
}

func (a *authUS) Register(ctx context.Context, user *dtoAuth.InRegisters) (*models.User, error) {
	const op = "AuthServiceUser.Register"

	result, err := a.userRepo.FindByEmail(ctx, user.Email)
	if err == nil && result != nil {
		return nil, errors.New("user already exists")
	}
	user1 := &models.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
	}
	resultate, err := a.userRepo.Create(ctx, user1)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать аккаунт %s", err)
	}

	return resultate, nil
}

func (a *authUS) Login(ctx context.Context, req *dtoAuth.LoginRequest, metadata *models.TokenMetadata) (*models.TokenPair, error) {
	user, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("такого пользователя не существует")
		}
		return nil, fmt.Errorf("не правильная записись: %w", err)
	}

	tokens, err := a.redisRepo.GenerateJWTToken(user, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil

}

func (a *authUS) Logout(ctx context.Context, accessToken, refreshToken string, userID int) error {
	claims, err := a.redisRepo.ParseToken(accessToken)
	if err == nil && claims.SessionID != "" {
		// Удаляем сессию
		err = a.redisRepo.DeleteUserSession(userID, claims.SessionID)

		// Блокируем access token
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			err = a.redisRepo.BlockAccessToken(accessToken, ttl)
		}
	}
	if refreshToken != "" {
		token, err := jwt2.Parse(refreshToken, func(token *jwt2.Token) (interface{}, error) {
			return []byte(a.cfg.Jwt.SecretKey), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt2.MapClaims); ok {
				if tokenID, ok := claims["jti"].(string); ok {
					a.redisRepo.DeleteRefreshToken(tokenID)
				}
			}
		}
	}

	return nil
}
func (a *authUS) RefreshTokens(ctx context.Context, refreshToken string, metadata *models.TokenMetadata) (*models.TokenPair, error) {
	// Парсим refresh token
	token, err := jwt2.Parse(refreshToken, func(token *jwt2.Token) (interface{}, error) {
		return []byte(a.cfg.Jwt.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt2.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Получаем ID refresh token
	tokenID, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("invalid token ID")
	}

	// Проверяем существование refresh token в Redis
	tokenData, err := a.redisRepo.GetRefreshToken(tokenID)
	if err != nil || tokenData == nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// Находим пользователя
	user, err := a.userRepo.FindById(ctx, tokenData.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Генерируем новую пару токенов
	newTokens, err := a.redisRepo.GenerateJWTToken(user, metadata)
	if err != nil {
		return nil, err
	}

	// Удаляем старый refresh token
	err = a.redisRepo.DeleteRefreshToken(tokenID)
	if err != nil {
		fmt.Println("service.RefreshToken.delete token invalid")
	}
	return newTokens, nil
}
func (a *authUS) ValidateToken(ctx context.Context, token string) (*models.Claims, error) {
	return a.redisRepo.ParseToken(token)
}

func (a *authUS) GetUserSessions(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	return a.redisRepo.GetUserSessions(userID)
}
