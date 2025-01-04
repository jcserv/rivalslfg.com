package env

import "os"

func GetBytes(key string, fallback []byte) []byte {
	if value, ok := os.LookupEnv(key); ok {
		return []byte(value)
	}
	return fallback
}

func GetString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}
	return fallback
}
