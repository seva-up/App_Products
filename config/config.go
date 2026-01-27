package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App   App        `mapstructure:"app"`
	Db    DbPostgres `mapstructure:"db"`
	Jwt   *Jwt       `mapstructure:"jwt"`
	Redis Redis      `mapstructure:"redis"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
	Port    string `mapstructure:"port"`
	Host    string `mapstructure:"host"`
}
type DbPostgres struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

type Jwt struct {
	SecretKey    string        `mapstructure:"secret_key"`
	AccessToken  string        `mapstructure:"access_token"`
	RefreshToken string        `mapstructure:"refresh_token"`
	Issuer       string        `mapstructure:"issuer"`
	AccessTTL    time.Duration `mapstructure:"access_ttl"`
	RefreshTTL   time.Duration `mapstructure:"refresh_ttl"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

func Load() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	bindEnvVariables()
	viper.AutomaticEnv()
	viper.SetDefault("jwt.access_ttl", "15m")
	viper.SetDefault("jwt.refresh_ttl", "168h")
	viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")
	viper.BindEnv("jwt.access_ttl", "JWT_ACCESS_TTL")
	viper.BindEnv("jwt.refresh_ttl", "JWT_REFRESH_TTL")
	viper.SetDefault("host", "localhost")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Конфиг не найден в функции config.go/Load()")
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Конфиг не смог за анмаршалиться")
	}

	setDefaults(&config)

	// Проверка
	if config.Jwt.SecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET must be set in config file or environment variable JWT_SECRET_KEY")
	}

	return &config, nil
}
func bindEnvVariables() {
	// Связываем переменные окружения с полями конфига

	// App
	viper.BindEnv("app.port", "APP_PORT")
	viper.BindEnv("app.env", "APP_ENV")

	// Database
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.dbname", "DB_NAME")

	// Redis
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	// JWT
	viper.BindEnv("jwt.secret_key", "JWT_SECRET")
	viper.BindEnv("jwt.access_ttl", "JWT_ACCESS_TTL")
	viper.BindEnv("jwt.refresh_ttl", "JWT_REFRESH_TTL")
}

func setDefaults(config *Config) {
	// App defaults
	if config.App.Port == "" {
		config.App.Port = "8080"
	}
	if config.App.Env == "" {
		config.App.Env = "development"
	}

	// Database defaults
	if config.Db.Host == "" {
		config.Db.Host = "localhost"
	}
	if config.Db.Port == "" {
		config.Db.Port = "5432"
	}
	// Redis defaults
	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
	}
	if config.Redis.Port == "" {
		config.Redis.Port = "6379"
	}

	// JWT defaults
	if config.Jwt.AccessTTL == 0 {
		config.Jwt.AccessTTL = 15 * time.Minute
	}
	if config.Jwt.RefreshTTL == 0 {
		config.Jwt.RefreshTTL = 168 * time.Hour // 7 дней
	}

	// Если JWT секрет не установлен, используем дефолтный для разработки
	if config.Jwt.SecretKey == "" && config.App.Env == "development" {
		config.Jwt.SecretKey = "dev-secret-key-change-in-production"
		log.Println("⚠️  WARNING: Using default JWT secret for development")
	}
}
