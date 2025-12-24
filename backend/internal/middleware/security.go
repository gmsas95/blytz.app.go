package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders middleware adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Content Security Policy
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://www.googletagmanager.com; " +
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
			"font-src 'self' https://fonts.gstatic.com; " +
			"img-src 'self' data: https:; " +
			"connect-src 'self' wss: ws:; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.Header("Content-Security-Policy", csp)
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Permissions Policy
		permissionsPolicy := "camera=(), microphone=(), geolocation=(), " +
			"payment=(), usb=(), magnetometer=(), gyroscope=()"
		c.Header("Permissions-Policy", permissionsPolicy)
		
		// Strict Transport Security (only on HTTPS)
		if c.Request.TLS != nil || strings.HasPrefix(c.Request.Header.Get("X-Forwarded-Proto"), "https") {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		
		// Hide server information
		c.Header("Server", "")
		
		c.Next()
	}
}

// HTTPSRedirect middleware redirects HTTP to HTTPS in production
func HTTPSRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only redirect in production
		if gin.Mode() == gin.ReleaseMode {
			// Check if the request is already HTTPS
			if c.Request.TLS == nil && !strings.HasPrefix(c.Request.Header.Get("X-Forwarded-Proto"), "https") {
				// Get the host
				host := c.Request.Host
				// Redirect to HTTPS
				httpsURL := "https://" + host + c.Request.RequestURI
				c.Redirect(http.StatusMovedPermanently, httpsURL)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// APISecurity middleware adds API-specific security headers
func APISecurity() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rate limit headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", "99")
		c.Header("X-RateLimit-Reset", "3600")
		
		// API version
		c.Header("API-Version", "v1")
		
		// Cache control for API responses
		if c.Request.Method == "GET" {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		
		c.Next()
	}
}