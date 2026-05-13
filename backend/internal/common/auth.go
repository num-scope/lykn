package common

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authContextKey string

const (
	userIDContextKey   authContextKey = "user_id"
	usernameContextKey authContextKey = "username"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, secret string, ttl time.Duration) (string, time.Time, error) {
	expireAt := time.Now().Add(ttl)
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expireAt, nil
}

func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func WithUserContext(ctx context.Context, userID uint, username string) context.Context {
	ctx = context.WithValue(ctx, userIDContextKey, userID)
	return context.WithValue(ctx, usernameContextKey, username)
}

func UserIDFromContext(ctx context.Context) (uint, bool) {
	value, ok := ctx.Value(userIDContextKey).(uint)
	return value, ok
}

func UsernameFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(usernameContextKey).(string)
	return value, ok
}
