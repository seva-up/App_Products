package repositoryAuth

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/seva-up/App_Products/config"
)

type Client struct {
	*redis.Client
	ctx context.Context
}

var (
	redisClient *Client
)

func NewRedisClient(cfg *config.Redis) (*Client, error) {
	opts, err := redis.ParseURL("redis://localhost:6379") //поменять на url а то какой хост нах
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Проверяем подключение
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	redisClient = &Client{
		Client: client,
		ctx:    ctx,
	}

	log.Println("Redis connected successfully")
	return redisClient, nil
}

func GetRedisClient() *Client {
	return redisClient
}

// HealthCheck проверяет состояние Redis
func (c *Client) HealthCheck() error {
	return c.Ping(c.ctx).Err()
}

// Close закрывает подключение к Redis
func (c *Client) Close() error {
	return c.Client.Close()
}
