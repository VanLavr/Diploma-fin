package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/VanLavr/Diploma-fin/utils/config"
)

type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{secret: config.Secret}
}

// CustomClaims contains custom JWT claims along with standard registered claims
type CustomClaims struct {
	UserUUID string `json:"user_uuid"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new JWT with user UUID in claims
func (a *AuthMiddleware) GenerateAccessToken(userUUID string) (string, error) {
	// Create the claims with user UUID and expiration time
	claims := CustomClaims{
		UserUUID: userUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(120 * time.Minute)), // Token expires in 120 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateAccessToken is a Gin middleware that validates JWT and checks user UUID
func (a *AuthMiddleware) ValidateAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Extract token from Bearer scheme
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token malformed"})
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(a.secret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check if token is valid
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Type assert claims to access custom claims
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Check if UUID exists and is valid
		if claims.UserUUID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user UUID in token"})
			return
		}

		// You can add additional UUID validation logic here
		// For example: check against database or route parameters

		// Set user UUID in context for subsequent handlers
		c.Set("user_uuid", claims.UserUUID)
		c.Next()
	}
}
