package v1

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jcserv/rivalslfg/internal/transport/http/httputil"
	"github.com/jcserv/rivalslfg/internal/types"

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

type PlayerRequirements struct {
	Gamemode  string `json:"gamemode,omitempty"`
	Region    string `json:"region,omitempty"`
	Platform  string `json:"platform,omitempty"`
	Role      string `json:"role,omitempty"`
	RankID    string `json:"rank,omitempty"`
	VoiceChat bool   `json:"voiceChat,omitempty"`
	Mic       bool   `json:"mic,omitempty"`
}

func (p *PlayerRequirements) Validate() error {
	if p == nil {
		return nil
	}

	if p.Gamemode != "" {
		if err := types.ValidateGamemode(p.Gamemode); err != nil {
			return err
		}
	}

	if p.Region != "" {
		if err := types.ValidateRegion(p.Region); err != nil {
			return err
		}
	}

	if p.Platform != "" {
		if err := types.ValidatePlatform(p.Platform); err != nil {
			return err
		}
	}

	if p.Role != "" {
		if err := types.ValidateRole(p.Role); err != nil {
			return err
		}
	}

	if p.RankID != "" && !types.IsValidRankID(p.RankID) {
		return fmt.Errorf("invalid rank %s", p.RankID)
	}

	return nil
}

func (p *PlayerRequirements) ToParams() (*repository.GetGroupsParams, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	var gamemode string
	if p.Gamemode != "" {
		gamemode = p.Gamemode
	}

	var region string
	if p.Region != "" {
		region = p.Region
	}

	var rankVal *int32
	if p.RankID != "" {
		val := int32(types.RankIDToRankVal[p.RankID])
		rankVal = &val
	}

	var platform *string
	if p.Platform != "" {
		platform = &p.Platform
	}

	var role *string
	if p.Role != "" {
		r := strings.ToLower(p.Role)
		role = &r
	}

	voiceChat := &p.VoiceChat
	mic := &p.Mic

	return &repository.GetGroupsParams{
		GamemodeFilter: gamemode,
		RegionFilter:   region,
		Platform:       platform,
		Role:           role,
		RankVal:        rankVal,
		VoiceChat:      voiceChat,
		Mic:            mic,
	}, nil
}

type CreateGroup struct {
	PlayerID int    `json:"player_id"`
	GroupID  string `json:"group_id"`
	Owner    string `json:"owner"`
	Region   string `json:"region"`
	Gamemode string `json:"gamemode"`
	Open     bool   `json:"open"`

	Platform   string   `json:"platform"`
	Role       string   `json:"role"`
	RankID     string   `json:"rankId"`
	Characters []string `json:"characters"`
	VoiceChat  bool     `json:"voiceChat"`
	Mic        bool     `json:"mic"`

	Vanguards   int `json:"vanguards"`
	Duelists    int `json:"duelists"`
	Strategists int `json:"strategists"`

	GroupPlatform  string `json:"groupPlatform"`
	GroupVoiceChat bool   `json:"groupVoiceChat"`
	GroupMic       bool   `json:"groupMic"`
}

type CreateGroupResult struct {
	GroupID  string `json:"groupId"`
	PlayerID string `json:"playerId"`
}

func (c *CreateGroup) validate() error {
	if c.Owner == "" {
		return fmt.Errorf("owner is required")
	}

	if err := types.ValidatePlatform(c.Platform); err != nil {
		return err
	}

	if err := types.ValidateRole(c.Role); err != nil {
		return err
	}

	if !types.IsValidRankID(c.RankID) {
		return fmt.Errorf("invalid rank %s", c.RankID)
	}

	if !types.Regions.Contains(c.Region) {
		return fmt.Errorf("region %s is not supported", c.Region)
	}

	if !types.Gamemodes.Contains(c.Gamemode) {
		return fmt.Errorf("gamemode %s is not supported", c.Gamemode)
	}

	if err := types.ValidateRoleQueue(c.Vanguards, c.Duelists, c.Strategists); err != nil {
		return err
	}

	if c.GroupPlatform != "" {
		if err := types.ValidatePlatform(c.GroupPlatform); err != nil {
			return err
		}
	}
	return nil
}

func (c *CreateGroup) Parse() (*repository.CreateGroupParams, error) {
	params := &repository.CreateGroupParams{}

	if err := c.validate(); err != nil {
		return nil, err
	}

	params.PlayerID = int32(c.PlayerID)
	params.GroupID = c.GroupID
	params.Owner = c.Owner
	params.Platform = c.Platform
	params.Role = strings.ToLower(c.Role)
	params.RankVal = int32(types.RankIDToRankVal[c.RankID])
	params.Characters = c.Characters
	params.VoiceChat = c.VoiceChat
	params.Mic = c.Mic
	params.Region = c.Region
	params.Gamemode = c.Gamemode
	params.Open = c.Open

	params.Vanguards = int32(c.Vanguards)
	params.Duelists = int32(c.Duelists)
	params.Strategists = int32(c.Strategists)

	params.Platform = c.Platform
	params.GroupVoiceChat = pgtype.Bool{Bool: c.GroupVoiceChat, Valid: true}
	params.GroupMic = pgtype.Bool{Bool: c.GroupMic, Valid: true}

	return params, nil
}

type JoinGroup struct {
	GroupID  string `json:"groupId"`
	PlayerID int    `json:"playerId"`

	Name        string   `json:"name"`
	Passcode    string   `json:"passcode"`
	Platform    string   `json:"platform"`
	Gamemode    string   `json:"gamemode"`
	Region      string   `json:"region"`
	Role        string   `json:"role"`
	RankID      string   `json:"rankId"`
	Characters  []string `json:"characters"`
	VoiceChat   bool     `json:"voiceChat"`
	Mic         bool     `json:"mic"`
	Vanguards   int      `json:"vanguards"`
	Duelists    int      `json:"duelists"`
	Strategists int      `json:"strategists"`
}

func (c *JoinGroup) validate() error {
	if c.GroupID == "" {
		return fmt.Errorf("groupId is required")
	}

	if c.Name == "" {
		return fmt.Errorf("playerName is required")
	}

	if err := types.ValidateGamemode(c.Gamemode); err != nil {
		return err
	}

	if err := types.ValidateRegion(c.Region); err != nil {
		return err
	}

	if err := types.ValidatePlatform(c.Platform); err != nil {
		return err
	}

	if err := types.ValidateRole(c.Role); err != nil {
		return err
	}

	if valid := types.IsValidRankID(c.RankID); !valid {
		return fmt.Errorf("rankId %s is invalid", c.RankID)
	}

	if err := types.ValidateRoleQueue(c.Vanguards, c.Duelists, c.Strategists); err != nil {
		return err
	}
	return nil
}

func (c *JoinGroup) Parse() (*repository.JoinGroupParams, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	params := &repository.JoinGroupParams{}
	params.GroupID = c.GroupID
	params.PlayerID = int32(c.PlayerID)
	params.Gamemode = c.Gamemode
	params.Region = c.Region
	params.Platform = c.Platform
	params.Role = strings.ToLower(c.Role)
	params.RankVal = int32(types.RankIDToRankVal[c.RankID])
	params.Name = c.Name
	params.Passcode = c.Passcode
	params.Characters = c.Characters
	params.VoiceChat = c.VoiceChat
	params.Mic = c.Mic
	params.Vanguards = int32(c.Vanguards)
	params.Duelists = int32(c.Duelists)
	params.Strategists = int32(c.Strategists)
	return params, nil
}

type RemovePlayer struct {
	GroupID          string `json:"groupId"`
	PlayerToRemoveID int    `json:"playerToRemoveId"`
}

func (c *RemovePlayer) validate() error {
	if c.GroupID == "" {
		return fmt.Errorf("groupId is required")
	}
	if c.PlayerToRemoveID <= 0 {
		return fmt.Errorf("playerId is required")
	}
	return nil
}

func (c *RemovePlayer) Parse() (*repository.RemovePlayerParams, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	params := &repository.RemovePlayerParams{}
	params.GroupID = c.GroupID
	params.PlayerID = int32(c.PlayerToRemoveID)
	return params, nil
}
