package repository

type Group struct {
	ID            string         `json:"id"`
	CommunityID   int32          `json:"communityId"`
	Owner         string         `json:"owner"`
	Region        string         `json:"region"`
	Gamemode      string         `json:"gamemode"`
	Open          bool           `json:"open"`
	Passcode      string         `json:"passcode"`
	RoleQueue     *RoleQueue     `json:"roleQueue"`
	GroupSettings *GroupSettings `json:"groupSettings"`
}

type GroupMember struct {
	GroupID  string `json:"group_id"`
	PlayerID int32  `json:"player_id"`
	Leader   bool   `json:"leader"`
}

type GroupWithPlayers struct {
	Group
	Name    string   `json:"name"`
	Players []Player `json:"players"`
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

type Player struct {
	ID            int32          `json:"id"`
	Name          string         `json:"name"`
	DisplayName   string         `json:"displayName"`
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
