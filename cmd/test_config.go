package main

import (
	"fmt"
	"log"
	"os"

	"github.com/seva-up/App_Products/config"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("=== Тест конфигурации Viper ===")

	// Проверка переменных окружения
	fmt.Println("\n1. Проверка переменных окружения:")
	envVars := []string{
		"APP_PORT",
		"JWT_SECRET",
		"REDIS_HOST",
		"DB_HOST",
	}

	for _, env := range envVars {
		val := os.Getenv(env)
		if val == "" {
			fmt.Printf("   ⚠️  %s: не установлен\n", env)
		} else {
			if env == "JWT_SECRET" {
				fmt.Printf("   ✅ %s: %s\n", env, maskSecret(val))
			} else {
				fmt.Printf("   ✅ %s: %s\n", env, val)
			}
		}
	}

	// Проверка наличия конфиг файла
	fmt.Println("\n2. Поиск конфиг файлов:")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	paths := []string{".", "./config", "/etc/app/"}
	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("   ⚠️  Конфиг файл не найден")
		} else {
			fmt.Printf("   ❌ Ошибка чтения конфига: %v\n", err)
		}
	} else {
		fmt.Printf("   ✅ Конфиг найден: %s\n", viper.ConfigFileUsed())
	}

	// Загрузка конфигурации
	fmt.Println("\n3. Загрузка конфигурации...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Ошибка загрузки конфига: %v", err)
	}

	// Вывод конфигурации
	fmt.Println("\n4. Загруженная конфигурация:")
	fmt.Printf("   App Port: %s\n", cfg.App.Port)
	fmt.Printf("   App Env: %s\n", cfg.App.Env)
	fmt.Printf("   JWT Secret: %s\n", maskSecret(cfg.Jwt.SecretKey))
	fmt.Printf("   JWT Access TTL: %v\n", cfg.Jwt.AccessTTL)
	fmt.Printf("   JWT Refresh TTL: %v\n", cfg.Jwt.RefreshTTL)
	fmt.Printf("   Redis Host: %s\n", cfg.Redis.Host)
	fmt.Printf("   Redis Port: %s\n", cfg.Redis.Port)
	fmt.Printf("   DB Host: %s\n", cfg.Db.Host)
	fmt.Printf("   DB Port: %s\n", cfg.Db.Port)

	fmt.Println("\n✅ Все проверки пройдены!")
}

func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "***" + secret[len(secret)-4:]
}
