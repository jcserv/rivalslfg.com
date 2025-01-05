package utils

import (
	"strconv"
	"strings"
)

func StringSliceToLower(s []string) []string {
	out := make([]string, len(s))
	for i, v := range s {
		out[i] = strings.ToLower(v)
	}
	return out
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
