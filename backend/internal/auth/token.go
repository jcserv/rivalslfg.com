package auth

import (
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jcserv/rivalslfg/internal/utils/env"
)

func getSecretKey() []byte {
	return env.GetBytes("JWT_SECRET_KEY", []byte("you cant skip lunch"))
}

func GenerateToken(subject string, additionalClaims map[string]string, additionalRights ...Right) (string, error) {
	sessionID := make([]byte, 32)
	if _, err := rand.Read(sessionID); err != nil {
		return "", err
	}

	rights := append(baseRights, additionalRights...)
	claims := jwt.MapClaims{
		"sub":    subject,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"rights": rights,
	}
	for k, v := range additionalClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return getSecretKey(), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

func HasRight(claims map[string]interface{}, requiredRight Right) bool {
	rights, ok := claims["rights"].([]interface{})
	if !ok {
		return false
	}

	for _, right := range rights {
		if IsEqual(right.(string), requiredRight) {
			return true
		}
	}
	return false
}
