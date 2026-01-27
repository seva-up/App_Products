package routesAuth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/auth/deliveryAuth/httpAuth"
)

func NewGinRouter(authService auth.UserService) *gin.Engine {
	router := gin.Default()

	authHandler := httpAuth.NewAuthDelivery(authService)

	publicApi := router.Group("/api/v1")
	{
		publicApi.POST("/register", authHandler.Register)
		publicApi.POST("/login", authHandler.Login)
		publicApi.POST("/logout", authHandler.Logout)
		publicApi.POST("/refresh", authHandler.Refresh)
		publicApi.GET("/session", authHandler.GetSession)
		publicApi.GET("/health", healthCheck)
	}

	return router
}
func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
