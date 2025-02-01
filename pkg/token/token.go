package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

const TokenDuration = 24 * time.Hour

type JwtService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJwtService(secretKey []byte) *JwtService {
	return &JwtService{
		secretKey:     secretKey,
		tokenDuration: TokenDuration,
	}
}

func (j *JwtService) GenerateTokenWithID(userID uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   fmt.Sprint(userID),
		ExpiresAt: time.Now().Add(j.tokenDuration).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString(j.secretKey)
}

func (j *JwtService) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (j *JwtService) CheckJWTAndGetUserID(tokenString string) (uint64, error) {
	token, err := j.VerifyToken(tokenString)
	if err != nil {
		return 0, err
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	if exp, ok := mapClaims["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return 0, errors.New("token is expired")
		}
	} else {
		return 0, errors.New("expiration time missing from token claims")
	}

	userID, err := strconv.ParseUint(mapClaims["sub"].(string), 10, 64)
	if err != nil {
		return 0, errors.New("invalid user ID in token claims")
	}

	return userID, nil
}
