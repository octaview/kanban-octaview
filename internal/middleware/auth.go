package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/service"
)

type AuthMiddleware struct {
	authService service.AuthServiceInterface
}

func NewAuthMiddleware(authService service.AuthServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := m.authService.ParseToken(token)
		if err != nil {
			status := http.StatusUnauthorized
			errMsg := "Invalid or expired token"

			if _, ok := err.(*models.AuthError); ok {
				errMsg = err.Error()
			}

			c.JSON(status, gin.H{
				"error": errMsg,
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		userID, err := m.authService.ParseToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}