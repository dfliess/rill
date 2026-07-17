package mcpconn_test

import (
	"net"
	"testing"

	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/stretchr/testify/require"
)

func TestBlockedIP(t *testing.T) {
	cases := []struct {
		ip      string
		blocked bool
	}{
		{"127.0.0.1", true},             // loopback
		{"10.0.0.1", true},              // RFC1918
		{"172.16.5.4", true},            // RFC1918
		{"172.31.255.255", true},        // RFC1918 upper bound
		{"192.168.1.1", true},           // RFC1918
		{"169.254.10.10", true},         // link-local
		{"0.0.0.0", true},               // unspecified
		{"224.0.0.1", true},             // multicast
		{"::1", true},                   // IPv6 loopback
		{"fe80::1", true},               // IPv6 link-local
		{"fc00::1", true},               // IPv6 unique-local
		{"::", true},                    // IPv6 unspecified
		{"::ffff:10.0.0.1", true},       // IPv4-mapped private
		{"100.64.0.1", true},            // CGNAT 100.64/10
		{"100.127.255.255", true},       // CGNAT upper bound
		{"192.0.0.1", true},             // IETF protocol assignments 192.0.0/24
		{"192.0.2.5", true},             // TEST-NET-1
		{"198.18.0.5", true},            // benchmarking 198.18/15
		{"198.51.100.9", true},          // TEST-NET-2
		{"203.0.113.9", true},           // TEST-NET-3
		{"240.0.0.1", true},             // reserved (former Class E)
		{"0.0.0.1", true},               // "this host on this network" 0.0.0.0/8 (not just the unspecified /32)
		{"0.255.255.255", true},         // 0.0.0.0/8 upper bound
		{"192.88.99.1", true},           // 6to4 relay anycast
		{"2001:db8::1", true},           // IPv6 documentation
		{"64:ff9b::808:808", true},      // NAT64 well-known prefix
		{"64:ff9b:1::1", true},          // NAT64 local-use
		{"100::1", true},                // discard-only address block
		{"2001::1", true},               // Teredo
		{"2001:2::1", true},             // IPv6 benchmarking
		{"2002::1", true},               // 6to4
		{"3fff::1", true},               // IPv6 documentation (RFC9637)
		{"172.15.0.1", false},           // just below RFC1918 172.16/12
		{"172.32.0.1", false},           // just above RFC1918 172.16/12
		{"100.63.255.255", false},       // just below CGNAT 100.64/10
		{"100.128.0.1", false},          // just above CGNAT 100.64/10
		{"1.0.0.1", false},              // public, just above 0.0.0.0/8
		{"8.8.8.8", false},              // public
		{"1.1.1.1", false},              // public
		{"2001:4860:4860::8888", false}, // public IPv6 (outside 2001::/32 Teredo)
		{"2606:4700:4700::1111", false}, // public IPv6
	}
	for _, c := range cases {
		ip := net.ParseIP(c.ip)
		require.NotNil(t, ip, "parse %s", c.ip)
		require.Equalf(t, c.blocked, mcpconn.BlockedIP(ip), "ip %s", c.ip)
	}
}

func TestCheckDialAddress(t *testing.T) {
	// Under deny, a private address is refused with ErrBlockedHost.
	err := mcpconn.CheckDialAddress("10.0.0.1:443", false)
	require.ErrorIs(t, err, mcpconn.ErrBlockedHost)

	// A public address passes.
	require.NoError(t, mcpconn.CheckDialAddress("8.8.8.8:443", false))

	// When private ranges are allowed, loopback passes (dev/local and the in-process mock).
	require.NoError(t, mcpconn.CheckDialAddress("127.0.0.1:443", true))
}
