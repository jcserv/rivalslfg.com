package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestTokenService(t *testing.T) {
	service := NewTokenService("test-secret")

	t.Run("GenerateToken", func(t *testing.T) {
		// Test valid token generation
		t.Run("ValidInput", func(t *testing.T) {
			token, err := service.GenerateToken("user123", true)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		})
	})

	t.Run("ValidateToken", func(t *testing.T) {
		// Test valid token validation
		t.Run("ValidToken", func(t *testing.T) {
			token, _ := service.GenerateToken("user123", true)
			claims, err := service.ValidateToken(token)
			assert.NoError(t, err)

			subject, err := claims.GetSubject()
			assert.NoError(t, err)
			assert.Equal(t, "user123", subject)

			// assert.NoError(t, err)
			// assert.NotNil(t, claims["rights"])
		})

		// Test expired token
		t.Run("ExpiredToken", func(t *testing.T) {
			expiredToken := createExpiredToken(service.secretKey)
			_, err := service.ValidateToken(expiredToken)
			assert.Error(t, err)
		})

		// Test invalid token
		t.Run("InvalidToken", func(t *testing.T) {
			_, err := service.ValidateToken("invalid-token")
			assert.Error(t, err)
		})
	})
}

func createExpiredToken(secretKey []byte) string {
	claims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(-24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(secretKey)
	return signedToken
}
