package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type client interface {
	CheckJWTAndGetUserID(token string) (uint64, error)
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

func (a *AuthMiddleware) JwtUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(16, "missing metadata")
	}

	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, grpc.Errorf(16, "missing token")
	}

	tokenString := strings.TrimPrefix(authHeader[0], BearerPrefix)

	userID, err := a.client.CheckJWTAndGetUserID(tokenString)
	if err != nil || userID == 0 {
		return nil, grpc.Errorf(16, "invalid token")
	}

	newCtx := context.WithValue(ctx, ContextUserIDKey, userID)
	return handler(newCtx, req)
}
