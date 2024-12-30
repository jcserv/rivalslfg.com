// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Inventory struct {
	PlayerID int64     `json:"player_id"`
	ItemID   uuid.UUID `json:"item_id"`
}

type Item struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value int32     `json:"value"`
}

type Player struct {
	ID        int32       `json:"id"`
	Name      string      `json:"name"`
	Level     int32       `json:"level"`
	Class     string      `json:"class"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Gold      pgtype.Int8 `json:"gold"`
}
