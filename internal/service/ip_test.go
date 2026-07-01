package service

import (
	"net/http"
	"testing"
)

func TestGetIPFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		xRealIP        string
		trustedProxies []string
		expected       string
	}{
		{
			name:           "no trusted proxies, no headers",
			remoteAddr:     "192.168.1.1:12345",
			trustedProxies: nil,
			expected:       "192.168.1.1",
		},
		{
			name:           "no trusted proxies, ignore X-Forwarded-For",
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "10.0.0.1",
			trustedProxies: nil,
			expected:       "192.168.1.1",
		},
		{
			name:           "no trusted proxies, ignore X-Real-IP",
			remoteAddr:     "192.168.1.1:12345",
			xRealIP:        "10.0.0.1",
			trustedProxies: nil,
			expected:       "192.168.1.1",
		},
		{
			name:           "trusted proxy, use X-Forwarded-For",
			remoteAddr:     "10.0.0.1:12345",
			xForwardedFor:  "192.168.1.1, 10.0.0.2",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "192.168.1.1",
		},
		{
			name:           "trusted proxy, use X-Real-IP",
			remoteAddr:     "10.0.0.1:12345",
			xRealIP:        "192.168.1.1",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "192.168.1.1",
		},
		{
			name:           "trusted proxy, X-Forwarded-For takes precedence",
			remoteAddr:     "10.0.0.1:12345",
			xForwardedFor:  "192.168.1.1",
			xRealIP:        "192.168.1.2",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "192.168.1.1",
		},
		{
			name:           "trusted proxy with wildcard",
			remoteAddr:     "10.0.0.1:12345",
			xForwardedFor:  "192.168.1.1",
			trustedProxies: []string{"*"},
			expected:       "192.168.1.1",
		},
		{
			name:           "untrusted proxy, ignore headers",
			remoteAddr:     "10.0.0.2:12345",
			xForwardedFor:  "192.168.1.1",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "10.0.0.2",
		},
		{
			name:           "IPv6 remote address",
			remoteAddr:     "[::1]:12345",
			trustedProxies: nil,
			expected:       "::1",
		},
		{
			name:           "X-Forwarded-For with spaces",
			remoteAddr:     "10.0.0.1:12345",
			xForwardedFor:  "  192.168.1.1  ,  10.0.0.2  ",
			trustedProxies: []string{"10.0.0.1"},
			expected:       "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			result := GetIPFromRequest(req, tt.trustedProxies)
			if result != tt.expected {
				t.Errorf("GetIPFromRequest() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIsTrustedProxy(t *testing.T) {
	tests := []struct {
		name           string
		ip             string
		trustedProxies []string
		expected       bool
	}{
		{"empty trusted list", "192.168.1.1", nil, false},
		{"ip in trusted list", "192.168.1.1", []string{"192.168.1.1", "10.0.0.1"}, true},
		{"ip not in trusted list", "192.168.1.2", []string{"192.168.1.1"}, false},
		{"wildcard trusted", "192.168.1.1", []string{"*"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTrustedProxy(tt.ip, tt.trustedProxies)
			if result != tt.expected {
				t.Errorf("isTrustedProxy(%q, %v) = %v, want %v", tt.ip, tt.trustedProxies, result, tt.expected)
			}
		})
	}
}
