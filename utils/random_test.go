package utils_test

import (
	"strings"
	"testing"

	"github.com/codeflames/cizzors/utils"
)

func TestRandomUrl(t *testing.T) {
	url := utils.RandomUrl()

	// Check that the URL has the correct length
	if len(url) != 6 {
		t.Errorf("Expected URL length of 6, but got %d", len(url))
	}

	// Check that the URL only contains valid characters
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range url {
		if !strings.ContainsRune(validChars, char) {
			t.Errorf("Invalid character '%c' found in URL", char)
		}
	}
}
