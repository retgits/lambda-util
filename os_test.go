// Package util implements utility methods
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvKeyExistingKey(t *testing.T) {
	result := GetEnvKey("PATH", "fallback")
	assert.NotEqual(t, result, "fallback")
	assert.NotEmpty(t, result)
}

func TestGetEnvKeyNonExistingKey(t *testing.T) {
	result := GetEnvKey("SOMEPATH", "fallback")
	assert.Equal(t, result, "fallback")
	assert.NotEmpty(t, result)
}

func TestGetCurrentDirectory(t *testing.T) {
	result, err := GetCurrentDirectory()
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
}
