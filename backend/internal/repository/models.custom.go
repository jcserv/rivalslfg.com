package repository

import "time"

type GroupDTO struct {
	ID            string         `json:"id"`
	CommunityID   int32          `json:"communityId"`
	OwnerID       int32          `json:"ownerId"`
	Owner         string         `json:"owner"`
	Region        string         `json:"region"`
	Gamemode      string         `json:"gamemode"`
	Open          bool           `json:"open"`
	Passcode      string         `json:"passcode"`
	RoleQueue     *RoleQueue     `json:"roleQueue"`
	GroupSettings *GroupSettings `json:"groupSettings"`
	LastActiveAt  time.Time      `json:"lastActiveAt"`
}

type GroupWithPlayers struct {
	GroupDTO

	// Computed fields
	Name string `json:"name"`
	Size int    `json:"size"`

	Players []PlayerInGroup `json:"players"`
}

type RoleQueue struct {
	Vanguards   int `json:"vanguards"`
	Duelists    int `json:"duelists"`
	Strategists int `json:"strategists"`
}

type GroupSettings struct {
	Platforms []string `json:"platforms"`
	VoiceChat bool     `json:"voiceChat"`
	Mic       bool     `json:"mic"`
}

type PlayerInGroup struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Leader     bool     `json:"leader"`
	Platform   string   `json:"platform"`
	Role       string   `json:"role"`
	Rank       string   `json:"rank"`
	Characters []string `json:"characters"`
	VoiceChat  bool     `json:"voiceChat"`
	Mic        bool     `json:"mic"`
}
