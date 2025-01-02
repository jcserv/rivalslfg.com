package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jcserv/rivalslfg/internal/repository"
	"github.com/jcserv/rivalslfg/internal/services"
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
	return nil
}

type UpsertGroup struct {
	ID            string                     `json:"id"`
	Owner         string                     `json:"owner"`
	Region        string                     `json:"region"`
	Gamemode      string                     `json:"gamemode"`
	Players       []repository.PlayerInGroup `json:"players"`
	Open          bool                       `json:"open"`
	RoleQueue     *repository.RoleQueue
	GroupSettings *repository.GroupSettings
}

func (c *UpsertGroup) Parse() (*repository.UpsertGroupParams, error) {
	if c.Owner == "" {
		return nil, fmt.Errorf("owner is required")
	}
	if c.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	if c.Gamemode == "" {
		return nil, fmt.Errorf("gamemode is required")
	}
	if len(c.Players) == 0 {
		return nil, fmt.Errorf("one or more players are required")
	}

	params := &repository.UpsertGroupParams{
		ID:       c.ID,
		Owner:    c.Owner,
		Region:   c.Region,
		Gamemode: c.Gamemode,
		Players:  c.Players,
		Open:     c.Open,
	}

	if c.RoleQueue != nil {
		params.Vanguards = pgtype.Int4{Int32: int32(c.RoleQueue.Vanguards), Valid: true}
		params.Duelists = pgtype.Int4{Int32: int32(c.RoleQueue.Duelists), Valid: true}
		params.Strategists = pgtype.Int4{Int32: int32(c.RoleQueue.Strategists), Valid: true}
	}
	if c.GroupSettings != nil {
		params.Platforms = c.GroupSettings.Platforms
		params.VoiceChat = pgtype.Bool{Bool: c.GroupSettings.VoiceChat, Valid: true}
		params.Mic = pgtype.Bool{Bool: c.GroupSettings.Mic, Valid: true}
	}

	return params, nil
}

type JoinGroup struct {
	GroupID  string              `json:"groupId"`
	Player   *repository.Profile `json:"player"`
	Passcode string              `json:"passcode"`
}

func (c *JoinGroup) Parse() (*services.JoinGroupArgs, error) {
	if c.GroupID == "" {
		return nil, fmt.Errorf("groupId is required")
	}
	if c.Player == nil {
		return nil, fmt.Errorf("player is required")
	}

	player, err := json.Marshal(c.Player)
	if err != nil {
		return nil, err
	}

	args := &services.JoinGroupArgs{
		CheckCanJoinGroup: repository.CheckCanJoinGroupParams{
			ID:         c.GroupID,
			Passcode:   c.Passcode,
			PlayerName: c.Player.Name,
		},
		JoinGroup: repository.JoinGroupParams{
			Player: player,
			ID:     c.GroupID,
		},
	}
	return args, nil
}

type RemovePlayer struct {
	GroupID       string `json:"groupId"`
	PlayerName    string `json:"playerName"`
	RequesterName string `json:"requesterName"` // Name of user making request
}

func (c *RemovePlayer) Parse() (*repository.RemovePlayerFromGroupParams, error) {
	if c.GroupID == "" {
		return nil, fmt.Errorf("groupId is required")
	}
	if c.PlayerName == "" {
		return nil, fmt.Errorf("playerName is required")
	}

	// TODO[AUTH]: This should be provided via the cookie/session
	if c.RequesterName == "" {
		return nil, fmt.Errorf("requesterName is required")
	}

	params := &repository.RemovePlayerFromGroupParams{
		ID:            c.GroupID,
		PlayerName:    c.PlayerName,
		RequesterName: c.RequesterName,
	}

	return params, nil
}
