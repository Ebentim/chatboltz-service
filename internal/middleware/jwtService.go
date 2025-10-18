package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user entity.Users, jwtSecret []byte) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   entity.UserRole(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Chatboltz",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)), //!2 hours expiry
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader string
		authHeader = c.GetHeader("Authorization")
		if authHeader == "" {
			token, err := c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
				c.Abort()
				return
			}
			authHeader = "Bearer " + token
		}

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if err := writeClaimsToLog(claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log claims"})
			c.Abort()
			return
		} // Attach claims to context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", string(claims.Role))

		c.Next()
	}
}

func writeClaimsToLog(claims *JWTClaims) error {
	// Format claims data as string
	claimsData := fmt.Sprintf("UserID: %s\nEmail: %s\nRole: %s\nIssuer: %s\nIssuedAt: %s\nExpiresAt: %s\n",
		claims.UserID,
		claims.Email,
		claims.Role,
		claims.Issuer,
		claims.IssuedAt.Time.String(),
		claims.ExpiresAt.Time.String())

	// Write to log file
	err := os.WriteFile("user.log", []byte(claimsData), 0644)
	if err != nil {
		return fmt.Errorf("failed to write claims to log: %v", err)
	}
	return nil
}
