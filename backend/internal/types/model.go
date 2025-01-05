package types

var Gamemodes = NewSet("competitive", "quickplay")

var Platforms = NewSet("xb", "ps", "pc")

var Regions = NewSet("na", "eu", "ap", "sa", "me")

var Roles = NewSet("vanguard", "duelist", "strategist")

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
