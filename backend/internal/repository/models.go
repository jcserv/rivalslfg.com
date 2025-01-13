// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Community struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type Group struct {
	ID           string      `json:"id"`
	CommunityID  int32       `json:"community_id"`
	Owner        pgtype.Text `json:"owner"`
	Region       string      `json:"region"`
	Gamemode     string      `json:"gamemode"`
	Open         bool        `json:"open"`
	Passcode     string      `json:"passcode"`
	Vanguards    int32       `json:"vanguards"`
	Duelists     int32       `json:"duelists"`
	Strategists  int32       `json:"strategists"`
	Platform     string      `json:"platform"`
	VoiceChat    pgtype.Bool `json:"voice_chat"`
	Mic          pgtype.Bool `json:"mic"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	LastActiveAt time.Time   `json:"last_active_at"`
}

type Groupmember struct {
	GroupID  string `json:"group_id"`
	PlayerID int32  `json:"player_id"`
	Leader   bool   `json:"leader"`
}

type Player struct {
	ID          int32    `json:"id"`
	Name        string   `json:"name"`
	Platform    string   `json:"platform"`
	Role        string   `json:"role"`
	Rank        int32    `json:"rank"`
	Characters  []string `json:"characters"`
	VoiceChat   bool     `json:"voice_chat"`
	Mic         bool     `json:"mic"`
	Vanguards   int32    `json:"vanguards"`
	Duelists    int32    `json:"duelists"`
	Strategists int32    `json:"strategists"`
}

type Rank struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value int32  `json:"value"`
}
