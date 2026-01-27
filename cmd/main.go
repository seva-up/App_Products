package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seva-up/App_Products/config"
	"github.com/seva-up/App_Products/internal/auth/deliveryAuth/routesAuth"
	"github.com/seva-up/App_Products/internal/auth/repositoryAuth"
	"github.com/seva-up/App_Products/internal/auth/serviceAuth"
	"github.com/seva-up/App_Products/internal/middleware"
)

func main() {
	fmt.Println("Hello,world")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Что то в конфиге пошло не так:", err)
	}

	fmt.Printf("Name:%s,Host:%s,Port:%s", cfg.Db.Name, cfg.Db.Host, cfg.Db.Port)

	dbURL := formatDBURL(cfg)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("ошибка подключения к БД", err)
	}
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatal("БД не отвечает: ", err)
	}

	redisClient, err := repositoryAuth.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatalf("Редис не отвечает на стадии подключения: %v", err)
	}
	defer redisClient.Close()
	log.Println("Подключение к бд установлено")

	redisRepo := repositoryAuth.NewAuthRedisRepository(redisClient, cfg)
	userRepo := repositoryAuth.NewAuthRepository(pool)
	authServ := serviceAuth.NewUserService(userRepo, redisRepo, cfg)
	routerGin := routesAuth.NewGinRouter(authServ)
	routerGin.Use(middleware.AuthMiddleware(redisRepo))
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      routerGin,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("...Сервер запущен на порту %s", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка сервера:", err)
		}
	}()

	<-stop
	log.Println("...Получен сигнал остановки...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Принудительное завершение:", err)
	}

	log.Println("...Сервер остановлен корректно")
}
func formatDBURL(cfg *config.Config) string {
	return "postgresql://" + cfg.Db.User + ":" + cfg.Db.Password +
		"@" + cfg.Db.Host + ":" + cfg.Db.Port + "/" + cfg.Db.Name
}

func formatDBREdis(cfg *config.Config) string {
	redisURL := fmt.Sprintf("redis://%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	return redisURL
}
