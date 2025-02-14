package utils

import (
	"os"
	"strconv"
)

// GetEnv reads an environment variable and falls back to the default value if not set.
func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetIntEnv reads an environment variable as an integer and falls back to the default value if not set or invalid.
func GetIntEnv(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
