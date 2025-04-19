package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var ErrNoToken = errors.New("no bearer token")
var ErrInvalidToken = errors.New("invalid token")

// JWTAuthenticator uses HMAC‑signed JWTs in the Authorization header.
type JWTAuthenticator struct {
	secret []byte
}

// NewJWTAuthenticator constructs one.
func NewJWTAuthenticator(secret string) *JWTAuthenticator {
	return &JWTAuthenticator{secret: []byte(secret)}
}

func (j *JWTAuthenticator) Authenticate(r *http.Request) (int64, error) {
	hdr := r.Header.Get("Authorization")
	if !strings.HasPrefix(hdr, "Bearer ") {
		return 0, ErrNoToken
	}
	tokenStr := strings.TrimPrefix(hdr, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}
	return int64(sub), nil
}
