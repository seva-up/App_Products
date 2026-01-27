package httpAuth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/auth/deliveryAuth/dtoAuth"
	"github.com/seva-up/App_Products/internal/models"
)

type authDelivery struct {
	authUS auth.UserService
}

func NewAuthDelivery(authUS auth.UserService) auth.UserDelivery {
	return &authDelivery{authUS: authUS}
}

func (h *authDelivery) Register(c *gin.Context) {
	var req dtoAuth.InRegisters

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.authUS.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *authDelivery) Login(c *gin.Context) {
	var req dtoAuth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	metadata := extractMetadata(c.Request)
	tokens, err := h.authUS.Login(c.Request.Context(), &req, metadata)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/v1",
		HttpOnly: true,
		Secure:   c.Request.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(7 * 24 * 60 * 60),
	})
	c.JSON(http.StatusOK, tokens)
}
func (h *authDelivery) Logout(c *gin.Context) {
	accessToken := extractToken(c.Request)
	refreshToken := ""
	if cookie, err := c.Request.Cookie("refresh_token"); err == nil {
		refreshToken = cookie.Value
	}

	userID := c.GetInt("user_id")
	err := h.authUS.Logout(c.Request.Context(), accessToken, refreshToken, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1",
		HttpOnly: true,
		Secure:   c.Request.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	c.JSON(http.StatusOK, "Successfully logged out")

}

func (h *authDelivery) Refresh(c *gin.Context) {
	var refreshToken string

	if cookie, err := c.Request.Cookie("refresh_token"); err == nil {
		refreshToken = cookie.Value
	} else {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err = c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		refreshToken = req.RefreshToken
	}
	metadata := extractMetadata(c.Request)

	tokens, err := h.authUS.RefreshTokens(c.Request.Context(), refreshToken, metadata)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/v1",
		HttpOnly: true,
		Secure:   c.Request.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(7 * 24 * 60 * 60),
	})

	c.JSON(http.StatusOK, tokens)
}

func (h *authDelivery) GetSession(c *gin.Context) {
	userID := c.GetInt("user_id")

	session, err := h.authUS.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

func extractMetadata(r *http.Request) *models.TokenMetadata {
	return &models.TokenMetadata{
		DeviceID:  r.Header.Get("X-Device-ID"),
		UserAgent: r.Header.Get("User-Agent"),
		IPAddress: getIPAddress(r),
	}
}
func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
	}
	return ip
}
