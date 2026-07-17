package mcpconn_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// mockMCP is a benign in-process MCP server used by the tests. It exposes a small, fixed tool set over
// Streamable HTTP via httptest, and can swap the schema of a tool at runtime so the schema-pinning path
// can be exercised: because NewStreamableHTTPHandler asks for a fresh server per request, mutating the
// build closure changes what the next tools/list observes.
type mockMCP struct {
	server *httptest.Server

	mu    sync.Mutex
	build func(s *mcp.Server)
}

// weatherIn/weatherOut back get_weather: a read-only tool with a stable schema.
type weatherIn struct {
	City string `json:"city" jsonschema:"the city to look up"`
}

type weatherOut struct {
	TempC      float64 `json:"temp_c"`
	Conditions string  `json:"conditions"`
}

// ticketIn/ticketOut back create_ticket: a write tool (no read-only hint).
type ticketIn struct {
	Title string `json:"title" jsonschema:"the ticket title"`
	Body  string `json:"body" jsonschema:"the ticket body"`
}

type ticketOut struct {
	ID string `json:"id"`
}

// emitIn/emitOut back emit: a read-only tool that returns a payload of a caller-chosen size, used to trip
// the response payload cap.
type emitIn struct {
	Size int `json:"size" jsonschema:"number of bytes to return"`
}

type emitOut struct {
	Data string `json:"data"`
}

// weatherInV2 is a schema-incompatible variant of weatherIn: it renames the field, so a server rebuilt with
// it produces a different input schema and therefore a different pin hash.
type weatherInV2 struct {
	Location string `json:"location" jsonschema:"the location to look up"`
}

func boolPtr(b bool) *bool { return &b }

// defaultBuild registers the standard tool set: a read-only tool, a write tool, and a payload emitter.
func defaultBuild(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_weather",
		Description: "Return the current weather for a city.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: boolPtr(true)},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in weatherIn) (*mcp.CallToolResult, weatherOut, error) {
		return nil, weatherOut{TempC: 21.5, Conditions: "sunny in " + in.City}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_ticket",
		Description: "Create a ticket.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(false)},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in ticketIn) (*mcp.CallToolResult, ticketOut, error) {
		return nil, ticketOut{ID: "TCK-1:" + in.Title}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "emit",
		Description: "Return a payload of the requested size.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in emitIn) (*mcp.CallToolResult, emitOut, error) {
		return nil, emitOut{Data: strings.Repeat("x", in.Size)}, nil
	})
}

// changedWeatherBuild is defaultBuild but with get_weather using the incompatible weatherInV2 schema.
func changedWeatherBuild(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_weather",
		Description: "Return the current weather for a city.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: boolPtr(true)},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in weatherInV2) (*mcp.CallToolResult, weatherOut, error) {
		return nil, weatherOut{TempC: 21.5, Conditions: "sunny in " + in.Location}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_ticket",
		Description: "Create a ticket.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: boolPtr(false)},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in ticketIn) (*mcp.CallToolResult, ticketOut, error) {
		return nil, ticketOut{ID: "TCK-1:" + in.Title}, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "emit",
		Description: "Return a payload of the requested size.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in emitIn) (*mcp.CallToolResult, emitOut, error) {
		return nil, emitOut{Data: strings.Repeat("x", in.Size)}, nil
	})
}

// newMockMCP starts a mock server with the default tool set and registers cleanup.
func newMockMCP(t *testing.T) *mockMCP {
	t.Helper()
	m := &mockMCP{server: nil, mu: sync.Mutex{}, build: defaultBuild}

	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		srv := mcp.NewServer(&mcp.Implementation{Name: "mock", Title: "Mock MCP", Version: "0.0.1"}, nil)
		m.mu.Lock()
		build := m.build
		m.mu.Unlock()
		build(srv)
		return srv
	}, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true, DisableLocalhostProtection: true})

	m.server = httptest.NewServer(handler)
	t.Cleanup(m.server.Close)
	return m
}

// setBuild swaps the server's tool-registration closure. The next request observes the new schema.
func (m *mockMCP) setBuild(f func(s *mcp.Server)) {
	m.mu.Lock()
	m.build = f
	m.mu.Unlock()
}

// url is the endpoint the client should connect to.
func (m *mockMCP) url() string { return m.server.URL }
