// Package util implements utility methods
package util

import (
	"os"
	"path/filepath"
)

// GetEnvKey tries to get the specified key from the OS environment and returns either the
// value or the fallback that was provided
func GetEnvKey(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetCurrentDirectory gets the directory in which the app was started and returns either
// the full directory or an error
func GetCurrentDirectory() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
