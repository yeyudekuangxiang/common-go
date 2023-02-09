package apitool

import (
	"net"
	"net/http"
	"strings"
)

var (
	// TrustedPlatform when running on Google App Engine. Trust X-Appengine-Remote-Addr
	// for determining the client's IP
	TrustedPlatform = "X-Appengine-Remote-Addr"

	RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	TrustedCIDRs    = []*net.IPNet{
		{ // 0.0.0.0/0 (IPv4)
			IP:   net.IP{0x0, 0x0, 0x0, 0x0},
			Mask: net.IPMask{0x0, 0x0, 0x0, 0x0},
		},
		{ // ::/0 (IPv6)
			IP:   net.IP{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			Mask: net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
	}
)

func ClientIp(r *http.Request) string {

	if TrustedPlatform != "" {
		if addr := r.Header.Get(TrustedPlatform); addr != "" {
			return addr
		}
	}

	// It also checks if the remoteIP is a trusted proxy or not.
	// In order to perform this validation, it will see if the IP is contained within at least one of the CIDR blocks
	// defined by Engine.SetTrustedProxies()
	remoteIP := net.ParseIP(RemoteIP(r))
	if remoteIP == nil {
		return ""
	}
	trusted := isTrustedProxy(remoteIP)

	if trusted && RemoteIPHeaders != nil {
		for _, headerName := range RemoteIPHeaders {
			ip, valid := validateHeader(r.Header.Get(headerName))
			if valid {
				return ip
			}
		}
	}
	return remoteIP.String()
}

// RemoteIP parses the IP from Request.RemoteAddr, normalizes and returns the IP (without the port).
func RemoteIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

// isTrustedProxy will check whether the IP address is included in the trusted list according to Engine.trustedCIDRs
func isTrustedProxy(ip net.IP) bool {
	for _, cidr := range TrustedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// validateHeader will parse X-Forwarded-For header and return the trusted client IP address
func validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, ",")
	for i := len(items) - 1; i >= 0; i-- {
		ipStr := strings.TrimSpace(items[i])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		// X-Forwarded-For is appended by proxy
		// Check IPs in reverse order and stop when find untrusted proxy
		if (i == 0) || (!isTrustedProxy(ip)) {
			return ipStr, true
		}
	}
	return "", false
}
