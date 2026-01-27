package repositoryAuth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/seva-up/App_Products/config"
	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/models"
)

type authRedisRepository struct {
	client *Client
	cfg    *config.Config
}

func NewAuthRedisRepository(clients *Client, cfg *config.Config) auth.RedisRepository {
	if cfg == nil {
		panic("config cannot be nil")
	}
	if cfg.Jwt == nil {
		panic("jwt config cannot be nil")
	}
	if cfg.Jwt.SecretKey == "" {
		panic("jwt secret cannot be empty")
	}

	return &authRedisRepository{
		client: clients,
		cfg:    cfg,
	}
}

func (r *authRedisRepository) SaveRefreshToken(tokenID string, userID int, email string, ttl time.Duration, metadata map[string]string) error {
	data := models.RefreshTokenData{
		UserID:    userID,
		Email:     email,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}
	if deviceID, ok := metadata["device_id"]; ok {
		data.DeviceID = deviceID
	}
	if userAgent, ok := metadata["user_agent"]; ok {
		data.DeviceID = userAgent
	}
	if ipAddress, ok := metadata["ip_address"]; ok {
		data.DeviceID = ipAddress
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to mashal token data: %w", err)
	}

	key := fmt.Sprintf("refresh:%s", tokenID)

	return r.client.Set(r.client.ctx, key, jsonData, ttl).Err()
}

func (r *authRedisRepository) GetRefreshToken(tokenID string) (*models.RefreshTokenData, error) {
	key := fmt.Sprintf("refresh:%s", tokenID)

	data, err := r.client.Get(r.client.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var tokenData models.RefreshTokenData
	if err = json.Unmarshal(data, &tokenData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}
	return &tokenData, nil
}

func (r *authRedisRepository) DeleteRefreshToken(tokenID string) error {
	key := fmt.Sprintf("refresh:%s", tokenID)
	return r.client.Del(r.client.ctx, key).Err()
}

func (r *authRedisRepository) BlockAccessToken(token string, ttl time.Duration) error {
	tokenHash := uuid.NewSHA1(uuid.Nil, []byte(token)).String()
	key := fmt.Sprintf("blacklist:access:%s", tokenHash)

	return r.client.Set(r.client.ctx, key, "blocked", ttl).Err()
}

func (r *authRedisRepository) IsAccessTokenBlocked(token string) (bool, error) {
	tokenHash := uuid.NewSHA1(uuid.Nil, []byte(token)).String()
	key := fmt.Sprintf("blacklist:access:%s", tokenHash)

	exists, err := r.client.Exists(r.client.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *authRedisRepository) SaveUserSession(userID int, sessionID string, sessionData map[string]interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("session:user:%d:%s", userID, sessionID)

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	return r.client.Set(r.client.ctx, key, jsonData, ttl).Err()
}

func (r *authRedisRepository) GetUserSessions(userID int) ([]map[string]interface{}, error) {
	pattern := fmt.Sprintf("session:user:%d:*", userID)
	var cursor uint64
	var sessions []map[string]interface{}

	for {
		keys, nextCursor, err := r.client.Scan(r.client.ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			data, err := r.client.Get(r.client.ctx, key).Bytes()
			if err != nil {
				continue
			}

			var session map[string]interface{}
			if err := json.Unmarshal(data, &session); err == nil {
				sessions = append(sessions, session)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return sessions, nil
}

func (r *authRedisRepository) DeleteUserSession(userID int, sessionID string) error {
	key := fmt.Sprintf("session:user:%d:%s", userID, sessionID)
	return r.client.Del(r.client.ctx, key).Err()
}

// Tokens

func (r *authRedisRepository) GenerateJWTToken(user *models.User, metadata *models.TokenMetadata) (*models.TokenPair, error) {
	sessionId := uuid.New().String()
	if r.cfg == nil {
		return nil, fmt.Errorf("jwt configuration is not set1")
	}
	if r.cfg.Jwt.SecretKey == "" {
		return nil, fmt.Errorf("jwt configuration is not set2")
	}

	if r.cfg.Jwt.AccessTTL == 0 {
		return nil, fmt.Errorf("access TTL is not configured")
	}
	accessToken, accessClaims, err := r.CreateAccessToken(user, sessionId)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshClaims, err := r.CreateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	tokenMetadata := map[string]string{
		"device_id":  metadata.DeviceID,
		"user_agent": metadata.UserAgent,
		"ip_address": metadata.IPAddress,
	}
	err = r.SaveRefreshToken(refreshClaims.ID, user.ID, user.Email, r.cfg.Jwt.RefreshTTL, tokenMetadata)
	if err != nil {
		return nil, err
	}

	sessionData := map[string]interface{}{
		"session_id":    sessionId,
		"device_info":   metadata,
		"created_at":    time.Now(),
		"last_activity": time.Now(),
	}
	err = r.SaveUserSession(user.ID, sessionId, sessionData, r.cfg.Jwt.RefreshTTL)

	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(r.cfg.Jwt.AccessTTL.Seconds()),
		ExpiresAt:    accessClaims.ExpiresAt.Time,
	}, nil
}
func (r *authRedisRepository) CreateAccessToken(user *models.User, sessionID string) (string, *models.Claims, error) {
	claims := &models.Claims{
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.cfg.Jwt.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    r.cfg.Jwt.Issuer,
			Subject:   user.Email,
		},
	}
	secret := "default-fallback-secret-change-me"
	if r.cfg != nil && r.cfg.Jwt != nil && r.cfg.Jwt.SecretKey != "" {
		secret = r.cfg.Jwt.SecretKey
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, claims, err
}
func (r *authRedisRepository) CreateRefreshToken(user *models.User) (string, *jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.cfg.Jwt.RefreshTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    r.cfg.Jwt.Issuer,
		ID:        uuid.New().String(),
		Subject:   user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(r.cfg.Jwt.SecretKey))

	return tokenString, claims, err
}
func (r *authRedisRepository) ParseToken(tokenString string) (*models.Claims, error) {
	blocked, err := r.IsAccessTokenBlocked(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check token status: %w", err)
	}
	if blocked {
		return nil, errors.New("token has been invalidated")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(r.cfg.Jwt.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")

}
