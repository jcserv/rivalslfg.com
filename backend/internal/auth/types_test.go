package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_HasOwnership(t *testing.T) {
	t.Run("GroupOwnership", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":      "1",
			"playerId": "1",
			"groupId":  "AAAA",
		}

		assert.True(t, HasOwnership(claims, "group", "AAAA"))
		assert.False(t, HasOwnership(claims, "group", "AAAB"))
	})
	t.Run("PlayerOwnership", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":      "1",
			"playerId": "1",
			"groupId":  "AAAA",
		}

		assert.True(t, HasOwnership(claims, "player", "1"))
		assert.False(t, HasOwnership(claims, "player", "2"))
	})
}
