package mcpconn_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/stretchr/testify/require"
)

func TestParseConnectorConfig_ValidHTTPS(t *testing.T) {
	cfg, err := mcpconn.ParseConnectorConfig("jira_mcp", map[string]any{
		"transport": "streamable_http",
		"url":       "https://mcp.example.net/mcp",
		"auth":      map[string]any{"secret": "jira_mcp_oauth"},
		"network": map[string]any{
			"allowed_hosts": []string{"mcp.example.net"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "jira_mcp", cfg.Name)
	require.Equal(t, mcpconn.TransportStreamableHTTP, cfg.Transport)
	require.Equal(t, "jira_mcp_oauth", cfg.Auth.Secret)
	// private_ranges defaults to deny (fail-closed) when not set.
	require.Equal(t, mcpconn.PrivateRangesDeny, cfg.Network.PrivateRanges)
	require.False(t, cfg.TrustReadOnlyHint)
}

func TestParseConnectorConfig_RejectsUnknownTransport(t *testing.T) {
	_, err := mcpconn.ParseConnectorConfig("c", map[string]any{
		"transport": "stdio",
		"url":       "https://mcp.example.net/mcp",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "transport")
}

func TestParseConnectorConfig_RejectsPlainHTTPUnderDeny(t *testing.T) {
	// http is only permitted when private ranges are allowed (dev/local); under the default deny a Tenant
	// connector must be https.
	_, err := mcpconn.ParseConnectorConfig("c", map[string]any{
		"transport": "streamable_http",
		"url":       "http://mcp.example.net/mcp",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "https")
}

func TestParseConnectorConfig_RejectsPrivateRangesAllow(t *testing.T) {
	// "allow" is no longer a tenant-settable value: private-range access is a build-time Go trust flag
	// (Options.AllowPrivateNetworks), not connector config. A definition asking for it is rejected fail-closed.
	_, err := mcpconn.ParseConnectorConfig("c", map[string]any{
		"transport": "streamable_http",
		"url":       "https://mcp.example.net/mcp",
		"network":   map[string]any{"private_ranges": "allow"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "private_ranges")
}

func TestParseConnectorConfig_RejectsHostOutsideAllowlist(t *testing.T) {
	_, err := mcpconn.ParseConnectorConfig("c", map[string]any{
		"transport": "streamable_http",
		"url":       "https://evil.example.com/mcp",
		"network":   map[string]any{"allowed_hosts": []string{"mcp.example.net"}},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "allowed_hosts")
}

func TestParseConnectorConfig_RejectsBadPrivateRanges(t *testing.T) {
	_, err := mcpconn.ParseConnectorConfig("c", map[string]any{
		"transport": "streamable_http",
		"url":       "https://mcp.example.net/mcp",
		"network":   map[string]any{"private_ranges": "maybe"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "private_ranges")
}

func TestParseConnectorConfig_RejectsRelativeURL(t *testing.T) {
	_, err := mcpconn.ParseConnectorConfig("c", map[string]any{
		"transport": "streamable_http",
		"url":       "/mcp",
	})
	require.Error(t, err)
}
