package service

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid short", "abc", false, ""},
		{"valid long", "abcdefghijklmnopqrstuvwx", false, ""},
		{"valid with numbers", "user123", false, ""},
		{"valid with underscore", "user_name", false, ""},
		{"valid with hyphen", "user-name", false, ""},
		{"valid Chinese", "测试用户", false, ""},
		{"valid mixed", "user123_测试", false, ""},
		{"too short", "ab", true, "username must be at least 3 characters"},
		{"too long", "abcdefghijklmnopqrstuvwxy", true, "username must be at most 24 characters"},
		{"empty", "", true, "username must be at least 3 characters"},
		{"with space", "user name", true, "username can only contain"},
		{"with special char", "user@name", true, "username can only contain"},
		{"reserved admin", "admin", true, "this username is reserved"},
		{"reserved root", "root", true, "this username is reserved"},
		{"reserved system", "system", true, "this username is reserved"},
		{"reserved guest", "guest", true, "this username is reserved"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateUsername(%q) error = %v, want err containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid strong", "Password123", false, ""},
		{"valid with special", "Pass123!@#", false, ""},
		{"too short", "Pass12", true, "password must be at least 8 characters"},
		{"empty", "", true, "password must be at least 8 characters"},
		{"no uppercase", "password123", true, "password must contain uppercase, lowercase letters and numbers"},
		{"no lowercase", "PASSWORD123", true, "password must contain uppercase, lowercase letters and numbers"},
		{"no digit", "Password", true, "password must contain uppercase, lowercase letters and numbers"},
		{"only letters", "Password", true, "password must contain uppercase, lowercase letters and numbers"},
		{"too long", "", true, ""},
	}

	longPass := make([]byte, 129)
	for i := range longPass {
		longPass[i] = 'a'
	}
	longPass[0] = 'A'
	longPass[1] = '1'
	tests[len(tests)-1].input = string(longPass)
	tests[len(tests)-1].errMsg = "password must be at most 128 characters"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidatePassword(%q) error = %v, want err containing %q", tt.input, err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "test@example.com", false},
		{"valid with dots", "test.user@example.com", false},
		{"valid with plus", "test+tag@example.com", false},
		{"valid subdomain", "test@mail.example.com", false},
		{"valid hyphen domain", "test@example-site.com", false},
		{"empty", "", true},
		{"no @", "testexample.com", true},
		{"no domain", "test@", true},
		{"no local", "@example.com", true},
		{"space", "test @example.com", true},
		{"special char", "test!@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
