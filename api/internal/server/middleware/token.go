package firestore_middleware

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken creates a JWT validation middleware with configurable options
func ValidateToken(config domain.MiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token
		tokenString, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid authorization token",
			})
			return
		}

		// Parse and validate JWT
		claims, err := parseAndValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			return
		}

		// Set claims in context (always available)
		c.Set("claims", claims)
		c.Set("user_id", claims.ID)
		c.Set("user_email", claims.Email)

		// Skip user lookup if not needed (for performance)
		if config.SkipUserLookup {
			c.Next()
			return
		}

		// Optional: Check cache first
		var user *domain.User
		cacheKey := fmt.Sprintf("user:%s", claims.ID)

		if config.Cache != nil {
			if cached, found := config.Cache.Get(cacheKey); found {
				if cachedUser, ok := cached.(*domain.User); ok {
					user = cachedUser
				}
			}
		}

		// Fetch user if not in cache
		if user == nil {
			user, err = config.UserService.GetUserByID(c.Request.Context(), claims.ID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "User verification failed",
				})
				return
			}

			if user == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "User not found",
				})
				return
			}

			// Cache user if cache is available
			if config.Cache != nil {
				config.Cache.Set(cacheKey, user, config.CacheDuration)
			}
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}

// extractToken extracts JWT token from cookie or Authorization header
func extractToken(c *gin.Context) (string, error) {
	// Try cookie first
	if tokenString, err := c.Cookie("Authorization"); err == nil && tokenString != "" {
		return tokenString, nil
	}

	// Try Authorization header as fallback
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization token provided")
	}

	// Handle "Bearer <token>" format
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return authHeader, nil
}

// parseAndValidateJWT parses and validates the JWT token
func parseAndValidateJWT(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get signing key (consider using a key rotation system)
		key := cmn.GetEnvVariable("JWT_SIGN_KEY")
		if key == "" {
			return nil, fmt.Errorf("JWT signing key not configured")
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate required fields
	if claims.ID == "" {
		return nil, fmt.Errorf("user ID not found in token")
	}
	if claims.Email == "" {
		return nil, fmt.Errorf("email not found in token")
	}

	return claims, nil
}

func AuthenticateOnly() gin.HandlerFunc {
	return ValidateToken(domain.MiddlewareConfig{
		SkipUserLookup: true,
	})
}

func AuthenticateAndLoadUser(userService domain.UserService, cache domain.TokenCache) gin.HandlerFunc {
	return ValidateToken(domain.MiddlewareConfig{
		UserService:    userService,
		Cache:          cache,
		CacheDuration:  15 * time.Minute,
		SkipUserLookup: false,
	})
}
