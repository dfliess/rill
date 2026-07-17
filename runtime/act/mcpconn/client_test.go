package mcpconn_test

import (
	"encoding/json"
	"testing"

	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/stretchr/testify/require"
)

// newTestClient wires a client to the in-process mock. The mock is plaintext http on loopback, the posture
// packaged/local development connectors use; private-range and http access is granted out-of-band via the Go
// trust flag (AllowPrivateNetworks), never through tenant config.
func newTestClient(t *testing.T, m *mockMCP, opts mcpconn.Options) *mcpconn.Client {
	t.Helper()
	cfg := &mcpconn.ConnectorConfig{
		Name:      "jira_mcp",
		Transport: mcpconn.TransportStreamableHTTP,
		URL:       m.url(),
	}
	opts.AllowPrivateNetworks = true

	c, err := mcpconn.NewClient(cfg, opts)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func toolByName(tools []mcpconn.RemoteTool, name string) (mcpconn.RemoteTool, bool) {
	for _, tt := range tools {
		if tt.Name == name {
			return tt, true
		}
	}
	return mcpconn.RemoteTool{}, false
}

func TestListToolsDiscoversAndNamespaces(t *testing.T) {
	m := newMockMCP(t)
	c := newTestClient(t, m, mcpconn.Options{})

	tools, err := c.ListTools(t.Context())
	require.NoError(t, err)
	require.Len(t, tools, 3)

	weather, ok := toolByName(tools, "mcp.jira_mcp.get_weather")
	require.True(t, ok, "effective name is namespaced as mcp.<connector>.<tool>")
	require.Equal(t, "get_weather", weather.RawName)
	require.Equal(t, "jira_mcp", weather.Connector)
	require.NotEmpty(t, weather.SchemaHash)
	require.NotEmpty(t, weather.InputSchema)
	// get_weather is announced read-only; the hint is surfaced (untrusted) for the gateway to weigh.
	require.True(t, weather.Hints.ReadOnly)

	ticket, ok := toolByName(tools, "mcp.jira_mcp.create_ticket")
	require.True(t, ok)
	require.False(t, ticket.Hints.ReadOnly)

	// Distinct tools hash differently: the pin identifies the specific schema, not just any schema.
	require.NotEqual(t, weather.SchemaHash, ticket.SchemaHash)
}

func TestCallToolReadOnly(t *testing.T) {
	m := newMockMCP(t)
	c := newTestClient(t, m, mcpconn.Options{})

	tools, err := c.ListTools(t.Context())
	require.NoError(t, err)
	weather, ok := toolByName(tools, "mcp.jira_mcp.get_weather")
	require.True(t, ok)

	res, err := c.CallTool(t.Context(), "mcp.jira_mcp.get_weather", weather.SchemaHash, json.RawMessage(`{"city":"Madrid"}`))
	require.NoError(t, err)
	require.False(t, res.IsError)

	var out struct {
		TempC      float64 `json:"temp_c"`
		Conditions string  `json:"conditions"`
	}
	require.NotEmpty(t, res.StructuredContent)
	require.NoError(t, json.Unmarshal(res.StructuredContent, &out))
	require.Equal(t, "sunny in Madrid", out.Conditions)
}

func TestCallToolWithoutPinIsRefused(t *testing.T) {
	m := newMockMCP(t)
	c := newTestClient(t, m, mcpconn.Options{})

	// An empty expected schema hash means the action carries no pin, so the call is refused rather than dispatched.
	_, err := c.CallTool(t.Context(), "mcp.jira_mcp.get_weather", "", json.RawMessage(`{"city":"Madrid"}`))
	require.ErrorIs(t, err, mcpconn.ErrToolNotDiscovered)
}

// TestSchemaChangeStopsCall covers the pinning guarantee: after discovery pins get_weather, the server is
// rebuilt with an incompatible schema for it. The next call to get_weather is stopped with ErrSchemaChanged
// and never dispatched, while a tool whose schema did not change still runs.
func TestSchemaChangeStopsCall(t *testing.T) {
	m := newMockMCP(t)
	c := newTestClient(t, m, mcpconn.Options{})

	tools, err := c.ListTools(t.Context())
	require.NoError(t, err)
	weather, ok := toolByName(tools, "mcp.jira_mcp.get_weather")
	require.True(t, ok)
	ticket, ok := toolByName(tools, "mcp.jira_mcp.create_ticket")
	require.True(t, ok)

	m.setBuild(changedWeatherBuild)

	// The pinned hash the action carries no longer matches the live schema, so the call is stopped.
	_, err = c.CallTool(t.Context(), "mcp.jira_mcp.get_weather", weather.SchemaHash, json.RawMessage(`{"city":"Madrid"}`))
	require.ErrorIs(t, err, mcpconn.ErrSchemaChanged)

	// create_ticket was not changed, so its pin still matches and it dispatches normally.
	res, err := c.CallTool(t.Context(), "mcp.jira_mcp.create_ticket", ticket.SchemaHash, json.RawMessage(`{"title":"hi","body":"x"}`))
	require.NoError(t, err)
	require.False(t, res.IsError)
}

// TestPayloadLimitRejected covers the response cap: a tool that returns more than MaxPayloadBytes fails the
// call. tools/list stays well under the cap, so discovery and the pin re-check succeed and only the oversized
// result trips the guard.
func TestPayloadLimitRejected(t *testing.T) {
	m := newMockMCP(t)
	c := newTestClient(t, m, mcpconn.Options{MaxPayloadBytes: 32 * 1024})

	tools, err := c.ListTools(t.Context())
	require.NoError(t, err)
	emit, ok := toolByName(tools, "mcp.jira_mcp.emit")
	require.True(t, ok)

	_, err = c.CallTool(t.Context(), "mcp.jira_mcp.emit", emit.SchemaHash, json.RawMessage(`{"size":524288}`))
	require.Error(t, err)
}

// TestClientBlocksPrivateRangeHost exercises the end-to-end egress guard: a valid connector whose host
// resolves into a denied range never establishes a session, so discovery fails closed.
func TestClientBlocksPrivateRangeHost(t *testing.T) {
	cfg, err := mcpconn.ParseConnectorConfig("blocked", map[string]any{
		"transport": "streamable_http",
		"url":       "https://10.0.0.1:9443/mcp", // RFC1918, https so config validation passes; dial is blocked
	})
	require.NoError(t, err)

	c, err := mcpconn.NewClient(cfg, mcpconn.Options{})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	_, err = c.ListTools(t.Context())
	require.Error(t, err)
	// The SDK's transport wraps the dial error without preserving its identity via %w, so match on the
	// guard's message rather than errors.Is here. The sentinel identity itself is asserted in
	// TestCheckDialAddress, which calls the dial guard directly.
	require.ErrorContains(t, err, "blocked by egress policy")
}
