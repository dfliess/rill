package mcpconn

import "net"

// BlockedIP exposes the internal SSRF address classifier to the external test package so the many
// denied ranges can be asserted directly, without standing up a server per range.
func BlockedIP(ip net.IP) bool {
	return isBlockedIP(ip)
}

// CheckDialAddress exposes the dial-time guard so a test can assert an address:port is refused.
func CheckDialAddress(address string, allowPrivate bool) error {
	return checkDialAddress(address, allowPrivate)
}
