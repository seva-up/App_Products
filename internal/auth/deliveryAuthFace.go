package auth

import (
	"github.com/gin-gonic/gin"
)

type UserDelivery interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Refresh(c *gin.Context)
	GetSession(c *gin.Context)
}
