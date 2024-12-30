package repository

type GroupWithPlayers struct {
	Group
	Players []Player `json:"players"`
}
