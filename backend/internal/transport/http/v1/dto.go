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

func (dto *PlayerRequirements) Validate() error {
	if dto == nil {
		return nil
	}

	if dto.Gamemode != "" {
		if err := types.ValidateGamemode(dto.Gamemode); err != nil {
			return err
		}
	}

	if dto.Region != "" {
		if err := types.ValidateRegion(dto.Region); err != nil {
			return err
		}
	}

	if dto.Platform != "" {
		if err := types.ValidatePlatform(dto.Platform); err != nil {
			return err
		}
	}

	if dto.Role != "" {
		if err := types.ValidateRole(dto.Role); err != nil {
			return err
		}
	}

	if dto.RankID != "" && !types.IsValidRankID(dto.RankID) {
		return fmt.Errorf("invalid rank %s", dto.RankID)
	}

	return nil
}

func (dto *PlayerRequirements) ToParams() (*repository.GetGroupsParams, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	var gamemode string
	if dto.Gamemode != "" {
		gamemode = dto.Gamemode
	}

	var region string
	if dto.Region != "" {
		region = dto.Region
	}

	var rankVal *int32
	if dto.RankID != "" {
		val := int32(types.RankIDToRankVal[dto.RankID])
		rankVal = &val
	}

	var platform *string
	if dto.Platform != "" {
		platform = &dto.Platform
	}

	var role *string
	if dto.Role != "" {
		r := strings.ToLower(dto.Role)
		role = &r
	}

	voiceChat := &dto.VoiceChat
	mic := &dto.Mic

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

func (dto *CreateGroup) validate() error {
	if dto.Owner == "" {
		return fmt.Errorf("owner is required")
	}

	if err := types.ValidatePlatform(dto.Platform); err != nil {
		return err
	}

	if err := types.ValidateRole(dto.Role); err != nil {
		return err
	}

	if !types.IsValidRankID(dto.RankID) {
		return fmt.Errorf("invalid rank %s", dto.RankID)
	}

	if !types.Regions.Contains(dto.Region) {
		return fmt.Errorf("region %s is not supported", dto.Region)
	}

	if !types.Gamemodes.Contains(dto.Gamemode) {
		return fmt.Errorf("gamemode %s is not supported", dto.Gamemode)
	}

	if err := types.ValidateRoleQueue(dto.Vanguards, dto.Duelists, dto.Strategists); err != nil {
		return err
	}

	if dto.GroupPlatform != "" {
		if err := types.ValidatePlatform(dto.GroupPlatform); err != nil {
			return err
		}
	}
	return nil
}

func (dto *CreateGroup) Parse() (*repository.CreateGroupParams, error) {
	params := &repository.CreateGroupParams{}

	if err := dto.validate(); err != nil {
		return nil, err
	}

	params.PlayerID = int32(dto.PlayerID)
	params.GroupID = dto.GroupID
	params.Owner = dto.Owner
	params.Platform = dto.Platform
	params.Role = strings.ToLower(dto.Role)
	params.RankVal = int32(types.RankIDToRankVal[dto.RankID])
	params.Characters = dto.Characters
	params.VoiceChat = dto.VoiceChat
	params.Mic = dto.Mic
	params.Region = dto.Region
	params.Gamemode = dto.Gamemode
	params.Open = dto.Open

	params.Vanguards = int32(dto.Vanguards)
	params.Duelists = int32(dto.Duelists)
	params.Strategists = int32(dto.Strategists)

	params.Platform = dto.Platform
	params.GroupVoiceChat = pgtype.Bool{Bool: dto.GroupVoiceChat, Valid: true}
	params.GroupMic = pgtype.Bool{Bool: dto.GroupMic, Valid: true}

	return params, nil
}

type PatchGroup struct {
	ID   string `json:"id"`
	Open bool   `json:"open"`
}

func (dto *PatchGroup) validate() error {
	if dto.ID == "" {
		return fmt.Errorf("groupId is required")
	}
	return nil
}

func (dto *PatchGroup) Parse() (*repository.PatchGroupParams, error) {
	if err := dto.validate(); err != nil {
		return nil, err
	}

	params := &repository.PatchGroupParams{}
	params.ID = dto.ID
	params.Open = dto.Open
	return params, nil
}

type JoinGroup struct {
	GroupID  string `json:"groupId"`
	PlayerID int    `json:"playerId"`

	Name             string   `json:"name"`
	Passcode         string   `json:"passcode"`
	Platform         string   `json:"platform"`
	Gamemode         string   `json:"gamemode"`
	Region           string   `json:"region"`
	Role             string   `json:"role"`
	RankID           string   `json:"rankId"`
	Characters       []string `json:"characters"`
	VoiceChat        bool     `json:"voiceChat"`
	Mic              bool     `json:"mic"`
	RoleQueueEnabled bool     `json:"roleQueueEnabled"`
}

func (dto *JoinGroup) validate() error {
	if dto.GroupID == "" {
		return fmt.Errorf("groupId is required")
	}

	if dto.Name == "" {
		return fmt.Errorf("playerName is required")
	}

	if err := types.ValidateGamemode(dto.Gamemode); err != nil {
		return err
	}

	if err := types.ValidateRegion(dto.Region); err != nil {
		return err
	}

	if err := types.ValidatePlatform(dto.Platform); err != nil {
		return err
	}

	if err := types.ValidateRole(dto.Role); err != nil {
		return err
	}

	if valid := types.IsValidRankID(dto.RankID); !valid {
		return fmt.Errorf("rankId %s is invalid", dto.RankID)
	}
	return nil
}

func (dto *JoinGroup) Parse() (*repository.JoinGroupParams, error) {
	if err := dto.validate(); err != nil {
		return nil, err
	}
	params := &repository.JoinGroupParams{}
	params.GroupID = dto.GroupID
	params.PlayerID = int32(dto.PlayerID)
	params.Gamemode = dto.Gamemode
	params.Region = dto.Region
	params.Platform = dto.Platform
	params.Role = strings.ToLower(dto.Role)
	params.RankVal = int32(types.RankIDToRankVal[dto.RankID])
	params.Name = dto.Name
	params.Passcode = dto.Passcode
	params.Characters = dto.Characters
	params.VoiceChat = dto.VoiceChat
	params.Mic = dto.Mic
	params.RoleQueueEnabled = dto.RoleQueueEnabled
	return params, nil
}

type RemovePlayer struct {
	GroupID          string `json:"groupId"`
	PlayerToRemoveID int    `json:"playerToRemoveId"`
}

func (dto *RemovePlayer) validate() error {
	if dto.GroupID == "" {
		return fmt.Errorf("groupId is required")
	}
	if dto.PlayerToRemoveID <= 0 {
		return fmt.Errorf("playerId is required")
	}
	return nil
}

func (dto *RemovePlayer) Parse() (*repository.RemovePlayerParams, error) {
	if err := dto.validate(); err != nil {
		return nil, err
	}
	params := &repository.RemovePlayerParams{}
	params.GroupID = dto.GroupID
	params.PlayerID = int32(dto.PlayerToRemoveID)
	return params, nil
}
