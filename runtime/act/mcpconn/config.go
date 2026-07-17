// Package mcpconn implements Act's outbound MCP role: it parses and validates an MCP connector
// definition, connects to a remote MCP server over Streamable HTTP, discovers its tools with a
// pinned schema hash, and enforces egress protection (host allowlist + private-range deny) and
// payload limits on every request. The internal MCP server that the fork already exposes is a
// separate concern (see runtime/ai/mcp.go); this package is only ever the client.
//
// The MCPClient interface here is the seam the Tool Gateway consumes: the gateway owns policy,
// identity and approval; this package owns transport, discovery, schema pinning and network
// safety. Nothing here decides whether a tool may run, only whether it can be reached and whether
// its schema still matches what was discovered.
package mcpconn

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Transport values accepted in a connector definition. Only Streamable HTTP over HTTPS is admitted
// for v1: stdio would let a Tenant hand an arbitrary command to the shared worker, so it is rejected
// at parse time rather than filtered later.
const (
	TransportStreamableHTTP = "streamable_http"
)

// PrivateRangesDeny is the only value a connector definition may carry for network.private_ranges: a
// Tenant connector always reaches a public, allowlisted host. Reaching loopback, private or special-purpose
// ranges is NOT a tenant-settable option; it is an out-of-band build-time trust decision expressed as a Go
// flag (Options.AllowPrivateNetworks) on client construction, so a Tenant-supplied definition can never opt
// itself out of the egress protection. Only packaged/local development connectors and the in-process test
// server, constructed in trusted Go code, set that flag.
const PrivateRangesDeny = "deny"

// ConnectorConfig is the parsed and validated form of an MCP connector definition. It mirrors the
// illustrative YAML in the design (§16.3): transport, url, an auth secret reference, and a network
// block. The struct is deliberately a plain value with a Validate method rather than a registered
// Rill driver: wiring a formal driver would mean touching the driver registry, and for v1 a Config
// type owned by the act package is enough for the gateway to consume. See the package README for the
// rationale.
type ConnectorConfig struct {
	// Name is the connector's local name, used to namespace effective tool names as mcp.<name>.<tool>.
	// It is not decoded from the property map: it comes from the resource name, so ParseConnectorConfig
	// takes it as a separate argument.
	Name string `mapstructure:"-"`

	// Transport must be TransportStreamableHTTP. Kept as a string (not a bool/enum) so an unknown value
	// surfaces as a clear validation error instead of silently defaulting.
	Transport string `mapstructure:"transport"`

	// URL is the remote MCP endpoint. It must be absolute and, unless private ranges are allowed, https.
	URL string `mapstructure:"url"`

	// Auth carries the reference to the credential, never the credential itself. The value is resolved
	// server-side by the caller (the gateway/executor) and injected as a bearer token; this package
	// never reads the secret manager.
	Auth AuthConfig `mapstructure:"auth"`

	// Network holds the egress guard configuration.
	Network NetworkConfig `mapstructure:"network"`

	// TrustReadOnlyHint opts the connector in to treating a tool's server-advertised readOnlyHint as
	// meaningful. It changes nothing in this package: discovery always surfaces the hint as untrusted
	// data on RemoteTool. The flag is carried so the gateway can decide, per §16.3, whether a
	// read-only-hinted tool may auto-execute. Off by default (fail-closed).
	TrustReadOnlyHint bool `mapstructure:"trust_read_only_hint"`
}

// AuthConfig references the credential to present to the remote server.
type AuthConfig struct {
	// Secret is the name of the secret in the platform secret manager. It is a reference only.
	Secret string `mapstructure:"secret"`
}

// NetworkConfig configures the egress guard applied to every request the client makes.
type NetworkConfig struct {
	// AllowedHosts is the egress allowlist of hostnames. When non-empty, every request host (including
	// redirect targets) must match one entry exactly, case-insensitively. When empty, no host allowlist
	// is enforced, but the private-range deny below still applies to the resolved IP.
	AllowedHosts []string `mapstructure:"allowed_hosts"`

	// PrivateRanges may only be PrivateRangesDeny (the default when empty). Any other value — including the
	// old "allow" — is rejected: private-range access is a build-time trust decision (Options.
	// AllowPrivateNetworks), never a tenant config knob. The field is kept so a definition can state the
	// deny posture explicitly, matching the design YAML.
	PrivateRanges string `mapstructure:"private_ranges"`
}

// ParseConnectorConfig decodes a connector property map into a ConnectorConfig, applies defaults and
// validates it. name is the connector's resource name (used for namespacing); props is the decoded
// YAML body of the connector definition.
func ParseConnectorConfig(name string, props map[string]any) (*ConnectorConfig, error) {
	cfg := &ConnectorConfig{
		Name:              name,
		Transport:         "",
		URL:               "",
		Auth:              AuthConfig{Secret: ""},
		Network:           NetworkConfig{AllowedHosts: nil, PrivateRanges: ""},
		TrustReadOnlyHint: false,
	}
	if err := mapstructure.WeakDecode(props, cfg); err != nil {
		return nil, fmt.Errorf("mcpconn: decode connector config: %w", err)
	}
	// Name is decoded as ignored, so restore it after WeakDecode would have zeroed it.
	cfg.Name = name
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate rejects any tenant-supplied definition that is not safe to connect to. It is the strict,
// fail-closed path ParseConnectorConfig uses: https is mandatory, private_ranges may only be deny, and it
// never grants private-range or plaintext access. Trusted construction that legitimately reaches a loopback
// endpoint goes through NewClient with Options.AllowPrivateNetworks, which relaxes only the scheme check.
func (c *ConnectorConfig) Validate() error {
	return c.validateWith(false)
}

// validateWith applies defaults and validates the definition, relaxing ONLY the https requirement when
// allowPrivate is set (a build-time trust flag, never a tenant knob). private_ranges is still required to be
// deny regardless: the private-range gating lives entirely in the Go flag, so a connector definition can
// never carry an opt-out.
func (c *ConnectorConfig) validateWith(allowPrivate bool) error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("mcpconn: connector name is required")
	}

	if c.Transport != TransportStreamableHTTP {
		return fmt.Errorf("mcpconn: unsupported transport %q: only %q is allowed in v1", c.Transport, TransportStreamableHTTP)
	}

	// private_ranges is fixed to deny: it is not a tenant-settable trust knob. An empty value defaults to
	// deny; anything else — including the old "allow" — is rejected fail-closed.
	if c.Network.PrivateRanges == "" {
		c.Network.PrivateRanges = PrivateRangesDeny
	}
	if c.Network.PrivateRanges != PrivateRangesDeny {
		return fmt.Errorf("mcpconn: network.private_ranges %q is not settable; only %q is allowed (private-range access is a build-time trust flag, not tenant config)", c.Network.PrivateRanges, PrivateRangesDeny)
	}

	u, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("mcpconn: parse url: %w", err)
	}
	if !u.IsAbs() || u.Host == "" {
		return fmt.Errorf("mcpconn: url must be an absolute http(s) URL, got %q", c.URL)
	}
	// HTTPS is mandatory for a Tenant connector. Plaintext http is permitted only under the trusted
	// AllowPrivateNetworks flag, for packaged/local development connectors and the in-process test server
	// where a loopback endpoint is legitimate.
	if u.Scheme != "https" {
		if !(u.Scheme == "http" && allowPrivate) {
			return fmt.Errorf("mcpconn: url scheme must be https, got %q", u.Scheme)
		}
	}

	// When an allowlist is set, the connector's own URL host must be in it, otherwise the connector could
	// never reach its own endpoint and the misconfiguration is better caught now than at first call.
	if len(c.Network.AllowedHosts) > 0 && !hostAllowed(u.Hostname(), c.Network.AllowedHosts) {
		return fmt.Errorf("mcpconn: url host %q is not in network.allowed_hosts", u.Hostname())
	}

	return nil
}

// hostAllowed reports whether host matches any entry in allowed, comparing case-insensitively. An
// empty allowlist means "not enforced" and is handled by callers, so this returns false for it.
func hostAllowed(host string, allowed []string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	return slices.ContainsFunc(allowed, func(h string) bool {
		return strings.ToLower(strings.TrimSpace(h)) == host
	})
}
