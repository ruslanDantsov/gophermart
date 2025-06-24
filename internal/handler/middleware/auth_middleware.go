package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CtxUserIDKey struct{}

func AuthMiddleware(jwtSecret string, logger *zap.Logger) gin.HandlerFunc {
	return func(gContext *gin.Context) {
		var tokenString string
		authHeader := gContext.GetHeader("Authorization")
		if len(authHeader) >= 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			if len(authHeader) > 0 {
				logger.Error(fmt.Sprintf("Value of auth header: %s", authHeader))
			}
			// Если нет в заголовке, пробуем получить из cookie
			cookie, err := gContext.Cookie("Authorization")
			if err == nil && cookie != "" {
				tokenString = cookie
			} else {
				logger.Error("Missing or invalid token in request")
				gContext.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
				gContext.Abort()
				return
			}
		}

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
			userID, _ := uuid.Parse(claims["id"].(string))
			gContext.Request = gContext.Request.WithContext(context.WithValue(gContext.Request.Context(), CtxUserIDKey{}, userID))
		}
		gContext.Next()
	}
}
