package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func AuthMiddleware(jwtSecret string, logger *zap.Logger) gin.HandlerFunc {
	return func(gContext *gin.Context) {
		authHeader := gContext.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			logger.Error("Missing or invalid token in request")
			gContext.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			gContext.Abort()
			return
		}

		tokenString := authHeader[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Error("Failed to parse token", zap.Error(err))
			gContext.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			gContext.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		exp, _ := claims["exp"].(float64)
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			logger.Error("Token has expired")
			gContext.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			gContext.Abort()
			return
		}

		if ok {
			gContext.Set("username", claims["username"])
			gContext.Set("id", claims["id"])
		}

		gContext.Next()
	}
}
