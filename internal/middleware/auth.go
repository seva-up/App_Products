package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seva-up/App_Products/internal/auth"
)

func AuthMiddleware(authService auth.RedisRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/v1/login" ||
			c.Request.URL.Path == "/api/v1/refresh" ||
			c.Request.URL.Path == "/api/v1/register" {
			c.Next()
			return
		}

		token := extractToken(c.Request)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := authService.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контексте
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("session_id", claims.SessionID)

		c.Next()
	}
}

func extractToken(r *http.Request) string {
	// 1. Пробуем из заголовка Authorization
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:]
	}

	// 2. Пробуем из query параметра
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	// 3. Пробуем из cookie
	cookie, err := r.Cookie("access_token")
	if err == nil {
		return cookie.Value
	}

	return ""
}
