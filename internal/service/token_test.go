package service

import (
	"testing"
)

func TestGenerateRandomToken(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 16", 16},
		{"length 32", 32},
		{"length 48", 48},
		{"length 64", 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := generateRandomToken(tt.length)
			if len(token) != tt.length {
				t.Errorf("generateRandomToken(%d) length = %d, want %d", tt.length, len(token), tt.length)
			}

			for _, c := range token {
				if !isValidTokenChar(c) {
					t.Errorf("generateRandomToken(%d) contains invalid char: %c", tt.length, c)
				}
			}
		})
	}
}

func TestGenerateRandomTokenUniqueness(t *testing.T) {
	tokens := make(map[string]bool)
	n := 100
	for i := 0; i < n; i++ {
		token := generateRandomToken(32)
		if tokens[token] {
			t.Errorf("duplicate token generated: %s", token)
		}
		tokens[token] = true
	}
	if len(tokens) != n {
		t.Errorf("expected %d unique tokens, got %d", n, len(tokens))
	}
}

func isValidTokenChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}
