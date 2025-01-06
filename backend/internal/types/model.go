package types

import (
	"fmt"

	"github.com/jcserv/rivalslfg/internal/utils"
)

var Gamemodes = NewSet("competitive", "quickplay")

func ValidateGamemode(gamemode string) error {
	if gamemode == "" {
		return fmt.Errorf("gamemode is required")
	}

	if !Gamemodes.Contains(gamemode) {
		return fmt.Errorf("gamemode %s is not supported", gamemode)
	}
	return nil
}

var Platforms = NewSet("xb", "ps", "pc")

func ValidatePlatform(platform string) error {
	if platform == "" {
		return fmt.Errorf("platform is required")
	}

	if !Platforms.Contains(platform) {
		return fmt.Errorf("platform %s is not supported", platform)
	}
	return nil
}

func ValidatePlatforms(platforms []string) error {
	if len(platforms) > 0 && len(Platforms.Intersection(NewSet(utils.StringSliceToLower(platforms)...))) != len(platforms) {
		return fmt.Errorf("one or more provided platforms %v is not supported", platforms)
	}
	return nil
}

var Regions = NewSet("na", "eu", "ap", "sa", "me")

func ValidateRegion(region string) error {
	if region == "" {
		return fmt.Errorf("region is required")
	}

	if !Regions.Contains(region) {
		return fmt.Errorf("region %s is not supported", region)
	}
	return nil
}

var Roles = NewSet("vanguard", "duelist", "strategist")

func ValidateRoles(roles []string) error {
	if len(roles) > 0 && len(Roles.Intersection(NewSet(utils.StringSliceToLower(roles)...))) != len(roles) {
		return fmt.Errorf("one or more provided roles %v is not supported", roles)
	}
	return nil
}

var RankIDToRankVal = map[string]int{
	"b3":  0,
	"b2":  1,
	"b1":  2,
	"s3":  10,
	"s2":  11,
	"s1":  12,
	"g3":  20,
	"g2":  21,
	"g1":  22,
	"p3":  30,
	"p2":  31,
	"p1":  32,
	"d3":  40,
	"d2":  41,
	"d1":  42,
	"gm3": 50,
	"gm2": 51,
	"gm1": 52,
	"e":   60,
	"oa":  70,
}

var RankValToRankID = map[int]string{
	0:  "b3",
	1:  "b2",
	2:  "b1",
	10: "s3",
	11: "s2",
	12: "s1",
	20: "g3",
	21: "g2",
	22: "g1",
	30: "p3",
	31: "p2",
	32: "p1",
	40: "d3",
	41: "d2",
	42: "d1",
	50: "gm3",
	51: "gm2",
	52: "gm1",
	60: "e",
	70: "oa",
}

func IsValidRankID(value string) bool {
	_, exists := RankIDToRankVal[value]
	return exists
}

func IsValidRankValue(value int) bool {
	_, exists := RankValToRankID[value]
	return exists
}

func ValidateRoleQueue(vanguards, duelists, strategists int) error {
	if vanguards < 0 || vanguards > 6 {
		return fmt.Errorf("vanguards must be between 0 and 6")
	}
	if duelists < 0 || duelists > 6 {
		return fmt.Errorf("duelists must be between 0 and 6")
	}
	if strategists < 0 || strategists > 6 {
		return fmt.Errorf("strategists must be between 0 and 6")
	}
	return nil
}
