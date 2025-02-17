package middleware

import (
	"fmt"
	"net/http"
	"os"
	"porty-go/models"
	"porty-go/services"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Status:  "error",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Status:  "error",
				Message: "Bearer token is required",
			})
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.ParseWithClaims(tokenString, &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*services.CustomClaims)
		fmt.Println("claims", claims.Email)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Status:  "error",
				Message: "Failed to parse token claims",
			})
			c.Abort()
			return
		}

		currentTimestamp := time.Now().Unix()
		if claims.ExpiresAt < currentTimestamp {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token has expired",
			})
			c.Abort()
			return
		}

		c.Set("user", claims)

		c.Next()
	}
}
