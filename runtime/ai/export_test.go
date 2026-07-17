package ai

import "github.com/rilldata/rill/runtime/act/mcpconn"

// SetNewMCPClientForTest overrides the outbound MCP client factory and returns a function that restores the
// previous factory. It is the seam an external test uses to point a dynamic agent's connectors at an in-process
// server (granting loopback/http via mcpconn.Options.AllowPrivateNetworks), which production never does.
func SetNewMCPClientForTest(fn func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error)) func() {
	prev := newMCPClient
	newMCPClient = fn
	return func() { newMCPClient = prev }
}
