package repository

import "time"

type GroupDTO struct {
	ID            string         `json:"id"`
	CommunityID   int32          `json:"communityId"`
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
	Name         string    `json:"name"`
	Size         int       `json:"size"`
	LastActiveAt time.Time `json:"lastActiveAt"`

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

type Profile struct {
	Name          string         `json:"name"`
	Region        string         `json:"region"`
	Platform      string         `json:"platform"`
	Gamemode      string         `json:"gamemode"`
	Roles         []string       `json:"roles"`
	Rank          Rankid         `json:"rank"`
	Characters    []string       `json:"characters"`
	VoiceChat     bool           `json:"voiceChat"`
	Mic           bool           `json:"mic"`
	RoleQueue     *RoleQueue     `json:"roleQueue"`
	GroupSettings *GroupSettings `json:"groupSettings"`
}

type PlayerInGroup struct {
	Name       string   `json:"name"`
	Leader     bool     `json:"leader"`
	Platform   string   `json:"platform"`
	Roles      []string `json:"roles"`
	Rank       string   `json:"rank"`
	Characters []string `json:"characters"`
	VoiceChat  bool     `json:"voiceChat"`
	Mic        bool     `json:"mic"`
}
