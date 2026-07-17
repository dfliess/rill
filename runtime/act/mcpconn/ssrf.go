package mcpconn

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

// ErrBlockedHost is returned when a request targets a host outside the connector's allowlist, or when
// its resolved address falls in a denied range. Both checks run BEFORE the request bytes leave the client
// (the host/scheme check precedes the base round trip; the dial guard runs during connect), so it is a
// pre-send signal: the original request never reached the server. It is deliberately a sentinel so the
// gateway can tell a pre-send egress refusal apart from a generic transport failure.
var ErrBlockedHost = errors.New("mcpconn: host blocked by egress policy")

// ErrRedirectBlocked is returned when a response is a redirect the guard refuses to follow. Unlike
// ErrBlockedHost it is NOT pre-send: the original request already reached the server and returned a 3xx
// before the client refused the hop, so a write may already have landed. The gateway must classify it as
// indeterminate, never as a definitive non-write.
var ErrRedirectBlocked = errors.New("mcpconn: redirect blocked by egress policy")

// ErrPayloadTooLarge is returned when a response body exceeds the configured maximum. Like ErrBlockedHost
// it is a sentinel so an over-limit response reads as a guard trip, not a malformed server.
var ErrPayloadTooLarge = errors.New("mcpconn: response payload exceeds limit")

// guardedHTTPClient builds an *http.Client whose every request is subject to the connector's egress
// guard: the host allowlist and scheme are checked per request (including redirect hops), the resolved
// IP is checked at dial time so DNS-rebinding cannot slip a private address past the hostname check, and
// each response body is capped at maxPayloadBytes. Redirects are refused outright to keep the reachable
// surface equal to the validated endpoint. bearerToken, if non-empty, is attached inside the guard so it
// is only set on a request that has already passed the host/scheme check.
func guardedHTTPClient(cfg NetworkConfig, allowPrivate bool, dialTimeout time.Duration, maxPayloadBytes int64, bearerToken string) *http.Client {
	dialer := &net.Dialer{
		Timeout: dialTimeout,
		// Control runs after DNS resolution with the concrete address about to be dialed, so it sees the
		// real IP even if the hostname resolved to a private one. This is the anti-rebinding checkpoint.
		Control: func(_, address string, _ syscall.RawConn) error {
			return checkDialAddress(address, allowPrivate)
		},
	}

	base := &http.Transport{
		DialContext:           dialer.DialContext,
		Proxy:                 nil, // no proxy: dial directly so Control always sees the true destination IP
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   dialTimeout,
		ExpectContinueTimeout: time.Second,
	}

	// Auth sits below the guard: the guarded round tripper only calls into it once the host/scheme check has
	// passed, so the credential is never attached to a request bound for a blocked host.
	var inner http.RoundTripper = base
	if bearerToken != "" {
		inner = &authRoundTripper{token: bearerToken, base: inner}
	}

	return &http.Client{
		Transport: &guardedRoundTripper{
			base:            inner,
			allowedHosts:    cfg.AllowedHosts,
			allowPrivate:    allowPrivate,
			maxPayloadBytes: maxPayloadBytes,
		},
		// Refuse every redirect: an MCP endpoint served over Streamable HTTP has no reason to redirect,
		// and a redirect is a classic way to bounce a request onto an unvalidated host. This fires only
		// after the original request got a 3xx response, so it is ErrRedirectBlocked (post-send), not
		// ErrBlockedHost: the gateway treats it as indeterminate because the first request may have written.
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return fmt.Errorf("%w: redirects are not allowed", ErrRedirectBlocked)
		},
		Jar:     nil,
		Timeout: 0, // per-operation deadlines are applied via context; a client-wide timeout would also kill streams
	}
}

// guardedRoundTripper enforces the host allowlist and scheme on the way out and the payload cap on the
// way in. It wraps the real transport rather than replacing it so redirect hops (each a fresh RoundTrip)
// are validated individually.
type guardedRoundTripper struct {
	base            http.RoundTripper
	allowedHosts    []string
	allowPrivate    bool
	maxPayloadBytes int64
}

func (g *guardedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := g.checkRequestHost(req); err != nil {
		return nil, err
	}

	resp, err := g.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Cap the body so a hostile or runaway server cannot stream an unbounded response into the worker.
	if g.maxPayloadBytes > 0 && resp.Body != nil {
		resp.Body = &limitedReadCloser{rc: resp.Body, remaining: g.maxPayloadBytes + 1}
	}
	return resp, nil
}

// checkRequestHost validates the scheme and hostname of an outgoing request before it is dialed.
func (g *guardedRoundTripper) checkRequestHost(req *http.Request) error {
	scheme := strings.ToLower(req.URL.Scheme)
	if scheme != "https" && !(scheme == "http" && g.allowPrivate) {
		return fmt.Errorf("%w: scheme %q not allowed", ErrBlockedHost, req.URL.Scheme)
	}
	if len(g.allowedHosts) > 0 && !hostAllowed(req.URL.Hostname(), g.allowedHosts) {
		return fmt.Errorf("%w: host %q not in allowlist", ErrBlockedHost, req.URL.Hostname())
	}
	return nil
}

// checkDialAddress rejects a dial target whose IP falls in a denied range unless private ranges are
// allowed for this connector. address is "host:port"; the host part is expected to already be an IP
// because the dialer resolves before invoking Control.
func checkDialAddress(address string, allowPrivate bool) error {
	if allowPrivate {
		return nil
	}
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("%w: cannot parse dial address %q", ErrBlockedHost, address)
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("%w: dial address %q is not an IP", ErrBlockedHost, host)
	}
	if isBlockedIP(ip) {
		return fmt.Errorf("%w: %s is in a denied range", ErrBlockedHost, ip)
	}
	return nil
}

// specialPurposeCIDRs are non-public / special-purpose prefixes a Tenant connector must never reach, beyond
// the ranges the net.IP predicate methods already cover. They are IANA special-purpose registrations that
// are nonetheless internally routable or reserved, so a deny posture must exclude them too. Parsed once at
// init; the list is a constant, so a parse failure is a programmer error.
//
// This deny-set is DEFENSE IN DEPTH, not the primary egress control: the primary control is the per-connector
// allowed_hosts allowlist (NetworkConfig.AllowedHosts), which pins reachable hosts to an explicit set. The deny-set
// exists so a connector configured without an allowlist, or one whose allowed host resolves (via DNS rebinding or
// misconfiguration) to a non-public address, still cannot reach internal or reserved space. It mirrors the IANA
// IPv4 and IPv6 special-purpose registries for prefixes the net.IP predicates do not already cover.
var specialPurposeCIDRs = mustParseCIDRs(
	"0.0.0.0/8",       // RFC1122 "this host on this network" (IsUnspecified only covers 0.0.0.0/32)
	"100.64.0.0/10",   // RFC6598 carrier-grade NAT
	"192.0.0.0/24",    // RFC6890 IETF protocol assignments
	"192.0.2.0/24",    // RFC5737 TEST-NET-1
	"192.88.99.0/24",  // RFC7526 6to4 relay anycast (deprecated)
	"198.18.0.0/15",   // RFC2544 benchmarking
	"198.51.100.0/24", // RFC5737 TEST-NET-2
	"203.0.113.0/24",  // RFC5737 TEST-NET-3
	"240.0.0.0/4",     // RFC1112 reserved (former Class E), also covers 255.255.255.255 limited broadcast
	"2001::/32",       // RFC4380 Teredo tunneling
	"2001:2::/48",     // RFC5180 IPv6 benchmarking
	"2001:db8::/32",   // RFC3849 IPv6 documentation
	"2002::/16",       // RFC3056 6to4
	"3fff::/20",       // RFC9637 IPv6 documentation
	"64:ff9b::/96",    // RFC6052 NAT64 well-known prefix
	"64:ff9b:1::/48",  // RFC8215 NAT64 local-use
	"100::/64",        // RFC6666 discard-only address block
)

// isBlockedIP reports whether ip is one that a Tenant connector must never reach. It covers loopback
// (127.0.0.0/8, ::1), RFC1918 private (10/8, 172.16/12, 192.168/16) and unique-local (fc00::/7), link-local
// (169.254/16, fe80::/10), the unspecified address (0.0.0.0, ::) and multicast via the net.IP predicates,
// plus the special-purpose CIDRs above (this-network 0.0.0.0/8, CGNAT, benchmarking, protocol assignments, the
// TEST-NET and 6to4/Teredo/NAT64/discard blocks, reserved space and IPv6 documentation). It is defense in depth
// behind the allowed_hosts allowlist, which is the primary egress control.
func isBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	// Normalise IPv4-mapped IPv6 (::ffff:a.b.c.d) so the IPv4 predicates apply to it as well.
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}
	if ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified() ||
		ip.IsInterfaceLocalMulticast() {
		return true
	}
	for _, n := range specialPurposeCIDRs {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// mustParseCIDRs parses a fixed list of CIDR prefixes, panicking on a malformed entry. The input is a
// package constant, so a failure is a programming error surfaced at init, not a runtime condition.
func mustParseCIDRs(cidrs ...string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		_, n, err := net.ParseCIDR(c)
		if err != nil {
			panic(fmt.Sprintf("mcpconn: invalid special-purpose CIDR %q: %v", c, err))
		}
		nets = append(nets, n)
	}
	return nets
}

// limitedReadCloser caps the number of bytes read from a response body. remaining starts at limit+1, so
// reading the (limit+1)-th byte trips ErrPayloadTooLarge while a body of exactly limit bytes passes.
type limitedReadCloser struct {
	rc        io.ReadCloser
	remaining int64
}

func (l *limitedReadCloser) Read(p []byte) (int, error) {
	if l.remaining <= 0 {
		return 0, ErrPayloadTooLarge
	}
	if int64(len(p)) > l.remaining {
		p = p[:l.remaining]
	}
	n, err := l.rc.Read(p)
	l.remaining -= int64(n)
	if l.remaining <= 0 && err == nil {
		// We have now consumed limit+1 bytes without the underlying reader signalling EOF, so the body is
		// over the limit.
		err = ErrPayloadTooLarge
	}
	return n, err
}

func (l *limitedReadCloser) Close() error {
	return l.rc.Close()
}

// authRoundTripper injects the bearer credential resolved by the caller. It sits below the egress guard so
// the credential is only attached once the host/scheme check has passed.
type authRoundTripper struct {
	token string
	base  http.RoundTripper
}

func (a *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone before mutating: RoundTrippers must not modify the caller's request.
	r := req.Clone(req.Context())
	r.Header.Set("Authorization", "Bearer "+a.token)
	return a.base.RoundTrip(r)
}
