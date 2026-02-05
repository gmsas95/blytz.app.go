package middleware

import (
	"net/http"
	"strings"

	"github.com/blytz/live/backend/internal/domain/user"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(tokenManager user.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "UNAUTHORIZED",
				"message": "authorization header required",
			})
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "UNAUTHORIZED",
				"message": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]
		claims, err := tokenManager.Validate(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "UNAUTHORIZED",
				"message": "invalid or expired token",
			})
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID.String())
		c.Set("user_email", claims.Email)
		c.Set("user_role", string(claims.Role))

		c.Next()
	}
}

// RequireRole creates role-based authorization middleware
func RequireRole(roles ...user.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "UNAUTHORIZED",
				"message": "authentication required",
			})
			return
		}

		userRole := user.Role(roleStr.(string))
		
		// Check if user has required role
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "FORBIDDEN",
			"message": "insufficient permissions",
		})
	}
}

// OptionalAuth creates optional authentication middleware
func OptionalAuth(tokenManager user.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := tokenManager.Validate(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("user_email", claims.Email)
		c.Set("user_role", string(claims.Role))

		c.Next()
	}
}