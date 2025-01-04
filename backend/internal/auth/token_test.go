package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var mockSecretKey = []byte("test-secret")

func TestTokenService(t *testing.T) {
	t.Run("GenerateToken", func(t *testing.T) {
		t.Run("ValidInput", func(t *testing.T) {
			token, err := GenerateToken("1", map[string]string{
				"playerId": "1",
				"groupId":  "AAAA",
			})
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		})
	})

	t.Run("ValidateToken", func(t *testing.T) {
		t.Run("ValidToken", func(t *testing.T) {
			token, _ := GenerateToken("1", map[string]string{
				"playerId": "1",
				"groupId":  "AAAA",
			})
			claims, err := ValidateToken(token)
			assert.NoError(t, err)

			subject, err := claims.GetSubject()
			assert.NoError(t, err)
			assert.Equal(t, "1", subject)
		})

		t.Run("ExpiredToken", func(t *testing.T) {
			expiredToken := createExpiredToken(mockSecretKey)
			_, err := ValidateToken(expiredToken)
			assert.Error(t, err)
		})

		t.Run("InvalidToken", func(t *testing.T) {
			_, err := ValidateToken("invalid-token")
			assert.Error(t, err)
		})
	})
}

func createExpiredToken(secretKey []byte) string {
	claims := jwt.MapClaims{
		"sub": "1",
		"exp": time.Now().Add(-24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(secretKey)
	return signedToken
}
