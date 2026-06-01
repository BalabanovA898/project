package delivery

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	userIDContextKey = "userID"
)

// AuthMiddleware проверяет JWT токен в Authorization хедере
func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Парсим токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(s.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Извлекаем claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		// Получаем userID из 'sub' claim
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id in token")
			return
		}

		// Сохраняем userID в контексте
		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// GetUserIDFromContext извлекает userID из контекста
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}
