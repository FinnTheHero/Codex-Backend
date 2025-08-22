package token

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/service"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetClaimsFromToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/api/user/refresh" || path == "/api/user/logout" || path == "/api/user/login" {
			c.Next()
			return
		}

		tokenString, err := ExtractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid authorization token",
			})
			return
		}

		// Parse and validate JWT
		claims, err := ParseAndValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			return
		}

		// Set claims in context (always available)
		c.Set("claims", claims)
	}
}

// ValidateToken creates a JWT validation middleware with configurable options
func LookupUser(config domain.LookupUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.Next()
			return
		}

		userClaims, ok := claims.(*domain.Claims)
		if !ok {
			c.Next()
			return
		}

		// Check cache first
		var user *domain.User
		cacheKey := fmt.Sprintf("user:%s", userClaims.ID)

		if config.Cache != nil {
			if cached, found := config.Cache.Get(cacheKey); found {
				if cachedUser, ok := cached.(*domain.User); ok {
					user = cachedUser
				}
			}
		}

		// Fetch user if not in cache
		if user == nil {
			user, err := service.GetUserByID(userClaims.ID, c.Request.Context())
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
func ExtractToken(c *gin.Context) (string, error) {
	// Try cookie first
	if tokenString, err := c.Cookie("access_token"); err == nil && tokenString != "" {
		return tokenString, nil
	}

	// Try Authorization header as fallback
	authHeader := c.GetHeader("access_token")
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
func ParseAndValidateJWT(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := cmn.GetEnvVariable("JWT_SIGN_KEY")
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

type FirestoreUserService struct {
	client *firestore.Client
}

func (mf *IMTokenCache) LoadUser() gin.HandlerFunc {
	return LookupUser(domain.LookupUser{
		Cache:         mf.cache,
		CacheDuration: 1 * time.Hour,
	})
}
