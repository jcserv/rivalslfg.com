package v1

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/types"
	"github.com/jcserv/rivalslfg/internal/utils"

	"github.com/jcserv/rivalslfg/internal/repository"
)

func Parse(params *httputil.QueryParams) (*repository.GetGroupsParams, error) {
	args := &repository.GetGroupsParams{
		RegionFilter:   "",
		GamemodeFilter: "",
		OpenFilter:     "",
		SizeSort:       "",
		Limit:          250,
		Offset:         0,
	}
	if params == nil {
		return args, nil
	}

	if err := parseFilters(args, params.FilterBy); err != nil {
		return nil, err
	}

	if err := parseSorting(args, params.SortBy); err != nil {
		return nil, err
	}

	if err := parsePagination(args, params.PaginateBy); err != nil {
		return nil, err
	}
	return args, nil
}

func parseFilters(args *repository.GetGroupsParams, filterBy []httputil.Filter) error {
	for _, filter := range filterBy {
		if filter.Field == "region" {
			switch filter.Value.(type) {
			case string:
				args.RegionFilter = strings.ToLower(filter.Value.(string))
			default:
				return fmt.Errorf("invalid type value for region filter value")
			}
		}
		if filter.Field == "gamemode" {
			switch filter.Value.(type) {
			case string:
				val := strings.ToLower(filter.Value.(string))
				if val != "competitive" && val != "quickplay" {
					return fmt.Errorf("invalid value for gamemode filter")
				}
				args.GamemodeFilter = val
			default:
				return fmt.Errorf("invalid type value for gamemode filter value")
			}
		}
		if filter.Field == "open" {
			switch filter.Value.(type) {
			case bool:
				if filter.Value.(bool) {
					args.OpenFilter = "true"
					break
				}
				args.OpenFilter = "false"
			default:
				return fmt.Errorf("invalid type value for open filter value")
			}
		}
	}
	return nil
}

func parseSorting(args *repository.GetGroupsParams, sorters []httputil.Sort) error {
	for _, sorter := range sorters {
		field := strings.ToLower(sorter.Field)
		if field == "size" {
			if sorter.Ascending {
				args.SizeSort = "asc"
				continue
			}
			args.SizeSort = "desc"
		}
	}
	return nil
}

func parsePagination(args *repository.GetGroupsParams, paginateBy *httputil.OffsetPagination) error {
	args.Limit = paginateBy.Limit
	args.Offset = paginateBy.Offset
	args.Count = paginateBy.Count
	return nil
}

type CreateGroup struct {
	Owner    string `json:"owner"`
	Region   string `json:"region"`
	Gamemode string `json:"gamemode"`
	Open     bool   `json:"open"`

	Platform   string   `json:"platform"`
	Roles      []string `json:"roles"`
	RankID     string   `json:"rankId"`
	Characters []string `json:"characters"`
	VoiceChat  bool     `json:"voiceChat"`
	Mic        bool     `json:"mic"`

	Vanguards   int `json:"vanguards"`
	Duelists    int `json:"duelists"`
	Strategists int `json:"strategists"`

	Platforms      []string `json:"platforms"`
	GroupVoiceChat bool     `json:"groupVoiceChat"`
	GroupMic       bool     `json:"groupMic"`
}

type CreateGroupResult struct {
	GroupID  string `json:"groupId"`
	PlayerID string `json:"playerId"`
}

func (c *CreateGroup) Parse() (*repository.CreateGroupParams, error) {
	params := &repository.CreateGroupParams{}
	if c.Owner == "" {
		return nil, fmt.Errorf("owner is required")
	}

	if c.Platform == "" {
		return nil, fmt.Errorf("platform is required")
	}

	if !types.Platforms.Contains(c.Platform) {
		return nil, fmt.Errorf("platform %s is not supported", c.Platform)
	}

	if len(c.Roles) > 0 && len(types.Roles.Intersection(types.NewSet(utils.StringSliceToLower(c.Roles)...))) != len(c.Roles) {
		return nil, fmt.Errorf("one or more provided roles %v is not supported", c.Roles)
	}

	if !types.IsValidRankID(c.RankID) {
		return nil, fmt.Errorf("invalid rank %s", c.RankID)
	}

	if !types.Regions.Contains(c.Region) {
		return nil, fmt.Errorf("region %s is not supported", c.Region)
	}

	if !types.Gamemodes.Contains(c.Gamemode) {
		return nil, fmt.Errorf("gamemode %s is not supported", c.Gamemode)
	}

	if c.Vanguards < 0 || c.Vanguards > 6 {
		return nil, fmt.Errorf("vanguards must be between 0 and 6")
	}

	if c.Duelists < 0 || c.Duelists > 6 {
		return nil, fmt.Errorf("duelists must be between 0 and 6")
	}

	if c.Strategists < 0 || c.Strategists > 6 {
		return nil, fmt.Errorf("strategists must be between 0 and 6")
	}

	if (c.Vanguards + c.Duelists + c.Strategists) != 6 {
		return nil, fmt.Errorf("vanguards, duelists, and strategists must add up to 6")
	}

	if len(types.Platforms.Intersection(types.NewSet(c.Platforms...))) != len(c.Platforms) {
		return nil, fmt.Errorf("one or more provided platforms %v is not supported", c.Platforms)
	}

	params.OwnerName = c.Owner
	params.Platform = c.Platform
	params.Roles = c.Roles
	params.RankValue = int32(types.RankIDToRankVal[c.RankID])
	params.Characters = c.Characters
	params.VoiceChat = c.VoiceChat
	params.Mic = c.Mic
	params.Vanguards = pgtype.Int4{Int32: int32(c.Vanguards), Valid: true}
	params.Duelists = pgtype.Int4{Int32: int32(c.Duelists), Valid: true}
	params.Strategists = pgtype.Int4{Int32: int32(c.Strategists), Valid: true}
	params.Platforms = c.Platforms
	params.GroupVoiceChat = pgtype.Bool{Bool: c.GroupVoiceChat, Valid: true}
	params.GroupMic = pgtype.Bool{Bool: c.GroupMic, Valid: true}

	return params, nil
}

type UpsertGroup struct {
	ID            string                     `json:"id"`
	Owner         string                     `json:"owner"`
	Region        string                     `json:"region"`
	Gamemode      string                     `json:"gamemode"`
	Players       []repository.PlayerInGroup `json:"players"`
	Open          bool                       `json:"open"`
	RoleQueue     *repository.RoleQueue      `json:"roleQueue"`
	GroupSettings *repository.GroupSettings  `json:"groupSettings"`
}

// // GetID returns the ID of the group, and also allows this to implement the RequestWithID interface
// func (c *UpsertGroup) GetID() string {
// 	return c.ID
// }

// func (c *UpsertGroup) ToMap() map[string]any {
// 	out := map[string]any{}
// 	_ = mapstructure.Decode(c, &out)
// 	return out
// }

// func (c *UpsertGroup) Parse() (*repository.UpsertGroupParams, error) {
// 	if c.Owner == "" {
// 		return nil, fmt.Errorf("owner is required")
// 	}
// 	if c.Region == "" {
// 		return nil, fmt.Errorf("region is required")
// 	}
// 	if c.Gamemode == "" {
// 		return nil, fmt.Errorf("gamemode is required")
// 	}
// 	if len(c.Players) == 0 {
// 		return nil, fmt.Errorf("one or more players are required")
// 	}

// 	params := &repository.UpsertGroupParams{
// 		ID:       c.ID,
// 		Owner:    c.Owner,
// 		Region:   c.Region,
// 		Gamemode: c.Gamemode,
// 		Players:  c.Players,
// 		Open:     c.Open,
// 	}

// 	if c.RoleQueue != nil {
// 		params.Vanguards = pgtype.Int4{Int32: int32(c.RoleQueue.Vanguards), Valid: true}
// 		params.Duelists = pgtype.Int4{Int32: int32(c.RoleQueue.Duelists), Valid: true}
// 		params.Strategists = pgtype.Int4{Int32: int32(c.RoleQueue.Strategists), Valid: true}
// 	}
// 	if c.GroupSettings != nil {
// 		params.Platforms = c.GroupSettings.Platforms
// 		params.VoiceChat = pgtype.Bool{Bool: c.GroupSettings.VoiceChat, Valid: true}
// 		params.Mic = pgtype.Bool{Bool: c.GroupSettings.Mic, Valid: true}
// 	}

// 	return params, nil
// }

// type JoinGroup struct {
// 	GroupID  string              `json:"groupId"`
// 	Player   *repository.Profile `json:"player"`
// 	Passcode string              `json:"passcode"`
// }

// func (c *JoinGroup) Parse() (*services.JoinGroupArgs, error) {
// 	if c.GroupID == "" {
// 		return nil, fmt.Errorf("groupId is required")
// 	}
// 	if c.Player == nil {
// 		return nil, fmt.Errorf("player is required")
// 	}

// 	player, err := json.Marshal(c.Player)
// 	if err != nil {
// 		return nil, err
// 	}

// 	args := &services.JoinGroupArgs{
// 		CheckCanJoinGroup: repository.CheckCanJoinGroupParams{
// 			ID:         c.GroupID,
// 			Passcode:   c.Passcode,
// 			PlayerName: c.Player.Name,
// 		},
// 		JoinGroup: repository.JoinGroupParams{
// 			Player: player,
// 			ID:     c.GroupID,
// 		},
// 	}
// 	return args, nil
// }

// type RemovePlayer struct {
// 	GroupID          string `json:"groupId"`
// 	RequesterID      string `json:"requesterId"`
// 	PlayerToRemoveID string `json:"playerToRemoveId"`
// }

// func (c *RemovePlayer) Parse() (*repository.RemovePlayerFromGroupParams, error) {
// 	if c.GroupID == "" {
// 		return nil, fmt.Errorf("groupId is required")
// 	}
// 	if c.RequesterID == "" {
// 		return nil, fmt.Errorf("requesterId is required")
// 	}
// 	if c.PlayerToRemoveID == "" {
// 		return nil, fmt.Errorf("playerToRemoveId is required")
// 	}

// 	params := &repository.RemovePlayerFromGroupParams{
// 		ID:               c.GroupID,
// 		PlayerToRemoveID: c.PlayerToRemoveID,
// 	}

// 	return params, nil
// }
