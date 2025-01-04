package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jcserv/rivalslfg/internal/auth"
	"github.com/jcserv/rivalslfg/internal/store"
	"github.com/jcserv/rivalslfg/internal/utils/log"
)

type Token string

type Auth struct {
	store store.Store
}

func NewAuth(store store.Store) *Auth {
	return &Auth{
		store: store,
	}
}

type PlayerAuth struct {
	PlayerID string    `json:"playerId"`
	Token    string    `json:"token"`
	LastSeen time.Time `json:"lastSeen"`
}

func (p *PlayerAuth) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"playerId": p.PlayerID,
		"token":    p.Token,
		"lastSeen": p.LastSeen,
	}
}

func (s *Auth) CreateAuth(ctx context.Context, playerID string) (Token, error) {
	token, err := auth.GenerateToken(playerID, map[string]string{
		"playerId": playerID,
	})
	if err != nil {
		return "", err
	}

	auth := PlayerAuth{
		PlayerID: playerID,
		Token:    token,
		LastSeen: time.Now(),
	}
	s.store.Set(ctx, token, auth.ToMap())
	return Token(token), nil
}

func (s *Auth) ValidateToken(ctx context.Context, token Token) (*PlayerAuth, error) {
	val, err := s.store.Get(ctx, string(token))
	if err != nil {
		return nil, err
	}

	log.Info(ctx, fmt.Sprintf("value: %v", val))
	s.store.Expire(ctx, string(token), 1*time.Hour)
	return &PlayerAuth{}, nil //auth, nil
}
