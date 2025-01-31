package middleware

import (
	"context"
	"net/http"
	"strings"
)

type client interface {
	VerifyToken(ctx context.Context, token string) uint64
}

type ContextKey string

const BearerPrefix = "Bearer "
const ContextUserIDKey ContextKey = "userID"

type AuthMiddleware struct {
	client client
}

func NewAuthMiddleware(client client) *AuthMiddleware {
	return &AuthMiddleware{
		client: client,
	}
}

func (a *AuthMiddleware) JwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		tokenString = strings.TrimPrefix(tokenString, BearerPrefix)

		userId := a.client.VerifyToken(context.Background(), tokenString)

		if userId == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIDKey, userId)
		next(w, r.WithContext(ctx))
	}
}
