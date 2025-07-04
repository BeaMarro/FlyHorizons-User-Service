package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIPAddress(r *http.Request) string {
	// Prefer the X-Forwarded-For header if it's present
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Handle comma-separated list if behind multiple proxies
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fall back to X-Real-IP header (used by some proxies)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	// Normalize IPv6 loopback to IPv4 for clarity
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	return ip
}
