package rest_server_user

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserMiddlewares struct {
	jwtSecret []byte
}

func NewUserMiddlewares(jwtSecret []byte) *UserMiddlewares {
	return &UserMiddlewares{jwtSecret: jwtSecret}
}

type contextKey string

const (
	userIDKey  contextKey = "user_id"
	emailKey   contextKey = "email"
	isAdminKey contextKey = "is_admin"
)

func (s *UserMiddlewares) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token não fornecido", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return s.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "token inválido", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "token inválido", http.StatusUnauthorized)
			return
		}

		r.Header.Set("owner_id", claims["user_id"].(string))
		r.Header.Set("user_id", claims["user_id"].(string))
		ctx := context.WithValue(r.Context(), userIDKey, claims["user_id"])
		ctx = context.WithValue(ctx, emailKey, claims["email"])
		ctx = context.WithValue(ctx, isAdminKey, claims["is_admin"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
